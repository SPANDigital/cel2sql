// Package bigquery implements the BigQuery SQL dialect for cel2sql.
package bigquery

import (
	"fmt"
	"strings"

	"github.com/spandigital/cel2sql/v3/dialect"
)

// Dialect implements dialect.Dialect for BigQuery.
type Dialect struct{}

// New creates a new BigQuery dialect.
func New() *Dialect {
	return &Dialect{}
}

func init() {
	dialect.Register(dialect.BigQuery, func() dialect.Dialect { return New() })
}

// Ensure Dialect implements dialect.Dialect at compile time.
var _ dialect.Dialect = (*Dialect)(nil)

// Name returns the dialect name.
func (d *Dialect) Name() dialect.Name { return dialect.BigQuery }

// --- Literals ---

// WriteStringLiteral writes a BigQuery string literal with ” escaping.
func (d *Dialect) WriteStringLiteral(w *strings.Builder, value string) {
	escaped := strings.ReplaceAll(value, "'", "\\'")
	w.WriteString("'")
	w.WriteString(escaped)
	w.WriteString("'")
}

// WriteBytesLiteral writes a BigQuery octal-encoded byte literal (b"...").
func (d *Dialect) WriteBytesLiteral(w *strings.Builder, value []byte) error {
	w.WriteString("b\"")
	for _, b := range value {
		fmt.Fprintf(w, "\\%03o", b)
	}
	w.WriteString("\"")
	return nil
}

// WriteParamPlaceholder writes a BigQuery named parameter (@p1, @p2, ...).
func (d *Dialect) WriteParamPlaceholder(w *strings.Builder, paramIndex int) {
	fmt.Fprintf(w, "@p%d", paramIndex)
}

// --- Operators ---

// WriteStringConcat writes BigQuery string concatenation using ||.
func (d *Dialect) WriteStringConcat(w *strings.Builder, writeLHS, writeRHS func() error) error {
	if err := writeLHS(); err != nil {
		return err
	}
	w.WriteString(" || ")
	return writeRHS()
}

// WriteRegexMatch writes a BigQuery regex match using REGEXP_CONTAINS.
func (d *Dialect) WriteRegexMatch(w *strings.Builder, writeTarget func() error, pattern string, _ bool) error {
	w.WriteString("REGEXP_CONTAINS(")
	if err := writeTarget(); err != nil {
		return err
	}
	w.WriteString(", '")
	escaped := strings.ReplaceAll(pattern, "'", "\\'")
	w.WriteString(escaped)
	w.WriteString("')")
	return nil
}

// WriteLikeEscape is a no-op for BigQuery.
// BigQuery uses backslash as the default escape character in LIKE patterns
// and does not support the ESCAPE keyword.
func (d *Dialect) WriteLikeEscape(_ *strings.Builder) {
}

// WriteArrayMembership writes a BigQuery array membership test using IN UNNEST().
func (d *Dialect) WriteArrayMembership(w *strings.Builder, writeElem func() error, writeArray func() error) error {
	if err := writeElem(); err != nil {
		return err
	}
	w.WriteString(" IN UNNEST(")
	if err := writeArray(); err != nil {
		return err
	}
	w.WriteString(")")
	return nil
}

// --- Type Casting ---

// WriteCastToNumeric writes a BigQuery numeric cast (CAST(... AS FLOAT64)).
func (d *Dialect) WriteCastToNumeric(w *strings.Builder) {
	// BigQuery doesn't have a ::type cast syntax; this is used after expressions.
	// For BigQuery, the converter should use CAST(expr AS FLOAT64) instead.
	w.WriteString("::FLOAT64")
}

// WriteTypeName writes a BigQuery type name for CAST expressions.
func (d *Dialect) WriteTypeName(w *strings.Builder, celTypeName string) {
	switch celTypeName {
	case "bool":
		w.WriteString("BOOL")
	case "bytes":
		w.WriteString("BYTES")
	case "double":
		w.WriteString("FLOAT64")
	case "int":
		w.WriteString("INT64")
	case "string":
		w.WriteString("STRING")
	case "uint":
		w.WriteString("INT64")
	default:
		w.WriteString(strings.ToUpper(celTypeName))
	}
}

