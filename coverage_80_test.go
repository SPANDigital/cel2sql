package cel2sql_test

import (
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/ext"
	"github.com/spandigital/cel2sql/v3"
	"github.com/spandigital/cel2sql/v3/pg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test coverage improvements for issue #91
// Target: Improve coverage from 72% to 80%

// TestCallContains_EdgeCases tests contains() edge cases to improve callContains from 56.5%
func TestCallContains_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected string
		schema   pg.Schema
	}{
		{
			name:     "contains with NULL handling - string",
			expr:     "person.name.contains('test')",
			expected: "POSITION('test' IN person.name) > 0",
			schema: pg.NewSchema([]pg.FieldSchema{
				{Name: "name", Type: "text"},
			}),
		},
		{
			name:     "contains with empty string",
			expr:     "'hello'.contains('')",
			expected: "POSITION('' IN 'hello') > 0",
			schema:   pg.NewSchema([]pg.FieldSchema{}),
		},
		{
			name:     "contains with tab character",
			expr:     "person.path.contains('\\t')",
			expected: "POSITION('\t' IN person.path) > 0",
			schema: pg.NewSchema([]pg.FieldSchema{
				{Name: "path", Type: "text"},
			}),
		},
		{
			name:     "contains with single quote",
			expr:     "person.name.contains(\"'\")",
			expected: "POSITION('''' IN person.name) > 0",
			schema: pg.NewSchema([]pg.FieldSchema{
				{Name: "name", Type: "text"},
			}),
		},
		{
			name:     "contains on field comparison",
			expr:     "person.email.contains(person.domain)",
			expected: "POSITION(person.domain IN person.email) > 0",
			schema: pg.NewSchema([]pg.FieldSchema{
				{Name: "email", Type: "text"},
				{Name: "domain", Type: "text"},
			}),
		},
		{
			name:     "contains in text field",
			expr:     "person.address.contains('street')",
			expected: "POSITION('street' IN person.address) > 0",
			schema: pg.NewSchema([]pg.FieldSchema{
				{Name: "address", Type: "text"},
			}),
		},
		{
			name:     "contains with newline character",
			expr:     "person.text.contains('\\n')",
			expected: "POSITION('\n' IN person.text) > 0",
			schema: pg.NewSchema([]pg.FieldSchema{
				{Name: "text", Type: "text"},
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schemas := map[string]pg.Schema{"person": tt.schema}
			env, err := cel.NewEnv(
				cel.CustomTypeProvider(pg.NewTypeProvider(schemas)),
				cel.Variable("person", cel.ObjectType("person")),
				ext.Strings(),
			)
			require.NoError(t, err)

			ast, issues := env.Compile(tt.expr)
			require.Nil(t, issues)

			sql, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
			require.NoError(t, err)
			assert.Equal(t, tt.expected, sql)
		})
	}
}

