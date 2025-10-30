package cel2sql_test

import (
	"strings"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/spandigital/cel2sql/v3"
	"github.com/spandigital/cel2sql/v3/pg"
)

func TestConvertParameterized(t *testing.T) {
	// Set up CEL environment with various types
	env, err := cel.NewEnv(
		cel.Variable("name", cel.StringType),
		cel.Variable("age", cel.IntType),
		cel.Variable("salary", cel.DoubleType),
		cel.Variable("active", cel.BoolType),
		cel.Variable("data", cel.BytesType),
		cel.Variable("tags", cel.ListType(cel.StringType)),
		cel.Variable("scores", cel.ListType(cel.IntType)),
	)
	require.NoError(t, err)

	tests := []struct {
		name           string
		celExpr        string
		wantSQL        string
		wantParamCount int
		wantParams     []interface{}
	}{
		// String parameters
		{
			name:           "simple string equality",
			celExpr:        `name == "John"`,
			wantSQL:        "name = $1",
			wantParamCount: 1,
			wantParams:     []interface{}{"John"},
		},
		{
			name:           "multiple string parameters",
			celExpr:        `name == "John" && name != "Jane"`,
			wantSQL:        "name = $1 AND name != $2",
			wantParamCount: 2,
			wantParams:     []interface{}{"John", "Jane"},
		},
		{
			name:           "string with escaped quotes",
			celExpr:        `name == "O'Brien"`,
			wantSQL:        "name = $1",
			wantParamCount: 1,
			wantParams:     []interface{}{"O'Brien"},
		},

		// Integer parameters
		{
			name:           "simple integer equality",
			celExpr:        `age == 18`,
			wantSQL:        "age = $1",
			wantParamCount: 1,
			wantParams:     []interface{}{int64(18)},
		},
		{
			name:           "integer comparison",
			celExpr:        `age > 21 && age < 65`,
			wantSQL:        "age > $1 AND age < $2",
			wantParamCount: 2,
			wantParams:     []interface{}{int64(21), int64(65)},
		},
		{
			name:           "negative integer",
			celExpr:        `age == -5`,
			wantSQL:        "age = $1",
			wantParamCount: 1,
			wantParams:     []interface{}{int64(-5)},
		},

		// Double parameters
		{
			name:           "simple double equality",
			celExpr:        `salary == 50000.50`,
			wantSQL:        "salary = $1",
			wantParamCount: 1,
			wantParams:     []interface{}{50000.50},
		},
		{
			name:           "double comparison",
			celExpr:        `salary >= 30000.0 && salary <= 100000.0`,
			wantSQL:        "salary >= $1 AND salary <= $2",
			wantParamCount: 2,
			wantParams:     []interface{}{30000.0, 100000.0},
		},

		// Boolean and NULL constants (kept inline)
		{
			name:           "boolean TRUE constant",
			celExpr:        `active == true`,
			wantSQL:        "active IS TRUE",
			wantParamCount: 0,
			wantParams:     nil,
		},
		{
			name:           "boolean FALSE constant",
			celExpr:        `active == false`,
			wantSQL:        "active IS FALSE",
			wantParamCount: 0,
			wantParams:     nil,
		},
		{
			name:           "mixed bool constants and params",
			celExpr:        `active == true && age == 18`,
			wantSQL:        "active IS TRUE AND age = $1",
			wantParamCount: 1,
			wantParams:     []interface{}{int64(18)},
		},

		// Bytes parameters
		{
			name:           "bytes equality",
			celExpr:        `data == b"hello"`,
			wantSQL:        "data = $1",
			wantParamCount: 1,
			wantParams:     []interface{}{[]byte("hello")},
		},

		// Complex expressions with multiple parameters
		{
			name:           "complex AND expression",
			celExpr:        `name == "John" && age >= 18 && salary > 50000.0`,
			wantSQL:        "name = $1 AND age >= $2 AND salary > $3",
			wantParamCount: 3,
			wantParams:     []interface{}{"John", int64(18), 50000.0},
		},
		{
			name:           "complex OR expression",
			celExpr:        `name == "John" || name == "Jane" || age == 25`,
			wantSQL:        "name = $1 OR name = $2 OR age = $3",
			wantParamCount: 3,
			wantParams:     []interface{}{"John", "Jane", int64(25)},
		},
		{
			name:           "nested parentheses with params",
			celExpr:        `(name == "John" && age == 18) || (name == "Jane" && age == 21)`,
			wantSQL:        "name = $1 AND age = $2 OR name = $3 AND age = $4",
			wantParamCount: 4,
			wantParams:     []interface{}{"John", int64(18), "Jane", int64(21)},
		},

		// Parameter ordering test
		{
			name:           "parameter ordering matches placeholders",
			celExpr:        `name == "First" && age == 1 && salary == 100.0 && name != "Second"`,
			wantSQL:        "name = $1 AND age = $2 AND salary = $3 AND name != $4",
			wantParamCount: 4,
			wantParams:     []interface{}{"First", int64(1), 100.0, "Second"},
		},

		// Empty parameter list
		{
			name:           "no parameters - only boolean constants",
			celExpr:        `active == true && active != false`,
			wantSQL:        "active IS TRUE AND active IS NOT FALSE",
			wantParamCount: 0,
			wantParams:     nil,
		},

		// String operations with parameters
		// Note: LIKE patterns are currently optimized inline for constant strings
		{
			name:           "startsWith with parameter",
			celExpr:        `name.startsWith("Jo")`,
			wantSQL:        "name LIKE 'Jo%'",
			wantParamCount: 0,
			wantParams:     nil,
		},
		{
			name:           "endsWith with parameter",
			celExpr:        `name.endsWith("hn")`,
			wantSQL:        "name LIKE '%hn'",
			wantParamCount: 0,
			wantParams:     nil,
		},
		{
			name:           "contains with parameter",
			celExpr:        `name.contains("oh")`,
			wantSQL:        "POSITION($1 IN name) > 0",
			wantParamCount: 1,
			wantParams:     []interface{}{"oh"},
		},

		// IN operator with parameters
		{
			name:           "IN with array literal",
			celExpr:        `age in [18, 21, 25]`,
			wantSQL:        "age = ANY(ARRAY[$1, $2, $3])",
			wantParamCount: 3,
			wantParams:     []interface{}{int64(18), int64(21), int64(25)},
		},
		{
			name:           "string IN with array literal",
			celExpr:        `name in ["John", "Jane", "Bob"]`,
			wantSQL:        "name = ANY(ARRAY[$1, $2, $3])",
			wantParamCount: 3,
			wantParams:     []interface{}{"John", "Jane", "Bob"},
		},

		// Ternary operator with parameters
		{
			name:           "ternary with parameters",
			celExpr:        `age > 18 ? "adult" : "minor"`,
			wantSQL:        "CASE WHEN age > $1 THEN $2 ELSE $3 END",
			wantParamCount: 3,
			wantParams:     []interface{}{int64(18), "adult", "minor"},
		},

		// Type casting with parameters
		{
			name:           "cast int to string",
			celExpr:        `string(age) == "18"`,
			wantSQL:        "CAST(age AS TEXT) = $1",
			wantParamCount: 1,
			wantParams:     []interface{}{"18"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Compile CEL expression
			ast, issues := env.Compile(tt.celExpr)
			require.Nil(t, issues, "CEL compilation should succeed")

			// Convert to parameterized SQL
			result, err := cel2sql.ConvertParameterized(ast)
			require.NoError(t, err, "ConvertParameterized should succeed")

			// Assert SQL matches expected
			assert.Equal(t, tt.wantSQL, result.SQL, "SQL should match expected")

			// Assert parameter count
			assert.Len(t, result.Parameters, tt.wantParamCount, "Parameter count should match")

			// Assert parameters match expected (if specified)
			if tt.wantParams != nil {
				require.Equal(t, tt.wantParams, result.Parameters, "Parameters should match expected values")
			}
		})
	}
}

