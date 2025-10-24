# Parameterized Queries Guide

This guide covers parameterized query support in cel2sql, including performance optimization, security benefits, and best practices.

## Table of Contents

- [Overview](#overview)
- [Why Use Parameterized Queries?](#why-use-parameterized-queries)
- [API Reference](#api-reference)
- [What Gets Parameterized?](#what-gets-parameterized)
- [Performance Optimization](#performance-optimization)
- [Security Considerations](#security-considerations)
- [Integration with database/sql](#integration-with-databasesql)
- [Examples](#examples)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)

## Overview

Parameterized queries (also known as prepared statements) separate SQL structure from data values. Instead of embedding values directly in SQL strings, they use placeholders (`$1`, `$2`, etc.) and pass values separately.

**Without parameterization:**
```sql
SELECT * FROM users WHERE age = 30 AND name = 'John'
```

**With parameterization:**
```sql
SELECT * FROM users WHERE age = $1 AND name = $2
-- Parameters: [30, "John"]
```

## Why Use Parameterized Queries?

### 1. Performance - Query Plan Caching

PostgreSQL caches execution plans for parameterized queries. When you execute the same query structure with different parameters, PostgreSQL reuses the cached plan instead of re-planning from scratch.

**Impact:**
- First execution: Parse → Plan → Execute
- Subsequent executions: Execute (reuses cached plan)
- Significant performance improvement for frequently executed queries

**Example:**
```go
// Same SQL structure, different values → PostgreSQL reuses the plan
result, _ := cel2sql.ConvertParameterized(ast)  // age > $1

db.Query(result.SQL, 25)  // Parse + Plan + Execute
db.Query(result.SQL, 30)  // Execute only (plan cached)
db.Query(result.SQL, 35)  // Execute only (plan cached)
```

### 2. Security - Defense in Depth

While cel2sql already escapes values properly, parameterized queries provide an additional layer of security by completely separating SQL structure from data.

**How it helps:**
- Parameters are sent to PostgreSQL separately from the SQL text
- PostgreSQL treats parameters as pure data, never as SQL code
- No possibility of SQL injection through parameter values
- Even if cel2sql had an escaping bug, parameters would remain safe

**Example:**
```go
// Malicious input is safely handled as data
maliciousInput := "John'; DROP TABLE users--"

result, _ := cel2sql.ConvertParameterized(ast)
// SQL: users.name = $1
// Parameters: ["John'; DROP TABLE users--"]

// PostgreSQL treats the entire string as data, not SQL
db.Query(result.SQL, maliciousInput)  // Safe!
```

### 3. Monitoring - Better Query Pattern Analysis

Parameterized queries produce consistent SQL patterns in logs and monitoring systems, making it easier to:
- Identify frequently executed queries
- Analyze query performance trends
- Set up alerting for specific query patterns
- Aggregate metrics across parameter variations

**Example:**
```
-- Non-parameterized (3 different log entries):
LOG: SELECT * FROM users WHERE age = 25
LOG: SELECT * FROM users WHERE age = 30
LOG: SELECT * FROM users WHERE age = 35

-- Parameterized (1 query pattern, easy to track):
LOG: SELECT * FROM users WHERE age = $1  (executed 3 times)
```

## API Reference

### ConvertParameterized Function

```go
func ConvertParameterized(ast *cel.Ast, opts ...ConvertOption) (*Result, error)
```

Converts a CEL AST to a parameterized PostgreSQL SQL WHERE clause with placeholders and parameter values.

**Parameters:**
- `ast` - Compiled CEL AST from `env.Compile()`
- `opts` - Optional configuration (same as `Convert()`)

**Returns:**
- `*Result` - Contains SQL string and parameter slice
- `error` - Conversion error, if any

**Options:**
All functional options from `Convert()` are supported:
- `WithSchemas(schemas)` - Required for JSON/JSONB support
- `WithContext(ctx)` - Enable cancellation and timeouts
- `WithLogger(logger)` - Enable structured logging
- `WithMaxDepth(depth)` - Set recursion depth limit

### Result Type

```go
type Result struct {
    SQL        string        // Generated SQL with placeholders ($1, $2, etc.)
    Parameters []interface{} // Parameter values in order
}
```

**Fields:**
- `SQL` - PostgreSQL SQL WHERE clause condition with `$n` placeholders
- `Parameters` - Slice of parameter values matching placeholder order

**Parameter Types:**
- `int64` - For integer literals
- `uint64` - For unsigned integer literals
- `float64` - For double/float literals
- `string` - For string literals
- `[]byte` - For byte literals

## What Gets Parameterized?

### Parameterized Constants

The following CEL constants are converted to placeholders with corresponding parameters:

#### String Literals
```go
// CEL: user.name == "John"
// SQL: user.name = $1
// Parameters: ["John"]
```

#### Integer Literals (int64)
```go
// CEL: user.age > 18
// SQL: user.age > $1
// Parameters: [18]
```

#### Unsigned Integer Literals (uint64)
```go
// CEL: user.id == 42u
// SQL: user.id = $1
// Parameters: [uint64(42)]
```

#### Float Literals (float64)
```go
// CEL: product.price <= 99.99
// SQL: product.price <= $1
// Parameters: [99.99]
```

#### Byte Literals ([]byte)
```go
// CEL: file.content == b"data"
// SQL: file.content = $1
// Parameters: [[]byte("data")]
```

### Constants Kept Inline

For PostgreSQL query plan optimization, these constants remain inline (not parameterized):

#### TRUE and FALSE
```go
// CEL: user.active == true && user.deleted == false
// SQL: user.active IS TRUE AND user.deleted IS FALSE
// Parameters: [] (empty - booleans kept inline)
```

**Reason:** PostgreSQL's query planner can make better optimization decisions when it knows boolean values at plan time. For example:
- `active IS TRUE` might use a partial index on `active = true`
- `active IS FALSE` might skip that index entirely

#### NULL
```go
// CEL: user.deleted_at == null
// SQL: user.deleted_at IS NULL
// Parameters: [] (empty - NULL kept inline)
```

**Reason:** NULL handling in PostgreSQL is special-cased with different execution paths. The planner optimizes better when it knows NULL is involved.

### Multiple Parameters

Parameters are numbered sequentially in the order they appear:

```go
// CEL: user.age > 18 && user.salary < 100000.0 && user.name == "John"
// SQL: user.age > $1 AND user.salary < $2 AND user.name = $3
// Parameters: [18, 100000.0, "John"]
```

## Performance Optimization

### Query Plan Caching

PostgreSQL caches query plans based on the SQL text. Parameterized queries maximize cache hits:

```go
result, _ := cel2sql.ConvertParameterized(ast)
// SQL: user.age > $1

// First execution: full planning
db.Query(result.SQL, 25)  // Parse → Plan → Cache → Execute

// Subsequent executions: plan reused
db.Query(result.SQL, 30)  // Execute (plan from cache)
db.Query(result.SQL, 35)  // Execute (plan from cache)
```

**Monitoring plan cache:**
```sql
-- View cached query plans
SELECT * FROM pg_prepared_statements;

-- View query statistics (requires pg_stat_statements extension)
SELECT query, calls, total_exec_time, mean_exec_time
FROM pg_stat_statements
WHERE query LIKE '%$1%'
ORDER BY calls DESC;
```

### Prepared Statements

For maximum performance with frequently executed queries, use prepared statements:

```go
result, _ := cel2sql.ConvertParameterized(ast)

// Prepare once (parsing + planning happens here)
stmt, err := db.Prepare("SELECT * FROM users WHERE " + result.SQL)
if err != nil {
    log.Fatal(err)
}
defer stmt.Close()

// Execute many times (only execution, no parsing or planning)
for _, age := range []int{18, 21, 25, 30, 35, 40} {
    rows, err := stmt.Query(age)
    // ... process rows
    rows.Close()
}
```

**Performance benefit:**
- **Without prepared statement:** Each execution parses and plans
- **With prepared statement:** Parse and plan once, execute many times
- **Speedup:** 2-10x for simple queries, more for complex queries

### Connection Pooling

Prepared statements are connection-scoped. Use connection pooling for best results:

```go
// Configure connection pool
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(25)
db.SetConnMaxLifetime(5 * time.Minute)

// Prepare statement (attached to a connection from pool)
stmt, err := db.Prepare("SELECT * FROM users WHERE " + result.SQL)
defer stmt.Close()

// Executions may use different connections
// But each connection caches the plan independently
stmt.Query(25)
stmt.Query(30)
```

### Benchmarking

Compare performance with and without parameterization:

```go
func BenchmarkParameterized(b *testing.B) {
    result, _ := cel2sql.ConvertParameterized(ast)
    stmt, _ := db.Prepare("SELECT * FROM users WHERE " + result.SQL)
    defer stmt.Close()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        rows, _ := stmt.Query(i % 100)
        rows.Close()
    }
}

func BenchmarkInline(b *testing.B) {
    for i := 0; i < b.N; i++ {
        sql, _ := cel2sql.Convert(ast)
        rows, _ := db.Query("SELECT * FROM users WHERE " + sql)
        rows.Close()
    }
}
```

Expected results:
- **Inline:** ~1000 ns/op (must parse and plan each time)
- **Parameterized with prepared statement:** ~100 ns/op (10x faster)

## Security Considerations

### Defense in Depth

Parameterized queries provide multiple layers of protection:

1. **CEL Type Safety** - Compile-time type checking prevents invalid expressions
2. **Field Name Validation** - `validateFieldName()` prevents SQL injection via field names
3. **JSON Field Escaping** - `escapeJSONFieldName()` prevents injection via JSON paths
4. **Value Escaping** - Non-parameterized mode properly escapes all values
5. **Parameterization** - Values sent separately from SQL (this layer)

Even if layers 4 fails, layer 5 prevents SQL injection.

### SQL Injection Prevention

Parameterized queries eliminate SQL injection for constant values:

```go
// Attack attempt
maliciousName := "John'; DELETE FROM users WHERE '1'='1"

// Non-parameterized (still safe due to escaping)
sql, _ := cel2sql.Convert(ast)
// SQL: user.name = 'John''; DELETE FROM users WHERE ''1''=''1'
// Safe due to quote escaping

// Parameterized (safe by design)
result, _ := cel2sql.ConvertParameterized(ast)
// SQL: user.name = $1
// Parameters: ["John'; DELETE FROM users WHERE '1'='1"]
// PostgreSQL treats entire string as data, not SQL
```

### Field Name Security

Field names are still validated even with parameterization:

```go
// Invalid field name (contains SQL)
celExpr := `user.name; DROP TABLE users-- == "test"`

// Both functions reject this during conversion
_, err1 := cel2sql.Convert(ast)              // Error: invalid field name
_, err2 := cel2sql.ConvertParameterized(ast) // Error: invalid field name
```

**Field name validation rules:**
- Maximum 63 characters (PostgreSQL NAMEDATALEN-1)
- Must start with letter or underscore
- Only alphanumeric characters and underscores
- Cannot be SQL reserved keywords

### Type Safety

Parameters maintain type information:

```go
// Type mismatch detected at execution time
result, _ := cel2sql.ConvertParameterized(ast)
// SQL: user.age > $1
// Parameters: ["not a number"]  // Wrong type!

// PostgreSQL returns error:
// ERROR: invalid input syntax for type integer: "not a number"
```

**Best practice:** Use CEL's type system to catch type errors before conversion:

```go
env, _ := cel.NewEnv(
    cel.Variable("user", cel.ObjectType("User")),
)

// This won't compile in CEL (caught early)
ast, issues := env.Compile(`user.age > "not a number"`)
if issues != nil {
    // Type error: cannot compare int to string
}
```

## Integration with database/sql

### Direct Query Execution

```go
result, err := cel2sql.ConvertParameterized(ast)
if err != nil {
    return err
}

query := "SELECT * FROM users WHERE " + result.SQL
rows, err := db.Query(query, result.Parameters...)
if err != nil {
    return err
}
defer rows.Close()

for rows.Next() {
    // Process rows
}
```

### Prepared Statement Execution

```go
result, err := cel2sql.ConvertParameterized(ast)
if err != nil {
    return err
}

stmt, err := db.Prepare("SELECT id, name FROM users WHERE " + result.SQL)
if err != nil {
    return err
}
defer stmt.Close()

rows, err := stmt.Query(result.Parameters...)
if err != nil {
    return err
}
defer rows.Close()
```

### Transaction Support

```go
tx, err := db.Begin()
if err != nil {
    return err
}
defer tx.Rollback()

result, _ := cel2sql.ConvertParameterized(ast)
rows, err := tx.Query("SELECT * FROM users WHERE "+result.SQL, result.Parameters...)
if err != nil {
    return err
}
defer rows.Close()

// Process rows...

if err := tx.Commit(); err != nil {
    return err
}
```

### Context Support

Combine parameterized queries with context for timeouts:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Conversion with context
result, err := cel2sql.ConvertParameterized(ast, cel2sql.WithContext(ctx))
if err != nil {
    return err
}

// Query execution with context
rows, err := db.QueryContext(ctx, "SELECT * FROM users WHERE "+result.SQL, result.Parameters...)
```

## Examples

### Example 1: Simple Filter

```go
// Schema
schema := pg.Schema{
    {Name: "age", Type: "int"},
    {Name: "name", Type: "text"},
}
provider := pg.NewTypeProvider(map[string]pg.Schema{"User": schema})

// CEL environment
env, _ := cel.NewEnv(
    cel.CustomTypeProvider(provider),
    cel.Variable("user", cel.ObjectType("User")),
)

// CEL expression
ast, _ := env.Compile(`user.age > 18 && user.name == "John"`)

// Convert
result, _ := cel2sql.ConvertParameterized(ast)
fmt.Println(result.SQL)         // user.age > $1 AND user.name = $2
fmt.Println(result.Parameters)  // [18 "John"]

// Execute
rows, _ := db.Query("SELECT * FROM users WHERE "+result.SQL, result.Parameters...)
```

### Example 2: Complex Expression

```go
celExpr := `user.age >= 21 && user.salary > 50000.0 && user.department == "Engineering" && user.active == true`
ast, _ := env.Compile(celExpr)

result, _ := cel2sql.ConvertParameterized(ast)
// SQL: user.age >= $1 AND user.salary > $2 AND user.department = $3 AND user.active IS TRUE
// Parameters: [21, 50000.0, "Engineering"]
// Note: TRUE is kept inline
```

### Example 3: JSON/JSONB Fields

```go
// Schema with JSONB field
schema := pg.Schema{
    {Name: "metadata", Type: "jsonb", IsJSON: true, IsJSONB: true},
}

celExpr := `user.metadata.age > 25`
ast, _ := env.Compile(celExpr)

result, _ := cel2sql.ConvertParameterized(ast,
    cel2sql.WithSchemas(map[string]pg.Schema{"User": schema}))
// SQL: (user.metadata->>'age')::numeric > $1
// Parameters: [25]
```

### Example 4: Array Operations

```go
celExpr := `"admin" in user.roles`
ast, _ := env.Compile(celExpr)

result, _ := cel2sql.ConvertParameterized(ast)
// SQL: $1 = ANY(user.roles)
// Parameters: ["admin"]
```

### Example 5: Comprehensions

```go
celExpr := `user.orders.all(o, o.total > 100.0)`
ast, _ := env.Compile(celExpr)

result, _ := cel2sql.ConvertParameterized(ast)
// SQL: NOT EXISTS (SELECT 1 FROM UNNEST(user.orders) AS o WHERE NOT (o.total > $1))
// Parameters: [100.0]
```

### Example 6: Logging Integration

```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,
}))

result, _ := cel2sql.ConvertParameterized(ast,
    cel2sql.WithSchemas(schemas),
    cel2sql.WithLogger(logger))

// Logs include:
// - Conversion start/completion
// - Parameter count
// - Duration
// - Any errors
```

## Best Practices

### 1. Use Parameterization by Default

For most use cases, prefer `ConvertParameterized()` over `Convert()`:

```go
// ✅ Good (default choice)
result, _ := cel2sql.ConvertParameterized(ast)
db.Query("SELECT * FROM users WHERE "+result.SQL, result.Parameters...)

// ⚠️ Use only when parameterization isn't needed
sql, _ := cel2sql.Convert(ast)
db.Query("SELECT * FROM users WHERE " + sql)
```

### 2. Reuse Prepared Statements

Prepare statements once and reuse them:

```go
// ✅ Good
stmt, _ := db.Prepare(query)
defer stmt.Close()
for _, param := range params {
    stmt.Query(param)
}

// ❌ Bad (prepares each time)
for _, param := range params {
    stmt, _ := db.Prepare(query)
    stmt.Query(param)
    stmt.Close()
}
```

### 3. Use Connection Pooling

Configure appropriate pool settings:

```go
db.SetMaxOpenConns(25)        // Max connections
db.SetMaxIdleConns(25)        // Keep connections ready
db.SetConnMaxLifetime(5*time.Minute)  // Recycle old connections
```

### 4. Monitor Query Performance

Track query execution and plan cache hits:

```sql
-- Enable pg_stat_statements
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

-- View most executed queries
SELECT
    query,
    calls,
    mean_exec_time,
    stddev_exec_time
FROM pg_stat_statements
ORDER BY calls DESC
LIMIT 20;
```

### 5. Provide Schemas for JSON Support

Always pass schemas when using JSON/JSONB fields:

```go
// ✅ Good
result, _ := cel2sql.ConvertParameterized(ast,
    cel2sql.WithSchemas(schemas))

// ❌ Bad (JSON fields won't work correctly)
result, _ := cel2sql.ConvertParameterized(ast)
```

### 6. Handle Errors Properly

Check conversion and execution errors:

```go
result, err := cel2sql.ConvertParameterized(ast)
if err != nil {
    log.Printf("Conversion failed: %v", err)
    return err
}

rows, err := db.Query(query, result.Parameters...)
if err != nil {
    log.Printf("Query failed: %v", err)
    return err
}
defer rows.Close()
```

### 7. Use Context for Timeouts

Protect against long-running conversions and queries:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

result, err := cel2sql.ConvertParameterized(ast, cel2sql.WithContext(ctx))
if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        log.Println("Conversion timeout")
    }
    return err
}

rows, err := db.QueryContext(ctx, query, result.Parameters...)
```

## Troubleshooting

### Issue: "sql: converting argument type: unsupported type"

**Cause:** Parameter type doesn't match PostgreSQL column type.

**Solution:** Use appropriate CEL types:
```go
// ❌ Bad
ast, _ := env.Compile(`user.age > "25"`)  // String instead of int

// ✅ Good
ast, _ := env.Compile(`user.age > 25`)    // Integer
```

### Issue: Query plan not being cached

**Possible causes:**

1. **SQL text changes:** Make sure SQL structure stays the same
   ```go
   // ❌ Bad (different SQL each time)
   for age := range ages {
       ast, _ := env.Compile(fmt.Sprintf(`user.age > %d`, age))
       result, _ := cel2sql.ConvertParameterized(ast)
       // SQL changes: user.age > $1 with different AST IDs
   }

   // ✅ Good (same SQL)
   ast, _ := env.Compile(`user.age > 18`)
   result, _ := cel2sql.ConvertParameterized(ast)
   for age := range ages {
       db.Query(result.SQL, age)  // Same SQL, different params
   }
   ```

2. **Not using prepared statements:** Direct queries don't cache as effectively
   ```go
   // ⚠️ Plan caching depends on PostgreSQL's automatic caching
   db.Query(result.SQL, param1)
   db.Query(result.SQL, param2)

   // ✅ Better: explicit prepared statement
   stmt, _ := db.Prepare(result.SQL)
   stmt.Query(param1)
   stmt.Query(param2)
   ```

### Issue: Parameter ordering mismatch

**Cause:** Parameters don't match placeholder order.

**Solution:** Always use `result.Parameters` as-is:
```go
result, _ := cel2sql.ConvertParameterized(ast)

// ✅ Good
db.Query(query, result.Parameters...)

// ❌ Bad (manual parameter ordering)
db.Query(query, param2, param1)  // Wrong order!
```

### Issue: NULL values not working

**Cause:** NULL is kept inline, not parameterized.

**Expected behavior:**
```go
ast, _ := env.Compile(`user.deleted_at == null`)
result, _ := cel2sql.ConvertParameterized(ast)
// SQL: user.deleted_at IS NULL
// Parameters: []  (empty - NULL is inline)

// ✅ Correct usage
db.Query("SELECT * FROM users WHERE " + result.SQL)  // No parameters needed
```

### Issue: Performance not improving

**Diagnosis:**
1. Verify query plan caching:
   ```sql
   SELECT * FROM pg_prepared_statements;
   ```

2. Check query statistics:
   ```sql
   SELECT query, calls FROM pg_stat_statements
   WHERE query LIKE '%$1%'
   ORDER BY calls DESC;
   ```

3. Compare execution times:
   ```sql
   EXPLAIN ANALYZE SELECT * FROM users WHERE age > $1;
   ```

**Common fixes:**
- Use prepared statements instead of direct queries
- Increase connection pool size
- Ensure PostgreSQL statistics are up to date (`ANALYZE` tables)

### Issue: "invalid field name" errors

**Cause:** Field name contains invalid characters or SQL keywords.

**Solution:**
- Use valid identifiers (alphanumeric + underscore)
- Avoid SQL reserved keywords
- Keep field names under 64 characters

```go
// ❌ Bad
`user.select == "value"`  // 'select' is SQL keyword

// ✅ Good
`user.status == "value"`
```

## See Also

- [Getting Started Guide](getting-started.md)
- [Security Guide](security.md)
- [Operators Reference](operators-reference.md)
- [PostgreSQL Prepared Statements Documentation](https://www.postgresql.org/docs/current/sql-prepare.html)
- [PostgreSQL pg_stat_statements](https://www.postgresql.org/docs/current/pgstatstatements.html)
