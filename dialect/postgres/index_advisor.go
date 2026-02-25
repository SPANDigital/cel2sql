package postgres

import (
	"fmt"
	"strings"

	"github.com/spandigital/cel2sql/v3/dialect"
)

// PostgreSQL index type constants.
const (
	IndexTypeBTree = "BTREE"
	IndexTypeGIN   = "GIN"
	IndexTypeGIST  = "GIST"
)

// RecommendIndex generates a PostgreSQL-specific index recommendation for the given pattern.
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
			IndexType: IndexTypeGIN,
			Expression: fmt.Sprintf("CREATE INDEX idx_%s_gin ON %s USING GIN (%s);",
				safeName, table, col),
			Reason: fmt.Sprintf("JSON path operations on '%s' benefit from GIN index for efficient nested field access", col),
		}

	case dialect.PatternRegexMatch:
		return &dialect.IndexRecommendation{
			Column:    col,
			IndexType: IndexTypeGIN,
			Expression: fmt.Sprintf("CREATE INDEX idx_%s_gin_trgm ON %s USING GIN (%s gin_trgm_ops);",
				safeName, table, col),
			Reason: fmt.Sprintf("Regex matching on '%s' benefits from GIN index with pg_trgm extension for pattern matching", col),
		}

	case dialect.PatternArrayMembership:
		return &dialect.IndexRecommendation{
			Column:    col,
			IndexType: IndexTypeGIN,
			Expression: fmt.Sprintf("CREATE INDEX idx_%s_gin ON %s USING GIN (%s);",
				safeName, table, col),
			Reason: fmt.Sprintf("Array membership tests on '%s' benefit from GIN index for efficient element lookups", col),
		}

	case dialect.PatternArrayComprehension:
		return &dialect.IndexRecommendation{
			Column:    col,
			IndexType: IndexTypeGIN,
			Expression: fmt.Sprintf("CREATE INDEX idx_%s_gin ON %s USING GIN (%s);",
				safeName, table, col),
			Reason: fmt.Sprintf("Array comprehension on '%s' benefits from GIN index for efficient array operations", col),
		}

	case dialect.PatternJSONArrayComprehension:
		return &dialect.IndexRecommendation{
			Column:    col,
			IndexType: IndexTypeGIN,
			Expression: fmt.Sprintf("CREATE INDEX idx_%s_gin ON %s USING GIN (%s);",
				safeName, table, col),
			Reason: fmt.Sprintf("JSONB array comprehension on '%s' benefits from GIN index for efficient array element access", col),
		}
	}

	return nil
}

// SupportedPatterns returns all pattern types supported by PostgreSQL.
func (d *Dialect) SupportedPatterns() []dialect.PatternType {
	return []dialect.PatternType{
		dialect.PatternComparison,
		dialect.PatternJSONAccess,
		dialect.PatternRegexMatch,
		dialect.PatternArrayMembership,
		dialect.PatternArrayComprehension,
		dialect.PatternJSONArrayComprehension,
	}
}

// sanitizeIndexName creates a safe index name from a column name.
func sanitizeIndexName(column string) string {
	sanitized := strings.ReplaceAll(column, ".", "_")
	sanitized = strings.ReplaceAll(sanitized, " ", "_")
	sanitized = strings.ReplaceAll(sanitized, "-", "_")

	// PostgreSQL index names are limited to 63 characters
	if len(sanitized) > 50 {
		sanitized = sanitized[:50]
	}

	return sanitized
}
