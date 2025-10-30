# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

cel2sql converts CEL (Common Expression Language) expressions to PostgreSQL SQL conditions. It specifically targets PostgreSQL standard SQL and was recently migrated from BigQuery.

**Module**: `github.com/spandigital/cel2sql/v3`
**Go Version**: 1.24+
**Current Version**: v3.0.0

## Common Development Commands

### Build and Test
```bash
make build              # Build the project
make test               # Run tests with race detection and coverage
make test-coverage      # Generate HTML coverage report
make bench              # Run performance benchmarks
make bench-compare      # Run benchmarks and save for comparison
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

**Text and String Types:**
- `text`, `varchar`, `char`, `character varying`, `character` → `decls.String`
- `xml` → `decls.String`
- `inet`, `cidr` → `decls.String` (network addresses)
- `macaddr`, `macaddr8` → `decls.String` (MAC addresses)
- `tsvector`, `tsquery` → `decls.String` (full-text search)

**Numeric Types:**
- `bigint`, `integer`, `int`, `int4`, `int8`, `smallint`, `int2` → `decls.Int`
- `double precision`, `real`, `float4`, `float8`, `numeric`, `decimal` → `decls.Double`
- `money` → `decls.Double`

**Boolean and Binary:**
- `boolean`, `bool` → `decls.Bool`
- `bytea` → `decls.Bytes`
- `uuid` → `decls.Bytes`

**Temporal Types:**
- `timestamp`, `timestamptz`, `timestamp with time zone`, `timestamp without time zone` → `decls.Timestamp`
- `date` → `sqltypes.Date`
- `time`, `timetz`, `time with time zone`, `time without time zone` → `sqltypes.Time`

**Structured Types:**
- `json`, `jsonb` → `decls.Dyn` (with automatic JSON path support)
- Arrays: Set `Repeated: true` in schema
- Composite types: Use nested `Schema` fields

**Unsupported Types:**
Unknown PostgreSQL types (e.g., `point`, `polygon`, `box`, custom enums) will cause `FindStructFieldType()` to return `found=false`. This prevents silent type mismatches. Add explicit support for custom types or use composite type definitions.

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

### Regex Pattern Matching (v2.8.0+)

Supports CEL `matches()` function with automatic RE2 to POSIX regex conversion:
- `field.matches(r"pattern")` → `field ~ 'pattern'`
- `field.matches(r"(?i)pattern")` → `field ~* 'pattern'` (case-insensitive, `(?i)` stripped from pattern)
- `field.matches(r"(?:abc)")` → `field ~ '(abc)'` (non-capturing groups converted to regular groups)

**Automatic Conversions:**
- Case-insensitive flag `(?i)` → Uses `~*` operator, flag stripped from pattern
- Non-capturing groups `(?:...)` → Converted to regular groups `(...)`
- Character classes: `\d` → `[[:digit:]]`, `\w` → `[[:alnum:]_]`, `\s` → `[[:space:]]`
- Word boundaries: `\b` → `\y`

**Unsupported RE2 Features (will return errors):**
- Lookahead assertions: `(?=...)`, `(?!...)`
- Lookbehind assertions: `(?<=...)`, `(?<!...)`
- Named capture groups: `(?P<name>...)`
- Inline flags other than `(?i)`: `(?m)`, `(?s)`, `(?-i)`, etc.

These validations prevent PostgreSQL syntax errors and ensure predictable behavior.

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

### Performance Benchmarks

cel2sql includes comprehensive benchmarks to track performance and detect regressions. Benchmarks are automatically run as part of the CI/CD pipeline.

#### Running Benchmarks Locally

```bash
# Run all benchmarks
make bench

# Run benchmarks and save output for comparison
make bench-compare

# Run specific benchmark
go test -bench=BenchmarkConvertSimple -benchmem ./...

