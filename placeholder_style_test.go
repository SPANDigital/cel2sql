package cel2sql_test

import (
	"strings"
	"testing"

	"cel.dev/cel-go/cel"
	"cel.dev/cel-go/common/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/spandigital/cel2sql/v3"
	"github.com/spandigital/cel2sql/v3/dialect"
	"github.com/spandigital/cel2sql/v3/dialect/bigquery"
	"github.com/spandigital/cel2sql/v3/dialect/duckdb"
	"github.com/spandigital/cel2sql/v3/dialect/mysql"
	"github.com/spandigital/cel2sql/v3/dialect/postgres"
	"github.com/spandigital/cel2sql/v3/dialect/spark"
	"github.com/spandigital/cel2sql/v3/dialect/sqlite"
	"github.com/spandigital/cel2sql/v3/pg"
	"github.com/spandigital/cel2sql/v3/schema"
)

func placeholderTestEnv(t *testing.T) *cel.Env {
	t.Helper()
	env, err := cel.NewEnv(
		cel.CustomTypeAdapter(types.DefaultTypeAdapter),
		cel.Variable("name", cel.StringType),
		cel.Variable("age", cel.IntType),
	)
	require.NoError(t, err)
	return env
}

func placeholderTestAST(t *testing.T) *cel.Ast {
	t.Helper()
	ast, issues := placeholderTestEnv(t).Compile(`name == "Alice" && age > 30`)
	require.NoError(t, issues.Err())
	return ast
}

var placeholderDialects = []struct {
	name        dialect.Name
	dialect     dialect.Dialect
	wantDefault string
}{
	{dialect.PostgreSQL, postgres.New(), `name = $1 AND age > $2`},
	{dialect.DuckDB, duckdb.New(), `name = $1 AND age > $2`},
	{dialect.BigQuery, bigquery.New(), `name = @p1 AND age > @p2`},
	{dialect.MySQL, mysql.New(), `name = ? AND age > ?`},
	{dialect.SQLite, sqlite.New(), `name = ? AND age > ?`},
	{dialect.Spark, spark.New(), `name = ? AND age > ?`},
}

// The default must stay byte-identical for every dialect — this is the
// regression guard for threading the style through visitConst.
func TestPlaceholderStyle_DefaultUnchanged(t *testing.T) {
	ast := placeholderTestAST(t)

	for _, d := range placeholderDialects {
		t.Run(string(d.name), func(t *testing.T) {
			result, err := cel2sql.ConvertParameterized(ast, cel2sql.WithDialect(d.dialect))
			require.NoError(t, err)

			assert.Equal(t, d.wantDefault, result.SQL)
			assert.Equal(t, []any{"Alice", int64(30)}, result.Parameters)

			explicit, err := cel2sql.ConvertParameterized(ast,
				cel2sql.WithDialect(d.dialect),
				cel2sql.WithPlaceholderStyle(cel2sql.PlaceholderDialect))
			require.NoError(t, err)
			assert.Equal(t, result.SQL, explicit.SQL, "PlaceholderDialect is the default")
		})
	}
}

func TestPlaceholderStyle_Question(t *testing.T) {
	ast := placeholderTestAST(t)

	for _, d := range placeholderDialects {
		t.Run(string(d.name), func(t *testing.T) {
			result, err := cel2sql.ConvertParameterized(ast,
				cel2sql.WithDialect(d.dialect),
				cel2sql.WithPlaceholderStyle(cel2sql.PlaceholderQuestion))
			require.NoError(t, err)

			assert.Equal(t, `name = ? AND age > ?`, result.SQL)
			assert.Equal(t, []any{"Alice", int64(30)}, result.Parameters,
				"parameters are unaffected by placeholder style")
			assert.Equal(t, len(result.Parameters), strings.Count(result.SQL, "?"))
		})
	}
}

