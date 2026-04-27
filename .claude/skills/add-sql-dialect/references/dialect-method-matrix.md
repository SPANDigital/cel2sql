# Dialect Method Matrix

## Contents

- Literals
- Operators
- Type casting
- Arrays
- JSON
- Timestamps
- String functions
- Comprehensions
- Validation + capabilities
- How to use this matrix

## How to use this matrix

When implementing a new dialect, scan each row to find which existing dialect's form your target engine matches, then copy that implementation. The full method list (~47 methods) lives in `dialect/dialect.go`; only the most-divergent ones are tabled here.

## Literals

| Method | PostgreSQL | DuckDB | MySQL | SQLite | BigQuery | Spark |
|---|---|---|---|---|---|---|
| `WriteStringLiteral` | `'v'` (`''` esc) | `'v'` (`''` esc) | `'v'` (`''` esc) | `'v'` (`''` esc) | `'v'` (`\'` esc) | `'v'` (`''` esc) |
| `WriteBytesLiteral` | `'\xHH…'` | `'\xhh…'` | `X'HH…'` | `X'HH…'` | `b"\OOO…"` (octal) | `X'HH…'` |
| `WriteParamPlaceholder` | `$1, $2` | `$1, $2` | `?, ?` | `?, ?` | `@p1, @p2` | `?, ?` |

## Operators

| Method | PostgreSQL | DuckDB | MySQL | SQLite | BigQuery | Spark |
|---|---|---|---|---|---|---|
| `WriteStringConcat` | `lhs \|\| rhs` | `lhs \|\| rhs` | `CONCAT(lhs, rhs)` | `lhs \|\| rhs` | `lhs \|\| rhs` | `concat(lhs, rhs)` |
| `WriteRegexMatch` | `t ~ 'p'` / `~* 'p'` | `t ~ 'p'` / `~* 'p'` | `t REGEXP 'p'` | unsupported | `REGEXP_CONTAINS(t, 'p')` | `t RLIKE 'p'` |
| `WriteLikeEscape` | `ESCAPE E'\\\\'` | `ESCAPE '\\\\'` | `ESCAPE '\\\\'` | `ESCAPE '\\'` | (no-op; default backslash) | `ESCAPE '\\\\'` |
| `WriteArrayMembership` | `e = ANY(a)` | `e = ANY(a)` | `JSON_CONTAINS(a, CAST(e AS JSON))` | `e IN (SELECT value FROM json_each(a))` | `e IN UNNEST(a)` | `array_contains(a, e)` |

## Type casting

| Method | PostgreSQL | DuckDB | MySQL | SQLite | BigQuery | Spark |
|---|---|---|---|---|---|---|
| `WriteCastToNumeric` | `::numeric` | `::DOUBLE` | ` + 0` | ` + 0` | `::FLOAT64` | ` + 0` |
| `WriteTypeName` (int) | `BIGINT` | `BIGINT` | `SIGNED` | `INTEGER` | `INT64` | `BIGINT` |
| `WriteTypeName` (string) | `TEXT` | `VARCHAR` | `CHAR` | `TEXT` | `STRING` | `STRING` |
| `WriteTypeName` (bool) | `BOOLEAN` | `BOOLEAN` | `UNSIGNED` | `INTEGER` | `BOOL` | `BOOLEAN` |
| `WriteTypeName` (bytes) | `BYTEA` | `BLOB` | `BINARY` | `BLOB` | `BYTES` | `BINARY` |
| `WriteTimestampCast` | `CAST(e AS TIMESTAMP WITH TIME ZONE)` | `CAST(e AS TIMESTAMPTZ)` | `CAST(e AS DATETIME)` | `datetime(e)` | `CAST(e AS TIMESTAMP)` | `CAST(e AS TIMESTAMP)` |
| `WriteEpochExtract` | `EXTRACT(EPOCH FROM e)::bigint` | `EXTRACT(EPOCH FROM e)::BIGINT` | `UNIX_TIMESTAMP(e)` | `CAST(strftime('%s', e) AS INTEGER)` | `UNIX_SECONDS(e)` | `UNIX_TIMESTAMP(e)` |

## Arrays

| Method | PostgreSQL | DuckDB | BigQuery | Spark | MySQL/SQLite |
|---|---|---|---|---|---|
| `WriteArrayLiteralOpen` | `ARRAY[` | `[` | `[` | `array(` | n/a (JSON arrays) |
| `WriteArrayLiteralClose` | `]` | `]` | `]` | `)` | n/a |
| `WriteArrayLength` | `COALESCE(ARRAY_LENGTH(e, dim), 0)` | `COALESCE(array_length(e), 0)` | `ARRAY_LENGTH(e)` | `COALESCE(size(e), 0)` (errors on dim>1) | varies |
| `WriteListIndex` | `a[i + 1]` | `a[i + 1]` | `a[OFFSET(i)]` | `a[i]` | varies |
| `WriteListIndexConst` | `a[i+1]` | `a[i+1]` | `a[OFFSET(i)]` | `a[i]` | varies |
| `WriteEmptyTypedArray` | `ARRAY[]::T[]` | `[]::T[]` | `ARRAY<T>[]` | `CAST(array() AS ARRAY<T>)` | `JSON_ARRAY()` (MySQL); `json_array()` (SQLite) |

