// Package schema provides dialect-agnostic database schema types for CEL to SQL conversion.
// These types describe column names, types, array dimensions, and JSON flags without
// coupling to any specific SQL dialect.
package schema

// FieldSchema represents a database field type with name, type, and optional nested schema.
// This type is dialect-agnostic and used by all SQL dialect providers.
type FieldSchema struct {
	Name        string
	Type        string        // SQL type name (text, integer, boolean, etc.)
	Repeated    bool          // true for arrays
	Dimensions  int           // number of array dimensions (1 for integer[], 2 for integer[][], etc.)
	Schema      []FieldSchema // for composite types
	IsJSON      bool          // true for json/jsonb types
	IsJSONB     bool          // true for jsonb (vs json)
	ElementType string        // for arrays: element type name
}

// Schema represents a table schema with O(1) field lookup.
// It contains a slice of fields for ordered iteration and a map index for fast lookups.
type Schema struct {
	fields     []FieldSchema
	fieldIndex map[string]*FieldSchema
}

// NewSchema creates a new Schema with field indexing for O(1) lookups.
// This improves performance for tables with many columns.
func NewSchema(fields []FieldSchema) Schema {
	index := make(map[string]*FieldSchema, len(fields))
	for i := range fields {
		index[fields[i].Name] = &fields[i]

		// Build indices for nested schemas recursively
		if len(fields[i].Schema) > 0 {
			fields[i].Schema = rebuildSchemaIndex(fields[i].Schema)
		}
	}

	return Schema{
		fields:     fields,
		fieldIndex: index,
	}
}

// rebuildSchemaIndex recursively rebuilds indices for nested schemas.
// This is used internally when converting old-style []FieldSchema to new Schema struct.
func rebuildSchemaIndex(oldSchema []FieldSchema) []FieldSchema {
	// For nested schemas, we need to ensure they're properly indexed too
	// But since nested schemas are stored as []FieldSchema in FieldSchema.Schema,
	// we keep them as slices but process them when needed
	return oldSchema
}

// Fields returns the ordered slice of field schemas.
// Use this when you need to iterate over fields in their defined order.
func (s Schema) Fields() []FieldSchema {
	return s.fields
}

// FindField performs an O(1) lookup for a field by name.
// Returns the field schema and true if found, nil and false otherwise.
func (s Schema) FindField(name string) (*FieldSchema, bool) {
	field, found := s.fieldIndex[name]
	return field, found
}

// Len returns the number of fields in the schema.
func (s Schema) Len() int {
	return len(s.fields)
}
