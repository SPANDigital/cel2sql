// Copyright 2026 Observe, Inc. Licensed under Apache 2.0.
// Tests for list- and map-valued JSON variables, added by Observe, Inc.
package cel2sql_test

import (
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/observeinc/cel2sql/v3"
	dialectpkg "github.com/observeinc/cel2sql/v3/dialect"
	bigqueryDialect "github.com/observeinc/cel2sql/v3/dialect/bigquery"
	duckdbDialect "github.com/observeinc/cel2sql/v3/dialect/duckdb"
	mysqlDialect "github.com/observeinc/cel2sql/v3/dialect/mysql"
	postgresDialect "github.com/observeinc/cel2sql/v3/dialect/postgres"
	sqliteDialect "github.com/observeinc/cel2sql/v3/dialect/sqlite"
)

// newJSONContainerEnv declares tags with the object-tag shape, a map whose values
// are lists, beside the map<string, string> shape WithJSONVariables started with.
func newJSONContainerEnv(t *testing.T) *cel.Env {
	t.Helper()
	env, err := cel.NewEnv(
		cel.CustomTypeAdapter(types.DefaultTypeAdapter),
		cel.Variable("tags", cel.MapType(cel.StringType, cel.ListType(cel.StringType))),
		cel.Variable("counts", cel.MapType(cel.StringType, cel.ListType(cel.IntType))),
		cel.Variable("nested", cel.MapType(cel.StringType, cel.MapType(cel.StringType, cel.StringType))),
		cel.Variable("context", cel.MapType(cel.StringType, cel.StringType)),
		cel.Variable("name", cel.StringType),
	)
	require.NoError(t, err)
	return env
}

func convertJSONContainer(t *testing.T, expr string, opts ...cel2sql.ConvertOption) (*cel2sql.Result, error) {
	t.Helper()
	ast, issues := newJSONContainerEnv(t).Compile(expr)
	require.NoError(t, issues.Err())
	opts = append([]cel2sql.ConvertOption{cel2sql.WithJSONVariables("tags", "counts", "nested", "context")}, opts...)
	return cel2sql.ConvertParameterized(ast, opts...)
}

// Filters shared with the integration test.
const (
	exprEnvKeyPresent = `"env" in tags`
	exprProdInEnv     = `"prod" in tags.env`
)

// guardEnv is the PostgreSQL array guard the converter wraps a read of tags.env in.
func guardEnv(body string) string {
	return `CASE WHEN jsonb_typeof(tags->'env') = 'array' THEN ` + body + ` END`
}

func TestJSONVariableContainers_PostgreSQL(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		wantSQL  string
		wantArgs []any
	}{
		{
			name:    "key in map",
			expr:    exprEnvKeyPresent,
			wantSQL: `tags ? 'env'`,
		},
		{
			name:    "key in a map<string, string> variable",
			expr:    `"host" in context`,
			wantSQL: `context ? 'host'`,
		},
		{
			name:    "has on a key",
			expr:    `has(tags.env)`,
			wantSQL: `tags ? 'env'`,
		},
		{
			name:    "non-identifier key in map",
			expr:    `"k8s.cluster" in tags`,
			wantSQL: `tags ? 'k8s.cluster'`,
		},
		{
			name:     "value in list, bracket",
			expr:     `"prod" in tags["env"]`,
			wantSQL:  guardEnv(`(tags->'env') ? ($1)`),
			wantArgs: []any{"prod"},
		},
		{
			name:     "value in list, dot",
			expr:     exprProdInEnv,
			wantSQL:  guardEnv(`(tags->'env') ? ($1)`),
			wantArgs: []any{"prod"},
		},
		{
			name:    "column in list",
			expr:    `name in tags.env`,
			wantSQL: guardEnv(`(tags->'env') ? (name)`),
		},
		{
			name: "exists over a list value",
			expr: `tags.env.exists(v, v.startsWith("pr"))`,
			wantSQL: guardEnv(
				`EXISTS (SELECT 1 FROM jsonb_array_elements_text(tags->'env') AS v WHERE v LIKE 'pr%' ESCAPE E'\\')`),
		},
		{
			name: "all over a list value",
			expr: `tags["env"].all(v, v != "dev")`,
			wantSQL: guardEnv(
				`NOT EXISTS (SELECT 1 FROM jsonb_array_elements_text(tags->'env') AS v WHERE NOT (v != $1))`),
			wantArgs: []any{"dev"},
		},
		{
			name:    "exists over the map's keys",
			expr:    `tags.exists(k, k.matches("^team"))`,
			wantSQL: `EXISTS (SELECT 1 FROM jsonb_object_keys(tags) AS k WHERE k ~ '^team')`,
		},
		{
			name:     "size of the map",
			expr:     `size(tags) > 1`,
			wantSQL:  `(SELECT COUNT(*) FROM jsonb_object_keys(tags)) > $1`,
			wantArgs: []any{int64(1)},
		},
		{
			name:     "size of a list value",
			expr:     `tags.env.size() == 2`,
			wantSQL:  guardEnv(`COALESCE(jsonb_array_length(tags->'env'), 0)`) + ` = $1`,
			wantArgs: []any{int64(2)},
		},
		{
			name:     "combined with a column",
			expr:     `"prod" in tags.env && name == "checkout"`,
			wantSQL:  guardEnv(`(tags->'env') ? ($1)`) + ` AND name = $2`,
			wantArgs: []any{"prod", "checkout"},
		},
		{
			name:     "a scalar value is still read as text",
			expr:     `context.host == "web-1"`,
			wantSQL:  `context->>'host' = $1`,
			wantArgs: []any{"web-1"},
		},
		{
			name:     "size of a scalar value is still its length",
			expr:     `size(context.host) > 3`,
			wantSQL:  `LENGTH(context->>'host') > $1`,
			wantArgs: []any{int64(3)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertJSONContainer(t, tt.expr)
			require.NoError(t, err)
			assert.Equal(t, tt.wantSQL, result.SQL)
			if tt.wantArgs != nil {
				assert.Equal(t, tt.wantArgs, result.Parameters)
			}
		})
	}
}

