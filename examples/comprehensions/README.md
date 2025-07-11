# CEL Comprehensions Example

This example demonstrates how to use CEL comprehensions with cel2sql to generate PostgreSQL-compatible SQL.

## Features Demonstrated

- **Simple Array Comprehensions**: Basic `all()`, `exists()`, `exists_one()`, `filter()`, and `map()` operations
- **PostgreSQL Schema Integration**: Using comprehensions with `pg.Schema` and `pg.TypeProvider`
- **Nested Comprehensions**: Complex nested operations like `employees.exists(e, e.skills.exists(s, s == 'Go'))`
- **Composite Types**: Working with nested PostgreSQL composite types
- **Array Fields**: Direct comprehensions on array fields in schemas

## Supported Comprehension Types

| CEL Expression | Description | Generated SQL Pattern |
|----------------|-------------|----------------------|
| `list.all(x, condition)` | All elements satisfy condition | `NOT EXISTS (SELECT 1 FROM UNNEST(list) AS x WHERE NOT (condition))` |
| `list.exists(x, condition)` | At least one element satisfies condition | `EXISTS (SELECT 1 FROM UNNEST(list) AS x WHERE condition)` |
| `list.exists_one(x, condition)` | Exactly one element satisfies condition | `(SELECT COUNT(*) FROM UNNEST(list) AS x WHERE condition) = 1` |
| `list.filter(x, condition)` | Return elements that satisfy condition | `ARRAY(SELECT x FROM UNNEST(list) AS x WHERE condition)` |
| `list.map(x, transform)` | Transform all elements | `ARRAY(SELECT transform FROM UNNEST(list) AS x)` |

## Running the Example

```bash
cd examples/comprehensions
go run main.go
```

## Expected Output

The example will show three categories of comprehensions:

1. **Simple Array Comprehensions** - Working with literal arrays
2. **PostgreSQL Schema-based Comprehensions** - Using employee schemas with array fields
3. **Nested and Complex Comprehensions** - Advanced patterns with composite types

Each comprehension is converted to efficient PostgreSQL SQL using `UNNEST()` and array functions.

## Key Concepts

### PostgreSQL UNNEST Integration

All comprehensions use PostgreSQL's `UNNEST()` function to efficiently process arrays:

```sql
-- CEL: [1,2,3,4,5].all(x, x > 0)
-- SQL: NOT EXISTS (SELECT 1 FROM UNNEST(ARRAY[1,2,3,4,5]) AS x WHERE NOT (x > 0))
```

### Schema-based Operations

When using PostgreSQL schemas, comprehensions work directly with table columns:

```sql
-- CEL: employees.filter(e, e.active)
-- SQL: ARRAY(SELECT e FROM UNNEST(employees) AS e WHERE e.active)
```

### Nested Comprehensions

Complex nested operations are supported for advanced filtering:

```sql
-- CEL: employees.exists(e, e.skills.exists(s, s == 'Go'))
-- SQL: EXISTS (SELECT 1 FROM UNNEST(employees) AS e WHERE EXISTS (SELECT 1 FROM UNNEST(e.skills) AS s WHERE s = 'Go'))
```

## Performance Notes

- **UNNEST efficiency**: PostgreSQL's `UNNEST()` is optimized for array processing
- **Index considerations**: Consider indexing strategies for large arrays
- **Memory usage**: Map operations create new arrays; consider streaming for large datasets
