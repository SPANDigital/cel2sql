// Package opensearch provides an OpenSearch type provider for CEL type system integration.
//
// An OpenSearch index mapping is loaded with GET /<index>/_mapping and turned
// into a schema: scalar field types map to CEL scalars, object fields become
// nested schemas addressed by dot path (p.dims.width) and nested fields become
// lists of nested schemas that CEL comprehensions iterate
// (p.variants.exists(v, v.sku == "A")).
package opensearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"cel.dev/cel-go/checker/decls"
	"cel.dev/cel-go/common/types"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"

	"github.com/spandigital/cel2sql/v3/internal/celprovider"
	"github.com/spandigital/cel2sql/v3/schema"
)

// Sentinel errors for the opensearch package.
var (
	// ErrInvalidSchema indicates a problem with the provided schema or index introspection.
	ErrInvalidSchema = errors.New("invalid schema")
)

// FieldSchema is an alias for schema.FieldSchema.
type FieldSchema = schema.FieldSchema

// Schema is an alias for schema.Schema.
type Schema = schema.Schema

// NewSchema creates a new Schema. This is an alias for schema.NewSchema.
var NewSchema = schema.NewSchema

// validIndexName matches OpenSearch index naming rules: lowercase, must not
// start with _, - or +, and none of the characters \ / * ? " < > | , # : or
// whitespace. The name is only ever used as a URL path segment, but it is
// validated anyway so that an alias pattern or wildcard never reaches the API.
var validIndexName = regexp.MustCompile(`^[a-z0-9][a-z0-9_.+-]{0,254}$`)

// maxMappingDepth bounds recursion through object/nested properties.
const maxMappingDepth = 32

// TypeProvider is the OpenSearch type provider interface.
type TypeProvider interface {
	types.Provider
	LoadTableSchema(ctx context.Context, indexName string) error
	GetSchemas() map[string]Schema
	Close()
}

type typeProvider struct {
	celprovider.Base
	client *opensearchapi.Client
}

// NewTypeProvider creates a new OpenSearch type provider with pre-defined schemas.
func NewTypeProvider(schemas map[string]Schema) TypeProvider {
	return &typeProvider{Base: celprovider.Base{Schemas: schemas, Mapper: opensearchTypeToCELExprType}}
}

// NewTypeProviderWithClient creates a new OpenSearch type provider that can
// introspect index mappings. The caller owns the *opensearchapi.Client and is
// responsible for its lifetime; Close is a no-op.
func NewTypeProviderWithClient(_ context.Context, client *opensearchapi.Client) (TypeProvider, error) {
	if client == nil {
		return nil, fmt.Errorf("%w: OpenSearch client must not be nil", ErrInvalidSchema)
	}
	return &typeProvider{
		Base:   celprovider.Base{Schemas: make(map[string]Schema), Mapper: opensearchTypeToCELExprType},
		client: client,
	}, nil
}

// LoadTableSchema loads the mapping of an index (GET /<index>/_mapping) and
// registers it as a schema under the index name.
func (tp *typeProvider) LoadTableSchema(ctx context.Context, indexName string) error {
	if !validIndexName.MatchString(indexName) {
		return fmt.Errorf("%w: invalid index name %q", ErrInvalidSchema, indexName)
	}
	if tp.client == nil {
		return fmt.Errorf("%w: no OpenSearch client available", ErrInvalidSchema)
	}

	resp, err := tp.client.Indices.Mapping.Get(ctx, &opensearchapi.MappingGetReq{Indices: []string{indexName}})
	if err != nil {
		return fmt.Errorf("%w: failed to get mapping for index %q: %w", ErrInvalidSchema, indexName, err)
	}

	indices := resp.GetIndices()
	entry, ok := indices[indexName]
	if !ok {
		// The name may be an alias; the response is keyed by the concrete index.
		if len(indices) != 1 {
			return fmt.Errorf("%w: mapping response for %q contains %d indices", ErrInvalidSchema, indexName, len(indices))
		}
		for _, entry = range indices {
			break
		}
	}

	fields, err := parseMapping(entry.Mappings)
	if err != nil {
		return fmt.Errorf("%w: index %q: %w", ErrInvalidSchema, indexName, err)
	}
	if len(fields) == 0 {
		return fmt.Errorf("%w: index %q has no mapped fields", ErrInvalidSchema, indexName)
	}

	tp.Schemas[indexName] = NewSchema(fields)
	return nil
}

// mappingNode is the subset of an OpenSearch field mapping needed to build a schema.
type mappingNode struct {
	Type       string                 `json:"type"`
	Properties map[string]mappingNode `json:"properties"`
}