// MySQL, SQLite and Spark already emit ? natively, so the option is a no-op there.
func TestPlaceholderStyle_NoOpForQuestionMarkDialects(t *testing.T) {
	ast := placeholderTestAST(t)

	for _, d := range []dialect.Dialect{mysql.New(), sqlite.New(), spark.New()} {
		t.Run(string(d.Name()), func(t *testing.T) {
			def, err := cel2sql.ConvertParameterized(ast, cel2sql.WithDialect(d))
			require.NoError(t, err)

			question, err := cel2sql.ConvertParameterized(ast,
				cel2sql.WithDialect(d),
				cel2sql.WithPlaceholderStyle(cel2sql.PlaceholderQuestion))
			require.NoError(t, err)

			assert.Equal(t, def.SQL, question.SQL)
			assert.Equal(t, def.Parameters, question.Parameters)
		})
	}
}

// The two options address opposite situations, so combining them is not an error —
// there is simply no index left to offset.
func TestPlaceholderStyle_QuestionIgnoresParamStartIndex(t *testing.T) {
	ast := placeholderTestAST(t)

	result, err := cel2sql.ConvertParameterized(ast,
		cel2sql.WithParamStartIndex(5),
		cel2sql.WithPlaceholderStyle(cel2sql.PlaceholderQuestion))
	require.NoError(t, err)

	assert.Equal(t, `name = ? AND age > ?`, result.SQL)
	assert.Equal(t, []any{"Alice", int64(30)}, result.Parameters)
}

// This is why the option exists rather than leaving callers to rewrite $n to ?:
// matches() inlines its pattern as a string literal, so a $1 in the pattern is
// data, not a placeholder, and only the converter can tell them apart.
func TestPlaceholderStyle_RegexPatternKeepsLiteralDollar(t *testing.T) {
	env, err := cel.NewEnv(
		cel.CustomTypeAdapter(types.DefaultTypeAdapter),
		cel.Variable("name", cel.StringType),
	)
	require.NoError(t, err)

	ast, issues := env.Compile(`name.matches("a$1b") && name == "x"`)
	require.NoError(t, issues.Err())

	result, err := cel2sql.ConvertParameterized(ast,
		cel2sql.WithPlaceholderStyle(cel2sql.PlaceholderQuestion))
	require.NoError(t, err)

	assert.Contains(t, result.SQL, `a$1b`, "the regex pattern keeps its literal $1")
	assert.Equal(t, 1, strings.Count(result.SQL, "?"), "only the bound literal is a placeholder")
	assert.Equal(t, []any{"x"}, result.Parameters)
}

func TestPlaceholderStyle_ConvertUnaffected(t *testing.T) {
	ast := placeholderTestAST(t)

	def, err := cel2sql.Convert(ast)
	require.NoError(t, err)

	styled, err := cel2sql.Convert(ast, cel2sql.WithPlaceholderStyle(cel2sql.PlaceholderQuestion))
	require.NoError(t, err)

	assert.Equal(t, def, styled, "Convert inlines literals and has no placeholders")
	assert.NotContains(t, styled, "?")
}

// PostgreSQL spells jsonb existence as the ? operator, which is indistinguishable
// from a placeholder to a consumer that scans for ?. Converting must fail rather
// than hand back SQL that would bind a value to an operator.
func TestPlaceholderStyle_QuestionRejectsPostgresJSONBExistence(t *testing.T) {
	// A deliberately unremarkable column name: the ? operator form is reached
	// because the schema says jsonb, not because of what the column is called.
	recordSchema := pg.NewSchema([]schema.FieldSchema{
		{Name: "payload", Type: "jsonb", IsJSON: true, IsJSONB: true},
	})
	provider := pg.NewTypeProvider(map[string]pg.Schema{"record": recordSchema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("record", cel.ObjectType("record")),
	)
	require.NoError(t, err)

	ast, issues := env.Compile(`has(record.payload.active)`)
	require.NoError(t, issues.Err())

	schemas := provider.GetSchemas()

	dialectResult, err := cel2sql.ConvertParameterized(ast, cel2sql.WithSchemas(schemas))
	require.NoError(t, err, "the default style has no collision")
	require.Contains(t, dialectResult.SQL, "?", "precondition: postgres emits the ? operator")

	_, err = cel2sql.ConvertParameterized(ast,
		cel2sql.WithSchemas(schemas),
		cel2sql.WithPlaceholderStyle(cel2sql.PlaceholderQuestion))
	require.Error(t, err)
	assert.ErrorIs(t, err, cel2sql.ErrUnsupportedDialectFeature)
	assert.Contains(t, err.Error(), "PlaceholderQuestion")
}
