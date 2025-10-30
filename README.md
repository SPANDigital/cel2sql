# cel2sql

> Convert [CEL (Common Expression Language)](https://cel.dev/) expressions to PostgreSQL SQL

[![Go Version](https://img.shields.io/badge/Go-1.24%2B-blue)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17-blue)](https://www.postgresql.org)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Benchmarks](https://img.shields.io/badge/benchmarks-performance%20tracking-green)](https://spandigital.github.io/cel2sql/dev/bench/)

**cel2sql** makes it easy to build dynamic SQL queries using CEL expressions. Write type-safe, expressive filters in CEL and automatically convert them to PostgreSQL-compatible SQL.

## Quick Start

### Installation

```bash
go get github.com/spandigital/cel2sql/v3
```

### Basic Example

```go
package main

import (
    "fmt"
    "github.com/google/cel-go/cel"
    "github.com/spandigital/cel2sql/v3"
    "github.com/spandigital/cel2sql/v3/pg"
)

func main() {
    // 1. Define your database table schema
    userSchema := pg.NewSchema([]pg.FieldSchema{
        {Name: "name", Type: "text"},
        {Name: "age", Type: "integer"},
        {Name: "active", Type: "boolean"},
    })

    // 2. Create CEL environment
    env, _ := cel.NewEnv(
        cel.CustomTypeProvider(pg.NewTypeProvider(map[string]pg.Schema{
            "User": userSchema,
        })),
        cel.Variable("user", cel.ObjectType("User")),
    )

    // 3. Write your filter expression in CEL
    ast, _ := env.Compile(`user.age >= 18 && user.active`)

    // 4. Convert to SQL
    sqlWhere, _ := cel2sql.Convert(ast)

    fmt.Println(sqlWhere)
    // Output: user.age >= 18 AND user.active IS TRUE

    // 5. Use in your query
    query := "SELECT * FROM users WHERE " + sqlWhere
}
```

## Why cel2sql?

✅ **Type-Safe**: Catch errors at compile time, not runtime
✅ **PostgreSQL 17**: Fully compatible with the latest PostgreSQL
✅ **Rich Features**: JSON/JSONB, arrays, regex, timestamps, and more
✅ **Well-Tested**: 100+ tests including integration tests with real PostgreSQL
✅ **Easy to Use**: Simple API, comprehensive documentation
✅ **Secure by Default**: Built-in protections against SQL injection and ReDoS attacks
✅ **Performance Tracked**: [Continuous benchmark monitoring](https://spandigital.github.io/cel2sql/dev/bench/) to prevent regressions

## Security Features

cel2sql includes comprehensive security protections:

- 🛡️ **Field Name Validation** - Prevents SQL injection via field names
- 🔒 **JSON Field Escaping** - Automatic quote escaping in JSON paths
- 🚫 **ReDoS Protection** - Validates regex patterns to prevent catastrophic backtracking
- 🔄 **Recursion Depth Limits** - Prevents stack overflow from deeply nested expressions (default: 100)
- 📏 **SQL Output Length Limits** - Prevents memory exhaustion from extremely large SQL queries (default: 50,000 chars)
- 🔢 **Byte Array Length Limits** - Prevents memory exhaustion from large hex-encoded byte arrays (max: 10,000 bytes)
- ⏱️ **Context Timeouts** - Optional timeout protection for complex expressions

All security features are enabled by default with zero configuration required.

## Advanced Options

cel2sql supports optional advanced features via functional options:

```go
import (
    "context"
    "log/slog"
    "github.com/spandigital/cel2sql/v3"
)

// Basic conversion
sql, err := cel2sql.Convert(ast)

// With schemas for JSON/JSONB support
sql, err := cel2sql.Convert(ast,
    cel2sql.WithSchemas(schemas))

// With context for timeouts
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
sql, err := cel2sql.Convert(ast,
    cel2sql.WithContext(ctx),
    cel2sql.WithSchemas(schemas))

// With logging for observability
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
sql, err := cel2sql.Convert(ast,
    cel2sql.WithContext(ctx),
    cel2sql.WithSchemas(schemas),
    cel2sql.WithLogger(logger))
```

**Available Options:**
- `WithSchemas(map[string]pg.Schema)` - Provide table schemas for JSON detection
- `WithContext(context.Context)` - Enable cancellation and timeouts
- `WithLogger(*slog.Logger)` - Enable structured logging
- `WithMaxDepth(int)` - Set custom recursion depth limit (default: 100)

## Query Analysis and Index Recommendations

cel2sql can analyze your CEL queries and recommend database indexes to optimize performance. The `AnalyzeQuery()` function returns both the converted SQL and actionable index recommendations.

### How It Works

`AnalyzeQuery()` examines your CEL expression and detects patterns that would benefit from specific PostgreSQL index types:

- **JSON/JSONB path operations** (`->>, ?`) → GIN indexes
- **Array operations** (comprehensions, `IN` clauses) → GIN indexes
- **Regex matching** (`matches()`) → GIN indexes with `pg_trgm` extension
- **Comparison operations** (`==, >, <, >=, <=`) → B-tree indexes

### Usage

```go
sql, recommendations, err := cel2sql.AnalyzeQuery(ast,
    cel2sql.WithSchemas(schemas))
if err != nil {
    log.Fatal(err)
}

// Use the generated SQL
rows, err := db.Query("SELECT * FROM users WHERE " + sql)

// Review and apply index recommendations
for _, rec := range recommendations {
    fmt.Printf("Column: %s\n", rec.Column)
    fmt.Printf("Type: %s\n", rec.IndexType)
    fmt.Printf("Reason: %s\n", rec.Reason)
    fmt.Printf("Execute: %s\n\n", rec.Expression)

    // Apply the recommendation
    // _, err := db.Exec(rec.Expression)
}
```

### Example

```go
// Query with multiple index-worthy patterns
celExpr := `person.age > 18 &&
            person.email.matches(r"@example\.com$") &&
            person.metadata.verified == true`

ast, _ := env.Compile(celExpr)
sql, recs, _ := cel2sql.AnalyzeQuery(ast, cel2sql.WithSchemas(schemas))

// Generated SQL:
// person.age > 18 AND person.email ~ '@example\.com$'
// AND person.metadata->>'verified' = 'true'

// Recommendations:
// 1. CREATE INDEX idx_person_age_btree ON table_name (person.age);
//    Reason: Comparison operations benefit from B-tree for range queries
//
// 2. CREATE INDEX idx_person_email_gin_trgm ON table_name
//    USING GIN (person.email gin_trgm_ops);
//    Reason: Regex matching benefits from GIN index with pg_trgm
//
// 3. CREATE INDEX idx_person_metadata_gin ON table_name
//    USING GIN (person.metadata);
//    Reason: JSON path operations benefit from GIN index
```

### When to Use

- **Development**: Discover which indexes your queries need
- **Performance tuning**: Identify missing indexes causing slow queries
- **Production monitoring**: Analyze user-generated filter expressions

See `examples/index_analysis/` for a complete working example.

## Parameterized Queries

cel2sql supports **parameterized queries** (prepared statements) for improved performance, security, and monitoring.

### Benefits

🚀 **Performance** - PostgreSQL caches query plans for parameterized queries, enabling plan reuse across executions
🔒 **Security** - Parameters are passed separately from SQL, providing defense-in-depth SQL injection protection
📊 **Monitoring** - Same query pattern appears in logs/metrics, making analysis easier

### Usage

```go
// Convert to parameterized SQL
result, err := cel2sql.ConvertParameterized(ast)
if err != nil {
    log.Fatal(err)
}

fmt.Println(result.SQL)         // "user.age > $1 AND user.name = $2"
fmt.Println(result.Parameters)  // [18 "John"]

// Execute with database/sql
rows, err := db.Query(
    "SELECT * FROM users WHERE " + result.SQL,
    result.Parameters...,
)
```

### What Gets Parameterized?

**Parameterized** (values become placeholders):
- ✅ String literals: `'John'` → `$1`
- ✅ Numeric literals: `42`, `3.14` → `$1`, `$2`
- ✅ Byte literals: `b"data"` → `$1`

**Kept Inline** (for query plan optimization):
- ❌ `TRUE`, `FALSE` - Boolean constants
- ❌ `NULL` - Null values

PostgreSQL's query planner optimizes better when it knows boolean and null values at plan time.

### Example Comparison

```go
celExpr := `user.age > 18 && user.active == true && user.name == "John"`
ast, _ := env.Compile(celExpr)

// Non-parameterized (inline values)
sql, _ := cel2sql.Convert(ast)
// SQL: user.age > 18 AND user.active IS TRUE AND user.name = 'John'

// Parameterized (placeholders + parameters)
result, _ := cel2sql.ConvertParameterized(ast)
// SQL: user.age > $1 AND user.active IS TRUE AND user.name = $2
// Parameters: [18 "John"]
// Note: TRUE is kept inline for query plan efficiency
```

### Prepared Statements

For maximum performance with repeated queries, use prepared statements:

```go
result, _ := cel2sql.ConvertParameterized(ast)

// Prepare once
stmt, err := db.Prepare("SELECT * FROM users WHERE " + result.SQL)
defer stmt.Close()

// Execute multiple times with different parameters
rows1, _ := stmt.Query(25)  // age > 25
rows2, _ := stmt.Query(30)  // age > 30
rows3, _ := stmt.Query(35)  // age > 35 (reuses cached plan!)
```

See the [parameterized example](examples/parameterized/) for a complete working demo with PostgreSQL integration.

## Common Use Cases

### 1. User Filters

```go
// CEL: Simple comparison
user.age > 21 && user.country == "USA"
// SQL: user.age > 21 AND user.country = 'USA'
```

### 2. Text Search

```go
// CEL: String operations
user.email.startsWith("admin") || user.name.contains("John")
// SQL: user.email LIKE 'admin%' OR POSITION('John' IN user.name) > 0
```

### 3. Date Filters

```go
// CEL: Date comparisons
user.created_at > timestamp("2024-01-01T00:00:00Z")
// SQL: user.created_at > CAST('2024-01-01T00:00:00Z' AS TIMESTAMP WITH TIME ZONE)
```

### 4. JSON/JSONB Fields

```go
// CEL: JSON field access
user.preferences.theme == "dark"
// SQL: user.preferences->>'theme' = 'dark'
```

### 5. Array Operations

```go
// CEL: Check if all items match
user.scores.all(s, s >= 60)
// SQL: NOT EXISTS (SELECT 1 FROM UNNEST(user.scores) AS s WHERE NOT (s >= 60))
```

## Documentation

- 📖 **[Getting Started Guide](docs/getting-started.md)** - Step-by-step tutorial
- 🔧 **[JSON/JSONB Support](docs/json-support.md)** - Working with JSON data
- 🎯 **[Array Comprehensions](docs/comprehensions.md)** - Advanced array operations
- 🔍 **[Regex Matching](docs/regex-matching.md)** - Pattern matching with regex
- 🛡️ **[Security Guide](docs/security.md)** - Security features and best practices
- 📚 **[Operators Reference](docs/operators-reference.md)** - Complete operator list
- 💡 **[Examples](examples/)** - More code examples

## Supported Features

| Feature | CEL Example | PostgreSQL SQL |
|---------|-------------|----------------|
| Comparisons | `age > 18` | `age > 18` |
| Logic | `active && verified` | `active IS TRUE AND verified IS TRUE` |
| Strings | `name.startsWith("A")` | `name LIKE 'A%'` |
| Lists | `"admin" in roles` | `'admin' IN UNNEST(roles)` |
| JSON | `data.key == "value"` | `data->>'key' = 'value'` |
| Regex | `email.matches(r".*@test\.com")` | `email ~ '.*@test\.com'` |
| Dates | `created_at.getFullYear() == 2024` | `EXTRACT(YEAR FROM created_at) = 2024` |
| Conditionals | `age > 30 ? "senior" : "junior"` | `CASE WHEN age > 30 THEN 'senior' ELSE 'junior' END` |

### Regex Matching Limitations

cel2sql automatically converts CEL's RE2 regex patterns to PostgreSQL POSIX regex. While most common patterns work, some RE2 features are **not supported** and will return errors:

**Supported:**
- ✅ Basic patterns: `.*`, `[a-z]+`, `\d{3}`
- ✅ Case-insensitive flag: `(?i)pattern` → Uses `~*` operator
- ✅ Character classes: `\d`, `\w`, `\s` (converted to POSIX)
- ✅ Non-capturing groups: `(?:...)` (converted to regular groups)

**Unsupported:**
- ❌ Lookahead assertions: `(?=...)`, `(?!...)`
- ❌ Lookbehind assertions: `(?<=...)`, `(?<!...)`
- ❌ Named capture groups: `(?P<name>...)`
- ❌ Inline flags (except `(?i)`): `(?m)`, `(?s)`, `(?-i)`, etc.

**ReDoS Protection:**
cel2sql includes automatic validation to prevent Regular Expression Denial of Service attacks:
- Pattern length limited to 500 characters
- Nested quantifiers blocked: `(a+)+` ❌
- Quantified alternation blocked: `(a|a)*` ❌
- Capture group limit: 20 maximum
- Nesting depth limit: 10 levels

See [Regex Matching documentation](docs/regex-matching.md) for complete details, safe pattern examples, and performance tips.

## Type Mapping

| CEL Type | PostgreSQL Type |
|----------|-----------------|
| `int` | `bigint` |
| `double` | `double precision` |
| `bool` | `boolean` |
| `string` | `text` |
| `bytes` | `bytea` |
| `list` | `ARRAY` |
| `timestamp` | `timestamp with time zone` |
| `duration` | `INTERVAL` |

## Dynamic Schema Loading

Load table schemas directly from your PostgreSQL database:

```go
// Connect to database and load schema
provider, _ := pg.NewTypeProviderWithConnection(ctx, "postgres://user:pass@localhost/db")
defer provider.Close()

// Load table schema dynamically
provider.LoadTableSchema(ctx, "users")

// Use with CEL
env, _ := cel.NewEnv(
    cel.CustomTypeProvider(provider),
    cel.Variable("user", cel.ObjectType("users")),
)
```

See [Getting Started Guide](docs/getting-started.md) for more details.

## Requirements

- Go 1.24 or higher
- PostgreSQL 17 (also compatible with PostgreSQL 15+)

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

Apache 2.0 - See [LICENSE](LICENSE) for details.

## Related Projects

- [CEL-Go](https://github.com/google/cel-go) - Common Expression Language implementation in Go
- [CEL Spec](https://github.com/google/cel-spec) - Common Expression Language specification

## Need Help?

- 📖 Check the [documentation](docs/)
- 💬 [Open an issue](https://github.com/SPANDigital/cel2sql/issues)
- 🌟 Star the repo if you find it useful!
