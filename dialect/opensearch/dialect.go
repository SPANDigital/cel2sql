// Package opensearch implements the OpenSearch SQL dialect for cel2sql.
//
// OpenSearch SQL is the query language of the OpenSearch SQL plugin
// (POST _plugins/_sql). Its grammar is MySQL-flavoured, but it is a search
// engine, not a relational database: there are no array values, no JSON
// functions, no array-building subqueries and no user-managed indexes.
// Object fields are addressed by dot path (a.b.c) and nested fields are
// queried PartiQL-style, with EXISTS (SELECT * FROM parent.nested AS n WHERE …).
package opensearch

import (
	"fmt"
	"strings"

	"github.com/spandigital/cel2sql/v3/dialect"
)

// Dialect implements dialect.Dialect for OpenSearch SQL.
type Dialect struct{}

// New creates a new OpenSearch SQL dialect.
func New() *Dialect { return &Dialect{} }

func init() {
	dialect.Register(dialect.OpenSearch, func() dialect.Dialect { return New() })
}

// Compile-time interface check.
var _ dialect.Dialect = (*Dialect)(nil)

// Name returns the dialect name.
func (d *Dialect) Name() dialect.Name { return dialect.OpenSearch }

// --- Literals ---

// WriteStringLiteral writes an OpenSearch SQL string literal with ” escaping.
func (d *Dialect) WriteStringLiteral(w *strings.Builder, value string) {
	escaped := strings.ReplaceAll(value, "'", "''")
	w.WriteString("'")
	w.WriteString(escaped)
	w.WriteString("'")
}

// WriteBytesLiteral returns an error: OpenSearch SQL has no binary literal syntax.
func (d *Dialect) WriteBytesLiteral(_ *strings.Builder, _ []byte) error {
	return fmt.Errorf("%w: bytes literal", dialect.ErrUnsupportedFeature)
}

// WriteParamPlaceholder writes an OpenSearch SQL positional parameter (?),
// bound through the request's "parameters" array or a JDBC PreparedStatement.
func (d *Dialect) WriteParamPlaceholder(w *strings.Builder, _ int) {
	w.WriteString("?")
}

// --- Operators ---

// WriteStringConcat writes string concatenation using the two-argument CONCAT().
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

// WriteRegexMatch writes (target REGEXP 'pattern') = 1. The plugin's REGEXP
// returns an integer, not a boolean, so a bare REGEXP in WHERE fails with
// "wrong result type instead of boolean"; the comparison makes it a predicate.
// The pattern keeps any inline (?i) flag, which the Java regex engine honours,
// so caseInsensitive is ignored.
func (d *Dialect) WriteRegexMatch(w *strings.Builder, writeTarget func() error, pattern string, _ bool) error {
	w.WriteString("(")
	if err := writeTarget(); err != nil {
		return err
	}
	w.WriteString(" REGEXP ")
	d.WriteStringLiteral(w, pattern)
	w.WriteString(") = 1")
	return nil
}

// WriteBooleanFieldPredicate writes field = true. The SQL plugin silently
// ignores a bare boolean column used as a predicate (every document matches)
// and rejects NOT column, so the comparison is always spelled out.
func (d *Dialect) WriteBooleanFieldPredicate(w *strings.Builder, writeField func() error) error {
	if err := writeField(); err != nil {
		return err
	}
	w.WriteString(" = true")
	return nil
}

// WriteLikeEscape writes nothing: the OpenSearch SQL grammar has no ESCAPE
// clause. LIKE patterns are translated to wildcard queries, which treat a
// backslash-prefixed % or _ literally.
func (d *Dialect) WriteLikeEscape(_ *strings.Builder) {}

// WriteArrayMembership writes elem IN (v1, v2, …).
//
// OpenSearch has no array values, so the only list a membership test can be
// written against is a literal one, which WriteArrayLiteralOpen/Close render
// as a parenthesised value list. Anything else that reaches here (a list
// variable, a nested field, a concatenation of two literals) has no OpenSearch
// SQL equivalent, so the array operand is checked after it is written and the
// expression is rejected when it is not a single parenthesised group.
func (d *Dialect) WriteArrayMembership(w *strings.Builder, writeElem func() error, writeArray func() error) error {
	if err := writeElem(); err != nil {
		return err
	}
	w.WriteString(" IN ")
	start := w.Len()
	if err := writeArray(); err != nil {
		return err
	}
	if !isSingleParenGroup(w.String()[start:]) {
		return fmt.Errorf("%w: membership test against a non-literal list", dialect.ErrUnsupportedFeature)
	}
	return nil
}

