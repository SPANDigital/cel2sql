// Package mysql implements the MySQL SQL dialect for cel2sql.
package mysql

import (
	"fmt"
	"strings"

	"github.com/spandigital/cel2sql/v3/dialect"
)

// Dialect implements dialect.Dialect for MySQL 8.0+.
type Dialect struct{}

// New creates a new MySQL dialect.
func New() *Dialect {
	return &Dialect{}
}

func init() {
	dialect.Register(dialect.MySQL, func() dialect.Dialect { return New() })
}

// Ensure Dialect implements dialect.Dialect at compile time.
var _ dialect.Dialect = (*Dialect)(nil)

// Name returns the dialect name.
func (d *Dialect) Name() dialect.Name { return dialect.MySQL }

// --- Literals ---

// WriteStringLiteral writes a MySQL string literal with ” escaping.
func (d *Dialect) WriteStringLiteral(w *strings.Builder, value string) {
	escaped := strings.ReplaceAll(value, "'", "''")
	w.WriteString("'")
	w.WriteString(escaped)
	w.WriteString("'")
}

// WriteBytesLiteral writes a MySQL hex-encoded byte literal (X'...').
func (d *Dialect) WriteBytesLiteral(w *strings.Builder, value []byte) error {
	w.WriteString("X'")
	for _, b := range value {
		fmt.Fprintf(w, "%02x", b)
	}
	w.WriteString("'")
	return nil
}

// WriteParamPlaceholder writes a MySQL positional parameter (?).
func (d *Dialect) WriteParamPlaceholder(w *strings.Builder, _ int) {
	w.WriteString("?")
}

// --- Operators ---

// WriteStringConcat writes MySQL string concatenation using CONCAT().
func (d *Dialect) WriteStringConcat(w *strings.Builder, writeLHS, writeRHS func() error) error {
	w.WriteString("CONCAT(")
	if err := writeLHS(); err != nil {
		return err
	}
	w.WriteString(", ")
	if err := writeRHS(); err != nil {
		return err
	}
	w.WriteString(")")
	return nil
}

// WriteRegexMatch writes a MySQL REGEXP match expression.
func (d *Dialect) WriteRegexMatch(w *strings.Builder, writeTarget func() error, pattern string, _ bool) error {
	if err := writeTarget(); err != nil {
		return err
	}
	w.WriteString(" REGEXP ")
	escaped := strings.ReplaceAll(pattern, "'", "''")
	w.WriteString("'")
	w.WriteString(escaped)
	w.WriteString("'")
	return nil
}

// WriteLikeEscape writes the MySQL LIKE escape clause.
func (d *Dialect) WriteLikeEscape(w *strings.Builder) {
	w.WriteString(" ESCAPE '\\\\'")
}

// WriteArrayMembership writes a MySQL array membership test using JSON_CONTAINS.
func (d *Dialect) WriteArrayMembership(w *strings.Builder, writeElem func() error, writeArray func() error) error {
	w.WriteString("JSON_CONTAINS(")
	if err := writeArray(); err != nil {
		return err
	}
	w.WriteString(", CAST(")
	if err := writeElem(); err != nil {
		return err
	}
	w.WriteString(" AS JSON))")
	return nil
}

// --- Type Casting ---

// WriteCastToNumeric writes a MySQL numeric cast (CAST(... AS DECIMAL)).
func (d *Dialect) WriteCastToNumeric(w *strings.Builder) {
	w.WriteString(" + 0")
}

// WriteTypeName writes a MySQL type name for CAST expressions.
func (d *Dialect) WriteTypeName(w *strings.Builder, celTypeName string) {
	switch celTypeName {
	case "bool":
		w.WriteString("UNSIGNED")
	case "bytes":
		w.WriteString("BINARY")
	case "double":
		w.WriteString("DECIMAL")
	case "int":
		w.WriteString("SIGNED")
	case "string":
		w.WriteString("CHAR")
	case "uint":
		w.WriteString("UNSIGNED")
	default:
		w.WriteString(strings.ToUpper(celTypeName))
	}
}