// TestJSONVariableContainers_Refused pins the uses of a list or map value that
// have no faithful SQL form. Each used to convert to text-reading SQL that either
// did not run or answered the wrong question; now each is a conversion error.
func TestJSONVariableContainers_Refused(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		wantErr error
	}{
		{name: "comparing a list value", expr: `tags.env == ["prod"]`},
		{name: "indexing a list value", expr: `tags["env"][0] == "prod"`},
		{name: "a list value inside a larger expression", expr: `size(tags.env + ["x"]) > 1`},
		{name: "membership in a non-string list", expr: `1 in counts.ids`},
		{name: "exists over a non-string list", expr: `counts.ids.exists(n, n > 1)`},
		{name: "size of a map value", expr: `size(nested.labels) > 0`},
		{name: "a key that is not a literal", expr: `name in tags`, wantErr: cel2sql.ErrInvalidFieldName},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := convertJSONContainer(t, tt.expr)
			require.Error(t, err)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}

// TestJSONVariableContainers_OtherDialects pins that a dialect without
// dialect.JSONContainerWriter refuses the list and map forms rather than emit
// SQL nobody has run against it, and that its scalar reads are unchanged.
func TestJSONVariableContainers_OtherDialects(t *testing.T) {
	containerForms := []string{
		exprEnvKeyPresent,
		`has(tags.env)`,
		exprProdInEnv,
		`tags.env.exists(v, v == "prod")`,
		`tags.exists(k, k == "env")`,
		`size(tags) > 0`,
		`size(tags.env) > 0`,
	}

	for _, d := range []dialectpkg.Dialect{mysqlDialect.New(), sqliteDialect.New(), duckdbDialect.New(), bigqueryDialect.New()} {
		t.Run(string(d.Name()), func(t *testing.T) {
			_, ok := dialectpkg.GetJSONContainerWriter(d)
			require.False(t, ok, "a dialect that opts in needs its own integration test before this one changes")

			for _, expr := range containerForms {
				_, err := convertJSONContainer(t, expr, cel2sql.WithDialect(d))
				assert.ErrorIs(t, err, cel2sql.ErrUnsupportedDialectFeature, expr)
			}

			_, err := convertJSONContainer(t, `context.host == "web-1"`, cel2sql.WithDialect(d))
			assert.NoError(t, err, "a scalar read needs no container support")
		})
	}

	t.Run("postgresql opts in", func(t *testing.T) {
		_, ok := dialectpkg.GetJSONContainerWriter(postgresDialect.New())
		require.True(t, ok)
	})
}
