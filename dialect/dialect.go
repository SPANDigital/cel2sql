// Package dialect defines the interface for SQL dialect-specific code generation.
// Each supported database implements this interface to produce correct SQL syntax.
package dialect

import (
	"errors"
	"strings"
)

// Name represents a SQL dialect name.
type Name string

// Supported SQL dialect names.
const (
	PostgreSQL Name = "postgresql"
	MySQL      Name = "mysql"
	SQLite     Name = "sqlite"
	DuckDB     Name = "duckdb"
	BigQuery   Name = "bigquery"
)

// ErrUnsupportedFeature indicates that the requested feature is not supported by this dialect.
var ErrUnsupportedFeature = errors.New("unsupported dialect feature")

// Dialect defines the interface for SQL dialect-specific code generation.
// The converter calls these methods at every point where SQL syntax diverges
// between databases. Methods receive a *strings.Builder that shares the
// converter's output buffer, and callback functions for writing sub-expressions.
type Dialect interface {
	// Name returns the dialect name.
	Name() Name

	// --- Literals ---

	// WriteStringLiteral writes a string literal in the dialect's syntax.
	// For PostgreSQL: 'value' with '' escaping.
	WriteStringLiteral(w *strings.Builder, value string)

	// WriteBytesLiteral writes a byte array literal in the dialect's syntax.
	// For PostgreSQL: '\xDEADBEEF'.
	WriteBytesLiteral(w *strings.Builder, value []byte) error

	// WriteParamPlaceholder writes a parameter placeholder.
	// For PostgreSQL: $1, $2. For MySQL: ?, ?. For BigQuery: @p1, @p2.
	WriteParamPlaceholder(w *strings.Builder, paramIndex int)

	// --- Operators ---

	// WriteStringConcat writes a string concatenation expression.
	// For PostgreSQL: lhs || rhs. For MySQL: CONCAT(lhs, rhs).
	WriteStringConcat(w *strings.Builder, writeLHS, writeRHS func() error) error

	// WriteRegexMatch writes a regex match expression.
	// For PostgreSQL: expr ~ 'pattern' or expr ~* 'pattern'.
	WriteRegexMatch(w *strings.Builder, writeTarget func() error, pattern string, caseInsensitive bool) error

	// WriteLikeEscape writes the LIKE escape clause.
	// For PostgreSQL: ESCAPE E'\\'. For MySQL: ESCAPE '\\'.
	WriteLikeEscape(w *strings.Builder)

	// WriteArrayMembership writes an array membership test.
	// For PostgreSQL: elem = ANY(array). For MySQL: JSON_CONTAINS(array, elem).
	WriteArrayMembership(w *strings.Builder, writeElem func() error, writeArray func() error) error

	// --- Type Casting ---

	// WriteCastToNumeric writes a cast to numeric type.
	// For PostgreSQL: ::numeric. For MySQL: CAST(... AS DECIMAL).
	WriteCastToNumeric(w *strings.Builder)

	// WriteTypeName writes a type name for CAST expressions.
	// For PostgreSQL: BOOLEAN, BYTEA, DOUBLE PRECISION, BIGINT, TEXT.
	WriteTypeName(w *strings.Builder, celTypeName string)

	// WriteEpochExtract writes extraction of epoch from a timestamp.
	// For PostgreSQL: EXTRACT(EPOCH FROM expr)::bigint.
	WriteEpochExtract(w *strings.Builder, writeExpr func() error) error

	// WriteTimestampCast writes a cast to timestamp type.
	// For PostgreSQL: CAST(expr AS TIMESTAMP WITH TIME ZONE).
	WriteTimestampCast(w *strings.Builder, writeExpr func() error) error

	// --- Arrays ---

	// WriteArrayLiteralOpen writes the opening of an array literal.
	// For PostgreSQL: ARRAY[. For DuckDB: [.
	WriteArrayLiteralOpen(w *strings.Builder)

	// WriteArrayLiteralClose writes the closing of an array literal.
	// For PostgreSQL: ]. For DuckDB: ].
	WriteArrayLiteralClose(w *strings.Builder)

	// WriteArrayLength writes an array length expression.
	// For PostgreSQL: COALESCE(ARRAY_LENGTH(expr, dimension), 0).
	WriteArrayLength(w *strings.Builder, dimension int, writeExpr func() error) error

	// WriteListIndex writes a list index expression.
	// For PostgreSQL: array[index + 1] (1-indexed).
	WriteListIndex(w *strings.Builder, writeArray func() error, writeIndex func() error) error

	// WriteListIndexConst writes a constant list index.
	// For PostgreSQL: array[idx+1] (converts 0-indexed to 1-indexed).
	WriteListIndexConst(w *strings.Builder, writeArray func() error, index int64) error

	// WriteEmptyTypedArray writes an empty typed array literal.
	// For PostgreSQL: ARRAY[]::text[].
	WriteEmptyTypedArray(w *strings.Builder, typeName string)

	// --- JSON ---

	// WriteJSONFieldAccess writes JSON field access.
	// For PostgreSQL: base->>'field' (text) or base->'field' (json).
	// For SQLite: json_extract(base, '$.field').
	// writeBase writes the base expression; the dialect wraps or appends as needed.
	WriteJSONFieldAccess(w *strings.Builder, writeBase func() error, fieldName string, isFinal bool) error

	// WriteJSONExistence writes a JSON key existence check.
	// For PostgreSQL (JSONB): ? 'key'. For PostgreSQL (JSON): ->'key' IS NOT NULL.
	WriteJSONExistence(w *strings.Builder, isJSONB bool, fieldName string, writeBase func() error) error

	// WriteJSONArrayElements writes a call to extract JSON array elements.
	// For PostgreSQL: jsonb_array_elements(expr) or json_array_elements(expr).
	WriteJSONArrayElements(w *strings.Builder, isJSONB bool, asText bool, writeExpr func() error) error

	// WriteJSONArrayLength writes a JSON array length expression.
	// For PostgreSQL: COALESCE(jsonb_array_length(expr), 0).
	WriteJSONArrayLength(w *strings.Builder, writeExpr func() error) error

	// WriteJSONExtractPath writes a JSON path extraction function.
	// For PostgreSQL: jsonb_extract_path_text(root, 'seg1', 'seg2').
	WriteJSONExtractPath(w *strings.Builder, pathSegments []string, writeRoot func() error) error

	// WriteJSONArrayMembership writes a JSON array membership test for the IN operator.
	// For PostgreSQL: ANY(ARRAY(SELECT jsonb_func(expr))).
	WriteJSONArrayMembership(w *strings.Builder, jsonFunc string, writeExpr func() error) error

	// WriteNestedJSONArrayMembership writes a nested JSON array membership test.
	// For PostgreSQL: ANY(ARRAY(SELECT jsonb_array_elements_text(expr))).
	WriteNestedJSONArrayMembership(w *strings.Builder, writeExpr func() error) error

	// --- Timestamps ---

	// WriteDuration writes a duration/interval literal.
	// For PostgreSQL: INTERVAL N UNIT.
	WriteDuration(w *strings.Builder, value int64, unit string)

	// WriteInterval writes an INTERVAL expression from a variable.
	// For PostgreSQL: INTERVAL expr UNIT.
	WriteInterval(w *strings.Builder, writeValue func() error, unit string) error

	// WriteExtract writes a timestamp field extraction expression.
	// Handles DOW conversion, Month/DOY adjustment, and timezone support.
	WriteExtract(w *strings.Builder, part string, writeExpr func() error, writeTZ func() error) error

	// WriteTimestampArithmetic writes timestamp arithmetic.
	// For PostgreSQL: timestamp +/- interval.
	WriteTimestampArithmetic(w *strings.Builder, op string, writeTS, writeDur func() error) error

	// --- String Functions ---

	// WriteContains writes a string contains expression.
	// For PostgreSQL: POSITION(needle IN haystack) > 0.
	WriteContains(w *strings.Builder, writeHaystack, writeNeedle func() error) error

	// WriteSplit writes a string split expression.
	// For PostgreSQL: STRING_TO_ARRAY(string, delimiter).
	WriteSplit(w *strings.Builder, writeStr, writeDelim func() error) error

	// WriteSplitWithLimit writes a string split expression with a limit.
	// For PostgreSQL: (STRING_TO_ARRAY(string, delimiter))[1:limit].
	WriteSplitWithLimit(w *strings.Builder, writeStr, writeDelim func() error, limit int64) error

	// WriteJoin writes an array join expression.
	// For PostgreSQL: ARRAY_TO_STRING(array, delimiter, '').
	WriteJoin(w *strings.Builder, writeArray, writeDelim func() error) error

	// --- Comprehensions ---

	// WriteUnnest writes the UNNEST source for comprehensions.
	// For PostgreSQL: UNNEST(array). For MySQL: JSON_TABLE(...).
	WriteUnnest(w *strings.Builder, writeSource func() error) error

	// WriteArraySubqueryOpen writes the prefix before the transform expression
	// in an array-building subquery.
	// For PostgreSQL: "ARRAY(SELECT ". For SQLite: "(SELECT json_group_array(".
	WriteArraySubqueryOpen(w *strings.Builder)

	// WriteArraySubqueryExprClose writes the suffix after the transform expression
	// and before FROM in an array-building subquery.
	// For PostgreSQL: "" (nothing). For SQLite: ")".
	WriteArraySubqueryExprClose(w *strings.Builder)

	// --- Struct ---

	// WriteStructOpen writes the opening of a struct/row literal.
	// For PostgreSQL: ROW(. For BigQuery: STRUCT(.
	WriteStructOpen(w *strings.Builder)

	// WriteStructClose writes the closing of a struct/row literal.
	// For PostgreSQL: ). For BigQuery: ).
	WriteStructClose(w *strings.Builder)

	// --- Validation ---

	// MaxIdentifierLength returns the maximum identifier length for this dialect.
	// For PostgreSQL: 63. For MySQL: 64. For SQLite: unlimited (0).
	MaxIdentifierLength() int

	// ValidateFieldName validates a field name for this dialect.
	ValidateFieldName(name string) error

	// ReservedKeywords returns the set of reserved SQL keywords for this dialect.
	ReservedKeywords() map[string]bool

	// --- Regex ---

	// ConvertRegex converts an RE2 regex pattern to the dialect's native format.
	// Returns: (convertedPattern, caseInsensitive, error).
	ConvertRegex(re2Pattern string) (pattern string, caseInsensitive bool, err error)

	// SupportsRegex indicates whether this dialect supports regex matching.
	SupportsRegex() bool

	// --- Capabilities ---

	// SupportsNativeArrays indicates whether this dialect has native array types.
	SupportsNativeArrays() bool

	// SupportsJSONB indicates whether this dialect has a distinct JSONB type.
	SupportsJSONB() bool

	// SupportsIndexAnalysis indicates whether index analysis is supported.
	SupportsIndexAnalysis() bool
}
