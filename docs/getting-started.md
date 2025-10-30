# Getting Started with cel2sql

This guide will help you get up and running with cel2sql in just a few minutes.

## What is cel2sql?

cel2sql converts CEL (Common Expression Language) expressions into PostgreSQL SQL. This is useful when you want to:

- Build dynamic filters for database queries
- Allow users to create custom filters safely
- Write type-safe query conditions
- Avoid SQL injection vulnerabilities

## Installation

```bash
go get github.com/spandigital/cel2sql/v3
```

## Your First cel2sql Program

Let's build a simple program that converts CEL expressions to SQL:

```go
package main

import (
    "fmt"
    "log"

    "github.com/google/cel-go/cel"
    "github.com/spandigital/cel2sql/v3"
    "github.com/spandigital/cel2sql/v3/pg"
)

func main() {
    // Step 1: Define your table schema
    // This matches your PostgreSQL table structure
    productSchema := pg.Schema{
        {Name: "id", Type: "bigint"},
        {Name: "name", Type: "text"},
        {Name: "price", Type: "double precision"},
        {Name: "in_stock", Type: "boolean"},
        {Name: "category", Type: "text"},
    }

    // Step 2: Create a type provider
    // This tells CEL about your database types
    provider := pg.NewTypeProvider(map[string]pg.Schema{
        "Product": productSchema,
    })

    // Step 3: Create CEL environment
    env, err := cel.NewEnv(
        cel.CustomTypeProvider(provider),
        cel.Variable("product", cel.ObjectType("Product")),
    )
    if err != nil {
        log.Fatalf("Failed to create CEL environment: %v", err)
    }

    // Step 4: Compile CEL expression
    celExpression := `product.price < 100 && product.in_stock`
    ast, issues := env.Compile(celExpression)
    if issues != nil && issues.Err() != nil {
        log.Fatalf("CEL compilation failed: %v", issues.Err())
    }

    // Step 5: Get schemas for JSON/JSONB detection
    schemas := provider.GetSchemas()

    // Step 6: Convert to SQL
    sqlWhere, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
    if err != nil {
        log.Fatalf("Failed to convert to SQL: %v", err)
    }

    // Step 7: Use in your SQL query
    query := "SELECT * FROM products WHERE " + sqlWhere
    fmt.Println("Generated SQL:")
    fmt.Println(query)
    // Output: SELECT * FROM products WHERE product.price < 100 AND product.in_stock IS TRUE
}
```

## Understanding the Steps

### 1. Define Your Schema

The schema tells cel2sql about your database table structure:

```go
schema := pg.Schema{
    {Name: "column_name", Type: "postgresql_type"},
    {Name: "column_name", Type: "postgresql_type", Repeated: true}, // for arrays
}
```

Common PostgreSQL types:
- `"text"` - Text/string columns
- `"bigint"` - Integer columns
- `"double precision"` - Decimal/float columns
- `"boolean"` - Boolean columns
- `"timestamp with time zone"` - Timestamp columns
- `"jsonb"` - JSON columns
- Set `Repeated: true` for array columns

### 2. Create Type Provider

The type provider connects your schema to CEL:

```go
provider := pg.NewTypeProvider(map[string]pg.Schema{
    "TableName": schema,
})
```

You can define multiple tables:

```go
provider := pg.NewTypeProvider(map[string]pg.Schema{
    "User": userSchema,
    "Product": productSchema,
    "Order": orderSchema,
})
```

### 3. Create CEL Environment

The CEL environment is where you define variables and compile expressions:

```go
env, err := cel.NewEnv(
    cel.CustomTypeProvider(provider),
    cel.Variable("product", cel.ObjectType("Product")),
)
```

### 4. Compile Expression

Compile your CEL expression to check for syntax errors:

```go
ast, issues := env.Compile(`product.price < 100`)
if issues != nil && issues.Err() != nil {
    // Handle compilation error
}
```

### 5. Convert to SQL

Convert the compiled CEL expression to SQL using functional options:

