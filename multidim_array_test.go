package cel2sql

import (
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/spandigital/cel2sql/v3/pg"
)

// TestMultiDimensionalArrays tests size() function with multi-dimensional arrays
func TestMultiDimensionalArrays(t *testing.T) {
	// Create schema with arrays of different dimensions
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "tags", Type: "text", Repeated: true, Dimensions: 1},         // 1D array
		{Name: "matrix", Type: "integer", Repeated: true, Dimensions: 2},    // 2D array
		{Name: "cube", Type: "integer", Repeated: true, Dimensions: 3},      // 3D array
		{Name: "hypercube", Type: "integer", Repeated: true, Dimensions: 4}, // 4D array
		{Name: "simple_array", Type: "text", Repeated: true, Dimensions: 0}, // Dimension not set (defaults to 1)
		{Name: "id", Type: "integer", Repeated: false, Dimensions: 0},       // Not an array
	})

	provider := pg.NewTypeProvider(map[string]pg.Schema{"TestTable": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("data", cel.ObjectType("TestTable")),
	)
	if err != nil {
		t.Fatalf("failed to create CEL environment: %v", err)
	}

	tests := []struct {
		name        string
		cel         string
		expectedSQL string
		expectedDim int
		description string
	}{
		{
			name:        "1D array size",
			cel:         "size(data.tags) > 0",
			expectedSQL: "COALESCE(ARRAY_LENGTH(data.tags, 1), 0) > 0",
			expectedDim: 1,
			description: "1D array should use dimension 1",
		},
		{
			name:        "2D array size",
			cel:         "size(data.matrix) > 0",
			expectedSQL: "COALESCE(ARRAY_LENGTH(data.matrix, 2), 0) > 0",
			expectedDim: 2,
			description: "2D array should use dimension 2",
		},
		{
			name:        "3D array size",
			cel:         "size(data.cube) == 5",
			expectedSQL: "COALESCE(ARRAY_LENGTH(data.cube, 3), 0) = 5",
			expectedDim: 3,
			description: "3D array should use dimension 3",
		},
		{
			name:        "4D array size",
			cel:         "size(data.hypercube) < 10",
			expectedSQL: "COALESCE(ARRAY_LENGTH(data.hypercube, 4), 0) < 10",
			expectedDim: 4,
			description: "4D array should use dimension 4",
		},
		{
			name:        "array with no dimension set defaults to 1",
			cel:         "size(data.simple_array) > 0",
			expectedSQL: "COALESCE(ARRAY_LENGTH(data.simple_array, 1), 0) > 0",
			expectedDim: 1,
			description: "Array with Dimensions=0 should default to dimension 1",
		},
		{
			name:        "complex expression with 2D array",
			cel:         "size(data.matrix) >= 3 && size(data.matrix) <= 10",
			expectedSQL: "COALESCE(ARRAY_LENGTH(data.matrix, 2), 0) >= 3 AND COALESCE(ARRAY_LENGTH(data.matrix, 2), 0) <= 10",
			expectedDim: 2,
			description: "Complex expressions should maintain correct dimension",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.cel)
			if issues != nil && issues.Err() != nil {
				t.Fatalf("failed to compile CEL: %v", issues.Err())
			}

			sql, err := Convert(ast, WithSchemas(map[string]pg.Schema{"TestTable": schema}))
			if err != nil {
				t.Fatalf("failed to convert to SQL: %v", err)
			}

			if sql != tt.expectedSQL {
				t.Errorf("unexpected SQL for %s:\n  got:  %s\n  want: %s", tt.description, sql, tt.expectedSQL)
			}

			t.Logf("✓ %s: ARRAY_LENGTH(..., %d)", tt.description, tt.expectedDim)
		})
	}
}

