package mysql

import (
	"fmt"
	"strings"

	"github.com/spandigital/cel2sql/v3/dialect"
)

// MySQL index type constants.
const (
	IndexTypeBTree    = "BTREE"
	IndexTypeFullText = "FULLTEXT"
)

// RecommendIndex generates a MySQL-specific index recommendation for the given pattern.
// Returns nil if no applicable index exists for this pattern.
func (d *Dialect) RecommendIndex(pattern dialect.IndexPattern) *dialect.IndexRecommendation {
	table := pattern.TableHint
	if table == "" {
		table = "table_name"
	}
	col := pattern.Column
	safeName := sanitizeIndexName(col)

	switch pattern.Pattern {
	case dialect.PatternComparison:
		return &dialect.IndexRecommendation{
			Column:    col,
			IndexType: IndexTypeBTree,
			Expression: fmt.Sprintf("CREATE INDEX idx_%s_btree ON %s (%s);",
				safeName, table, col),
			Reason: fmt.Sprintf("Comparison operations on '%s' benefit from B-tree index for efficient range queries and equality checks", col),
		}

	case dialect.PatternJSONAccess:
		return &dialect.IndexRecommendation{
			Column:    col,
			IndexType: IndexTypeBTree,
			Expression: fmt.Sprintf("CREATE INDEX idx_%s_json ON %s ((CAST(%s->>'$.path' AS CHAR(255))));",
				safeName, table, col),
			Reason: fmt.Sprintf("JSON field access on '%s' benefits from a functional B-tree index on extracted JSON paths", col),
		}

	case dialect.PatternRegexMatch:
		return &dialect.IndexRecommendation{
			Column:    col,
			IndexType: IndexTypeFullText,
			Expression: fmt.Sprintf("CREATE FULLTEXT INDEX idx_%s_fulltext ON %s (%s);",
				safeName, table, col),
			Reason: fmt.Sprintf("Regex matching on '%s' may benefit from FULLTEXT index for text search patterns", col),
		}

	case dialect.PatternArrayMembership, dialect.PatternArrayComprehension:
		// MySQL does not have native array types; skip
		return nil

	case dialect.PatternJSONArrayComprehension:
		return &dialect.IndexRecommendation{
			Column:    col,
			IndexType: IndexTypeBTree,
			Expression: fmt.Sprintf("CREATE INDEX idx_%s_json ON %s ((CAST(%s->>'$.path' AS CHAR(255))));",
				safeName, table, col),
			Reason: fmt.Sprintf("JSON array operations on '%s' may benefit from a functional index on extracted JSON values", col),
		}
	}

	return nil
}

// SupportedPatterns returns the pattern types supported by MySQL.
func (d *Dialect) SupportedPatterns() []dialect.PatternType {
	return []dialect.PatternType{
		dialect.PatternComparison,
		dialect.PatternJSONAccess,
		dialect.PatternRegexMatch,
		dialect.PatternJSONArrayComprehension,
	}
}

// sanitizeIndexName creates a safe index name from a column name.
func sanitizeIndexName(column string) string {
	sanitized := strings.ReplaceAll(column, ".", "_")
	sanitized = strings.ReplaceAll(sanitized, " ", "_")
	sanitized = strings.ReplaceAll(sanitized, "-", "_")

	// MySQL index names are limited to 64 characters
	if len(sanitized) > 50 {
		sanitized = sanitized[:50]
	}

	return sanitized
}
