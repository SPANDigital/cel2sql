# Regex Pattern Matching

cel2sql supports CEL's `matches()` function with automatic conversion from RE2 (CEL's regex flavor) to PostgreSQL POSIX regex.

## Quick Start

```go
// Define schema
schema := pg.Schema{
    {Name: "email", Type: "text"},
    {Name: "phone", Type: "text"},
}

// CEL: Email validation
email.matches(r"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$")

// Generated SQL:
email ~ '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'
```

## Basic Syntax

### Method Style

```go
// CEL: Call matches() on string
field.matches(r"pattern")

// Generated SQL:
field ~ 'pattern'
```

### Function Style

```go
// CEL: Use matches() as function
matches(field, r"pattern")

// Generated SQL:
field ~ 'pattern'
```

## PostgreSQL Operators

| Operator | Description | CEL Example |
|----------|-------------|-------------|
| `~` | Case-sensitive match | `field.matches(r"pattern")` |
| `~*` | Case-insensitive match | `field.matches(r"(?i)pattern")` |
| `!~` | Case-sensitive non-match | `!field.matches(r"pattern")` |
| `!~*` | Case-insensitive non-match | `!field.matches(r"(?i)pattern")` |

## Common Patterns

### Email Validation

```go
// Simple email pattern
email.matches(r"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$")
// SQL: email ~ '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'

// Email from specific domain
email.matches(r"^.*@company\.com$")
// SQL: email ~ '^.*@company\.com$'
```

### Phone Numbers

```go
// US phone number (123-456-7890)
phone.matches(r"^\d{3}-\d{3}-\d{4}$")
// SQL: phone ~ '^[[:digit:]]{3}-[[:digit:]]{3}-[[:digit:]]{4}$'

// International format (+1-234-567-8900)
phone.matches(r"^\+\d{1,3}-\d{3}-\d{3}-\d{4}$")
// SQL: phone ~ '^\+[[:digit:]]{1,3}-[[:digit:]]{3}-[[:digit:]]{3}-[[:digit:]]{4}$'
```

### URLs

```go
// HTTP/HTTPS URLs
url.matches(r"^https?://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}")
// SQL: url ~ '^https?://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}'

// Specific domain
url.matches(r"^https://.*\.example\.com")
// SQL: url ~ '^https://.*\.example\.com'
```

### Usernames

```go
// Alphanumeric, 3-20 characters
username.matches(r"^[a-zA-Z0-9]{3,20}$")
// SQL: username ~ '^[a-zA-Z0-9]{3,20}$'

// Allow underscores and hyphens
username.matches(r"^[a-zA-Z0-9_-]{3,20}$")
// SQL: username ~ '^[a-zA-Z0-9_-]{3,20}$'
```

### Zip Codes

```go
// US ZIP code (12345 or 12345-6789)
zipcode.matches(r"^\d{5}(-\d{4})?$")
// SQL: zipcode ~ '^[[:digit:]]{5}(-[[:digit:]]{4})?$'
```

## Case-Insensitive Matching

Use the `(?i)` flag for case-insensitive matching:

```go
// CEL: Case-insensitive search
description.matches(r"(?i)urgent|priority|important")

// Generated SQL:
description ~* 'urgent|priority|important'
```

## Character Classes

CEL uses RE2 character classes, which are automatically converted:

| RE2 Class | POSIX Class | Description |
|-----------|-------------|-------------|
| `\d` | `[[:digit:]]` | Digits (0-9) |
| `\w` | `[[:alnum:]_]` | Word characters (a-z, A-Z, 0-9, _) |
| `\s` | `[[:space:]]` | Whitespace |
| `\D` | `[^[:digit:]]` | Non-digits |
| `\W` | `[^[:alnum:]_]` | Non-word characters |
| `\S` | `[^[:space:]]` | Non-whitespace |

### Examples

```go
// Digits only
code.matches(r"^\d+$")
// SQL: code ~ '^[[:digit:]]+$'

// Word characters
identifier.matches(r"^\w+$")
// SQL: identifier ~ '^[[:alnum:]_]+$'

// Contains whitespace
text.matches(r".*\s+.*")
// SQL: text ~ '.*[[:space:]]+.*'
```