# Run benchmarks with custom duration
go test -bench=. -benchmem -benchtime=5s ./...
```

#### Benchmark Categories

The benchmark suite covers all major conversion features:

**Simple Operations** (`BenchmarkConvertSimple`)
- Field comparisons (equality, greater than, less than)
- String operations (equality checks)
- Boolean checks

**Operators** (`BenchmarkConvertOperators`)
- Logical operators (AND, OR)
- Arithmetic operators (+, -, *, /, %)
- String concatenation
- Complex nested expressions

**Comprehensions** (`BenchmarkConvertComprehensions`)
- `all()` - Universal quantification
- `exists()` - Existential quantification
- `exists_one()` - Unique existence
- `filter()` - Array filtering
- `map()` - Array transformation

**JSON/JSONB** (`BenchmarkConvertJSONPath`)
- Simple and nested JSON field access
- JSON field existence checks (`has()`)
- JSON field comparisons
- Complex JSON expressions

**Regex** (`BenchmarkConvertRegex`)
- Simple patterns
- Case-insensitive patterns
- Character classes (\\d, \\w, \\s)
- Word boundaries (\\b)

**Complex Expressions** (`BenchmarkConvertDeeplyNested`, `BenchmarkConvertLargeExpression`)
- Deeply nested AND/OR chains
- Nested parentheses
- Ternary operators
- Large expressions with many conditions

**Timestamps** (`BenchmarkConvertTimestamps`)
- Timestamp comparisons
- Date/time function calls
- DateTime operations

**String Operations** (`BenchmarkConvertStringOperations`)
- startsWith, endsWith, contains
- String concatenation
- Multiple string operations combined

**Query Analysis** (`BenchmarkAnalyzeQuery`)
- Index recommendation generation
- Pattern detection for optimization

**Options** (`BenchmarkConvertWithOptions`)
- Various conversion option combinations
- Schema usage overhead
- Max depth and output length limits

#### Benchmark Output

Benchmarks report:
- **Iterations**: Number of times the operation was executed
- **ns/op**: Nanoseconds per operation (lower is better)
- **B/op**: Bytes allocated per operation (lower is better)
- **allocs/op**: Number of allocations per operation (lower is better)

Example output:
```
BenchmarkConvertSimple/equality-12     1412060    848.7 ns/op    1593 B/op    25 allocs/op
BenchmarkConvertOperators/logical_and-12   943438    1255 ns/op    2154 B/op    36 allocs/op
```

#### CI/CD Integration

Benchmarks run automatically on every PR and push to main:
- Runs on Go 1.24.x
- Stores benchmark results in gh-pages branch for historical tracking
- Comments on PRs when performance changes >150%
- Generates visual charts at https://spandigital.github.io/cel2sql/dev/bench/
- Provides summary of all benchmark results

#### Comparing Benchmark Results

To compare benchmarks between two runs:

```bash
# Run benchmarks and save baseline
make bench-compare  # Saves to bench-new.txt

# Make changes to code

# Run benchmarks again and compare
mv bench-new.txt bench-old.txt
make bench-compare

# Install benchstat for detailed comparison (optional)
go install golang.org/x/perf/cmd/benchstat@latest
benchstat bench-old.txt bench-new.txt
```

#### When to Run Benchmarks

- Before and after performance optimizations
- When modifying core conversion logic
- When adding new features that may impact performance
- To validate that changes don't cause regressions

## Common Patterns

### Creating Type Providers
```go
schema := pg.NewSchema([]pg.FieldSchema{
    {Name: "field_name", Type: "text", Repeated: false},
    {Name: "array_field", Type: "text", Repeated: true},
    {Name: "json_field", Type: "jsonb", Repeated: false},
    {Name: "composite_field", Type: "composite", Schema: []pg.FieldSchema{...}},
})
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

### Query Analysis and Index Recommendations

cel2sql can analyze CEL expressions and recommend database indexes to optimize performance.

#### Using AnalyzeQuery

```go
ast, issues := env.Compile(`person.age > 18 && person.metadata.verified == true`)
if issues != nil && issues.Err() != nil {
    return issues.Err()
}

sql, recommendations, err := cel2sql.AnalyzeQuery(ast,
    cel2sql.WithSchemas(schemas))
if err != nil {
    return err
}

// Use the generated SQL
rows, err := db.Query("SELECT * FROM people WHERE " + sql)

// Review and apply index recommendations
for _, rec := range recommendations {
    fmt.Printf("Column: %s, Type: %s\n", rec.Column, rec.IndexType)
    fmt.Printf("Reason: %s\n", rec.Reason)
    fmt.Printf("Execute: %s\n\n", rec.Expression)
}
```

#### Index Recommendation Types

AnalyzeQuery detects patterns and recommends appropriate index types:

- **B-tree indexes**: Comparison operations (`==, >, <, >=, <=`)
  - Best for: Equality checks, range queries, sorting
  - Example: `person.age > 18` → B-tree on `person.age`

- **GIN indexes**: JSON/JSONB path operations, array operations
  - Best for: JSON field access, array membership, containment
  - Example: `person.metadata.verified == true` → GIN on `person.metadata`
  - Example: `"premium" in person.tags` → GIN on `person.tags`

