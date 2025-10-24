package cel2sql_test

import (
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/spandigital/cel2sql/v2"
	"github.com/spandigital/cel2sql/v2/pg"
)

// TestConstEdgeCases tests visitConst with edge case values
func TestConstEdgeCases(t *testing.T) {
	env, err := cel.NewEnv()
	require.NoError(t, err)

	tests := []struct {
		name        string
		expression  string
		expectedSQL string
		description string
	}{
		{
			name:        "max_int64",
			expression:  "9223372036854775807 > 0",
			expectedSQL: "9223372036854775807 > 0",
			description: "Maximum int64 value",
		},
		{
			name:        "min_int64",
			expression:  "-9223372036854775808 < 0",
			expectedSQL: "-9223372036854775808 < 0",
			description: "Minimum int64 value",
		},
		{
			name:        "zero_integer",
			expression:  "0 == 0",
			expectedSQL: "0 = 0",
			description: "Zero integer value",
		},
		{
			name:        "zero_double",
			expression:  "0.0 == 0.0",
			expectedSQL: "0 = 0",
			description: "Zero double value",
		},
		{
			name:        "negative_zero",
			expression:  "-0.0 < 1.0",
			expectedSQL: "-0 < 1",
			description: "Negative zero",
		},
		{
			name:        "very_large_double",
			expression:  "1.7976931348623157e308 > 0.0",
			expectedSQL: "1.7976931348623157e+308 > 0",
			description: "Very large double near MAX_DOUBLE",
		},
		{
			name:        "very_small_positive_double",
			expression:  "2.2250738585072014e-308 > 0.0",
			expectedSQL: "2.2250738585072014e-308 > 0",
			description: "Very small positive double",
		},
		{
			name:        "negative_double",
			expression:  "-123.456 < 0.0",
			expectedSQL: "-123.456 < 0",
			description: "Negative double",
		},
		{
			name:        "empty_string",
			expression:  `"" == ""`,
			expectedSQL: `'' = ''`,
			description: "Empty string constant",
		},
		{
			name:        "string_with_quotes",
			expression:  `"it's working" == "it's working"`,
			expectedSQL: `'it''s working' = 'it''s working'`,
			description: "String with single quotes (escaped)",
		},
		{
			name:        "string_with_backslash",
			expression:  `"path\\to\\file" != ""`,
			expectedSQL: `'path\to\file' != ''`,
			description: "String with backslashes",
		},
		{
			name:        "unicode_string",
			expression:  `"Hello 世界 🌍" != ""`,
			expectedSQL: `'Hello 世界 🌍' != ''`,
			description: "Unicode string with emojis",
		},
		{
			name:        "null_value",
			expression:  "null == null",
			expectedSQL: "NULL IS NULL",
			description: "NULL constant",
		},
		{
			name:        "true_bool",
			expression:  "true == true",
			expectedSQL: "TRUE IS TRUE",
			description: "Boolean true",
		},
		{
			name:        "false_bool",
			expression:  "false == false",
			expectedSQL: "FALSE IS FALSE",
			description: "Boolean false",
		},
		{
			name:        "bytes_value",
			expression:  `b"hello" != b""`,
			expectedSQL: `'\x68656c6c6f' != '\x'`,
			description: "Bytes constant",
		},
		{
			name:        "uint64_value",
			expression:  "18446744073709551615u > 0u",
			expectedSQL: "18446744073709551615 > 0",
			description: "Maximum uint64 value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.expression)
			require.NoError(t, issues.Err(), "CEL compilation should succeed")

			sql, err := cel2sql.Convert(ast)
			require.NoError(t, err, "Conversion should succeed for %s", tt.description)
			assert.Equal(t, tt.expectedSQL, sql, "SQL should match for %s", tt.description)
			t.Logf("Generated SQL: %s", sql)
		})
	}
}