## Word Boundaries

```go
// Match whole word "test"
text.matches(r"\btest\b")
// SQL: text ~ '\ytest\y'

// Word starting with "pre"
text.matches(r"\bpre\w+")
// SQL: text ~ '\ypre[[:alnum:]_]+'
```

## Anchors

| Anchor | Description | Example |
|--------|-------------|---------|
| `^` | Start of string | `^abc` matches "abc..." |
| `$` | End of string | `xyz$` matches "...xyz" |
| `\b` or `\y` | Word boundary | `\bword\b` matches "word" |

```go
// Start of string
field.matches(r"^Hello")
// SQL: field ~ '^Hello'

// End of string
field.matches(r"\.com$")
// SQL: field ~ '\.com$'

// Exact match
field.matches(r"^exact$")
// SQL: field ~ '^exact$'
```

## Quantifiers

| Quantifier | Description | Example |
|------------|-------------|---------|
| `*` | 0 or more | `a*` |
| `+` | 1 or more | `a+` |
| `?` | 0 or 1 | `a?` |
| `{n}` | Exactly n | `a{3}` |
| `{n,}` | n or more | `a{3,}` |
| `{n,m}` | Between n and m | `a{3,5}` |

```go
// One or more digits
field.matches(r"^\d+$")
// SQL: field ~ '^[[:digit:]]+$'

// 3 to 5 letters
field.matches(r"^[a-z]{3,5}$")
// SQL: field ~ '^[a-z]{3,5}$'

// Optional dash
field.matches(r"^ABC-?\d{4}$")
// SQL: field ~ '^ABC-?[[:digit:]]{4}$'
```

## Groups and Alternation

```go
// Alternation (OR)
field.matches(r"^(cat|dog|bird)$")
// SQL: field ~ '^(cat|dog|bird)$'

// Grouped pattern
field.matches(r"^(Mr|Ms|Dr)\.?\s+[A-Z][a-z]+$")
// SQL: field ~ '^(Mr|Ms|Dr)\.?[[:space:]]+[A-Z][a-z]+$'
```

## Complete Example

```go
package main

import (
    "fmt"
    "log"

    "github.com/google/cel-go/cel"
    "github.com/spandigital/cel2sql/v3"
    "github.com/spandigital/cel2sql/v3/pg"
)

func main() {
    // Define schema
    contactSchema := pg.Schema{
        {Name: "id", Type: "bigint"},
        {Name: "email", Type: "text"},
        {Name: "phone", Type: "text"},
        {Name: "website", Type: "text"},
    }

    provider := pg.NewTypeProvider(map[string]pg.Schema{
        "Contact": contactSchema,
    })

    env, err := cel.NewEnv(
        cel.CustomTypeProvider(provider),
        cel.Variable("contact", cel.ObjectType("Contact")),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Example 1: Email validation
    ast, _ := env.Compile(`contact.email.matches(r"^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$")`)
    sql, _ := cel2sql.Convert(ast)
    fmt.Println("Email:", sql)
    // Output: contact.email ~ '^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$'

    // Example 2: Phone validation
    ast, _ = env.Compile(`contact.phone.matches(r"^\d{3}-\d{3}-\d{4}$")`)
    sql, _ = cel2sql.Convert(ast)
    fmt.Println("Phone:", sql)
    // Output: contact.phone ~ '^[[:digit:]]{3}-[[:digit:]]{3}-[[:digit:]]{4}$'

    // Example 3: Case-insensitive domain check
    ast, _ = env.Compile(`contact.website.matches(r"(?i).*\.com$")`)
    sql, _ = cel2sql.Convert(ast)
    fmt.Println("Website:", sql)
    // Output: contact.website ~* '.*\.com$'

    // Example 4: Combined conditions
    ast, _ = env.Compile(`
        contact.email.matches(r"@company\.com$") &&
        contact.phone.matches(r"^\+1-")
    `)
    sql, _ = cel2sql.Convert(ast)
    fmt.Println("Combined:", sql)
    // Output: contact.email ~ '@company\.com$' AND contact.phone ~ '^\+1-'
}
```

