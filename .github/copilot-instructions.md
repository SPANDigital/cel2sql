# Copilot Instructions for cel2sql

## Project Overview

This project converts [CEL (Common Expression Language)](https://opensource.google/projects/cel) expressions to PostgreSQL SQL conditions. It was recently migrated from BigQuery to PostgreSQL using the latest pgx v5 driver.

## Key Architecture

### Core Components

1. **`cel2sql.go`** - Main conversion engine that transforms CEL AST to SQL strings
2. **`pg/provider.go`** - PostgreSQL type provider for CEL type system integration
3. **`sqltypes/types.go`** - Custom SQL type definitions for CEL (Date, Time, DateTime)
4. **`test/testdata.go`** - PostgreSQL schema definitions for testing
5. **`json.go`** - JSON/JSONB nested path handling and PostgreSQL operator generation

### JSON/JSONB Support

The system provides comprehensive support for nested JSON/JSONB field access:

- **Nested Path Generation**: Converts CEL expressions like `table.metadata.corpus.section` to PostgreSQL JSON paths `table.metadata->'corpus'->>'section'`
- **Automatic Type Casting**: Detects numeric comparisons and adds `::numeric` casting where needed
- **Array Operations**: Handles array membership with `jsonb_array_elements_text()` functions
- **Mixed JSON Types**: Supports both JSON and JSONB columns with appropriate operators

#### JSON Path Operators

- `->` - JSON object field access (returns JSON)
- `->>` - JSON object field access (returns text)
- Rule: Use `->` for intermediate navigation, `->>` for final text extraction

#### Supported Patterns

```cel
// String comparisons
information_assets.metadata.corpus.section == "Getting Started"
// → information_assets.metadata->'corpus'->>'section' = 'Getting Started'

// Numeric comparisons (with automatic casting)
information_assets.metadata.version.major > 1
// → (information_assets.metadata->'version'->>'major')::numeric > 1

// Array membership
"documentation" in information_assets.metadata.corpus.tags
// → 'documentation' = ANY(ARRAY(SELECT jsonb_array_elements_text(...)))

// Complex nested conditions
info.metadata.corpus.section == "Guide" && info.metadata.version.major == 2
```

### Type System Integration

- Uses CEL's protobuf-based type system (`exprpb.Type`, `exprpb.Expr`)
- Maps PostgreSQL types to CEL types through the `pg.TypeProvider`
- Supports composite types, arrays, and nested schemas
- Handles JSON/JSONB fields with automatic path detection and operator selection

## Development Guidelines

### Code Style

- Use Go 1.24+ features
- Follow standard Go naming conventions
- Prefer explicit error handling over panics
- Use context.Context for database operations

### PostgreSQL Integration

- Always use `pgxpool.Pool` for connection pooling
- Map PostgreSQL types properly:
  - `text` → `decls.String`
  - `bigint` → `decls.Int`
  - `boolean` → `decls.Bool`
  - `double precision` → `decls.Double`
  - `timestamp with time zone` → `decls.Timestamp`
- Support arrays with `Repeated: true`
- Handle composite types with nested `Schema` fields

### Testing

- Test files should use PostgreSQL schemas, not BigQuery
- Use `pg.NewTypeProvider()` with `pg.Schema` definitions
- Include tests for nested types and arrays
- Verify SQL output matches PostgreSQL syntax
- Test JSON/JSONB nested path expressions with `testcontainers-go`
- Include comprehensive tests for complex nested structures

### Dependencies

- **CEL**: `github.com/google/cel-go` - Core CEL functionality
- **PostgreSQL**: `github.com/jackc/pgx/v5` - Database driver
- **Protobuf**: Required for CEL (don't remove these dependencies)
- **Testing**: `github.com/stretchr/testify`
- **Test Containers**: `github.com/testcontainers/testcontainers-go` - PostgreSQL integration tests

## Common Patterns

### Creating Type Providers

```go
schema := pg.Schema{
    {Name: "field_name", Type: "text", Repeated: false},
    {Name: "array_field", Type: "text", Repeated: true},
    {Name: "composite_field", Type: "composite", Schema: []pg.FieldSchema{...}},
    {Name: "json_field", Type: "jsonb", Repeated: false}, // For JSON/JSONB fields
}
provider := pg.NewTypeProvider(map[string]pg.Schema{"TableName": schema})
```

### CEL Environment Setup

```go
env, err := cel.NewEnv(
    cel.CustomTypeProvider(provider),
    cel.Variable("table", cel.ObjectType("TableName")),
)
```

### JSON/JSONB Field Configuration

```go
// Configure tables with JSON/JSONB fields for nested path detection
jsonFields := map[string][]string{
    "information_assets": {"metadata", "properties", "classification"},
    "documents": {"content", "structure", "taxonomy", "analytics"},
}
```

### Adding New SQL Functions

1. Add function mapping in `cel2sql.go` conversion logic
2. Add corresponding tests in `cel2sql_test.go`
3. Update README documentation

## Migration Context

This project was recently migrated from BigQuery to PostgreSQL:

- **Removed**: All `cloud.google.com/go/bigquery` dependencies
- **Removed**: `bq/` package entirely
- **Added**: `pg/` package with PostgreSQL-specific logic
- **Updated**: All tests to use PostgreSQL schemas
- **Updated**: Documentation to reflect PostgreSQL usage

## Things to Avoid

- **Don't** add BigQuery dependencies back
- **Don't** remove protobuf dependencies (required by CEL)
- **Don't** use direct SQL string concatenation (use proper escaping)
- **Don't** ignore context cancellation in database operations

## When Adding Features

1. Consider PostgreSQL-specific SQL syntax differences
2. Add comprehensive tests with realistic PostgreSQL schemas
3. Update type mappings in `pg/provider.go` if needed
4. Document new CEL operators/functions in README
5. Ensure backward compatibility with existing CEL expressions

## Debugging Tips

- Use `cel.AstToCheckedExpr()` to inspect CEL AST structure
- Check `typeMap` in converter for type resolution issues
- PostgreSQL arrays use `[]` suffix in type names
- Composite types require proper nested schema navigation
- For JSON path issues, check `shouldUseJSONPath()` and `hasJSONFieldInChain()` functions
- Verify JSON field detection in `isJSONArrayField()` and `isJSONBField()` functions

## Security Considerations

- Always use parameterized queries when integrating with actual databases
- Validate CEL expressions before conversion
- Sanitize field names and table names in SQL output
- Be cautious with user-provided schema definitions