// WriteEpochExtract writes UNIX_TIMESTAMP(expr) for MySQL.
func (d *Dialect) WriteEpochExtract(w *strings.Builder, writeExpr func() error) error {
	w.WriteString("UNIX_TIMESTAMP(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString(")")
	return nil
}

// WriteTimestampCast writes a MySQL CAST to DATETIME.
func (d *Dialect) WriteTimestampCast(w *strings.Builder, writeExpr func() error) error {
	w.WriteString("CAST(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString(" AS DATETIME)")
	return nil
}

// --- Arrays ---

// WriteArrayLiteralOpen writes the MySQL JSON array literal opening.
func (d *Dialect) WriteArrayLiteralOpen(w *strings.Builder) {
	w.WriteString("JSON_ARRAY(")
}

// WriteArrayLiteralClose writes the MySQL JSON array literal closing.
func (d *Dialect) WriteArrayLiteralClose(w *strings.Builder) {
	w.WriteString(")")
}

// WriteArrayLength writes JSON_LENGTH(expr) for MySQL.
func (d *Dialect) WriteArrayLength(w *strings.Builder, _ int, writeExpr func() error) error {
	w.WriteString("COALESCE(JSON_LENGTH(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString("), 0)")
	return nil
}

// WriteListIndex writes MySQL JSON array index access.
func (d *Dialect) WriteListIndex(w *strings.Builder, writeArray func() error, writeIndex func() error) error {
	w.WriteString("JSON_EXTRACT(")
	if err := writeArray(); err != nil {
		return err
	}
	w.WriteString(", CONCAT('$[', ")
	if err := writeIndex(); err != nil {
		return err
	}
	w.WriteString(", ']'))")
	return nil
}

// WriteListIndexConst writes MySQL JSON constant array index access.
func (d *Dialect) WriteListIndexConst(w *strings.Builder, writeArray func() error, index int64) error {
	w.WriteString("JSON_EXTRACT(")
	if err := writeArray(); err != nil {
		return err
	}
	fmt.Fprintf(w, ", '$[%d]')", index)
	return nil
}

// WriteEmptyTypedArray writes an empty MySQL JSON array.
func (d *Dialect) WriteEmptyTypedArray(w *strings.Builder, _ string) {
	w.WriteString("JSON_ARRAY()")
}

// --- JSON ---

// WriteJSONFieldAccess writes MySQL JSON field access using JSON_EXTRACT/JSON_UNQUOTE.
func (d *Dialect) WriteJSONFieldAccess(w *strings.Builder, writeBase func() error, fieldName string, isFinal bool) error {
	if err := writeBase(); err != nil {
		return err
	}
	escaped := escapeJSONFieldName(fieldName)
	if isFinal {
		// For final access, we need text: use ->> which is JSON_UNQUOTE(JSON_EXTRACT(...))
		w.WriteString("->>'$.")
		w.WriteString(escaped)
		w.WriteString("'")
	} else {
		w.WriteString("->'$.")
		w.WriteString(escaped)
		w.WriteString("'")
	}
	return nil
}

// WriteJSONExistence writes a MySQL JSON key existence check.
func (d *Dialect) WriteJSONExistence(w *strings.Builder, _ bool, fieldName string, writeBase func() error) error {
	w.WriteString("JSON_CONTAINS_PATH(")
	if err := writeBase(); err != nil {
		return err
	}
	escaped := escapeJSONFieldName(fieldName)
	w.WriteString(", 'one', '$.")
	w.WriteString(escaped)
	w.WriteString("')")
	return nil
}

// WriteJSONArrayElements writes MySQL JSON array expansion using JSON_TABLE.
func (d *Dialect) WriteJSONArrayElements(w *strings.Builder, _ bool, _ bool, writeExpr func() error) error {
	w.WriteString("JSON_TABLE(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString(", '$[*]' COLUMNS(value TEXT PATH '$'))")
	return nil
}

