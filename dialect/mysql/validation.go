package mysql

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	// maxMySQLIdentifierLength is the maximum length for MySQL identifiers.
	maxMySQLIdentifierLength = 64
)

var (
	// fieldNameRegexp validates MySQL identifier format.
	fieldNameRegexp = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

	// reservedSQLKeywords contains MySQL reserved keywords.
	reservedSQLKeywords = map[string]bool{
		"accessible": true, "add": true, "all": true, "alter": true, "analyze": true,
		"and": true, "as": true, "asc": true, "asensitive": true, "before": true,
		"between": true, "bigint": true, "binary": true, "blob": true, "both": true,
		"by": true, "call": true, "cascade": true, "case": true, "change": true,
		"char": true, "character": true, "check": true, "collate": true, "column": true,
		"condition": true, "constraint": true, "continue": true, "convert": true,
		"create": true, "cross": true, "current_date": true, "current_time": true,
		"current_timestamp": true, "current_user": true, "cursor": true, "database": true,
		"databases": true, "day_hour": true, "day_microsecond": true, "day_minute": true,
		"day_second": true, "dec": true, "decimal": true, "declare": true, "default": true,
		"delayed": true, "delete": true, "desc": true, "describe": true, "deterministic": true,
		"distinct": true, "distinctrow": true, "div": true, "double": true, "drop": true,
		"dual": true, "each": true, "else": true, "elseif": true, "enclosed": true,
		"escaped": true, "exists": true, "exit": true, "explain": true, "false": true,
		"fetch": true, "float": true, "float4": true, "float8": true, "for": true,
		"force": true, "foreign": true, "from": true, "fulltext": true, "grant": true,
		"group": true, "having": true, "high_priority": true, "hour_microsecond": true,
		"hour_minute": true, "hour_second": true, "if": true, "ignore": true, "in": true,
		"index": true, "infile": true, "inner": true, "inout": true, "insensitive": true,
		"insert": true, "int": true, "int1": true, "int2": true, "int3": true,
		"int4": true, "int8": true, "integer": true, "interval": true, "into": true,
		"is": true, "iterate": true, "join": true, "key": true, "keys": true,
		"kill": true, "leading": true, "leave": true, "left": true, "like": true,
		"limit": true, "linear": true, "lines": true, "load": true, "localtime": true,
		"localtimestamp": true, "lock": true, "long": true, "longblob": true,
		"longtext": true, "loop": true, "low_priority": true, "match": true,
		"mediumblob": true, "mediumint": true, "mediumtext": true, "middleint": true,
		"minute_microsecond": true, "minute_second": true, "mod": true, "modifies": true,
		"natural": true, "not": true, "null": true, "numeric": true, "on": true,
		"optimize": true, "option": true, "optionally": true, "or": true, "order": true,
		"out": true, "outer": true, "outfile": true, "precision": true, "primary": true,
		"procedure": true, "purge": true, "range": true, "read": true, "reads": true,
		"real": true, "references": true, "regexp": true, "release": true, "rename": true,
		"repeat": true, "replace": true, "require": true, "restrict": true, "return": true,
		"revoke": true, "right": true, "rlike": true, "schema": true, "schemas": true,
		"second_microsecond": true, "select": true, "sensitive": true, "separator": true,
		"set": true, "show": true, "signal": true, "smallint": true, "spatial": true,
		"specific": true, "sql": true, "sqlexception": true, "sqlstate": true,
		"sqlwarning": true, "sql_big_result": true, "sql_calc_found_rows": true,
		"sql_small_result": true, "ssl": true, "starting": true, "straight_join": true,
		"table": true, "terminated": true, "then": true, "tinyblob": true, "tinyint": true,
		"tinytext": true, "to": true, "trailing": true, "trigger": true, "true": true,
		"undo": true, "union": true, "unique": true, "unlock": true, "unsigned": true,
		"update": true, "usage": true, "use": true, "using": true, "utc_date": true,
		"utc_time": true, "utc_timestamp": true, "values": true, "varbinary": true,
		"varchar": true, "varcharacter": true, "varying": true, "when": true,
		"where": true, "while": true, "with": true, "write": true, "xor": true,
		"year_month": true, "zerofill": true,
	}
)

// validateFieldName validates that a field name follows MySQL naming conventions.
func validateFieldName(name string) error {
	if len(name) == 0 {
		return errors.New("field name cannot be empty")
	}

	if len(name) > maxMySQLIdentifierLength {
		return fmt.Errorf("field name %q exceeds MySQL maximum identifier length of %d characters", name, maxMySQLIdentifierLength)
	}

	if !fieldNameRegexp.MatchString(name) {
		return fmt.Errorf("field name %q must start with a letter or underscore and contain only alphanumeric characters and underscores", name)
	}

	if reservedSQLKeywords[strings.ToLower(name)] {
		return fmt.Errorf("field name %q is a reserved SQL keyword and cannot be used without quoting", name)
	}

	return nil
}
