package duckdb

import (
	"fmt"
	"strings"

	"github.com/spandigital/cel2sql/v3/dialect"
)

// DuckDB index type constants.
const (
	IndexTypeART = "ART"
)

// RecommendIndex generates a DuckDB-specific index recommendation for the given pattern.
// DuckDB uses ART (Adaptive Radix Tree) indexes by default. Returns nil for unsupported patterns.
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
			IndexType: IndexTypeART,
			Expression: fmt.Sprintf("CREATE INDEX idx_%s ON %s (%s);",
				safeName, table, col),
			Reason: fmt.Sprintf("Comparison operations on '%s' benefit from an ART index for efficient range queries and equality checks", col),
		}

	case dialect.PatternJSONAccess:
		return &dialect.IndexRecommendation{
			Column:    col,
			IndexType: IndexTypeART,
			Expression: fmt.Sprintf("CREATE INDEX idx_%s_json ON %s (%s);",
				safeName, table, col),
			Reason: fmt.Sprintf("JSON field access on '%s' may benefit from an ART index", col),
		}

	case dialect.PatternRegexMatch:
		// DuckDB does not have specialized regex indexes
		return nil

	case dialect.PatternArrayMembership, dialect.PatternArrayComprehension:
		return &dialect.IndexRecommendation{
			Column:    col,
			IndexType: IndexTypeART,
			Expression: fmt.Sprintf("CREATE INDEX idx_%s ON %s (%s);",
				safeName, table, col),
			Reason: fmt.Sprintf("Array operations on '%s' may benefit from an ART index", col),
		}

	case dialect.PatternJSONArrayComprehension:
		return &dialect.IndexRecommendation{
			Column:    col,
			IndexType: IndexTypeART,
			Expression: fmt.Sprintf("CREATE INDEX idx_%s_json ON %s (%s);",
				safeName, table, col),
			Reason: fmt.Sprintf("JSON array comprehension on '%s' may benefit from an ART index", col),
		}
	}

	return nil
}

// SupportedPatterns returns the pattern types supported by DuckDB.
func (d *Dialect) SupportedPatterns() []dialect.PatternType {
	return []dialect.PatternType{
		dialect.PatternComparison,
		dialect.PatternJSONAccess,
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

	if len(sanitized) > 50 {
		sanitized = sanitized[:50]
	}

	return sanitized
}