```go
// Basic conversion
sqlWhere, err := cel2sql.Convert(ast)

// With schemas for JSON/JSONB support (recommended)
schemas := provider.GetSchemas()
sqlWhere, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))

// With multiple options
sqlWhere, err := cel2sql.Convert(ast,
    cel2sql.WithSchemas(schemas),
    cel2sql.WithContext(ctx),
    cel2sql.WithLogger(logger))
```

## Dynamic Schema Loading

Instead of manually defining schemas, you can load them directly from your PostgreSQL database:

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/google/cel-go/cel"
    "github.com/spandigital/cel2sql/v3"
    "github.com/spandigital/cel2sql/v3/pg"
)

func main() {
    ctx := context.Background()

    // Connect to PostgreSQL and create provider
    connString := "postgres://user:password@localhost:5432/mydb?sslmode=disable"
    provider, err := pg.NewTypeProviderWithConnection(ctx, connString)
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer provider.Close()

    // Load table schema from database
    err = provider.LoadTableSchema(ctx, "products")
    if err != nil {
        log.Fatalf("Failed to load schema: %v", err)
    }

    // Create CEL environment with loaded schema
    env, err := cel.NewEnv(
        cel.CustomTypeProvider(provider),
        cel.Variable("product", cel.ObjectType("products")),
    )
    if err != nil {
        log.Fatalf("Failed to create environment: %v", err)
    }

    // Now use it like before
    ast, issues := env.Compile(`product.price < 100`)
    if issues != nil && issues.Err() != nil {
        log.Fatalf("Compilation error: %v", issues.Err())
    }

    // Get schemas for JSON/JSONB detection
    schemas := provider.GetSchemas()

    sqlWhere, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
    if err != nil {
        log.Fatalf("Conversion error: %v", err)
    }

    fmt.Println(sqlWhere)
    // Output: product.price < 100
}
```

**Benefits of dynamic schema loading:**
- No manual schema definition needed
- Automatically stays in sync with database changes
- Supports composite types and complex structures
- Easier to maintain

## Common Examples

### Simple Comparisons

```go
// CEL
product.price > 50
// SQL: product.price > 50

// CEL
product.name == "Laptop"
// SQL: product.name = 'Laptop'
```

### Logical Operators

```go
// CEL
product.price < 100 && product.in_stock
// SQL: product.price < 100 AND product.in_stock IS TRUE

// CEL
product.category == "Electronics" || product.category == "Computers"
// SQL: product.category = 'Electronics' OR product.category = 'Computers'
```

### String Operations

```go
// CEL
product.name.startsWith("Mac")
// SQL: product.name LIKE 'Mac%'

// CEL
product.description.contains("wireless")
// SQL: POSITION('wireless' IN product.description) > 0
```

### Date Comparisons

```go
// CEL
product.created_at > timestamp("2024-01-01T00:00:00Z")
// SQL: product.created_at > CAST('2024-01-01T00:00:00Z' AS TIMESTAMP WITH TIME ZONE)
```

### Lists/Arrays

```go
// CEL
"electronics" in product.categories
// SQL: 'electronics' IN UNNEST(product.categories)

// CEL
size(product.tags) > 3
// SQL: ARRAY_LENGTH(product.tags, 1) > 3
```

## Error Handling

Always check for errors at each step:

```go
// Check CEL compilation
ast, issues := env.Compile(expression)
if issues != nil && issues.Err() != nil {
    return fmt.Errorf("CEL compilation error: %w", issues.Err())
}

// Check SQL conversion
sqlWhere, err := cel2sql.Convert(ast)
if err != nil {
    return fmt.Errorf("SQL conversion error: %w", err)
}
```

Common errors:
- **Unknown field**: Field not in schema definition
- **Type mismatch**: Wrong type for operation (e.g., comparing string to number)
- **Syntax error**: Invalid CEL expression syntax

## Advanced Options

cel2sql supports optional advanced features via functional options:

### Context Support

Add timeout and cancellation support:

```go
import "context"
import "time"

// With timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

sqlWhere, err := cel2sql.Convert(ast,
    cel2sql.WithContext(ctx),
    cel2sql.WithSchemas(schemas))

