package postgres

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	// maxPostgreSQLIdentifierLength is the maximum length for PostgreSQL identifiers
	// PostgreSQL's NAMEDATALEN is 64 bytes (including null terminator), so max usable length is 63
	maxPostgreSQLIdentifierLength = 63
)

var (
	// fieldNameRegexp validates PostgreSQL identifier format
	fieldNameRegexp = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

	// reservedSQLKeywords contains SQL keywords that should not be used as unquoted identifiers
	reservedSQLKeywords = map[string]bool{
		"all": true, "analyse": true, "analyze": true, "and": true, "any": true,
		"array": true, "as": true, "asc": true, "asymmetric": true, "both": true,
		"case": true, "cast": true, "check": true, "collate": true, "column": true,
		"constraint": true, "create": true, "cross": true, "current_catalog": true,
		"current_date": true, "current_role": true, "current_time": true,
		"current_timestamp": true, "current_user": true, "default": true,
		"deferrable": true, "desc": true, "distinct": true, "do": true, "else": true,
		"end": true, "except": true, "false": true, "fetch": true, "for": true,
		"foreign": true, "from": true, "grant": true, "group": true, "having": true,
		"in": true, "initially": true, "inner": true, "intersect": true, "into": true,
		"is": true, "join": true, "leading": true, "left": true, "like": true,
		"limit": true, "localtime": true, "localtimestamp": true, "natural": true,
		"not": true, "null": true, "offset": true, "on": true, "only": true,
		"or": true, "order": true, "outer": true, "overlaps": true, "placing": true,
		"primary": true, "references": true, "returning": true, "right": true,
		"select": true, "session_user": true, "similar": true, "some": true,
		"symmetric": true, "table": true, "then": true, "to": true, "trailing": true,
		"true": true, "union": true, "unique": true, "user": true, "using": true,
		"variadic": true, "when": true, "where": true, "window": true, "with": true,
		// Additional keywords that commonly cause issues
		"alter": true, "delete": true, "drop": true, "insert": true, "update": true,
	}
)

// validateFieldName validates that a field name follows PostgreSQL naming conventions
// and is safe to use in SQL queries without quoting.
func validateFieldName(name string) error {
	if len(name) == 0 {
		return errors.New("field name cannot be empty")
	}

	if len(name) > maxPostgreSQLIdentifierLength {
		return fmt.Errorf("field name %q exceeds PostgreSQL maximum identifier length of %d characters", name, maxPostgreSQLIdentifierLength)
	}

	if !fieldNameRegexp.MatchString(name) {
		return fmt.Errorf("field name %q must start with a letter or underscore and contain only alphanumeric characters and underscores", name)
	}

	if reservedSQLKeywords[strings.ToLower(name)] {
		return fmt.Errorf("field name %q is a reserved SQL keyword and cannot be used without quoting", name)
	}

	return nil
}
