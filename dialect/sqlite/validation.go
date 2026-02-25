package sqlite

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	// fieldNameRegexp validates SQLite identifier format.
	fieldNameRegexp = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

	// reservedSQLKeywords contains SQLite reserved keywords.
	reservedSQLKeywords = map[string]bool{
		"abort": true, "action": true, "add": true, "after": true, "all": true,
		"alter": true, "always": true, "analyze": true, "and": true, "as": true,
		"asc": true, "attach": true, "autoincrement": true, "before": true,
		"begin": true, "between": true, "by": true, "cascade": true, "case": true,
		"cast": true, "check": true, "collate": true, "column": true, "commit": true,
		"conflict": true, "constraint": true, "create": true, "cross": true,
		"current": true, "current_date": true, "current_time": true,
		"current_timestamp": true, "database": true, "default": true,
		"deferrable": true, "deferred": true, "delete": true, "desc": true,
		"detach": true, "distinct": true, "do": true, "drop": true, "each": true,
		"else": true, "end": true, "escape": true, "except": true, "exclude": true,
		"exclusive": true, "exists": true, "explain": true, "fail": true,
		"filter": true, "first": true, "following": true, "for": true,
		"foreign": true, "from": true, "full": true, "glob": true, "group": true,
		"groups": true, "having": true, "if": true, "ignore": true, "immediate": true,
		"in": true, "index": true, "indexed": true, "initially": true, "inner": true,
		"insert": true, "instead": true, "intersect": true, "into": true, "is": true,
		"isnull": true, "join": true, "key": true, "last": true, "left": true,
		"like": true, "limit": true, "match": true, "materialized": true,
		"natural": true, "no": true, "not": true, "nothing": true, "notnull": true,
		"null": true, "nulls": true, "of": true, "offset": true, "on": true,
		"or": true, "order": true, "others": true, "outer": true, "over": true,
		"partition": true, "plan": true, "pragma": true, "preceding": true,
		"primary": true, "query": true, "raise": true, "range": true,
		"recursive": true, "references": true, "regexp": true, "reindex": true,
		"release": true, "rename": true, "replace": true, "restrict": true,
		"returning": true, "right": true, "rollback": true, "row": true,
		"rows": true, "savepoint": true, "select": true, "set": true, "table": true,
		"temp": true, "temporary": true, "then": true, "ties": true, "to": true,
		"transaction": true, "trigger": true, "unbounded": true, "union": true,
		"unique": true, "update": true, "using": true, "vacuum": true, "values": true,
		"view": true, "virtual": true, "when": true, "where": true, "window": true,
		"with": true, "without": true,
	}
)

// validateFieldName validates that a field name follows SQLite naming conventions.
func validateFieldName(name string) error {
	if len(name) == 0 {
		return errors.New("field name cannot be empty")
	}

	// SQLite has no hard limit on identifier length but we use a reasonable limit
	if !fieldNameRegexp.MatchString(name) {
		return fmt.Errorf("field name %q must start with a letter or underscore and contain only alphanumeric characters and underscores", name)
	}

	if reservedSQLKeywords[strings.ToLower(name)] {
		return fmt.Errorf("field name %q is a reserved SQL keyword and cannot be used without quoting", name)
	}

	return nil
}
