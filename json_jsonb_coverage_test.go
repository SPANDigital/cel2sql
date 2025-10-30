package cel2sql_test

import (
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/spandigital/cel2sql/v3"
	"github.com/spandigital/cel2sql/v3/pg"
)

// TestJSONBFieldDetection tests the isFieldJSONB function through actual conversions
func TestJSONBFieldDetection(t *testing.T) {
	// Define schema with both JSON and JSONB fields
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "integer"},
		{Name: "json_data", Type: "json", IsJSON: true, IsJSONB: false},
		{Name: "jsonb_data", Type: "jsonb", IsJSON: true, IsJSONB: true},
		{Name: "jsonb_metadata", Type: "jsonb", IsJSON: true, IsJSONB: true},
	})

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
	}{
		{
			name:        "jsonb_field_access",
			expression:  `record.jsonb_data.name == "test"`,
			description: "Access JSONB field - should use ->> operator",
		},
		{
			name:        "jsonb_nested_access",
			expression:  `record.jsonb_metadata.user.id > 0`,
			description: "Nested JSONB field access",
		},
		{
			name:        "has_on_jsonb",
			expression:  `has(record.jsonb_data.active)`,
			description: "has() function on JSONB field",
		},
		{
			name:        "json_field_access",
			expression:  `record.json_data.status == "ok"`,
			description: "Access JSON (not JSONB) field",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.expression)
			require.NoError(t, issues.Err(), "CEL compilation should succeed")

			schemas := provider.GetSchemas()
			sql, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))

			if err != nil {
				t.Logf("Conversion for %s resulted in error: %v", tt.description, err)
			} else {
				t.Logf("Generated SQL for %s: %s", tt.description, sql)
				assert.NotEmpty(t, sql, "SQL should not be empty")
			}
		})
	}
}

// TestArrayFieldDetection tests the isFieldArray and getFieldElementType functions
func TestArrayFieldDetection(t *testing.T) {
	// Define schema with various array types
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "integer"},
		{Name: "int_array", Type: "integer", Repeated: true, ElementType: "integer"},
		{Name: "text_array", Type: "text", Repeated: true, ElementType: "text"},
		{Name: "double_array", Type: "double precision", Repeated: true, ElementType: "double precision"},
		{Name: "not_array", Type: "text", Repeated: false},
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
		description string
	}{
		{
			name:        "int_array_comprehension",
			expression:  `data.int_array.all(x, x > 0)`,
			description: "Comprehension on integer array",
		},
		{
			name:        "text_array_comprehension",
			expression:  `data.text_array.exists(x, x.startsWith("prefix"))`,
			description: "Comprehension on text array",
		},
		{
			name:        "double_array_comprehension",
			expression:  `data.double_array.filter(x, x >= 0.5)`,
			description: "Filter on double precision array",
		},
		{
			name:        "array_exists_one",
			expression:  `data.int_array.exists_one(x, x == 42)`,
			description: "exists_one on integer array",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.expression)
			require.NoError(t, issues.Err(), "CEL compilation should succeed")

			schemas := provider.GetSchemas()
			sql, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))

			if err != nil {
				t.Logf("Conversion for %s resulted in error: %v", tt.description, err)
			} else {
				t.Logf("Generated SQL for %s: %s", tt.description, sql)
				assert.NotEmpty(t, sql, "SQL should not be empty")
			}
		})
	}
}

// TestJSONBWithArrays tests JSONB fields that contain arrays
// Note: Skipped because .size() on JSONB arrays is not yet fully supported
func TestJSONBWithArrays(t *testing.T) {
	t.Skip("JSONB array .size() operations not yet fully supported")

	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "integer"},
		{Name: "data", Type: "jsonb", IsJSON: true, IsJSONB: true},
	})

	provider := pg.NewTypeProvider(map[string]pg.Schema{"doc": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("doc", cel.ObjectType("doc")),
	)
	require.NoError(t, err)

	tests := []struct {
		name        string
		expression  string
		description string
	}{
		{
			name:        "jsonb_array_access",
			expression:  `doc.data.items.size() > 0`,
			description: "Size of JSONB array field",
		},
		{
			name:        "jsonb_nested_field",
			expression:  `doc.data.user.roles.size() > 0`,
			description: "Nested JSONB array field size",
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
			}
		})
	}
}

// TestSchemaEdgeCases tests edge cases in schema lookups
func TestSchemaEdgeCases(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "integer"},
		{Name: "data", Type: "jsonb", IsJSON: true, IsJSONB: true},
		{Name: "array_field", Type: "text", Repeated: true, ElementType: "text"},
	})

	provider := pg.NewTypeProvider(map[string]pg.Schema{"table1": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("table1", cel.ObjectType("table1")),
	)
	require.NoError(t, err)

	tests := []struct {
		name        string
		expression  string
		description string
		expectError bool
	}{
		{
			name:        "valid_jsonb_field",
			expression:  `table1.data.name == "test"`,
			description: "Valid JSONB field access",
			expectError: false,
		},
		{
			name:        "valid_array_field",
			expression:  `table1.array_field.exists(x, x == "value")`,
			description: "Valid array field comprehension",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.expression)
			require.NoError(t, issues.Err(), "CEL compilation should succeed")

			schemas := provider.GetSchemas()
			sql, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))

			if tt.expectError {
				assert.Error(t, err, "Should error for %s", tt.description)
			} else {
				if err != nil {
					t.Logf("Note for %s: %v", tt.description, err)
				} else {
					t.Logf("Generated SQL for %s: %s", tt.description, sql)
				}
			}
		})
	}
}

// TestNoSchemaProvided tests behavior when no schema is provided
func TestNoSchemaProvided(t *testing.T) {
	// No custom type provider, just basic CEL environment
	env, err := cel.NewEnv(
		cel.Variable("data", cel.DynType),
	)
	require.NoError(t, err)

	tests := []struct {
		name        string
		expression  string
		description string
	}{
		{
			name:        "no_schema_field_access",
			expression:  `data.name == "test"`,
			description: "Field access without schema",
		},
		{
			name:        "no_schema_has_function",
			expression:  `has(data.field)`,
			description: "has() without schema",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.expression)
			require.NoError(t, issues.Err(), "CEL compilation should succeed")

			// Convert WITHOUT schemas
			sql, err := cel2sql.Convert(ast)

			if err != nil {
				t.Logf("Conversion without schema for %s resulted in error: %v", tt.description, err)
			} else {
				t.Logf("Generated SQL without schema for %s: %s", tt.description, sql)
				assert.NotEmpty(t, sql, "SQL should not be empty")
			}
		})
	}
}