// TestVisitExistsOneComprehension_EdgeCases tests exists_one() edge cases to improve from 64.0%
func TestVisitExistsOneComprehension_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected string
	}{
		{
			name:     "exists_one with zero matches",
			expr:     "data.numbers.exists_one(x, x > 1000)",
			expected: "(SELECT COUNT(*) FROM UNNEST(data.numbers) AS x WHERE x > 1000) = 1",
		},
		{
			name:     "exists_one with exactly one match",
			expr:     "data.numbers.exists_one(x, x == 42)",
			expected: "(SELECT COUNT(*) FROM UNNEST(data.numbers) AS x WHERE x = 42) = 1",
		},
		{
			name:     "exists_one with multiple potential matches",
			expr:     "data.numbers.exists_one(x, x < 10)",
			expected: "(SELECT COUNT(*) FROM UNNEST(data.numbers) AS x WHERE x < 10) = 1",
		},
		{
			name:     "exists_one with complex predicate",
			expr:     "data.numbers.exists_one(x, x > 50 && x < 60 && x % 2 == 0)",
			expected: "(SELECT COUNT(*) FROM UNNEST(data.numbers) AS x WHERE x > 50 AND x < 60 AND MOD(x, 2) = 0) = 1",
		},
		{
			name:     "exists_one with OR condition",
			expr:     "data.numbers.exists_one(x, x == 1 || x == 100)",
			expected: "(SELECT COUNT(*) FROM UNNEST(data.numbers) AS x WHERE x = 1 OR x = 100) = 1",
		},
		{
			name:     "exists_one with string array",
			expr:     "data.tags.exists_one(t, t.startsWith('prod'))",
			expected: "(SELECT COUNT(*) FROM UNNEST(data.tags) AS t WHERE t LIKE 'prod%' ESCAPE E'\\\\') = 1",
		},
		{
			name:     "exists_one with negation in predicate",
			expr:     "data.numbers.exists_one(x, x != 0 && x < 5)",
			expected: "(SELECT COUNT(*) FROM UNNEST(data.numbers) AS x WHERE x != 0 AND x < 5) = 1",
		},
		{
			name:     "exists_one used in larger expression",
			expr:     "data.numbers.exists_one(x, x == 5) && data.numbers.size() > 3",
			expected: "(SELECT COUNT(*) FROM UNNEST(data.numbers) AS x WHERE x = 5) = 1 AND COALESCE(ARRAY_LENGTH(data.numbers, 1), 0) > 3",
		},
	}

	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "numbers", Type: "integer", Repeated: true},
		{Name: "tags", Type: "text", Repeated: true},
	})
	schemas := map[string]pg.Schema{"data": schema}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env, err := cel.NewEnv(
				cel.CustomTypeProvider(pg.NewTypeProvider(schemas)),
				cel.Variable("data", cel.ObjectType("data")),
				ext.Strings(),
			)
			require.NoError(t, err)

			ast, issues := env.Compile(tt.expr)
			require.Nil(t, issues)

			sql, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
			require.NoError(t, err)
			assert.Equal(t, tt.expected, sql)
		})
	}
}

// TestVisitMapComprehension_EdgeCases tests map() edge cases to improve from 66.7%
func TestVisitMapComprehension_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected string
	}{
		{
			name:     "map with complex arithmetic transformation",
			expr:     "data.numbers.map(x, (x * 2 + 10) / 3)",
			expected: "ARRAY(SELECT (x * 2 + 10) / 3 FROM UNNEST(data.numbers) AS x)",
		},
		{
			name:     "map with ternary operator",
			expr:     "data.numbers.map(x, x > 0 ? x : -x)",
			expected: "ARRAY(SELECT CASE WHEN x > 0 THEN x ELSE -x END FROM UNNEST(data.numbers) AS x)",
		},
		{
			name:     "map with string concatenation",
			expr:     "data.tags.map(t, 'prefix_' + t)",
			expected: "ARRAY(SELECT 'prefix_' || t FROM UNNEST(data.tags) AS t)",
		},
		{
			name:     "map with nested function calls",
			expr:     "data.tags.map(t, t.upperAscii() + '_SUFFIX')",
			expected: "ARRAY(SELECT UPPER(t) || '_SUFFIX' FROM UNNEST(data.tags) AS t)",
		},
		{
			name:     "map with boolean transformation",
			expr:     "data.numbers.map(x, x > 10)",
			expected: "ARRAY(SELECT x > 10 FROM UNNEST(data.numbers) AS x)",
		},
		{
			name:     "map with multiple operations",
			expr:     "data.numbers.map(x, x * x - x + 1)",
			expected: "ARRAY(SELECT x * x - x + 1 FROM UNNEST(data.numbers) AS x)",
		},
		{
			name:     "map with division and modulo",
			expr:     "data.numbers.map(x, x / 10 + x % 10)",
			expected: "ARRAY(SELECT x / 10 + MOD(x, 10) FROM UNNEST(data.numbers) AS x)",
		},
		{
			name:     "map identity transformation",
			expr:     "data.numbers.map(x, x)",
			expected: "ARRAY(SELECT x FROM UNNEST(data.numbers) AS x)",
		},
		{
			name:     "map with size() result",
			expr:     "data.tags.map(t, t.size())",
			expected: "ARRAY(SELECT LENGTH(t) FROM UNNEST(data.tags) AS t)",
		},
		{
			name:     "nested map operations",
			expr:     "data.numbers.map(x, x * 2).map(y, y + 1)",
			expected: "ARRAY(SELECT y + 1 FROM UNNEST(ARRAY(SELECT x * 2 FROM UNNEST(data.numbers) AS x)) AS y)",
		},
	}

	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "numbers", Type: "integer", Repeated: true},
		{Name: "tags", Type: "text", Repeated: true},
	})
	schemas := map[string]pg.Schema{"data": schema}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env, err := cel.NewEnv(
				cel.CustomTypeProvider(pg.NewTypeProvider(schemas)),
				cel.Variable("data", cel.ObjectType("data")),
				ext.Strings(),
			)
			require.NoError(t, err)

			ast, issues := env.Compile(tt.expr)
			require.Nil(t, issues)

			sql, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
			require.NoError(t, err)
			assert.Equal(t, tt.expected, sql)
		})
	}
}

