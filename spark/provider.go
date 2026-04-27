// Package spark provides a Spark SQL type provider for CEL type system integration.
package spark

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
	"github.com/spandigital/cel2sql/v3/sqltypes"
)

// Sentinel errors for the spark package.
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

// validTableName matches Spark identifier rules. The Spark catalog APIs
// (DESCRIBE TABLE) do not accept positional parameters, so we validate the
// table name string before interpolating it into the query.
var validTableName = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_.]*$`)

// TypeProvider is the Spark SQL type provider interface.
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

// NewTypeProvider creates a new Spark SQL type provider with pre-defined schemas.
func NewTypeProvider(schemas map[string]Schema) TypeProvider {
	return &typeProvider{schemas: schemas}
}

// NewTypeProviderWithConnection creates a new Spark SQL type provider that can
// introspect schemas via DESCRIBE TABLE. The caller owns the *sql.DB and is
// responsible for closing it. Works with any database/sql driver (e.g. gohive
// for Hive/SparkThrift, or a JDBC bridge).
func NewTypeProviderWithConnection(_ context.Context, db *sql.DB) (TypeProvider, error) {
	if db == nil {
		return nil, fmt.Errorf("%w: db connection must not be nil", ErrInvalidSchema)
	}
	return &typeProvider{
		schemas: make(map[string]Schema),
		db:      db,
	}, nil
}

// LoadTableSchema loads schema information for a table from the Spark catalog
// using DESCRIBE TABLE. The table name must be a valid Spark identifier (letters,
// digits, underscore, dot for catalog.schema.table); this is enforced because
// DESCRIBE TABLE does not accept parameterized queries.
func (tp *typeProvider) LoadTableSchema(ctx context.Context, tableName string) error {
	if !validTableName.MatchString(tableName) {
		return fmt.Errorf("%w: invalid table name %q", ErrInvalidSchema, tableName)
	}
	if tp.db == nil {
		return fmt.Errorf("%w: no database connection available", ErrInvalidSchema)
	}

	query := "DESCRIBE TABLE " + tableName
	rows, err := tp.db.QueryContext(ctx, query)
	if err != nil {
		return fmt.Errorf("%w: failed to describe table", ErrInvalidSchema)
	}
	defer func() { _ = rows.Close() }()

	var fields []FieldSchema
	for rows.Next() {
		var colName, dataType, comment sql.NullString
		if err := rows.Scan(&colName, &dataType, &comment); err != nil {
			return fmt.Errorf("%w: failed to scan row", ErrInvalidSchema)
		}
		// DESCRIBE TABLE in Spark emits trailing partition info rows whose col_name
		// starts with "#" or is empty; skip them.
		name := strings.TrimSpace(colName.String)
		if name == "" || strings.HasPrefix(name, "#") {
			continue
		}
		fields = append(fields, sparkColumnToFieldSchema(name, dataType.String))
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

// sparkColumnToFieldSchema converts a Spark DESCRIBE TABLE row into a FieldSchema,
// detecting array<T>, struct<...>, and primitive types.
func sparkColumnToFieldSchema(name, dataType string) FieldSchema {
	dt := strings.TrimSpace(strings.ToLower(dataType))

	if isArr, elem := detectSparkArray(dt); isArr {
		return FieldSchema{
			Name:        name,
			Type:        normalizeSparkType(elem),
			Repeated:    true,
			Dimensions:  1,
			ElementType: normalizeSparkType(elem),
		}
	}
	if strings.HasPrefix(dt, "struct<") {
		nested, ok := parseSparkStruct(dt)
		if ok {
			return FieldSchema{
				Name:   name,
				Type:   "struct",
				Schema: nested,
			}
		}
	}
	return FieldSchema{
		Name: name,
		Type: normalizeSparkType(dt),
	}
}

// detectSparkArray reports whether a Spark type string is array<T> and returns
// the element type T. Nested arrays (array<array<T>>) collapse to dimension 1
// in v1; multi-dimensional arrays are uncommon in Spark and unsupported by
// the dialect emitter.
func detectSparkArray(dt string) (isArray bool, elementType string) {
	if !strings.HasPrefix(dt, "array<") || !strings.HasSuffix(dt, ">") {
		return false, ""
	}
	return true, dt[len("array<") : len(dt)-1]
}

// parseSparkStruct parses a Spark struct type string of the form
// `struct<field1:type1,field2:type2>` into nested FieldSchemas. Returns
// false if the input cannot be parsed (e.g., malformed nesting); callers
// should fall back to treating the field as opaque.
func parseSparkStruct(dt string) ([]FieldSchema, bool) {
	if !strings.HasPrefix(dt, "struct<") || !strings.HasSuffix(dt, ">") {
		return nil, false
	}
	body := dt[len("struct<") : len(dt)-1]
	if body == "" {
		return nil, true
	}
	parts := splitTopLevel(body, ',')
	out := make([]FieldSchema, 0, len(parts))
	for _, p := range parts {
		colon := indexTopLevel(p, ':')
		if colon < 0 {
			return nil, false
		}
		fname := strings.TrimSpace(p[:colon])
		ftype := strings.TrimSpace(p[colon+1:])
		out = append(out, sparkColumnToFieldSchema(fname, ftype))
	}
	return out, true
}

// splitTopLevel splits s on sep, ignoring sep inside angle-bracket pairs.
func splitTopLevel(s string, sep byte) []string {
	var (
		out   []string
		depth int
		start int
	)
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '<':
			depth++
		case '>':
			if depth > 0 {
				depth--
			}
		case sep:
			if depth == 0 {
				out = append(out, s[start:i])
				start = i + 1
			}
		}
	}
	out = append(out, s[start:])
	return out
}

// indexTopLevel returns the index of the first sep at depth 0, or -1.
func indexTopLevel(s string, sep byte) int {
	depth := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '<':
			depth++
		case '>':
			if depth > 0 {
				depth--
			}
		case sep:
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

// normalizeSparkType lowercases and strips parameters from a Spark type string,
// e.g., "varchar(64)" → "varchar", "decimal(10,2)" → "decimal".
func normalizeSparkType(t string) string {
	t = strings.ToLower(strings.TrimSpace(t))
	if i := strings.IndexByte(t, '('); i >= 0 {
		return strings.TrimSpace(t[:i])
	}
	return t
}

// Close is a no-op since we don't own the *sql.DB.
func (tp *typeProvider) Close() {
	// No-op: caller owns the *sql.DB connection.
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
	exprType := sparkTypeToCELExprType(field)
	celType, err := types.ExprTypeToType(exprType)
	if err != nil {
		return nil, false
	}
	return &types.FieldType{Type: celType}, true
}

// NewValue implements types.Provider.
func (tp *typeProvider) NewValue(_ string, _ map[string]ref.Val) ref.Val {
	return types.NewErr("unknown type in schema")
}

// sparkTypeToCELExprType converts a Spark field schema to a CEL expression type.
func sparkTypeToCELExprType(field *schema.FieldSchema) *exprpb.Type {
	if field.Repeated {
		return decls.NewListType(sparkBaseTypeToCEL(field.ElementType))
	}
	return sparkBaseTypeToCEL(field.Type)
}

// sparkBaseTypeToCEL converts a Spark type name to a CEL expression type.
func sparkBaseTypeToCEL(typeName string) *exprpb.Type {
	switch typeName {
	case "string", "varchar", "char":
		return decls.String
	case "byte", "tinyint", "short", "smallint", "int", "integer", "long", "bigint":
		return decls.Int
	case "float", "double", "real", "decimal", "numeric":
		return decls.Double
	case "boolean", "bool":
		return decls.Bool
	case "binary":
		return decls.Bytes
	case "date":
		return sqltypes.Date
	case "timestamp", "timestamp_ntz":
		return decls.Timestamp
	default:
		return decls.Dyn
	}
}
