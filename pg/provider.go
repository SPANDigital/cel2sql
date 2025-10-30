// Package pg provides PostgreSQL type provider for CEL type system integration.
package pg

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/jackc/pgx/v5/pgxpool"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"

	"github.com/spandigital/cel2sql/v3/sqltypes"
)

// Sentinel errors specific to the pg package
var (
	// ErrInvalidSchema indicates a problem with the provided schema
	ErrInvalidSchema = errors.New("invalid schema")
)

const (
	typeJSON  = "json"
	typeJSONB = "jsonb"

	// MaxConnectionStringLength limits connection string length to prevent
	// resource exhaustion and align with ODBC standard (1024 chars).
	// Legitimate PostgreSQL connection strings rarely exceed a few hundred characters.
	MaxConnectionStringLength = 1000

	// Error messages (sanitized for end users per CWE-209)
	errMsgUnknownEnum = "unknown enum value"
	errMsgUnknownType = "unknown type in schema"
)

// FieldSchema represents a PostgreSQL field type with name, type, and optional nested schema.
type FieldSchema struct {
	Name        string
	Type        string        // PostgreSQL type name (text, integer, boolean, etc.)
	Repeated    bool          // true for arrays
	Schema      []FieldSchema // for composite types
	IsJSON      bool          // true for json/jsonb types
	IsJSONB     bool          // true for jsonb (vs json)
	ElementType string        // for arrays: element type name
}

// Schema represents a PostgreSQL table schema with O(1) field lookup.
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

// TypeProvider interface for PostgreSQL type providers
type TypeProvider interface {
	types.Provider
	LoadTableSchema(ctx context.Context, tableName string) error
	GetSchemas() map[string]Schema
	Close()
}

type typeProvider struct {
	schemas map[string]Schema
	pool    *pgxpool.Pool
}

// NewTypeProvider creates a new PostgreSQL type provider with pre-defined schemas
func NewTypeProvider(schemas map[string]Schema) TypeProvider {
	return &typeProvider{schemas: schemas}
}

// NewTypeProviderWithConnection creates a new PostgreSQL type provider that can introspect database schemas
func NewTypeProviderWithConnection(ctx context.Context, connectionString string) (TypeProvider, error) {
	// Validate connection string length to prevent DoS and align with industry standards
	if len(connectionString) > MaxConnectionStringLength {
		return nil, fmt.Errorf("%w: connection string exceeds maximum length of %d characters", ErrInvalidSchema, MaxConnectionStringLength)
	}

	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		// Security: Don't wrap the error with %w to prevent exposing connection details
		// (credentials, hostnames, database names) in error messages or logs.
		// See pgx issues #1271 and #1428, CWE-209, CWE-532.
		return nil, fmt.Errorf("%w: failed to create connection pool", ErrInvalidSchema)
	}

	// Validate connection works immediately rather than on first query
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		// Security: Same sanitized error approach as above
		return nil, fmt.Errorf("%w: failed to connect to database", ErrInvalidSchema)
	}

	return &typeProvider{
		schemas: make(map[string]Schema),
		pool:    pool,
	}, nil
}

// LoadTableSchema loads schema information for a table from the database
func (p *typeProvider) LoadTableSchema(ctx context.Context, tableName string) error {
	if p.pool == nil {
		return fmt.Errorf("%w: no database connection available", ErrInvalidSchema)
	}

	query := `
		SELECT 
			column_name, 
			data_type, 
			is_nullable, 
			column_default,
			CASE 
				WHEN data_type = 'ARRAY' THEN 
					(SELECT data_type FROM information_schema.element_types 
					 WHERE object_name = $1 
					 AND collection_type_identifier = (
						SELECT dtd_identifier FROM information_schema.columns 
						WHERE table_name = $1 AND column_name = c.column_name
					))
				ELSE data_type
			END as element_type
		FROM information_schema.columns c
		WHERE table_name = $1 
		ORDER BY ordinal_position
	`

	rows, err := p.pool.Query(ctx, query, tableName)
	if err != nil {
		return fmt.Errorf("%w: failed to query table schema: %w", ErrInvalidSchema, err)
	}
	defer rows.Close()

	var fields []FieldSchema
	for rows.Next() {
		var columnName, dataType, isNullable string
		var columnDefault *string
		var elementType string

		err := rows.Scan(&columnName, &dataType, &isNullable, &columnDefault, &elementType)
		if err != nil {
			return fmt.Errorf("%w: failed to scan row: %w", ErrInvalidSchema, err)
		}

		isArray := dataType == "ARRAY"
		isJSON := elementType == typeJSON || elementType == typeJSONB
		isJSONB := elementType == typeJSONB

		field := FieldSchema{
			Name:        columnName,
			Type:        elementType, // Use element type for arrays, or data_type for non-arrays
			Repeated:    isArray,
			IsJSON:      isJSON,
			IsJSONB:     isJSONB,
			ElementType: "", // Will be set below for arrays
		}

		// For arrays, elementType is the array element type, so store it
		if isArray {
			field.ElementType = elementType
		}

		fields = append(fields, field)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("%w: error iterating rows: %w", ErrInvalidSchema, err)
	}

	p.schemas[tableName] = NewSchema(fields)
	return nil
}