// TestCallIndexOf_EdgeCases tests indexOf() edge cases to improve from 67.7%
func TestCallIndexOf_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected string
	}{
		{
			name:     "indexOf not found returns -1",
			expr:     "person.email.indexOf('notfound') == -1",
			expected: "CASE WHEN POSITION('notfound' IN person.email) > 0 THEN POSITION('notfound' IN person.email) - 1 ELSE -1 END = -1",
		},
		{
			name:     "indexOf with zero offset",
			expr:     "person.text.indexOf('test', 0)",
			expected: "CASE WHEN POSITION('test' IN SUBSTRING(person.text, 1)) > 0 THEN POSITION('test' IN SUBSTRING(person.text, 1)) + 0 - 1 ELSE -1 END",
		},
		{
			name:     "indexOf with large offset",
			expr:     "person.text.indexOf('test', 100) >= 0",
			expected: "CASE WHEN POSITION('test' IN SUBSTRING(person.text, 101)) > 0 THEN POSITION('test' IN SUBSTRING(person.text, 101)) + 100 - 1 ELSE -1 END >= 0",
		},
		{
			name:     "indexOf with negative result comparison",
			expr:     "person.email.indexOf('@') < 0",
			expected: "CASE WHEN POSITION('@' IN person.email) > 0 THEN POSITION('@' IN person.email) - 1 ELSE -1 END < 0",
		},
		{
			name:     "indexOf at beginning",
			expr:     "person.text.indexOf('start') == 0",
			expected: "CASE WHEN POSITION('start' IN person.text) > 0 THEN POSITION('start' IN person.text) - 1 ELSE -1 END = 0",
		},
		{
			name:     "indexOf with offset beyond string length",
			expr:     "person.email.indexOf('test', 500)",
			expected: "CASE WHEN POSITION('test' IN SUBSTRING(person.email, 501)) > 0 THEN POSITION('test' IN SUBSTRING(person.email, 501)) + 500 - 1 ELSE -1 END",
		},
		{
			name:     "indexOf in conditional expression",
			expr:     "person.email.indexOf('@') > 0 && person.email.indexOf('.') > 0",
			expected: "CASE WHEN POSITION('@' IN person.email) > 0 THEN POSITION('@' IN person.email) - 1 ELSE -1 END > 0 AND CASE WHEN POSITION('.' IN person.email) > 0 THEN POSITION('.' IN person.email) - 1 ELSE -1 END > 0",
		},
	}

	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "email", Type: "text"},
		{Name: "text", Type: "text"},
	})
	schemas := map[string]pg.Schema{"person": schema}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env, err := cel.NewEnv(
				cel.CustomTypeProvider(pg.NewTypeProvider(schemas)),
				cel.Variable("person", cel.ObjectType("person")),
				ext.Strings(),
			)
			require.NoError(t, err)

			ast, issues := env.Compile(tt.expr)
			require.Nil(t, issues)

			sql, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
			require.NoError(t, err)
			assert.Equal(t, tt.expected, sql)
		})
	}
}

