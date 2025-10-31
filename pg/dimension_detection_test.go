package pg

import "testing"

// TestDetectArrayDimensions tests the detectArrayDimensions function with various PostgreSQL type strings
func TestDetectArrayDimensions(t *testing.T) {
	tests := []struct {
		name     string
		pgType   string
		expected int
	}{
		// Non-array types
		{
			name:     "scalar integer",
			pgType:   "integer",
			expected: 0,
		},
		{
			name:     "scalar text",
			pgType:   "text",
			expected: 0,
		},
		{
			name:     "scalar int4",
			pgType:   "int4",
			expected: 0,
		},
		{
			name:     "empty string",
			pgType:   "",
			expected: 0,
		},

		// 1D arrays - bracket notation
		{
			name:     "integer[] - 1D array",
			pgType:   "integer[]",
			expected: 1,
		},
		{
			name:     "text[] - 1D array",
			pgType:   "text[]",
			expected: 1,
		},
		{
			name:     "varchar[] - 1D array",
			pgType:   "varchar[]",
			expected: 1,
		},

		// 1D arrays - underscore notation
		{
			name:     "_int4 - 1D array (underscore notation)",
			pgType:   "_int4",
			expected: 1,
		},
		{
			name:     "_text - 1D array (underscore notation)",
			pgType:   "_text",
			expected: 1,
		},
		{
			name:     "_varchar - 1D array (underscore notation)",
			pgType:   "_varchar",
			expected: 1,
		},
		{
			name:     "_int8 - 1D array (bigint)",
			pgType:   "_int8",
			expected: 1,
		},

		// 2D arrays
		{
			name:     "integer[][] - 2D array",
			pgType:   "integer[][]",
			expected: 2,
		},
		{
			name:     "text[][] - 2D array",
			pgType:   "text[][]",
			expected: 2,
		},
		{
			name:     "_int4[] - 2D array (underscore + bracket)",
			pgType:   "_int4[]",
			expected: 2,
		},
		{
			name:     "_text[] - 2D array (underscore + bracket)",
			pgType:   "_text[]",
			expected: 2,
		},

		// 3D arrays
		{
			name:     "integer[][][] - 3D array",
			pgType:   "integer[][][]",
			expected: 3,
		},
		{
			name:     "text[][][] - 3D array",
			pgType:   "text[][][]",
			expected: 3,
		},
		{
			name:     "_int4[][] - 3D array (underscore + 2 brackets)",
			pgType:   "_int4[][]",
			expected: 3,
		},

		// 4D arrays (edge case)
		{
			name:     "integer[][][][] - 4D array",
			pgType:   "integer[][][][]",
			expected: 4,
		},

		// Edge cases
		{
			name:     "just underscore",
			pgType:   "_",
			expected: 1, // Underscore prefix = 1D array
		},
		{
			name:     "underscore with number",
			pgType:   "_123",
			expected: 1,
		},
		{
			name:     "brackets in middle (invalid but should handle gracefully)",
			pgType:   "int[]eger",
			expected: 0, // Won't match trailing [] pattern
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detectArrayDimensions(tt.pgType)
			if result != tt.expected {
				t.Errorf("detectArrayDimensions(%q) = %d, expected %d", tt.pgType, result, tt.expected)
			}
		})
	}
}

// TestFieldSchemaWithDimensions tests that FieldSchema properly stores dimension information
func TestFieldSchemaWithDimensions(t *testing.T) {
	tests := []struct {
		name       string
		field      FieldSchema
		wantDims   int
		wantRepeat bool
	}{
		{
			name: "1D array",
			field: FieldSchema{
				Name:       "tags",
				Type:       "text",
				Repeated:   true,
				Dimensions: 1,
			},
			wantDims:   1,
			wantRepeat: true,
		},
		{
			name: "2D array",
			field: FieldSchema{
				Name:       "matrix",
				Type:       "integer",
				Repeated:   true,
				Dimensions: 2,
			},
			wantDims:   2,
			wantRepeat: true,
		},
		{
			name: "3D array",
			field: FieldSchema{
				Name:       "cube",
				Type:       "double precision",
				Repeated:   true,
				Dimensions: 3,
			},
			wantDims:   3,
			wantRepeat: true,
		},
		{
			name: "scalar field",
			field: FieldSchema{
				Name:       "id",
				Type:       "integer",
				Repeated:   false,
				Dimensions: 0,
			},
			wantDims:   0,
			wantRepeat: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.field.Dimensions != tt.wantDims {
				t.Errorf("field.Dimensions = %d, expected %d", tt.field.Dimensions, tt.wantDims)
			}
			if tt.field.Repeated != tt.wantRepeat {
				t.Errorf("field.Repeated = %v, expected %v", tt.field.Repeated, tt.wantRepeat)
			}
		})
	}
}

// TestSchemaWithMultiDimensionalArrays tests that Schema properly handles fields with various dimensions
func TestSchemaWithMultiDimensionalArrays(t *testing.T) {
	fields := []FieldSchema{
		{Name: "id", Type: "integer", Repeated: false, Dimensions: 0},
		{Name: "tags", Type: "text", Repeated: true, Dimensions: 1},
		{Name: "matrix", Type: "integer", Repeated: true, Dimensions: 2},
		{Name: "cube", Type: "float", Repeated: true, Dimensions: 3},
	}

	schema := NewSchema(fields)

	// Test that we can find each field
	tests := []struct {
		fieldName string
		wantDims  int
		wantFound bool
	}{
		{"id", 0, true},
		{"tags", 1, true},
		{"matrix", 2, true},
		{"cube", 3, true},
		{"nonexistent", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.fieldName, func(t *testing.T) {
			field, found := schema.FindField(tt.fieldName)
			if found != tt.wantFound {
				t.Errorf("FindField(%q) found = %v, expected %v", tt.fieldName, found, tt.wantFound)
				return
			}
			if found && field.Dimensions != tt.wantDims {
				t.Errorf("field %q dimensions = %d, expected %d", tt.fieldName, field.Dimensions, tt.wantDims)
			}
		})
	}
}
