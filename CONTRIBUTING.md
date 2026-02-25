# Contributing to cel2sql

Thank you for your interest in contributing to cel2sql! This document provides guidelines and information for contributors.

## Development Setup

### Prerequisites

- Go 1.24 or later
- Git
- Make (optional, but recommended)

### Getting Started

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/your-username/cel2sql.git
   cd cel2sql
   ```

3. Install development tools:
   ```bash
   make install-tools
   ```

4. Download dependencies:
   ```bash
   make deps
   ```

5. Run tests to ensure everything works:
   ```bash
   make test
   ```

## Development Workflow

### Code Style

We use standard Go formatting and linting tools:

- **Format**: `make fmt` - formats code using `go fmt` and `goimports`
- **Lint**: `make lint` - runs `golangci-lint` with our configuration
- **Test**: `make test` - runs all tests with race detection and coverage

### Before Submitting

Run the full CI pipeline locally:
```bash
make ci
```

This will run formatting, linting, testing, and vulnerability checks.

### Testing

- Write tests for new functionality
- Ensure all tests pass: `make test`
- Check test coverage: `make test-coverage`
- Tests should use PostgreSQL schemas (not BigQuery)

Example test structure:
```go
func TestNewFeature(t *testing.T) {
    schema := pg.Schema{
        {Name: "field_name", Type: "text"},
        {Name: "array_field", Type: "text", Repeated: true},
    }
    provider := pg.NewTypeProvider(map[string]pg.Schema{"TableName": schema})
    
    // Test your feature
}
```

#### Integration Testing with PostgreSQL Testcontainers

For integration tests that require a real PostgreSQL database, use testcontainers:

```go
func TestLoadTableSchema_WithPostgresContainer(t *testing.T) {
    ctx := context.Background()

    // Create a PostgreSQL container
    container, err := postgres.Run(ctx, 
        "postgres:15",
        postgres.WithDatabase("testdb"),
        postgres.WithUsername("testuser"),
        postgres.WithPassword("testpass"),
        postgres.WithInitScripts("create_test_table.sql"),
        testcontainers.WithWaitStrategy(
            wait.ForLog("database system is ready to accept connections").
                WithOccurrence(2).
                WithStartupTimeout(time.Second * 60),
        ),
    )
    require.NoError(t, err)
    defer container.Terminate(ctx)

    // Get connection string
    connStr, err := container.ConnectionString(ctx, "sslmode=disable")
    require.NoError(t, err)

    // Create type provider with database connection
    provider, err := pg.NewTypeProviderWithConnection(ctx, connStr)
    require.NoError(t, err)
    defer provider.Close()

    // Test LoadTableSchema
    err = provider.LoadTableSchema(ctx, "users")
    require.NoError(t, err)

    // Verify the schema was loaded correctly
    foundType, found := provider.FindStructType("users")
    assert.True(t, found)
    assert.NotNil(t, foundType)
}
```

Required imports for testcontainer tests:
```go
import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
    "github.com/testcontainers/testcontainers-go/wait"
)
```

### Adding New CEL Functions

1. Add the function mapping in `cel2sql.go`
2. Add comprehensive tests in `cel2sql_test.go`
3. Update the README with documentation
4. Ensure the function works with dialect abstraction

### Multi-Dialect Architecture

cel2sql supports PostgreSQL (default), MySQL, SQLite, DuckDB, and BigQuery. When adding features:

- Call `con.dialect.*` methods for any SQL that differs between databases
- Standard SQL (AND, OR, =, !=, etc.) stays inline in the converter
- Add expected SQL for all dialects in `testcases/*.go`
- Run `make ci` to verify all dialects pass

### Adding a New Dialect

To add support for a new SQL dialect:

1. **Create the dialect package**: `dialect/<name>/dialect.go`
   - Implement the `dialect.Dialect` interface (~40 methods)
   - Register with `dialect.Register()` in `init()`

2. **Create regex conversion** (if applicable): `dialect/<name>/regex.go`
   - Convert RE2 patterns to the dialect's regex format
   - Include ReDoS protection (pattern length, nesting limits)

3. **Create validation**: `dialect/<name>/validation.go`
   - Field name validation, reserved keywords

4. **Create type provider**: `<name>/provider.go`
   - Map native database types to CEL types

5. **Add env factory**: `testutil/env.go`
   - Add `<Name>EnvFactory()` function
   - Update `DialectEnvFactory()` switch

6. **Add test runner**: `testutil/runner_<name>_test.go`

7. **Add expected SQL to all test case files** in `testcases/`

8. **Update `dialect/dialect.go`** to add the dialect name constant

## Pull Request Process

1. Create a feature branch from `main`
2. Make your changes
3. Run `make ci` to ensure everything passes
4. Update documentation if needed
5. Create a pull request with:
   - Clear description of changes
   - Link to any related issues
   - Test results

### PR Requirements

- [ ] Tests pass locally
- [ ] Code is formatted (`make fmt`)
- [ ] Linting passes (`make lint`)
- [ ] No security vulnerabilities (`make vuln-check`)
- [ ] Documentation updated (if applicable)
- [ ] Commit messages follow conventional format

## Code Organization

### Project Structure

```
cel2sql/
├── cel2sql.go              # Main conversion engine (uses dialect interface)
├── cel2sql_test.go         # Main tests
├── dialect/                # Dialect interface + implementations
│   ├── dialect.go          # Interface definition + Name type
│   ├── registry.go         # Name→Dialect lookup
│   ├── postgres/           # PostgreSQL dialect
│   ├── mysql/              # MySQL dialect
│   ├── sqlite/             # SQLite dialect
│   ├── duckdb/             # DuckDB dialect
│   └── bigquery/           # BigQuery dialect
├── pg/                     # PostgreSQL type provider
├── mysql/                  # MySQL type provider
├── sqlite/                 # SQLite type provider
├── duckdb/                 # DuckDB type provider
├── bigquery/               # BigQuery type provider
├── schema/                 # Dialect-agnostic schema types
├── sqltypes/               # Custom SQL types for CEL
├── testcases/              # Shared test cases with per-dialect expected SQL
├── testutil/               # Test runner + env factories
└── examples/               # Usage examples
```

### Key Components

- **cel2sql.go**: Core conversion logic from CEL AST to SQL (calls dialect methods)
- **dialect/dialect.go**: Dialect interface defining all SQL generation points
- **dialect/*/dialect.go**: Per-dialect SQL generation implementations
- **pg/provider.go**, **mysql/provider.go**, etc.: Type system integration per dialect
- **sqltypes/types.go**: Custom SQL type definitions for CEL
- **testcases/*.go**: Shared test cases with expected SQL for all dialects

## Debugging

### Common Issues

1. **Type resolution**: Check `typeMap` in converter
2. **PostgreSQL arrays**: Use `[]` suffix in type names
3. **Composite types**: Ensure proper nested schema navigation

### Debug Tips

- Use `cel.AstToCheckedExpr()` to inspect CEL AST
- Check type mappings in `pg.TypeProvider`
- Validate schema definitions match PostgreSQL structure

## Security

- Use parameterized queries when integrating with databases
- Validate CEL expressions before conversion
- Sanitize field names and table names in SQL output
- Be cautious with user-provided schema definitions

## Release Process

Releases are automated through GitHub Actions:

1. Create and push a tag: `git tag v1.2.3 && git push origin v1.2.3`
2. GitHub Actions will run tests and create a release
3. Update changelog and release notes

## Getting Help

- Check existing issues and discussions
- Create an issue for bugs or feature requests
- Ask questions in discussions

## Code of Conduct

- Be respectful and inclusive
- Focus on constructive feedback
- Help maintain a welcoming environment

Thank you for contributing to cel2sql!
