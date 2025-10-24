# Security Guide

cel2sql includes comprehensive security protections against common attack vectors. All security features are enabled by default with zero configuration required.

## Table of Contents

- [Overview](#overview)
- [Field Name Validation](#field-name-validation)
- [JSON Field Escaping](#json-field-escaping)
- [ReDoS Protection](#redos-protection)
- [Recursion Depth Limits](#recursion-depth-limits)
- [Context Timeouts](#context-timeouts)
- [Best Practices](#best-practices)
- [Security Checklist](#security-checklist)

## Overview

cel2sql protects against:

| Threat | Protection | Status |
|--------|------------|--------|
| SQL Injection (field names) | Field name validation | ✅ Automatic |
| SQL Injection (JSON fields) | Quote escaping | ✅ Automatic |
| ReDoS (Regex DoS) | Pattern validation | ✅ Automatic |
| Stack overflow (deep nesting) | Recursion depth limits | ✅ Automatic |
| Resource exhaustion | Context timeouts | ⚙️ Optional |
| Information disclosure | Schema minimization | 📋 Manual |

All automatic protections require **no configuration** and are always active.

## Field Name Validation

### What It Protects Against

SQL injection via malicious field names:

```go
// ❌ Attack attempt
field'; DROP TABLE users; --

// ❌ Attack attempt
SELECT * FROM sensitive_data

// ❌ Attack attempt
user OR 1=1
```

### How It Works

All field names are validated against strict rules before SQL generation:

**Validation Rules:**
1. **Length**: Maximum 63 characters (PostgreSQL `NAMEDATALEN-1`)
2. **Format**: Must start with letter or underscore (`[a-zA-Z_]`)
3. **Characters**: Only alphanumeric and underscore (`[a-zA-Z0-9_]*`)
4. **Reserved words**: Blocks 60+ SQL reserved keywords
5. **Empty**: Empty strings are rejected

**Implementation:**
- Validation occurs in `utils.go:validateFieldName()`
- Applied in `cel2sql.go:visitSelect()` and `visitIdent()`

### Examples

```go
// ✅ Valid field names
user.name
product.price_usd
order_item_123
_internal_field

// ❌ Invalid field names
field'; DROP TABLE users--     // SQL injection attempt
123invalid                      // Starts with number
field-name                      // Contains hyphen
SELECT                          // Reserved keyword
user.field_with_very_long_name_that_exceeds_the_postgresql_maximum_identifier_length  // Too long
```

### Error Messages

```go
ast, _ := env.Compile(`field'; DROP TABLE users--.value == 1`)
sql, err := cel2sql.Convert(ast)
// err: "invalid field name 'field'; DROP TABLE users--': must start with letter or underscore and contain only alphanumeric characters and underscores"

ast, _ = env.Compile(`SELECT.value == 1`)
sql, err = cel2sql.Convert(ast)
// err: "invalid field name 'SELECT': reserved SQL keyword"
```

### Reserved Keywords

The following SQL keywords are blocked:

```
SELECT, INSERT, UPDATE, DELETE, DROP, CREATE, ALTER, TRUNCATE,
FROM, WHERE, JOIN, UNION, HAVING, GROUP, ORDER, BY, AS, ON,
AND, OR, NOT, NULL, TRUE, FALSE, CASE, WHEN, THEN, ELSE, END,
IN, EXISTS, BETWEEN, LIKE, IS, DISTINCT, ALL, ANY, SOME,
PRIMARY, FOREIGN, KEY, REFERENCES, CONSTRAINT, INDEX, TABLE,
DATABASE, SCHEMA, VIEW, PROCEDURE, FUNCTION, TRIGGER, GRANT,
REVOKE, COMMIT, ROLLBACK, TRANSACTION
```

## JSON Field Escaping

### What It Protects Against

SQL injection via malicious JSON field names:

```go
// ❌ Attack attempt
user.preferences.field' OR '1'='1

// ❌ Attack attempt
user.data.key'; DROP TABLE users--
```

### How It Works

Single quotes in JSON field names are automatically escaped by doubling them:

```
' → ''
```

This follows PostgreSQL's standard for escaping quotes in string literals.

**Implementation:**
- Escaping occurs in `utils.go:escapeJSONFieldName()`
- Applied in multiple locations:
  - `cel2sql.go:visitSelect()` - JSON path operators
  - `cel2sql.go:visitHasFunction()` - Existence checks
  - `cel2sql.go:visitNestedJSONHas()` - Path extraction
  - `json.go:buildJSONPath*()` - All JSON path construction

### Examples

```go
// Field name with single quote
user.preferences.user'name == "test"
// Generated SQL: user.preferences->>'user''name' = 'test'

// Malicious field name (neutralized)
user.data.field' OR '1'='1 == "value"
// Generated SQL: user.data->>'field'' OR ''1''=''1' = 'value'

// Multiple quotes
user.settings.key''with''quotes == "test"
// Generated SQL: user.settings->>'key''''with''''quotes' = 'test'
```

### JSON Operations Protected

All JSON operations are protected:

| Operation | Operator | Example |
|-----------|----------|---------|
| Field access | `->>` | `user.data->>'key''name'` |
| Nested access | `->` | `user.data->'nested''key'` |
| Existence check | `?` | `user.data ? 'key''name'` |
| Path extraction | `jsonb_extract_path_text()` | `jsonb_extract_path_text(user.data, 'key''name')` |

### Defense in Depth

JSON field escaping is applied at multiple stages:
1. **visitSelect()** - When constructing field access
2. **buildJSONPath()** - When building JSON path expressions
3. **visitHasFunction()** - When checking field existence
4. **visitNestedJSONHas()** - When extracting nested paths

This ensures protection even if one layer is bypassed.

## ReDoS Protection

### What It Protects Against

Regular Expression Denial of Service (ReDoS) attacks that cause catastrophic backtracking:

```go
// ❌ Dangerous pattern - can freeze your application
field.matches(r"(a+)+b")
// Input "aaaaaaaaaaaaaaaaaaX" → HOURS of processing time

// ❌ Quantified alternation
field.matches(r"(a|a)*b")
// Can cause exponential time complexity
```

**CVE References:**
- Similar to CVE-2019-16155 (ReDoS in third-party libraries)
- CWE-1333: Inefficient Regular Expression Complexity

### How It Works

All regex patterns are validated before conversion using a multi-layered approach:

**Implementation:**
- Pattern validation in `cel2sql.go:convertRE2ToPOSIX()`
- Validation occurs before pattern conversion
- Descriptive errors returned for blocked patterns

### Pattern Limits

#### 1. Length Limit (500 characters)

Prevents resource exhaustion from extremely long patterns:

```go
// ✅ Normal patterns allowed
email.matches(r"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$")  // 62 chars

// ❌ Excessive length blocked
field.matches(r"very long pattern..." * 100)  // > 500 chars
// Error: pattern exceeds 500 character limit
```

#### 2. Nested Quantifiers Detection

Uses state machine to detect nested quantifiers:

```go
// ❌ Blocked: Nested + quantifiers
field.matches(r"(a+)+")
// Error: nested quantifiers detected at position 4

// ❌ Blocked: Nested * quantifiers
field.matches(r"(a*)*")
// Error: nested quantifiers detected

// ❌ Blocked: Mixed nesting
field.matches(r"(a+)*b")
// Error: nested quantifiers detected

// ❌ Blocked: Deep nesting
field.matches(r"((a+)+)+")
// Error: nested quantifiers detected

// ✅ Allowed: Non-nested quantifiers
field.matches(r"(abc)+")     // OK
field.matches(r"a+b+c+")     // OK
field.matches(r"[a-z]+")     // OK
```

#### 3. Capture Group Limit (20 groups)

Prevents memory exhaustion from excessive groups:

```go
// ✅ Reasonable groups allowed
field.matches(r"(a)(b)(c)(d)(e)")  // 5 groups

// ❌ Too many groups blocked
field.matches(r"(a)(b)(c)...(z)(aa)(ab)(ac)")  // > 20 groups
// Error: pattern exceeds 20 capture group limit
```

#### 4. Quantified Alternation Detection

Detects problematic alternation patterns:

```go
// ❌ Blocked: Quantified alternation
field.matches(r"(a|a)+b")
// Error: quantified alternation detected

// ❌ Blocked: Subtle alternation issue
field.matches(r"(x|xy)+y")
// Error: quantified alternation detected

// ❌ Blocked: Overlapping alternatives
field.matches(r"(ab|a)+c")
// Error: quantified alternation detected

// ✅ Allowed: Safe alternation without quantifiers
field.matches(r"(cat|dog|bird)")        // OK
field.matches(r"^(active|inactive)$")   // OK
```

#### 5. Nesting Depth Limit (10 levels)

Prevents stack exhaustion from deeply nested patterns:

```go
// ✅ Reasonable nesting allowed
field.matches(r"((abc))")         // 2 levels

// ❌ Excessive nesting blocked
field.matches(r"((((((((((a))))))))))")  // 10+ levels
// Error: pattern exceeds nesting depth limit of 10
```

### Safe Patterns

These patterns are validated as safe:

```go
// ✅ Email validation
email.matches(r"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$")

// ✅ Phone number (US format)
phone.matches(r"^\d{3}-\d{3}-\d{4}$")

// ✅ URL validation
url.matches(r"^https?://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}")

// ✅ Username (alphanumeric, 3-20 chars)
username.matches(r"^[a-zA-Z0-9_-]{3,20}$")

// ✅ Zip code (US)
zipcode.matches(r"^\d{5}(-\d{4})?$")

// ✅ Case-insensitive search
description.matches(r"(?i)urgent|priority|important")

// ✅ UUID
id.matches(r"^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$")
```

### Attack Examples

Real-world ReDoS attack patterns (all blocked):

```go
// ❌ Classic nested quantifier attack
field.matches(r"(a+)+b")
field.matches(r"(a*)*b")

// ❌ Email ReDoS (similar to CVE-2019-16155)
field.matches(r"([a-zA-Z0-9_\-\.]+)+@example\.com")

// ❌ URL ReDoS
field.matches(r"^(https?://)?([a-z0-9-]+\.)*[a-z0-9-]+\.[a-z]{2,}$")

// ❌ Quantified alternation
field.matches(r"(a|a)*b")
field.matches(r"(x|xy)+y")

// ❌ Excessive backtracking
field.matches(r"(a|ab)*c")
field.matches(r"(a*)*b")
```

### Error Messages

Descriptive errors help developers understand why patterns are rejected:

```go
ast, _ := env.Compile(`field.matches(r"(a+)+")`)
sql, err := cel2sql.Convert(ast)
// err: "invalid regex pattern: nested quantifiers detected at position 4"

ast, _ = env.Compile(`field.matches(r"(a|a)*b")`)
sql, err = cel2sql.Convert(ast)
// err: "invalid regex pattern: quantified alternation detected"

ast, _ = env.Compile(`field.matches(r"((((((((((a))))))))))")`)
sql, err = cel2sql.Convert(ast)
// err: "invalid regex pattern: pattern exceeds nesting depth limit of 10"
```

## Recursion Depth Limits

### What It Protects Against

Deeply nested CEL expressions that could cause stack overflow or excessive resource consumption:

```go
// ❌ Attack: Deeply nested arithmetic (1000+ levels)
((((((((((x + 1) + 1) + 1) + 1) + 1) + 1)...)))))

// ❌ Attack: Nested boolean conditions
(((((x > 0 && x < 100) && (x > 0 && x < 100))...)))))

// ❌ Attack: Nested ternary operators
(((x > 0 ? (x > 0 ? (x > 0 ? 1 : 0) : 0) : 0)...)))
```

### How It Works

cel2sql tracks recursion depth during CEL expression tree traversal:

1. **Depth Counter**: Increments on each `visit()` call to AST nodes
2. **Automatic Limit**: Default maximum depth of **100** (configurable)
3. **Fast Rejection**: Expressions exceeding limit are rejected immediately
4. **CWE-674 Protection**: Guards against uncontrolled recursion

**Implementation:**

```go
func (con *converter) visit(expr *exprpb.Expr) error {
    con.depth++
    defer func() { con.depth-- }()

    if con.depth > con.maxDepth {
        return fmt.Errorf("expression exceeds maximum recursion depth of %d", con.maxDepth)
    }
    // ... convert expression
}
```

### Configuration

The default limit of 100 provides strong protection while supporting realistic expressions:

```go
// Default: Automatic depth limit of 100
sql, err := cel2sql.Convert(ast)

// Custom: Adjust depth limit for specific use cases
sql, err := cel2sql.Convert(ast,
    cel2sql.WithMaxDepth(200))  // Allow deeper nesting

// Conservative: Lower limit for untrusted input
sql, err := cel2sql.Convert(ast,
    cel2sql.WithMaxDepth(50))   // Stricter protection
```

### What Counts Toward Depth

Every AST node visit increments the depth counter:

| Expression Type | Approximate Depth per Level |
|----------------|----------------------------|
| Simple arithmetic `x + 1` | 3-5 nodes |
| Nested calls `((x + 1) + 1)` | 2-3x per nesting |
| Boolean logic `a && b` | 3-4 nodes |
| Ternary `a ? b : c` | 4-6 nodes |
| Function calls `int(x)` | 2-3 nodes |
| Comprehensions `list.all(x, ...)` | 5-10 nodes |

**Example:** A CEL expression like `((x + 1) + 1)` nested 50 times will consume approximately 150-200 depth units.

### Defense in Depth

Combine depth limits with other protections:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

sql, err := cel2sql.Convert(ast,
    cel2sql.WithContext(ctx),        // Timeout protection
    cel2sql.WithMaxDepth(100),       // Recursion protection
    cel2sql.WithSchemas(schemas))    // Type safety

if err != nil {
    // Handle: depth exceeded, timeout, or validation error
}
```

### Error Messages

Clear error messages indicate depth limit violations:

```go
_, err := cel2sql.Convert(deeplyNestedAST)
// err: "expression exceeds maximum recursion depth of 100"

_, err := cel2sql.Convert(deeplyNestedAST, cel2sql.WithMaxDepth(50))
// err: "expression exceeds maximum recursion depth of 50"
```

### Best Practices

1. **Use Default Limits**: The default depth of 100 handles realistic expressions
2. **Validate User Input**: Reject obviously malicious patterns before CEL compilation
3. **Monitor Depth Errors**: Log depth limit violations to detect attack patterns
4. **Combine Protections**: Use depth limits WITH context timeouts
5. **Test Your Limits**: Verify legitimate expressions work within configured limits

### CWE Reference

- **CWE-674**: Uncontrolled Recursion - Mitigated by automatic depth tracking

## Context Timeouts

### What It Protects Against

- Resource exhaustion from complex expressions
- Long-running conversions blocking your application
- Denial of service via computational complexity

### How It Works

Use Go's `context.Context` to add timeout and cancellation support:

**Implementation:**
- Context checking in `cel2sql.go:checkContext()`
- Checked at key recursion points:
  - `visit()` - Main traversal entry
  - `visitCall()` - Every function call
  - `visitComprehension()` - Before processing comprehensions

### Usage

```go
import (
    "context"
    "time"
)

// With timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

sql, err := cel2sql.Convert(ast,
    cel2sql.WithContext(ctx),
    cel2sql.WithSchemas(schemas))

if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        log.Println("Conversion timed out - expression too complex")
    } else if errors.Is(err, context.Canceled) {
        log.Println("Conversion cancelled")
    }
}
```

### Recommended Timeouts

| Use Case | Recommended Timeout | Rationale |
|----------|-------------------|-----------|
| User-provided expressions | 5 seconds | Prevents abuse |
| Admin/trusted expressions | 30 seconds | More complex allowed |
| Simple expressions | 1 second | Fast fail |
| Batch processing | 10 seconds | Balance throughput/safety |

### Example: User-Facing API

```go
func handleUserFilter(w http.ResponseWriter, r *http.Request) {
    celExpr := r.FormValue("filter")

    // 1. Basic validation
    if len(celExpr) > 1000 {
        http.Error(w, "Expression too long", http.StatusBadRequest)
        return
    }

    // 2. Compile CEL
    ast, issues := env.Compile(celExpr)
    if issues != nil && issues.Err() != nil {
        http.Error(w, "Invalid expression", http.StatusBadRequest)
        return
    }

    // 3. Convert with timeout
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

    schemas := provider.GetSchemas()
    sql, err := cel2sql.Convert(ast,
        cel2sql.WithContext(ctx),
        cel2sql.WithSchemas(schemas))

    if err != nil {
        if errors.Is(err, context.DeadlineExceeded) {
            http.Error(w, "Expression too complex", http.StatusBadRequest)
            return
        }
        http.Error(w, "Conversion failed", http.StatusInternalServerError)
        return
    }

    // 4. Execute query (with timeout)
    // ...
}
```

## Best Practices

### 1. Defense in Depth

Use multiple layers of protection:

```go
func convertUserExpression(celExpr string) (string, error) {
    // Layer 1: Input validation
    if len(celExpr) > 1000 {
        return "", errors.New("expression too long")
    }
    if strings.Count(celExpr, "(") > 20 {
        return "", errors.New("too many nested groups")
    }

    // Layer 2: CEL compilation
    ast, issues := env.Compile(celExpr)
    if issues != nil && issues.Err() != nil {
        return "", fmt.Errorf("compilation failed: %w", issues.Err())
    }

    // Layer 3: Context timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Layer 4: cel2sql conversion (automatic protections)
    schemas := provider.GetSchemas()
    sql, err := cel2sql.Convert(ast,
        cel2sql.WithContext(ctx),
        cel2sql.WithSchemas(schemas))

    return sql, err
}
```

### 2. Schema Minimization

Only expose necessary fields in your schema:

```go
// ✅ Good: Minimal schema with only needed fields
userSchema := pg.Schema{
    {Name: "id", Type: "bigint"},
    {Name: "name", Type: "text"},
    {Name: "email", Type: "text"},
    {Name: "active", Type: "boolean"},
    {Name: "created_at", Type: "timestamp with time zone"},
}

// ❌ Bad: Exposing sensitive fields
userSchema := pg.Schema{
    {Name: "id", Type: "bigint"},
    {Name: "name", Type: "text"},
    {Name: "password_hash", Type: "text"},        // ❌ Don't expose
    {Name: "ssn", Type: "text"},                  // ❌ Don't expose
    {Name: "credit_card", Type: "text"},          // ❌ Don't expose
    {Name: "api_token", Type: "text"},            // ❌ Don't expose
    {Name: "internal_notes", Type: "text"},       // ❌ Don't expose
}
```

### 3. Use Prepared Statements

Always use prepared statements when executing generated SQL:

```go
// ✅ Good: Prepared statement
sqlWhere, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
if err != nil {
    return err
}

query := "SELECT * FROM users WHERE " + sqlWhere
stmt, err := db.Prepare(query)
if err != nil {
    return err
}
defer stmt.Close()

rows, err := stmt.Query()

// ❌ Bad: String concatenation with user values
// (cel2sql generates safe SQL, but this pattern is risky if extended)
query := "SELECT * FROM users WHERE " + sqlWhere + " AND extra=" + userValue
```

### 4. Enable Logging

Monitor conversion patterns in production:

```go
import "log/slog"

// JSON logging for production
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,
}))

sql, err := cel2sql.Convert(ast,
    cel2sql.WithContext(ctx),
    cel2sql.WithSchemas(schemas),
    cel2sql.WithLogger(logger))

// Logs include:
// - JSON path detection decisions
// - Comprehension types identified
// - Schema lookups
// - Regex pattern conversions
// - Performance metrics
// - Error contexts
```

### 5. Rate Limiting

Implement rate limiting for user-provided expressions:

```go
import "golang.org/x/time/rate"

// Per-user rate limiter
limiter := rate.NewLimiter(rate.Limit(10), 20)  // 10 req/sec, burst 20

func handleFilter(w http.ResponseWriter, r *http.Request) {
    userID := getUserID(r)

    if !limiter.Allow() {
        http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
        return
    }

    // ... convert expression
}
```

### 6. Input Validation

Validate user input before CEL compilation:

```go
func validateCELExpression(expr string) error {
    // Length check
    if len(expr) > 1000 {
        return errors.New("expression too long")
    }

    // Complexity check
    if strings.Count(expr, "(") > 20 {
        return errors.New("too many nested groups")
    }

    // Character whitelist (optional)
    for _, r := range expr {
        if !unicode.IsPrint(r) && !unicode.IsSpace(r) {
            return errors.New("invalid characters in expression")
        }
    }

    return nil
}
```

### 7. Error Handling

Handle errors securely without leaking information:

```go
sql, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
if err != nil {
    // Log detailed error for debugging
    log.Printf("Conversion error: %v", err)

    // Return generic error to user
    return "", errors.New("invalid filter expression")
}
```

### 8. Testing

Test your security assumptions:

```go
func TestSecurityProtections(t *testing.T) {
    tests := []struct {
        name        string
        expression  string
        shouldError bool
    }{
        {
            name:        "SQL injection via field name",
            expression:  `field'; DROP TABLE users--.value == 1`,
            shouldError: true,
        },
        {
            name:        "ReDoS nested quantifiers",
            expression:  `field.matches(r"(a+)+")`,
            shouldError: true,
        },
        {
            name:        "Valid safe expression",
            expression:  `user.age >= 18 && user.active`,
            shouldError: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ast, issues := env.Compile(tt.expression)
            if issues != nil && issues.Err() != nil {
                if !tt.shouldError {
                    t.Errorf("Expected success but got error: %v", issues.Err())
                }
                return
            }

            _, err := cel2sql.Convert(ast)
            if tt.shouldError && err == nil {
                t.Errorf("Expected error but got success")
            } else if !tt.shouldError && err != nil {
                t.Errorf("Expected success but got error: %v", err)
            }
        })
    }
}
```

## Security Checklist

Use this checklist when deploying cel2sql in production:

### Configuration

- [ ] Using `WithSchemas()` for JSON/JSONB detection
- [ ] Schemas include only necessary fields
- [ ] No sensitive fields (passwords, tokens, SSN, etc.) in schemas
- [ ] Context timeout configured (5 seconds recommended)
- [ ] Structured logging enabled (with `WithLogger()`)

### Input Validation

- [ ] Expression length limited (1000 chars recommended)
- [ ] Character whitelist applied if needed
- [ ] Complexity checks implemented (nesting depth, etc.)
- [ ] Rate limiting implemented for user-provided expressions

### Error Handling

- [ ] Detailed errors logged for debugging
- [ ] Generic errors returned to users (no information leakage)
- [ ] Context timeout errors handled gracefully
- [ ] Conversion errors don't cause panic

### Deployment

- [ ] Using prepared statements for query execution
- [ ] Database user has minimal required permissions
- [ ] Monitoring/alerting configured for errors
- [ ] Security tests included in CI/CD pipeline
- [ ] Documentation updated with security considerations

### Testing

- [ ] SQL injection tests (field names, JSON fields)
- [ ] ReDoS pattern tests (nested quantifiers, etc.)
- [ ] Context timeout tests
- [ ] Edge case tests (long expressions, deep nesting)
- [ ] Integration tests with real PostgreSQL database

### Monitoring

- [ ] Conversion errors tracked
- [ ] Timeout frequency monitored
- [ ] Pattern complexity tracked (via logging)
- [ ] Performance metrics collected
- [ ] Security alerts configured

## See Also

- [Getting Started Guide](getting-started.md) - Basic usage and setup
- [JSON/JSONB Support](json-support.md) - JSON security details
- [Regex Matching](regex-matching.md) - ReDoS protection details
- [Operators Reference](operators-reference.md) - Operator security

## Security Contact

To report security vulnerabilities, please open an issue on [GitHub](https://github.com/SPANDigital/cel2sql/issues) with the "security" label.
