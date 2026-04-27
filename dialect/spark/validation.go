package spark

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	// fieldNameRegexp validates Spark identifier format (letter/underscore start,
	// alphanumerics + underscore body). Spark also accepts backtick-quoted
	// identifiers but we restrict to the unquoted form for safety.
	fieldNameRegexp = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

	// reservedSQLKeywords contains Spark SQL reserved keywords.
	// Source: Apache Spark docs (sql-ref-ansi-compliance.html#sql-keywords)
	// plus the standard SQL set already covered by other dialects. Lowercased.
	reservedSQLKeywords = map[string]bool{
		"all": true, "alter": true, "and": true, "anti": true, "any": true,
		"array": true, "as": true, "asc": true, "between": true, "both": true,
		"by": true, "case": true, "cast": true, "check": true, "cluster": true,
		"collate": true, "column": true, "create": true, "cross": true,
		"cube": true, "current": true, "current_date": true, "current_time": true,
		"current_timestamp": true, "current_user": true, "default": true,
		"delete": true, "desc": true, "describe": true, "distinct": true,
		"drop": true, "else": true, "end": true, "escape": true, "except": true,
		"exists": true, "false": true, "fetch": true, "filter": true,
		"for": true, "foreign": true, "from": true, "full": true, "function": true,
		"grant": true, "group": true, "grouping": true, "having": true,
		"hour": true, "in": true, "inner": true, "insert": true, "intersect": true,
		"interval": true, "into": true, "is": true, "join": true, "lateral": true,
		"leading": true, "left": true, "like": true, "limit": true, "local": true,
		"map": true, "minute": true, "month": true, "natural": true, "no": true,
		"not": true, "null": true, "of": true, "on": true, "only": true,
		"or": true, "order": true, "outer": true, "overlaps": true, "primary": true,
		"references": true, "right": true, "rollup": true, "row": true,
		"rows": true, "second": true, "select": true, "semi": true,
		"session_user": true, "set": true, "some": true, "struct": true,
		"table": true, "tablesample": true, "then": true, "time": true,
		"to": true, "trailing": true, "true": true, "union": true, "unique": true,
		"unknown": true, "update": true, "user": true, "using": true,
		"values": true, "when": true, "where": true, "window": true, "with": true,
		"year": true,
	}
)

// validateFieldName validates that a field name follows Spark naming conventions.
func validateFieldName(name string) error {
	if len(name) == 0 {
		return errors.New("field name cannot be empty")
	}

	if !fieldNameRegexp.MatchString(name) {
		return fmt.Errorf("field name %q must start with a letter or underscore and contain only alphanumeric characters and underscores", name)
	}

	if reservedSQLKeywords[strings.ToLower(name)] {
		return fmt.Errorf("field name %q is a reserved SQL keyword and cannot be used without quoting", name)
	}

	return nil
}
