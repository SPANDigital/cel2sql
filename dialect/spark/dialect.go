// Package spark implements the Apache Spark SQL dialect for cel2sql.
package spark

import (
	"fmt"
	"strings"

	"github.com/spandigital/cel2sql/v3/dialect"
)

// Dialect implements dialect.Dialect for Apache Spark SQL.
type Dialect struct{}

// New creates a new Spark SQL dialect.
func New() *Dialect {
	return &Dialect{}
}

func init() {
	dialect.Register(dialect.Spark, func() dialect.Dialect { return New() })
}

// Ensure Dialect implements dialect.Dialect at compile time.
var _ dialect.Dialect = (*Dialect)(nil)

// Name returns the dialect name.
func (d *Dialect) Name() dialect.Name { return dialect.Spark }

// --- Literals ---

// WriteStringLiteral writes a Spark SQL string literal with ” escaping.
func (d *Dialect) WriteStringLiteral(w *strings.Builder, value string) {
	escaped := strings.ReplaceAll(value, "'", "''")
	w.WriteString("'")
	w.WriteString(escaped)
	w.WriteString("'")
}

// WriteBytesLiteral writes a Spark SQL byte literal as X'HEX'.
func (d *Dialect) WriteBytesLiteral(w *strings.Builder, value []byte) error {
	w.WriteString("X'")
	for _, b := range value {
		fmt.Fprintf(w, "%02X", b)
	}
	w.WriteString("'")
	return nil
}

// WriteParamPlaceholder writes a Spark SQL positional parameter (?).
// The paramIndex argument is intentionally unused: Spark JDBC uses
// positional ? placeholders, so the converter relies on parameter order.
func (d *Dialect) WriteParamPlaceholder(w *strings.Builder, _ int) {
	w.WriteString("?")
}

// --- Operators ---

// WriteStringConcat writes Spark string concatenation using the concat() function.
// concat() works in all Spark versions; the || operator was added in 3.0+.
func (d *Dialect) WriteStringConcat(w *strings.Builder, writeLHS, writeRHS func() error) error {
	w.WriteString("concat(")
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

// WriteRegexMatch writes a Spark SQL regex match using RLIKE.
// Spark regex uses Java pattern syntax; (?i) inline flag is supported,
// so caseInsensitive is folded into the pattern by ConvertRegex.
func (d *Dialect) WriteRegexMatch(w *strings.Builder, writeTarget func() error, pattern string, _ bool) error {
	if err := writeTarget(); err != nil {
		return err
	}
	w.WriteString(" RLIKE '")
	escaped := strings.ReplaceAll(pattern, "'", "''")
	w.WriteString(escaped)
	w.WriteString("'")
	return nil
}

// WriteLikeEscape writes the Spark SQL LIKE escape clause.
func (d *Dialect) WriteLikeEscape(w *strings.Builder) {
	w.WriteString(" ESCAPE '\\\\'")
}

// WriteArrayMembership writes a Spark array membership test using array_contains().
func (d *Dialect) WriteArrayMembership(w *strings.Builder, writeElem, writeArray func() error) error {
	w.WriteString("array_contains(")
	if err := writeArray(); err != nil {
		return err
	}
	w.WriteString(", ")
	if err := writeElem(); err != nil {
		return err
	}
	w.WriteString(")")
	return nil
}

// --- Type Casting ---

// WriteCastToNumeric writes a Spark numeric coercion suffix.
// Spark does not support PostgreSQL-style `::TYPE` postfix casts, so we use the
// arithmetic coercion `+ 0` (same convention MySQL and SQLite use): `'5' + 0`
// evaluates as a number in Spark, ensuring JSON text extractions are compared
// numerically rather than lexicographically.
func (d *Dialect) WriteCastToNumeric(w *strings.Builder) {
	w.WriteString(" + 0")
}

// WriteTypeName writes a Spark SQL type name for CAST expressions.
func (d *Dialect) WriteTypeName(w *strings.Builder, celTypeName string) {
	switch celTypeName {
	case "bool":
		w.WriteString("BOOLEAN")
	case "bytes":
		w.WriteString("BINARY")
	case "double":
		w.WriteString("DOUBLE")
	case "int":
		w.WriteString("BIGINT")
	case "string":
		w.WriteString("STRING")
	case "uint":
		w.WriteString("BIGINT")
	default:
		w.WriteString(strings.ToUpper(celTypeName))
	}
}

// WriteEpochExtract writes UNIX_TIMESTAMP(expr) for Spark.
func (d *Dialect) WriteEpochExtract(w *strings.Builder, writeExpr func() error) error {
	w.WriteString("UNIX_TIMESTAMP(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString(")")
	return nil
}

// WriteTimestampCast writes a Spark CAST to TIMESTAMP.
func (d *Dialect) WriteTimestampCast(w *strings.Builder, writeExpr func() error) error {
	w.WriteString("CAST(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString(" AS TIMESTAMP)")
	return nil
}

// --- Arrays ---

// WriteArrayLiteralOpen writes the Spark array literal opening (array().
func (d *Dialect) WriteArrayLiteralOpen(w *strings.Builder) {
	w.WriteString("array(")
}

// WriteArrayLiteralClose writes the Spark array literal closing.
func (d *Dialect) WriteArrayLiteralClose(w *strings.Builder) {
	w.WriteString(")")
}

// WriteArrayLength writes COALESCE(size(expr), 0) for Spark.
// Spark's size() returns -1 for null; COALESCE+ifnull-style handling matches cel2sql semantics
// where size(null) should be 0.
func (d *Dialect) WriteArrayLength(w *strings.Builder, dimension int, writeExpr func() error) error {
	if dimension > 1 {
		return fmt.Errorf("%w: spark dialect does not support multi-dimensional arrays (dimension=%d)", dialect.ErrUnsupportedFeature, dimension)
	}
	w.WriteString("COALESCE(size(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString("), 0)")
	return nil
}

