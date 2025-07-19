# cel2sql

cel2sql converts [CEL (Common Expression Language)](https://cel.dev/) [Github Reoi](https://opensource.google/projects/cel) to SQL condition.
It is specifically targeting PostgreSQL standard SQL.

## Latest Release - v2.8.0

🚀 **Latest New Features:**
- **🔥 Regex Pattern Matching Support**: Full support for CEL `matches()` function with RE2 to POSIX regex conversion
- **Enhanced JSON Field Existence**: Improved `has()` macro support for JSON/JSONB fields with proper path checking
- **PostgreSQL Pattern Compatibility**: Automatic conversion of RE2 regex patterns to PostgreSQL-compatible POSIX format
- **Advanced Text Search**: Support for complex regex operations with proper escaping and PostgreSQL syntax

**Recent Features (v2.6.0):**
- **🔥 JSON/JSONB Comprehensions Support**: Full support for CEL comprehensions on JSON/JSONB arrays
- **Advanced JSON Array Operations**: Support for `exists()`, `all()`, `exists_one()` on JSON/JSONB arrays  
- **Numeric JSON Field Casting**: Automatic casting of numeric JSON fields (e.g., `(score)::numeric`)
- **Nested JSON Array Access**: Support for comprehensions on nested JSON arrays (e.g., `settings.permissions`)
- **JSON Type Safety**: Null and type checks for JSON/JSONB comprehensions using `jsonb_typeof()`

**Previous Features (v2.5.0):**
- **🔥 CEL Comprehensions Support**: Full support for `all()`, `exists()`, `exists_one()`, `filter()`, and `map()` comprehensions
- **PostgreSQL UNNEST Integration**: Comprehensions are converted to efficient PostgreSQL SQL using `UNNEST()` and array functions
- **Nested Comprehensions**: Support for complex nested comprehensions like `employees.exists(e, e.skills.exists(s, s == 'Go'))`
- **Array Field Access**: Direct comprehensions on array fields in PostgreSQL schemas

**Previous Features (v2.4.0):**
- **Comprehensive JSON/JSONB Support**: Native PostgreSQL JSON path operations
- **Dynamic Schema Loading**: Load table schemas directly from PostgreSQL databases
- **Enhanced Type System**: Improved PostgreSQL type mappings and array support
- **Testcontainer Integration**: Full test coverage with real PostgreSQL databases

**Key Improvements:**
- PostgreSQL-optimized SQL generation (single quotes, proper functions)
- JSON field access: `user.preferences.theme` → `user.preferences->>'theme'`
- Array operations: `size(array)` → `ARRAY_LENGTH(array, 1)`
- String operations: `contains()` → `POSITION(...) > 0`
- CEL comprehensions: `list.all(x, x > 0)` → `NOT EXISTS (SELECT 1 FROM UNNEST(list) AS x WHERE NOT (x > 0))`
- All tests pass with comprehensive integration coverage

## Usage

```go
import (
    "context"
    "fmt"
    
    "github.com/spandigital/cel2sql/v2"
    "github.com/spandigital/cel2sql/v2/pg"
    "github.com/spandigital/cel2sql/v2/sqltypes"
    "github.com/google/cel-go/cel"
)

// PostgreSQL table schema definition
employeeSchema := pg.Schema{
    {Name: "name", Type: "text", Repeated: false},
    {Name: "hired_at", Type: "timestamp with time zone", Repeated: false},
    {Name: "age", Type: "integer", Repeated: false},
    {Name: "active", Type: "boolean", Repeated: false},
}

// Prepare CEL environment
env, _ := cel.NewEnv(
    cel.CustomTypeProvider(pg.NewTypeProvider(map[string]pg.Schema{
        "Employee": employeeSchema,
    })),
    cel.Variable("employee", cel.ObjectType("Employee")),
)

// Convert CEL to SQL
ast, _ := env.Compile(`employee.name == "John Doe" && employee.hired_at >= current_timestamp() - duration("24h")`)
sqlCondition, _ := cel2sql.Convert(ast)

fmt.Println(sqlCondition) // employee.name = 'John Doe' AND employee.hired_at >= CURRENT_TIMESTAMP - INTERVAL '1 DAY'
```

## Dynamic Schema Loading

cel2sql supports dynamically loading table schemas from a PostgreSQL database:

```go
import (
    "context"
    "fmt"
    
    "github.com/spandigital/cel2sql/v2"
    "github.com/spandigital/cel2sql/v2/pg"
    "github.com/spandigital/cel2sql/v2/sqltypes"
    "github.com/google/cel-go/cel"
)

func main() {
    ctx := context.Background()
    
    // Create a type provider with database connection
    provider, err := pg.NewTypeProviderWithConnection(ctx, "postgres://user:pass@localhost/mydb?sslmode=disable")
    if err != nil {
        panic(err)
    }
    defer provider.Close()
    
    // Load table schema dynamically from database
    err = provider.LoadTableSchema(ctx, "employees")
    if err != nil {
        panic(err)
    }
    
    // Use the loaded schema in CEL environment
    env, err := cel.NewEnv(
        cel.CustomTypeProvider(provider),
        cel.Variable("employee", cel.ObjectType("employees")),
    )
    if err != nil {
        panic(err)
    }
    
    // Convert CEL to SQL using the dynamically loaded schema
    ast, issues := env.Compile(`employee.name == "John Doe" && employee.age > 30`)
    if issues != nil && issues.Err() != nil {
        panic(issues.Err())
    }
    
    sqlCondition, err := cel2sql.Convert(ast)
    if err != nil {
        panic(err)
    }
    
    fmt.Println(sqlCondition)
    // Output: employee.name = 'John Doe' AND employee.age > 30
}
```

This approach is particularly useful when:
- Database schemas change frequently
- You want to avoid manually defining schemas
- Working with multiple tables with different structures
- Building dynamic query builders

## Type Conversion

CEL Type    | PostgreSQL Data Type
----------- | ----------------------------------
`int`       | `bigint`
`uint`      | Unsupported but treated as `bigint`
`double`    | `double precision`
`bool`      | `boolean`
`string`    | `text`
`bytes`     | `bytea`
`list`      | `ARRAY`
`map`       | `JSONB` (for complex objects)
`null_type` | `NULL`
`timestamp` | `timestamp with time zone`
`duration`  | `INTERVAL` 

## JSON/JSONB Support

cel2sql provides comprehensive support for PostgreSQL JSON and JSONB columns:

```go
// Example with JSON/JSONB fields
userSchema := pg.Schema{
    {Name: "name", Type: "text", Repeated: false},
    {Name: "preferences", Type: "jsonb", Repeated: false},
    {Name: "profile", Type: "json", Repeated: false},
}

env, _ := cel.NewEnv(
    cel.CustomTypeProvider(pg.NewTypeProvider(map[string]pg.Schema{
        "User": userSchema,
    })),
    cel.Variable("user", cel.ObjectType("User")),
)

// CEL expressions automatically convert to PostgreSQL JSON path operations
ast, _ := env.Compile(`user.preferences.theme == "dark"`)
sqlCondition, _ := cel2sql.Convert(ast)
fmt.Println(sqlCondition) // user.preferences->>'theme' = 'dark'

// Nested JSON access
ast, _ = env.Compile(`user.profile.settings.notifications == "enabled"`)
sqlCondition, _ = cel2sql.Convert(ast)
fmt.Println(sqlCondition) // user.profile->>'settings'->>'notifications' = 'enabled'
```

**Supported JSON Operations:**
- Field access: `user.preferences.theme` → `user.preferences->>'theme'`
- Nested access: `user.profile.settings.key` → `user.profile->>'settings'->>'key'`
- Works with both `json` and `jsonb` column types
- Automatically detects JSON columns and applies proper PostgreSQL syntax 

## Regex Pattern Matching

cel2sql provides comprehensive support for CEL `matches()` function with automatic RE2 to POSIX regex conversion:

```go
// Define schema with text fields for pattern matching
schema := pg.Schema{
    {Name: "email", Type: "text", Repeated: false},
    {Name: "phone", Type: "text", Repeated: false},
    {Name: "description", Type: "text", Repeated: false},
}

env, _ := cel.NewEnv(
    cel.CustomTypeProvider(pg.NewTypeProvider(map[string]pg.Schema{
        "Contact": schema,
    })),
    cel.Variable("contact", cel.ObjectType("Contact")),
)

// Email validation pattern
ast, _ := env.Compile(`contact.email.matches(r"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$")`)
sqlCondition, _ := cel2sql.Convert(ast)
fmt.Println(sqlCondition) 
// Output: contact.email ~ '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'

// Phone number pattern matching
ast, _ = env.Compile(`contact.phone.matches(r"^\+?1?[-.\s]?\(?([0-9]{3})\)?[-.\s]?([0-9]{3})[-.\s]?([0-9]{4})$")`)
sqlCondition, _ = cel2sql.Convert(ast)
fmt.Println(sqlCondition)
// Output: contact.phone ~ '^\+?1?[-.\s]?\(?([0-9]{3})\)?[-.\s]?([0-9]{3})[-.\s]?([0-9]{4})$'

// Case-insensitive pattern matching
ast, _ = env.Compile(`contact.description.matches(r"(?i)urgent|priority")`)
sqlCondition, _ = cel2sql.Convert(ast)
fmt.Println(sqlCondition)
// Output: contact.description ~* 'urgent|priority'
```

**Supported Regex Features:**
- **RE2 to POSIX Conversion**: Automatic conversion of CEL RE2 patterns to PostgreSQL-compatible POSIX regex
- **Case-insensitive Matching**: `(?i)` flag converts to PostgreSQL `~*` operator
- **Standard Patterns**: Email validation, phone numbers, URLs, and custom text patterns
- **Escape Handling**: Proper escaping of special characters for PostgreSQL
- **Pattern Optimization**: Efficient regex compilation and execution in PostgreSQL

**Regex Operators:**
- `~` - Case-sensitive POSIX regex match
- `~*` - Case-insensitive POSIX regex match (when `(?i)` flag detected)
- `!~` - Case-sensitive POSIX regex non-match
- `!~*` - Case-insensitive POSIX regex non-match

## Enhanced JSON Field Existence

The `has()` macro provides enhanced support for checking JSON/JSONB field existence:

```go
// Check if JSON field exists
ast, _ := env.Compile(`has(user.preferences.theme)`)
sqlCondition, _ := cel2sql.Convert(ast)
fmt.Println(sqlCondition) 
// Output: user.preferences ? 'theme'

// Check nested JSON field existence
ast, _ = env.Compile(`has(user.profile.settings.notifications)`)
sqlCondition, _ := cel2sql.Convert(ast)
fmt.Println(sqlCondition)
// Output: user.profile->'settings' ? 'notifications'

// Combined existence and value check
ast, _ = env.Compile(`has(user.preferences.theme) && user.preferences.theme == "dark"`)
sqlCondition, _ = cel2sql.Convert(ast)
fmt.Println(sqlCondition)
// Output: user.preferences ? 'theme' AND user.preferences->>'theme' = 'dark'
```

**JSON Field Existence Features:**
- **Path-aware Checking**: Properly handles nested JSON path existence
- **JSONB Optimization**: Uses PostgreSQL's efficient `?` operator for JSONB fields
- **Null Safety**: Prevents errors when accessing non-existent JSON fields
- **Combined Operations**: Works seamlessly with value comparisons and other JSON operations

## Supported CEL Operators/Functions

<table style="width: 100%; border: solid 1px;">
  <col style="width: 15%;">
  <col style="width: 40%;">
  <col style="width: 45%;">
  <tr>
    <th>Symbol</th>
    <th>Type</th>
    <th>SQL</th>
  </tr>
  <tr>
    <th rowspan="1">
      !
    </th>
    <td>
      (bool) -> bool
    </td>
    <td>
      <code>NOT</code> bool
    </td>
  </tr>
  <tr>
    <th rowspan="2">
      - (unary)
    </th>
    <td>
      (int) -> int
    </td>
    <td>
      <code>-</code>int
    </td>
  </tr>
  <tr>
    <td>
      (double) -> double
    </td>
    <td>
      <code>-</code>double
    </td>
  </tr>
  <tr>
    <th rowspan="3">
      !=
    </th>
    <td>
      (A, A) -> bool
    </td>
    <td>
      A <code>!=</code> A
    </td>
  </tr>
  <tr>
    <td>
      (bool, bool) -> bool
    </td>
    <td>
      bool <code>IS NOT</code> bool
    </td>
  </tr>
  <tr>
    <td>
      (A, null) -> bool
    </td>
    <td>
      A <code>IS NOT NULL</code>
    </td>
  </tr>
  <tr>
    <th rowspan="1">
      %
    </th>
    <td>
      (int, int) -> int
    </td>
    <td>
      <code>MOD(</code>int<code>, </code>int<code>)</code>
    </td>
  </tr>
  <tr>
    <th rowspan="1">
      &&
    </th>
    <td>
      (bool, bool) -> bool
    </td>
    <td>
      bool <code>AND</code> bool
    </td>
  </tr>
  <tr>
    <th rowspan="2">
      *
    </th>
    <td>
      (int, int) -> int
    </td>
    <td>
      int <code>*</code> int
    </td>
  </tr>
  <tr>
    <td>
      (double, double) -> double
    </td>
    <td>
      double <code>*</code> double
    </td>
  </tr>
  <tr>
    <th rowspan="7">
      +
    </th>
    <td>
      (int, int) -> int
    </td>
    <td>
      int <code>+</code> int 
    </td>
  </tr>
  <tr>
    <td>
      (double, double) -> double
    </td>
    <td>
      double <code>+</code> double 
    </td>
  </tr>
  <tr>
    <td>
      (string, string) -> string
    </td>
    <td>
      string <code>||</code> string 
    </td>
  </tr>
  <tr>
    <td>
      (bytes, bytes) -> bytes
    </td>
    <td>
      bytes <code>||</code> bytes 
    </td>
  </tr>
  <tr>
    <td>
      (list(A), list(A)) -> list(A)
    </td>
    <td>
      list(A) <code>||</code> list(A) 
    </td>
  </tr>
  <tr>
    <td>
      (google.protobuf.Timestamp, google.protobuf.Duration) -> google.protobuf.Timestamp
    </td>
    <td>
      <code>TIMESTAMP_ADD(</code>timestamp<code>, INTERVAL </code>duration<code> date_part)</code>
    </td>
  </tr>
  <tr>
    <td>
      (google.protobuf.Duration, google.protobuf.Timestamp) -> google.protobuf.Timestamp
    </td>
    <td>
      <code>TIMESTAMP_ADD(</code>timestamp<code>, INTERVAL </code>duration<code> date_part)</code>
    </td>
  </tr>
  <tr>
    <th rowspan="3">
      - (binary)
    </th>
    <td>
      (int, int) -> int
    </td>
    <td>
      int <code>-</code> int
    </td>
  </tr>
  <tr>
    <td>
      (double, double) -> double
    </td>
    <td>
      double <code>-</code> double
    </td>
  </tr>
  <tr>
    <td>
      (google.protobuf.Timestamp, google.protobuf.Duration) -> google.protobuf.Timestamp
    </td>
    <td>
      <code>TIMESTAMP_SUB(</code>timestamp<code>, INTERVAL </code>duration<code> date_part)</code>
    </td>
  </tr>
  <tr>
    <th rowspan="2">
      /
    </th>
    <td>
      (int, int) -> int
    </td>
    <td>
      int <code>/</code> int
    </td>
  </tr>
  <tr>
    <td>
      (double, double) -> double
    </td>
    <td>
      double <code>/</code> double
    </td>
  </tr>
  <tr>
    <th rowspan="6">
      &lt;=
    </th>
    <td>
      (bool, bool) -> bool
    </td>
    <td>
      bool <code><=</code> bool
    </td>
  </tr>
  <tr>
    <td>
      (int, int) -> bool
    </td>
    <td>
      int <code><=</code> int
    </td>
  </tr>
  <tr>
    <td>
      (double, double) -> bool
    </td>
    <td>
      double <code><=</code> double
    </td>
  </tr>
  <tr>
    <td>
      (string, string) -> bool
    </td>
    <td>
      string <code><=</code> string
    </td>
  </tr>
  <tr>
    <td>
      (bytes, bytes) -> bool
    </td>
    <td>
      bytes <code><=</code> bytes
    </td>
  </tr>
  <tr>
    <td>
      (google.protobuf.Timestamp, google.protobuf.Timestamp) -> bool
    </td>
    <td>
      timestamp <code><=</code> timestamp
    </td>
  </tr>
  <tr>
    <th rowspan="6">
      &lt;
    </th>
    <td>
      (bool, bool) -> bool
    </td>
    <td>
      bool <code><</code> bool
    </td>
  </tr>
  <tr>
    <td>
      (int, int) -> bool
    </td>
    <td>
      int <code><</code> int
    </td>
  </tr>
  <tr>
    <td>
      (double, double) -> bool
    </td>
    <td>
      double <code><</code> double
    </td>
  </tr>
  <tr>
    <td>
      (string, string) -> bool
    </td>
    <td>
      string <code><</code> string
    </td>
  </tr>
  <tr>
    <td>
      (bytes, bytes) -> bool
    </td>
    <td>
      bytes <code><</code> bytes
    </td>
  </tr>
  <tr>
    <td>
      (google.protobuf.Timestamp, google.protobuf.Timestamp) -> bool
    </td>
    <td>
      timestamp <code><</code> timestamp
    </td>
  </tr>
  <tr>
    <th rowspan="3">
      ==
    </th>
    <td>
      (A, A) -> bool
    </td>
    <td>
      A <code>=</code> A
    </td>
  </tr>
  <tr>
    <td>
      (bool, bool) -> bool
    </td>
    <td>
      A <code>IS</code> A
    </td>
  </tr>
  <tr>
    <td>
      (A, null) -> bool
    </td>
    <td>
      A <code>IS NULL</code>
    </td>
  </tr>
  <tr>
    <th rowspan="6">
      &gt;=
    </th>
    <td>
      (bool, bool) -> bool
    </td>
    <td>
      bool <code>>=</code> bool
    </td>
  </tr>
  <tr>
    <td>
      (int, int) -> bool
    </td>
    <td>
      int <code>>=</code> int
    </td>
  </tr>
  <tr>
    <td>
      (double, double) -> bool
    </td>
    <td>
      double <code>>=</code> double
    </td>
  </tr>
  <tr>
    <td>
      (string, string) -> bool
    </td>
    <td>
      string <code>>=</code> string
    </td>
  </tr>
  <tr>
    <td>
      (bytes, bytes) -> bool
    </td>
    <td>
      bytes <code>>=</code> bytes
    </td>
  </tr>
  <tr>
    <td>
      (google.protobuf.Timestamp, google.protobuf.Timestamp) -> bool
    </td>
    <td>
      timestamp <code>>=</code> timestamp
    </td>
  </tr>
  <tr>
    <th rowspan="6">
      &gt;
    </th>
    <td>
      (bool, bool) -> bool
    </td>
    <td>
      bool <code>></code> bool
    </td>
  </tr>
  <tr>
    <td>
      (int, int) -> bool
    </td>
    <td>
      int <code>></code> int
    </td>
  </tr>
  <tr>
    <td>
      (double, double) -> bool
    </td>
    <td>
      double <code>></code> double
    </td>
  </tr>
  <tr>
    <td>
      (string, string) -> bool
    </td>
    <td>
      string <code>></code> string
    </td>
  </tr>
  <tr>
    <td>
      (bytes, bytes) -> bool
    </td>
    <td>
      bytes <code>></code> bytes
    </td>
  </tr>
  <tr>
    <td>
      (google.protobuf.Timestamp, google.protobuf.Timestamp) -> bool
    </td>
    <td>
      timestamp <code>></code> timestamp
    </td>
  </tr>
  <tr>
    <th rowspan="1">
      ? :
    </th>
    <td>
      (bool, A, A) -> A
    </td>
    <td>
      <code>IF(</code>bool<code>, </code>A<code>, </code>A<code>)</code>
    </td>
  </tr>
  <tr>
    <th rowspan="2">
      [ ]
    </th>
    <td>
      (list(A), int) -> A
    </td>
    <td>
      list<code>[OFFSET(</code>int<code>)]</code>
    </td>
  </tr>
  <tr>
    <td>
      (map(A, B), A) -> B
    </td>
    <td>
      map<code>.`</code>A<code>`</code>
    </td>
  </tr>
  <tr>
    <th rowspan="1">
      in
    </th>
    <td>
      (A, list(A)) -> bool
    </td>
    <td>
      A <code>IN UNNEST(</code>list<code>)</code>
    </td>
  </tr>
  <tr>
    <th rowspan="1">
      ||
    </th>
    <td>
      (bool, bool) -> bool
    </td>
    <td>
      bool <code>OR</code> bool
    </td>
  </tr>
  <tr>
    <th rowspan="2">
      bool
    </th>
    <td>
      (int) -> bool
    </td>
    <td>
      <code>CAST(</code>int<code> AS BOOL)</code>
    </td>
  </tr>
  <tr>
    <td>
      (string) -> bool
    </td>
    <td>
      <code>CAST(</code>string<code> AS BOOL)</code>
    </td>
  </tr>
  <tr>
    <th rowspan="1">
      bytes
    </th>
    <td>
      (string) -> bytes
    </td>
    <td>
      <code>CAST(</code>string<code>AS BYTES)</code>
    </td>
  </tr>
  <tr>
    <th rowspan="1">
      contains
    </th>
    <td>
      string.(string) -> bool
    </td>
    <td>
      <code>POSITION(</code>string<code> IN </code>string<code>) > 0</code>
    </td>
  </tr>
  <tr>
    <th rowspan="2">
      double
    </th>
    <td>
      (int) -> double
    </td>
    <td>
      <code>CAST(</code>int<code> AS FLOAT64)</code>
    </td>
  </tr>
  <tr>
    <td>
      (string) -> double
    </td>
    <td>
      <code>CAST(</code>string<code> AS FLOAT64)</code>
    </td>
  </tr>
  <tr>
    <th rowspan="1">
      duration
    </th>
    <td>
      (string) -> google.protobuf.Duration
    </td>
    <td>
      <code>INTERVAL </code>duration<code> date_part</code>
    </td>
  </tr>
  <tr>
    <th rowspan="1">
      endsWith
    </th>
    <td>
      string.(string) -> bool
    </td>
    <td>
      <code>ENDS_WITH(</code>string<code>, </code>string<code>)</code>
    </td>
  </tr>
  <tr>
    <th rowspan="2">
      getDate
    </th>
    <td>
      google.protobuf.Timestamp.() -> int
    </td>
    <td>
      <code>EXTRACT(DAY FROM </code>timestamp<code>)</code>
    </td>
  </tr>
  <tr>
    <td>
      google.protobuf.Timestamp.(string) -> int
    </td>
    <td>
      <code>EXTRACT(DAY FROM </code>timestamp<code> AT </code>string<code>)</code>
    </td>
  </tr>
  <tr>
    <th rowspan="2">
      getDayOfMonth
    </th>
    <td>
      google.protobuf.Timestamp.() -> int
    </td>
    <td>
      <code>EXTRACT(DAY FROM </code>timestamp<code>) - 1</code>
    </td>
  </tr>
  <tr>
    <td>
      google.protobuf.Timestamp.(string) -> int
    </td>
    <td>
      <code>EXTRACT(DAY FROM </code>timestamp<code> AT </code>string<code>) - 1</code>
    </td>
  </tr>
  <tr>
    <th rowspan="2">
      getDayOfWeek
    </th>
    <td>
      google.protobuf.Timestamp.() -> int
    </td>
    <td>
      <code>EXTRACT(DAYOFWEEK FROM </code>timestamp<code>) - 1</code>
    </td>
  </tr>
  <tr>
    <td>
      google.protobuf.Timestamp.(string) -> int
    </td>
    <td>
      <code>EXTRACT(DAYOFWEEK FROM </code>timestamp<code> AT </code>string<code>) - 1</code>
    </td>
  </tr>
  <tr>
    <th rowspan="2">
      getDayOfYear
    </th>
    <td>
      google.protobuf.Timestamp.() -> int
    </td>
    <td>
      <code>EXTRACT(DAYOFYEAR FROM </code>timestamp<code>) - 1</code>
    </td>
  </tr>
  <tr>
    <td>
      google.protobuf.Timestamp.(string) -> int
    </td>
    <td>
      <code>EXTRACT(DAYOFYEAR FROM </code>timestamp<code> AT </code>string<code>) - 1</code>
    </td>
  </tr>
  <tr>
    <th rowspan="2">
      getFullYear
    </th>
    <td>
      google.protobuf.Timestamp.() -> int
    </td>
    <td>
      <code>EXTRACT(YEAR FROM </code>timestamp<code>)</code>
    </td>
  </tr>
  <tr>
    <td>
      google.protobuf.Timestamp.(string) -> int
    </td>
    <td>
      <code>EXTRACT(YEAR FROM </code>timestamp<code> AT </code>string<code>)</code>
    </td>
  </tr>
  <tr>
    <th rowspan="2">
      getHours
    </th>
    <td>
      google.protobuf.Timestamp.() -> int
    </td>
    <td>
      <code>EXTRACT(HOUR FROM </code>timestamp<code>)</code>
    </td>
  </tr>
  <tr>
    <td>
      google.protobuf.Timestamp.(string) -> int
    </td>
    <td>
      <code>EXTRACT(HOUR FROM </code>timestamp<code> AT </code>string<code>)</code>
    </td>
  </tr>
  <tr>
    <th rowspan="2">
      getMilliseconds
    </th>
    <td>
      google.protobuf.Timestamp.() -> int
    </td>
    <td>
      <code>EXTRACT(MILLISECOND FROM </code>timestamp<code>)</code>
    </td>
  </tr>
  <tr>
    <td>
      google.protobuf.Timestamp.(string) -> int
    </td>
    <td>
      <code>EXTRACT(MILLISECOND FROM </code>timestamp<code> AT </code>string<code>)</code>
    </td>
  </tr>
  <tr>
    <th rowspan="2">
      getMinutes
    </th>
    <td>
      google.protobuf.Timestamp.() -> int
    </td>
    <td>
      <code>EXTRACT(MINUTE FROM </code>timestamp<code>)</code>
    </td>
  </tr>
  <tr>
    <td>
      google.protobuf.Timestamp.(string) -> int
    </td>
    <td>
      <code>EXTRACT(MINUTE FROM </code>timestamp<code> AT </code>string<code>)</code>
    </td>
  </tr>
  <tr>
    <th rowspan="2">
      getMonth
    </th>
    <td>
      google.protobuf.Timestamp.() -> int
    </td>
    <td>
      <code>EXTRACT(MONTH FROM </code>timestamp<code>) - 1</code>
    </td>
  </tr>
  <tr>
    <td>
      google.protobuf.Timestamp.(string) -> int
    </td>
    <td>
      <code>EXTRACT(MONTH FROM </code>timestamp<code> AT </code>string<code>) - 1</code>
    </td>
  </tr>
  <tr>
    <th rowspan="2">
      getSeconds
    </th>
    <td>
      google.protobuf.Timestamp.() -> int
    </td>
    <td>
      <code>EXTRACT(SECOND FROM </code>timestamp<code>)</code>
    </td>
  </tr>
  <tr>
    <td>
      google.protobuf.Timestamp.(string) -> int
    </td>
    <td>
      <code>EXTRACT(SECOND FROM </code>timestamp<code> AT </code>string<code>)</code>
    </td>
  </tr>
  <tr>
    <th rowspan="1">
      has
    </th>
    <td>
      (map) -> bool
    </td>
    <td>
      JSON/JSONB field <code>?</code> 'key' <code>OR</code> nested_path <code>-></code> 'parent' <code>?</code> 'key'
    </td>
  </tr>
  <tr>
    <th rowspan="4">
      int
    </th>
    <td>
      (bool) -> int
    </td>
    <td>
      <code>CAST(</code>bool<code> AS INT64)</code>
    </td>
  </tr>
  <tr>
    <td>
      (double) -> int
    </td>
    <td>
      <code>CAST(</code>double<code> AS INT64)</code>
    </td>
  </tr>
  <tr>
    <td>
      (string) -> int
    </td>
    <td>
      <code>CAST(</code>string<code> AS INT64)</code>
    </td>
  </tr>
  <tr>
    <td>
      (google.protobuf.Timestamp) -> int
    </td>
    <td>
      <code>UNIX_SECONDS(</code>timestamp<code>)</code>
    </td>
  </tr>
  <tr>
    <th rowspan="1">
      matches
    </th>
    <td>
      string.(string) -> bool
    </td>
    <td>
      string <code>~</code> regex_pattern <code>OR</code> string <code>~*</code> regex_pattern (case-insensitive)
    </td>
  </tr>
  <tr>
    <th rowspan="3">
      size
    </th>
    <td>
      (string) -> int
    </td>
    <td>
      <code>CHAR_LENGTH(</code>string<code>)</code>
    </td>
  </tr>
  <tr>
    <td>
      (bytes) -> int
    </td>
    <td>
      <code>BYTE_LENGTH(</code>bytes<code>)</code>
    </td>
  </tr>
  <tr>
    <td>
      (list(A)) -> int
    </td>
    <td>
      <code>ARRAY_LENGTH(</code>list<code>, 1)</code>
    </td>
  </tr>
  <tr>
    <th rowspan="1">
      startsWith
    </th>
    <td>
      string.(string) -> bool
    </td>
    <td>
      <code>STARTS_WITH</code>string<code>, </code>string<code>)</code>
    </td>
  </tr>
  <tr>
    <th rowspan="5">
      string
    </th>
    <td>
      (bool) -> string
    </td>
    <td>
      <code>CAST(</code>bool<code> AS STRING)</code>
    </td>
  </tr>
  <tr>
    <td>
      (int) -> string
    </td>
    <td>
      <code>CAST(</code>int<code> AS STRING)</code>
    </td>
  </tr>
  <tr>
    <td>
      (double) -> string
    </td>
    <td>
      <code>CAST(</code>double<code> AS STRING)</code>
    </td>
  </tr>
  <tr>
    <td>
      (bytes) -> string
    </td>
    <td>
      <code>CAST(</code>bytes<code> AS STRING)</code>
    </td>
  </tr>
  <tr>
    <td>
      (timestamp) -> string
    </td>
    <td>
      <code>CAST(</code>timestamp<code> AS STRING)</code>
    </td>
  </tr>
  <tr>
    <th rowspan="1">
      timestamp
    </th>
    <td>
      (string) -> google.protobuf.Timestamp
    </td>
    <td>
      <code>TIMESTAMP(</code>string<code>)</code>
    </td>
  </tr>
