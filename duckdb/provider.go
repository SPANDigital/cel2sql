// Package duckdb provides DuckDB type provider for CEL type system integration.
package duckdb

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
)

// Sentinel errors for the duckdb package.
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

// TypeProvider interface for DuckDB type providers.
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

// NewTypeProvider creates a new DuckDB type provider with pre-defined schemas.
func NewTypeProvider(schemas map[string]Schema) TypeProvider {
	return &typeProvider{schemas: schemas}
}

// NewTypeProviderWithConnection creates a new DuckDB type provider that can introspect database schemas.
// The caller owns the *sql.DB and is responsible for closing it.
// This works with any DuckDB driver that implements database/sql (e.g., github.com/marcboeker/go-duckdb).
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
		SELECT column_name, data_type, is_nullable
		FROM information_schema.columns
		WHERE table_name = ?
		ORDER BY ordinal_position
	`

	rows, err := tp.db.QueryContext(ctx, query, tableName)
	if err != nil {
		return fmt.Errorf("%w: failed to query table schema", ErrInvalidSchema)
	}
	defer func() { _ = rows.Close() }()

	var fields []FieldSchema
	for rows.Next() {
		var columnName, dataType, isNullable string

		if err := rows.Scan(&columnName, &dataType, &isNullable); err != nil {
			return fmt.Errorf("%w: failed to scan row", ErrInvalidSchema)
		}

		field := duckdbColumnToFieldSchema(columnName, dataType)
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

// duckdbColumnToFieldSchema converts DuckDB column metadata to a FieldSchema.
func duckdbColumnToFieldSchema(columnName, dataType string) FieldSchema {
	// DuckDB array types appear as "INTEGER[]", "VARCHAR[]", etc.
	isArray, elementType, dimensions := detectDuckDBArray(dataType)
	isJSON := strings.EqualFold(dataType, "json")

	if isArray {
		return FieldSchema{
			Name:        columnName,
			Type:        strings.ToLower(elementType),
			Repeated:    true,
			Dimensions:  dimensions,
			ElementType: strings.ToLower(elementType),
		}
	}

	return FieldSchema{
		Name:   columnName,
		Type:   normalizeDuckDBType(dataType),
		IsJSON: isJSON,
	}
}

// detectDuckDBArray detects if a DuckDB data type is an array and returns element type and dimensions.
func detectDuckDBArray(dataType string) (isArray bool, elementType string, dimensions int) {
	// Count trailing [] pairs
	remaining := dataType
	dims := 0
	for strings.HasSuffix(remaining, "[]") {
		dims++
		remaining = strings.TrimSuffix(remaining, "[]")
	}

	if dims > 0 {
		return true, remaining, dims
	}
	return false, "", 0
}

// normalizeDuckDBType normalizes a DuckDB type name to lowercase.
func normalizeDuckDBType(dataType string) string {
	return strings.ToLower(dataType)
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

	exprType := duckdbTypeToCELExprType(field)
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

// duckdbTypeToCELExprType converts a DuckDB field schema to a CEL expression type.
func duckdbTypeToCELExprType(field *schema.FieldSchema) *exprpb.Type {
	baseType := duckdbBaseTypeToCEL(field.Type)
	if field.Repeated {
		return decls.NewListType(baseType)
	}
	return baseType
}

// duckdbBaseTypeToCEL converts a DuckDB type name to a CEL expression type.
func duckdbBaseTypeToCEL(typeName string) *exprpb.Type {
	switch typeName {
	case "varchar", "text", "char", "bpchar", "name":
		return decls.String
	case "bigint", "integer", "int", "int4", "int8", "smallint", "int2", "tinyint", "hugeint":
		return decls.Int
	case "double", "float", "real", "float4", "float8", "numeric", "decimal":
		return decls.Double
	case "boolean", "bool":
		return decls.Bool
	case "blob", "bytea":
		return decls.Bytes
	case "json":
		return decls.Dyn
	case "timestamp", "timestamptz", "timestamp with time zone", "timestamp without time zone":
		return decls.Timestamp
	default:
		return decls.Dyn
	}
}
