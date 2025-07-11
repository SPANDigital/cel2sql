# cel2sql

cel2sql converts [CEL (Common Expression Language)](https://opensource.google/projects/cel) to SQL condition.
It is specifically targeting PostgreSQL standard SQL.

## Latest Release - v2.4.0

🚀 **New Features:**
- **Comprehensive JSON/JSONB Support**: Native PostgreSQL JSON path operations
- **Dynamic Schema Loading**: Load table schemas directly from PostgreSQL databases
- **Enhanced Type System**: Improved PostgreSQL type mappings and array support
- **Testcontainer Integration**: Full test coverage with real PostgreSQL databases

**Key Improvements:**
- PostgreSQL-optimized SQL generation (single quotes, proper functions)
- JSON field access: `user.preferences.theme` → `user.preferences->>'theme'`
- Array operations: `size(array)` → `ARRAY_LENGTH(array, 1)`
- String operations: `contains()` → `POSITION(...) > 0`
- All tests pass with comprehensive integration coverage

## Usage

```go
import (
    "context"
    "fmt"
    
    "github.com/spandigital/cel2sql"
    "github.com/spandigital/cel2sql/pg"
    "github.com/spandigital/cel2sql/sqltypes"
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
    
    "github.com/spandigital/cel2sql"
    "github.com/spandigital/cel2sql/pg"
    "github.com/spandigital/cel2sql/sqltypes"
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
      !_
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
      -_
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
      _!=_
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
      _%_
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
      _&&_
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
      _*_
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
      _+_
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
      _-_
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
      _/_
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
      _<=_
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
      _<_
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
      _==_
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
      _>=_
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
      _>_
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
      _?_:_
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
      _[_]
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
      _||_
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
      <code>REGEXP_LIKE(</code>string<code>, </code>string<code>)</code>
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

## Acknowledgments

This project is based on the excellent work by [@cockscomb](https://github.com/cockscomb) at [cockscomb/cel2sql](https://github.com/cockscomb/cel2sql). We thank the original author for creating the foundation that made this PostgreSQL-focused implementation possible.