</table>

## Standard SQL Types/Functions

cel2sql supports time related types bellow.

- `DATE`
- `TIME`
- `DATETIME`

cel2sql contains time related functions bellow.

- `current_date()`
- `current_time()`
- `current_datetime()`
- `current_timestamp()`
- `interval(N, date_part)`

## CEL Comprehensions

cel2sql now supports CEL comprehensions for working with lists and arrays. Comprehensions are converted to PostgreSQL-compatible SQL using `UNNEST()` and various array functions.

### Supported Comprehension Types

| CEL Expression | Description | Generated SQL Pattern |
|----------------|-------------|----------------------|
| `list.all(x, condition)` | All elements satisfy condition | `NOT EXISTS (SELECT 1 FROM UNNEST(list) AS x WHERE NOT (condition))` |
| `list.exists(x, condition)` | At least one element satisfies condition | `EXISTS (SELECT 1 FROM UNNEST(list) AS x WHERE condition)` |
| `list.exists_one(x, condition)` | Exactly one element satisfies condition | `(SELECT COUNT(*) FROM UNNEST(list) AS x WHERE condition) = 1` |
| `list.filter(x, condition)` | Return elements that satisfy condition | `ARRAY(SELECT x FROM UNNEST(list) AS x WHERE condition)` |
| `list.map(x, transform)` | Transform all elements | `ARRAY(SELECT transform FROM UNNEST(list) AS x)` |