// TestStringFunctionEdgeCases tests callContains, callStartsWith, callEndsWith with edge cases
func TestStringFunctionEdgeCases(t *testing.T) {
	env, err := cel.NewEnv(
		cel.Variable("text", cel.StringType),
		cel.Variable("empty", cel.StringType),
	)
	require.NoError(t, err)

	tests := []struct {
		name        string
		expression  string
		expectedSQL string
		description string
	}{
		{
			name:        "contains_empty_string",
			expression:  `text.contains("")`,
			expectedSQL: `POSITION('' IN text) > 0`,
			description: "Contains empty string (always true)",
		},
		{
			name:        "contains_special_chars",
			expression:  `text.contains("100%")`,
			expectedSQL: `POSITION('100%' IN text) > 0`,
			description: "Contains string with percent sign",
		},
		{
			name:        "contains_with_underscore",
			expression:  `text.contains("_test")`,
			expectedSQL: `POSITION('_test' IN text) > 0`,
			description: "Contains string with underscore",
		},
		{
			name:        "startsWith_empty",
			expression:  `text.startsWith("")`,
			expectedSQL: `text LIKE '' || '%'`,
			description: "StartsWith empty string",
		},
		{
			name:        "startsWith_percent",
			expression:  `text.startsWith("50%")`,
			expectedSQL: `text LIKE '50\%%'`,
			description: "StartsWith string containing % (LIKE special char)",
		},
		{
			name:        "startsWith_underscore",
			expression:  `text.startsWith("_prefix")`,
			expectedSQL: `text LIKE '\_prefix%'`,
			description: "StartsWith string containing _ (LIKE special char)",
		},
		{
			name:        "startsWith_backslash",
			expression:  `text.startsWith("\\path")`,
			expectedSQL: `text LIKE '\\path%'`,
			description: "StartsWith string containing backslash",
		},
		{
			name:        "endsWith_empty",
			expression:  `text.endsWith("")`,
			expectedSQL: `text LIKE '%' || ''`,
			description: "EndsWith empty string",
		},
		{
			name:        "endsWith_percent",
			expression:  `text.endsWith("100%")`,
			expectedSQL: `text LIKE '%100\%'`,
			description: "EndsWith string containing %",
		},
		{
			name:        "endsWith_underscore",
			expression:  `text.endsWith("suffix_")`,
			expectedSQL: `text LIKE '%suffix\_'`,
			description: "EndsWith string containing _",
		},
		{
			name:        "contains_unicode",
			expression:  `text.contains("世界")`,
			expectedSQL: `POSITION('世界' IN text) > 0`,
			description: "Contains Unicode characters",
		},
		{
			name:        "multiple_contains",
			expression:  `text.contains("foo") && text.contains("bar")`,
			expectedSQL: `POSITION('foo' IN text) > 0 AND POSITION('bar' IN text) > 0`,
			description: "Multiple contains in AND",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.expression)
			require.NoError(t, issues.Err(), "CEL compilation should succeed")

			sql, err := cel2sql.Convert(ast)
			require.NoError(t, err, "Conversion should succeed for %s", tt.description)
			assert.Equal(t, tt.expectedSQL, sql, "SQL should match for %s", tt.description)
			t.Logf("Generated SQL: %s", sql)
		})
	}
}

// TestBinaryOperatorEdgeCases tests visitCallBinary with various operator combinations
func TestBinaryOperatorEdgeCases(t *testing.T) {
	schema := pg.Schema{
		{Name: "id", Type: "integer"},
		{Name: "name", Type: "text"},
		{Name: "score", Type: "double precision"},
		{Name: "tags", Type: "text", Repeated: true, ElementType: "text"},
		{Name: "metadata", Type: "jsonb", IsJSON: true, IsJSONB: true},
	}

	provider := pg.NewTypeProvider(map[string]pg.Schema{"data": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("data", cel.ObjectType("data")),
		cel.Variable("a", cel.IntType),
		cel.Variable("b", cel.IntType),
	)
	require.NoError(t, err)

	tests := []struct {
		name        string
		expression  string
		expectedSQL string
		description string
	}{
		{
			name:        "bool_comparison_equals",
			expression:  "true == false",
			expectedSQL: "TRUE IS FALSE",
			description: "Bool comparison with == becomes IS",
		},
		{
			name:        "bool_comparison_not_equals",
			expression:  "true != false",
			expectedSQL: "TRUE IS NOT FALSE",
			description: "Bool comparison with != becomes IS NOT",
		},
		{
			name:        "string_concatenation",
			expression:  `"hello" + " world"`,
			expectedSQL: `'hello' || ' world'`,
			description: "String concatenation with +",
		},
		{
			name:        "list_concatenation",
			expression:  "[1, 2] + [3, 4]",
			expectedSQL: "ARRAY[1, 2] || ARRAY[3, 4]",
			description: "List concatenation",
		},
		{
			name:        "modulo_operator",
			expression:  "a % b == 0",
			expectedSQL: "MOD(a, b) = 0",
			description: "Modulo operator",
		},
		{
			name:        "division_operator",
			expression:  "a / b > 0",
			expectedSQL: "a / b > 0",
			description: "Division operator",
		},
		{
			name:        "logical_or",
			expression:  "a > 10 || b < 5",
			expectedSQL: "a > 10 OR b < 5",
			description: "Logical OR",
		},
		{
			name:        "logical_and",
			expression:  "a > 10 && b < 5",
			expectedSQL: "a > 10 AND b < 5",
			description: "Logical AND",
		},
		{
			name:        "nested_operators",
			expression:  "(a + b) * (a - b)",
			expectedSQL: "(a + b) * (a - b)",
			description: "Nested arithmetic operators",
		},
		{
			name:        "comparison_chain",
			expression:  "a < b && b < 100",
			expectedSQL: "a < b AND b < 100",
			description: "Comparison chain",
		},
		{
			name:        "in_operator_with_list",
			expression:  `data.name in ["admin", "user"]`,
			expectedSQL: `data.name = ANY(ARRAY['admin', 'user'])`,
			description: "IN operator with list literal",
		},
		{
			name:        "in_operator_with_array_field",
			expression:  `"test" in data.tags`,
			expectedSQL: `'test' = ANY(data.tags)`,
			description: "IN operator with array field",
		},
		{
			name:        "negation_operator",
			expression:  "!(a > 10)",
			expectedSQL: "NOT (a > 10)",
			description: "Negation operator",
		},
		{
			name:        "less_than_or_equal",
			expression:  "a <= b",
			expectedSQL: "a <= b",
			description: "Less than or equal",
		},
		{
			name:        "greater_than_or_equal",
			expression:  "a >= b",
			expectedSQL: "a >= b",
			description: "Greater than or equal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.expression)
			require.NoError(t, issues.Err(), "CEL compilation should succeed")

			schemas := provider.GetSchemas()
			sql, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
			require.NoError(t, err, "Conversion should succeed for %s", tt.description)
			assert.Equal(t, tt.expectedSQL, sql, "SQL should match for %s", tt.description)
			t.Logf("Generated SQL: %s", sql)
		})
	}
}