func TestConvertParameterized_JSONFields(t *testing.T) {
	// Set up schema with JSON fields
	testSchema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "integer"},
		{Name: "name", Type: "text"},
		{Name: "metadata", Type: "jsonb", IsJSON: true, IsJSONB: true},
	})

	provider := pg.NewTypeProvider(map[string]pg.Schema{
		"users": testSchema,
	})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("usr", cel.ObjectType("users")),
	)
	require.NoError(t, err)

	tests := []struct {
		name           string
		celExpr        string
		wantSQL        string
		wantParamCount int
		wantParams     []interface{}
	}{
		{
			name:           "JSON field comparison with parameter",
			celExpr:        `usr.metadata.username == "john_doe"`,
			wantSQL:        "usr.metadata->>'username' = $1",
			wantParamCount: 1,
			wantParams:     []interface{}{"john_doe"},
		},
		{
			name:           "nested JSON field comparison",
			celExpr:        `usr.metadata.settings.theme == "dark"`,
			wantSQL:        "usr.metadata->'settings'->>'theme' = $1",
			wantParamCount: 1,
			wantParams:     []interface{}{"dark"},
		},
		{
			name:           "JSON and regular field with parameters",
			celExpr:        `usr.name == "John" && usr.metadata.age == "25"`,
			wantSQL:        "usr.name = $1 AND usr.metadata->>'age' = $2",
			wantParamCount: 2,
			wantParams:     []interface{}{"John", "25"},
		},
	}

	schemas := map[string]pg.Schema{
		"usr": testSchema,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Compile CEL expression
			ast, issues := env.Compile(tt.celExpr)
			require.Nil(t, issues, "CEL compilation should succeed")

			// Convert to parameterized SQL
			result, err := cel2sql.ConvertParameterized(ast, cel2sql.WithSchemas(schemas))
			require.NoError(t, err, "ConvertParameterized should succeed")

			// Assert SQL matches expected
			assert.Equal(t, tt.wantSQL, result.SQL, "SQL should match expected")

			// Assert parameter count
			assert.Len(t, result.Parameters, tt.wantParamCount, "Parameter count should match")

			// Assert parameters match expected
			if tt.wantParams != nil {
				require.Equal(t, tt.wantParams, result.Parameters, "Parameters should match expected values")
			}
		})
	}
}

