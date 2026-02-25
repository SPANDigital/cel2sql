package bigquery

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	// fieldNameRegexp validates BigQuery identifier format.
	fieldNameRegexp = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

	// reservedSQLKeywords contains BigQuery reserved keywords.
	reservedSQLKeywords = map[string]bool{
		"all": true, "and": true, "any": true, "array": true, "as": true,
		"asc": true, "assert_rows_modified": true, "at": true, "between": true,
		"by": true, "case": true, "cast": true, "collate": true, "contains": true,
		"create": true, "cross": true, "cube": true, "current": true,
		"default": true, "define": true, "desc": true, "distinct": true,
		"else": true, "end": true, "enum": true, "escape": true, "except": true,
		"exclude": true, "exists": true, "extract": true, "false": true,
		"fetch": true, "following": true, "for": true, "from": true, "full": true,
		"group": true, "grouping": true, "groups": true, "hash": true,
		"having": true, "if": true, "ignore": true, "in": true, "inner": true,
		"intersect": true, "interval": true, "into": true, "is": true,
		"join": true, "lateral": true, "left": true, "like": true, "limit": true,
		"lookup": true, "merge": true, "natural": true, "new": true, "no": true,
		"not": true, "null": true, "nulls": true, "of": true, "on": true,
		"or": true, "order": true, "outer": true, "over": true,
		"partition": true, "preceding": true, "proto": true, "range": true,
		"recursive": true, "respect": true, "right": true, "rollup": true,
		"rows": true, "select": true, "set": true, "some": true, "struct": true,
		"tablesample": true, "then": true, "to": true, "treat": true,
		"true": true, "unbounded": true, "union": true, "unnest": true,
		"using": true, "when": true, "where": true, "window": true,
		"with": true, "within": true,
	}
)

// validateFieldName validates that a field name follows BigQuery naming conventions.
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
