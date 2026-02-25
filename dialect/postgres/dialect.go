// Package postgres implements the PostgreSQL SQL dialect for cel2sql.
package postgres

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/spandigital/cel2sql/v3/dialect"
)

// Dialect implements dialect.Dialect for PostgreSQL.
type Dialect struct{}

// New creates a new PostgreSQL dialect.
func New() *Dialect {
	return &Dialect{}
}

func init() {
	dialect.Register(dialect.PostgreSQL, func() dialect.Dialect { return New() })
}

// Ensure Dialect implements dialect.Dialect at compile time.
var _ dialect.Dialect = (*Dialect)(nil)

// Name returns the dialect name.
func (d *Dialect) Name() dialect.Name { return dialect.PostgreSQL }

// --- Literals ---

// WriteStringLiteral writes a PostgreSQL string literal with ” escaping.
func (d *Dialect) WriteStringLiteral(w *strings.Builder, value string) {
	escaped := strings.ReplaceAll(value, "'", "''")
	w.WriteString("'")
	w.WriteString(escaped)
	w.WriteString("'")
}

// WriteBytesLiteral writes a PostgreSQL hex-encoded byte literal.
func (d *Dialect) WriteBytesLiteral(w *strings.Builder, value []byte) error {
	w.WriteString("'\\x")
	w.WriteString(hex.EncodeToString(value))
	w.WriteString("'")
	return nil
}

// WriteParamPlaceholder writes a PostgreSQL positional parameter ($1, $2, ...).
func (d *Dialect) WriteParamPlaceholder(w *strings.Builder, paramIndex int) {
	fmt.Fprintf(w, "$%d", paramIndex)
}

// --- Operators ---

// WriteStringConcat writes a PostgreSQL string concatenation using ||.
func (d *Dialect) WriteStringConcat(w *strings.Builder, writeLHS, writeRHS func() error) error {
	if err := writeLHS(); err != nil {
		return err
	}
	w.WriteString(" || ")
	return writeRHS()
}

// WriteRegexMatch writes a PostgreSQL regex match using ~ or ~* operators.
func (d *Dialect) WriteRegexMatch(w *strings.Builder, writeTarget func() error, pattern string, caseInsensitive bool) error {
	if err := writeTarget(); err != nil {
		return err
	}
	if caseInsensitive {
		w.WriteString(" ~* ")
	} else {
		w.WriteString(" ~ ")
	}
	escaped := strings.ReplaceAll(pattern, "'", "''")
	w.WriteString("'")
	w.WriteString(escaped)
	w.WriteString("'")
	return nil
}

// WriteLikeEscape writes the PostgreSQL LIKE escape clause.
func (d *Dialect) WriteLikeEscape(w *strings.Builder) {
	w.WriteString(" ESCAPE E'\\\\'")
}

// WriteArrayMembership writes a PostgreSQL array membership test using = ANY().
func (d *Dialect) WriteArrayMembership(w *strings.Builder, writeElem func() error, writeArray func() error) error {
	if err := writeElem(); err != nil {
		return err
	}
	w.WriteString(" = ANY(")
	if err := writeArray(); err != nil {
		return err
	}
	w.WriteString(")")
	return nil
}

// --- Type Casting ---

// WriteCastToNumeric writes a PostgreSQL numeric cast suffix (::numeric).
func (d *Dialect) WriteCastToNumeric(w *strings.Builder) {
	w.WriteString("::numeric")
}

// WriteTypeName writes a PostgreSQL type name for CAST expressions.
func (d *Dialect) WriteTypeName(w *strings.Builder, celTypeName string) {
	switch celTypeName {
	case "bool":
		w.WriteString("BOOLEAN")
	case "bytes":
		w.WriteString("BYTEA")
	case "double":
		w.WriteString("DOUBLE PRECISION")
	case "int":
		w.WriteString("BIGINT")
	case "string":
		w.WriteString("TEXT")
	case "uint":
		w.WriteString("BIGINT")
	default:
		w.WriteString(strings.ToUpper(celTypeName))
	}
}