### Examples

#### Simple Array Comprehensions

```go
// Check if all numbers are positive
cel: [1, 2, 3, 4, 5].all(x, x > 0)
sql: NOT EXISTS (SELECT 1 FROM UNNEST(ARRAY[1, 2, 3, 4, 5]) AS x WHERE NOT (x > 0))

// Check if any number is greater than 3
cel: [1, 2, 3, 4, 5].exists(x, x > 3)
sql: EXISTS (SELECT 1 FROM UNNEST(ARRAY[1, 2, 3, 4, 5]) AS x WHERE x > 3)

// Filter even numbers
cel: [1, 2, 3, 4, 5].filter(x, x % 2 == 0)
sql: ARRAY(SELECT x FROM UNNEST(ARRAY[1, 2, 3, 4, 5]) AS x WHERE MOD(x, 2) = 0)

// Double all numbers
cel: [1, 2, 3, 4, 5].map(x, x * 2)
sql: ARRAY(SELECT x * 2 FROM UNNEST(ARRAY[1, 2, 3, 4, 5]) AS x)
```

#### PostgreSQL Schema-based Comprehensions

```go
// Define schema with array fields
schema := pg.Schema{
    {Name: "name", Type: "text", Repeated: false},
    {Name: "age", Type: "bigint", Repeated: false},
    {Name: "skills", Type: "text", Repeated: true}, // Array field
}

// CEL expressions with comprehensions
cel: employees.all(e, e.age >= 18)
sql: NOT EXISTS (SELECT 1 FROM UNNEST(employees) AS e WHERE NOT (e.age >= 18))

cel: employees.filter(e, e.age > 30).map(e, e.name)
sql: ARRAY(SELECT e.name FROM UNNEST(ARRAY(SELECT e FROM UNNEST(employees) AS e WHERE e.age > 30)) AS e)

cel: emp.skills.exists(s, s == 'Go')
sql: EXISTS (SELECT 1 FROM UNNEST(emp.skills) AS s WHERE s = 'Go')
```