// TestCallSubstring_EdgeCases tests substring() edge cases to improve from 68.0%
func TestCallSubstring_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected string
	}{
		{
			name:     "substring with negative start index",
			expr:     "person.name.substring(-5)",
			expected: "SUBSTRING(person.name, -4)",
		},
		{
			name:     "substring with start beyond string length",
			expr:     "person.name.substring(1000)",
			expected: "SUBSTRING(person.name, 1001)",
		},
		{
			name:     "substring with zero start",
			expr:     "person.name.substring(0, 5)",
			expected: "SUBSTRING(person.name, 1, 5)",
		},
		{
			name:     "substring with end less than start",
			expr:     "person.name.substring(10, 5)",
			expected: "SUBSTRING(person.name, 11, 0)",
		},
		{
			name:     "substring with equal start and end",
			expr:     "person.name.substring(5, 5)",
			expected: "SUBSTRING(person.name, 6, 0)",
		},
		{
			name:     "substring from middle to end (single arg)",
			expr:     "person.text.substring(10) == 'result'",
			expected: "SUBSTRING(person.text, 11) = 'result'",
		},
		{
			name:     "substring with field references",
			expr:     "person.name.substring(person.start_pos, person.end_pos)",
			expected: "SUBSTRING(person.name, person.start_pos + 1, person.end_pos - (person.start_pos))",
		},
		{
			name:     "substring in comparison",
			expr:     "person.email.substring(0, 5).lowerAscii() == 'admin'",
			expected: "LOWER(SUBSTRING(person.email, 1, 5)) = 'admin'",
		},
		{
			name:     "substring with large indices",
			expr:     "person.text.substring(100, 200)",
			expected: "SUBSTRING(person.text, 101, 100)",
		},
		{
			name:     "nested substring calls",
			expr:     "person.name.substring(1, 10).substring(0, 3)",
			expected: "SUBSTRING(SUBSTRING(person.name, 2, 9), 1, 3)",
		},
	}

	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "name", Type: "text"},
		{Name: "email", Type: "text"},
		{Name: "text", Type: "text"},
		{Name: "start_pos", Type: "integer"},
		{Name: "end_pos", Type: "integer"},
	})
	schemas := map[string]pg.Schema{"person": schema}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env, err := cel.NewEnv(
				cel.CustomTypeProvider(pg.NewTypeProvider(schemas)),
				cel.Variable("person", cel.ObjectType("person")),
				ext.Strings(),
			)
			require.NoError(t, err)

			ast, issues := env.Compile(tt.expr)
			require.Nil(t, issues)

			sql, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
			require.NoError(t, err)
			assert.Equal(t, tt.expected, sql)
		})
	}
}