// isSingleParenGroup reports whether s is exactly one parenthesised group,
// ignoring parentheses inside single-quoted string literals.
func isSingleParenGroup(s string) bool {
	if len(s) < 2 || s[0] != '(' || s[len(s)-1] != ')' {
		return false
	}
	depth := 0
	inString := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if inString {
			if c == '\'' {
				// A doubled quote is an escaped quote inside the literal.
				if i+1 < len(s) && s[i+1] == '\'' {
					i++
					continue
				}
				inString = false
			}
			continue
		}
		switch c {
		case '\'':
			inString = true
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 && i != len(s)-1 {
				// The opening group closed before the end: "(1) || (2, 3)".
				return false
			}
		}
	}
	return depth == 0 && !inString
}

// --- Type Casting ---

// WriteCastToNumeric writes nothing. The converter only requests this cast when
// comparing a value extracted from a JSON document; OpenSearch object
// sub-fields are typed by the index mapping, so no coercion is needed.
func (d *Dialect) WriteCastToNumeric(_ *strings.Builder) {}

// WriteTypeName writes an OpenSearch SQL type name for CAST expressions.
func (d *Dialect) WriteTypeName(w *strings.Builder, celTypeName string) {
	switch celTypeName {
	case "bool":
		w.WriteString("BOOLEAN")
	case "double":
		w.WriteString("DOUBLE")
	case "int", "uint":
		w.WriteString("LONG")
	case "string":
		w.WriteString("STRING")
	default:
		w.WriteString(strings.ToUpper(celTypeName))
	}
}

// WriteEpochExtract writes CAST(UNIX_TIMESTAMP(expr) AS LONG); the plugin's
// unix_timestamp returns a double.
func (d *Dialect) WriteEpochExtract(w *strings.Builder, writeExpr func() error) error {
	w.WriteString("CAST(UNIX_TIMESTAMP(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString(") AS LONG)")
	return nil
}

// WriteTimestampCast writes CAST(expr AS TIMESTAMP).
func (d *Dialect) WriteTimestampCast(w *strings.Builder, writeExpr func() error) error {
	w.WriteString("CAST(")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString(" AS TIMESTAMP)")
	return nil
}

// --- Arrays ---

// WriteArrayLiteralOpen opens a parenthesised value list, the only list form
// OpenSearch SQL has (as the right-hand side of IN).
func (d *Dialect) WriteArrayLiteralOpen(w *strings.Builder) {
	w.WriteString("(")
}

// WriteArrayLiteralClose closes the value list.
func (d *Dialect) WriteArrayLiteralClose(w *strings.Builder) {
	w.WriteString(")")
}

// WriteArrayLength returns an error: OpenSearch SQL has no array length function.
func (d *Dialect) WriteArrayLength(_ *strings.Builder, _ int, _ func() error) error {
	return fmt.Errorf("%w: array length", dialect.ErrUnsupportedFeature)
}

// WriteListIndex returns an error: OpenSearch SQL has no array indexing.
func (d *Dialect) WriteListIndex(_ *strings.Builder, _ func() error, _ func() error) error {
	return fmt.Errorf("%w: list indexing", dialect.ErrUnsupportedFeature)
}

// WriteListIndexConst returns an error: OpenSearch SQL has no array indexing.
func (d *Dialect) WriteListIndexConst(_ *strings.Builder, _ func() error, _ int64) error {
	return fmt.Errorf("%w: list indexing", dialect.ErrUnsupportedFeature)
}

// WriteEmptyTypedArray writes (NULL): x IN (NULL) is never true, which matches
// membership in an empty list.
func (d *Dialect) WriteEmptyTypedArray(w *strings.Builder, _ string) {
	w.WriteString("(NULL)")
}

// --- JSON ---

// WriteJSONFieldAccess writes base.field: object and flat_object sub-fields
// are addressed by dot path in OpenSearch SQL, whether or not the access is
// the final step of the path.
func (d *Dialect) WriteJSONFieldAccess(w *strings.Builder, writeBase func() error, fieldName string, _ bool) error {
	if err := validateFieldName(fieldName); err != nil {
		return err
	}
	if err := writeBase(); err != nil {
		return err
	}
	w.WriteString(".")
	w.WriteString(fieldName)
	return nil
}

// WriteJSONExistence writes base.field IS NOT NULL. OpenSearch does not
// distinguish a missing field from one explicitly set to null.
func (d *Dialect) WriteJSONExistence(w *strings.Builder, _ bool, fieldName string, writeBase func() error) error {
	if err := validateFieldName(fieldName); err != nil {
		return err
	}
	if err := writeBase(); err != nil {
		return err
	}
	w.WriteString(".")
	w.WriteString(fieldName)
	w.WriteString(" IS NOT NULL")
	return nil
}