#### Nested Comprehensions

```go
// Check if any employee has Go skills (nested comprehension)
cel: employees.exists(e, e.skills.exists(s, s == 'Go'))
sql: EXISTS (SELECT 1 FROM UNNEST(employees) AS e WHERE EXISTS (SELECT 1 FROM UNNEST(e.skills) AS s WHERE s = 'Go'))

// Filter employees with all high scores
cel: employees.filter(e, e.scores.all(s, s >= 80))
sql: ARRAY(SELECT e FROM UNNEST(employees) AS e WHERE NOT EXISTS (SELECT 1 FROM UNNEST(e.scores) AS s WHERE NOT (s >= 80)))
```

#### Working with Composite Types

```go
// Define nested schema
addressSchema := pg.Schema{
    {Name: "city", Type: "text", Repeated: false},
    {Name: "country", Type: "text", Repeated: false},
}

employeeSchema := pg.Schema{
    {Name: "name", Type: "text", Repeated: false},
    {Name: "address", Type: "composite", Schema: addressSchema},
}

// CEL with nested field access
cel: employees.filter(e, e.address.city == 'New York')
sql: ARRAY(SELECT e FROM UNNEST(employees) AS e WHERE e.address.city = 'New York')

cel: employees.map(e, e.address.city)
sql: ARRAY(SELECT e.address.city FROM UNNEST(employees) AS e)
```

