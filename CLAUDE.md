# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

cel2sql converts CEL (Common Expression Language) expressions to PostgreSQL SQL conditions. It specifically targets PostgreSQL standard SQL and was recently migrated from BigQuery.

**Module**: `github.com/spandigital/cel2sql/v2`
**Go Version**: 1.24+
**Current Version**: v2.8.0

## Common Development Commands

### Build and Test
```bash
make build              # Build the project
make test               # Run tests with race detection and coverage
make test-coverage      # Generate HTML coverage report
```

### Code Quality
```bash
make lint               # Run golangci-lint (required before commits)
make fmt                # Format code with go fmt and goimports
make vuln-check         # Run security vulnerability checks
make ci                 # Run all checks (fmt, lint, test, vuln-check)
```

### Dependencies
```bash
make deps               # Download and verify dependencies
make update-deps        # Update all dependencies
make install-tools      # Install development tools (golangci-lint, goimports, govulncheck)
```

### Run a single test
```bash
go test -v -run TestFunctionName ./...
```

## Core Architecture

### Main Components

1. **`cel2sql.go`** - Main conversion engine that transforms CEL AST to PostgreSQL SQL strings
   - Entry point: `Convert(ast *cel.Ast) (string, error)`
   - Uses visitor pattern to traverse CEL expression tree
   - Handles operators, functions, comprehensions, and type conversions

2. **`comprehensions.go`** - CEL comprehensions support (all, exists, exists_one, filter, map)
   - Converts to PostgreSQL UNNEST() patterns
   - Supports nested comprehensions
   - Works with both schema-based arrays and JSON arrays

3. **`json.go`** - JSON/JSONB field handling
   - Detects JSON columns and applies PostgreSQL path operators (->>, ?)
   - Supports nested JSON field access
   - Handles has() macro for JSON field existence checks

4. **`operators.go`** - Operator conversion logic
   - Maps CEL operators to PostgreSQL equivalents
   - Handles special cases (IS NULL, boolean comparisons, etc.)

5. **`timestamps.go`** - Timestamp and duration handling
   - Converts CEL timestamp operations to PostgreSQL TIMESTAMP operations
   - Handles INTERVAL conversions

6. **`pg/provider.go`** - PostgreSQL type provider for CEL type system
   - Maps PostgreSQL types to CEL types
   - Supports dynamic schema loading from live databases
   - Handles composite types and arrays

7. **`sqltypes/types.go`** - Custom SQL type definitions for CEL (DATE, TIME, DATETIME, INTERVAL)

### Type System Integration

The library uses CEL's protobuf-based type system (`exprpb.Type`, `exprpb.Expr`). PostgreSQL types are mapped to CEL types through `pg.TypeProvider`:

- `text` → `decls.String`
- `bigint` / `integer` → `decls.Int`
- `boolean` → `decls.Bool`
- `double precision` → `decls.Double`
- `timestamp with time zone` → `decls.Timestamp`
- `json` / `jsonb` → `decls.String` (with automatic JSON path support)
- Arrays: Set `Repeated: true` in schema
- Composite types: Use nested `Schema` fields

### JSON/JSONB Support

CEL field access on JSON/JSONB columns automatically converts to PostgreSQL JSON path operations:
- `user.preferences.theme` → `user.preferences->>'theme'`
- `has(user.preferences.theme)` → `user.preferences ? 'theme'`
- Works with nested paths: `user.profile.settings.key` → `user.profile->'settings'->>'key'`

Detection happens in `shouldUseJSONPath()` and `visitSelect()` functions in `json.go`.

### CEL Comprehensions Support

