package testcases

import (
	"github.com/spandigital/cel2sql/v3/pg"
	"github.com/spandigital/cel2sql/v3/schema"
)

// EnvDefault is the default environment setup name (basic types, no schema).
const EnvDefault = ""

// EnvWithSchema is an environment with a schema-based type provider.
const EnvWithSchema = "schema"

// EnvWithJSON is an environment with JSON/JSONB schema fields.
const EnvWithJSON = "json_schema"

// EnvWithTimestamp is an environment for timestamp operations.
const EnvWithTimestamp = "timestamp"

// NewPersonSchema returns a dialect-agnostic schema for the "person" table,
// suitable for basic, operator, and string tests.
func NewPersonSchema() schema.Schema {
	return schema.NewSchema([]schema.FieldSchema{
		{Name: "name", Type: "text"},
		{Name: "age", Type: "integer"},
		{Name: "adult", Type: "boolean"},
		{Name: "height", Type: "double precision"},
		{Name: "email", Type: "text"},
		{Name: "tags", Type: "text", Repeated: true},
		{Name: "scores", Type: "integer", Repeated: true},
	})
}

// NewPersonPGSchema returns a PostgreSQL-specific schema for the "person" table.
func NewPersonPGSchema() pg.Schema {
	return pg.NewSchema([]pg.FieldSchema{
		{Name: "name", Type: "text"},
		{Name: "age", Type: "integer"},
		{Name: "adult", Type: "boolean"},
		{Name: "height", Type: "double precision"},
		{Name: "email", Type: "text"},
		{Name: "tags", Type: "text", Repeated: true},
		{Name: "scores", Type: "integer", Repeated: true},
	})
}

// NewProductSchema returns a dialect-agnostic schema for the "product" table,
// with JSON/JSONB fields for JSON-related tests.
func NewProductSchema() schema.Schema {
	return schema.NewSchema([]schema.FieldSchema{
		{Name: "name", Type: "text"},
		{Name: "price", Type: "double precision"},
		{Name: "metadata", Type: "jsonb", IsJSON: true, IsJSONB: true},
		{Name: "attributes", Type: "json", IsJSON: true},
		{Name: "tags", Type: "jsonb", IsJSON: true, IsJSONB: true, Repeated: true, ElementType: "text"},
	})
}

// NewProductPGSchema returns a PostgreSQL-specific schema for the "product" table.
func NewProductPGSchema() pg.Schema {
	return pg.NewSchema([]pg.FieldSchema{
		{Name: "name", Type: "text"},
		{Name: "price", Type: "double precision"},
		{Name: "metadata", Type: "jsonb", IsJSON: true, IsJSONB: true},
		{Name: "attributes", Type: "json", IsJSON: true},
		{Name: "tags", Type: "jsonb", IsJSON: true, IsJSONB: true, Repeated: true, ElementType: "text"},
	})
}

// NewOrderSchema returns a dialect-agnostic schema for the "orders" table,
// with array and timestamp fields.
func NewOrderSchema() schema.Schema {
	return schema.NewSchema([]schema.FieldSchema{
		{Name: "order_id", Type: "bigint"},
		{Name: "customer_name", Type: "text"},
		{Name: "total", Type: "double precision"},
		{Name: "items", Type: "text", Repeated: true},
		{Name: "created_at", Type: "timestamp with time zone"},
		{Name: "status", Type: "text"},
	})
}

// NewOrderPGSchema returns a PostgreSQL-specific schema for the "orders" table.
func NewOrderPGSchema() pg.Schema {
	return pg.NewSchema([]pg.FieldSchema{
		{Name: "order_id", Type: "bigint"},
		{Name: "customer_name", Type: "text"},
		{Name: "total", Type: "double precision"},
		{Name: "items", Type: "text", Repeated: true},
		{Name: "created_at", Type: "timestamp with time zone"},
		{Name: "status", Type: "text"},
	})
}