// WriteListIndex writes Spark 0-indexed array access.
func (d *Dialect) WriteListIndex(w *strings.Builder, writeArray, writeIndex func() error) error {
	if err := writeArray(); err != nil {
		return err
	}
	w.WriteString("[")
	if err := writeIndex(); err != nil {
		return err
	}
	w.WriteString("]")
	return nil
}

// WriteListIndexConst writes a Spark constant array index (0-indexed).
func (d *Dialect) WriteListIndexConst(w *strings.Builder, writeArray func() error, index int64) error {
	if err := writeArray(); err != nil {
		return err
	}
	fmt.Fprintf(w, "[%d]", index)
	return nil
}

// WriteEmptyTypedArray writes an empty Spark typed array.
func (d *Dialect) WriteEmptyTypedArray(w *strings.Builder, typeName string) {
	w.WriteString("CAST(array() AS ARRAY<")
	w.WriteString(sparkTypeName(typeName))
	w.WriteString(">)")
}

// --- JSON ---

// WriteJSONFieldAccess writes Spark JSON field access using get_json_object.
// Spark's get_json_object always returns a string; the same function is used
// for both intermediate and final access (Spark has no JSON_QUERY equivalent).
func (d *Dialect) WriteJSONFieldAccess(w *strings.Builder, writeBase func() error, fieldName string, _ bool) error {
	w.WriteString("get_json_object(")
	if err := writeBase(); err != nil {
		return err
	}
	w.WriteString(", '$.")
	w.WriteString(escapeJSONFieldName(fieldName))
	w.WriteString("')")
	return nil
}

// WriteJSONExistence writes a Spark JSON key existence check.
func (d *Dialect) WriteJSONExistence(w *strings.Builder, _ bool, fieldName string, writeBase func() error) error {
	w.WriteString("get_json_object(")
	if err := writeBase(); err != nil {
		return err
	}
	w.WriteString(", '$.")
	w.WriteString(escapeJSONFieldName(fieldName))
	w.WriteString("') IS NOT NULL")
	return nil
}

