package cel2sql_test

import (
	"testing"

	"cel.dev/cel-go/cel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/spandigital/cel2sql/v3"
	"github.com/spandigital/cel2sql/v3/pg"
)

// has() detection is driven by the schema, not by the column's name. Before this
// was fixed, a column had to be named one of seven hardcoded values to be treated
// as JSON; anything else walked past the column boundary and emitted the table
// alias as the JSON document, which PostgreSQL rejects.
func TestHasUsesSchemaNotColumnName(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "integer"},
		{Name: "payload", Type: "jsonb", IsJSON: true, IsJSONB: true},
		{Name: "legacy_doc", Type: "json", IsJSON: true},
		{Name: "title", Type: "text"},
	})
	provider := pg.NewTypeProvider(map[string]pg.Schema{"rec": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("rec", cel.ObjectType("rec")),
	)
	require.NoError(t, err)

	tests := []struct {
		name    string
		celExpr string
		wantSQL string
	}{
		{
			name:    "jsonb column not on the old name list",
			celExpr: `has(rec.payload.active)`,
			wantSQL: `rec.payload ? 'active'`,
		},
		{
			name:    "json column uses the arrow form",
			celExpr: `has(rec.legacy_doc.active)`,
			wantSQL: `rec.legacy_doc->'active' IS NOT NULL`,
		},
		{
			name:    "nested path below a non-listed jsonb column",
			celExpr: `has(rec.payload.user.name)`,
			wantSQL: `jsonb_extract_path_text(rec.payload, 'user', 'name') IS NOT NULL`,
		},
		{
			name:    "non-JSON column falls through to a null check",
			celExpr: `has(rec.title)`,
			wantSQL: `rec.title IS NOT NULL`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.celExpr)
			require.NoError(t, issues.Err())

			sql, err := cel2sql.Convert(ast, cel2sql.WithSchemas(provider.GetSchemas()))
			require.NoError(t, err)
			assert.Equal(t, tt.wantSQL, sql)
		})
	}
}

// A column named like one of the seven formerly-hardcoded values gets no special
// treatment: with no schema declaring it JSON, it is an ordinary column.
func TestHasWithoutSchemaIgnoresColumnName(t *testing.T) {
	env, err := cel.NewEnv(
		cel.Variable("rec", cel.MapType(cel.StringType, cel.DynType)),
	)
	require.NoError(t, err)

	ast, issues := env.Compile(`has(rec.metadata)`)
	require.NoError(t, issues.Err())

	sql, err := cel2sql.Convert(ast)
	require.NoError(t, err)
	assert.Equal(t, `rec.metadata IS NOT NULL`, sql)
	assert.NotContains(t, sql, "?", "no schema means no JSON treatment")
}

// WithJSONVariables declares a bare column as JSONB. The name-based
// implementation never consulted it, so these never reached the ? operator.
func TestHasHonoursJSONVariables(t *testing.T) {
	env, err := cel.NewEnv(
		cel.Variable("tags", cel.MapType(cel.StringType, cel.DynType)),
	)
	require.NoError(t, err)

	ast, issues := env.Compile(`has(tags.colour)`)
	require.NoError(t, issues.Err())

	sql, err := cel2sql.Convert(ast, cel2sql.WithJSONVariables("tags"))
	require.NoError(t, err)
	assert.Equal(t, `tags ? 'colour'`, sql)
}