// TestIdentifierEdgeCases tests visitIdent with various identifier patterns
func TestIdentifierEdgeCases(t *testing.T) {
	schema := pg.Schema{
		{Name: "id", Type: "integer"},
		{Name: "score", Type: "integer"},
		{Name: "value", Type: "integer"},
		{Name: "count", Type: "integer"},
		{Name: "amount", Type: "integer"},
	}

	provider := pg.NewTypeProvider(map[string]pg.Schema{"record": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("record", cel.ObjectType("record")),
	)
	require.NoError(t, err)

	tests := []struct {
		name        string
		expression  string
		description string
		checkSQL    func(t *testing.T, sql string)
	}{
		{
			name:        "identifier_score",
			expression:  "record.score > 10",
			description: "Identifier 'score' (may need numeric casting in JSON context)",
			checkSQL: func(t *testing.T, sql string) {
				// Should be simple field access
				assert.Contains(t, sql, "record.score")
			},
		},
		{
			name:        "identifier_value",
			expression:  "record.value > 0",
			description: "Identifier 'value' (may need numeric casting in JSON context)",
			checkSQL: func(t *testing.T, sql string) {
				assert.Contains(t, sql, "record.value")
			},
		},
		{
			name:        "identifier_count",
			expression:  "record.count == 5",
			description: "Identifier 'count'",
			checkSQL: func(t *testing.T, sql string) {
				assert.Contains(t, sql, "record.count")
			},
		},
		{
			name:        "identifier_amount",
			expression:  "record.amount >= 100",
			description: "Identifier 'amount'",
			checkSQL: func(t *testing.T, sql string) {
				assert.Contains(t, sql, "record.amount")
			},
		},
		{
			name:        "simple_identifier",
			expression:  "record.id > 0",
			description: "Simple identifier access",
			checkSQL: func(t *testing.T, sql string) {
				assert.Equal(t, "record.id > 0", sql)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.expression)
			require.NoError(t, issues.Err(), "CEL compilation should succeed")

			schemas := provider.GetSchemas()
			sql, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
			require.NoError(t, err, "Conversion should succeed for %s", tt.description)

			tt.checkSQL(t, sql)
			t.Logf("Generated SQL for %s: %s", tt.description, sql)
		})
	}
}

