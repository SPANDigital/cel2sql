package pg

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/jackc/pgx/v5/pgxpool"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"

	"github.com/spandigital/cel2sql/sqltypes"
)

// PostgreSQL field type representation
type FieldSchema struct {
	Name     string
	Type     string        // PostgreSQL type name (text, integer, boolean, etc.)
	Repeated bool          // true for arrays
	Schema   []FieldSchema // for composite types
}

// PostgreSQL table schema
type Schema []FieldSchema

type typeProvider struct {
	schemas map[string]Schema
	pool    *pgxpool.Pool
}

// NewTypeProvider creates a new PostgreSQL type provider with pre-defined schemas
func NewTypeProvider(schemas map[string]Schema) *typeProvider {
	return &typeProvider{schemas: schemas}
}

// NewTypeProviderWithConnection creates a new PostgreSQL type provider that can introspect database schemas
func NewTypeProviderWithConnection(ctx context.Context, connectionString string) (*typeProvider, error) {
	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	return &typeProvider{
		schemas: make(map[string]Schema),
		pool:    pool,
	}, nil
}

// LoadTableSchema loads schema information for a table from the database
func (p *typeProvider) LoadTableSchema(ctx context.Context, tableName string) error {
	if p.pool == nil {
		return fmt.Errorf("no database connection available")
	}

	query := `
		SELECT column_name, data_type, is_nullable, column_default
		FROM information_schema.columns 
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

		err := rows.Scan(&columnName, &dataType, &isNullable, &columnDefault)
		if err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}

		field := FieldSchema{
			Name:     columnName,
			Type:     dataType,
			Repeated: strings.HasSuffix(dataType, "[]"), // PostgreSQL array notation
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

func (p *typeProvider) EnumValue(enumName string) ref.Val {
	return types.NewErr("unknown enum name '%s'", enumName)
}

func (p *typeProvider) FindIdent(identName string) (ref.Val, bool) {
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

func (p *typeProvider) FindType(typeName string) (*exprpb.Type, bool) {
	_, found := p.findSchema(typeName)
	if !found {
		return nil, false
	}
	return decls.NewTypeType(decls.NewObjectType(typeName)), true
}

func (p *typeProvider) FindFieldType(messageType string, fieldName string) (*ref.FieldType, bool) {
	schema, found := p.findSchema(messageType)
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

	var typ *exprpb.Type
	switch field.Type {
	case "text", "varchar", "char", "character varying", "character":
		typ = decls.String
	case "bytea":
		typ = decls.Bytes
	case "boolean", "bool":
		typ = decls.Bool
	case "integer", "int", "int4", "bigint", "int8", "smallint", "int2":
		typ = decls.Int
	case "real", "float4", "double precision", "float8", "numeric", "decimal":
		typ = decls.Double
	case "timestamp", "timestamptz", "timestamp with time zone", "timestamp without time zone":
		typ = decls.Timestamp
	case "date":
		typ = sqltypes.Date
	case "time", "timetz", "time with time zone", "time without time zone":
		typ = sqltypes.Time
	default:
		// Handle composite types
		if strings.Contains(field.Type, "composite") || len(field.Schema) > 0 {
			typ = decls.NewObjectType(strings.Join([]string{messageType, fieldName}, "."))
		} else {
			// Default to string for unknown types
			typ = decls.String
		}
	}

	if field.Repeated {
		typ = decls.NewListType(typ)
	}

	return &ref.FieldType{
		Type: typ,
	}, true
}

func (p *typeProvider) NewValue(typeName string, fields map[string]ref.Val) ref.Val {
	return types.NewErr("unknown type '%s'", typeName)
}

var _ ref.TypeProvider = new(typeProvider)
