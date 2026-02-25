// Package duckdb implements the DuckDB SQL dialect for cel2sql.
package duckdb

import (
	"fmt"
	"strings"

	"github.com/spandigital/cel2sql/v3/dialect"
)

// Dialect implements dialect.Dialect for DuckDB.
type Dialect struct{}

// New creates a new DuckDB dialect.
func New() *Dialect {
	return &Dialect{}
}

func init() {
	dialect.Register(dialect.DuckDB, func() dialect.Dialect { return New() })
}

// Ensure Dialect implements dialect.Dialect at compile time.
var _ dialect.Dialect = (*Dialect)(nil)

// Name returns the dialect name.
func (d *Dialect) Name() dialect.Name { return dialect.DuckDB }

// --- Literals ---

// WriteStringLiteral writes a DuckDB string literal with ” escaping.
func (d *Dialect) WriteStringLiteral(w *strings.Builder, value string) {
	escaped := strings.ReplaceAll(value, "'", "''")
	w.WriteString("'")
	w.WriteString(escaped)
	w.WriteString("'")
}

// WriteBytesLiteral writes a DuckDB hex-encoded byte literal ('\x...').
func (d *Dialect) WriteBytesLiteral(w *strings.Builder, value []byte) error {
	w.WriteString("'\\x")
	for _, b := range value {
		fmt.Fprintf(w, "%02x", b)
	}
	w.WriteString("'")
	return nil
}

// WriteParamPlaceholder writes a DuckDB positional parameter ($1, $2, ...).
func (d *Dialect) WriteParamPlaceholder(w *strings.Builder, paramIndex int) {
	fmt.Fprintf(w, "$%d", paramIndex)
}

// --- Operators ---

// WriteStringConcat writes DuckDB string concatenation using ||.
func (d *Dialect) WriteStringConcat(w *strings.Builder, writeLHS, writeRHS func() error) error {
	if err := writeLHS(); err != nil {
		return err
	}
	w.WriteString(" || ")
	return writeRHS()
}

// WriteRegexMatch writes a DuckDB regex match expression using ~ or ~*.
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

// WriteLikeEscape writes the DuckDB LIKE escape clause.
func (d *Dialect) WriteLikeEscape(w *strings.Builder) {
	w.WriteString(" ESCAPE '\\\\'")
}

// WriteArrayMembership writes a DuckDB array membership test using = ANY().
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

// WriteCastToNumeric writes a DuckDB numeric cast (::DOUBLE).
func (d *Dialect) WriteCastToNumeric(w *strings.Builder) {
	w.WriteString("::DOUBLE")
}

// WriteTypeName writes a DuckDB type name for CAST expressions.
func (d *Dialect) WriteTypeName(w *strings.Builder, celTypeName string) {
	switch celTypeName {
	case "bool":
		w.WriteString("BOOLEAN")
	case "bytes":
		w.WriteString("BLOB")
	case "double":
		w.WriteString("DOUBLE")
	case "int":
		w.WriteString("BIGINT")
	case "string":
		w.WriteString("VARCHAR")
	case "uint":
		w.WriteString("UBIGINT")
	default:
		w.WriteString(strings.ToUpper(celTypeName))
	}
}

// WriteEpochExtract writes EXTRACT(EPOCH FROM expr)::BIGINT for DuckDB.
func (d *Dialect) WriteEpochExtract(w *strings.Builder, writeExpr func() error) error {
	w.WriteString("EXTRACT(EPOCH FROM ")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString(")::BIGINT")
	return nil
}

// WriteTimestampCast writes a DuckDB CAST to TIMESTAMPTZ.
func (d *Dialect) WriteTimestampCast(w *strings.Builder, writeExpr func() error) error {
	w.WriteString("CAST(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString(" AS TIMESTAMPTZ)")
	return nil
}

// --- Arrays ---

// WriteArrayLiteralOpen writes the DuckDB array literal opening ([).
func (d *Dialect) WriteArrayLiteralOpen(w *strings.Builder) {
	w.WriteString("[")
}

// WriteArrayLiteralClose writes the DuckDB array literal closing (]).
func (d *Dialect) WriteArrayLiteralClose(w *strings.Builder) {
	w.WriteString("]")
}