if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        log.Println("Conversion timed out")
    }
}
```

**Benefits:**
- Protect against complex expressions
- Enable cancellation of long-running conversions
- Integration with distributed tracing

### Structured Logging

Enable observability with structured logging:

```go
import "log/slog"
import "os"

// Create logger (JSON for production)
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,
}))

// Convert with logging
sqlWhere, err := cel2sql.Convert(ast,
    cel2sql.WithSchemas(schemas),
    cel2sql.WithLogger(logger))
```

**What gets logged:**
- JSON path detection decisions
- Comprehension type identification
- Schema lookups
- Performance metrics
- Error contexts

### Combining Options

You can combine all options together:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
schemas := provider.GetSchemas()

sqlWhere, err := cel2sql.Convert(ast,
    cel2sql.WithContext(ctx),
    cel2sql.WithSchemas(schemas),
    cel2sql.WithLogger(logger))
```

See [examples/context](../examples/context/) and [examples/logging](../examples/logging/) for complete working examples.

## Security Best Practices

cel2sql includes built-in security protections that are enabled by default:

### 1. Field Name Validation

All field names are automatically validated to prevent SQL injection:

```go
// ✅ Safe - valid field names pass through
product.name == "Laptop"

// ❌ Blocked - prevents SQL injection
field'; DROP TABLE users--
```

**Protections:**
- Maximum field name length (63 chars)
- Alphanumeric + underscore only
- Blocks SQL reserved keywords
- Prevents common injection patterns

### 2. JSON Field Escaping

Single quotes in JSON field names are automatically escaped:

```go
// CEL with quote in field name
user.preferences.user'name == "test"

// Generated SQL (safely escaped)
user.preferences->>'user''name' = 'test'
```

### 3. ReDoS Protection

Regex patterns are validated to prevent catastrophic backtracking:

```go
// ✅ Safe patterns allowed
email.matches(r"[a-z]+@[a-z]+\.[a-z]+")

// ❌ Dangerous patterns blocked
field.matches(r"(a+)+b")  // Nested quantifiers
```

**Pattern limits:**
- Maximum 500 characters
- No nested quantifiers
- Maximum 20 capture groups
- Maximum 10 nesting levels

### 4. Recursion Depth Limits

Expressions are automatically protected from excessive nesting:

```go
// ✅ Normal expressions work fine
(product.price > 100) && (product.stock > 0)

// ❌ Excessive nesting blocked (default limit: 100 depth)
((((((((((x + 1) + 1) + 1)...)))))))) // 150+ levels deep
```

**Automatic protection:**
- Default maximum depth: 100 (sufficient for realistic expressions)
- Prevents stack overflow from deeply nested expressions
- Configurable limit with `WithMaxDepth()` option
- Protects against CWE-674 (Uncontrolled Recursion)

**Custom limits:**
```go
// Allow deeper nesting for complex expressions
sql, err := cel2sql.Convert(ast, cel2sql.WithMaxDepth(200))

// Stricter limit for untrusted input
sql, err := cel2sql.Convert(ast, cel2sql.WithMaxDepth(50))
```

### 5. Use Context Timeouts

Add defense-in-depth with context timeouts:

```go
// Protect against complex user expressions
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

sqlWhere, err := cel2sql.Convert(ast,
    cel2sql.WithContext(ctx),
    cel2sql.WithSchemas(schemas))
```

### 6. Additional Best Practices

- **Validate user input** before passing to CEL
- **Use prepared statements** when executing generated SQL
- **Keep schemas minimal** - only expose necessary fields
- **Enable logging** in production to monitor patterns
- **Test edge cases** with your specific field names

For detailed security information, see the [Security Guide](security.md).

## Next Steps

- Learn about [JSON/JSONB Support](json-support.md)
- Explore [Array Comprehensions](comprehensions.md)
- Try [Regex Matching](regex-matching.md)
- See more [Examples](../examples/)

## Best Practices

1. **Always validate user input** before compiling CEL expressions
2. **Use type providers** to catch type errors early
3. **Test your expressions** with sample data
4. **Use dynamic schema loading** for production applications
5. **Handle errors gracefully** with proper error messages