// WriteEpochExtract writes EXTRACT(EPOCH FROM expr)::bigint for PostgreSQL.
func (d *Dialect) WriteEpochExtract(w *strings.Builder, writeExpr func() error) error {
	w.WriteString("EXTRACT(EPOCH FROM ")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString(")::bigint")
	return nil
}

// WriteTimestampCast writes CAST(expr AS TIMESTAMP WITH TIME ZONE) for PostgreSQL.
func (d *Dialect) WriteTimestampCast(w *strings.Builder, writeExpr func() error) error {
	w.WriteString("CAST(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString(" AS TIMESTAMP WITH TIME ZONE)")
	return nil
}

// --- Arrays ---

// WriteArrayLiteralOpen writes the PostgreSQL array literal opening (ARRAY[).
func (d *Dialect) WriteArrayLiteralOpen(w *strings.Builder) {
	w.WriteString("ARRAY[")
}

// WriteArrayLiteralClose writes the PostgreSQL array literal closing (]).
func (d *Dialect) WriteArrayLiteralClose(w *strings.Builder) {
	w.WriteString("]")
}

// WriteArrayLength writes COALESCE(ARRAY_LENGTH(expr, dimension), 0) for PostgreSQL.
func (d *Dialect) WriteArrayLength(w *strings.Builder, dimension int, writeExpr func() error) error {
	w.WriteString("COALESCE(ARRAY_LENGTH(")
	if err := writeExpr(); err != nil {
		return err
	}
	fmt.Fprintf(w, ", %d), 0)", dimension)
	return nil
}

// WriteListIndex writes a PostgreSQL 1-indexed array access (array[index + 1]).
func (d *Dialect) WriteListIndex(w *strings.Builder, writeArray func() error, writeIndex func() error) error {
	if err := writeArray(); err != nil {
		return err
	}
	w.WriteString("[")
	if err := writeIndex(); err != nil {
		return err
	}
	w.WriteString(" + 1]")
	return nil
}

// WriteListIndexConst writes a PostgreSQL constant array index (0-indexed to 1-indexed).
func (d *Dialect) WriteListIndexConst(w *strings.Builder, writeArray func() error, index int64) error {
	if err := writeArray(); err != nil {
		return err
	}
	fmt.Fprintf(w, "[%d]", index+1)
	return nil
}

// WriteEmptyTypedArray writes an empty PostgreSQL typed array (ARRAY[]::type[]).
func (d *Dialect) WriteEmptyTypedArray(w *strings.Builder, typeName string) {
	fmt.Fprintf(w, "ARRAY[]::%s[]", typeName)
}

// --- JSON ---

// WriteJSONFieldAccess writes PostgreSQL JSON field access using -> or ->> operators.
func (d *Dialect) WriteJSONFieldAccess(w *strings.Builder, writeBase func() error, fieldName string, isFinal bool) error {
	if err := writeBase(); err != nil {
		return err
	}
	escapedField := escapeJSONFieldName(fieldName)
	if isFinal {
		w.WriteString("->>'")
	} else {
		w.WriteString("->'")
	}
	w.WriteString(escapedField)
	w.WriteString("'")
	return nil
}

// WriteJSONExistence writes a PostgreSQL JSON key existence check (? or IS NOT NULL).
func (d *Dialect) WriteJSONExistence(w *strings.Builder, isJSONB bool, fieldName string, writeBase func() error) error {
	if err := writeBase(); err != nil {
		return err
	}
	escapedField := escapeJSONFieldName(fieldName)
	if isJSONB {
		w.WriteString(" ? '")
		w.WriteString(escapedField)
		w.WriteString("'")
	} else {
		w.WriteString("->'")
		w.WriteString(escapedField)
		w.WriteString("' IS NOT NULL")
	}
	return nil
}

// WriteJSONArrayElements writes a PostgreSQL JSON array expansion function.
func (d *Dialect) WriteJSONArrayElements(w *strings.Builder, isJSONB bool, asText bool, writeExpr func() error) error {
	if isJSONB {
		if asText {
			w.WriteString("jsonb_array_elements_text(")
		} else {
			w.WriteString("jsonb_array_elements(")
		}
	} else {
		if asText {
			w.WriteString("json_array_elements_text(")
		} else {
			w.WriteString("json_array_elements(")
		}
	}
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString(")")
	return nil
}

