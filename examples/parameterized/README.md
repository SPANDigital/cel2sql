# Parameterized Queries Example

This example demonstrates how to use cel2sql with parameterized queries for improved performance, security, and monitoring.

## What are Parameterized Queries?

Parameterized queries (also called prepared statements) separate SQL structure from data values:

**Non-parameterized:**
```sql
SELECT * FROM users WHERE age = 30 AND name = 'John'
```

**Parameterized:**
```sql
SELECT * FROM users WHERE age = $1 AND name = $2
-- Parameters: [30, "John"]
```

## Benefits

### 1. **Performance - Query Plan Caching**
PostgreSQL caches the execution plan for parameterized queries. When you execute the same query with different parameters, PostgreSQL reuses the cached plan instead of re-planning.

```go
// Same SQL structure, different values → plan reuse
result1, _ := cel2sql.ConvertParameterized(ast)  // age > $1
db.Query(result1.SQL, 25)  // Uses cached plan
db.Query(result1.SQL, 30)  // Reuses same plan
db.Query(result1.SQL, 35)  // Reuses same plan
```

### 2. **Security - SQL Injection Protection**
Parameters are passed separately from SQL text, providing defense-in-depth protection:

```go
result, _ := cel2sql.ConvertParameterized(ast)
// SQL: "users.name = $1"
// Parameters: ["John'; DROP TABLE users--"]
// PostgreSQL treats the entire string as data, not SQL
```

### 3. **Monitoring - Better Query Pattern Analysis**
Same query structure appears consistently in logs and metrics:

```
-- Non-parameterized (3 different log entries):
SELECT * FROM users WHERE age = 25
SELECT * FROM users WHERE age = 30
SELECT * FROM users WHERE age = 35

-- Parameterized (1 query pattern):
SELECT * FROM users WHERE age = $1
```

## API Usage

### Basic Parameterized Conversion

```go
import (
    "github.com/spandigital/cel2sql/v3"
)

// Compile CEL expression
ast, _ := env.Compile(`user.age > 18 && user.name == "John"`)

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

### With Options

```go
import (
    "context"
    "log/slog"
    "github.com/spandigital/cel2sql/v3"
)

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

result, err := cel2sql.ConvertParameterized(ast,
    cel2sql.WithSchemas(schemas),
    cel2sql.WithContext(ctx),
    cel2sql.WithLogger(logger),
)
```

### Prepared Statements

For maximum performance, use prepared statements:

```go
result, _ := cel2sql.ConvertParameterized(ast)

// Prepare once
stmt, err := db.Prepare("SELECT * FROM users WHERE " + result.SQL)
if err != nil {
    log.Fatal(err)
}
defer stmt.Close()

// Execute multiple times with different parameters
rows1, _ := stmt.Query(25)   // age > 25
rows2, _ := stmt.Query(30)   // age > 30
rows3, _ := stmt.Query(35)   // age > 35
```

## What Gets Parameterized?

### Parameterized Constants
- ✅ **String literals**: `'John'` → `$1`
- ✅ **Integer literals**: `42` → `$1`
- ✅ **Float literals**: `3.14` → `$1`
- ✅ **Byte literals**: `b"data"` → `$1`

### Kept Inline (For Query Plan Optimization)
- ❌ **TRUE/FALSE**: Boolean constants stay inline
- ❌ **NULL**: Null values stay inline

PostgreSQL's query planner optimizes better when it knows boolean and null values at plan time.

### Example

```go
// CEL: user.age > 18 && user.active == true && user.name != null
// SQL: user.age > $1 AND user.active IS TRUE AND user.name IS NOT NULL
// Parameters: [18]
//
// Note: TRUE and NULL are inline, only 18 is parameterized
```

## Running This Example

```bash
# From the examples/parameterized directory:
go run main.go

# Or from the project root:
go run ./examples/parameterized
```

The example will:
1. Start a PostgreSQL 17 container using testcontainers
2. Create a test schema and insert sample data
3. Demonstrate various parameterized query patterns
4. Show query plan caching benefits
5. Compare parameterized vs non-parameterized approaches
6. Clean up the container automatically

## Expected Output

```
Starting PostgreSQL 17 container...
✓ Inserted 5 test users

================================================================================
PARAMETERIZED QUERIES DEMONSTRATION
================================================================================

--------------------------------------------------------------------------------
Example 1: Simple Parameterized Query
--------------------------------------------------------------------------------
CEL Expression: users.age > 28

Generated SQL:    users.age > $1
Parameters:       [28]
Parameter Types:  int64

Results:
--------------------------------------------------
ID: 1, Name: Alice Smith         Age: 30
ID: 3, Name: Carol Williams      Age: 35
ID: 5, Name: Eve Davis           Age: 32

...
```

## Integration with database/sql

The `Result.Parameters` field is `[]interface{}` to match Go's `database/sql` API:

```go
func (db *DB) Query(query string, args ...interface{}) (*Rows, error)
func (db *DB) Exec(query string, args ...interface{}) (Result, error)
func (stmt *Stmt) Query(args ...interface{}) (*Rows, error)
```

This allows seamless integration:

```go
result, _ := cel2sql.ConvertParameterized(ast)

// Direct query
db.Query("SELECT * FROM t WHERE " + result.SQL, result.Parameters...)

// Prepared statement
stmt, _ := db.Prepare("SELECT * FROM t WHERE " + result.SQL)
stmt.Query(result.Parameters...)
```

## Performance Tips

1. **Use prepared statements** for queries executed multiple times
2. **Keep the same SQL structure** to maximize plan cache hits
3. **Consider connection pooling** with `sql.DB` to reuse prepared statements
4. **Monitor pg_stat_statements** to verify plan caching is working

## See Also

- [Basic Example](../basic/) - Simple cel2sql usage
- [Load Table Schema Example](../load_table_schema/) - Dynamic schema loading
- [Context Example](../context/) - Context cancellation and timeouts
- [Logging Example](../logging/) - Observability integration
- [PostgreSQL Documentation on Prepared Statements](https://www.postgresql.org/docs/current/sql-prepare.html)