// WriteJSONArrayElements writes the nested field itself: the comprehension
// source for a nested field is `parent.field AS alias`, the PartiQL form the
// plugin rewrites to a nested query.
func (d *Dialect) WriteJSONArrayElements(_ *strings.Builder, _ bool, _ bool, writeExpr func() error) error {
	return writeExpr()
}

// WriteJSONArrayLength returns an error: OpenSearch SQL has no array length function.
func (d *Dialect) WriteJSONArrayLength(_ *strings.Builder, _ func() error) error {
	return fmt.Errorf("%w: array length", dialect.ErrUnsupportedFeature)
}

// WriteJSONExtractPath writes root.seg1.seg2 IS NOT NULL.
func (d *Dialect) WriteJSONExtractPath(w *strings.Builder, pathSegments []string, writeRoot func() error) error {
	for _, segment := range pathSegments {
		if err := validateFieldName(segment); err != nil {
			return err
		}
	}
	if err := writeRoot(); err != nil {
		return err
	}
	for _, segment := range pathSegments {
		w.WriteString(".")
		w.WriteString(segment)
	}
	w.WriteString(" IS NOT NULL")
	return nil
}

// WriteJSONArrayMembership returns an error. A multi-valued OpenSearch field
// is not an array to the SQL plugin: expressions over it only see the first
// value, so `elem IN field` cannot be expressed faithfully.
func (d *Dialect) WriteJSONArrayMembership(_ *strings.Builder, _ string, _ func() error, _ func() error) error {
	return fmt.Errorf("%w: membership test on a multi-value field", dialect.ErrUnsupportedFeature)
}

// WriteNestedJSONArrayMembership returns an error for the same reason as
// WriteJSONArrayMembership.
func (d *Dialect) WriteNestedJSONArrayMembership(_ *strings.Builder, _ func() error, _ func() error) error {
	return fmt.Errorf("%w: membership test on a multi-value field", dialect.ErrUnsupportedFeature)
}

// --- Timestamps ---

// WriteDuration writes INTERVAL n UNIT. OpenSearch has no MILLISECOND unit, so
// millisecond durations are written in microseconds.
func (d *Dialect) WriteDuration(w *strings.Builder, value int64, unit string) {
	if unit == "MILLISECOND" {
		value *= 1000
		unit = "MICROSECOND"
	}
	fmt.Fprintf(w, "INTERVAL %d %s", value, unit)
}

// WriteInterval writes INTERVAL expr UNIT.
func (d *Dialect) WriteInterval(w *strings.Builder, writeValue func() error, unit string) error {
	w.WriteString("INTERVAL ")
	if err := writeValue(); err != nil {
		return err
	}
	w.WriteString(" ")
	w.WriteString(unit)
	return nil
}

// WriteExtract writes a date-part extraction.
//
// DOW: OpenSearch's DAYOFWEEK is MySQL's (1=Sunday … 7=Saturday) and CEL's
// getDayOfWeek is ISO (0=Monday … 6=Sunday), so (DAYOFWEEK(x) + 5) % 7.
// DOY and MILLISECONDS have no EXTRACT part and use functions instead.
// Time zone conversion is rejected: the grammar has no AT TIME ZONE.
func (d *Dialect) WriteExtract(w *strings.Builder, part string, writeExpr func() error, writeTZ func() error) error {
	if writeTZ != nil {
		return fmt.Errorf("%w: time zone conversion in date extraction", dialect.ErrUnsupportedFeature)
	}

	switch part {
	case "DOW":
		w.WriteString("(DAYOFWEEK(")
		if err := writeExpr(); err != nil {
			return err
		}
		w.WriteString(") + 5) % 7")
		return nil
	case "DOY":
		w.WriteString("DAYOFYEAR(")
		if err := writeExpr(); err != nil {
			return err
		}
		w.WriteString(")")
		return nil
	case "MILLISECONDS":
		w.WriteString("FLOOR(MICROSECOND(")
		if err := writeExpr(); err != nil {
			return err
		}
		w.WriteString(") / 1000)")
		return nil
	}

	w.WriteString("EXTRACT(")
	w.WriteString(part)
	w.WriteString(" FROM ")
	if err := writeExpr(); err != nil {
		return err
	}
	w.WriteString(")")
	return nil
}