- **GIN indexes with pg_trgm**: Regex pattern matching
  - Best for: Text search, pattern matching, fuzzy matching
  - Requires: PostgreSQL pg_trgm extension
  - Example: `person.email.matches(r"@example\.com$")` → GIN on `person.email`

#### IndexRecommendation Structure

```go
type IndexRecommendation struct {
    Column     string  // Full column name (e.g., "person.metadata")
    IndexType  string  // "BTREE", "GIN", or "GIST"
    Expression string  // Complete CREATE INDEX statement
    Reason     string  // Explanation of why this index is recommended
}
```

#### When to Use

- **Development**: Discover which indexes your queries need
- **Performance tuning**: Identify missing indexes causing slow queries
- **Production monitoring**: Analyze user-generated filter expressions

See `examples/index_analysis/` for a complete working example.

### Logging and Observability

cel2sql supports structured logging using Go's standard `log/slog` package (Go 1.21+).

Logging is optional and has **zero overhead** when not enabled (uses `slog.DiscardHandler` by default).

#### Enable Logging

```go
import "log/slog"

// JSON handler for production/machine parsing
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,
}))

// Text handler for development/debugging
logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

sql, err := cel2sql.Convert(ast,
    cel2sql.WithSchemas(schemas),
    cel2sql.WithLogger(logger))
```

#### What Gets Logged

- **JSON path detection decisions** - Table, field, operator selection (->>, ?)
- **Comprehension type identification** - all, exists, exists_one, filter, map
- **Schema lookups** - Hits/misses, field type detection
- **Performance metrics** - Conversion duration
- **Regex pattern transformations** - RE2 to POSIX conversion
- **Operator mapping decisions** - CEL to SQL operator conversion
- **Error contexts** - Full details when conversions fail

#### Log Levels

- **Debug**: Detailed conversion steps, operator mappings, schema lookups
- **Error**: Conversion failures with full context

#### Example Usage

```go
// Without logger - zero overhead (default)
sql, err := cel2sql.Convert(ast)

// With logging - detailed observability
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,
}))

sql, err := cel2sql.Convert(ast,
    cel2sql.WithSchemas(schemas),
    cel2sql.WithContext(ctx),
    cel2sql.WithLogger(logger))
```

See `examples/logging/` for a complete working example with both JSON and text handlers.

### Context Support (v2.10.0)

cel2sql supports context propagation for cancellation, timeouts, and observability integration.

Context support is **optional** and uses the functional options pattern.

#### Enable Context

```go
import "context"

// With timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

sql, err := cel2sql.Convert(ast,
    cel2sql.WithContext(ctx),
    cel2sql.WithSchemas(schemas))

// With cancellation
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

sql, err := cel2sql.Convert(ast,
    cel2sql.WithContext(ctx))
```

#### When Context Is Checked

Context cancellation is checked at key recursion points:
- **visit()** - Main traversal entry point
- **visitCall()** - Every function call
- **visitComprehension()** - Before processing comprehensions
- Individual comprehension handlers

#### Error Handling

If context is cancelled or times out during conversion:
```go
sql, err := cel2sql.Convert(ast, cel2sql.WithContext(ctx))
if err != nil {
    if errors.Is(err, context.Canceled) {
        // Conversion was cancelled
    } else if errors.Is(err, context.DeadlineExceeded) {
        // Conversion timed out
    }
}
```

#### Benefits

- **Cancellation**: Stop long-running conversions
- **Timeouts**: Protect against complex expressions
- **Observability**: Integrate with distributed tracing
- **Resource Cleanup**: Automatic cleanup on cancellation

#### Example Usage

```go
// Without context (default) - backward compatible
sql, err := cel2sql.Convert(ast)

// With context and other options
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

sql, err := cel2sql.Convert(ast,
    cel2sql.WithContext(ctx),
    cel2sql.WithSchemas(schemas),
    cel2sql.WithLogger(logger))
```

See `examples/context/` for a complete working example with timeouts and cancellation.

### Security Features (v2.10.0)

cel2sql includes comprehensive security protections against common attack vectors.

#### Field Name Validation

All field names are validated to prevent SQL injection:

**Validation Rules:**
- Maximum length: 63 characters (PostgreSQL NAMEDATALEN-1)
- Format: Must start with letter/underscore, contain only alphanumeric + underscore
- Reserved keywords: 60+ SQL keywords are rejected
- Empty strings: Not allowed