### Performance Considerations

- **UNNEST with large arrays**: PostgreSQL's `UNNEST()` function is efficient but consider indexing strategies for large datasets
- **Nested comprehensions**: May generate complex SQL; consider restructuring data or using materialized views for frequently accessed patterns
- **Map operations**: Return new arrays which may use memory; consider streaming for large results

### Usage in Practice

```go
package main

import (
    "fmt"
    "github.com/google/cel-go/cel"
    "github.com/spandigital/cel2sql/v2"
    "github.com/spandigital/cel2sql/v2/pg"
)

func main() {
    // Define PostgreSQL schema
    schema := pg.Schema{
        {Name: "id", Type: "bigint", Repeated: false},
        {Name: "name", Type: "text", Repeated: false},
        {Name: "skills", Type: "text", Repeated: true},
        {Name: "scores", Type: "bigint", Repeated: true},
    }

    provider := pg.NewTypeProvider(map[string]pg.Schema{"Employee": schema})

    env, _ := cel.NewEnv(
        cel.CustomTypeProvider(provider),
        cel.Variable("employees", cel.ListType(cel.ObjectType("Employee"))),
    )

    // Compile and convert CEL comprehension to SQL
    ast, _ := env.Compile(`employees.filter(e, e.scores.all(s, s >= 80)).map(e, e.name)`)
    sqlCondition, _ := cel2sql.Convert(ast)
    
    fmt.Println(sqlCondition)
    // Output: ARRAY(SELECT e.name FROM UNNEST(ARRAY(SELECT e FROM UNNEST(employees) AS e WHERE NOT EXISTS (SELECT 1 FROM UNNEST(e.scores) AS s WHERE NOT (s >= 80)))) AS e)
}
```