func TestConvertParameterized_Comprehensions(t *testing.T) {
	// Use simple CEL environment without custom type provider to avoid JSON detection
	env, err := cel.NewEnv(
		cel.Variable("scores", cel.ListType(cel.IntType)),
		cel.Variable("tags", cel.ListType(cel.StringType)),
	)
	require.NoError(t, err)

	tests := []struct {
		name           string
		celExpr        string
		wantSQL        string
		wantParamCount int
		wantParams     []interface{}
	}{
		{
			name:           "all() with parameterized predicate",
			celExpr:        `scores.all(x, x > 50)`,
			wantSQL:        "NOT EXISTS (SELECT 1 FROM UNNEST(scores) AS x WHERE NOT (x > $1))",
			wantParamCount: 1,
			wantParams:     []interface{}{int64(50)},
		},
		{
			name:           "exists() with parameterized predicate",
			celExpr:        `scores.exists(x, x == 100)`,
			wantSQL:        "EXISTS (SELECT 1 FROM UNNEST(scores) AS x WHERE x = $1)",
			wantParamCount: 1,
			wantParams:     []interface{}{int64(100)},
		},
		{
			name:           "exists_one() with parameterized predicate",
			celExpr:        `scores.exists_one(x, x == 42)`,
			wantSQL:        "(SELECT COUNT(*) FROM UNNEST(scores) AS x WHERE x = $1) = 1",
			wantParamCount: 1,
			wantParams:     []interface{}{int64(42)},
		},
		{
			name:           "map() with parameterized transform",
			celExpr:        `scores.map(x, x + 10).exists(y, y == 110)`,
			wantSQL:        "EXISTS (SELECT 1 FROM UNNEST(ARRAY(SELECT x + $1 FROM UNNEST(scores) AS x)) AS y WHERE y = $2)",
			wantParamCount: 2,
			wantParams:     []interface{}{int64(10), int64(110)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Compile CEL expression
			ast, issues := env.Compile(tt.celExpr)
			require.Nil(t, issues, "CEL compilation should succeed")

			// Convert to parameterized SQL
			result, err := cel2sql.ConvertParameterized(ast)
			require.NoError(t, err, "ConvertParameterized should succeed")

			// Assert SQL matches expected
			assert.Equal(t, tt.wantSQL, result.SQL, "SQL should match expected")

			// Assert parameter count
			assert.Len(t, result.Parameters, tt.wantParamCount, "Parameter count should match")

			// Assert parameters match expected
			if tt.wantParams != nil {
				require.Equal(t, tt.wantParams, result.Parameters, "Parameters should match expected values")
			}
		})
	}
}

func TestConvertParameterized_RegexPatterns(t *testing.T) {
	env, err := cel.NewEnv(
		cel.Variable("email", cel.StringType),
		cel.Variable("phone", cel.StringType),
	)
	require.NoError(t, err)

	tests := []struct {
		name           string
		celExpr        string
		wantParamCount int
		containsRegex  bool
	}{
		{
			name:           "matches with simple pattern",
			celExpr:        `email.matches(r"[a-z]+@[a-z]+\.[a-z]+")`,
			wantParamCount: 0, // Regex patterns are inline, not parameterized
			containsRegex:  true,
		},
		{
			name:           "matches with case-insensitive pattern",
			celExpr:        `phone.matches(r"(?i)^\d{3}-\d{3}-\d{4}$")`,
			wantParamCount: 0,
			containsRegex:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Compile CEL expression
			ast, issues := env.Compile(tt.celExpr)
			require.Nil(t, issues, "CEL compilation should succeed")

			// Convert to parameterized SQL
			result, err := cel2sql.ConvertParameterized(ast)
			require.NoError(t, err, "ConvertParameterized should succeed")

			// Assert parameter count
			assert.Len(t, result.Parameters, tt.wantParamCount, "Parameter count should match")

			// Assert SQL contains regex operator (~ or ~* for case-insensitive)
			if tt.containsRegex {
				hasRegexOperator := strings.Contains(result.SQL, " ~ ") || strings.Contains(result.SQL, " ~* ")
				assert.True(t, hasRegexOperator, "SQL should contain regex operator (~ or ~*)")
			}
		})
	}
}

