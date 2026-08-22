package cel2sql_test

import (
	"testing"

	"cel.dev/cel-go/cel"
	"cel.dev/cel-go/common/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/spandigital/cel2sql/v3"
)

func TestWithJSONVariables_DotNotation(t *testing.T) {
	env, err := cel.NewEnv(
		cel.CustomTypeAdapter(types.DefaultTypeAdapter),
		cel.Variable("status", cel.StringType),
		cel.Variable("tags", cel.MapType(cel.StringType, cel.StringType)),
	)
	require.NoError(t, err)

	tests := []struct {
		name     string
		expr     string
		wantSQL  string
		wantArgs []any
	}{
		{
			name:     "equality",
			expr:     `tags.color == "blue"`,
			wantSQL:  `tags->>'color' = $1`,
			wantArgs: []any{"blue"},
		},
		{
			name:     "not equal",
			expr:     `tags.color != "red"`,
			wantSQL:  `tags->>'color' != $1`,
			wantArgs: []any{"red"},
		},
		{
			name:     "contains",
			expr:     `tags.color.contains("lu")`,
			wantSQL:  `POSITION($1 IN tags->>'color') > 0`,
			wantArgs: []any{"lu"},
		},
		{
			name:     "startsWith",
			expr:     `tags.color.startsWith("bl")`,
			wantSQL:  `tags->>'color' LIKE 'bl%' ESCAPE E'\\'`,
			wantArgs: nil,
		},
		{
			name:     "endsWith",
			expr:     `tags.color.endsWith("ue")`,
			wantSQL:  `tags->>'color' LIKE '%ue' ESCAPE E'\\'`,
			wantArgs: nil,
		},
		{
			name:     "combined with flat variable",
			expr:     `status == "ok" && tags.color == "blue"`,
			wantSQL:  `status = $1 AND tags->>'color' = $2`,
			wantArgs: []any{"ok", "blue"},
		},
		{
			name:     "multiple keys",
			expr:     `tags.color == "blue" && tags.shape == "circle"`,
			wantSQL:  `tags->>'color' = $1 AND tags->>'shape' = $2`,
			wantArgs: []any{"blue", "circle"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.expr)
			require.NoError(t, issues.Err())

			result, err := cel2sql.ConvertParameterized(ast, cel2sql.WithJSONVariables("tags"))
			require.NoError(t, err)

			assert.Equal(t, tt.wantSQL, result.SQL)
			if tt.wantArgs == nil {
				assert.Empty(t, result.Parameters, "expected no parameters")
			} else {
				assert.Equal(t, tt.wantArgs, result.Parameters)
			}
		})
	}
}

func TestWithJSONVariables_BracketNotation(t *testing.T) {
	env, err := cel.NewEnv(
		cel.CustomTypeAdapter(types.DefaultTypeAdapter),
		cel.Variable("metadata", cel.MapType(cel.StringType, cel.StringType)),
	)
	require.NoError(t, err)

	tests := []struct {
		name     string
		expr     string
		wantSQL  string
		wantArgs []any
	}{
		{
			name:     "equality",
			expr:     `metadata["key"] == "value"`,
			wantSQL:  `metadata->>'key' = $1`,
			wantArgs: []any{"value"},
		},
		{
			name:     "contains",
			expr:     `metadata["key"].contains("val")`,
			wantSQL:  `POSITION($1 IN metadata->>'key') > 0`,
			wantArgs: []any{"val"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.expr)
			require.NoError(t, issues.Err())

			result, err := cel2sql.ConvertParameterized(ast, cel2sql.WithJSONVariables("metadata"))
			require.NoError(t, err)

			assert.Equal(t, tt.wantSQL, result.SQL)
			if tt.wantArgs == nil {
				assert.Empty(t, result.Parameters, "expected no parameters")
			} else {
				assert.Equal(t, tt.wantArgs, result.Parameters)
			}
		})
	}
}

func TestWithJSONVariables_MultipleVariables(t *testing.T) {
	env, err := cel.NewEnv(
		cel.CustomTypeAdapter(types.DefaultTypeAdapter),
		cel.Variable("tags", cel.MapType(cel.StringType, cel.StringType)),
		cel.Variable("metadata", cel.MapType(cel.StringType, cel.StringType)),
	)
	require.NoError(t, err)

	ast, issues := env.Compile(`tags.color == "blue" && metadata.source == "api"`)
	require.NoError(t, issues.Err())

	result, err := cel2sql.ConvertParameterized(ast, cel2sql.WithJSONVariables("tags", "metadata"))
	require.NoError(t, err)

	assert.Equal(t, `tags->>'color' = $1 AND metadata->>'source' = $2`, result.SQL)
	assert.Equal(t, []any{"blue", "api"}, result.Parameters)
}

func TestWithJSONVariables_BackwardCompatible(t *testing.T) {
	env, err := cel.NewEnv(
		cel.CustomTypeAdapter(types.DefaultTypeAdapter),
		cel.Variable("data", cel.MapType(cel.StringType, cel.StringType)),
	)
	require.NoError(t, err)

	ast, issues := env.Compile(`data.key == "value"`)
	require.NoError(t, issues.Err())

	result, err := cel2sql.ConvertParameterized(ast)
	require.NoError(t, err)

	assert.Equal(t, `data.key = $1`, result.SQL, "without WithJSONVariables, should produce dot notation")
}

