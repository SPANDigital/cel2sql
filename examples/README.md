# cel2sql Examples

This directory contains examples demonstrating various features and use cases of cel2sql.

## Available Examples

### 1. [Basic Example](./basic/)
Demonstrates the fundamental usage of cel2sql:
- Manual schema definition
- CEL environment setup
- Expression compilation and SQL conversion
- Basic PostgreSQL type mappings

**Run**: `cd basic && go run main.go`

### 2. [LoadTableSchema Example](./load_table_schema/)
Demonstrates dynamic schema loading from PostgreSQL databases:
- Connecting to PostgreSQL databases
- Loading table schemas dynamically using `LoadTableSchema`
- Using loaded schemas in CEL expressions
- Integration testing patterns

**Run**: `cd load_table_schema && go run main.go`

## Getting Started

Each example directory contains:
- `main.go` - The runnable example code
- `README.md` - Detailed documentation and explanation

### Prerequisites

- Go 1.24 or later
- For the LoadTableSchema example: PostgreSQL database (optional for static demo)

### Running Examples

Navigate to any example directory and run:

```bash
go run main.go
```

## Example Progression

We recommend exploring the examples in this order:

1. **Basic** - Start here to understand the core concepts
2. **LoadTableSchema** - Learn about dynamic schema loading and database integration

## Key Concepts Demonstrated

- **Schema Definition**: Both manual and dynamic schema creation
- **Type Safety**: CEL type checking with PostgreSQL types
- **SQL Generation**: Converting CEL expressions to PostgreSQL-compatible SQL
- **PostgreSQL Integration**: Working with real database connections
- **Best Practices**: Error handling, resource management, and testing patterns

## Adding New Examples

When adding new examples:

1. Create a new subdirectory under `examples/`
2. Include a `main.go` with a complete, runnable example
3. Add a comprehensive `README.md` explaining the example
4. Update this main README to include the new example
5. Ensure the example follows the established patterns and includes proper error handling

## Related Documentation

- [Main README](../README.md) - Project overview and usage
- [Contributing Guide](../CONTRIBUTING.md) - Development setup and testing patterns
- [API Documentation](https://pkg.go.dev/github.com/spandigital/cel2sql/v2) - Detailed API reference