// WriteEpochExtract writes UNIX_SECONDS(expr) for BigQuery.
func (d *Dialect) WriteEpochExtract(w *strings.Builder, writeExpr func() error) error {
	w.WriteString("UNIX_SECONDS(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString(")")
	return nil
}

// WriteTimestampCast writes a BigQuery CAST to TIMESTAMP.
func (d *Dialect) WriteTimestampCast(w *strings.Builder, writeExpr func() error) error {
	w.WriteString("CAST(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString(" AS TIMESTAMP)")
	return nil
}

// --- Arrays ---

// WriteArrayLiteralOpen writes the BigQuery array literal opening ([).
func (d *Dialect) WriteArrayLiteralOpen(w *strings.Builder) {
	w.WriteString("[")
}

// WriteArrayLiteralClose writes the BigQuery array literal closing (]).
func (d *Dialect) WriteArrayLiteralClose(w *strings.Builder) {
	w.WriteString("]")
}

// WriteArrayLength writes ARRAY_LENGTH(expr) for BigQuery.
func (d *Dialect) WriteArrayLength(w *strings.Builder, _ int, writeExpr func() error) error {
	w.WriteString("ARRAY_LENGTH(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString(")")
	return nil
}

// WriteListIndex writes BigQuery 0-indexed array access using OFFSET.
func (d *Dialect) WriteListIndex(w *strings.Builder, writeArray func() error, writeIndex func() error) error {
	if err := writeArray(); err != nil {
		return err
	}
	w.WriteString("[OFFSET(")
	if err := writeIndex(); err != nil {
		return err
	}
	w.WriteString(")]")
	return nil
}

// WriteListIndexConst writes BigQuery constant array index access using OFFSET.
func (d *Dialect) WriteListIndexConst(w *strings.Builder, writeArray func() error, index int64) error {
	if err := writeArray(); err != nil {
		return err
	}
	fmt.Fprintf(w, "[OFFSET(%d)]", index)
	return nil
}

// WriteEmptyTypedArray writes an empty BigQuery typed array.
func (d *Dialect) WriteEmptyTypedArray(w *strings.Builder, typeName string) {
	w.WriteString("ARRAY<")
	w.WriteString(bigqueryTypeName(typeName))
	w.WriteString(">[]")
}

// --- JSON ---

// WriteJSONFieldAccess writes BigQuery JSON field access using JSON_VALUE.
func (d *Dialect) WriteJSONFieldAccess(w *strings.Builder, writeBase func() error, fieldName string, isFinal bool) error {
	escaped := escapeJSONFieldName(fieldName)
	if isFinal {
		w.WriteString("JSON_VALUE(")
	} else {
		w.WriteString("JSON_QUERY(")
	}
	if err := writeBase(); err != nil {
		return err
	}
	w.WriteString(", '$.")
	w.WriteString(escaped)
	w.WriteString("')")
	return nil
}

// WriteJSONExistence writes a BigQuery JSON key existence check.
func (d *Dialect) WriteJSONExistence(w *strings.Builder, _ bool, fieldName string, writeBase func() error) error {
	escaped := escapeJSONFieldName(fieldName)
	w.WriteString("JSON_VALUE(")
	if err := writeBase(); err != nil {
		return err
	}
	w.WriteString(", '$.")
	w.WriteString(escaped)
	w.WriteString("') IS NOT NULL")
	return nil
}

// WriteJSONArrayElements writes BigQuery JSON array expansion.
func (d *Dialect) WriteJSONArrayElements(w *strings.Builder, _ bool, _ bool, writeExpr func() error) error {
	w.WriteString("UNNEST(JSON_QUERY_ARRAY(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString("))")
	return nil
}

// WriteJSONArrayLength writes ARRAY_LENGTH(JSON_QUERY_ARRAY(expr)) for BigQuery.
func (d *Dialect) WriteJSONArrayLength(w *strings.Builder, writeExpr func() error) error {
	w.WriteString("ARRAY_LENGTH(JSON_QUERY_ARRAY(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString("))")
	return nil
}

