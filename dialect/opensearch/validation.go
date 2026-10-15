package opensearch

import (
	"regexp"

	"github.com/spandigital/cel2sql/v3/dialect/internal/identsafe"
)

var (
	// fieldNameRegexp validates OpenSearch SQL identifier format. Field names in
	// an index mapping may contain other characters, but cel2sql never quotes
	// identifiers, so only names that are valid bare identifiers are accepted.
	fieldNameRegexp = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

	// reservedSQLKeywords contains the keywords of the OpenSearch SQL grammar
	// (OpenSearchSQLLexer.g4) that cannot be used as a bare identifier. The
	// grammar's keywordsCanBeId rule lets function names such as nested, score,
	// query and type-ish words such as host or message be used as field names,
	// so those are deliberately absent.
	reservedSQLKeywords = map[string]bool{
		"all": true, "and": true, "as": true, "asc": true, "between": true,
		"boolean": true, "by": true, "case": true, "cast": true, "cross": true,
		"date": true, "datetime": true, "day": true, "delete": true, "desc": true,
		"describe": true, "distinct": true, "double": true, "else": true, "end": true,
		"escape": true, "exists": true, "false": true, "first": true, "float": true,
		"from": true, "group": true, "having": true, "hour": true, "in": true,
		"inner": true, "int": true, "integer": true, "interval": true, "is": true,
		"join": true, "last": true, "left": true, "like": true, "limit": true,
		"long": true, "match": true, "microsecond": true, "minus": true, "minute": true,
		"missing": true, "month": true, "natural": true, "not": true, "null": true,
		"nulls": true, "on": true, "or": true, "order": true, "outer": true,
		"over": true, "partition": true, "quarter": true, "regexp": true, "right": true,
		"second": true, "select": true, "show": true, "string": true, "then": true,
		"time": true, "timestamp": true, "true": true, "union": true, "using": true,
		"week": true, "when": true, "where": true, "year": true,
	}
)

// validateFieldName validates that a field name is a bare OpenSearch SQL
// identifier. OpenSearch documents no identifier length limit, so none is
// enforced.
func validateFieldName(name string) error {
	return identsafe.ValidateFieldName(name, "OpenSearch", 0, fieldNameRegexp, reservedSQLKeywords)
}
