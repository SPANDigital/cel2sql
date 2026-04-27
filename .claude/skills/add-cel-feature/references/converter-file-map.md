# Converter File Map

Where each CEL surface area lives in the converter. Use this to find the visitor that needs updating when adding a feature.

## Contents

- Top-level dispatch
- Operators and special-cased operators
- Comprehensions
- JSON / JSONB
- Timestamps and durations
- Index analysis
- Where it doesn't go

## Top-level dispatch

| File | Owns |
|---|---|
| `cel2sql.go` | The `converter` struct, `Convert` / `ConvertParameterized` entry points, the `visit` dispatch (`visitCall`, `visitSelect`, `visitIdent`, `visitList`, `visitConst`, `visitStruct`, `visitComprehension`), parameter-placeholder writing, identifier/struct-field validation, recursion-depth and output-length guards. |

`visitCall` is the dispatch hub for nearly every CEL function and operator. New function support starts here, almost always by adding a new branch to the overload-ID switch (search for `overloads.Add`, `overloads.Contains`, etc., to see existing wiring).

## Operators and special-cased operators

| File | Owns |
|---|---|
| `operators.go` | `standardSQLBinaryOperators` map (CEL operator → SQL token), reverse lookup, helper predicates like `isStringLiteral`, `isNullLiteral`, `isBoolLiteral`, `isListType`, `isNumericComparison`, `isFieldAccessExpression`. |
| `cel2sql.go` `visitCallBinary` | Binary operator emission with special cases (`||` for string/list concat, `IS`/`IS NOT` for null/bool comparisons, JSON-text numeric coercion via `WriteCastToNumeric`, IN-list dispatch to `WriteArrayMembership` or `WriteJSONArrayMembership`). |

Adding a new operator overload (e.g., timestamp + duration) usually requires updates here.

## Comprehensions

| File | Owns |
|---|---|
| `comprehensions.go` | The `ComprehensionInfo` struct, `identifyComprehension`, comprehension-shape predicates. |
| `cel2sql.go` `visitComprehension`, `visitAllComprehension`, `visitExistsComprehension`, `visitExistsOneComprehension`, `visitFilterComprehension`, `visitMapComprehension`, `writeComprehensionSource` | Per-shape SQL emission. |

The comprehension source is dispatched to either `dialect.WriteUnnest` (native arrays) or `dialect.WriteJSONArrayElements` (JSON arrays).

## JSON / JSONB

| File | Owns |
|---|---|
| `json.go` | JSON-path detection (`shouldUseJSONPath`, `hasJSONFieldInChain`), JSON-variable handling (`isJSONVariable` for `WithJSONVariables`), `has()` macro emission, JSON path construction (`buildJSONPath`, `buildJSONPathInternal`, `buildJSONPathForArray`), JSON field name escaping. |

The detection chain is: schema-driven (via `isFieldJSON`) → JSON-variable (via `WithJSONVariables`) → fallback to chain-walk (`hasJSONFieldInChain`). Adding new JSON behaviour usually plugs into `shouldUseJSONPath` plus a new `dialect.Dialect` method.

## Timestamps and durations

| File | Owns |
|---|---|
| `timestamps.go` | Custom CEL types for `DATE` / `TIME` / `DATETIME` / `INTERVAL`, duration parsing, EXTRACT helpers, timestamp arithmetic dispatch. |

Adding a new timestamp accessor (e.g., `getQuarter()`) goes here. Day-of-week conversion is handled per-dialect via `Dialect.WriteExtract` (CEL convention: Sunday=0).

## Index analysis

| File | Owns |
|---|---|
| `analysis.go` | `AnalyzeQuery` entry point; pattern detection (`PatternComparison`, `PatternJSONAccess`, `PatternRegexMatch`, `PatternArrayMembership`, `PatternArrayComprehension`, `PatternJSONArrayComprehension`); pattern → DDL delegation to `dialect.IndexAdvisor`. |

Adding a new pattern type means: extend `analysis.go` to detect it, add a `PatternType` constant in `dialect/index_advisor.go`, then implement the new pattern in each dialect's `index_advisor.go`.

## Where it doesn't go

- Per-dialect SQL syntax → `dialect/<name>/dialect.go`. The converter is dialect-agnostic.
- Reserved keywords / identifier rules → `dialect/<name>/validation.go`.
- Regex pattern translation (RE2 → engine-native) → `dialect/<name>/regex.go`.
- Type provider / database-introspection → `<name>/provider.go`.