## JSON

| Method | PostgreSQL | DuckDB | MySQL | SQLite | BigQuery | Spark |
|---|---|---|---|---|---|---|
| `WriteJSONFieldAccess` (final) | `b->>'f'` | `b->>'f'` | `b->>'$.f'` | `json_extract(b, '$.f')` | `JSON_VALUE(b, '$.f')` | `get_json_object(b, '$.f')` |
| `WriteJSONFieldAccess` (intermediate) | `b->'f'` | `b->'f'` | `b->'$.f'` | `json_extract(b, '$.f')` | `JSON_QUERY(b, '$.f')` | `get_json_object(b, '$.f')` |
| `WriteJSONExistence` | `b ? 'f'` | `json_exists(b, '$.f')` | `JSON_CONTAINS_PATH(b, 'one', '$.f')` | `json_type(b, '$.f') IS NOT NULL` | `JSON_VALUE(b, '$.f') IS NOT NULL` | `get_json_object(b, '$.f') IS NOT NULL` |
| `WriteJSONArrayElements` | `jsonb_array_elements(e)` (must be set-returning!) | `json_each(e)` | varies | `json_each(e)` | `UNNEST(JSON_QUERY_ARRAY(e))` | `EXPLODE(from_json(e, 'ARRAY<STRING>'))` |
| `WriteJSONArrayMembership` (RHS of `lhs = `) | `ANY(ARRAY(SELECT jsonb_func(e)))` | `(SELECT value FROM json_each(e))` | `JSON_CONTAINS(e, CAST(? AS JSON))` | `(SELECT value FROM json_each(e))` | `UNNEST(JSON_VALUE_ARRAY(e))` | `(SELECT col FROM (SELECT EXPLODE(from_json(e, 'ARRAY<STRING>')) AS col) t)` |

## Timestamps

| Method | PostgreSQL | DuckDB | MySQL | SQLite | BigQuery | Spark |
|---|---|---|---|---|---|---|
| `WriteDuration` | `INTERVAL N UNIT` | `INTERVAL N UNIT` | `INTERVAL N UNIT` | `'+N units'` | `INTERVAL N UNIT` | `INTERVAL N UNIT` |
| `WriteTimestampArithmetic` | `ts +/- iv` | `ts +/- iv` | `ts +/- iv` | `datetime(ts, '...')` | `TIMESTAMP_ADD/SUB(ts, iv)` | `ts +/- iv` |
| `WriteExtract` (DOW) | `EXTRACT(DOW FROM e)` | `(EXTRACT(DOW FROM e) + 6) % 7` | varies | `strftime('%w', e)` | `EXTRACT(DAYOFWEEK FROM e) - 1` | `(dayofweek(e) - 1)` |

DOW convention: CEL wants `Sunday=0..Saturday=6`. PostgreSQL emits Sunday=0 already; DuckDB returns Monday=1..Sunday=0 (off by one in the wrong direction, hence the `+6) % 7` adjustment); BigQuery `DAYOFWEEK` returns Sunday=1..Saturday=7 (subtract 1); Spark `dayofweek()` matches BigQuery.

## String functions

| Method | PostgreSQL | DuckDB | MySQL | SQLite | BigQuery | Spark |
|---|---|---|---|---|---|---|
| `WriteContains` | `POSITION(n IN h) > 0` | `CONTAINS(h, n)` | `LOCATE(n, h) > 0` | `INSTR(h, n) > 0` | `STRPOS(h, n) > 0` | `LOCATE(n, h) > 0` |
| `WriteSplit` | `STRING_TO_ARRAY(s, d)` | `STRING_SPLIT(s, d)` | varies (JSON arrays) | varies | `SPLIT(s, d)` | `split(s, d)` |
| `WriteJoin` | `ARRAY_TO_STRING(a, d, '')` | `ARRAY_TO_STRING(a, d)` | varies | varies | `ARRAY_TO_STRING(a, d)` | `array_join(a, d)` |

## Comprehensions

| Method | PostgreSQL | DuckDB | BigQuery | Spark | SQLite |
|---|---|---|---|---|---|
| `WriteUnnest` | `UNNEST(s)` | `UNNEST(s)` | `UNNEST(s)` | `EXPLODE(s)` | `json_each(s)` |
| `WriteArraySubqueryOpen` | `ARRAY(SELECT ` | `ARRAY(SELECT ` | `ARRAY(SELECT ` | `(SELECT collect_list(` | `(SELECT json_group_array(` |
| `WriteArraySubqueryExprClose` | `` | `` | `` | `)` | `)` |

## Validation + capabilities

| Method | PostgreSQL | DuckDB | MySQL | SQLite | BigQuery | Spark |
|---|---|---|---|---|---|---|
| `MaxIdentifierLength` | 63 | 0 (unlimited) | 64 | 0 (unlimited) | 300 | 128 |
| `SupportsRegex` | true | true | true | **false** | true | true |
| `SupportsNativeArrays` | true | true | false | false | true | true |
| `SupportsJSONB` | true | false | false | false | false | false |
| `SupportsIndexAnalysis` | true | true | true | true | true | **false** |
