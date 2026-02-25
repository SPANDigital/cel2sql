package bigquery

import (
	"fmt"
	"strings"

	"github.com/spandigital/cel2sql/v3/dialect"
)

// BigQuery index type constants.
const (
	IndexTypeClustering  = "CLUSTERING"
	IndexTypeSearchIndex = "SEARCH_INDEX"
)

// RecommendIndex generates a BigQuery-specific index recommendation for the given pattern.
// BigQuery uses clustering keys and search indexes. Returns nil for unsupported patterns.
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
			IndexType: IndexTypeClustering,
			Expression: fmt.Sprintf("ALTER TABLE %s SET OPTIONS (clustering_columns=['%s']);",
				table, col),
			Reason: fmt.Sprintf("Comparison operations on '%s' benefit from clustering for efficient partition pruning and range scans", col),
		}

	case dialect.PatternJSONAccess:
		return &dialect.IndexRecommendation{
			Column:    col,
			IndexType: IndexTypeSearchIndex,
			Expression: fmt.Sprintf("CREATE SEARCH INDEX idx_%s ON %s (%s);",
				safeName, table, col),
			Reason: fmt.Sprintf("JSON field access on '%s' benefits from a search index for efficient nested field lookups", col),
		}

	case dialect.PatternRegexMatch:
		// BigQuery does not have specialized regex indexes
		return nil

	case dialect.PatternArrayMembership, dialect.PatternArrayComprehension:
		// BigQuery arrays do not benefit from standalone indexes
		return nil

	case dialect.PatternJSONArrayComprehension:
		return &dialect.IndexRecommendation{
			Column:    col,
			IndexType: IndexTypeSearchIndex,
			Expression: fmt.Sprintf("CREATE SEARCH INDEX idx_%s ON %s (%s);",
				safeName, table, col),
			Reason: fmt.Sprintf("JSON array operations on '%s' may benefit from a search index", col),
		}
	}

	return nil
}

// SupportedPatterns returns the pattern types supported by BigQuery.
func (d *Dialect) SupportedPatterns() []dialect.PatternType {
	return []dialect.PatternType{
		dialect.PatternComparison,
		dialect.PatternJSONAccess,
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