// TestCallFuncEdgeCases tests visitCallFunc with various function types
func TestCallFuncEdgeCases(t *testing.T) {
	schema := pg.Schema{
		{Name: "id", Type: "integer"},
		{Name: "name", Type: "text"},
		{Name: "created", Type: "timestamp with time zone"},
		{Name: "data", Type: "jsonb", IsJSON: true, IsJSONB: true},
		{Name: "tags", Type: "text", Repeated: true, ElementType: "text"},
		{Name: "scores", Type: "integer", Repeated: true, ElementType: "integer"},
	}

	provider := pg.NewTypeProvider(map[string]pg.Schema{"record": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("record", cel.ObjectType("record")),
		cel.Variable("text", cel.StringType),
	)
	require.NoError(t, err)

	tests := []struct {
		name        string
		expression  string
		description string
		checkSQL    func(t *testing.T, sql string)
	}{
		{
			name:        "type_conversion_int",
			expression:  `int("123") > 100`,
			description: "int() type conversion",
			checkSQL: func(t *testing.T, sql string) {
				assert.Contains(t, sql, "CAST")
				assert.Contains(t, sql, "BIGINT")
			},
		},
		{
			name:        "type_conversion_double",
			expression:  `double("3.14") > 3.0`,
			description: "double() type conversion",
			checkSQL: func(t *testing.T, sql string) {
				assert.Contains(t, sql, "CAST")
				assert.Contains(t, sql, "DOUBLE PRECISION")
			},
		},
		{
			name:        "type_conversion_string",
			expression:  "string(123) == \"123\"",
			description: "string() type conversion",
			checkSQL: func(t *testing.T, sql string) {
				assert.Contains(t, sql, "CAST")
				assert.Contains(t, sql, "TEXT")
			},
		},
		{
			name:        "type_conversion_bool",
			expression:  `bool("true") == true`,
			description: "bool() type conversion",
			checkSQL: func(t *testing.T, sql string) {
				assert.Contains(t, sql, "CAST")
				assert.Contains(t, sql, "BOOLEAN")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.expression)
			require.NoError(t, issues.Err(), "CEL compilation should succeed")

			schemas := provider.GetSchemas()
			sql, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
			require.NoError(t, err, "Conversion should succeed for %s", tt.description)

			tt.checkSQL(t, sql)
			t.Logf("Generated SQL for %s: %s", tt.description, sql)
		})
	}
}

// TestSpecialDoubleValues tests special double values (Inf, NaN)
func TestSpecialDoubleValues(t *testing.T) {
	// Note: CEL doesn't directly support Inf/NaN literals, but we can test with very large numbers
	env, err := cel.NewEnv(
		cel.Variable("x", cel.DoubleType),
	)
	require.NoError(t, err)

	tests := []struct {
		name        string
		expression  string
		expectedSQL string
		description string
	}{
		{
			name:        "positive_infinity_comparison",
			expression:  "x < 1.7976931348623157e308",
			expectedSQL: "x < 1.7976931348623157e+308",
			description: "Comparison with near-infinity value",
		},
		{
			name:        "negative_infinity_comparison",
			expression:  "x > -1.7976931348623157e308",
			expectedSQL: "x > -1.7976931348623157e+308",
			description: "Comparison with near-negative-infinity value",
		},
		{
			name:        "scientific_notation_small",
			expression:  "x > 1.23e-100",
			expectedSQL: "x > 1.23e-100",
			description: "Small scientific notation",
		},
		{
			name:        "scientific_notation_large",
			expression:  "x < 9.87e+200",
			expectedSQL: "x < 9.87e+200",
			description: "Large scientific notation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.expression)
			require.NoError(t, issues.Err(), "CEL compilation should succeed")

			sql, err := cel2sql.Convert(ast)
			require.NoError(t, err, "Conversion should succeed for %s", tt.description)
			assert.Equal(t, tt.expectedSQL, sql, "SQL should match for %s", tt.description)
			t.Logf("Generated SQL: %s", sql)
		})
	}
}

// TestComplexNestedExpressions tests complex nested expressions to improve coverage
func TestComplexNestedExpressions(t *testing.T) {
	env, err := cel.NewEnv(
		cel.Variable("a", cel.IntType),
		cel.Variable("b", cel.IntType),
		cel.Variable("c", cel.IntType),
		cel.Variable("text", cel.StringType),
	)
	require.NoError(t, err)

	tests := []struct {
		name        string
		expression  string
		description string
	}{
		{
			name:        "deeply_nested_arithmetic",
			expression:  "((a + b) * (c - a)) / ((b * c) + 1)",
			description: "Deeply nested arithmetic operations",
		},
		{
			name:        "complex_boolean_logic",
			expression:  "(a > 10 && b < 20) || (c == 5 && a != b)",
			description: "Complex boolean logic with mixed operators",
		},
		{
			name:        "mixed_string_and_numeric",
			expression:  `(text.contains("hello") && a > 0) || (text.contains("test") && b < 100)`,
			description: "Mixed string and numeric operations",
		},
		{
			name:        "nested_ternary_style",
			expression:  "(a > b ? a : b) > c",
			description: "Ternary-style expression with comparison",
		},
		{
			name:        "parenthesized_expressions",
			expression:  "((a + b) > (c * 2)) && ((a - c) < (b + 1))",
			description: "Multiple levels of parenthesized expressions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.expression)
			require.NoError(t, issues.Err(), "CEL compilation should succeed")

			sql, err := cel2sql.Convert(ast)
			require.NoError(t, err, "Conversion should succeed for %s", tt.description)
			assert.NotEmpty(t, sql, "SQL should not be empty")
			t.Logf("Generated SQL for %s: %s", tt.description, sql)
		})
	}
}
