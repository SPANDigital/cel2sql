// Package sqlite implements the SQLite SQL dialect for cel2sql.
package sqlite

import (
	"fmt"
	"strings"

	"github.com/spandigital/cel2sql/v3/dialect"
)

// Dialect implements dialect.Dialect for SQLite.
type Dialect struct{}

// New creates a new SQLite dialect.
func New() *Dialect {
	return &Dialect{}
}

func init() {
	dialect.Register(dialect.SQLite, func() dialect.Dialect { return New() })
}

// Ensure Dialect implements dialect.Dialect at compile time.
var _ dialect.Dialect = (*Dialect)(nil)

// Name returns the dialect name.
func (d *Dialect) Name() dialect.Name { return dialect.SQLite }

// --- Literals ---

// WriteStringLiteral writes a SQLite string literal with ” escaping.
func (d *Dialect) WriteStringLiteral(w *strings.Builder, value string) {
	escaped := strings.ReplaceAll(value, "'", "''")
	w.WriteString("'")
	w.WriteString(escaped)
	w.WriteString("'")
}

// WriteBytesLiteral writes a SQLite hex-encoded byte literal (X'...').
func (d *Dialect) WriteBytesLiteral(w *strings.Builder, value []byte) error {
	w.WriteString("X'")
	for _, b := range value {
		fmt.Fprintf(w, "%02x", b)
	}
	w.WriteString("'")
	return nil
}

// WriteParamPlaceholder writes a SQLite positional parameter (?).
func (d *Dialect) WriteParamPlaceholder(w *strings.Builder, _ int) {
	w.WriteString("?")
}

// --- Operators ---

// WriteStringConcat writes SQLite string concatenation using ||.
func (d *Dialect) WriteStringConcat(w *strings.Builder, writeLHS, writeRHS func() error) error {
	if err := writeLHS(); err != nil {
		return err
	}
	w.WriteString(" || ")
	return writeRHS()
}

// WriteRegexMatch returns an error as SQLite does not natively support regex.
func (d *Dialect) WriteRegexMatch(_ *strings.Builder, _ func() error, _ string, _ bool) error {
	return fmt.Errorf("%w: regex matching", dialect.ErrUnsupportedFeature)
}

// WriteLikeEscape writes the SQLite LIKE escape clause.
// SQLite does not use backslash escaping in string literals, so '\' is a single character.
func (d *Dialect) WriteLikeEscape(w *strings.Builder) {
	w.WriteString(" ESCAPE '\\'")
}

// WriteArrayMembership writes a SQLite array membership test using json_each.
func (d *Dialect) WriteArrayMembership(w *strings.Builder, writeElem func() error, writeArray func() error) error {
	if err := writeElem(); err != nil {
		return err
	}
	w.WriteString(" IN (SELECT value FROM json_each(")
	if err := writeArray(); err != nil {
		return err
	}
	w.WriteString("))")
	return nil
}

// --- Type Casting ---

// WriteCastToNumeric writes a SQLite numeric cast (CAST(... AS REAL)).
func (d *Dialect) WriteCastToNumeric(w *strings.Builder) {
	w.WriteString(" + 0")
}

// WriteTypeName writes a SQLite type name for CAST expressions.
func (d *Dialect) WriteTypeName(w *strings.Builder, celTypeName string) {
	switch celTypeName {
	case "bool":
		w.WriteString("INTEGER")
	case "bytes":
		w.WriteString("BLOB")
	case "double":
		w.WriteString("REAL")
	case "int":
		w.WriteString("INTEGER")
	case "string":
		w.WriteString("TEXT")
	case "uint":
		w.WriteString("INTEGER")
	default:
		w.WriteString(strings.ToUpper(celTypeName))
	}
}

// WriteEpochExtract writes strftime('%s', expr) for SQLite.
func (d *Dialect) WriteEpochExtract(w *strings.Builder, writeExpr func() error) error {
	w.WriteString("CAST(strftime('%s', ")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString(") AS INTEGER)")
	return nil
}

// WriteTimestampCast writes a SQLite datetime cast.
func (d *Dialect) WriteTimestampCast(w *strings.Builder, writeExpr func() error) error {
	w.WriteString("datetime(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString(")")
	return nil
}

// --- Arrays ---

// WriteArrayLiteralOpen writes the SQLite JSON array literal opening.
func (d *Dialect) WriteArrayLiteralOpen(w *strings.Builder) {
	w.WriteString("json_array(")
}

// WriteArrayLiteralClose writes the SQLite JSON array literal closing.
func (d *Dialect) WriteArrayLiteralClose(w *strings.Builder) {
	w.WriteString(")")
}

// WriteArrayLength writes json_array_length(expr) for SQLite.
func (d *Dialect) WriteArrayLength(w *strings.Builder, _ int, writeExpr func() error) error {
	w.WriteString("COALESCE(json_array_length(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString("), 0)")
	return nil
}