// TestMultiDimensionalArrayBackwardCompatibility ensures that existing 1D array behavior is preserved
func TestMultiDimensionalArrayBackwardCompatibility(t *testing.T) {
	tests := []struct {
		name        string
		cel         string
		expectedSQL string
	}{
		{
			name:        "size without schema - defaults to 1D",
			cel:         "size(string_list) > 0",
			expectedSQL: "COALESCE(ARRAY_LENGTH(string_list, 1), 0) > 0",
		},
		{
			name:        "size in complex expression - defaults to 1D",
			cel:         "size(items) >= 2 && size(items) <= 5",
			expectedSQL: "COALESCE(ARRAY_LENGTH(items, 1), 0) >= 2 AND COALESCE(ARRAY_LENGTH(items, 1), 0) <= 5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create environment with variables declared but no schema
			env, err := cel.NewEnv(
				cel.Variable("string_list", cel.ListType(cel.StringType)),
				cel.Variable("items", cel.ListType(cel.IntType)),
			)
			if err != nil {
				t.Fatalf("failed to create CEL environment: %v", err)
			}

			ast, issues := env.Compile(tt.cel)
			if issues != nil && issues.Err() != nil {
				t.Fatalf("failed to compile CEL: %v", issues.Err())
			}

			// Convert WITHOUT providing schema - should default to dimension 1
			sql, err := Convert(ast)
			if err != nil {
				t.Fatalf("failed to convert to SQL: %v", err)
			}

			if sql != tt.expectedSQL {
				t.Errorf("backward compatibility failed:\n  got:  %s\n  want: %s", sql, tt.expectedSQL)
			}

			t.Logf("✓ Backward compatible: defaults to dimension 1 when no schema")
		})
	}
}

// TestExplicitDimensionOverridesDefault tests that explicitly setting Dimensions to 1 works correctly
func TestExplicitDimensionOverridesDefault(t *testing.T) {
	// Schema with explicitly set Dimensions=1
	schemaExplicit := pg.NewSchema([]pg.FieldSchema{
		{Name: "tags", Type: "text", Repeated: true, Dimensions: 1},
	})

	// Schema with Dimensions=0 (not set, should use 1 as default)
	schemaImplicit := pg.NewSchema([]pg.FieldSchema{
		{Name: "tags", Type: "text", Repeated: true, Dimensions: 0},
	})

	provider1 := pg.NewTypeProvider(map[string]pg.Schema{"TestTable": schemaExplicit})
	provider2 := pg.NewTypeProvider(map[string]pg.Schema{"TestTable": schemaImplicit})

	env1, _ := cel.NewEnv(
		cel.CustomTypeProvider(provider1),
		cel.Variable("data", cel.ObjectType("TestTable")),
	)

	env2, _ := cel.NewEnv(
		cel.CustomTypeProvider(provider2),
		cel.Variable("data", cel.ObjectType("TestTable")),
	)

	celExpr := "size(data.tags) > 0"
	expectedSQL := "COALESCE(ARRAY_LENGTH(data.tags, 1), 0) > 0"

	// Test with explicit Dimensions=1
	ast1, _ := env1.Compile(celExpr)
	sql1, err := Convert(ast1, WithSchemas(map[string]pg.Schema{"TestTable": schemaExplicit}))
	if err != nil {
		t.Fatalf("failed with explicit dimension: %v", err)
	}
	if sql1 != expectedSQL {
		t.Errorf("explicit Dimensions=1 failed:\n  got:  %s\n  want: %s", sql1, expectedSQL)
	}

	// Test with implicit Dimensions=0 (should default to 1)
	ast2, _ := env2.Compile(celExpr)
	sql2, err := Convert(ast2, WithSchemas(map[string]pg.Schema{"TestTable": schemaImplicit}))
	if err != nil {
		t.Fatalf("failed with implicit dimension: %v", err)
	}
	if sql2 != expectedSQL {
		t.Errorf("implicit Dimensions=0 (default) failed:\n  got:  %s\n  want: %s", sql2, expectedSQL)
	}

	t.Log("✓ Both explicit Dimensions=1 and implicit Dimensions=0 default to dimension 1")
}
