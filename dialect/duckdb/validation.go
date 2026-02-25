package duckdb

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	// fieldNameRegexp validates DuckDB identifier format.
	fieldNameRegexp = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

	// reservedSQLKeywords contains DuckDB reserved keywords.
	reservedSQLKeywords = map[string]bool{
		"all": true, "alter": true, "analyse": true, "analyze": true, "and": true,
		"any": true, "array": true, "as": true, "asc": true, "asymmetric": true,
		"between": true, "both": true, "case": true, "cast": true, "check": true,
		"collate": true, "column": true, "constraint": true, "create": true,
		"cross": true, "current_catalog": true, "current_date": true,
		"current_role": true, "current_schema": true, "current_time": true,
		"current_timestamp": true, "current_user": true, "default": true,
		"deferrable": true, "desc": true, "distinct": true, "do": true,
		"else": true, "end": true, "except": true, "exists": true, "false": true,
		"fetch": true, "for": true, "foreign": true, "from": true, "full": true,
		"grant": true, "group": true, "having": true, "in": true, "initially": true,
		"inner": true, "intersect": true, "into": true, "is": true, "isnull": true,
		"join": true, "lateral": true, "leading": true, "left": true, "like": true,
		"limit": true, "localtime": true, "localtimestamp": true, "natural": true,
		"not": true, "notnull": true, "null": true, "offset": true, "on": true,
		"only": true, "or": true, "order": true, "outer": true, "overlaps": true,
		"placing": true, "primary": true, "references": true, "returning": true,
		"right": true, "select": true, "session_user": true, "similar": true,
		"some": true, "symmetric": true, "table": true, "then": true, "to": true,
		"trailing": true, "true": true, "union": true, "unique": true, "using": true,
		"variadic": true, "when": true, "where": true, "window": true, "with": true,
	}
)

// validateFieldName validates that a field name follows DuckDB naming conventions.
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
