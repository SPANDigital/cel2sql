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

	"github.com/spandigital/cel2sql/v2/sqltypes"
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

// Schema represents a PostgreSQL table schema as a slice of field schemas.
type Schema []FieldSchema

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
		return nil, fmt.Errorf("connection string exceeds maximum length of %d characters", MaxConnectionStringLength)
	}

	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		// Security: Don't wrap the error with %w to prevent exposing connection details
		// (credentials, hostnames, database names) in error messages or logs.
		// See pgx issues #1271 and #1428, CWE-209, CWE-532.
		return nil, errors.New("failed to create connection pool")
	}

	return &typeProvider{
		schemas: make(map[string]Schema),
		pool:    pool,
	}, nil
}

// LoadTableSchema loads schema information for a table from the database
func (p *typeProvider) LoadTableSchema(ctx context.Context, tableName string) error {
	if p.pool == nil {
		return errors.New("no database connection available")
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
		return fmt.Errorf("failed to query table schema: %w", err)
	}
	defer rows.Close()

	var schema Schema
	for rows.Next() {
		var columnName, dataType, isNullable string
		var columnDefault *string
		var elementType string

		err := rows.Scan(&columnName, &dataType, &isNullable, &columnDefault, &elementType)
		if err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
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

		schema = append(schema, field)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating rows: %w", err)
	}

	p.schemas[tableName] = schema
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
		return nil, false
	}

	// For single-level types, return the schema directly
	if len(typeNames) == 1 {
		return schema, true
	}

	// For nested types, traverse the schema hierarchy
	for _, tn := range typeNames[1:] {
		var s Schema
		for _, fieldSchema := range schema {
			if fieldSchema.Name == tn {
				s = fieldSchema.Schema
				break
			}
		}
		if len(s) == 0 {
			return nil, false
		}
		schema = s
	}
	return schema, true
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

	fieldNames := make([]string, len(schema))
	for i, field := range schema {
		fieldNames[i] = field.Name
	}
	return fieldNames, true
}

func (p *typeProvider) FindStructFieldType(structType, fieldName string) (*types.FieldType, bool) {
	schema, found := p.findSchema(structType)
	if !found {
		return nil, false
	}
	var field *FieldSchema
	for _, fieldSchema := range schema {
		if fieldSchema.Name == fieldName {
			field = &fieldSchema
			break
		}
	}
	if field == nil {
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
	default:
		// Handle composite types
		if strings.Contains(field.Type, "composite") || len(field.Schema) > 0 {
			exprType = decls.NewObjectType(strings.Join([]string{structType, fieldName}, "."))
		} else {
			// Default to string for unknown types
			exprType = decls.String
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