// WriteJSONExtractPath writes BigQuery JSON path existence using JSON_VALUE.
func (d *Dialect) WriteJSONExtractPath(w *strings.Builder, pathSegments []string, writeRoot func() error) error {
	w.WriteString("JSON_VALUE(")
	if err := writeRoot(); err != nil {
		return err
	}
	w.WriteString(", '$")
	for _, segment := range pathSegments {
		w.WriteString(".")
		w.WriteString(escapeJSONFieldName(segment))
	}
	w.WriteString("') IS NOT NULL")
	return nil
}

// WriteJSONArrayMembership writes BigQuery JSON array membership.
func (d *Dialect) WriteJSONArrayMembership(w *strings.Builder, _ string, writeExpr func() error) error {
	w.WriteString("UNNEST(JSON_VALUE_ARRAY(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString("))")
	return nil
}

// WriteNestedJSONArrayMembership writes BigQuery nested JSON array membership.
func (d *Dialect) WriteNestedJSONArrayMembership(w *strings.Builder, writeExpr func() error) error {
	w.WriteString("UNNEST(JSON_VALUE_ARRAY(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString("))")
	return nil
}

// --- Timestamps ---

// WriteDuration writes a BigQuery INTERVAL literal.
func (d *Dialect) WriteDuration(w *strings.Builder, value int64, unit string) {
	fmt.Fprintf(w, "INTERVAL %d %s", value, unit)
}

// WriteInterval writes a BigQuery INTERVAL expression.
func (d *Dialect) WriteInterval(w *strings.Builder, writeValue func() error, unit string) error {
	w.WriteString("INTERVAL ")
	if err := writeValue(); err != nil {
		return err
	}
	w.WriteString(" ")
	w.WriteString(unit)
	return nil
}

// WriteExtract writes a BigQuery EXTRACT expression.
// BigQuery uses DAYOFWEEK (1=Sunday) instead of DOW (0=Sunday).
func (d *Dialect) WriteExtract(w *strings.Builder, part string, writeExpr func() error, writeTZ func() error) error {
	isDOW := part == "DOW"
	bqPart := part
	if isDOW {
		bqPart = "DAYOFWEEK"
		w.WriteString("(")
	}
	w.WriteString("EXTRACT(")
	w.WriteString(bqPart)
	w.WriteString(" FROM ")
	if err := writeExpr(); err != nil {
		return err
	}
	if writeTZ != nil {
		w.WriteString(" AT TIME ZONE ")
		if err := writeTZ(); err != nil {
			return err
		}
	}
	w.WriteString(")")
	if isDOW {
		// BigQuery DAYOFWEEK: 1=Sunday, 2=Monday, ..., 7=Saturday
		// CEL getDayOfWeek: 0=Sunday, 1=Monday, ..., 6=Saturday
		w.WriteString(" - 1)")
	}
	return nil
}

// WriteTimestampArithmetic writes BigQuery timestamp arithmetic using functions.
// BigQuery uses TIMESTAMP_ADD/TIMESTAMP_SUB instead of + / - operators.
func (d *Dialect) WriteTimestampArithmetic(w *strings.Builder, op string, writeTS, writeDur func() error) error {
	if op == "+" {
		w.WriteString("TIMESTAMP_ADD(")
	} else {
		w.WriteString("TIMESTAMP_SUB(")
	}
	if err := writeTS(); err != nil {
		return err
	}
	w.WriteString(", ")
	if err := writeDur(); err != nil {
		return err
	}
	w.WriteString(")")
	return nil
}

// --- String Functions ---

// WriteContains writes STRPOS(haystack, needle) > 0 for BigQuery.
func (d *Dialect) WriteContains(w *strings.Builder, writeHaystack, writeNeedle func() error) error {
	w.WriteString("STRPOS(")
	if err := writeHaystack(); err != nil {
		return err
	}
	w.WriteString(", ")
	if err := writeNeedle(); err != nil {
		return err
	}
	w.WriteString(") > 0")
	return nil
}

// WriteSplit writes BigQuery string split using SPLIT.
func (d *Dialect) WriteSplit(w *strings.Builder, writeStr, writeDelim func() error) error {
	w.WriteString("SPLIT(")
	if err := writeStr(); err != nil {
		return err
	}
	w.WriteString(", ")
	if err := writeDelim(); err != nil {
		return err
	}
	w.WriteString(")")
	return nil
}