// WriteTimestampArithmetic writes DATE_ADD(ts, interval) or DATE_SUB(ts, interval).
func (d *Dialect) WriteTimestampArithmetic(w *strings.Builder, op string, writeTS, writeDur func() error) error {
	switch op {
	case "+":
		w.WriteString("DATE_ADD(")
	case "-":
		w.WriteString("DATE_SUB(")
	default:
		return fmt.Errorf("%w: timestamp arithmetic operator %q", dialect.ErrUnsupportedFeature, op)
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

// WriteContains writes LOCATE(needle, haystack) > 0.
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

// WriteSplit returns an error: OpenSearch SQL has no array values to split into.
func (d *Dialect) WriteSplit(_ *strings.Builder, _, _ func() error) error {
	return fmt.Errorf("%w: string split", dialect.ErrUnsupportedFeature)
}

// WriteSplitWithLimit returns an error: OpenSearch SQL has no array values to split into.
func (d *Dialect) WriteSplitWithLimit(_ *strings.Builder, _, _ func() error, _ int64) error {
	return fmt.Errorf("%w: string split", dialect.ErrUnsupportedFeature)
}

// WriteJoin returns an error: OpenSearch SQL has no array values to join.
func (d *Dialect) WriteJoin(_ *strings.Builder, _, _ func() error) error {
	return fmt.Errorf("%w: array join", dialect.ErrUnsupportedFeature)
}

// --- Comprehensions ---

// WriteIterVarRef writes the alias unchanged: `FROM parent.nested AS n` binds
// n to one nested document.
func (d *Dialect) WriteIterVarRef(w *strings.Builder, alias string) error {
	w.WriteString(alias)
	return nil
}

// WriteComprehensionExists writes EXISTS (SELECT * FROM <body>), the form the
// plugin documents for nested fields and rewrites to a nested query.
func (d *Dialect) WriteComprehensionExists(w *strings.Builder, writeBody func() error) error {
	w.WriteString("EXISTS (SELECT * FROM ")
	if err := writeBody(); err != nil {
		return err
	}
	w.WriteString(")")
	return nil
}

// WriteComprehensionNotExists returns an error. NOT EXISTS itself is
// accepted, but the body of an all() comprehension is WHERE NOT (<predicate>)
// and the plugin's nested-query handler rejects NOT inside the subquery, so
// the form can never run.
func (d *Dialect) WriteComprehensionNotExists(_ *strings.Builder, _ func() error) error {
	return fmt.Errorf("%w: all() comprehension (NOT inside a nested subquery)", dialect.ErrUnsupportedFeature)
}

// WriteUnnest writes the source expression bare: a nested field is iterated
// as `parent.field AS alias` with no table function.
func (d *Dialect) WriteUnnest(_ *strings.Builder, writeSource func() error) error {
	return writeSource()
}

// WriteArraySubqueryOpen is unreachable: SupportsSubqueryExpressions is false,
// so the converter rejects filter/map comprehensions before writing one.
func (d *Dialect) WriteArraySubqueryOpen(w *strings.Builder) {
	w.WriteString("(SELECT ")
}

// WriteArraySubqueryExprClose is unreachable; see WriteArraySubqueryOpen.
func (d *Dialect) WriteArraySubqueryExprClose(_ *strings.Builder) {}

// --- Struct ---

// WriteStructOpen writes the opening of a struct literal. OpenSearch SQL has no
// struct constructor; this is only reached by unreachable map/filter output.
func (d *Dialect) WriteStructOpen(w *strings.Builder) {
	w.WriteString("(")
}

// WriteStructClose writes the closing of a struct literal.
func (d *Dialect) WriteStructClose(w *strings.Builder) {
	w.WriteString(")")
}

// --- Validation ---

// MaxIdentifierLength returns 0: OpenSearch documents no identifier length limit.
func (d *Dialect) MaxIdentifierLength() int { return 0 }

// ValidateFieldName validates a field name against OpenSearch SQL identifier rules.
func (d *Dialect) ValidateFieldName(name string) error {
	return validateFieldName(name)
}

// ReservedKeywords returns the set of reserved OpenSearch SQL keywords.
func (d *Dialect) ReservedKeywords() map[string]bool {
	return reservedSQLKeywords
}

// --- Regex ---

// ConvertRegex validates an RE2 pattern for use with REGEXP.
func (d *Dialect) ConvertRegex(re2Pattern string) (string, bool, error) {
	return convertRE2ToOpenSearch(re2Pattern)
}

// SupportsRegex returns true: OpenSearch SQL has a REGEXP operator.
func (d *Dialect) SupportsRegex() bool { return true }

// --- Capabilities ---

// SupportsNativeArrays returns false: OpenSearch has multi-valued fields, not arrays.
func (d *Dialect) SupportsNativeArrays() bool { return false }

// SupportsJSONB returns false.
func (d *Dialect) SupportsJSONB() bool { return false }

// SupportsIndexAnalysis returns false: OpenSearch indexes every mapped field
// itself and has no user-managed SQL indexes.
func (d *Dialect) SupportsIndexAnalysis() bool { return false }

// SupportsSubqueryExpressions returns false: the plugin only allows a subquery
// inside an EXISTS or IN predicate, never as a value.
func (d *Dialect) SupportsSubqueryExpressions() bool { return false }