// WriteArrayLength writes COALESCE(array_length(expr), 0) for DuckDB.
func (d *Dialect) WriteArrayLength(w *strings.Builder, _ int, writeExpr func() error) error {
	w.WriteString("COALESCE(array_length(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString("), 0)")
	return nil
}

// WriteListIndex writes DuckDB 1-indexed array access.
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

// WriteListIndexConst writes DuckDB constant array index access (0-indexed to 1-indexed).
func (d *Dialect) WriteListIndexConst(w *strings.Builder, writeArray func() error, index int64) error {
	if err := writeArray(); err != nil {
		return err
	}
	fmt.Fprintf(w, "[%d]", index+1)
	return nil
}

// WriteEmptyTypedArray writes an empty DuckDB typed array.
func (d *Dialect) WriteEmptyTypedArray(w *strings.Builder, typeName string) {
	w.WriteString("[]::") //nolint:gocritic
	w.WriteString(typeName)
	w.WriteString("[]")
}

// --- JSON ---

// WriteJSONFieldAccess writes DuckDB JSON field access using -> or ->> operators.
func (d *Dialect) WriteJSONFieldAccess(w *strings.Builder, writeBase func() error, fieldName string, isFinal bool) error {
	if err := writeBase(); err != nil {
		return err
	}
	escaped := escapeJSONFieldName(fieldName)
	if isFinal {
		w.WriteString("->>'")
	} else {
		w.WriteString("->'")
	}
	w.WriteString(escaped)
	w.WriteString("'")
	return nil
}

// WriteJSONExistence writes a DuckDB JSON key existence check using json_exists.
func (d *Dialect) WriteJSONExistence(w *strings.Builder, _ bool, fieldName string, writeBase func() error) error {
	w.WriteString("json_exists(")
	if err := writeBase(); err != nil {
		return err
	}
	escaped := escapeJSONFieldName(fieldName)
	w.WriteString(", '$.")
	w.WriteString(escaped)
	w.WriteString("')")
	return nil
}

// WriteJSONArrayElements writes DuckDB JSON array expansion using json_each.
func (d *Dialect) WriteJSONArrayElements(w *strings.Builder, _ bool, _ bool, writeExpr func() error) error {
	w.WriteString("json_each(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString(")")
	return nil
}

// WriteJSONArrayLength writes COALESCE(json_array_length(expr), 0) for DuckDB.
func (d *Dialect) WriteJSONArrayLength(w *strings.Builder, writeExpr func() error) error {
	w.WriteString("COALESCE(json_array_length(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString("), 0)")
	return nil
}

// WriteJSONExtractPath writes DuckDB JSON path existence using json_exists.
func (d *Dialect) WriteJSONExtractPath(w *strings.Builder, pathSegments []string, writeRoot func() error) error {
	w.WriteString("json_exists(")
	if err := writeRoot(); err != nil {
		return err
	}
	w.WriteString(", '$")
	for _, segment := range pathSegments {
		w.WriteString(".")
		w.WriteString(escapeJSONFieldName(segment))
	}
	w.WriteString("')")
	return nil
}

// WriteJSONArrayMembership writes DuckDB JSON array membership using json_each.
func (d *Dialect) WriteJSONArrayMembership(w *strings.Builder, _ string, writeExpr func() error) error {
	w.WriteString("(SELECT value FROM json_each(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString("))")
	return nil
}

// WriteNestedJSONArrayMembership writes DuckDB nested JSON array membership.
func (d *Dialect) WriteNestedJSONArrayMembership(w *strings.Builder, writeExpr func() error) error {
	w.WriteString("(SELECT value FROM json_each(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString("))")
	return nil
}

// --- Timestamps ---

// WriteDuration writes a DuckDB INTERVAL literal.
func (d *Dialect) WriteDuration(w *strings.Builder, value int64, unit string) {
	fmt.Fprintf(w, "INTERVAL %d %s", value, unit)
}

// WriteInterval writes a DuckDB INTERVAL expression.
func (d *Dialect) WriteInterval(w *strings.Builder, writeValue func() error, unit string) error {
	w.WriteString("INTERVAL ")
	if err := writeValue(); err != nil {
		return err
	}
	w.WriteString(" ")
	w.WriteString(unit)
	return nil
}

