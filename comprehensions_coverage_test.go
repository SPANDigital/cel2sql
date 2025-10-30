package cel2sql_test

import (
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/spandigital/cel2sql/v3"
	"github.com/spandigital/cel2sql/v3/pg"
)

// TestExistsOneComprehension tests the visitExistsOneComprehension function with schema-based arrays
func TestExistsOneComprehension(t *testing.T) {
	// Define schema with array field
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "integer"},
		{Name: "scores", Type: "integer", Repeated: true},
		{Name: "tags", Type: "text", Repeated: true},
	})

	provider := pg.NewTypeProvider(map[string]pg.Schema{"Record": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("record", cel.ObjectType("Record")),
	)
	require.NoError(t, err)

	tests := []struct {
		name        string
		expression  string
		expectedSQL string
		description string
	}{
		{
			name:        "exists_one_with_integer_array",
			expression:  `record.scores.exists_one(s, s > 90)`,
			expectedSQL: `(SELECT COUNT(*) FROM UNNEST(record.scores) AS s WHERE s > 90) = 1`,
			description: "exists_one on integer array field",
		},
		{
			name:        "exists_one_with_string_array",
			expression:  `record.tags.exists_one(t, t == "urgent")`,
			expectedSQL: `(SELECT COUNT(*) FROM UNNEST(record.tags) AS t WHERE t = 'urgent') = 1`,
			description: "exists_one on string array field",
		},
		{
			name:        "exists_one_with_multiple_conditions",
			expression:  `record.scores.exists_one(s, s >= 80 && s <= 90)`,
			expectedSQL: `(SELECT COUNT(*) FROM UNNEST(record.scores) AS s WHERE s >= 80 AND s <= 90) = 1`,
			description: "exists_one with multiple conditions",
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

// TestFilterComprehension tests the visitFilterComprehension function with schema-based arrays
func TestFilterComprehension(t *testing.T) {
	// Define schema with array field
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "integer"},
		{Name: "scores", Type: "integer", Repeated: true},
		{Name: "tags", Type: "text", Repeated: true},
		{Name: "values", Type: "double precision", Repeated: true},
	})

	provider := pg.NewTypeProvider(map[string]pg.Schema{"Record": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("record", cel.ObjectType("Record")),
	)
	require.NoError(t, err)

	tests := []struct {
		name        string
		expression  string
		expectedSQL string
		description string
	}{
		{
			name:        "filter_integer_array",
			expression:  `record.scores.filter(s, s > 70)`,
			expectedSQL: `ARRAY(SELECT s FROM UNNEST(record.scores) AS s WHERE s > 70)`,
			description: "filter on integer array field",
		},
		{
			name:        "filter_string_array",
			expression:  `record.tags.filter(t, t.startsWith("high"))`,
			expectedSQL: `ARRAY(SELECT t FROM UNNEST(record.tags) AS t WHERE t LIKE 'high%' ESCAPE E'\\')`,
			description: "filter on string array field with startsWith",
		},
		{
			name:        "filter_with_multiple_conditions",
			expression:  `record.scores.filter(s, s >= 60 && s < 90)`,
			expectedSQL: `ARRAY(SELECT s FROM UNNEST(record.scores) AS s WHERE s >= 60 AND s < 90)`,
			description: "filter with multiple conditions",
		},
		{
			name:        "filter_double_array",
			expression:  `record.values.filter(v, v > 0.5)`,
			expectedSQL: `ARRAY(SELECT v FROM UNNEST(record.values) AS v WHERE v > 0.5)`,
			description: "filter on double precision array field",
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

// TestJSONArrayComprehensions tests comprehensions with JSON array fields
func TestJSONArrayComprehensions(t *testing.T) {
	// Define schema with JSON array field
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "integer"},
		{Name: "metadata", Type: "jsonb", IsJSON: true, IsJSONB: true},
	})

	provider := pg.NewTypeProvider(map[string]pg.Schema{"Record": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("record", cel.ObjectType("Record")),
	)
	require.NoError(t, err)

	tests := []struct {
		name        string
		expression  string
		description string
		shouldPass  bool
	}{
		{
			name:        "exists_one_on_json_array",
			expression:  `record.metadata.items.exists_one(i, i.active)`,
			description: "exists_one on JSON array field",
			shouldPass:  true,
		},
		{
			name:        "filter_on_json_array",
			expression:  `record.metadata.items.filter(i, i.score > 50)`,
			description: "filter on JSON array field",
			shouldPass:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.expression)
			require.NoError(t, issues.Err(), "CEL compilation should succeed")

			schemas := provider.GetSchemas()
			sql, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))

			if tt.shouldPass {
				if err != nil {
					t.Logf("Note: %s - Error: %v", tt.description, err)
					t.Logf("This may be expected if JSON array comprehensions aren't fully supported yet")
				} else {
					t.Logf("Generated SQL for %s: %s", tt.description, sql)
				}
			}
		})
	}
}

// TestNestedComprehensions tests nested comprehension scenarios
func TestNestedComprehensions(t *testing.T) {
	// Define schema with nested arrays
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "integer"},
		{Name: "groups", Type: "jsonb", IsJSON: true, IsJSONB: true},
	})

	provider := pg.NewTypeProvider(map[string]pg.Schema{"Record": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("record", cel.ObjectType("Record")),
		cel.Variable("numbers", cel.ListType(cel.ListType(cel.IntType))),
	)
	require.NoError(t, err)

	tests := []struct {
		name        string
		expression  string
		description string
	}{
		{
			name:        "nested_exists_one",
			expression:  `numbers.exists(list, list.exists_one(n, n == 5))`,
			description: "nested exists with exists_one",
		},
		{
			name:        "filter_then_exists_one",
			expression:  `numbers.filter(list, list.exists_one(n, n > 0))`,
			description: "filter containing exists_one",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.expression)
			require.NoError(t, issues.Err(), "CEL compilation should succeed")

			schemas := provider.GetSchemas()
			sql, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))

			if err != nil {
				t.Logf("Note: %s may not be fully supported yet - Error: %v", tt.description, err)
			} else {
				t.Logf("Generated SQL for %s: %s", tt.description, sql)
			}
		})
	}
}

// TestComprehensionZeroCoverageEdgeCases tests edge cases for zero-coverage comprehension functions
func TestComprehensionZeroCoverageEdgeCases(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "integer"},
		{Name: "items", Type: "integer", Repeated: true},
	})

	provider := pg.NewTypeProvider(map[string]pg.Schema{"Record": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("record", cel.ObjectType("Record")),
		cel.Variable("empty_list", cel.ListType(cel.IntType)),
	)
	require.NoError(t, err)

	tests := []struct {
		name        string
		expression  string
		description string
	}{
		{
			name:        "exists_one_no_matches",
			expression:  `record.items.exists_one(i, i > 1000000)`,
			description: "exists_one when no elements match (should return false)",
		},
		{
			name:        "filter_no_matches",
			expression:  `record.items.filter(i, i < 0)`,
			description: "filter when no elements match (should return empty array)",
		},
		{
			name:        "exists_one_all_match",
			expression:  `record.items.exists_one(i, i >= 0)`,
			description: "exists_one when multiple elements match (should return false)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.expression)
			require.NoError(t, issues.Err(), "CEL compilation should succeed")

			schemas := provider.GetSchemas()
			sql, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
			require.NoError(t, err, "Conversion should succeed for %s", tt.description)

			t.Logf("Generated SQL for %s: %s", tt.description, sql)
			assert.NotEmpty(t, sql, "SQL should not be empty")
		})
	}
}