// WriteSplitWithLimit writes BigQuery SPLIT with array slice.
func (d *Dialect) WriteSplitWithLimit(w *strings.Builder, writeStr, writeDelim func() error, limit int64) error {
	w.WriteString("ARRAY(SELECT x FROM UNNEST(SPLIT(")
	if err := writeStr(); err != nil {
		return err
	}
	w.WriteString(", ")
	if err := writeDelim(); err != nil {
		return err
	}
	fmt.Fprintf(w, ")) AS x WITH OFFSET WHERE OFFSET < %d)", limit)
	return nil
}

// WriteJoin writes BigQuery array join using ARRAY_TO_STRING.
func (d *Dialect) WriteJoin(w *strings.Builder, writeArray, writeDelim func() error) error {
	w.WriteString("ARRAY_TO_STRING(")
	if err := writeArray(); err != nil {
		return err
	}
	w.WriteString(", ")
	if err := writeDelim(); err != nil {
		return err
	}
	w.WriteString(")")
	return nil
}

// --- Comprehensions ---

// WriteUnnest writes BigQuery UNNEST for array unnesting.
func (d *Dialect) WriteUnnest(w *strings.Builder, writeSource func() error) error {
	w.WriteString("UNNEST(")
	if err := writeSource(); err != nil {
		return err
	}
	w.WriteString(")")
	return nil
}

// WriteArraySubqueryOpen writes ARRAY(SELECT for BigQuery.
func (d *Dialect) WriteArraySubqueryOpen(w *strings.Builder) {
	w.WriteString("ARRAY(SELECT ")
}

// WriteArraySubqueryExprClose is a no-op for BigQuery.
func (d *Dialect) WriteArraySubqueryExprClose(_ *strings.Builder) {
}

// --- Struct ---

// WriteStructOpen writes the BigQuery struct literal opening.
func (d *Dialect) WriteStructOpen(w *strings.Builder) {
	w.WriteString("STRUCT(")
}

// WriteStructClose writes the BigQuery struct literal closing.
func (d *Dialect) WriteStructClose(w *strings.Builder) {
	w.WriteString(")")
}

// --- Validation ---

// MaxIdentifierLength returns 300 for BigQuery.
func (d *Dialect) MaxIdentifierLength() int {
	return 300
}

// ValidateFieldName validates a field name against BigQuery naming rules.
func (d *Dialect) ValidateFieldName(name string) error {
	return validateFieldName(name)
}

// ReservedKeywords returns the set of reserved SQL keywords for BigQuery.
func (d *Dialect) ReservedKeywords() map[string]bool {
	return reservedSQLKeywords
}

// --- Regex ---

// ConvertRegex converts an RE2 regex pattern for BigQuery.
// BigQuery uses RE2 natively, so minimal conversion is needed.
func (d *Dialect) ConvertRegex(re2Pattern string) (string, bool, error) {
	return convertRE2ToBigQuery(re2Pattern)
}

// SupportsRegex returns true as BigQuery supports RE2 regex natively.
func (d *Dialect) SupportsRegex() bool { return true }

// --- Capabilities ---

// SupportsNativeArrays returns true as BigQuery has native array types.
func (d *Dialect) SupportsNativeArrays() bool { return true }

// SupportsJSONB returns false as BigQuery has a single JSON type.
func (d *Dialect) SupportsJSONB() bool { return false }

// SupportsIndexAnalysis returns true as BigQuery index analysis is supported.
func (d *Dialect) SupportsIndexAnalysis() bool { return true }

// --- Internal helpers ---

// escapeJSONFieldName escapes special characters in JSON field names for BigQuery.
func escapeJSONFieldName(fieldName string) string {
	return strings.ReplaceAll(fieldName, "'", "\\'")
}

// bigqueryTypeName converts a CEL/common type name to a BigQuery type name.
func bigqueryTypeName(typeName string) string {
	switch strings.ToLower(typeName) {
	case "text", "string", "varchar":
		return "STRING"
	case "int", "integer", "bigint", "int64":
		return "INT64"
	case "double", "float", "real", "float64":
		return "FLOAT64"
	case "boolean", "bool":
		return "BOOL"
	case "bytes", "bytea", "blob":
		return "BYTES"
	default:
		return strings.ToUpper(typeName)
	}
}
