package cel2sql_test

import (
	"testing"

	"cel.dev/cel-go/cel"
	"cel.dev/cel-go/common/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/spandigital/cel2sql/v3"
)

func TestWithParamStartIndex_Default(t *testing.T) {
	env, err := cel.NewEnv(
		cel.CustomTypeAdapter(types.DefaultTypeAdapter),
		cel.Variable("name", cel.StringType),
		cel.Variable("age", cel.IntType),
	)
	require.NoError(t, err)

	ast, issues := env.Compile(`name == "Alice" && age > 30`)
	require.NoError(t, issues.Err())

	result, err := cel2sql.ConvertParameterized(ast)
	require.NoError(t, err)

	assert.Equal(t, `name = $1 AND age > $2`, result.SQL)
	assert.Equal(t, []any{"Alice", int64(30)}, result.Parameters)
}

func TestWithParamStartIndex_ExplicitOne(t *testing.T) {
	env, err := cel.NewEnv(
		cel.CustomTypeAdapter(types.DefaultTypeAdapter),
		cel.Variable("name", cel.StringType),
	)
	require.NoError(t, err)

	ast, issues := env.Compile(`name == "Alice"`)
	require.NoError(t, issues.Err())

	result, err := cel2sql.ConvertParameterized(ast, cel2sql.WithParamStartIndex(1))
	require.NoError(t, err)

	assert.Equal(t, `name = $1`, result.SQL, "WithParamStartIndex(1) is identical to default")
	assert.Equal(t, []any{"Alice"}, result.Parameters)
}

func TestWithParamStartIndex_Five(t *testing.T) {
	env, err := cel.NewEnv(
		cel.CustomTypeAdapter(types.DefaultTypeAdapter),
		cel.Variable("name", cel.StringType),
		cel.Variable("age", cel.IntType),
	)
	require.NoError(t, err)

	ast, issues := env.Compile(`name == "Alice" && age > 30`)
	require.NoError(t, issues.Err())

	result, err := cel2sql.ConvertParameterized(ast, cel2sql.WithParamStartIndex(5))
	require.NoError(t, err)

	assert.Equal(t, `name = $5 AND age > $6`, result.SQL,
		"placeholders should start at $5 and continue $6, $7, ...")
	assert.Equal(t, []any{"Alice", int64(30)}, result.Parameters,
		"Parameters slice contents and length are unchanged regardless of start index")
}

func TestWithParamStartIndex_ClampToOne(t *testing.T) {
	env, err := cel.NewEnv(
		cel.CustomTypeAdapter(types.DefaultTypeAdapter),
		cel.Variable("name", cel.StringType),
	)
	require.NoError(t, err)

	ast, issues := env.Compile(`name == "Alice"`)
	require.NoError(t, issues.Err())

	tests := []struct {
		name  string
		index int
	}{
		{"zero", 0},
		{"negative", -1},
		{"large negative", -100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := cel2sql.ConvertParameterized(ast, cel2sql.WithParamStartIndex(tt.index))
			require.NoError(t, err)

			assert.Equal(t, `name = $1`, result.SQL, "values < 1 should clamp to 1")
			assert.Equal(t, []any{"Alice"}, result.Parameters)
		})
	}
}

func TestWithParamStartIndex_ManyParams(t *testing.T) {
	env, err := cel.NewEnv(
		cel.CustomTypeAdapter(types.DefaultTypeAdapter),
		cel.Variable("a", cel.StringType),
		cel.Variable("b", cel.StringType),
		cel.Variable("c", cel.StringType),
		cel.Variable("d", cel.StringType),
	)
	require.NoError(t, err)

	ast, issues := env.Compile(`a == "1" && b == "2" && c == "3" && d == "4"`)
	require.NoError(t, issues.Err())

	result, err := cel2sql.ConvertParameterized(ast, cel2sql.WithParamStartIndex(10))
	require.NoError(t, err)

	assert.Equal(t, `a = $10 AND b = $11 AND c = $12 AND d = $13`, result.SQL)
	assert.Equal(t, []any{"1", "2", "3", "4"}, result.Parameters)
}