// Close closes the database connection pool
func (p *typeProvider) Close() {
	if p.pool != nil {
		p.pool.Close()
	}
}

// GetSchemas returns the schema map
func (p *typeProvider) GetSchemas() map[string]Schema {
	return p.schemas
}

func (p *typeProvider) EnumValue(_ string) ref.Val {
	// Don't expose enum names to users - return generic message
	// Internal details are not needed here since CEL runtime handles this
	return types.NewErr(errMsgUnknownEnum)
}

func (p *typeProvider) FindIdent(_ string) (ref.Val, bool) {
	return nil, false
}

func (p *typeProvider) findSchema(typeName string) (Schema, bool) {
	typeNames := strings.Split(typeName, ".")
	schema, found := p.schemas[typeNames[0]]
	if !found {
		return Schema{}, false
	}

	// For single-level types, return the schema directly
	if len(typeNames) == 1 {
		return schema, true
	}

	// For nested types, traverse the schema hierarchy using O(1) lookups
	currentFields := schema.fields
	for _, tn := range typeNames[1:] {
		// Use O(1) indexed lookup instead of linear search
		var nestedField *FieldSchema
		for i := range currentFields {
			if currentFields[i].Name == tn {
				nestedField = &currentFields[i]
				break
			}
		}

		if nestedField == nil || len(nestedField.Schema) == 0 {
			return Schema{}, false
		}
		currentFields = nestedField.Schema
	}

	// Convert the nested []FieldSchema to Schema for the return
	return NewSchema(currentFields), true
}

func (p *typeProvider) FindStructType(structType string) (*types.Type, bool) {
	_, found := p.findSchema(structType)
	if !found {
		return nil, false
	}
	return types.NewObjectType(structType), true
}

func (p *typeProvider) FindStructFieldNames(structType string) ([]string, bool) {
	schema, found := p.findSchema(structType)
	if !found {
		return nil, false
	}

	fields := schema.Fields()
	fieldNames := make([]string, len(fields))
	for i, field := range fields {
		fieldNames[i] = field.Name
	}
	return fieldNames, true
}

func (p *typeProvider) FindStructFieldType(structType, fieldName string) (*types.FieldType, bool) {
	schema, found := p.findSchema(structType)
	if !found {
		return nil, false
	}

	// Use O(1) indexed lookup instead of linear search
	field, found := schema.FindField(fieldName)
	if !found {
		return nil, false
	}

	var exprType *exprpb.Type
	switch field.Type {
	case "text", "varchar", "char", "character varying", "character":
		exprType = decls.String
	case "bytea":
		exprType = decls.Bytes
	case "boolean", "bool":
		exprType = decls.Bool
	case "integer", "int", "int4", "bigint", "int8", "smallint", "int2":
		exprType = decls.Int
	case "real", "float4", "double precision", "float8", "numeric", "decimal":
		exprType = decls.Double
	case "timestamp", "timestamptz", "timestamp with time zone", "timestamp without time zone":
		exprType = decls.Timestamp
	case "date":
		exprType = sqltypes.Date
	case "time", "timetz", "time with time zone", "time without time zone":
		exprType = sqltypes.Time
	case "json", "jsonb":
		// JSON and JSONB types are treated as dynamic objects in CEL
		exprType = decls.Dyn
	case "uuid":
		// UUID is represented as bytes for proper type handling
		exprType = decls.Bytes
	case "inet", "cidr":
		// Network address types are represented as strings
		// Note: Limited CEL operations available (equality, comparison)
		exprType = decls.String
	case "macaddr", "macaddr8":
		// MAC address types are represented as strings
		exprType = decls.String
	case "xml":
		// XML data is represented as string
		exprType = decls.String
	case "money":
		// Money type is represented as double for numeric operations
		exprType = decls.Double
	case "tsvector", "tsquery":
		// Full-text search types are represented as strings
		exprType = decls.String
	default:
		// Handle composite types
		if strings.Contains(field.Type, "composite") || len(field.Schema) > 0 {
			exprType = decls.NewObjectType(strings.Join([]string{structType, fieldName}, "."))
		} else {
			// Unknown type - return not found instead of defaulting to string
			// This prevents silent type mismatches and incorrect SQL generation
			return nil, false
		}
	}

	if field.Repeated {
		exprType = decls.NewListType(exprType)
	}

	// Convert exprpb.Type to types.Type
	celType, err := types.ExprTypeToType(exprType)
	if err != nil {
		return nil, false
	}

	return &types.FieldType{
		Type: celType,
	}, true
}

func (p *typeProvider) NewValue(_ string, _ map[string]ref.Val) ref.Val {
	// Don't expose type names to users - return generic message
	// Internal details are not needed here since CEL runtime handles this
	return types.NewErr(errMsgUnknownType)
}

var _ types.Provider = new(typeProvider)
