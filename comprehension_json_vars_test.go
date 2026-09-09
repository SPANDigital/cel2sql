package cel2sql_test

import (
	"testing"

	"cel.dev/cel-go/cel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/spandigital/cel2sql/v3"
	"github.com/spandigital/cel2sql/v3/pg"
)

// Whether a comprehension variable accesses JSON is decided by what it ranges
// over, not by what it is called. The converter used to key this off a list of
// variable names — attr, item, element, obj, feature, review — so a jsonb array
// iterated as `row` was treated as a composite, and a composite array iterated
// as `item` was treated as JSON.
func TestComprehensionVarJSONnessComesFromTheRange(t *testing.T) {
	employee := pg.NewSchema([]pg.FieldSchema{
		{Name: "name", Type: "text"},
		{Name: "salary", Type: "bigint"},
	})
	data := pg.NewSchema([]pg.FieldSchema{
		{Name: "docs", Type: "jsonb", Repeated: true},
		{Name: "items", Type: "jsonb"},
		{Name: "scores", Type: "bigint", Repeated: true},
	})

	provider := pg.NewTypeProvider(map[string]pg.Schema{
		"Employee": employee,
		"data":     data,
	})
	schemas := map[string]pg.Schema{"data": data}

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("data", cel.ObjectType("data")),
		cel.Variable("staff", cel.ListType(cel.ObjectType("Employee"))),
	)
	require.NoError(t, err)

	tests := []struct {
		name    string
		celExpr string
		wantSQL string
	}{
		{
			// jsonb[] unnests to documents, so field access extracts from them —
			// and "row" is not a name the old list knew about.
			name:    "json array with an unlisted variable name",
			celExpr: `data.docs.exists(row, row.status == "open")`,
			wantSQL: `EXISTS (SELECT 1 FROM UNNEST(data.docs) AS row WHERE row->>'status' = 'open')`,
		},
		{
			// A composite array unnests to rows, so the same access names a column —
			// even though "item" was on the old list.
			name:    "composite array with a formerly listed variable name",
			celExpr: `staff.exists(item, item.name == "admin")`,
			wantSQL: `EXISTS (SELECT 1 FROM UNNEST(staff) AS item WHERE item.name = 'admin')`,
		},
		{
			name:    "scalar array is untouched",
			celExpr: `data.scores.exists(x, x > 10)`,
			wantSQL: `EXISTS (SELECT 1 FROM UNNEST(data.scores) AS x WHERE x > 10)`,
		},
		{
			// Numeric casting comes from the compared type, not from the field name.
			name:    "numeric comparison casts without a name list",
			celExpr: `data.docs.exists(row, row.score > 100)`,
			wantSQL: `EXISTS (SELECT 1 FROM UNNEST(data.docs) AS row WHERE (row->>'score')::numeric > 100)`,
		},
		{
			// The documented example from docs/json-support.md. It produced this
			// same SQL before, but only because the variable was called "item";
			// now it is because items is declared jsonb. (Which unnest function
			// suits a jsonb array value is a separate, pre-existing question.)
			name:    "jsonb array value",
			celExpr: `data.items.filter(item, item.price > 10)`,
			wantSQL: `ARRAY(SELECT item FROM UNNEST(data.items) AS item WHERE (item->>'price')::numeric > 10)`,
		},
		{
			// Nested paths route through the JSON path builder rather than
			// extracting the first segment as text and then dotting into it.
			name:    "nested access on a JSON element",
			celExpr: `data.docs.exists(row, row.meta.active == "yes")`,
			wantSQL: `EXISTS (SELECT 1 FROM UNNEST(data.docs) AS row WHERE row->'meta'->>'active' = 'yes')`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.celExpr)
			require.NoError(t, issues.Err())

			sql, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
			require.NoError(t, err)
			assert.Equal(t, tt.wantSQL, sql)
		})
	}
}

// Outside a comprehension a variable is just a variable, whatever it is called.
func TestFormerlyListedVariableNamesGetNoSpecialTreatment(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{{Name: "name", Type: "text"}})
	schemas := map[string]pg.Schema{}

	for _, varName := range []string{"attr", "item", "element", "obj", "feature", "review"} {
		t.Run(varName, func(t *testing.T) {
			schemas[varName] = schema
			defer delete(schemas, varName)

			env, err := cel.NewEnv(
				cel.CustomTypeProvider(pg.NewTypeProvider(map[string]pg.Schema{varName: schema})),
				cel.Variable(varName, cel.ObjectType(varName)),
			)
			require.NoError(t, err)

			ast, issues := env.Compile(varName + `.name == "x"`)
			require.NoError(t, issues.Err())

			sql, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
			require.NoError(t, err)
			assert.Equal(t, varName+`.name = 'x'`, sql)
		})
	}
}

// The binding is scoped: an inner comprehension may reuse a name the outer one
// bound to a different kind of collection, and the outer meaning must survive.
func TestNestedComprehensionsRestoreOuterBinding(t *testing.T) {
	data := pg.NewSchema([]pg.FieldSchema{
		{Name: "docs", Type: "jsonb", Repeated: true},
		{Name: "scores", Type: "bigint", Repeated: true},
	})
	schemas := map[string]pg.Schema{"data": data}

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(pg.NewTypeProvider(schemas)),
		cel.Variable("data", cel.ObjectType("data")),
	)
	require.NoError(t, err)

	// The inner comprehension rebinds "x" to a scalar; the outer "x" is a JSON
	// document and must still extract after the inner one closes.
	ast, issues := env.Compile(`data.docs.exists(x, data.scores.exists(x, x > 1) && x.status == "open")`)
	require.NoError(t, issues.Err())

	sql, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
	require.NoError(t, err)

	assert.Contains(t, sql, `x->>'status' = 'open'`, "outer JSON binding restored after the inner one")
	assert.Contains(t, sql, `AS x WHERE x > 1`, "inner scalar binding not treated as JSON")
}
