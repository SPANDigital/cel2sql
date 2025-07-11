# Copilot Instructions for cel2sql

## Repository Information

**Repository Owner**: SPANDigital  
**Repository URL**: https://github.com/SPANDigital/cel2sql  
**Maintainer**: Richard Wooding (@richardwooding)

## Project Overview

This project converts [CEL (Common Expression Language)](https://opensource.google/projects/cel) expressions to PostgreSQL SQL conditions. It was recently migrated from BigQuery to PostgreSQL using the latest pgx v5 driver.

## Key Architecture

### Core Components

1. **`cel2sql.go`** - Main conversion engine that transforms CEL AST to SQL strings
2. **`pg/provider.go`** - PostgreSQL type provider for CEL type system integration
3. **`sqltypes/types.go`** - Custom SQL type definitions for CEL (Date, Time, DateTime)
4. **`test/testdata.go`** - PostgreSQL schema definitions for testing

### Type System Integration

- Uses CEL's protobuf-based type system (`exprpb.Type`, `exprpb.Expr`)
- Maps PostgreSQL types to CEL types through the `pg.TypeProvider`
- Supports composite types, arrays, and nested schemas

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
  - `json` → `decls.String` (with JSON path support)
  - `jsonb` → `decls.String` (with JSON path support)
- Support arrays with `Repeated: true`
- Handle composite types with nested `Schema` fields
- JSON/JSONB fields support PostgreSQL path operations (`->>`)

### JSON/JSONB Support

- CEL expressions like `user.preferences.theme` automatically convert to `user.preferences->>'theme'`
- The converter detects JSON/JSONB columns and applies proper PostgreSQL syntax
- Nested JSON access is supported: `user.profile.settings.key` → `user.profile->>'settings'->>'key'`
- JSON field detection happens in `shouldUseJSONPath()` and `visitSelect()` functions

### CEL Comprehensions Support

- **Full comprehension support**: `all()`, `exists()`, `exists_one()`, `filter()`, `map()`
- **PostgreSQL UNNEST integration**: All comprehensions use `UNNEST()` for array processing
- **Pattern recognition**: `comprehensions.go` handles AST pattern matching for comprehension types
- **Nested comprehensions**: Support for complex nested operations
- **Schema integration**: Works with `pg.Schema` including array fields and composite types

### Testing

- Test files should use PostgreSQL schemas, not BigQuery
- Use `pg.NewTypeProvider()` with `pg.Schema` definitions
- Include tests for nested types and arrays
- Verify SQL output matches PostgreSQL syntax
- Use testcontainers for integration testing

### Dependencies

- **CEL**: `github.com/google/cel-go` - Core CEL functionality
- **PostgreSQL**: `github.com/jackc/pgx/v5` - Database driver
- **Protobuf**: Required for CEL (don't remove these dependencies)
- **Testing**: `github.com/stretchr/testify`
- **Containers**: `github.com/testcontainers/testcontainers-go`

## Common Patterns

### Creating Type Providers

```go
schema := pg.Schema{
    {Name: "field_name", Type: "text", Repeated: false},
    {Name: "array_field", Type: "text", Repeated: true},
    {Name: "json_field", Type: "jsonb", Repeated: false},
    {Name: "composite_field", Type: "composite", Schema: []pg.FieldSchema{...}},
}
provider := pg.NewTypeProvider(map[string]pg.Schema{"TableName": schema})
```

### Dynamic Schema Loading

```go
// Load schema from PostgreSQL database
provider, err := pg.NewTypeProviderWithConnection(ctx, connectionString)
if err != nil {
    return err
}
defer provider.Close()

// Load specific table schema
err = provider.LoadTableSchema(ctx, "tableName")
if err != nil {
    return err
}
```

### CEL Environment Setup

```go
env, err := cel.NewEnv(
    cel.CustomTypeProvider(provider),
    cel.Variable("table", cel.ObjectType("TableName")),
)
```

### Adding New SQL Functions

1. Add function mapping in `cel2sql.go` conversion logic
2. Add corresponding tests in `cel2sql_test.go`
3. Update README documentation

## Migration Context

This project was recently migrated from BigQuery to PostgreSQL and modernized:

- **Removed**: All `cloud.google.com/go/bigquery` dependencies
- **Removed**: `bq/` package entirely
- **Added**: `pg/` package with PostgreSQL-specific logic
- **Updated**: All tests to use PostgreSQL schemas and testcontainers
- **Updated**: Documentation to reflect PostgreSQL usage
- **Added**: Comprehensive JSON/JSONB support with path operations
- **Enhanced**: Type system with dynamic schema loading
- **Improved**: SQL generation with PostgreSQL-specific syntax

## Current Version Features (v2.4.0)

- **JSON/JSONB Support**: Full PostgreSQL JSON path operations
- **Dynamic Schema Loading**: Load table schemas from live PostgreSQL databases
- **Enhanced Testing**: Comprehensive testcontainer integration tests
- **PostgreSQL Optimized**: Single quotes, POSITION(), ARRAY_LENGTH(,1), etc.
- **Type Safety**: Improved type mappings and error handling

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

## Security Considerations

- Always use parameterized queries when integrating with actual databases
- Validate CEL expressions before conversion
- Sanitize field names and table names in SQL output
- Be cautious with user-provided schema definitions

## Release Process

1. Run full test suite including integration tests
2. Update version in `go.mod` if needed
3. Create release notes documenting changes
4. Tag release following semantic versioning
5. Update documentation as needed

## Contact

For questions or issues, contact Richard Wooding or create an issue on the SPANDigital/cel2sql repository.

## Project Structure

- **Root files**: Core library files (`cel2sql.go`, `comprehensions.go`)
- **`pg/`**: PostgreSQL-specific type provider and schema handling
- **`sqltypes/`**: Custom SQL type definitions for CEL integration
- **`test/`**: Test data and schema definitions
- **`examples/`**: Example implementations in separate directories:
  - `examples/basic/`: Basic usage examples
  - `examples/load_table_schema/`: Dynamic schema loading examples
  - `examples/comprehensions/`: CEL comprehensions examples
  - Each example should be in its own directory with a `main.go` and `README.md`

### Example Directory Guidelines

- Each example must be in its own subdirectory under `examples/`
- Main file should be named `main.go` (not named after the feature)
- Include a comprehensive `README.md` explaining the example
- Examples should be runnable with `go run main.go` from their directory
- Document expected output and key concepts demonstrated

### Code Quality and Linting

- **golangci-lint**: All code must pass `golangci-lint run` without issues
- **Required before commits**: Run `golangci-lint run` and fix all issues
- **Common linting rules**:
  - Use `errors.New()` instead of `fmt.Errorf()` for static error messages
  - Rename unused parameters to `_` (e.g., `func foo(used string, _ string)`)
  - Add comments for exported constants and types
  - Include package comments for main packages in examples
- **Formatting**: Always run `go fmt ./...` before committing
- **Static analysis**: Ensure `go vet ./...` passes without warnings