// WriteExtract writes a DuckDB EXTRACT expression with DOW conversion.
func (d *Dialect) WriteExtract(w *strings.Builder, part string, writeExpr func() error, writeTZ func() error) error {
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
		w.WriteString(" + 6) % 7")
	}
	return nil
}

// WriteTimestampArithmetic writes DuckDB timestamp arithmetic.
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

// WriteContains writes CONTAINS(haystack, needle) for DuckDB.
func (d *Dialect) WriteContains(w *strings.Builder, writeHaystack, writeNeedle func() error) error {
	w.WriteString("CONTAINS(")
	if err := writeHaystack(); err != nil {
		return err
	}
	w.WriteString(", ")
	if err := writeNeedle(); err != nil {
		return err
	}
	w.WriteString(")")
	return nil
}

// WriteSplit writes DuckDB string split using STRING_SPLIT.
func (d *Dialect) WriteSplit(w *strings.Builder, writeStr, writeDelim func() error) error {
	w.WriteString("STRING_SPLIT(")
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

// WriteSplitWithLimit writes DuckDB string split with array slice.
func (d *Dialect) WriteSplitWithLimit(w *strings.Builder, writeStr, writeDelim func() error, limit int64) error {
	w.WriteString("STRING_SPLIT(")
	if err := writeStr(); err != nil {
		return err
	}
	w.WriteString(", ")
	if err := writeDelim(); err != nil {
		return err
	}
	fmt.Fprintf(w, ")[1:%d]", limit)
	return nil
}

// WriteJoin writes DuckDB array join using ARRAY_TO_STRING.
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

// WriteUnnest writes DuckDB UNNEST for array unnesting.
func (d *Dialect) WriteUnnest(w *strings.Builder, writeSource func() error) error {
	w.WriteString("UNNEST(")
	if err := writeSource(); err != nil {
		return err
	}
	w.WriteString(")")
	return nil
}

// WriteArraySubqueryOpen writes ARRAY(SELECT for DuckDB.
func (d *Dialect) WriteArraySubqueryOpen(w *strings.Builder) {
	w.WriteString("ARRAY(SELECT ")
}

// WriteArraySubqueryExprClose is a no-op for DuckDB (no wrapper around the expression).
func (d *Dialect) WriteArraySubqueryExprClose(_ *strings.Builder) {
}

// --- Struct ---

// WriteStructOpen writes the DuckDB struct literal opening.
func (d *Dialect) WriteStructOpen(w *strings.Builder) {
	w.WriteString("ROW(")
}

// WriteStructClose writes the DuckDB struct literal closing.
func (d *Dialect) WriteStructClose(w *strings.Builder) {
	w.WriteString(")")
}

// --- Validation ---

// MaxIdentifierLength returns 0 as DuckDB has no hard identifier length limit.
func (d *Dialect) MaxIdentifierLength() int {
	return 0
}

// ValidateFieldName validates a field name against DuckDB naming rules.
func (d *Dialect) ValidateFieldName(name string) error {
	return validateFieldName(name)
}

// ReservedKeywords returns the set of reserved SQL keywords for DuckDB.
func (d *Dialect) ReservedKeywords() map[string]bool {
	return reservedSQLKeywords
}

// --- Regex ---

// ConvertRegex converts an RE2 regex pattern to DuckDB-compatible format.
// DuckDB uses RE2 natively, so minimal conversion is needed.
func (d *Dialect) ConvertRegex(re2Pattern string) (string, bool, error) {
	return convertRE2ToDuckDB(re2Pattern)
}

// SupportsRegex returns true as DuckDB supports RE2 regex natively.
func (d *Dialect) SupportsRegex() bool { return true }

// --- Capabilities ---

// SupportsNativeArrays returns true as DuckDB has native array (LIST) types.
func (d *Dialect) SupportsNativeArrays() bool { return true }

// SupportsJSONB returns false as DuckDB has a single JSON type.
func (d *Dialect) SupportsJSONB() bool { return false }

// SupportsIndexAnalysis returns true as DuckDB index analysis is supported.
func (d *Dialect) SupportsIndexAnalysis() bool { return true }

// --- Internal helpers ---

// escapeJSONFieldName escapes special characters in JSON field names for DuckDB.
func escapeJSONFieldName(fieldName string) string {
	return strings.ReplaceAll(fieldName, "'", "''")
}