// parseMapping converts the "mappings" object of a GET _mapping response
// into field schemas. Multi-fields ("fields") are ignored, properties are
// emitted in name order, and recursion is capped at maxMappingDepth.
func parseMapping(raw json.RawMessage) ([]FieldSchema, error) {
	if len(raw) == 0 {
		return nil, errors.New("empty mapping")
	}
	var root mappingNode
	if err := json.Unmarshal(raw, &root); err != nil {
		return nil, fmt.Errorf("malformed mapping: %w", err)
	}
	return propertiesToFields(root.Properties, 0)
}

func propertiesToFields(props map[string]mappingNode, depth int) ([]FieldSchema, error) {
	if depth > maxMappingDepth {
		return nil, fmt.Errorf("mapping nesting exceeds %d levels", maxMappingDepth)
	}
	names := make([]string, 0, len(props))
	for name := range props {
		names = append(names, name)
	}
	sort.Strings(names)

	fields := make([]FieldSchema, 0, len(names))
	for _, name := range names {
		node := props[name]
		field := FieldSchema{Name: name, Type: strings.ToLower(node.Type)}

		switch field.Type {
		case "nested":
			children, err := propertiesToFields(node.Properties, depth+1)
			if err != nil {
				return nil, err
			}
			field.Repeated = true
			field.Dimensions = 1
			field.ElementType = "nested"
			field.Schema = children
		case "object", "":
			// A field with properties and no explicit type is an object.
			field.Type = "object"
			children, err := propertiesToFields(node.Properties, depth+1)
			if err != nil {
				return nil, err
			}
			field.Schema = children
		case "flat_object":
			field.IsJSON = true
		}

		fields = append(fields, field)
	}
	return fields, nil
}

// findSchema resolves a dotted struct name (index, or index.field.subfield)
// to the schema at that path.
func (tp *typeProvider) findSchema(typeName string) (Schema, bool) {
	parts := strings.Split(typeName, ".")
	s, found := tp.Schemas[parts[0]]
	if !found {
		return Schema{}, false
	}
	if len(parts) == 1 {
		return s, true
	}

	current := s
	for _, part := range parts[1:] {
		field, ok := current.FindField(part)
		if !ok || len(field.Schema) == 0 {
			return Schema{}, false
		}
		current = NewSchema(field.Schema)
	}
	return current, true
}

// FindStructType implements types.Provider, resolving dotted paths into
// object and nested fields.
func (tp *typeProvider) FindStructType(structType string) (*types.Type, bool) {
	if _, found := tp.findSchema(structType); !found {
		return nil, false
	}
	return types.NewObjectType(structType), true
}

// FindStructFieldNames implements types.Provider.
func (tp *typeProvider) FindStructFieldNames(structType string) ([]string, bool) {
	s, found := tp.findSchema(structType)
	if !found {
		return nil, false
	}
	fields := s.Fields()
	names := make([]string, len(fields))
	for i, f := range fields {
		names[i] = f.Name
	}
	return names, true
}

// FindStructFieldType implements types.Provider. Fields with a nested schema
// are typed as <structType>.<fieldName> objects, wrapped in a list when the
// field is a nested (multi-valued) mapping.
func (tp *typeProvider) FindStructFieldType(structType, fieldName string) (*types.FieldType, bool) {
	s, found := tp.findSchema(structType)
	if !found {
		return nil, false
	}
	field, found := s.FindField(fieldName)
	if !found {
		return nil, false
	}

	var exprType *exprpb.Type
	if len(field.Schema) > 0 {
		exprType = decls.NewObjectType(structType + "." + fieldName)
		if field.Repeated {
			exprType = decls.NewListType(exprType)
		}
	} else {
		exprType = opensearchTypeToCELExprType(field)
	}

	celType, err := types.ExprTypeToType(exprType)
	if err != nil {
		return nil, false
	}
	return &types.FieldType{Type: celType}, true
}

// opensearchTypeToCELExprType converts an OpenSearch field schema to a CEL expression type.
func opensearchTypeToCELExprType(field *schema.FieldSchema) *exprpb.Type {
	baseType := opensearchBaseTypeToCEL(field.Type)
	if field.Repeated {
		return decls.NewListType(baseType)
	}
	return baseType
}

// opensearchBaseTypeToCEL converts an OpenSearch mapping type name to a CEL expression type.
func opensearchBaseTypeToCEL(typeName string) *exprpb.Type {
	switch strings.ToLower(typeName) {
	case "keyword", "text", "ip", "constant_keyword", "wildcard", "match_only_text", "search_as_you_type":
		return decls.String
	case "byte", "short", "integer", "long", "unsigned_long":
		return decls.Int
	case "float", "half_float", "scaled_float", "double":
		return decls.Double
	case "boolean":
		return decls.Bool
	case "binary":
		return decls.Bytes
	case "date", "date_nanos":
		return decls.Timestamp
	default:
		// object without known properties, flat_object, geo types, …
		return decls.Dyn
	}
}