## Performance Tips

### 1. Use Anchors

```go
// ✅ Better: Anchored pattern
field.matches(r"^pattern$")

// ❌ Slower: Unanchored pattern
field.matches(r"pattern")
```

### 2. Keep Patterns Simple

```go
// ✅ Better: Simple pattern
field.matches(r"^[a-z]+$")

// ❌ Slower: Complex nested groups
field.matches(r"^((([a-z])+)+)+$")  // Catastrophic backtracking risk
```

### 3. Use Indexes When Possible

For frequently searched patterns, consider PostgreSQL indexes:

```sql
-- GIN index for trigram search (pg_trgm extension)
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX idx_field_trgm ON table USING GIN (field gin_trgm_ops);
```

## ReDoS Protection

cel2sql includes comprehensive protection against Regular Expression Denial of Service (ReDoS) attacks, also known as catastrophic backtracking (CWE-1333).

### What is ReDoS?

ReDoS attacks exploit poorly-constructed regex patterns that can cause exponential time complexity when matching certain inputs. A malicious pattern can consume 100% CPU and freeze your application.

**Example of dangerous pattern:**
```go
// ❌ DANGEROUS: Nested quantifiers cause catastrophic backtracking
field.matches(r"(a+)+b")
// Input "aaaaaaaaaaaaaaaaaaaX" can take HOURS to process
```

### Automatic Validation

All regex patterns are automatically validated before conversion. Dangerous patterns are rejected with descriptive error messages.

### Pattern Limits

#### 1. Length Limit (500 characters)

```go
// ✅ Allowed: Normal pattern
email.matches(r"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$")

// ❌ Rejected: Pattern too long
field.matches(r"very long pattern..." * 100)
// Error: pattern exceeds 500 character limit
```

#### 2. Nested Quantifiers Detection

Patterns with nested quantifiers are blocked:

```go
// ❌ Rejected: Nested + quantifiers
field.matches(r"(a+)+")
// Error: nested quantifiers detected at position X

// ❌ Rejected: Nested * quantifiers
field.matches(r"(a*)*")
// Error: nested quantifiers detected

// ❌ Rejected: Mixed nesting
field.matches(r"(a+)*b")
// Error: nested quantifiers detected

// ✅ Allowed: Non-nested quantifiers
field.matches(r"(abc)+")
field.matches(r"a+b+c+")
```

#### 3. Capture Group Limit (20 groups)

```go
// ✅ Allowed: Reasonable number of groups
field.matches(r"(a)(b)(c)(d)(e)")

// ❌ Rejected: Too many capture groups
field.matches(r"(a)(b)(c)...(z)(aa)(ab)")  // > 20 groups
// Error: pattern exceeds 20 capture group limit
```

#### 4. Quantified Alternation Detection

Patterns with quantified alternation are blocked:

```go
// ❌ Rejected: Quantified alternation
field.matches(r"(a|a)+b")
// Error: quantified alternation detected

// ❌ Rejected: Subtle alternation issue
field.matches(r"(x|xy)+y")
// Error: quantified alternation detected

// ✅ Allowed: Safe alternation
field.matches(r"(cat|dog|bird)")
field.matches(r"^(Mr|Ms|Dr)")
```

#### 5. Nesting Depth Limit (10 levels)

```go
// ✅ Allowed: Reasonable nesting
field.matches(r"((abc))")

// ❌ Rejected: Excessive nesting
field.matches(r"((((((((((a))))))))))")  // 10+ levels
// Error: pattern exceeds nesting depth limit of 10
```

### Safe Pattern Examples

These patterns are safe and allowed:

```go
// ✅ Email validation
email.matches(r"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$")

// ✅ Phone number
phone.matches(r"^\d{3}-\d{3}-\d{4}$")

// ✅ URL validation
url.matches(r"^https?://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}")

// ✅ Username (3-20 chars)
username.matches(r"^[a-zA-Z0-9_-]{3,20}$")

// ✅ Case-insensitive search
description.matches(r"(?i)urgent|priority|important")

// ✅ Multiple alternations
status.matches(r"^(active|inactive|pending|approved|rejected)$")
```