// WriteListIndex writes SQLite JSON array index access.
func (d *Dialect) WriteListIndex(w *strings.Builder, writeArray func() error, writeIndex func() error) error {
	w.WriteString("json_extract(")
	if err := writeArray(); err != nil {
		return err
	}
	w.WriteString(", '$[' || ")
	if err := writeIndex(); err != nil {
		return err
	}
	w.WriteString(" || ']')")
	return nil
}

// WriteListIndexConst writes SQLite JSON constant array index access.
func (d *Dialect) WriteListIndexConst(w *strings.Builder, writeArray func() error, index int64) error {
	w.WriteString("json_extract(")
	if err := writeArray(); err != nil {
		return err
	}
	fmt.Fprintf(w, ", '$[%d]')", index)
	return nil
}

// WriteEmptyTypedArray writes an empty SQLite JSON array.
func (d *Dialect) WriteEmptyTypedArray(w *strings.Builder, _ string) {
	w.WriteString("json_array()")
}

// --- JSON ---

// WriteJSONFieldAccess writes SQLite JSON field access using json_extract.
func (d *Dialect) WriteJSONFieldAccess(w *strings.Builder, writeBase func() error, fieldName string, _ bool) error {
	w.WriteString("json_extract(")
	if err := writeBase(); err != nil {
		return err
	}
	escaped := escapeJSONFieldName(fieldName)
	w.WriteString(", '$.")
	w.WriteString(escaped)
	w.WriteString("')")
	return nil
}

// WriteJSONExistence writes a SQLite JSON key existence check.
func (d *Dialect) WriteJSONExistence(w *strings.Builder, _ bool, fieldName string, writeBase func() error) error {
	w.WriteString("json_type(")
	if err := writeBase(); err != nil {
		return err
	}
	escaped := escapeJSONFieldName(fieldName)
	w.WriteString(", '$.")
	w.WriteString(escaped)
	w.WriteString("') IS NOT NULL")
	return nil
}

// WriteJSONArrayElements writes SQLite JSON array expansion using json_each.
func (d *Dialect) WriteJSONArrayElements(w *strings.Builder, _ bool, _ bool, writeExpr func() error) error {
	w.WriteString("json_each(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString(")")
	return nil
}

// WriteJSONArrayLength writes COALESCE(json_array_length(expr), 0) for SQLite.
func (d *Dialect) WriteJSONArrayLength(w *strings.Builder, writeExpr func() error) error {
	w.WriteString("COALESCE(json_array_length(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString("), 0)")
	return nil
}

// WriteJSONExtractPath writes SQLite JSON path extraction.
func (d *Dialect) WriteJSONExtractPath(w *strings.Builder, pathSegments []string, writeRoot func() error) error {
	w.WriteString("json_type(")
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

// WriteJSONArrayMembership writes SQLite JSON array membership using json_each.
func (d *Dialect) WriteJSONArrayMembership(w *strings.Builder, _ string, writeExpr func() error) error {
	w.WriteString("(SELECT value FROM json_each(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString("))")
	return nil
}

// WriteNestedJSONArrayMembership writes SQLite nested JSON array membership.
func (d *Dialect) WriteNestedJSONArrayMembership(w *strings.Builder, writeExpr func() error) error {
	w.WriteString("(SELECT value FROM json_each(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString("))")
	return nil
}

// --- Timestamps ---

// WriteDuration writes a SQLite duration as a string modifier.
func (d *Dialect) WriteDuration(w *strings.Builder, value int64, unit string) {
	// SQLite uses datetime modifiers like '+N seconds', '+N minutes', etc.
	fmt.Fprintf(w, "'%+d %s'", value, strings.ToLower(unit)+"s")
}

// WriteInterval writes a SQLite interval expression.
func (d *Dialect) WriteInterval(w *strings.Builder, writeValue func() error, unit string) error {
	w.WriteString("'+'||")
	if err := writeValue(); err != nil {
		return err
	}
	fmt.Fprintf(w, "||' %s'", strings.ToLower(unit)+"s")
	return nil
}

// WriteExtract writes a SQLite strftime extraction expression.
func (d *Dialect) WriteExtract(w *strings.Builder, part string, writeExpr func() error, _ func() error) error {
	format := sqliteExtractFormat(part)
	w.WriteString("CAST(strftime('")
	w.WriteString(format)
	w.WriteString("', ")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString(") AS INTEGER)")
	return nil
}

