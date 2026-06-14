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
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"

	"github.com/spandigital/cel2sql/v3/internal/celprovider"
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
	celprovider.Base
	db *sql.DB
}

// NewTypeProvider creates a new MySQL type provider with pre-defined schemas.
func NewTypeProvider(schemas map[string]Schema) TypeProvider {
	return &typeProvider{Base: celprovider.Base{Schemas: schemas, Mapper: mysqlTypeToCELExprType}}
}

// NewTypeProviderWithConnection creates a new MySQL type provider that can introspect database schemas.
// The caller owns the *sql.DB and is responsible for closing it.
func NewTypeProviderWithConnection(_ context.Context, db *sql.DB) (TypeProvider, error) {
	if db == nil {
		return nil, fmt.Errorf("%w: db connection must not be nil", ErrInvalidSchema)
	}

	return &typeProvider{
		Base: celprovider.Base{Schemas: make(map[string]Schema), Mapper: mysqlTypeToCELExprType},
		db:   db,
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

	tp.Schemas[tableName] = NewSchema(fields)
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