### Dangerous Pattern Examples

These patterns are blocked for security:

```go
// ❌ Nested quantifiers
field.matches(r"(a+)+b")
field.matches(r"(a*)*b")
field.matches(r"(a+)*b")
field.matches(r"(a?)+b")

// ❌ Quantified alternation
field.matches(r"(a|a)+")
field.matches(r"(x|xy)+y")
field.matches(r"(ab|a)+c")

// ❌ Excessive nesting
field.matches(r"((((((((((a))))))))))")

// ❌ Too many groups
field.matches(r"(a)(b)(c)...(z)(aa)(ab)(ac)")  // > 20 groups

// ❌ Pattern too long
field.matches(r"..." * 600)  // > 500 chars
```

### Error Messages

When a dangerous pattern is detected, you get a descriptive error:

```go
ast, _ := env.Compile(`field.matches(r"(a+)+")`)
sql, err := cel2sql.Convert(ast)
// err: "invalid regex pattern: nested quantifiers detected at position 4"

ast, _ = env.Compile(`field.matches(r"(a|a)*b")`)
sql, err = cel2sql.Convert(ast)
// err: "invalid regex pattern: quantified alternation detected"
```

### Testing Your Patterns

Before deploying regex patterns, test them:

```go
func TestRegexPattern(t *testing.T) {
    // Test pattern compiles
    ast, issues := env.Compile(`field.matches(r"^[a-z]+$")`)
    require.NoError(t, issues.Err())

    // Test pattern converts without error
    sql, err := cel2sql.Convert(ast)
    require.NoError(t, err)

    // Test against sample data
    // ... execute query with test data
}
```

### Additional Protection Layers

For defense-in-depth, combine ReDoS protection with:

#### Context Timeouts

```go
// Timeout conversion after 5 seconds
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

sql, err := cel2sql.Convert(ast, cel2sql.WithContext(ctx))
if errors.Is(err, context.DeadlineExceeded) {
    log.Println("Pattern conversion timed out")
}
```

#### Input Validation

```go
// Validate before compiling
func validateCELExpression(expr string) error {
    if len(expr) > 1000 {
        return errors.New("expression too long")
    }
    if strings.Count(expr, "(") > 20 {
        return errors.New("too many nested groups")
    }
    return nil
}
```

#### Logging

```go
// Monitor patterns in production
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
sql, err := cel2sql.Convert(ast,
    cel2sql.WithLogger(logger))
// Logs will show regex conversion details
```

For more security information, see the [Security Guide](security.md).

## Common Patterns Reference

### Validation Patterns

```go
// Email
r"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$"

// URL
r"^https?://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}"

// IPv4 Address
r"^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$"

// Credit Card (basic)
r"^\d{4}-?\d{4}-?\d{4}-?\d{4}$"

// Date (YYYY-MM-DD)
r"^\d{4}-\d{2}-\d{2}$"

// Time (HH:MM:SS)
r"^\d{2}:\d{2}:\d{2}$"

// Hex Color
r"^#[0-9A-Fa-f]{6}$"

// UUID
r"^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"
```

## Escaping Special Characters

Remember to escape regex special characters:

```go
// Escape dots for literal match
field.matches(r"example\.com")  // Literal dot
field.matches(r"example.com")   // Any character

// Common characters to escape: . * + ? ^ $ { } ( ) [ ] | \
```

## Error Handling

```go
// Invalid regex pattern
field.matches(r"[invalid")
// Error during compilation: invalid regex pattern

// Type mismatch
number_field.matches(r"pattern")
// Error: matches() requires string type

// Missing raw string prefix
field.matches("pattern")  // Works but CEL may interpret escapes
field.matches(r"pattern") // Better: raw string
```

## See Also

- [Getting Started Guide](getting-started.md)
- [Operators Reference](operators-reference.md)
- [PostgreSQL Regex Documentation](https://www.postgresql.org/docs/current/functions-matching.html)
