package cel2sql_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"cel.dev/cel-go/cel"
	"cel.dev/cel-go/common/types"
	"github.com/opensearch-project/opensearch-go/v4"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tcopensearch "github.com/testcontainers/testcontainers-go/modules/opensearch"

	"github.com/spandigital/cel2sql/v3"
	opensearchDialect "github.com/spandigital/cel2sql/v3/dialect/opensearch"
	osprovider "github.com/spandigital/cel2sql/v3/opensearch"
)

const opensearchImage = "opensearchproject/opensearch:3.5.0"

const productsIndexBody = `{
  "mappings": {
    "properties": {
      "name":        {"type": "keyword"},
      "description": {"type": "text"},
      "price":       {"type": "double"},
      "stock":       {"type": "long"},
      "active":      {"type": "boolean"},
      "created_at":  {"type": "date"},
      "dims": {
        "properties": {
          "width":  {"type": "double"},
          "height": {"type": "double"}
        }
      },
      "variants": {
        "type": "nested",
        "properties": {
          "sku": {"type": "keyword"},
          "qty": {"type": "long"}
        }
      }
    }
  }
}`

var productDocs = []string{
	`{"name":"Widget","description":"a small blue widget","price":9.5,"stock":3,"active":true,"created_at":"2024-01-15T10:00:00Z","dims":{"width":4,"height":2},"variants":[{"sku":"W-1","qty":2},{"sku":"W-2","qty":0}]}`,
	`{"name":"Gadget","description":"a large red gadget","price":25,"stock":0,"active":true,"created_at":"2024-03-10T08:30:00Z","dims":{"width":12,"height":6},"variants":[{"sku":"G-1","qty":5}]}`,
	`{"name":"Gizmo","description":"blue gizmo, discontinued","price":40,"stock":1,"active":false,"created_at":"2024-07-01T00:00:00Z","dims":{"width":8},"variants":[{"sku":"A-1","qty":1},{"sku":"A-2","qty":3}]}`,
	`{"name":"Doohickey","description":"plain doohickey","price":3,"stock":99,"active":true,"created_at":"2023-11-20T16:45:00Z","dims":{"width":1,"height":1},"variants":[]}`,
}

type sqlResponse struct {
	Datarows [][]any `json:"datarows"`
	Error    *struct {
		Reason  string `json:"reason"`
		Details string `json:"details"`
		Type    string `json:"type"`
	} `json:"error"`
	Status int `json:"status"`
}

// runSQL posts a query to the SQL plugin and returns the rows.
func runSQL(t *testing.T, addr, query string, params []map[string]any) ([][]any, error) {
	t.Helper()
	body := map[string]any{"query": query}
	if params != nil {
		body["parameters"] = params
	}
	payload, err := json.Marshal(body)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, addr+"/_plugins/_sql", bytes.NewReader(payload))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	raw, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var parsed sqlResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, fmt.Errorf("unparseable response %q: %w", raw, err)
	}
	if parsed.Error != nil {
		return nil, fmt.Errorf("%s: %s (%s)", parsed.Error.Type, parsed.Error.Reason, parsed.Error.Details)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, raw)
	}
	return parsed.Datarows, nil
}

func names(rows [][]any) []string {
	out := make([]string, 0, len(rows))
	for _, r := range rows {
		out = append(out, fmt.Sprint(r[0]))
	}
	return out
}

