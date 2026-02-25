// Package bigquery provides BigQuery type provider for CEL type system integration.
package bigquery

import (
	"context"
	"errors"
	"fmt"
	"strings"

	bq "cloud.google.com/go/bigquery"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"

	"github.com/spandigital/cel2sql/v3/schema"
)

// Sentinel errors for the bigquery package.
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

// TypeProvider interface for BigQuery type providers.
type TypeProvider interface {
	types.Provider
	LoadTableSchema(ctx context.Context, tableName string) error
	GetSchemas() map[string]Schema
	Close()
}

type typeProvider struct {
	schemas   map[string]Schema
	client    *bq.Client
	datasetID string
}

// NewTypeProvider creates a new BigQuery type provider with pre-defined schemas.
func NewTypeProvider(schemas map[string]Schema) TypeProvider {
	return &typeProvider{schemas: schemas}
}

// NewTypeProviderWithClient creates a new BigQuery type provider that can introspect database schemas.
// The caller owns the *bigquery.Client and is responsible for closing it.
func NewTypeProviderWithClient(_ context.Context, client *bq.Client, datasetID string) (TypeProvider, error) {
	if client == nil {
		return nil, fmt.Errorf("%w: BigQuery client must not be nil", ErrInvalidSchema)
	}
	if datasetID == "" {
		return nil, fmt.Errorf("%w: dataset ID must not be empty", ErrInvalidSchema)
	}

	return &typeProvider{
		schemas:   make(map[string]Schema),
		client:    client,
		datasetID: datasetID,
	}, nil
}

// LoadTableSchema loads schema information for a table from BigQuery using the client API.
func (tp *typeProvider) LoadTableSchema(ctx context.Context, tableName string) error {
	if tp.client == nil {
		return fmt.Errorf("%w: no BigQuery client available", ErrInvalidSchema)
	}

	meta, err := tp.client.Dataset(tp.datasetID).Table(tableName).Metadata(ctx)
	if err != nil {
		return fmt.Errorf("%w: failed to get table metadata", ErrInvalidSchema)
	}

	fields := bigquerySchemaToFieldSchemas(meta.Schema)
	if len(fields) == 0 {
		return fmt.Errorf("%w: table %q has no columns", ErrInvalidSchema, tableName)
	}

	tp.schemas[tableName] = NewSchema(fields)
	return nil
}

// bigquerySchemaToFieldSchemas converts a BigQuery schema to a slice of FieldSchemas.
func bigquerySchemaToFieldSchemas(bqSchema bq.Schema) []FieldSchema {
	fields := make([]FieldSchema, 0, len(bqSchema))
	for _, f := range bqSchema {
		fields = append(fields, bigqueryFieldToFieldSchema(f))
	}
	return fields
}

// bigqueryFieldToFieldSchema converts a BigQuery FieldSchema to our FieldSchema.
func bigqueryFieldToFieldSchema(f *bq.FieldSchema) FieldSchema {
	typeName := bigqueryFieldTypeToString(f.Type)
	isJSON := f.Type == bq.JSONFieldType
	repeated := f.Repeated

	field := FieldSchema{
		Name:     f.Name,
		Type:     typeName,
		Repeated: repeated,
		IsJSON:   isJSON,
	}

	// Handle nested RECORD types recursively
	if f.Type == bq.RecordFieldType && len(f.Schema) > 0 {
		field.Schema = bigquerySchemaToFieldSchemas(f.Schema)
	}

	if repeated {
		field.Dimensions = 1
		field.ElementType = typeName
	}

	return field
}

// bigqueryFieldTypeToString converts a BigQuery FieldType to a string type name.
func bigqueryFieldTypeToString(ft bq.FieldType) string {
	return strings.ToLower(string(ft))
}

// Close is a no-op since we don't own the *bigquery.Client.
func (tp *typeProvider) Close() {
	// No-op: caller owns the *bigquery.Client connection
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

	exprType := bigqueryTypeToCELExprType(field)
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

// bigqueryTypeToCELExprType converts a BigQuery field schema to a CEL expression type.
func bigqueryTypeToCELExprType(field *schema.FieldSchema) *exprpb.Type {
	baseType := bigqueryBaseTypeToCEL(field.Type)
	if field.Repeated {
		return decls.NewListType(baseType)
	}
	return baseType
}

// bigqueryBaseTypeToCEL converts a BigQuery type name to a CEL expression type.
func bigqueryBaseTypeToCEL(typeName string) *exprpb.Type {
	switch typeName {
	case "STRING", "string":
		return decls.String
	case "INT64", "int64", "INTEGER", "integer":
		return decls.Int
	case "FLOAT64", "float64", "FLOAT", "float", "NUMERIC", "numeric":
		return decls.Double
	case "BOOL", "bool", "BOOLEAN", "boolean":
		return decls.Bool
	case "BYTES", "bytes":
		return decls.Bytes
	case "JSON", "json":
		return decls.Dyn
	case "TIMESTAMP", "timestamp":
		return decls.Timestamp
	default:
		return decls.Dyn
	}
}