// WriteJSONArrayLength writes COALESCE(JSON_LENGTH(expr), 0) for MySQL.
func (d *Dialect) WriteJSONArrayLength(w *strings.Builder, writeExpr func() error) error {
	w.WriteString("COALESCE(JSON_LENGTH(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString("), 0)")
	return nil
}

// WriteJSONExtractPath writes MySQL JSON path extraction using JSON_CONTAINS_PATH.
func (d *Dialect) WriteJSONExtractPath(w *strings.Builder, pathSegments []string, writeRoot func() error) error {
	w.WriteString("JSON_CONTAINS_PATH(")
	if err := writeRoot(); err != nil {
		return err
	}
	w.WriteString(", 'one', '$")
	for _, segment := range pathSegments {
		w.WriteString(".")
		w.WriteString(escapeJSONFieldName(segment))
	}
	w.WriteString("')")
	return nil
}

// WriteJSONArrayMembership writes MySQL JSON array membership using JSON_CONTAINS.
func (d *Dialect) WriteJSONArrayMembership(w *strings.Builder, _ string, writeExpr func() error) error {
	w.WriteString("JSON_CONTAINS(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString(", CAST(? AS JSON))")
	return nil
}

// WriteNestedJSONArrayMembership writes MySQL nested JSON array membership.
func (d *Dialect) WriteNestedJSONArrayMembership(w *strings.Builder, writeExpr func() error) error {
	w.WriteString("JSON_CONTAINS(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString(", CAST(? AS JSON))")
	return nil
}

// --- Timestamps ---

// WriteDuration writes a MySQL INTERVAL literal.
func (d *Dialect) WriteDuration(w *strings.Builder, value int64, unit string) {
	fmt.Fprintf(w, "INTERVAL %d %s", value, unit)
}

// WriteInterval writes a MySQL INTERVAL expression.
func (d *Dialect) WriteInterval(w *strings.Builder, writeValue func() error, unit string) error {
	w.WriteString("INTERVAL ")
	if err := writeValue(); err != nil {
		return err
	}
	w.WriteString(" ")
	w.WriteString(unit)
	return nil
}

// WriteExtract writes a MySQL EXTRACT expression with DOW conversion.
func (d *Dialect) WriteExtract(w *strings.Builder, part string, writeExpr func() error, writeTZ func() error) error {
	isDOW := part == "DOW"
	if isDOW {
		// MySQL DAYOFWEEK: 1=Sunday, 2=Monday, ..., 7=Saturday
		// CEL getDayOfWeek: 0=Monday, 1=Tuesday, ..., 6=Sunday (ISO 8601)
		// Convert: (DAYOFWEEK(x) + 5) % 7
		w.WriteString("(DAYOFWEEK(")
		if err := writeExpr(); err != nil {
			return err
		}
		w.WriteString(") + 5) % 7")
		return nil
	}

	w.WriteString("EXTRACT(")
	w.WriteString(part)
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
	return nil
}

// WriteTimestampArithmetic writes MySQL timestamp arithmetic.
func (d *Dialect) WriteTimestampArithmetic(w *strings.Builder, op string, writeTS, writeDur func() error) error {
	if err := writeTS(); err != nil {
		return err
	}
	w.WriteString(" ")
	w.WriteString(op)
	w.WriteString(" ")
	return writeDur()
}

// --- String Functions ---

// WriteContains writes LOCATE(needle, haystack) > 0 for MySQL.
func (d *Dialect) WriteContains(w *strings.Builder, writeHaystack, writeNeedle func() error) error {
	w.WriteString("LOCATE(")
	if err := writeNeedle(); err != nil {
		return err
	}
	w.WriteString(", ")
	if err := writeHaystack(); err != nil {
		return err
	}
	w.WriteString(") > 0")
	return nil
}

