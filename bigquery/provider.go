// Package bigquery provides BigQuery type provider for CEL type system integration.
package bigquery

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"cel.dev/cel-go/checker/decls"
	"cel.dev/cel-go/common/types"
	bq "cloud.google.com/go/bigquery"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"

	"github.com/spandigital/cel2sql/v3/internal/celprovider"
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
	celprovider.Base
	client    *bq.Client
	datasetID string
}

// NewTypeProvider creates a new BigQuery type provider with pre-defined schemas.
func NewTypeProvider(schemas map[string]Schema) TypeProvider {
	return &typeProvider{Base: celprovider.Base{Schemas: schemas, Mapper: bigqueryTypeToCELExprType}}
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
		Base:      celprovider.Base{Schemas: make(map[string]Schema), Mapper: bigqueryTypeToCELExprType},
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

	tp.Schemas[tableName] = NewSchema(fields)
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
