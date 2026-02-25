// Package mysql provides MySQL type provider for CEL type system integration.
package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"

	"github.com/spandigital/cel2sql/v3/schema"
	"github.com/spandigital/cel2sql/v3/sqltypes"
)

// Sentinel errors for the mysql package.
var (
	// ErrInvalidSchema indicates a problem with the provided schema or database introspection.
	ErrInvalidSchema = errors.New("invalid schema")
)

// FieldSchema is an alias for schema.FieldSchema.
type FieldSchema = schema.FieldSchema

// Schema is an alias for schema.Schema.
type Schema = schema.Schema

// NewSchema creates a new Schema. This is an alias for schema.NewSchema.
var NewSchema = schema.NewSchema

// TypeProvider interface for MySQL type providers.
type TypeProvider interface {
	types.Provider
	LoadTableSchema(ctx context.Context, tableName string) error
	GetSchemas() map[string]Schema
	Close()
}

type typeProvider struct {
	schemas map[string]Schema
	db      *sql.DB
}

// NewTypeProvider creates a new MySQL type provider with pre-defined schemas.
func NewTypeProvider(schemas map[string]Schema) TypeProvider {
	return &typeProvider{schemas: schemas}
}

// NewTypeProviderWithConnection creates a new MySQL type provider that can introspect database schemas.
// The caller owns the *sql.DB and is responsible for closing it.
func NewTypeProviderWithConnection(_ context.Context, db *sql.DB) (TypeProvider, error) {
	if db == nil {
		return nil, fmt.Errorf("%w: db connection must not be nil", ErrInvalidSchema)
	}

	return &typeProvider{
		schemas: make(map[string]Schema),
		db:      db,
	}, nil
}

// LoadTableSchema loads schema information for a table from the database.
func (tp *typeProvider) LoadTableSchema(ctx context.Context, tableName string) error {
	if tp.db == nil {
		return fmt.Errorf("%w: no database connection available", ErrInvalidSchema)
	}

	query := `
		SELECT column_name, data_type, column_type, is_nullable
		FROM information_schema.columns
		WHERE table_schema = DATABASE() AND table_name = ?
		ORDER BY ordinal_position
	`

	rows, err := tp.db.QueryContext(ctx, query, tableName)
	if err != nil {
		return fmt.Errorf("%w: failed to query table schema", ErrInvalidSchema)
	}
	defer func() { _ = rows.Close() }()

	var fields []FieldSchema
	for rows.Next() {
		var columnName, dataType, columnType, isNullable string

		if err := rows.Scan(&columnName, &dataType, &columnType, &isNullable); err != nil {
			return fmt.Errorf("%w: failed to scan row", ErrInvalidSchema)
		}

		field := mysqlColumnToFieldSchema(columnName, dataType, columnType)
		fields = append(fields, field)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("%w: error iterating rows", ErrInvalidSchema)
	}

	if len(fields) == 0 {
		return fmt.Errorf("%w: table %q has no columns or does not exist", ErrInvalidSchema, tableName)
	}

	tp.schemas[tableName] = NewSchema(fields)
	return nil
}

// mysqlColumnToFieldSchema converts MySQL column metadata to a FieldSchema.
func mysqlColumnToFieldSchema(columnName, dataType, _ string) FieldSchema {
	// Normalize data type to lowercase
	dataType = strings.ToLower(dataType)

	isJSON := dataType == "json"

	return FieldSchema{
		Name:   columnName,
		Type:   dataType,
		IsJSON: isJSON,
	}
}

// Close is a no-op since we don't own the *sql.DB.
func (tp *typeProvider) Close() {
	// No-op: caller owns the *sql.DB connection
}

// GetSchemas returns the schemas known to this type provider.
func (tp *typeProvider) GetSchemas() map[string]Schema {
	return tp.schemas
}

// EnumValue implements types.Provider.
func (tp *typeProvider) EnumValue(_ string) ref.Val {
	return types.NewErr("unknown enum value")
}

// FindIdent implements types.Provider.
func (tp *typeProvider) FindIdent(_ string) (ref.Val, bool) {
	return nil, false
}

// FindStructType implements types.Provider.
func (tp *typeProvider) FindStructType(structType string) (*types.Type, bool) {
	if _, ok := tp.schemas[structType]; ok {
		return types.NewObjectType(structType), true
	}
	return nil, false
}

// FindStructFieldNames implements types.Provider.
func (tp *typeProvider) FindStructFieldNames(structType string) ([]string, bool) {
	s, ok := tp.schemas[structType]
	if !ok {
		return nil, false
	}
	fields := s.Fields()
	names := make([]string, len(fields))
	for i, f := range fields {
		names[i] = f.Name
	}
	return names, true
}

// FindStructFieldType implements types.Provider.
func (tp *typeProvider) FindStructFieldType(structType, fieldName string) (*types.FieldType, bool) {
	s, ok := tp.schemas[structType]
	if !ok {
		return nil, false
	}
	field, found := s.FindField(fieldName)
	if !found {
		return nil, false
	}

	exprType := mysqlTypeToCELExprType(field)

	celType, err := types.ExprTypeToType(exprType)
	if err != nil {
		return nil, false
	}

	return &types.FieldType{
		Type: celType,
	}, true
}

// NewValue implements types.Provider.
func (tp *typeProvider) NewValue(_ string, _ map[string]ref.Val) ref.Val {
	return types.NewErr("unknown type in schema")
}

// mysqlTypeToCELExprType converts a MySQL field schema to a CEL expression type.
func mysqlTypeToCELExprType(field *schema.FieldSchema) *exprpb.Type {
	baseType := mysqlBaseTypeToCEL(field.Type)
	if field.Repeated {
		return decls.NewListType(baseType)
	}
	return baseType
}

// mysqlBaseTypeToCEL converts a MySQL type name to a CEL expression type.
func mysqlBaseTypeToCEL(typeName string) *exprpb.Type {
	switch typeName {
	case "varchar", "char", "text", "tinytext", "mediumtext", "longtext", "enum", "set":
		return decls.String
	case "int", "integer", "tinyint", "smallint", "mediumint", "bigint":
		return decls.Int
	case "float", "double", "decimal", "numeric", "real":
		return decls.Double
	case "boolean", "bool":
		return decls.Bool
	case "blob", "binary", "varbinary", "tinyblob", "mediumblob", "longblob":
		return decls.Bytes
	case "json":
		return decls.Dyn
	case "datetime", "timestamp":
		return decls.Timestamp
	case "date":
		return sqltypes.Date
	case "time":
		return sqltypes.Time
	default:
		return decls.Dyn
	}
}