// WriteJSONArrayLength writes COALESCE(jsonb_array_length(expr), 0) for PostgreSQL.
func (d *Dialect) WriteJSONArrayLength(w *strings.Builder, writeExpr func() error) error {
	w.WriteString("COALESCE(jsonb_array_length(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString("), 0)")
	return nil
}

// WriteJSONExtractPath writes jsonb_extract_path_text() IS NOT NULL for PostgreSQL.
func (d *Dialect) WriteJSONExtractPath(w *strings.Builder, pathSegments []string, writeRoot func() error) error {
	w.WriteString("jsonb_extract_path_text(")
	if err := writeRoot(); err != nil {
		return err
	}
	for _, segment := range pathSegments {
		w.WriteString(", '")
		w.WriteString(escapeJSONFieldName(segment))
		w.WriteString("'")
	}
	w.WriteString(") IS NOT NULL")
	return nil
}

// WriteJSONArrayMembership writes ANY(ARRAY(SELECT json_func(expr))) for PostgreSQL.
func (d *Dialect) WriteJSONArrayMembership(w *strings.Builder, jsonFunc string, writeExpr func() error) error {
	w.WriteString("ANY(ARRAY(SELECT ")
	w.WriteString(jsonFunc)
	w.WriteString("(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString(")))")
	return nil
}

// WriteNestedJSONArrayMembership writes ANY(ARRAY(SELECT jsonb_array_elements_text(expr))) for PostgreSQL.
func (d *Dialect) WriteNestedJSONArrayMembership(w *strings.Builder, writeExpr func() error) error {
	w.WriteString("ANY(ARRAY(SELECT jsonb_array_elements_text(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString(")))")
	return nil
}

// --- Timestamps ---

// WriteDuration writes a PostgreSQL INTERVAL literal (INTERVAL N UNIT).
func (d *Dialect) WriteDuration(w *strings.Builder, value int64, unit string) {
	fmt.Fprintf(w, "INTERVAL %d %s", value, unit)
}

// WriteInterval writes a PostgreSQL INTERVAL expression (INTERVAL expr UNIT).
func (d *Dialect) WriteInterval(w *strings.Builder, writeValue func() error, unit string) error {
	w.WriteString("INTERVAL ")
	if err := writeValue(); err != nil {
		return err
	}
	w.WriteString(" ")
	w.WriteString(unit)
	return nil
}

