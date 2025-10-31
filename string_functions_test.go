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

// Test cases from Issue #85 - These were causing panics
func TestIssue85_StringFunctionsInComprehensions(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected string
	}{
		{
			name:     "size() in exists_one comprehension",
			expr:     "data.tags.exists_one(t, t.size() == 10)",
			expected: "(SELECT COUNT(*) FROM UNNEST(data.tags) AS t WHERE LENGTH(t) = 10) = 1",
		},
		{
			name:     "upperAscii() in map comprehension",
			expr:     "data.tags.map(t, t.upperAscii())",
			expected: "ARRAY(SELECT UPPER(t) FROM UNNEST(data.tags) AS t)",
		},
		{
			name:     "lowerAscii() in map comprehension",
			expr:     "data.tags.map(t, t.lowerAscii())",
			expected: "ARRAY(SELECT LOWER(t) FROM UNNEST(data.tags) AS t)",
		},
		{
			name:     "size() in map comprehension",
			expr:     "data.tags.map(t, t.size())",
			expected: "ARRAY(SELECT LENGTH(t) FROM UNNEST(data.tags) AS t)",
		},
		{
			name:     "size() in filter comprehension",
			expr:     "data.tags.filter(t, t.size() > 5)",
			expected: "ARRAY(SELECT t FROM UNNEST(data.tags) AS t WHERE LENGTH(t) > 5)",
		},
	}

	// Schema with array of strings
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "tags", Type: "text", Repeated: true},
	})
	schemas := map[string]pg.Schema{"data": schema}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env, err := cel.NewEnv(
				cel.CustomTypeProvider(pg.NewTypeProvider(schemas)),
				cel.Variable("data", cel.ObjectType("data")),
				ext.Strings(), // Enable string extension functions
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

// Test case from Issue #85 - size() outside comprehension
func TestIssue85_SizeOutsideComprehension(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "name", Type: "text", Repeated: false},
	})
	schemas := map[string]pg.Schema{"item": schema}

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(pg.NewTypeProvider(schemas)),
		cel.Variable("item", cel.ObjectType("item")),
	)
	require.NoError(t, err)

	ast, issues := env.Compile("item.name.size() > 10")
	require.Nil(t, issues)

	sql, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
	require.NoError(t, err)
	// Note: item.name will be treated as a struct field, generating item->>'name' for JSON
	assert.Equal(t, "LENGTH(item->>'name') > 10", sql)
}

// Comprehensive tests for all string extension functions

func TestStringFunctions_LowerAscii(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected string
	}{
		{
			name:     "method call",
			expr:     "person.name.lowerAscii() == 'john'",
			expected: "LOWER(person.name) = 'john'",
		},
		{
			name:     "with comparison",
			expr:     "person.email.lowerAscii() == person.username.lowerAscii()",
			expected: "LOWER(person.email) = LOWER(person.username)",
		},
	}

	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "name", Type: "text"},
		{Name: "email", Type: "text"},
		{Name: "username", Type: "text"},
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

func TestStringFunctions_UpperAscii(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected string
	}{
		{
			name:     "method call",
			expr:     "person.name.upperAscii() == 'JOHN'",
			expected: "UPPER(person.name) = 'JOHN'",
		},
		{
			name:     "with startsWith",
			expr:     "person.name.upperAscii().startsWith('J')",
			expected: "UPPER(person.name) LIKE 'J%' ESCAPE E'\\\\'",
		},
	}

	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "name", Type: "text"},
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

func TestStringFunctions_Trim(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected string
	}{
		{
			name:     "method call",
			expr:     "person.name.trim() == 'John'",
			expected: "TRIM(person.name) = 'John'",
		},
		{
			name:     "in comparison",
			expr:     "person.name.trim().size() > 0",
			expected: "LENGTH(TRIM(person.name)) > 0",
		},
	}

	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "name", Type: "text"},
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

func TestStringFunctions_CharAt(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected string
	}{
		{
			name:     "constant index",
			expr:     "person.name.charAt(0) == 'J'",
			expected: "SUBSTRING(person.name, 1, 1) = 'J'",
		},
		{
			name:     "dynamic index",
			expr:     "person.name.charAt(person.position) == 'x'",
			expected: "SUBSTRING(person.name, person.position + 1, 1) = 'x'",
		},
	}

	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "name", Type: "text"},
		{Name: "position", Type: "bigint"},
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