Full support for CEL comprehensions converted to PostgreSQL UNNEST patterns:
- `list.all(x, x > 0)` → `NOT EXISTS (SELECT 1 FROM UNNEST(list) AS x WHERE NOT (x > 0))`
- `list.exists(x, condition)` → `EXISTS (SELECT 1 FROM UNNEST(list) AS x WHERE condition)`
- `list.filter(x, condition)` → `ARRAY(SELECT x FROM UNNEST(list) AS x WHERE condition)`
- `list.map(x, transform)` → `ARRAY(SELECT transform FROM UNNEST(list) AS x)`
- Supports nested comprehensions and both schema arrays and JSON arrays

Pattern recognition and conversion logic is in `comprehensions.go`.

### Regex Pattern Matching (v2.8.0)

Supports CEL `matches()` function with automatic RE2 to POSIX regex conversion:
- `field.matches(r"pattern")` → `field ~ 'pattern'`
- `field.matches(r"(?i)pattern")` → `field ~* 'pattern'` (case-insensitive)
- Automatic conversion of RE2 patterns to PostgreSQL-compatible POSIX format

## Code Quality Requirements

### Pre-commit Checklist
1. Run `make fmt` - Format code
2. Run `make lint` - Must pass without errors (golangci-lint)
3. Run `make test` - All tests must pass
4. Common lint fixes:
   - Use `errors.New()` instead of `fmt.Errorf()` for static error messages
   - Rename unused parameters to `_`
   - Add comments for exported constants and types
   - Include package comments for main packages

### Testing Guidelines
- Use PostgreSQL schemas (`pg.Schema`) in tests, not BigQuery
- Use `pg.NewTypeProvider()` for schema definitions
- Include tests for nested types, arrays, and JSON fields
- Verify SQL output matches PostgreSQL syntax (single quotes, proper functions)
- Use testcontainers for integration tests with real PostgreSQL

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
provider, err := pg.NewTypeProviderWithConnection(ctx, connectionString)
if err != nil {
    return err
}
defer provider.Close()

err = provider.LoadTableSchema(ctx, "tableName")
```

### CEL Environment Setup
```go
env, err := cel.NewEnv(
    cel.CustomTypeProvider(provider),
    cel.Variable("table", cel.ObjectType("TableName")),
)
```

### Converting CEL to SQL
```go
ast, issues := env.Compile(`table.field == "value" && table.age > 30`)
if issues != nil && issues.Err() != nil {
    return issues.Err()
}

sqlCondition, err := cel2sql.Convert(ast)
// Returns: table.field = 'value' AND table.age > 30
```

## Important Notes

### Migration Context
This project was migrated from BigQuery to PostgreSQL in v2.0:
- All `cloud.google.com/go/bigquery` dependencies removed
- `bq/` package removed entirely
- PostgreSQL-specific syntax (single quotes, POSITION(), ARRAY_LENGTH(,1), etc.)
- Comprehensive JSON/JSONB support added
- Dynamic schema loading added

### Things to Avoid
- Do NOT add BigQuery dependencies back
- Do NOT remove protobuf dependencies (required by CEL)
- Do NOT use direct SQL string concatenation (use proper escaping)
- Do NOT ignore context cancellation in database operations

### When Adding Features
1. Consider PostgreSQL-specific SQL syntax
2. Add comprehensive tests with realistic schemas
3. Update type mappings in `pg/provider.go` if needed
4. Document new CEL operators/functions in README.md
5. Ensure backward compatibility
6. Run `make ci` before committing

## Project Structure
```
cel2sql/
├── cel2sql.go              # Main conversion engine
├── comprehensions.go       # CEL comprehensions support
├── json.go                 # JSON/JSONB handling
├── operators.go            # Operator conversion
├── timestamps.go           # Timestamp/duration handling
├── utils.go                # Utility functions
├── pg/                     # PostgreSQL type provider
│   └── provider.go
├── sqltypes/               # Custom SQL types for CEL
│   └── types.go
└── examples/               # Usage examples
    ├── basic/
    ├── comprehensions/
    └── load_table_schema/
```

Each example should be in its own subdirectory with `main.go` and `README.md`, runnable via `go run main.go`.
