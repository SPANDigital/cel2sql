// Package celprovider provides the shared cel-go types.Provider implementation
// embedded by the flat (non-nested-schema) SQL dialect type providers.
//
// The MySQL, SQLite, DuckDB, BigQuery and Spark providers previously each
// reimplemented the same types.Provider boilerplate over a map[string]Schema,
// differing only in their type-name → CEL mapping. Base centralises that
// boilerplate; a dialect embeds Base, supplies its Mapper and LoadTableSchema,
// and overrides Close only if it owns a connection.
//
// The PostgreSQL provider (pg) is intentionally not built on Base: it resolves
// nested/composite schemas through a dotted-path lookup and owns a connection
// pool, so it keeps its own implementation.
package celprovider

import (
	"cel.dev/cel-go/common/types"
	"cel.dev/cel-go/common/types/ref"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"

	"github.com/spandigital/cel2sql/v3/schema"
)

// TypeMapper converts a dialect field schema to a CEL expression type.
type TypeMapper func(*schema.FieldSchema) *exprpb.Type

// Base implements the portion of cel-go's types.Provider shared by the flat
// dialect type providers. Embed it and set Schemas and Mapper at construction.
type Base struct {
	// Schemas holds the table schemas known to the provider, keyed by table name.
	Schemas map[string]schema.Schema
	// Mapper converts a field schema to its CEL expression type.
	Mapper TypeMapper
}

// GetSchemas returns the schemas known to this provider.
func (b *Base) GetSchemas() map[string]schema.Schema {
	return b.Schemas
}

// EnumValue implements types.Provider. Schemas declare no enums.
func (b *Base) EnumValue(_ string) ref.Val {
	return types.NewErr("unknown enum value")
}

// FindIdent implements types.Provider. Schemas declare no identifiers.
func (b *Base) FindIdent(_ string) (ref.Val, bool) {
	return nil, false
}

// FindStructType implements types.Provider.
func (b *Base) FindStructType(structType string) (*types.Type, bool) {
	if _, ok := b.Schemas[structType]; ok {
		return types.NewObjectType(structType), true
	}
	return nil, false
}

// FindStructFieldNames implements types.Provider.
func (b *Base) FindStructFieldNames(structType string) ([]string, bool) {
	s, ok := b.Schemas[structType]
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

// FindStructFieldType implements types.Provider, mapping the field's SQL type
// to a CEL type via Mapper.
func (b *Base) FindStructFieldType(structType, fieldName string) (*types.FieldType, bool) {
	s, ok := b.Schemas[structType]
	if !ok {
		return nil, false
	}
	field, found := s.FindField(fieldName)
	if !found {
		return nil, false
	}

	celType, err := types.ExprTypeToType(b.Mapper(field))
	if err != nil {
		return nil, false
	}

	return &types.FieldType{
		Type: celType,
	}, true
}

// NewValue implements types.Provider. Schema-described types are not constructible.
func (b *Base) NewValue(_ string, _ map[string]ref.Val) ref.Val {
	return types.NewErr("unknown type in schema")
}

// Close is a no-op. Providers that own a connection override it.
func (b *Base) Close() {}