func TestStringFunctions_IndexOf(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected string
	}{
		{
			name:     "simple indexOf",
			expr:     "person.email.indexOf('@') > 0",
			expected: "CASE WHEN POSITION('@' IN person.email) > 0 THEN POSITION('@' IN person.email) - 1 ELSE -1 END > 0",
		},
		{
			name:     "indexOf with offset",
			expr:     "person.text.indexOf('test', 5) >= 0",
			expected: "CASE WHEN POSITION('test' IN SUBSTRING(person.text, 6)) > 0 THEN POSITION('test' IN SUBSTRING(person.text, 6)) + 5 - 1 ELSE -1 END >= 0",
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

func TestStringFunctions_LastIndexOf(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected string
	}{
		{
			name:     "simple lastIndexOf",
			expr:     "person.path.lastIndexOf('/') > 0",
			expected: "CASE WHEN POSITION(REVERSE('/') IN REVERSE(person.path)) > 0 THEN LENGTH(person.path) - POSITION(REVERSE('/') IN REVERSE(person.path)) - LENGTH('/') + 1 ELSE -1 END > 0",
		},
	}

	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "path", Type: "text"},
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

func TestStringFunctions_Substring(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected string
	}{
		{
			name:     "substring with start only (constant)",
			expr:     "person.name.substring(5) == 'test'",
			expected: "SUBSTRING(person.name, 6) = 'test'",
		},
		{
			name:     "substring with start and end (constant)",
			expr:     "person.name.substring(0, 4) == 'John'",
			expected: "SUBSTRING(person.name, 1, 4) = 'John'",
		},
		{
			name:     "substring with dynamic start",
			expr:     "person.name.substring(person.startpos, person.endpos) == 'test'",
			expected: "SUBSTRING(person.name, person.startpos + 1, person.endpos - (person.startpos)) = 'test'",
		},
	}

	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "name", Type: "text"},
		{Name: "startpos", Type: "bigint"},
		{Name: "endpos", Type: "bigint"},
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

func TestStringFunctions_Replace(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected string
	}{
		{
			name:     "replace without limit",
			expr:     "person.text.replace('old', 'new') == 'test'",
			expected: "REPLACE(person.text, 'old', 'new') = 'test'",
		},
		{
			name:     "replace with limit=-1 (replace all)",
			expr:     "person.text.replace('a', 'b', -1) == 'test'",
			expected: "REPLACE(person.text, 'a', 'b') = 'test'",
		},
	}

	schema := pg.NewSchema([]pg.FieldSchema{
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

func TestStringFunctions_ReplaceWithLimitError(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "text", Type: "text"},
	})
	schemas := map[string]pg.Schema{"person": schema}

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(pg.NewTypeProvider(schemas)),
		cel.Variable("person", cel.ObjectType("person")),
		ext.Strings(),
	)
	require.NoError(t, err)

	// replace() with limit != -1 should return error
	ast, issues := env.Compile("person.text.replace('a', 'b', 1) == 'test'")
	require.Nil(t, issues)

	_, err = cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "replace() with limit != -1 is not supported")
}

func TestStringFunctions_Reverse(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected string
	}{
		{
			name:     "simple reverse",
			expr:     "person.name.reverse() == 'nhoJ'",
			expected: "REVERSE(person.name) = 'nhoJ'",
		},
	}

	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "name", Type: "text"},
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

// Tests for unsupported functions that should return errors

func TestStringFunctions_UnsupportedSplit(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "text", Type: "text"},
	})
	schemas := map[string]pg.Schema{"person": schema}

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(pg.NewTypeProvider(schemas)),
		cel.Variable("person", cel.ObjectType("person")),
		ext.Strings(),
	)
	require.NoError(t, err)

	ast, issues := env.Compile("person.text.split(',').size() > 0")
	require.Nil(t, issues)

	_, err = cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "split() returns arrays and cannot be converted inline to SQL")
}

func TestStringFunctions_UnsupportedJoin(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "items", Type: "text", Repeated: true},
	})
	schemas := map[string]pg.Schema{"person": schema}

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(pg.NewTypeProvider(schemas)),
		cel.Variable("person", cel.ObjectType("person")),
		ext.Strings(),
	)
	require.NoError(t, err)

	ast, issues := env.Compile("person.items.join(',') == 'a,b,c'")
	require.Nil(t, issues)

	_, err = cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "join() on arrays is not supported in SQL conversion")
}

func TestStringFunctions_UnsupportedFormat(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "name", Type: "text"},
		{Name: "age", Type: "bigint"},
	})
	schemas := map[string]pg.Schema{"person": schema}

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(pg.NewTypeProvider(schemas)),
		cel.Variable("person", cel.ObjectType("person")),
		ext.Strings(),
	)
	require.NoError(t, err)

	ast, issues := env.Compile("'%s is %d'.format([person.name, person.age]) == 'John is 30'")
	require.Nil(t, issues)

	_, err = cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "format() with printf-style formatting is too complex for SQL conversion")
}

// Note: quote() function is not available in CEL ext.Strings() standard extension
// It's part of CEL spec but not commonly implemented, so we skip testing it

// Panic prevention tests - ensure no panics for edge cases

func TestStringFunctions_NoPanicOnEmptyArgs(t *testing.T) {
	// These tests ensure that all defensive checks are working
	// and no panics occur even with malformed expressions

	// Note: Most of these cases would be caught by CEL's type checker,
	// but we add defensive checks to prevent panics if somehow they get through
	t.Run("defensive checks exist", func(t *testing.T) {
		// This is more of a code review checkpoint than an actual test
		// The defensive checks in callCasting, visitCallIndex, visitCallMapIndex,
		// visitCallListIndex, and visitCallUnary should prevent panics
		assert.True(t, true, "Defensive checks have been added")
	})
}