**Protection Against:**
```go
// ❌ These will be rejected:
field'; DROP TABLE users--
SELECT * FROM sensitive
user OR 1=1
```

**Validation happens at:**
- `visitSelect()` - Field names in select expressions
- `visitIdent()` - Identifier names to prevent reserved keywords

#### JSON Field Escaping

Single quotes in JSON field names are automatically escaped:

**Automatic Escaping:**
```go
// CEL with quote in field name
user.preferences.user'name == "test"

// Generated SQL (quotes escaped)
user.preferences->>'user''name' = 'test'
```

**Protection Against:**
- SQL injection via malicious JSON field names
- Field names like: `user' OR '1'='1`

**Escaping applied in:**
- `visitSelect()` - JSON path operators (->>)
- `visitHasFunction()` - JSON existence operators (?)
- `visitNestedJSONHas()` - jsonb_extract_path_text()
- `buildJSONPath*()` - All JSON path construction

#### ReDoS Protection

Comprehensive validation prevents Regular Expression Denial of Service attacks:

**Pattern Validation:**
- **Length limit**: 500 characters maximum
- **Nested quantifiers**: Detects patterns like `(a+)+`, `(a*)*`
- **Capture group limit**: Maximum 20 groups
- **Quantified alternation**: Blocks patterns like `(a|a)*b`
- **Nesting depth limit**: Maximum 10 levels

**Examples:**
```go
// ✅ Safe patterns (allowed)
field.matches(r"[a-z]+@[a-z]+\.[a-z]+")
field.matches(r"(?i)^user_\d+$")

// ❌ Dangerous patterns (rejected)
field.matches(r"(a+)+b")           // Nested quantifiers
field.matches(r"(a|a)*b")          // Quantified alternation
field.matches(r"(((((((((((a"))    // Excessive nesting
```

**Protection Against:**
- Catastrophic backtracking (CWE-1333)
- CPU exhaustion from complex patterns
- Service disruption from malicious regex

#### Resource Exhaustion Protection

cel2sql includes multiple layers of protection against resource exhaustion attacks (CWE-400):

**Recursion Depth Limits:**
- **Default limit**: 100 levels of expression nesting
- **Configurable**: Use `WithMaxDepth()` to adjust
- **Protection**: Prevents stack overflow from deeply nested expressions (CWE-674)

**SQL Output Length Limits:**
- **Default limit**: 50,000 characters of generated SQL
- **Configurable**: Use `WithMaxOutputLength()` to adjust
- **Protection**: Prevents memory exhaustion from extremely large SQL queries

**Comprehension Depth Limits:**
- **Fixed limit**: 3 levels of nested comprehensions
- **Protection**: Prevents resource exhaustion from deeply nested UNNEST/subquery operations

**Byte Array Length Limits:**
- **Fixed limit**: 10,000 bytes maximum
- **Applies to**: Non-parameterized mode (hex encoding)
- **Protection**: Prevents memory exhaustion from large hex-encoded SQL strings (each byte → ~4 chars)
- **Note**: Parameterized mode bypasses this limit (bytes passed directly to database driver)

**Examples:**
```go
// Use default limits (recommended)
sql, err := cel2sql.Convert(ast)

// Custom recursion depth limit
sql, err := cel2sql.Convert(ast,
    cel2sql.WithMaxDepth(150))

// Custom SQL output length limit
sql, err := cel2sql.Convert(ast,
    cel2sql.WithMaxOutputLength(100000))

// Combine multiple limits
sql, err := cel2sql.Convert(ast,
    cel2sql.WithMaxDepth(75),
    cel2sql.WithMaxOutputLength(25000),
    cel2sql.WithContext(ctx))
```

**Protection Against:**
- Stack overflow from deeply nested expressions
- Memory exhaustion from extremely large SQL output
- CPU/memory exhaustion from deeply nested comprehensions
- Memory exhaustion from large hex-encoded byte arrays
- DoS attacks via resource consumption

#### Context Timeouts

Use context timeouts as defense-in-depth:

```go
// Protect against complex expressions
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

sql, err := cel2sql.Convert(ast,
    cel2sql.WithContext(ctx),
    cel2sql.WithSchemas(schemas))
```

#### Security Best Practices

1. **Always validate user input** before passing to CEL
2. **Use context timeouts** for user-provided expressions
3. **Enable logging** to monitor conversion patterns
4. **Keep schemas minimal** - only expose necessary fields
5. **Use prepared statements** when executing generated SQL
6. **Test edge cases** with your specific field names

For detailed security information, see the security documentation.

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
