// Package sqlite provides SQLite type provider for CEL type system integration.
package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"

	"github.com/spandigital/cel2sql/v3/schema"
)

// Sentinel errors for the sqlite package.
var (
	// ErrInvalidSchema indicates a problem with the provided schema or database introspection.
	ErrInvalidSchema = errors.New("invalid schema")
)

// validTableName matches safe SQLite table names (letters, digits, underscores).
var validTableName = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// FieldSchema is an alias for schema.FieldSchema.
type FieldSchema = schema.FieldSchema

// Schema is an alias for schema.Schema.
type Schema = schema.Schema

// NewSchema creates a new Schema. This is an alias for schema.NewSchema.
var NewSchema = schema.NewSchema

// TypeProvider interface for SQLite type providers.
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

// NewTypeProvider creates a new SQLite type provider with pre-defined schemas.
func NewTypeProvider(schemas map[string]Schema) TypeProvider {
	return &typeProvider{schemas: schemas}
}

// NewTypeProviderWithConnection creates a new SQLite type provider that can introspect database schemas.
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

// LoadTableSchema loads schema information for a table from the database using PRAGMA table_info.
func (tp *typeProvider) LoadTableSchema(ctx context.Context, tableName string) error {
	if tp.db == nil {
		return fmt.Errorf("%w: no database connection available", ErrInvalidSchema)
	}

	// Validate table name to prevent SQL injection (PRAGMA doesn't support parameterized queries)
	if !validTableName.MatchString(tableName) {
		return fmt.Errorf("%w: invalid table name %q", ErrInvalidSchema, tableName)
	}

	// PRAGMA doesn't support parameterized queries, but we've validated the table name above
	// #nosec G202 - table name is validated against strict regex pattern
	query := fmt.Sprintf("PRAGMA table_info(%s)", tableName)

	rows, err := tp.db.QueryContext(ctx, query)
	if err != nil {
		return fmt.Errorf("%w: failed to query table schema", ErrInvalidSchema)
	}
	defer func() { _ = rows.Close() }()

	var fields []FieldSchema
	for rows.Next() {
		var cid int
		var name, colType string
		var notnull int
		var dfltValue *string
		var pk int

		if err := rows.Scan(&cid, &name, &colType, &notnull, &dfltValue, &pk); err != nil {
			return fmt.Errorf("%w: failed to scan row", ErrInvalidSchema)
		}

		field := sqliteColumnToFieldSchema(name, colType)
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

// sqliteColumnToFieldSchema converts SQLite column metadata to a FieldSchema.
func sqliteColumnToFieldSchema(name, colType string) FieldSchema {
	normalizedType := normalizeSQLiteType(colType)
	isJSON := strings.EqualFold(colType, "json") || strings.EqualFold(colType, "jsonb")

	return FieldSchema{
		Name:   name,
		Type:   normalizedType,
		IsJSON: isJSON,
	}
}

// Normalized type constants used by normalizeSQLiteType.
const (
	sqliteTypeText     = "text"
	sqliteTypeInteger  = "integer"
	sqliteTypeReal     = "real"
	sqliteTypeBlob     = "blob"
	sqliteTypeJSON     = "json"
	sqliteTypeBool     = "boolean"
	sqliteTypeDatetime = "datetime"
)

// normalizeSQLiteType converts a SQLite column type declaration to a normalized type name.
// SQLite uses type affinity, so we map common type names to our internal types.
func normalizeSQLiteType(colType string) string {
	upper := strings.ToUpper(strings.TrimSpace(colType))

	// Check for exact matches first
	switch upper {
	case "TEXT", "VARCHAR", "CHAR", "CLOB":
		return sqliteTypeText
	case "INTEGER", "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "BIGINT":
		return sqliteTypeInteger
	case "REAL", "DOUBLE", "FLOAT", "NUMERIC", "DECIMAL":
		return sqliteTypeReal
	case "BOOLEAN", "BOOL":
		return sqliteTypeBool
	case "BLOB":
		return sqliteTypeBlob
	case "JSON", "JSONB":
		return sqliteTypeJSON
	case "DATETIME", "TIMESTAMP":
		return sqliteTypeDatetime
	}

	// Check for type names that contain known keywords (e.g., "VARCHAR(255)")
	if strings.Contains(upper, "INT") {
		return sqliteTypeInteger
	}
	if strings.Contains(upper, "CHAR") || strings.Contains(upper, "CLOB") || strings.Contains(upper, "TEXT") {
		return sqliteTypeText
	}
	if strings.Contains(upper, "BLOB") {
		return sqliteTypeBlob
	}
	if strings.Contains(upper, "REAL") || strings.Contains(upper, "FLOA") || strings.Contains(upper, "DOUBLE") {
		return sqliteTypeReal
	}

	// Default to text for unknown types (SQLite's flexible typing)
	return sqliteTypeText
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

	exprType := sqliteTypeToCELExprType(field)
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

// sqliteTypeToCELExprType converts a SQLite field schema to a CEL expression type.
func sqliteTypeToCELExprType(field *schema.FieldSchema) *exprpb.Type {
	baseType := sqliteBaseTypeToCEL(field.Type)
	if field.Repeated {
		return decls.NewListType(baseType)
	}
	return baseType
}

// sqliteBaseTypeToCEL converts a SQLite type name to a CEL expression type.
func sqliteBaseTypeToCEL(typeName string) *exprpb.Type {
	switch typeName {
	case "text", "varchar", "char", "clob":
		return decls.String
	case "integer", "int", "tinyint", "smallint", "mediumint", "bigint":
		return decls.Int
	case "real", "double", "float", "numeric", "decimal":
		return decls.Double
	case "boolean", "bool":
		return decls.Bool
	case "blob":
		return decls.Bytes
	case "json":
		return decls.Dyn
	case "datetime", "timestamp":
		return decls.Timestamp
	default:
		return decls.Dyn
	}
}