// WriteExtract writes a PostgreSQL EXTRACT expression with DOW conversion.
func (d *Dialect) WriteExtract(w *strings.Builder, part string, writeExpr func() error, writeTZ func() error) error {
	// For getDayOfWeek, we need to wrap the entire EXTRACT for modulo operation
	isDOW := part == "DOW"
	if isDOW {
		w.WriteString("(")
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
	if isDOW {
		// PostgreSQL DOW: 0=Sunday, 1=Monday, ..., 6=Saturday
		// CEL getDayOfWeek: 0=Monday, 1=Tuesday, ..., 6=Sunday (ISO 8601)
		// Convert: (DOW + 6) % 7
		w.WriteString(" + 6) % 7")
	}
	return nil
}

// WriteTimestampArithmetic writes PostgreSQL timestamp arithmetic (timestamp +/- interval).
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

// WriteContains writes POSITION(needle IN haystack) > 0 for PostgreSQL.
func (d *Dialect) WriteContains(w *strings.Builder, writeHaystack, writeNeedle func() error) error {
	w.WriteString("POSITION(")
	if err := writeNeedle(); err != nil {
		return err
	}
	w.WriteString(" IN ")
	if err := writeHaystack(); err != nil {
		return err
	}
	w.WriteString(") > 0")
	return nil
}

// WriteSplit writes STRING_TO_ARRAY(string, delimiter) for PostgreSQL.
func (d *Dialect) WriteSplit(w *strings.Builder, writeStr, writeDelim func() error) error {
	w.WriteString("STRING_TO_ARRAY(")
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

// WriteSplitWithLimit writes (STRING_TO_ARRAY(string, delimiter))[1:limit] for PostgreSQL.
func (d *Dialect) WriteSplitWithLimit(w *strings.Builder, writeStr, writeDelim func() error, limit int64) error {
	w.WriteString("(STRING_TO_ARRAY(")
	if err := writeStr(); err != nil {
		return err
	}
	w.WriteString(", ")
	if err := writeDelim(); err != nil {
		return err
	}
	fmt.Fprintf(w, "))[1:%d]", limit)
	return nil
}

// WriteJoin writes ARRAY_TO_STRING(array, delimiter, ”) for PostgreSQL.
func (d *Dialect) WriteJoin(w *strings.Builder, writeArray, writeDelim func() error) error {
	w.WriteString("ARRAY_TO_STRING(")
	if err := writeArray(); err != nil {
		return err
	}
	w.WriteString(", ")
	if writeDelim != nil {
		if err := writeDelim(); err != nil {
			return err
		}
	} else {
		w.WriteString("''")
	}
	w.WriteString(", '')")
	return nil
}

// --- Comprehensions ---

// WriteUnnest writes UNNEST(source) for PostgreSQL comprehensions.
func (d *Dialect) WriteUnnest(w *strings.Builder, writeSource func() error) error {
	w.WriteString("UNNEST(")
	if err := writeSource(); err != nil {
		return err
	}
	w.WriteString(")")
	return nil
}

// WriteArraySubqueryOpen writes ARRAY(SELECT for PostgreSQL.
func (d *Dialect) WriteArraySubqueryOpen(w *strings.Builder) {
	w.WriteString("ARRAY(SELECT ")
}

// WriteArraySubqueryExprClose is a no-op for PostgreSQL (no wrapper around the expression).
func (d *Dialect) WriteArraySubqueryExprClose(_ *strings.Builder) {
}

// --- Struct ---

// WriteStructOpen writes the PostgreSQL struct/row literal opening (ROW().
func (d *Dialect) WriteStructOpen(w *strings.Builder) {
	w.WriteString("ROW(")
}

// WriteStructClose writes the PostgreSQL struct/row literal closing ()).
func (d *Dialect) WriteStructClose(w *strings.Builder) {
	w.WriteString(")")
}

// --- Validation ---

// MaxIdentifierLength returns the PostgreSQL maximum identifier length (63).
func (d *Dialect) MaxIdentifierLength() int {
	return maxPostgreSQLIdentifierLength
}

// ValidateFieldName validates a field name against PostgreSQL naming rules.
func (d *Dialect) ValidateFieldName(name string) error {
	return validateFieldName(name)
}

// ReservedKeywords returns the set of reserved SQL keywords for PostgreSQL.
func (d *Dialect) ReservedKeywords() map[string]bool {
	return reservedSQLKeywords
}

// --- Regex ---

// ConvertRegex converts an RE2 regex pattern to PostgreSQL POSIX format.
func (d *Dialect) ConvertRegex(re2Pattern string) (string, bool, error) {
	return convertRE2ToPOSIX(re2Pattern)
}

// SupportsRegex returns true as PostgreSQL supports POSIX regex matching.
func (d *Dialect) SupportsRegex() bool { return true }

// --- Capabilities ---

// SupportsNativeArrays returns true as PostgreSQL has native array types.
func (d *Dialect) SupportsNativeArrays() bool { return true }

// SupportsJSONB returns true as PostgreSQL has a distinct JSONB type.
func (d *Dialect) SupportsJSONB() bool { return true }

// SupportsIndexAnalysis returns true as PostgreSQL index analysis is supported.
func (d *Dialect) SupportsIndexAnalysis() bool { return true }

// --- Internal helpers ---

// escapeJSONFieldName escapes single quotes in JSON field names for safe use in PostgreSQL JSON path operators.
func escapeJSONFieldName(fieldName string) string {
	return strings.ReplaceAll(fieldName, "'", "''")
}
