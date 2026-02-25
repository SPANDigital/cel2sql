package sqlite

import (
	"fmt"
	"strings"

	"github.com/spandigital/cel2sql/v3/dialect"
)

// SQLite index type constants.
const (
	IndexTypeBTree = "BTREE"
)

// RecommendIndex generates a SQLite-specific index recommendation for the given pattern.
// SQLite only supports standard B-tree indexes. Returns nil for unsupported patterns.
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
			Expression: fmt.Sprintf("CREATE INDEX idx_%s ON %s (%s);",
				safeName, table, col),
			Reason: fmt.Sprintf("Comparison operations on '%s' benefit from an index for efficient range queries and equality checks", col),
		}

	case dialect.PatternJSONAccess:
		// SQLite does not support indexes on JSON expressions directly
		return nil

	case dialect.PatternRegexMatch:
		// SQLite does not support native regex; no index recommendation
		return nil

	case dialect.PatternArrayMembership, dialect.PatternArrayComprehension:
		// SQLite does not have native array types
		return nil

	case dialect.PatternJSONArrayComprehension:
		// SQLite does not support indexes on JSON array operations
		return nil
	}

	return nil
}

// SupportedPatterns returns the pattern types supported by SQLite.
func (d *Dialect) SupportedPatterns() []dialect.PatternType {
	return []dialect.PatternType{
		dialect.PatternComparison,
	}
}

// sanitizeIndexName creates a safe index name from a column name.
func sanitizeIndexName(column string) string {
	sanitized := strings.ReplaceAll(column, ".", "_")
	sanitized = strings.ReplaceAll(sanitized, " ", "_")
	sanitized = strings.ReplaceAll(sanitized, "-", "_")

	if len(sanitized) > 50 {
		sanitized = sanitized[:50]
	}

	return sanitized
}
