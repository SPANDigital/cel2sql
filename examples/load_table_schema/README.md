# LoadTableSchema Example

This example demonstrates how to use the `LoadTableSchema` functionality to dynamically load PostgreSQL table schemas and use them with CEL expressions.

## Features Demonstrated

- **Static Schema Definition**: Manual definition of PostgreSQL table schemas
- **Dynamic Schema Loading**: Loading table schemas directly from a PostgreSQL database
- **CEL Integration**: Using loaded schemas in CEL expressions
- **Type Conversion**: Converting CEL expressions to SQL conditions

## Running the Example

### Prerequisites

- Go 1.24 or later
- PostgreSQL database (for dynamic schema loading example)

### Static Schema Example

```bash
go run main.go
```

This will demonstrate CEL to SQL conversion using a pre-defined schema.

### Dynamic Schema Example

1. Set up a PostgreSQL database with a `users` table:

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    age INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    tags TEXT[]
);
```

2. Update the connection string in `main.go`:

```go
connStr := "postgres://user:password@localhost:5432/mydb?sslmode=disable"
```

3. Uncomment the dynamic schema example call in `main()`:

```go
// exampleWithDynamicSchema(ctx)
```

4. Run the example:

```bash
go run main.go
```

## Sample Output

```
=== Example 1: Pre-defined Schema ===
CEL: user.name == "John Doe"
SQL: user.name = 'John Doe'

CEL: user.age > 30 && user.is_active
SQL: user.age > 30 AND user.is_active

CEL: user.email.contains("@example.com")
SQL: POSITION('@example.com' IN user.email) > 0

CEL: "admin" in user.tags
SQL: 'admin' = ANY(user.tags)

CEL: user.created_at > timestamp("2023-01-01T00:00:00Z")
SQL: user.created_at > '2023-01-01T00:00:00Z'
```

## Key Benefits

1. **Schema Flexibility**: Adapt to changing database schemas without code changes
2. **Type Safety**: CEL provides compile-time type checking
3. **SQL Generation**: Automatic conversion to PostgreSQL-compatible SQL
4. **PostgreSQL Features**: Full support for arrays, timestamps, and other PostgreSQL types

## Use Cases

- Dynamic query builders
- API filtering systems
- Configuration-driven data access
- Multi-tenant applications with varying schemas
