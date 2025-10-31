# Operators and Functions Reference

Complete reference for all CEL operators and functions supported by cel2sql.

## Table of Contents

- [Comparison Operators](#comparison-operators)
- [Logical Operators](#logical-operators)
- [Arithmetic Operators](#arithmetic-operators)
- [String Functions](#string-functions)
- [Type Conversion Functions](#type-conversion-functions)
- [Date/Time Functions](#datetime-functions)
- [Array Functions](#array-functions)
- [JSON Functions](#json-functions)
- [Other Functions](#other-functions)

## Comparison Operators

| CEL Operator | Description | Example CEL | Generated SQL |
|--------------|-------------|-------------|---------------|
| `==` | Equal | `age == 25` | `age = 25` |
| `!=` | Not equal | `status != "active"` | `status != 'active'` |
| `<` | Less than | `price < 100` | `price < 100` |
| `<=` | Less than or equal | `age <= 65` | `age <= 65` |
| `>` | Greater than | `score > 80` | `score > 80` |
| `>=` | Greater than or equal | `quantity >= 10` | `quantity >= 10` |

### Special Comparisons

| CEL Expression | Description | Generated SQL |
|----------------|-------------|---------------|
| `field == null` | Is null | `field IS NULL` |
| `field != null` | Is not null | `field IS NOT NULL` |
| `flag == true` | Boolean true | `flag IS TRUE` |
| `flag == false` | Boolean false | `flag IS FALSE` |

## Logical Operators

| CEL Operator | Description | Example CEL | Generated SQL |
|--------------|-------------|-------------|---------------|
| `&&` | Logical AND | `active && verified` | `active IS TRUE AND verified IS TRUE` |
| `\|\|` | Logical OR | `admin \|\| moderator` | `admin IS TRUE OR moderator IS TRUE` |
| `!` | Logical NOT | `!deleted` | `NOT deleted IS TRUE` |

### Examples

```go
// AND
age >= 18 && age <= 65
// SQL: age >= 18 AND age <= 65

// OR
status == "active" || status == "pending"
// SQL: status = 'active' OR status = 'pending'

// NOT
!archived
// SQL: NOT archived IS TRUE

// Combined
(age >= 18 && age <= 65) || vip
// SQL: (age >= 18 AND age <= 65) OR vip IS TRUE
```

## Arithmetic Operators

| CEL Operator | Types | Example CEL | Generated SQL |
|--------------|-------|-------------|---------------|
| `+` | int, double | `price + tax` | `price + tax` |
| `-` | int, double | `total - discount` | `total - discount` |
| `*` | int, double | `quantity * price` | `quantity * price` |
| `/` | int, double | `total / count` | `total / count` |
| `%` | int | `value % 10` | `MOD(value, 10)` |
| `-` (unary) | int, double | `-balance` | `-balance` |

## String Functions

| CEL Function | Description | Example CEL | Generated SQL |
|--------------|-------------|-------------|---------------|
| `startsWith()` | String starts with prefix | `name.startsWith("A")` | `name LIKE 'A%'` |
| `endsWith()` | String ends with suffix | `email.endsWith(".com")` | `email LIKE '%.com'` |
| `contains()` | String contains substring | `text.contains("search")` | `POSITION('search' IN text) > 0` |
| `matches()` | Regex pattern match | `email.matches(r".*@test\.com")` | `email ~ '.*@test\.com'` |
| `+` | String concatenation | `first + " " + last` | `first \|\| ' ' \|\| last` |
| `size()` | String length | `size(name)` | `CHAR_LENGTH(name)` |
| `split()` | Split string to array | `"a,b,c".split(",")` | `STRING_TO_ARRAY('a,b,c', ',')` |
| `split(limit)` | Split with limit | `text.split(",", 3)` | `(STRING_TO_ARRAY(text, ','))[1:3]` |
| `join()` | Join array to string | `tags.join(",")` | `ARRAY_TO_STRING(tags, ',', '')` |
| `format()` | Format string | `"%s: %d".format([name, age])` | `FORMAT('%s: %s', name, age)` |

### Examples

```go
// Starts with
name.startsWith("John")
// SQL: name LIKE 'John%'

// Ends with
filename.endsWith(".pdf")
// SQL: filename LIKE '%.pdf'

// Contains
description.contains("urgent")
// SQL: POSITION('urgent' IN description) > 0

// Pattern matching
email.matches(r"^[a-z]+@example\.com$")
// SQL: email ~ '^[a-z]+@example\.com$'

// Concatenation
first_name + " " + last_name
// SQL: first_name || ' ' || last_name

// Split string to array
"a,b,c".split(",")
// SQL: STRING_TO_ARRAY('a,b,c', ',')

// Split with limit
text.split(",", 3)
// SQL: (STRING_TO_ARRAY(text, ','))[1:3]

// Join array to string
tags.join(",")
// SQL: ARRAY_TO_STRING(tags, ',', '')

// Format string
"Name: %s, Age: %d".format([person.name, person.age])
// SQL: FORMAT('Name: %s, Age: %s', person.name, person.age)
```

### String Extension Functions (v3.4.0+)

cel2sql v3.4.0 adds support for three CEL `ext.Strings()` functions:

#### split(delimiter [, limit])

Splits a string into an array using a delimiter.

```go
// Basic split (unlimited)
"a,b,c".split(",")
// SQL: STRING_TO_ARRAY('a,b,c', ',')

// Split with limit
"a,b,c,d".split(",", 2)
// SQL: (STRING_TO_ARRAY('a,b,c,d', ','))[1:2]

// Special cases
"text".split(",", 0)  // Returns: ARRAY[]::text[] (empty array)
"text".split(",", 1)  // Returns: ARRAY['text'] (no split)
"text".split(",", -1) // Returns: STRING_TO_ARRAY('text', ',') (unlimited, default)
```

**Supported in comprehensions:**
```go
person.csv.split(',').exists(x, x == 'target')
// SQL: EXISTS (SELECT 1 FROM UNNEST(STRING_TO_ARRAY(person.csv, ',')) AS x WHERE x = 'target')
```

**Limitations:**
- Dynamic limits not supported (must be constant)
- Negative limits other than -1 not supported

#### join([delimiter])

Joins an array into a string using a delimiter.

```go
// Join with delimiter
["a", "b", "c"].join(",")
// SQL: ARRAY_TO_STRING(ARRAY['a', 'b', 'c'], ',', '')

// Join without delimiter (empty string)
["a", "b", "c"].join()
// SQL: ARRAY_TO_STRING(ARRAY['a', 'b', 'c'], '', '')

// Join array field
person.tags.join(", ")
// SQL: ARRAY_TO_STRING(person.tags, ', ', '')
```

**Supported in comprehensions:**
```go
person.tags.filter(t, t.startsWith('a')).join(',')
// SQL: ARRAY_TO_STRING(ARRAY(SELECT t FROM UNNEST(person.tags) AS t WHERE t LIKE 'a%'), ',', '')
```

**Note:** Null elements are replaced with empty strings.

#### format(args)

Formats a string using printf-style placeholders.

```go
// Basic formatting
"Hello %s".format(["World"])
// SQL: FORMAT('Hello %s', 'World')

// Multiple arguments
"%s is %d years old".format(["John", 30])
// SQL: FORMAT('%s is %s years old', 'John', 30)

// With field values
"User: %s, Email: %s".format([person.name, person.email])
// SQL: FORMAT('User: %s, Email: %s', person.name, person.email)
```

**Supported specifiers:**
- `%s`: String (stays as %s)
- `%d`: Decimal/integer (converted to %s)
- `%f`: Float (converted to %s)
- `%%`: Escaped percent sign

**Unsupported specifiers:** `%b`, `%x`, `%X`, `%o`, `%e`, `%E`, `%g`, `%G`

**Limitations:**
- Format string must be constant (not dynamic)
- Arguments must be constant list
- Format string max length: 1000 characters
- Argument count must match placeholder count

See `examples/string_extensions/` for complete working examples.

## Type Conversion Functions

| CEL Function | Description | Example CEL | Generated SQL |
|--------------|-------------|-------------|---------------|
| `bool()` | Convert to boolean | `bool(active)` | `CAST(active AS BOOLEAN)` |
| `int()` | Convert to integer | `int(score)` | `CAST(score AS BIGINT)` |
| `double()` | Convert to double | `double(price)` | `CAST(price AS DOUBLE PRECISION)` |
| `string()` | Convert to string | `string(age)` | `CAST(age AS TEXT)` |
| `bytes()` | Convert to bytes | `bytes(data)` | `CAST(data AS BYTEA)` |

### Special: Timestamp to Unix

```go
// Convert timestamp to Unix epoch
int(created_at)
// SQL: EXTRACT(EPOCH FROM created_at)::bigint
```

## Date/Time Functions

### Timestamp Functions

| CEL Function | Description | Example CEL | Generated SQL |
|--------------|-------------|-------------|---------------|
| `timestamp()` | Parse timestamp | `timestamp("2024-01-01T00:00:00Z")` | `CAST('2024-01-01T00:00:00Z' AS TIMESTAMP WITH TIME ZONE)` |
| `getFullYear()` | Get year | `created_at.getFullYear()` | `EXTRACT(YEAR FROM created_at)` |
| `getMonth()` | Get month (0-based) | `created_at.getMonth()` | `EXTRACT(MONTH FROM created_at) - 1` |
| `getDayOfMonth()` | Get day (0-based) | `created_at.getDayOfMonth()` | `EXTRACT(DAY FROM created_at) - 1` |
| `getDayOfWeek()` | Get weekday (0=Monday) | `created_at.getDayOfWeek()` | `(EXTRACT(DOW FROM created_at) + 6) % 7` |
| `getDayOfYear()` | Get day of year (0-based) | `created_at.getDayOfYear()` | `EXTRACT(DOY FROM created_at) - 1` |
| `getHours()` | Get hours | `created_at.getHours()` | `EXTRACT(HOUR FROM created_at)` |
| `getMinutes()` | Get minutes | `created_at.getMinutes()` | `EXTRACT(MINUTE FROM created_at)` |
| `getSeconds()` | Get seconds | `created_at.getSeconds()` | `EXTRACT(SECOND FROM created_at)` |

### Duration Functions

| CEL Function | Description | Example CEL | Generated SQL |
|--------------|-------------|-------------|---------------|
| `duration()` | Create duration | `duration("24h")` | `INTERVAL 24 HOUR` |
| `+` | Add duration | `timestamp + duration("1h")` | `timestamp + INTERVAL 1 HOUR` |
| `-` | Subtract duration | `timestamp - duration("1d")` | `timestamp - INTERVAL 1 DAY` |

### Current Time Functions

| CEL Function | Description | Generated SQL |
|--------------|-------------|---------------|
| `current_date()` | Current date | `CURRENT_DATE` |
| `current_time()` | Current time | `CURRENT_TIME` |
| `current_datetime()` | Current datetime | `CURRENT_TIMESTAMP` |
| `current_timestamp()` | Current timestamp | `CURRENT_TIMESTAMP` |

### Examples

```go
// Parse timestamp
timestamp("2024-01-01T00:00:00Z")
// SQL: CAST('2024-01-01T00:00:00Z' AS TIMESTAMP WITH TIME ZONE)

// Get year
created_at.getFullYear() == 2024
// SQL: EXTRACT(YEAR FROM created_at) = 2024

// Add duration
created_at + duration("7d")
// SQL: created_at + INTERVAL 7 DAY

// Recent records
created_at > current_timestamp() - duration("24h")
// SQL: created_at > CURRENT_TIMESTAMP - INTERVAL 24 HOUR
```

## Array Functions

| CEL Function | Description | Example CEL | Generated SQL |
|--------------|-------------|-------------|---------------|
| `size()` | Array length | `size(tags)` | `ARRAY_LENGTH(tags, 1)` |
| `in` | Element in array | `"admin" in roles` | `'admin' IN UNNEST(roles)` |
| `all()` | All elements match | `scores.all(s, s > 60)` | `NOT EXISTS (SELECT 1 FROM UNNEST(scores) AS s WHERE NOT (s > 60))` |
| `exists()` | Any element matches | `tags.exists(t, t == "urgent")` | `EXISTS (SELECT 1 FROM UNNEST(tags) AS t WHERE t = 'urgent')` |
| `exists_one()` | Exactly one matches | `items.exists_one(i, i.id == 5)` | `(SELECT COUNT(*) FROM UNNEST(items) AS i WHERE i.id = 5) = 1` |
| `filter()` | Filter elements | `scores.filter(s, s >= 80)` | `ARRAY(SELECT s FROM UNNEST(scores) AS s WHERE s >= 80)` |
| `map()` | Transform elements | `names.map(n, n + "!")` | `ARRAY(SELECT n \|\| '!' FROM UNNEST(names) AS n)` |

### Examples

```go
// Array length
size(tags) > 0
// SQL: ARRAY_LENGTH(tags, 1) > 0

// Element in array
"premium" in features
// SQL: 'premium' IN UNNEST(features)

// All elements match
prices.all(p, p > 0)
// SQL: NOT EXISTS (SELECT 1 FROM UNNEST(prices) AS p WHERE NOT (p > 0))

// Filter array
scores.filter(s, s >= 80)
// SQL: ARRAY(SELECT s FROM UNNEST(scores) AS s WHERE s >= 80)
```

## JSON Functions

| CEL Expression | Description | Example CEL | Generated SQL |
|----------------|-------------|-------------|---------------|
| `field.key` | Access JSON field | `data.name` | `data->>'name'` |
| `field.nested.key` | Nested JSON access | `user.profile.email` | `user.profile->'profile'->>'email'` |
| `has()` | Check field exists | `has(data.key)` | `data ? 'key'` |

### Examples

```go
// Simple JSON access
preferences.theme == "dark"
// SQL: preferences->>'theme' = 'dark'

// Nested JSON
user.profile.settings.language == "en"
// SQL: user.profile->'settings'->>'language' = 'en'

// Field existence
has(preferences.theme)
// SQL: preferences ? 'theme'

// Numeric JSON field
int(metadata.count) > 10
// SQL: (metadata->>'count')::bigint > 10
```

## Other Functions

### Conditional (Ternary)

| CEL Expression | Description | Generated SQL |
|----------------|-------------|---------------|
| `condition ? true_val : false_val` | If-then-else | `CASE WHEN condition THEN true_val ELSE false_val END` |

```go
// Ternary operator
age >= 18 ? "adult" : "minor"
// SQL: CASE WHEN age >= 18 THEN 'adult' ELSE 'minor' END
```

### List Construction

```go
// List literal
[1, 2, 3, 4, 5]
// SQL: ARRAY[1, 2, 3, 4, 5]

// String list
["a", "b", "c"]
// SQL: ARRAY['a', 'b', 'c']
```

## Type Mapping Reference

| CEL Type | PostgreSQL Type | Notes |
|----------|-----------------|-------|
| `int` | `bigint` | 64-bit integers |
| `uint` | `bigint` | Treated as signed |
| `double` | `double precision` | Floating point |
| `bool` | `boolean` | True/false |
| `string` | `text` | Variable-length text |
| `bytes` | `bytea` | Binary data |
| `list<T>` | `T[]` | Arrays |
| `map<K,V>` | `jsonb` | For complex objects |
| `timestamp` | `timestamp with time zone` | Timestamps |
| `duration` | `interval` | Time intervals |
| `null_type` | `NULL` | Null values |

## Security Considerations

cel2sql includes built-in security protections that are automatically applied to all operations:

### Field Name Validation

All field names used in operators are validated to prevent SQL injection:

```go
// ✅ Safe: Valid field names
user.name == "John"
product.price > 100

// ❌ Blocked: Malicious field names
field'; DROP TABLE users-- == "value"
// Error: invalid field name
```

**Validation rules:**
- Maximum 63 characters (PostgreSQL limit)
- Must start with letter or underscore
- Only alphanumeric and underscore characters
- SQL reserved keywords blocked

### String Escaping

All string literals in comparisons are properly escaped:

```go
// Strings with quotes are safely escaped
name == "O'Brien"
// SQL: name = 'O''Brien'

// Prevents injection via string values
field == "'; DROP TABLE users--"
// SQL: field = '''; DROP TABLE users--'
```

### JSON Field Security

JSON field names are automatically escaped:

```go
// Quotes in JSON field names are escaped
user.preferences.theme'name == "dark"
// SQL: user.preferences->>'theme''name' = 'dark'
```

See [JSON/JSONB Support](json-support.md) for more details.

### Regex Pattern Validation

Regex patterns are validated to prevent ReDoS attacks:

```go
// ✅ Safe patterns allowed
email.matches(r"^[a-z]+@[a-z]+\.[a-z]+$")

// ❌ Dangerous patterns blocked
field.matches(r"(a+)+")  // Nested quantifiers
// Error: nested quantifiers detected
```

**Pattern limits:**
- Maximum 500 characters
- No nested quantifiers
- Maximum 20 capture groups
- Maximum 10 nesting levels

See [Regex Matching](regex-matching.md) for more details.

### Operator Safety

All operators are converted using safe, parameterized patterns:

| Security Feature | Protection |
|------------------|------------|
| **Comparison operators** | Type-safe conversions |
| **Logical operators** | Proper boolean handling |
| **String operations** | Escape special characters |
| **List operations** | Array boundary checks |
| **Math operations** | Overflow protection via PostgreSQL |

### Best Practices

1. **Use context timeouts** for user-provided expressions:
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
sql, err := cel2sql.Convert(ast, cel2sql.WithContext(ctx))
```

2. **Validate expression complexity** before processing:
```go
if len(celExpression) > 1000 {
    return errors.New("expression too complex")
}
```

3. **Use prepared statements** when executing generated SQL:
```go
stmt, err := db.Prepare("SELECT * FROM table WHERE " + sqlCondition)
```

4. **Enable logging** to monitor patterns:
```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
sql, err := cel2sql.Convert(ast, cel2sql.WithLogger(logger))
```

5. **Keep schemas minimal** - only expose necessary fields:
```go
// ✅ Good: Only expose needed fields
schema := pg.Schema{
    {Name: "id", Type: "bigint"},
    {Name: "name", Type: "text"},
    {Name: "email", Type: "text"},
}

// ❌ Avoid: Exposing sensitive fields
// Don't expose: password_hash, ssn, credit_card, etc.
```

For comprehensive security information, see the [Security Guide](security.md).

## See Also

- [Getting Started Guide](getting-started.md)
- [JSON/JSONB Support](json-support.md)
- [Array Comprehensions](comprehensions.md)
- [Regex Matching](regex-matching.md)
