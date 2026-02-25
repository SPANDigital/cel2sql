// Package dialect defines the IndexAdvisor interface for dialect-specific index recommendations.
package dialect

// PatternType enumerates detected index-worthy query patterns.
type PatternType int

// Index-worthy pattern types detected during query analysis.
const (
	PatternComparison             PatternType = iota // Equality/range comparisons (==, >, <, >=, <=)
	PatternJSONAccess                                // JSON/JSONB field access
	PatternRegexMatch                                // Regex pattern matching
	PatternArrayMembership                           // Array IN/containment
	PatternArrayComprehension                        // Array comprehension (all, exists, filter, map)
	PatternJSONArrayComprehension                    // JSON array comprehension
)

// IndexPattern describes a detected query pattern that could benefit from indexing.
type IndexPattern struct {
	// Column is the full column name (e.g., "person.metadata").
	Column string

	// Pattern is the type of query pattern detected.
	Pattern PatternType

	// TableHint is an optional table name hint for generating CREATE INDEX statements.
	// If empty, "table_name" is used as the default placeholder.
	TableHint string
}

// IndexRecommendation represents a database index recommendation.
// It provides actionable guidance for optimizing query performance.
type IndexRecommendation struct {
	// Column is the database column that should be indexed.
	Column string

	// IndexType specifies the index type (e.g., "BTREE", "GIN", "ART", "CLUSTERING").
	IndexType string

	// Expression is the complete DDL statement that can be executed directly.
	Expression string

	// Reason explains why this index is recommended and what query patterns it optimizes.
	Reason string
}

// IndexAdvisor generates dialect-specific index recommendations.
// Dialects that support index analysis implement this interface on their Dialect struct.
type IndexAdvisor interface {
	// RecommendIndex generates an IndexRecommendation for the given pattern,
	// or returns nil if the dialect has no applicable index for this pattern.
	RecommendIndex(pattern IndexPattern) *IndexRecommendation

	// SupportedPatterns returns which PatternTypes this advisor can handle.
	SupportedPatterns() []PatternType
}

// GetIndexAdvisor returns the IndexAdvisor for a dialect, if it implements the interface.
func GetIndexAdvisor(d Dialect) (IndexAdvisor, bool) {
	advisor, ok := d.(IndexAdvisor)
	return advisor, ok
}