// TestOpenSearchIntegration runs generated SQL against a real OpenSearch node.
func TestOpenSearchIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	container, err := tcopensearch.Run(ctx, opensearchImage)
	if err != nil {
		t.Skipf("Skipping OpenSearch integration test: container unavailable (is Docker reachable?): %v", err)
	}
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}()

	addr, err := container.Address(ctx)
	require.NoError(t, err)

	client, err := opensearchapi.NewClient(opensearchapi.Config{
		Client: opensearch.Config{Addresses: []string{addr}},
	})
	require.NoError(t, err)

	_, err = client.Indices.Create(ctx, opensearchapi.IndicesCreateReq{
		Index: "products",
		Body:  strings.NewReader(productsIndexBody),
	})
	require.NoError(t, err)

	for i, doc := range productDocs {
		_, err = client.Document.Create(ctx, opensearchapi.DocumentCreateReq{
			Index:      "products",
			DocumentID: strconv.Itoa(i + 1),
			Body:       strings.NewReader(doc),
			Params:     opensearchapi.DocumentCreateParams{Refresh: "true"},
		})
		require.NoError(t, err)
	}

	// The SQL plugin reads from the index's search shards; give the refresh a moment.
	time.Sleep(500 * time.Millisecond)

	// --- Provider: load the mapping from the live cluster ---
	provider, err := osprovider.NewTypeProviderWithClient(ctx, client)
	require.NoError(t, err)
	require.NoError(t, provider.LoadTableSchema(ctx, "products"))

	ft, found := provider.FindStructFieldType("products", "dims")
	require.True(t, found)
	assert.True(t, types.NewObjectType("products.dims").IsExactType(ft.Type))
	ft, found = provider.FindStructFieldType("products", "variants")
	require.True(t, found)
	assert.True(t, types.NewListType(types.NewObjectType("products.variants")).IsExactType(ft.Type))

	require.ErrorIs(t, provider.LoadTableSchema(ctx, "does-not-exist"), osprovider.ErrInvalidSchema)

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("p", cel.ObjectType("products")),
	)
	require.NoError(t, err)

	// --- Generated SQL executes and returns the expected documents ---
	tests := []struct {
		name string
		cel  string
		want []string
	}{
		{"numeric and boolean", `p.price > 10.0 && p.active`, []string{"Gadget"}},
		{"bare boolean", `p.active`, []string{"Doohickey", "Gadget", "Widget"}},
		{"not active", `!p.active`, []string{"Gizmo"}},
		{"boolean equality", `p.active == false`, []string{"Gizmo"}},
		{"inequality", `p.name != "Widget" && p.stock > 0`, []string{"Doohickey", "Gizmo"}},
		{"ternary on bare boolean", `(p.active ? "y" : "n") == "n"`, []string{"Gizmo"}},
		{"startsWith", `p.name.startsWith("Wid")`, []string{"Widget"}},
		{"endsWith", `p.name.endsWith("et")`, []string{"Gadget", "Widget"}},
		{"in list", `p.name in ["Widget", "Gadget"]`, []string{"Gadget", "Widget"}},
		{"regex", `p.name.matches("^G.*o$")`, []string{"Gizmo"}},
		{"regex case-insensitive", `p.name.matches("(?i)^gadget$")`, []string{"Gadget"}},
		{"contains", `p.name.contains("dg")`, []string{"Gadget", "Widget"}},
		{"ternary", `(p.stock > 0 ? "in" : "out") == "out"`, []string{"Gadget"}},
		{"year extract", `p.created_at.getFullYear() == 2023`, []string{"Doohickey"}},
		{"month extract", `p.created_at.getMonth() == 2`, []string{"Gadget"}},
		{"timestamp literal compare", `p.created_at < timestamp("2024-02-01T00:00:00Z")`, []string{"Doohickey", "Widget"}},
		{"timestamp arithmetic", `p.created_at + duration("24h") > timestamp("2024-06-30T12:00:00Z")`, []string{"Gizmo"}},
		{"object dot path", `p.dims.width > 5.0`, []string{"Gadget", "Gizmo"}},
		{"has on object field", `has(p.dims.height)`, []string{"Doohickey", "Gadget", "Widget"}},
		{"nested exists", `p.variants.exists(v, v.sku == "A-1")`, []string{"Gizmo"}},
		{"nested exists with predicate", `p.variants.exists(v, v.qty > 2)`, []string{"Gadget", "Gizmo"}},
		{"string size", `size(p.name) == 5`, []string{"Gizmo"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ast, issues := env.Compile(tc.cel)
			require.NoError(t, issues.Err())
			where, err := cel2sql.Convert(ast, cel2sql.WithDialect(opensearchDialect.New()))
			require.NoError(t, err)

			query := "SELECT name FROM products AS p WHERE " + where + " ORDER BY name"
			rows, err := runSQL(t, addr, query, nil)
			require.NoError(t, err, "query: %s", query)
			assert.Equal(t, tc.want, names(rows), "query: %s", query)
		})
	}

	// --- Forms the engine cannot run are rejected before reaching it ---
	t.Run("all is rejected at conversion time", func(t *testing.T) {
		ast, issues := env.Compile(`p.variants.all(v, v.qty > 0)`)
		require.NoError(t, issues.Err())
		_, err := cel2sql.Convert(ast, cel2sql.WithDialect(opensearchDialect.New()))
		require.Error(t, err)
	})

	// --- Parameterized SQL through the request's parameters array ---
	t.Run("parameterized", func(t *testing.T) {
		ast, issues := env.Compile(`p.name == "Widget" && p.stock > 1`)
		require.NoError(t, issues.Err())
		res, err := cel2sql.ConvertParameterized(ast, cel2sql.WithDialect(opensearchDialect.New()))
		require.NoError(t, err)
		assert.Equal(t, "p.name = ? AND p.stock > ?", res.SQL)

		params := []map[string]any{
			{"type": "string", "value": res.Parameters[0]},
			{"type": "long", "value": res.Parameters[1]},
		}
		rows, err := runSQL(t, addr, "SELECT name FROM products AS p WHERE "+res.SQL, params)
		require.NoError(t, err)
		assert.Equal(t, []string{"Widget"}, names(rows))
	})
}