// WriteJSONArrayElements writes Spark JSON array expansion as EXPLODE(from_json(...)).
// The converter uses this in `FROM <here> AS iter`, so the result must be a
// set-returning expression. EXPLODE turns the parsed array into a relation of
// element rows. Element type is fixed to STRING in v1; comparisons coerce via
// arithmetic context (see WriteCastToNumeric).
func (d *Dialect) WriteJSONArrayElements(w *strings.Builder, _, _ bool, writeExpr func() error) error {
	w.WriteString("EXPLODE(from_json(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString(", 'ARRAY<STRING>'))")
	return nil
}

// WriteJSONArrayLength writes COALESCE(size(from_json(expr, 'ARRAY<STRING>')), 0) for Spark.
func (d *Dialect) WriteJSONArrayLength(w *strings.Builder, writeExpr func() error) error {
	w.WriteString("COALESCE(size(from_json(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString(", 'ARRAY<STRING>')), 0)")
	return nil
}

// WriteJSONExtractPath writes a Spark JSON path existence check using get_json_object.
func (d *Dialect) WriteJSONExtractPath(w *strings.Builder, pathSegments []string, writeRoot func() error) error {
	w.WriteString("get_json_object(")
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

// WriteJSONArrayMembership writes Spark JSON array membership as a scalar
// subquery that scans elements. The converter writes `lhs = ` before this,
// so the result is `lhs = (SELECT col FROM (SELECT EXPLODE(from_json(rhs,
// 'ARRAY<STRING>')) AS col) t)`. This mirrors SQLite's `lhs = (SELECT value
// FROM json_each(...))` pattern; both dialects rely on the subquery
// returning at most one match for the comparison to succeed.
func (d *Dialect) WriteJSONArrayMembership(w *strings.Builder, _ string, writeExpr func() error) error {
	w.WriteString("(SELECT col FROM (SELECT EXPLODE(from_json(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString(", 'ARRAY<STRING>')) AS col) t)")
	return nil
}

// WriteNestedJSONArrayMembership writes Spark nested JSON array membership.
func (d *Dialect) WriteNestedJSONArrayMembership(w *strings.Builder, writeExpr func() error) error {
	w.WriteString("(SELECT col FROM (SELECT EXPLODE(from_json(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString(", 'ARRAY<STRING>')) AS col) t)")
	return nil
}

// --- Timestamps ---

// WriteDuration writes a Spark INTERVAL literal.
func (d *Dialect) WriteDuration(w *strings.Builder, value int64, unit string) {
	fmt.Fprintf(w, "INTERVAL %d %s", value, unit)
}

// WriteInterval writes a Spark INTERVAL expression.
func (d *Dialect) WriteInterval(w *strings.Builder, writeValue func() error, unit string) error {
	w.WriteString("INTERVAL ")
	if err := writeValue(); err != nil {
		return err
	}
	w.WriteString(" ")
	w.WriteString(unit)
	return nil
}

