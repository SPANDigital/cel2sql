# CEL String Extension Functions Example

This example demonstrates the CEL string extension functions (`split()`, `join()`, and `format()`) that are now supported in cel2sql v3.4.0.

## String Extension Functions

### split(delimiter [, limit])

Splits a string into an array using a delimiter.

```cel
// Basic split (unlimited)
"a,b,c".split(",")              // → STRING_TO_ARRAY('a,b,c', ',')

// Split with limit
"a,b,c,d".split(",", 2)         // → (STRING_TO_ARRAY('a,b,c,d', ','))[1:2]

// Special cases
"text".split(",", 0)            // → ARRAY[]::text[] (empty array)
"text".split(",", 1)            // → ARRAY['text'] (no split)
```

**Parameters:**
- `delimiter`: String delimiter to split on
- `limit` (optional): Maximum number of splits
  - `-1` or omitted: Unlimited splits (default)
  - `0`: Return empty array
  - `1`: Return original string in array (no split)
  - `> 1`: Return first N elements

**Security:** Delimiter cannot contain null bytes.

### join([delimiter])

Joins an array into a string using a delimiter.

```cel
// Join with delimiter
["a", "b", "c"].join(",")       // → ARRAY_TO_STRING(ARRAY['a', 'b', 'c'], ',', '')

// Join without delimiter (empty string)
["a", "b", "c"].join()          // → ARRAY_TO_STRING(ARRAY['a', 'b', 'c'], '', '')
```

**Parameters:**
- `delimiter` (optional): String to join elements with (default: empty string)

**Null Handling:** Null elements are replaced with empty strings.

**Security:** Delimiter cannot contain null bytes.

### format(args)

Formats a string using printf-style placeholders.

```cel
// Basic formatting
"Hello %s".format(["World"])            // → FORMAT('Hello %s', 'World')

// Multiple arguments
"%s is %d years old".format(["John", 30]) // → FORMAT('%s is %s years old', 'John', 30)

// Escaped percent sign
"100%% complete".format([])             // → FORMAT('100%% complete')
```

**Supported Specifiers:**
- `%s`: String (stays as `%s` in PostgreSQL)
- `%d`: Decimal/integer (converted to `%s`)
- `%f`: Float (converted to `%s`)
- `%%`: Escaped percent sign

**Unsupported Specifiers:**
- `%b`, `%x`, `%X`, `%o`, `%e`, `%E`, `%g`, `%G`

**Requirements:**
- Format string must be a constant (not dynamic)
- Arguments must be a constant list
- Argument count must match placeholder count

**Security:**
- Format string limited to 1000 characters
- Only whitelisted specifiers allowed

## Running the Example

```bash
go run main.go
```

## Example Output

```
CEL String Extension Functions Examples
==========================================

1. split() - Basic String Splitting
   CEL: person.csv_data.split(',').size() > 0
   SQL: COALESCE(ARRAY_LENGTH(STRING_TO_ARRAY(person.csv_data, ','), 1), 0) > 0
   ✓ Success

2. split() with Limit
   CEL: person.csv_data.split(',', 3).size() == 3
   SQL: COALESCE(ARRAY_LENGTH((STRING_TO_ARRAY(person.csv_data, ','))[1:3], 1), 0) = 3
   ✓ Success

...
```

## Use Cases

### 1. CSV Parsing and Filtering

```cel
// Split CSV, filter non-empty values, rejoin
person.csv_data.split(',').filter(x, x.size() > 0).join('|')
```

### 2. Tag Management

```cel
// Find specific tag in comma-separated list
person.csv_data.split(',').exists(x, x == 'premium')

// Join tags with custom delimiter
person.tags.join(' | ')
```

### 3. Dynamic Message Generation

```cel
// Format user-friendly messages
"User %s (age %d) registered with email %s".format([person.name, person.age, person.email])
```

### 4. Data Transformation

```cel
// Split, transform, and rejoin
person.csv_data.split(',').map(x, x.upperAscii()).join(',')
```

## Integration with Comprehensions

All string extension functions work seamlessly with CEL comprehensions:

```cel
// split() in exists
person.csv_data.split(',').exists(x, x == 'target')

// split() in filter
person.csv_data.split(',').filter(x, x.startsWith('a')).size() > 0

// split() in map
person.csv_data.split(',').map(x, x.upperAscii())

// join() with filtered array
person.tags.filter(t, t.startsWith('a')).join(',')
```

## Related Documentation

- [CEL String Extensions](https://github.com/google/cel-go/tree/master/ext#strings)
- [PostgreSQL String Functions](https://www.postgresql.org/docs/current/functions-string.html)
- [PostgreSQL Array Functions](https://www.postgresql.org/docs/current/functions-array.html)
- [cel2sql Operators Reference](../../docs/operators-reference.md)

## Limitations

### split()
- Dynamic limits not supported (must be constant)
- Negative limits other than `-1` not supported

### join()
- Joins all non-null elements (nulls replaced with empty string)

### format()
- Format string must be constant (not dynamic)
- Arguments must be constant list (not dynamic)
- Limited specifier support (`%s`, `%d`, `%f` only)
- No width or precision modifiers

### quote()
- Not available (not part of CEL `ext.Strings()` standard extension)

## Security Considerations

All string extension functions include security validations:

- **Null Byte Protection**: Delimiters cannot contain null bytes
- **Format String Validation**:
  - Maximum length: 1000 characters
  - Only whitelisted specifiers
  - Argument count validation
- **Input Validation**: All field names validated for SQL injection

## Version History

- **v3.4.0**: Added `split()`, `join()`, and `format()` support