func TestWithJSONVariables_NonParameterized(t *testing.T) {
	env, err := cel.NewEnv(
		cel.CustomTypeAdapter(types.DefaultTypeAdapter),
		cel.Variable("props", cel.MapType(cel.StringType, cel.StringType)),
	)
	require.NoError(t, err)

	ast, issues := env.Compile(`props.key == "value"`)
	require.NoError(t, issues.Err())

	sql, err := cel2sql.Convert(ast, cel2sql.WithJSONVariables("props"))
	require.NoError(t, err)

	assert.Equal(t, `props->>'key' = 'value'`, sql)
}

func TestWithJSONVariables_NestedDotNotation(t *testing.T) {
	env, err := cel.NewEnv(
		cel.CustomTypeAdapter(types.DefaultTypeAdapter),
		cel.Variable("tags", cel.MapType(cel.StringType, cel.DynType)),
	)
	require.NoError(t, err)

	tests := []struct {
		name     string
		expr     string
		wantSQL  string
		wantArgs []any
	}{
		{
			name:     "two levels",
			expr:     `tags.corpus.section == "intro"`,
			wantSQL:  `tags->'corpus'->>'section' = $1`,
			wantArgs: []any{"intro"},
		},
		{
			name:     "three levels",
			expr:     `tags.corpus.section.subsection == "x"`,
			wantSQL:  `tags->'corpus'->'section'->>'subsection' = $1`,
			wantArgs: []any{"x"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.expr)
			require.NoError(t, issues.Err())

			result, err := cel2sql.ConvertParameterized(ast, cel2sql.WithJSONVariables("tags"))
			require.NoError(t, err)

			assert.Equal(t, tt.wantSQL, result.SQL)
			assert.Equal(t, tt.wantArgs, result.Parameters)
		})
	}
}

func TestWithJSONVariables_OnlyDeclaredVarsAffected(t *testing.T) {
	env, err := cel.NewEnv(
		cel.CustomTypeAdapter(types.DefaultTypeAdapter),
		cel.Variable("tags", cel.MapType(cel.StringType, cel.StringType)),
		cel.Variable("other", cel.MapType(cel.StringType, cel.StringType)),
	)
	require.NoError(t, err)

	ast, issues := env.Compile(`tags.color == "blue" && other.key == "val"`)
	require.NoError(t, issues.Err())

	result, err := cel2sql.ConvertParameterized(ast, cel2sql.WithJSONVariables("tags"))
	require.NoError(t, err)

	assert.Equal(t, `tags->>'color' = $1 AND other.key = $2`, result.SQL,
		"only 'tags' should use JSONB operators; 'other' should use dot notation")
}

func TestWithColumnAliases(t *testing.T) {
	env, err := cel.NewEnv(
		cel.CustomTypeAdapter(types.DefaultTypeAdapter),
		cel.Variable("name", cel.StringType),
		cel.Variable("active", cel.BoolType),
	)
	require.NoError(t, err)

	aliases := map[string]string{
		"name":   "tbl_name",
		"active": "tbl_active",
	}

	tests := []struct {
		name     string
		expr     string
		wantSQL  string
		wantArgs []any
	}{
		{
			name:     "simple equality",
			expr:     `name == "Alice"`,
			wantSQL:  `tbl_name = $1`,
			wantArgs: []any{"Alice"},
		},
		{
			name:    "boolean",
			expr:    `active == true`,
			wantSQL: `tbl_active IS TRUE`,
		},
		{
			name:     "string contains",
			expr:     `name.contains("li")`,
			wantSQL:  `POSITION($1 IN tbl_name) > 0`,
			wantArgs: []any{"li"},
		},
		{
			name:     "combined",
			expr:     `name == "Alice" && active == true`,
			wantSQL:  `tbl_name = $1 AND tbl_active IS TRUE`,
			wantArgs: []any{"Alice"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.expr)
			require.NoError(t, issues.Err())

			result, err := cel2sql.ConvertParameterized(ast, cel2sql.WithColumnAliases(aliases))
			require.NoError(t, err)

			assert.Equal(t, tt.wantSQL, result.SQL)
			if tt.wantArgs == nil {
				assert.Empty(t, result.Parameters, "expected no parameters")
			} else {
				assert.Equal(t, tt.wantArgs, result.Parameters)
			}
		})
	}
}

func TestWithColumnAliases_CombinedWithJSONVariables(t *testing.T) {
	env, err := cel.NewEnv(
		cel.CustomTypeAdapter(types.DefaultTypeAdapter),
		cel.Variable("status", cel.StringType),
		cel.Variable("tags", cel.MapType(cel.StringType, cel.StringType)),
	)
	require.NoError(t, err)

	ast, issues := env.Compile(`status == "ok" && tags.color == "blue"`)
	require.NoError(t, issues.Err())

	result, err := cel2sql.ConvertParameterized(ast,
		cel2sql.WithColumnAliases(map[string]string{
			"status": "tbl_status",
			"tags":   "tbl_tags",
		}),
		cel2sql.WithJSONVariables("tags"),
	)
	require.NoError(t, err)

	assert.Equal(t, `tbl_status = $1 AND tbl_tags->>'color' = $2`, result.SQL)
	assert.Equal(t, []any{"ok", "blue"}, result.Parameters)
}

func TestWithColumnAliases_BackwardCompatible(t *testing.T) {
	env, err := cel.NewEnv(
		cel.CustomTypeAdapter(types.DefaultTypeAdapter),
		cel.Variable("name", cel.StringType),
	)
	require.NoError(t, err)

	ast, issues := env.Compile(`name == "Alice"`)
	require.NoError(t, issues.Err())

	result, err := cel2sql.ConvertParameterized(ast)
	require.NoError(t, err)

	assert.Equal(t, `name = $1`, result.SQL, "without WithColumnAliases, should use original name")
}
