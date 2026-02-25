package cel2sql_test

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"runtime"
	"strings"
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/google/cel-go/cel"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcbigquery "github.com/testcontainers/testcontainers-go/modules/gcloud/bigquery"
	"google.golang.org/api/option"
	"google.golang.org/api/option/internaloption"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/spandigital/cel2sql/v3"
	bigqueryDialect "github.com/spandigital/cel2sql/v3/dialect/bigquery"
	"github.com/spandigital/cel2sql/v3/pg"
)

//go:embed testdata/bigquery_seed.yaml
var bigQuerySeedYAML []byte

const (
	bigQueryProjectID = "test-project"
	bigQueryDataset   = "testdataset"
)

// setupBigQueryContainer starts a BigQuery emulator container and returns a client.
// Returns nil container and client if the emulator cannot start (e.g., on arm64).
func setupBigQueryContainer(ctx context.Context, t *testing.T) (*tcbigquery.Container, *bigquery.Client) {
	t.Helper()

	container, err := tcbigquery.Run(ctx,
		"ghcr.io/goccy/bigquery-emulator:0.6.6",
		tcbigquery.WithProjectID(bigQueryProjectID),
		tcbigquery.WithDataYAML(bytes.NewReader(bigQuerySeedYAML)),
		testcontainers.WithImagePlatform("linux/amd64"),
	)
	if err != nil {
		// The BigQuery emulator only provides amd64 images. On arm64 (Apple Silicon),
		// it crashes under QEMU emulation due to Go runtime lfstack.push issues.
		if runtime.GOARCH == "arm64" || strings.Contains(err.Error(), "no image found") || strings.Contains(err.Error(), "container exited") {
			t.Skipf("Skipping BigQuery integration test: emulator not available on %s/%s: %v", runtime.GOOS, runtime.GOARCH, err)
		}
		t.Fatalf("Failed to start BigQuery emulator container: %v", err)
	}

	opts := []option.ClientOption{
		option.WithEndpoint(container.URI()),
		option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
		option.WithoutAuthentication(),
		internaloption.SkipDialSettingsValidation(),
	}

	client, err := bigquery.NewClient(ctx, container.ProjectID(), opts...)
	if err != nil {
		if termErr := container.Terminate(ctx); termErr != nil {
			t.Logf("failed to terminate container: %v", termErr)
		}
		t.Fatalf("Failed to create BigQuery client: %v", err)
	}

	return container, client
}

// bigQueryCount executes a count query and returns the result.
func bigQueryCount(ctx context.Context, t *testing.T, client *bigquery.Client, query string) int {
	t.Helper()

	q := client.Query(query)
	it, err := q.Read(ctx)
	require.NoError(t, err, "Failed to execute query: %s", query)

	var row []bigquery.Value
	err = it.Next(&row)
	require.NoError(t, err, "Failed to read query result: %s", query)
	require.Len(t, row, 1, "Expected exactly one column in COUNT(*) result")

	switch v := row[0].(type) {
	case int64:
		return int(v)
	case float64:
		return int(v)
	default:
		t.Fatalf("Unexpected type %T for COUNT(*) result: %v", row[0], row[0])
		return 0
	}
}

