# cel2sql

> Convert [CEL (Common Expression Language)](https://cel.dev/) expressions to SQL for PostgreSQL, MySQL, SQLite, DuckDB, and BigQuery

[![Go Version](https://img.shields.io/badge/Go-1.24%2B-blue)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17-336791)](https://www.postgresql.org)
[![MySQL](https://img.shields.io/badge/MySQL-8.0-4479A1)](https://www.mysql.com)
[![SQLite](https://img.shields.io/badge/SQLite-3-003B57)](https://www.sqlite.org)
[![DuckDB](https://img.shields.io/badge/DuckDB-1.x-FFF000)](https://duckdb.org)
[![BigQuery](https://img.shields.io/badge/BigQuery-GCP-4285F4)](https://cloud.google.com/bigquery)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Benchmarks](https://img.shields.io/badge/benchmarks-performance%20tracking-green)](https://spandigital.github.io/cel2sql/dev/bench/)

**cel2sql** makes it easy to build dynamic SQL queries using CEL expressions. Write type-safe, expressive filters in CEL and automatically convert them to SQL for your database of choice.

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

✅ **Multi-Dialect**: PostgreSQL, MySQL, SQLite, DuckDB, and BigQuery from a single API
✅ **Type-Safe**: Catch errors at compile time, not runtime
✅ **Rich Features**: JSON/JSONB, arrays, regex, timestamps, and more
✅ **Well-Tested**: 100+ tests including integration tests with real databases
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
- `WithDialect(dialect.Dialect)` - Select target SQL dialect (default: PostgreSQL)
- `WithSchemas(map[string]pg.Schema)` - Provide table schemas for JSON detection
- `WithContext(context.Context)` - Enable cancellation and timeouts
- `WithLogger(*slog.Logger)` - Enable structured logging
- `WithMaxDepth(int)` - Set custom recursion depth limit (default: 100)

## Multi-Dialect Support

cel2sql supports 5 SQL dialects. PostgreSQL is the default; select other dialects with `WithDialect()`:

```go
import (
    "github.com/spandigital/cel2sql/v3"
    "github.com/spandigital/cel2sql/v3/dialect/mysql"
    "github.com/spandigital/cel2sql/v3/dialect/sqlite"
    "github.com/spandigital/cel2sql/v3/dialect/duckdb"
    "github.com/spandigital/cel2sql/v3/dialect/bigquery"
)

// PostgreSQL (default - no option needed)
sql, err := cel2sql.Convert(ast)

// MySQL
sql, err := cel2sql.Convert(ast, cel2sql.WithDialect(mysql.New()))

// SQLite
sql, err := cel2sql.Convert(ast, cel2sql.WithDialect(sqlite.New()))

// DuckDB
sql, err := cel2sql.Convert(ast, cel2sql.WithDialect(duckdb.New()))

// BigQuery
sql, err := cel2sql.Convert(ast, cel2sql.WithDialect(bigquery.New()))
```

### Dialect Comparison

| Feature | PostgreSQL | MySQL | SQLite | DuckDB | BigQuery |
|---------|-----------|-------|--------|--------|----------|
| String concat | `\|\|` | `CONCAT()` | `\|\|` | `\|\|` | `\|\|` |
| Regex | `~ / ~*` | `REGEXP` | unsupported | `~ / ~*` | `REGEXP_CONTAINS()` |
| JSON access | `->>'f'` | `->>'$.f'` | `json_extract()` | `->>'f'` | `JSON_VALUE()` |
| Arrays | `ARRAY[...]` | JSON arrays | JSON arrays | `[...]` | `[...]` |
| UNNEST | `UNNEST(x)` | `JSON_TABLE(...)` | `json_each(x)` | `UNNEST(x)` | `UNNEST(x)` |
| Param placeholder | `$1, $2` | `?, ?` | `?, ?` | `$1, $2` | `@p1, @p2` |
| Timestamp cast | `TIMESTAMP WITH TIME ZONE` | `DATETIME` | `datetime()` | `TIMESTAMPTZ` | `TIMESTAMP` |
| Contains | `POSITION()` | `LOCATE()` | `INSTR()` | `CONTAINS()` | `STRPOS()` |
| Index analysis | BTREE, GIN, GIN+trgm | BTREE, FULLTEXT | BTREE | ART | CLUSTERING, SEARCH_INDEX |

### Per-Dialect Type Providers

Each dialect has its own type provider for mapping database types to CEL types. All providers support both pre-defined schemas (`NewTypeProvider`) and dynamic schema loading (`LoadTableSchema`):

```go
import "github.com/spandigital/cel2sql/v3/pg"       // PostgreSQL (pgxpool connection string)
import "github.com/spandigital/cel2sql/v3/mysql"     // MySQL (*sql.DB)
import "github.com/spandigital/cel2sql/v3/sqlite"    // SQLite (*sql.DB)
import "github.com/spandigital/cel2sql/v3/duckdb"    // DuckDB (*sql.DB)
import "github.com/spandigital/cel2sql/v3/bigquery"  // BigQuery (*bigquery.Client)
```

## Query Analysis and Index Recommendations

cel2sql can analyze your CEL queries and recommend database indexes to optimize performance. The `AnalyzeQuery()` function returns both the converted SQL and **dialect-specific** index recommendations.

### How It Works

`AnalyzeQuery()` examines your CEL expression and detects patterns that would benefit from indexing, then generates dialect-appropriate DDL:

- **Comparison operations** (`==, >, <, >=, <=`) → B-tree (PG/MySQL/SQLite), ART (DuckDB), Clustering (BigQuery)
- **JSON/JSONB path operations** (`->>, ?`) → GIN (PG), functional index (MySQL), Search Index (BigQuery), ART (DuckDB)
- **Regex matching** (`matches()`) → GIN with pg_trgm (PG), FULLTEXT (MySQL)
- **Array operations** (comprehensions, `IN` clauses) → GIN (PG), ART (DuckDB)

### Usage

```go
// PostgreSQL (default dialect)
sql, recommendations, err := cel2sql.AnalyzeQuery(ast,
    cel2sql.WithSchemas(schemas))

// Or specify a dialect
sql, recommendations, err := cel2sql.AnalyzeQuery(ast,
    cel2sql.WithSchemas(schemas),
    cel2sql.WithDialect(mysql.New()))

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
}
```

### Per-Dialect Index Types

| Pattern | PostgreSQL | MySQL | SQLite | DuckDB | BigQuery |
|---------|-----------|-------|--------|--------|----------|
| Comparison | BTREE | BTREE | BTREE | ART | CLUSTERING |
| JSON access | GIN | BTREE (functional) | _(skip)_ | ART | SEARCH_INDEX |
| Regex | GIN + pg_trgm | FULLTEXT | _(skip)_ | _(skip)_ | _(skip)_ |
| Array membership | GIN | _(skip)_ | _(skip)_ | ART | _(skip)_ |
| Comprehension | GIN | _(skip)_ | _(skip)_ | ART | _(skip)_ |

Unsupported patterns are silently skipped (no recommendation emitted).

### Example

```go
celExpr := `person.age > 18 && person.metadata.verified == true`
ast, _ := env.Compile(celExpr)

// PostgreSQL recommendations
sql, recs, _ := cel2sql.AnalyzeQuery(ast, cel2sql.WithSchemas(schemas))
// Recommendations:
// 1. CREATE INDEX idx_person_age_btree ON table_name (person.age);
// 2. CREATE INDEX idx_person_metadata_gin ON table_name USING GIN (person.metadata);

// MySQL recommendations
sql, recs, _ = cel2sql.AnalyzeQuery(ast,
    cel2sql.WithSchemas(schemas),
    cel2sql.WithDialect(mysql.New()))
// Recommendations:
// 1. CREATE INDEX idx_person_age_btree ON table_name (person.age);
// 2. CREATE INDEX idx_person_metadata_json ON table_name ((CAST(person.metadata->>'$.path' AS CHAR(255))));

// BigQuery recommendations
sql, recs, _ = cel2sql.AnalyzeQuery(ast,
    cel2sql.WithSchemas(schemas),
    cel2sql.WithDialect(bigquery.New()))
// Recommendations:
// 1. ALTER TABLE table_name SET OPTIONS (clustering_columns=['person.age']);
// 2. CREATE SEARCH INDEX idx_person_metadata ON table_name (person.metadata);
```

### When to Use

- **Development**: Discover which indexes your queries need
- **Performance tuning**: Identify missing indexes causing slow queries
- **Production monitoring**: Analyze user-generated filter expressions

See `examples/index_analysis/` for a complete working example with all 5 dialects.

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

### 6. Multi-Dimensional Arrays

cel2sql supports PostgreSQL multi-dimensional arrays (1D, 2D, 3D, 4D+) with automatic dimension detection:

```go
// Define schema with multi-dimensional arrays
schema := pg.NewSchema([]pg.FieldSchema{
    {Name: "tags", Type: "text", Repeated: true, Dimensions: 1},      // 1D: text[]
    {Name: "matrix", Type: "integer", Repeated: true, Dimensions: 2},  // 2D: integer[][]
    {Name: "cube", Type: "float", Repeated: true, Dimensions: 3},      // 3D: float[][][]
})

// CEL: size() automatically uses correct dimension
ast, _ := env.Compile("size(data.matrix) > 0")
// SQL: COALESCE(ARRAY_LENGTH(data.matrix, 2), 0) > 0

// Or load dimensions automatically from database
provider, _ := pg.NewTypeProviderWithConnection(ctx, connString)
provider.LoadTableSchema(ctx, "products")  // Dimensions detected from schema
```

**Dimension Detection:**
- Detects dimensions from PostgreSQL type strings (`integer[][]`, `_int4[]`)
- Works with both bracket notation and underscore notation
- Defaults to 1D for backward compatibility when no schema is provided

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
| Multi-Dim Arrays | `size(matrix) > 0` | `COALESCE(ARRAY_LENGTH(matrix, 2), 0) > 0` |
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

| CEL Type | PostgreSQL | MySQL | SQLite | DuckDB | BigQuery |
|----------|-----------|-------|--------|--------|----------|
| `int` | `bigint` | `SIGNED` | `INTEGER` | `BIGINT` | `INT64` |
| `double` | `double precision` | `DECIMAL` | `REAL` | `DOUBLE` | `FLOAT64` |
| `bool` | `boolean` | `UNSIGNED` | `INTEGER` | `BOOLEAN` | `BOOL` |
| `string` | `text` | `CHAR` | `TEXT` | `VARCHAR` | `STRING` |
| `bytes` | `bytea` | `BINARY` | `BLOB` | `BLOB` | `BYTES` |
| `list` | `ARRAY` | JSON array | JSON array | `LIST` | `ARRAY` |
| `timestamp` | `timestamptz` | `DATETIME` | `datetime()` | `TIMESTAMPTZ` | `TIMESTAMP` |
| `duration` | `INTERVAL` | `INTERVAL` | string modifier | `INTERVAL` | `INTERVAL` |

## Dynamic Schema Loading

Load table schemas directly from your database at runtime instead of defining them manually. Each dialect provider supports introspecting table schemas from a live database connection.

### PostgreSQL

```go
import "github.com/spandigital/cel2sql/v3/pg"

// PostgreSQL accepts a connection string and manages its own connection pool
provider, _ := pg.NewTypeProviderWithConnection(ctx, "postgres://user:pass@localhost/db")
defer provider.Close()

provider.LoadTableSchema(ctx, "users")

env, _ := cel.NewEnv(
    cel.CustomTypeProvider(provider),
    cel.Variable("user", cel.ObjectType("users")),
)
```

### MySQL

```go
import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "github.com/spandigital/cel2sql/v3/mysql"
)

// MySQL accepts a *sql.DB — you own the connection
db, _ := sql.Open("mysql", "user:pass@tcp(localhost:3306)/mydb?parseTime=true")
defer db.Close()

provider, _ := mysql.NewTypeProviderWithConnection(ctx, db)
provider.LoadTableSchema(ctx, "users")

env, _ := cel.NewEnv(
    cel.CustomTypeProvider(provider),
    cel.Variable("user", cel.ObjectType("users")),
)

sql, _ := cel2sql.Convert(ast, cel2sql.WithDialect(mysqlDialect.New()),
    cel2sql.WithSchemas(provider.GetSchemas()))
```

### SQLite

```go
import (
    "database/sql"
    _ "modernc.org/sqlite"
    "github.com/spandigital/cel2sql/v3/sqlite"
)

db, _ := sql.Open("sqlite", "mydb.sqlite")
defer db.Close()

provider, _ := sqlite.NewTypeProviderWithConnection(ctx, db)
provider.LoadTableSchema(ctx, "users")

env, _ := cel.NewEnv(
    cel.CustomTypeProvider(provider),
    cel.Variable("user", cel.ObjectType("users")),
)

sql, _ := cel2sql.Convert(ast, cel2sql.WithDialect(sqliteDialect.New()),
    cel2sql.WithSchemas(provider.GetSchemas()))
```

### DuckDB

```go
import (
    "database/sql"
    "github.com/spandigital/cel2sql/v3/duckdb"
)

// DuckDB accepts *sql.DB — works with any DuckDB driver (requires CGO)
db, _ := sql.Open("duckdb", "mydb.duckdb")
defer db.Close()

provider, _ := duckdb.NewTypeProviderWithConnection(ctx, db)
provider.LoadTableSchema(ctx, "users")

env, _ := cel.NewEnv(
    cel.CustomTypeProvider(provider),
    cel.Variable("user", cel.ObjectType("users")),
)

sql, _ := cel2sql.Convert(ast, cel2sql.WithDialect(duckdbDialect.New()),
    cel2sql.WithSchemas(provider.GetSchemas()))
```

### BigQuery

```go
import (
    "cloud.google.com/go/bigquery"
    bqprovider "github.com/spandigital/cel2sql/v3/bigquery"
)

// BigQuery uses the BigQuery client API (not database/sql)
client, _ := bigquery.NewClient(ctx, "my-project")
defer client.Close()

provider, _ := bqprovider.NewTypeProviderWithClient(ctx, client, "my_dataset")
provider.LoadTableSchema(ctx, "users")

env, _ := cel.NewEnv(
    cel.CustomTypeProvider(provider),
    cel.Variable("user", cel.ObjectType("users")),
)

sql, _ := cel2sql.Convert(ast, cel2sql.WithDialect(bigqueryDialect.New()),
    cel2sql.WithSchemas(provider.GetSchemas()))
```

### Notes

- **PostgreSQL** manages its own connection pool via `pgxpool` — call `provider.Close()` when done.
- **MySQL, SQLite, DuckDB** accept a `*sql.DB` you provide — you own the connection lifecycle. `Close()` is a no-op.
- **BigQuery** accepts a `*bigquery.Client` + dataset ID — you own the client lifecycle. `Close()` is a no-op.
- All providers also support pre-defined schemas via `NewTypeProvider(schemas)` if you don't need runtime introspection.

See [Getting Started Guide](docs/getting-started.md) for more details.

## Requirements

- Go 1.24 or higher

### CGO Requirement (DuckDB only)

The DuckDB dialect's `LoadTableSchema` requires a DuckDB Go driver (e.g., `github.com/marcboeker/go-duckdb`) which depends on **CGO** and a C/C++ compiler. This means:

- You must have `CGO_ENABLED=1` (the Go default on most platforms)
- A C/C++ compiler must be installed (GCC, Clang, or MSVC)
- Cross-compilation requires a C cross-compiler for the target platform

**All other dialects (PostgreSQL, MySQL, SQLite, BigQuery) use pure Go drivers and do not require CGO.**

If you only use DuckDB with pre-defined schemas via `duckdb.NewTypeProvider()` (no live database connection), CGO is **not** required.

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