// WriteExtract writes a Spark EXTRACT expression.
// Spark dayofweek() returns 1=Sunday..7=Saturday; CEL convention is 0=Sunday..6=Saturday.
func (d *Dialect) WriteExtract(w *strings.Builder, part string, writeExpr func() error, writeTZ func() error) error {
	isDOW := part == "DOW"
	if isDOW {
		w.WriteString("(dayofweek(")
		if err := writeExpr(); err != nil {
			return err
		}
		if writeTZ != nil {
			w.WriteString(" AT TIME ZONE ")
			if err := writeTZ(); err != nil {
				return err
			}
		}
		w.WriteString(") - 1)")
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

// WriteTimestampArithmetic writes Spark timestamp arithmetic.
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

// WriteContains writes Spark string contains using LOCATE() > 0.
// LOCATE(substr, str) returns 1-based position or 0 when not found.
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

// WriteSplit writes Spark string split using split().
func (d *Dialect) WriteSplit(w *strings.Builder, writeStr, writeDelim func() error) error {
	w.WriteString("split(")
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

// WriteSplitWithLimit writes Spark string split with limit (3-arg split, Spark 3.x+).
func (d *Dialect) WriteSplitWithLimit(w *strings.Builder, writeStr, writeDelim func() error, limit int64) error {
	w.WriteString("split(")
	if err := writeStr(); err != nil {
		return err
	}
	w.WriteString(", ")
	if err := writeDelim(); err != nil {
		return err
	}
	fmt.Fprintf(w, ", %d)", limit)
	return nil
}

// WriteJoin writes Spark array join using array_join().
func (d *Dialect) WriteJoin(w *strings.Builder, writeArray, writeDelim func() error) error {
	w.WriteString("array_join(")
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

// WriteUnnest writes Spark explode-style unnesting via lateral view replacement.
// Note: cel2sql wraps this in an ARRAY-building subquery; Spark uses array
// higher-order functions (transform/filter/exists/forall) which don't need UNNEST.
// For the SELECT FROM UNNEST() pattern Spark requires a lateral view. We emit
// EXPLODE() and rely on the converter's subquery scaffolding.
func (d *Dialect) WriteUnnest(w *strings.Builder, writeSource func() error) error {
	w.WriteString("EXPLODE(")
	if err := writeSource(); err != nil {
		return err
	}
	w.WriteString(")")
	return nil
}

// WriteArraySubqueryOpen writes the Spark array-building subquery prefix.
// Spark has no ARRAY(SELECT ...) constructor; collect_list() inside a subquery
// is the closest equivalent.
func (d *Dialect) WriteArraySubqueryOpen(w *strings.Builder) {
	w.WriteString("(SELECT collect_list(")
}

// WriteArraySubqueryExprClose closes the collect_list() argument list.
func (d *Dialect) WriteArraySubqueryExprClose(w *strings.Builder) {
	w.WriteString(")")
}

// --- Struct ---

// WriteStructOpen writes the Spark struct literal opening using struct().
func (d *Dialect) WriteStructOpen(w *strings.Builder) {
	w.WriteString("struct(")
}

// WriteStructClose writes the Spark struct literal closing.
func (d *Dialect) WriteStructClose(w *strings.Builder) {
	w.WriteString(")")
}

// --- Validation ---

// MaxIdentifierLength returns 128 for Spark (Hive-derived limit).
func (d *Dialect) MaxIdentifierLength() int {
	return 128
}

// ValidateFieldName validates a field name against Spark naming rules.
func (d *Dialect) ValidateFieldName(name string) error {
	return validateFieldName(name)
}

// ReservedKeywords returns the set of reserved SQL keywords for Spark.
func (d *Dialect) ReservedKeywords() map[string]bool {
	return reservedSQLKeywords
}

// --- Regex ---

// ConvertRegex converts an RE2 regex pattern to Spark/Java regex format.
// Spark uses java.util.regex.Pattern, which is largely a superset of RE2 for
// the safe patterns cel2sql accepts.
func (d *Dialect) ConvertRegex(re2Pattern string) (string, bool, error) {
	return convertRE2ToSpark(re2Pattern)
}

// SupportsRegex returns true as Spark supports regex via RLIKE.
func (d *Dialect) SupportsRegex() bool { return true }

// --- Capabilities ---

// SupportsNativeArrays returns true as Spark has native ARRAY<T> types.
func (d *Dialect) SupportsNativeArrays() bool { return true }

// SupportsJSONB returns false as Spark has no separate JSONB type.
func (d *Dialect) SupportsJSONB() bool { return false }

// SupportsIndexAnalysis returns false. Spark indexing is highly storage-layer-specific
// (Delta Z-order vs Iceberg sort vs plain Parquet) and out of scope for v1.
func (d *Dialect) SupportsIndexAnalysis() bool { return false }

// --- Internal helpers ---

// escapeJSONFieldName escapes single quotes in JSON field names for Spark.
func escapeJSONFieldName(fieldName string) string {
	return strings.ReplaceAll(fieldName, "'", "''")
}

// sparkTypeName converts a CEL/common type name to a Spark type name.
func sparkTypeName(typeName string) string {
	switch strings.ToLower(typeName) {
	case "text", "string", "varchar", "char":
		return "STRING"
	case "int", "integer", "bigint", "int64", "long":
		return "BIGINT"
	case "double", "float", "real", "float64":
		return "DOUBLE"
	case "boolean", "bool":
		return "BOOLEAN"
	case "bytes", "bytea", "blob", "binary":
		return "BINARY"
	default:
		return strings.ToUpper(typeName)
	}
}