// TestVisitAllComprehension_EdgeCases tests all() edge cases to improve from 69.2%
func TestVisitAllComprehension_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected string
	}{
		{
			name:     "all with empty array behavior",
			expr:     "data.numbers.all(x, x < 0)",
			expected: "NOT EXISTS (SELECT 1 FROM UNNEST(data.numbers) AS x WHERE NOT (x < 0))",
		},
		{
			name:     "all with complex OR condition",
			expr:     "data.numbers.all(x, x < 10 || x > 100)",
			expected: "NOT EXISTS (SELECT 1 FROM UNNEST(data.numbers) AS x WHERE NOT (x < 10 OR x > 100))",
		},
		{
			name:     "all with multiple AND conditions",
			expr:     "data.numbers.all(x, x > 0 && x < 100 && x % 2 == 0)",
			expected: "NOT EXISTS (SELECT 1 FROM UNNEST(data.numbers) AS x WHERE NOT (x > 0 AND x < 100 AND MOD(x, 2) = 0))",
		},
		{
			name:     "all with nested parentheses",
			expr:     "data.numbers.all(x, ((x > 5) && (x < 50)) || (x == 100))",
			expected: "NOT EXISTS (SELECT 1 FROM UNNEST(data.numbers) AS x WHERE NOT (x > 5 AND x < 50 OR x = 100))",
		},
		{
			name:     "all with string operations",
			expr:     "data.tags.all(t, t.startsWith('valid') && t.size() > 5)",
			expected: "NOT EXISTS (SELECT 1 FROM UNNEST(data.tags) AS t WHERE NOT (t LIKE 'valid%' ESCAPE E'\\\\' AND LENGTH(t) > 5))",
		},
		{
			name:     "all with negation",
			expr:     "data.numbers.all(x, !(x < 0))",
			expected: "NOT EXISTS (SELECT 1 FROM UNNEST(data.numbers) AS x WHERE NOT (NOT (x < 0)))",
		},
		{
			name:     "nested all operations",
			expr:     "data.numbers.all(x, x > 0) && data.numbers.all(y, y < 100)",
			expected: "NOT EXISTS (SELECT 1 FROM UNNEST(data.numbers) AS x WHERE NOT (x > 0)) AND NOT EXISTS (SELECT 1 FROM UNNEST(data.numbers) AS y WHERE NOT (y < 100))",
		},
		{
			name:     "all with comparison to field",
			expr:     "data.numbers.all(x, x >= data.min_value)",
			expected: "NOT EXISTS (SELECT 1 FROM UNNEST(data.numbers) AS x WHERE NOT (x >= data.min_value))",
		},
		{
			name:     "all with inequality",
			expr:     "data.numbers.all(x, x != 0)",
			expected: "NOT EXISTS (SELECT 1 FROM UNNEST(data.numbers) AS x WHERE NOT (x != 0))",
		},
		{
			name:     "all in ternary expression",
			expr:     "data.numbers.all(x, x > 0) ? 'valid' : 'invalid'",
			expected: "CASE WHEN NOT EXISTS (SELECT 1 FROM UNNEST(data.numbers) AS x WHERE NOT (x > 0)) THEN 'valid' ELSE 'invalid' END",
		},
	}

	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "numbers", Type: "integer", Repeated: true},
		{Name: "tags", Type: "text", Repeated: true},
		{Name: "min_value", Type: "integer"},
	})
	schemas := map[string]pg.Schema{"data": schema}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env, err := cel.NewEnv(
				cel.CustomTypeProvider(pg.NewTypeProvider(schemas)),
				cel.Variable("data", cel.ObjectType("data")),
				ext.Strings(),
			)
			require.NoError(t, err)

			ast, issues := env.Compile(tt.expr)
			require.Nil(t, issues)

			sql, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
			require.NoError(t, err)
			assert.Equal(t, tt.expected, sql)
		})
	}
}

// TestVisitFilterComprehension_EdgeCases tests filter() edge cases to improve from 70.4%
func TestVisitFilterComprehension_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected string
	}{
		{
			name:     "filter with complex AND/OR conditions",
			expr:     "data.numbers.filter(x, (x > 10 && x < 20) || (x > 30 && x < 40) || x == 50)",
			expected: "ARRAY(SELECT x FROM UNNEST(data.numbers) AS x WHERE x > 10 AND x < 20 OR x > 30 AND x < 40 OR x = 50)",
		},
		{
			name:     "filter with nested field access on JSON",
			expr:     "data.records.filter(r, r.metadata.active == true)",
			expected: "ARRAY(SELECT r FROM UNNEST(data.records) AS r WHERE r.metadata.active IS TRUE)",
		},
		{
			name:     "filter with multiple string conditions",
			expr:     "data.tags.filter(t, t.startsWith('a') && t.endsWith('z') && t.contains('test'))",
			expected: "ARRAY(SELECT t FROM UNNEST(data.tags) AS t WHERE t LIKE 'a%' ESCAPE E'\\\\' AND t LIKE '%z' ESCAPE E'\\\\' AND POSITION('test' IN t) > 0)",
		},
		{
			name:     "filter with arithmetic in predicate",
			expr:     "data.numbers.filter(x, (x * 2) % 3 == 0)",
			expected: "ARRAY(SELECT x FROM UNNEST(data.numbers) AS x WHERE MOD(x * 2, 3) = 0)",
		},
		{
			name:     "filter on JSON array",
			expr:     "data.json_items.filter(item, item.value > 100)",
			expected: "ARRAY(SELECT item FROM UNNEST(data.json_items) AS item WHERE (item->>'value')::numeric > 100)",
		},
		{
			name:     "filter with size comparison",
			expr:     "data.tags.filter(t, t.size() >= 5 && t.size() <= 10)",
			expected: "ARRAY(SELECT t FROM UNNEST(data.tags) AS t WHERE LENGTH(t) >= 5 AND LENGTH(t) <= 10)",
		},
		{
			name:     "filter then check size",
			expr:     "data.numbers.filter(x, x > 0).size() == 0",
			expected: "COALESCE(ARRAY_LENGTH(ARRAY(SELECT x FROM UNNEST(data.numbers) AS x WHERE x > 0), 1), 0) = 0",
		},
	}

	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "numbers", Type: "integer", Repeated: true},
		{Name: "tags", Type: "text", Repeated: true},
		{Name: "records", Type: "jsonb", Repeated: true},
		{Name: "json_items", Type: "jsonb", Repeated: true},
	})
	schemas := map[string]pg.Schema{"data": schema}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env, err := cel.NewEnv(
				cel.CustomTypeProvider(pg.NewTypeProvider(schemas)),
				cel.Variable("data", cel.ObjectType("data")),
				ext.Strings(),
			)
			require.NoError(t, err)

			ast, issues := env.Compile(tt.expr)
			require.Nil(t, issues)

			sql, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
			require.NoError(t, err)
			assert.Equal(t, tt.expected, sql)
		})
	}
}