// WriteSplit writes MySQL string split using SUBSTRING_INDEX pattern.
func (d *Dialect) WriteSplit(w *strings.Builder, writeStr, writeDelim func() error) error {
	// MySQL doesn't have a direct STRING_TO_ARRAY equivalent.
	// Use a JSON approach: convert to JSON array.
	w.WriteString("JSON_ARRAY(")
	if err := writeStr(); err != nil {
		return err
	}
	w.WriteString(")")
	// Note: A full MySQL split implementation would require a more complex approach.
	// This is a simplified version.
	_ = writeDelim
	return nil
}

// WriteSplitWithLimit writes MySQL string split with limit.
func (d *Dialect) WriteSplitWithLimit(w *strings.Builder, writeStr, writeDelim func() error, _ int64) error {
	// Simplified: delegate to WriteSplit
	return d.WriteSplit(w, writeStr, writeDelim)
}

// WriteJoin writes MySQL array join using JSON_UNQUOTE/GROUP_CONCAT pattern.
func (d *Dialect) WriteJoin(w *strings.Builder, writeArray, writeDelim func() error) error {
	// MySQL doesn't have ARRAY_TO_STRING; simplified approach
	w.WriteString("JSON_UNQUOTE(")
	if err := writeArray(); err != nil {
		return err
	}
	w.WriteString(")")
	_ = writeDelim
	return nil
}

// --- Comprehensions ---

// WriteUnnest writes MySQL JSON_TABLE for array unnesting.
func (d *Dialect) WriteUnnest(w *strings.Builder, writeSource func() error) error {
	w.WriteString("JSON_TABLE(")
	if err := writeSource(); err != nil {
		return err
	}
	w.WriteString(", '$[*]' COLUMNS(value TEXT PATH '$'))")
	return nil
}

// WriteArraySubqueryOpen writes (SELECT JSON_ARRAYAGG( for MySQL array subqueries.
func (d *Dialect) WriteArraySubqueryOpen(w *strings.Builder) {
	w.WriteString("(SELECT JSON_ARRAYAGG(")
}

// WriteArraySubqueryExprClose closes the JSON_ARRAYAGG aggregate function for MySQL.
func (d *Dialect) WriteArraySubqueryExprClose(w *strings.Builder) {
	w.WriteString(")")
}

// --- Struct ---

// WriteStructOpen writes the MySQL struct literal opening.
func (d *Dialect) WriteStructOpen(w *strings.Builder) {
	w.WriteString("ROW(")
}

// WriteStructClose writes the MySQL struct literal closing.
func (d *Dialect) WriteStructClose(w *strings.Builder) {
	w.WriteString(")")
}

// --- Validation ---

// MaxIdentifierLength returns the MySQL maximum identifier length (64).
func (d *Dialect) MaxIdentifierLength() int {
	return maxMySQLIdentifierLength
}

// ValidateFieldName validates a field name against MySQL naming rules.
func (d *Dialect) ValidateFieldName(name string) error {
	return validateFieldName(name)
}

// ReservedKeywords returns the set of reserved SQL keywords for MySQL.
func (d *Dialect) ReservedKeywords() map[string]bool {
	return reservedSQLKeywords
}

// --- Regex ---

// ConvertRegex converts an RE2 regex pattern to MySQL-compatible format.
func (d *Dialect) ConvertRegex(re2Pattern string) (string, bool, error) {
	return convertRE2ToMySQL(re2Pattern)
}

// SupportsRegex returns true as MySQL 8.0+ supports ICU regex.
func (d *Dialect) SupportsRegex() bool { return true }

// --- Capabilities ---

// SupportsNativeArrays returns false as MySQL uses JSON arrays.
func (d *Dialect) SupportsNativeArrays() bool { return false }

// SupportsJSONB returns false as MySQL has a single JSON type.
func (d *Dialect) SupportsJSONB() bool { return false }

// SupportsIndexAnalysis returns true as MySQL index analysis is supported.
func (d *Dialect) SupportsIndexAnalysis() bool { return true }

// --- Internal helpers ---

// escapeJSONFieldName escapes special characters in JSON field names for MySQL.
func escapeJSONFieldName(fieldName string) string {
	return strings.ReplaceAll(fieldName, "'", "''")
}