// WriteTimestampArithmetic writes SQLite timestamp arithmetic using datetime().
func (d *Dialect) WriteTimestampArithmetic(w *strings.Builder, op string, writeTS, writeDur func() error) error {
	if op == "-" {
		// For subtraction, negate the duration
		w.WriteString("datetime(")
		if err := writeTS(); err != nil {
			return err
		}
		w.WriteString(", '-'||")
		if err := writeDur(); err != nil {
			return err
		}
		w.WriteString(")")
	} else {
		w.WriteString("datetime(")
		if err := writeTS(); err != nil {
			return err
		}
		w.WriteString(", ")
		if err := writeDur(); err != nil {
			return err
		}
		w.WriteString(")")
	}
	return nil
}

// --- String Functions ---

// WriteContains writes INSTR(haystack, needle) > 0 for SQLite.
func (d *Dialect) WriteContains(w *strings.Builder, writeHaystack, writeNeedle func() error) error {
	w.WriteString("INSTR(")
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

// WriteSplit returns an error as SQLite does not have a native string split.
func (d *Dialect) WriteSplit(_ *strings.Builder, _, _ func() error) error {
	return fmt.Errorf("%w: string split", dialect.ErrUnsupportedFeature)
}

// WriteSplitWithLimit returns an error as SQLite does not have a native string split.
func (d *Dialect) WriteSplitWithLimit(_ *strings.Builder, _, _ func() error, _ int64) error {
	return fmt.Errorf("%w: string split with limit", dialect.ErrUnsupportedFeature)
}

// WriteJoin returns an error as SQLite does not have a native array join.
func (d *Dialect) WriteJoin(_ *strings.Builder, _, _ func() error) error {
	return fmt.Errorf("%w: array join", dialect.ErrUnsupportedFeature)
}

// --- Comprehensions ---

// WriteUnnest writes SQLite json_each for array unnesting.
func (d *Dialect) WriteUnnest(w *strings.Builder, writeSource func() error) error {
	w.WriteString("json_each(")
	if err := writeSource(); err != nil {
		return err
	}
	w.WriteString(")")
	return nil
}

// WriteArraySubqueryOpen writes (SELECT json_group_array( for SQLite array subqueries.
func (d *Dialect) WriteArraySubqueryOpen(w *strings.Builder) {
	w.WriteString("(SELECT json_group_array(")
}

// WriteArraySubqueryExprClose closes the json_group_array aggregate function for SQLite.
func (d *Dialect) WriteArraySubqueryExprClose(w *strings.Builder) {
	w.WriteString(")")
}

// --- Struct ---

// WriteStructOpen writes the SQLite struct literal opening.
func (d *Dialect) WriteStructOpen(w *strings.Builder) {
	w.WriteString("json_object(")
}

// WriteStructClose writes the SQLite struct literal closing.
func (d *Dialect) WriteStructClose(w *strings.Builder) {
	w.WriteString(")")
}

// --- Validation ---

// MaxIdentifierLength returns 0 as SQLite has no hard identifier length limit.
func (d *Dialect) MaxIdentifierLength() int {
	return 0
}

// ValidateFieldName validates a field name against SQLite naming rules.
func (d *Dialect) ValidateFieldName(name string) error {
	return validateFieldName(name)
}

// ReservedKeywords returns the set of reserved SQL keywords for SQLite.
func (d *Dialect) ReservedKeywords() map[string]bool {
	return reservedSQLKeywords
}

// --- Regex ---

// ConvertRegex returns an error as SQLite does not natively support regex.
func (d *Dialect) ConvertRegex(_ string) (string, bool, error) {
	return "", false, fmt.Errorf("%w: regex matching", dialect.ErrUnsupportedFeature)
}

// SupportsRegex returns false as SQLite does not natively support regex.
func (d *Dialect) SupportsRegex() bool { return false }

// --- Capabilities ---

// SupportsNativeArrays returns false as SQLite uses JSON arrays.
func (d *Dialect) SupportsNativeArrays() bool { return false }

// SupportsJSONB returns false as SQLite has a single JSON type.
func (d *Dialect) SupportsJSONB() bool { return false }

// SupportsIndexAnalysis returns true as SQLite index analysis is supported.
func (d *Dialect) SupportsIndexAnalysis() bool { return true }

// --- Internal helpers ---

// escapeJSONFieldName escapes special characters in JSON field names for SQLite.
func escapeJSONFieldName(fieldName string) string {
	return strings.ReplaceAll(fieldName, "'", "''")
}

// sqliteExtractFormat maps SQL EXTRACT parts to SQLite strftime format strings.
func sqliteExtractFormat(part string) string {
	switch part {
	case "YEAR":
		return "%Y"
	case "MONTH":
		return "%m"
	case "DAY":
		return "%d"
	case "HOUR":
		return "%H"
	case "MINUTE":
		return "%M"
	case "SECOND":
		return "%S"
	case "DOY":
		return "%j"
	case "DOW":
		return "%w"
	case "MILLISECONDS":
		return "%f"
	default:
		return "%Y"
	}
}