// TestVisitExistsComprehension_EdgeCases tests exists() edge cases to improve from 72.0%
func TestVisitExistsComprehension_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected string
	}{
		{
			name:     "exists with complex predicates",
			expr:     "data.numbers.exists(x, x > 100 && x < 200 && x % 5 == 0)",
			expected: "EXISTS (SELECT 1 FROM UNNEST(data.numbers) AS x WHERE x > 100 AND x < 200 AND MOD(x, 5) = 0)",
		},
		{
			name:     "nested exists operations",
			expr:     "data.numbers.exists(x, x > 50) && data.numbers.exists(y, y < 10)",
			expected: "EXISTS (SELECT 1 FROM UNNEST(data.numbers) AS x WHERE x > 50) AND EXISTS (SELECT 1 FROM UNNEST(data.numbers) AS y WHERE y < 10)",
		},
		{
			name:     "exists with OR conditions",
			expr:     "data.numbers.exists(x, x == 1 || x == 2 || x == 3 || x == 4)",
			expected: "EXISTS (SELECT 1 FROM UNNEST(data.numbers) AS x WHERE x = 1 OR x = 2 OR x = 3 OR x = 4)",
		},
		{
			name:     "exists with string matching",
			expr:     "data.tags.exists(t, t.contains('prod') && t.size() > 10)",
			expected: "EXISTS (SELECT 1 FROM UNNEST(data.tags) AS t WHERE POSITION('prod' IN t) > 0 AND LENGTH(t) > 10)",
		},
		{
			name:     "exists in negated expression",
			expr:     "!data.numbers.exists(x, x < 0) && data.numbers.size() > 0",
			expected: "NOT EXISTS (SELECT 1 FROM UNNEST(data.numbers) AS x WHERE x < 0) AND COALESCE(ARRAY_LENGTH(data.numbers, 1), 0) > 0",
		},
		{
			name:     "exists with field comparison",
			expr:     "data.numbers.exists(x, x >= data.threshold)",
			expected: "EXISTS (SELECT 1 FROM UNNEST(data.numbers) AS x WHERE x >= data.threshold)",
		},
	}

	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "numbers", Type: "integer", Repeated: true},
		{Name: "tags", Type: "text", Repeated: true},
		{Name: "threshold", Type: "integer"},
	})
	schemas := map[string]pg.Schema{"data": schema}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env, err := cel.NewEnv(
				cel.CustomTypeProvider(pg.NewTypeProvider(schemas)),
				cel.Variable("data", cel.ObjectType("data")),
				ext.Strings(),
			)
			require.NoError(t, err)

			ast, issues := env.Compile(tt.expr)
			require.Nil(t, issues)

			sql, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
			require.NoError(t, err)
			assert.Equal(t, tt.expected, sql)
		})
	}
}