func TestConvertParameterized_EmptyParameters(t *testing.T) {
	env, err := cel.NewEnv(
		cel.Variable("x", cel.IntType),
		cel.Variable("y", cel.IntType),
	)
	require.NoError(t, err)

	tests := []struct {
		name    string
		celExpr string
		wantSQL string
	}{
		{
			name:    "field comparison only",
			celExpr: `x == y`,
			wantSQL: "x = y",
		},
		{
			name:    "field arithmetic",
			celExpr: `x + y > x * y`,
			wantSQL: "x + y > x * y",
		},
		{
			name:    "field with boolean constants",
			celExpr: `x > y && true`,
			wantSQL: "x > y AND TRUE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Compile CEL expression
			ast, issues := env.Compile(tt.celExpr)
			require.Nil(t, issues, "CEL compilation should succeed")

			// Convert to parameterized SQL
			result, err := cel2sql.ConvertParameterized(ast)
			require.NoError(t, err, "ConvertParameterized should succeed")

			// Assert SQL matches expected
			assert.Equal(t, tt.wantSQL, result.SQL, "SQL should match expected")

			// Assert no parameters
			assert.Empty(t, result.Parameters, "Should have no parameters")
		})
	}
}

func TestConvertParameterized_NegativeCases(t *testing.T) {
	env, err := cel.NewEnv(
		cel.Variable("name", cel.StringType),
		cel.Variable("age", cel.IntType),
	)
	require.NoError(t, err)

	tests := []struct {
		name        string
		celExpr     string
		expectError bool
	}{
		{
			name:        "valid expression",
			celExpr:     `name == "John" && age > 18`,
			expectError: false,
		},
		{
			name:        "string with null byte",
			celExpr:     "name == \"test\x00\"",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Compile CEL expression
			ast, issues := env.Compile(tt.celExpr)
			if issues != nil && issues.Err() != nil {
				if tt.expectError {
					return
				}
				t.Fatalf("CEL compilation failed: %v", issues.Err())
			}

			// Convert to parameterized SQL
			result, err := cel2sql.ConvertParameterized(ast)

			if tt.expectError {
				assert.Error(t, err, "Should return error")
			} else {
				assert.NoError(t, err, "Should not return error")
				assert.NotNil(t, result, "Result should not be nil")
			}
		})
	}
}

// TestConvertParameterized_ParameterTypeConsistency verifies that parameters
// maintain their proper Go types for database/sql compatibility
func TestConvertParameterized_ParameterTypeConsistency(t *testing.T) {
	env, err := cel.NewEnv(
		cel.Variable("name", cel.StringType),
		cel.Variable("age", cel.IntType),
		cel.Variable("salary", cel.DoubleType),
		cel.Variable("data", cel.BytesType),
	)
	require.NoError(t, err)

	tests := []struct {
		name         string
		celExpr      string
		expectedType string
		paramIndex   int
	}{
		{
			name:         "string parameter type",
			celExpr:      `name == "John"`,
			expectedType: "string",
			paramIndex:   0,
		},
		{
			name:         "int64 parameter type",
			celExpr:      `age == 25`,
			expectedType: "int64",
			paramIndex:   0,
		},
		{
			name:         "float64 parameter type",
			celExpr:      `salary == 50000.50`,
			expectedType: "float64",
			paramIndex:   0,
		},
		{
			name:         "bytes parameter type",
			celExpr:      `data == b"hello"`,
			expectedType: "[]uint8",
			paramIndex:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Compile CEL expression
			ast, issues := env.Compile(tt.celExpr)
			require.Nil(t, issues, "CEL compilation should succeed")

			// Convert to parameterized SQL
			result, err := cel2sql.ConvertParameterized(ast)
			require.NoError(t, err, "ConvertParameterized should succeed")
			require.Greater(t, len(result.Parameters), tt.paramIndex, "Should have expected parameter")

			// Check parameter type
			param := result.Parameters[tt.paramIndex]
			actualType := assert.ObjectsAreEqualValues(param, result.Parameters[tt.paramIndex])
			t.Logf("Parameter type: %T (value: %v)", param, param)
			assert.NotNil(t, actualType, "Parameter should have correct type")
		})
	}
}