// TestBigQueryOperatorsIntegration validates operator conversions against a BigQuery emulator.
func TestBigQueryOperatorsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	container, client := setupBigQueryContainer(ctx, t)
	defer func() {
		if closeErr := client.Close(); closeErr != nil {
			t.Logf("failed to close client: %v", closeErr)
		}
		if err := container.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}()

	env, err := cel.NewEnv(
		cel.Variable("id", cel.IntType),
		cel.Variable("text_val", cel.StringType),
		cel.Variable("int_val", cel.IntType),
		cel.Variable("float_val", cel.DoubleType),
		cel.Variable("bool_val", cel.BoolType),
		cel.Variable("nullable_text", cel.StringType),
		cel.Variable("nullable_int", cel.IntType),
	)
	require.NoError(t, err)

	dialectOpt := cel2sql.WithDialect(bigqueryDialect.New())

	tests := []struct {
		name         string
		celExpr      string
		expectedRows int
		description  string
	}{
		// Comparison operators
		{
			name:         "Equality string",
			celExpr:      `text_val == "hello"`,
			expectedRows: 1,
			description:  "String equality comparison",
		},
		{
			name:         "Equality integer",
			celExpr:      `int_val == 20`,
			expectedRows: 1,
			description:  "Integer equality comparison",
		},
		{
			name:         "Equality float",
			celExpr:      `float_val == 10.5`,
			expectedRows: 1,
			description:  "Float equality comparison",
		},
		{
			name:         "Equality boolean",
			celExpr:      `bool_val == true`,
			expectedRows: 3,
			description:  "Boolean equality comparison",
		},
		{
			name:         "Not equal",
			celExpr:      `text_val != "hello"`,
			expectedRows: 4,
			description:  "Not equal comparison",
		},
		{
			name:         "Less than",
			celExpr:      `int_val < 15`,
			expectedRows: 2, // 10, 5
			description:  "Less than comparison",
		},
		{
			name:         "Less than or equal",
			celExpr:      `int_val <= 15`,
			expectedRows: 3, // 10, 5, 15
			description:  "Less than or equal comparison",
		},
		{
			name:         "Greater than",
			celExpr:      `int_val > 15`,
			expectedRows: 2, // 20, 30
			description:  "Greater than comparison",
		},
		{
			name:         "Greater than or equal",
			celExpr:      `int_val >= 15`,
			expectedRows: 3, // 20, 30, 15
			description:  "Greater than or equal comparison",
		},

		// Logical operators
		{
			name:         "Logical AND",
			celExpr:      `int_val > 10 && bool_val == true`,
			expectedRows: 2, // rows 3 (30,true) and 5 (15,true)
			description:  "Logical AND operator",
		},
		{
			name:         "Logical OR",
			celExpr:      `int_val < 10 || bool_val == false`,
			expectedRows: 2, // rows 2 (20,false) and 4 (5,false)
			description:  "Logical OR operator",
		},
		{
			name:         "Logical NOT",
			celExpr:      `!bool_val`,
			expectedRows: 2, // rows 2 and 4
			description:  "Logical NOT operator",
		},
		{
			name:         "Complex logical expression",
			celExpr:      `(int_val > 10 && bool_val) || int_val < 10`,
			expectedRows: 3, // rows 3, 5, 4
			description:  "Complex nested logical operators",
		},

		// Arithmetic operators
		{
			name:         "Addition",
			celExpr:      `int_val + 10 == 20`,
			expectedRows: 1, // 10 + 10 = 20
			description:  "Addition operator",
		},
		{
			name:         "Subtraction",
			celExpr:      `int_val - 5 == 15`,
			expectedRows: 1, // 20 - 5 = 15
			description:  "Subtraction operator",
		},
		{
			name:         "Multiplication",
			celExpr:      `int_val * 2 == 20`,
			expectedRows: 1, // 10 * 2 = 20
			description:  "Multiplication operator",
		},
		{
			name:         "Division",
			celExpr:      `int_val / 2 == 10`,
			expectedRows: 1, // 20 / 2 = 10
			description:  "Division operator",
		},
		{
			name:         "Modulo",
			celExpr:      `int_val % 10 == 0`,
			expectedRows: 3, // 10, 20, 30
			description:  "Modulo operator",
		},
		{
			name:         "Complex arithmetic",
			celExpr:      `(int_val * 2) + 5 > 30`,
			expectedRows: 3, // (20*2)+5=45, (30*2)+5=65, (15*2)+5=35
			description:  "Complex arithmetic expression",
		},

		// String operators
		{
			name:         "String concatenation",
			celExpr:      `text_val + "!" == "hello!"`,
			expectedRows: 1,
			description:  "String concatenation (||)",
		},
		{
			name:         "String contains",
			celExpr:      `text_val.contains("world")`,
			expectedRows: 2, // "world", "hello world"
			description:  "String contains function (STRPOS)",
		},
		{
			name:         "String startsWith",
			celExpr:      `text_val.startsWith("hello")`,
			expectedRows: 2, // "hello", "hello world"
			description:  "String startsWith function (LIKE)",
		},
		{
			name:         "String endsWith",
			celExpr:      `text_val.endsWith("world")`,
			expectedRows: 2, // "world", "hello world"
			description:  "String endsWith function (LIKE)",
		},

		// Regex (BigQuery uses REGEXP_CONTAINS with RE2)
		{
			name:         "Regex match",
			celExpr:      `text_val.matches(r"^hello")`,
			expectedRows: 2, // "hello", "hello world"
			description:  "Regex match (REGEXP_CONTAINS)",
		},
		{
			name:         "Regex simple pattern",
			celExpr:      `text_val.matches(r"test")`,
			expectedRows: 2, // "test", "testing"
			description:  "Regex simple pattern",
		},

		// Complex combined operators
		{
			name:         "Complex multi-operator expression",
			celExpr:      `int_val > 10 && bool_val && text_val.contains("test")`,
			expectedRows: 2, // rows 3 and 5
			description:  "Complex expression with multiple operator types",
		},
		{
			name:         "Nested parenthesized operators",
			celExpr:      `((int_val + 5) * 2 > 30) && (text_val.contains("test") || bool_val)`,
			expectedRows: 2, // rows 3 and 5
			description:  "Deeply nested operators with parentheses",
		},
		{
			name:         "Triple negation",
			celExpr:      `!!!bool_val`,
			expectedRows: 2,
			description:  "Multiple NOT operators",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.celExpr)
			if issues != nil && issues.Err() != nil {
				t.Fatalf("CEL compilation failed: %v", issues.Err())
			}

			sqlCondition, err := cel2sql.Convert(ast, dialectOpt)
			require.NoError(t, err, "Conversion should succeed for: %s", tt.description)

			t.Logf("CEL Expression: %s", tt.celExpr)
			t.Logf("Generated SQL WHERE clause: %s", sqlCondition)

			// #nosec G202 - This is a test validating SQL generation, not a security risk
			query := fmt.Sprintf("SELECT COUNT(*) FROM `%s.test_data` WHERE %s", bigQueryDataset, sqlCondition)
			t.Logf("Full SQL Query: %s", query)

			actualRows := bigQueryCount(ctx, t, client, query)

			require.Equal(t, tt.expectedRows, actualRows,
				"Query should return expected number of rows. %s\nCEL: %s\nSQL: %s",
				tt.description, tt.celExpr, sqlCondition)

			t.Logf("OK: %s (expected %d rows, got %d rows)",
				tt.description, tt.expectedRows, actualRows)
		})
	}
}

