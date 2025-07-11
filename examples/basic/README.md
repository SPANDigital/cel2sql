# Basic cel2sql Example

This example demonstrates the basic usage of cel2sql to convert CEL (Common Expression Language) expressions to PostgreSQL SQL conditions.

## What This Example Shows

- **Schema Definition**: How to define a PostgreSQL table schema using `pg.Schema`
- **CEL Environment Setup**: Creating a CEL environment with PostgreSQL type provider
- **Expression Compilation**: Compiling CEL expressions with type checking
- **SQL Conversion**: Converting compiled CEL expressions to PostgreSQL SQL

## Running the Example

```bash
go run main.go
```

## Expected Output

```
CEL Expression: employee.name == "John Doe" && employee.age >= 25 && employee.active
PostgreSQL SQL: employee.name = 'John Doe' AND employee.age >= 25 AND employee.active
```

## Code Breakdown

### 1. Schema Definition
```go
employeeSchema := pg.Schema{
    {Name: "name", Type: "text", Repeated: false},
    {Name: "age", Type: "integer", Repeated: false},
    {Name: "department", Type: "text", Repeated: false},
    {Name: "hired_at", Type: "timestamp with time zone", Repeated: false},
    {Name: "active", Type: "boolean", Repeated: false},
}
```

### 2. CEL Environment Setup
```go
env, err := cel.NewEnv(
    cel.CustomTypeProvider(pg.NewTypeProvider(map[string]pg.Schema{
        "Employee": employeeSchema,
    })),
    cel.Variable("employee", cel.ObjectType("Employee")),
)
```

### 3. Expression Compilation and Conversion
```go
ast, issues := env.Compile(`employee.name == "John Doe" && employee.age >= 25 && employee.active`)
sqlCondition, err := cel2sql.Convert(ast)
```

## Supported PostgreSQL Types

- `text` - String values
- `integer` - Integer numbers
- `boolean` - Boolean values
- `timestamp with time zone` - Timestamps
- Arrays (with `Repeated: true`)
- And many more...

## Next Steps

- Try modifying the CEL expression to test different operators
- Add more fields to the schema
- Experiment with array fields by setting `Repeated: true`
- Check out the [LoadTableSchema example](../load_table_schema/) for dynamic schema loading
