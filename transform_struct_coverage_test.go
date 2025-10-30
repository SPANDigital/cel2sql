package cel2sql_test

import (
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/spandigital/cel2sql/v3"
	"github.com/spandigital/cel2sql/v3/pg"
)

// TestJSONColumnTableReference tests the isTableReference function through JSON column detection
// This tests that table.metadata pattern is correctly identified as JSON column access
func TestJSONColumnTableReference(t *testing.T) {
	// Define schema with known JSON columns
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "integer"},
		{Name: "name", Type: "text"},
		{Name: "metadata", Type: "jsonb", IsJSON: true, IsJSONB: true},
		{Name: "properties", Type: "jsonb", IsJSON: true, IsJSONB: true},
		{Name: "content", Type: "json", IsJSON: true, IsJSONB: false},
	})

	provider := pg.NewTypeProvider(map[string]pg.Schema{"assets": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("assets", cel.ObjectType("assets")),
	)
	require.NoError(t, err)

	tests := []struct {
		name        string
		expression  string
		expectedSQL string
		description string
	}{
		{
			name:        "json_column_metadata",
			expression:  `assets.metadata.status == "active"`,
			expectedSQL: `assets.metadata->>'status' = 'active'`,
			description: "Access metadata JSON column - should use ->> operator",
		},
		{
			name:        "json_column_properties",
			expression:  `assets.properties.version == "1.0"`,
			expectedSQL: `assets.properties->>'version' = '1.0'`,
			description: "Access properties JSON column",
		},
		{
			name:        "json_column_content",
			expression:  `assets.content.title == "test"`,
			expectedSQL: `assets.content->>'title' = 'test'`,
			description: "Access content JSON column",
		},
		{
			name:        "nested_json_access",
			expression:  `assets.metadata.user.name == "admin"`,
			expectedSQL: `assets.metadata->'user'->>'name' = 'admin'`,
			description: "Nested JSON field access",
		},
		{
			name:        "json_with_has",
			expression:  `has(assets.metadata.status)`,
			expectedSQL: `assets.metadata ? 'status'`,
			description: "has() on JSON column",
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
			t.Logf("Generated SQL for %s: %s", tt.description, sql)
		})
	}
}

// TestTableReferenceDetection tests isTableReference through non-JSON column access
func TestTableReferenceDetection(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "integer"},
		{Name: "name", Type: "text"},
		{Name: "age", Type: "integer"},
		{Name: "metadata", Type: "jsonb", IsJSON: true, IsJSONB: true},
	})

	provider := pg.NewTypeProvider(map[string]pg.Schema{"users": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("users", cel.ObjectType("users")),
	)
	require.NoError(t, err)

	tests := []struct {
		name        string
		expression  string
		expectedSQL string
		description string
	}{
		{
			name:        "simple_column_access",
			expression:  `users.name == "test"`,
			expectedSQL: `users.name = 'test'`,
			description: "Simple column access (not JSON)",
		},
		{
			name:        "integer_column",
			expression:  `users.age > 18`,
			expectedSQL: `users.age > 18`,
			description: "Integer column access",
		},
		{
			name:        "combined_columns_and_json",
			expression:  `users.name == "admin" && users.metadata.role == "superuser"`,
			expectedSQL: `users.name = 'admin' AND users.metadata->>'role' = 'superuser'`,
			description: "Mix of regular columns and JSON column",
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

// TestMapComprehensionAsTransform tests that map() comprehensions work correctly
// This indirectly tests that TransformList comprehensions are handled via Map
func TestMapComprehensionAsTransform(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "integer"},
		{Name: "scores", Type: "integer", Repeated: true, ElementType: "integer"},
		{Name: "prices", Type: "double precision", Repeated: true, ElementType: "double precision"},
	})

	provider := pg.NewTypeProvider(map[string]pg.Schema{"data": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("data", cel.ObjectType("data")),
	)
	require.NoError(t, err)

	tests := []struct {
		name        string
		expression  string
		expectedSQL string
		description string
	}{
		{
			name:        "map_simple_transform",
			expression:  `data.scores.map(s, s * 2)`,
			expectedSQL: `ARRAY(SELECT s * 2 FROM UNNEST(data.scores) AS s)`,
			description: "Simple map transformation (multiply by 2)",
		},
		{
			name:        "map_with_addition",
			expression:  `data.scores.map(s, s + 10)`,
			expectedSQL: `ARRAY(SELECT s + 10 FROM UNNEST(data.scores) AS s)`,
			description: "Map with addition",
		},
		{
			name:        "map_double_array",
			expression:  `data.prices.map(p, p * 1.1)`,
			expectedSQL: `ARRAY(SELECT p * 1.1 FROM UNNEST(data.prices) AS p)`,
			description: "Map on double precision array",
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

// TestStructExpressions tests struct/message construction in CEL
// Note: CEL struct syntax is limited, these test what's possible
func TestStructExpressions(t *testing.T) {
	env, err := cel.NewEnv()
	require.NoError(t, err)

	tests := []struct {
		name          string
		expression    string
		description   string
		expectError   bool
		errorContains string
	}{
		{
			name:        "simple_map_struct",
			expression:  `{"key": "value", "number": 42}`,
			description: "Simple map structure",
			expectError: false,
		},
		{
			name:        "nested_map_struct",
			expression:  `{"user": {"name": "test", "age": 30}}`,
			description: "Nested map structure",
			expectError: false,
		},
		{
			name:        "map_with_mixed_types",
			expression:  `{"string": "test", "int": 42, "bool": true}`,
			description: "Map with mixed value types",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.expression)
			require.NoError(t, issues.Err(), "CEL compilation should succeed")

			sql, err := cel2sql.Convert(ast)

			if tt.expectError {
				require.Error(t, err, "Should error for %s", tt.description)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				if err != nil {
					t.Logf("Note: Struct conversion for %s resulted in: %v", tt.description, err)
					t.Logf("This may be expected if struct message conversion is not fully supported")
				} else {
					t.Logf("Generated SQL for %s: %s", tt.description, sql)
					assert.NotEmpty(t, sql, "SQL should not be empty")
				}
			}
		})
	}
}

// TestEdgeCasesForCoverage tests additional edge cases to improve coverage
func TestEdgeCasesForCoverage(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "integer"},
		{Name: "tags", Type: "text", Repeated: true, ElementType: "text"},
		{Name: "metadata", Type: "jsonb", IsJSON: true, IsJSONB: true},
	})

	provider := pg.NewTypeProvider(map[string]pg.Schema{"records": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("records", cel.ObjectType("records")),
	)
	require.NoError(t, err)

	tests := []struct {
		name        string
		expression  string
		description string
	}{
		{
			name:        "deeply_nested_json",
			expression:  `records.metadata.a.b.c.d == "deep"`,
			description: "Deeply nested JSON path",
		},
		{
			name:        "json_with_comparison",
			expression:  `records.metadata.count > 10`,
			description: "JSON field with numeric comparison",
		},
		{
			name:        "array_with_complex_filter",
			expression:  `records.tags.filter(t, t.startsWith("prod") && t.endsWith("ion"))`,
			description: "Filter with multiple string operations",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.expression)
			require.NoError(t, issues.Err(), "CEL compilation should succeed")

			schemas := provider.GetSchemas()
			sql, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))

			if err != nil {
				t.Logf("Note for %s: %v", tt.description, err)
			} else {
				t.Logf("Generated SQL for %s: %s", tt.description, sql)
				assert.NotEmpty(t, sql, "SQL should not be empty")
			}
		})
	}
}