// TestBigQueryJSONIntegration validates JSON operations against a BigQuery emulator.
func TestBigQueryJSONIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	container, client := setupBigQueryContainer(ctx, t)
	defer func() {
		if closeErr := client.Close(); closeErr != nil {
			t.Logf("failed to close client: %v", closeErr)
		}
		if err := container.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}()

	// Set up CEL environment with schema for JSON detection
	productSchema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "integer"},
		{Name: "name", Type: "text"},
		{Name: "price", Type: "double precision"},
		{Name: "metadata", Type: "json", IsJSON: true},
	})

	schemas := map[string]pg.Schema{
		"product": productSchema,
	}

	provider := pg.NewTypeProvider(schemas)

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("product", cel.ObjectType("product")),
	)
	require.NoError(t, err)

	dialectOpt := cel2sql.WithDialect(bigqueryDialect.New())
	schemaOpt := cel2sql.WithSchemas(schemas)

	tests := []struct {
		name         string
		celExpr      string
		expectedRows int
		description  string
	}{
		{
			name:         "JSON field access",
			celExpr:      `product.metadata.brand == "Acme"`,
			expectedRows: 2,
			description:  "JSON field access with JSON_VALUE",
		},
		{
			name:         "JSON field access different value",
			celExpr:      `product.metadata.color == "blue"`,
			expectedRows: 1,
			description:  "JSON field access with different value",
		},
		{
			name:         "JSON with regular field",
			celExpr:      `product.metadata.brand == "Acme" && product.price > 30.0`,
			expectedRows: 1, // Doohickey (Acme, 39.99)
			description:  "JSON field combined with regular field comparison",
		},
		{
			name:         "JSON field existence",
			celExpr:      `has(product.metadata.brand)`,
			expectedRows: 3, // All rows have 'brand'
			description:  "JSON field existence check",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.celExpr)
			if issues != nil && issues.Err() != nil {
				t.Fatalf("CEL compilation failed: %v", issues.Err())
			}

			sqlCondition, err := cel2sql.Convert(ast, dialectOpt, schemaOpt)
			require.NoError(t, err, "Conversion should succeed for: %s", tt.description)

			t.Logf("CEL Expression: %s", tt.celExpr)
			t.Logf("Generated SQL WHERE clause: %s", sqlCondition)

			// #nosec G202 - This is a test validating SQL generation, not a security risk
			query := fmt.Sprintf("SELECT COUNT(*) FROM `%s.products` product WHERE %s", bigQueryDataset, sqlCondition)
			t.Logf("Full SQL Query: %s", query)

			actualRows := bigQueryCount(ctx, t, client, query)

			require.Equal(t, tt.expectedRows, actualRows,
				"Query should return expected number of rows. %s\nCEL: %s\nSQL: %s",
				tt.description, tt.celExpr, sqlCondition)

			t.Logf("OK: %s", tt.description)
		})
	}
}
