// Package cel2sql provides query analysis and index recommendations
package cel2sql

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/common/overloads"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

// Index type constants for recommendations
const (
	// IndexTypeBTree represents a B-tree index for efficient range queries and equality checks
	IndexTypeBTree = "BTREE"
	// IndexTypeGIN represents a GIN (Generalized Inverted Index) for JSON, arrays, and full-text search
	IndexTypeGIN = "GIN"
	// IndexTypeGIST represents a GIST (Generalized Search Tree) index
	IndexTypeGIST = "GIST"
)

// IndexRecommendation represents a database index recommendation based on CEL query patterns.
// It provides actionable guidance for optimizing query performance through appropriate indexing strategies.
type IndexRecommendation struct {
	// Column is the database column that should be indexed
	Column string

	// IndexType specifies the PostgreSQL index type (e.g., "BTREE", "GIN", "GIST")
	IndexType string

	// Expression is the complete CREATE INDEX statement that can be executed directly
	Expression string

	// Reason explains why this index is recommended and what query patterns it optimizes
	Reason string
}

// analysisConverter extends converter with index pattern tracking capabilities
type analysisConverter struct {
	*converter
	recommendations map[string]*IndexRecommendation // Key: column name, Value: recommendation
	visitedColumns  map[string]bool                 // Track which columns have been accessed
}

// AnalyzeQuery converts a CEL AST to PostgreSQL SQL and provides index recommendations.
// It analyzes the query patterns to suggest indexes that would optimize performance.
//
// The function detects patterns that benefit from specific index types:
//   - JSON/JSONB path operations (->>, ?) → GIN indexes
//   - Array operations (UNNEST, comprehensions) → GIN indexes
//   - Regex matching (matches()) → GIN indexes with pg_trgm extension
//   - Frequently accessed fields in comparisons → B-tree indexes
//
// Example:
//
//	sql, recommendations, err := cel2sql.AnalyzeQuery(ast,
//	    cel2sql.WithSchemas(schemas))
//	if err != nil {
//	    return err
//	}
//
//	// Use the SQL query
//	rows, err := db.Query("SELECT * FROM table WHERE " + sql)
//
//	// Apply index recommendations
//	for _, rec := range recommendations {
//	    fmt.Printf("Recommendation: %s\n", rec.Reason)
//	    fmt.Printf("Execute: %s\n\n", rec.Expression)
//	}
func AnalyzeQuery(ast *cel.Ast, opts ...ConvertOption) (string, []IndexRecommendation, error) {
	start := time.Now()

	// Parse options
	options := &convertOptions{
		logger:       slog.New(slog.DiscardHandler), // Default: no-op logger with zero overhead
		maxDepth:     defaultMaxRecursionDepth,
		maxOutputLen: defaultMaxSQLOutputLength,
	}
	for _, opt := range opts {
		opt(options)
	}

	// Convert AST to CheckedExpr
	checkedExpr, err := cel.AstToCheckedExpr(ast)
	if err != nil {
		return "", nil, fmt.Errorf("failed to convert AST to CheckedExpr: %w", err)
	}

	// Pass 1: Analyze the AST to collect index recommendations
	baseConverter := &converter{
		typeMap:      checkedExpr.TypeMap,
		schemas:      options.schemas,
		ctx:          options.ctx,
		logger:       options.logger,
		maxDepth:     options.maxDepth,
		maxOutputLen: options.maxOutputLen,
	}

	analyzer := &analysisConverter{
		converter:       baseConverter,
		recommendations: make(map[string]*IndexRecommendation),
		visitedColumns:  make(map[string]bool),
	}

	// Analyze the expression tree to collect index patterns
	if err := analyzer.analyzeVisit(checkedExpr.Expr); err != nil {
		return "", nil, fmt.Errorf("analysis failed: %w", err)
	}

	// Pass 2: Generate SQL using the standard conversion
	sql, err := Convert(ast, opts...)
	if err != nil {
		return "", nil, fmt.Errorf("SQL generation failed: %w", err)
	}
	duration := time.Since(start)

	if options.logger != nil {
		options.logger.Debug("query analysis completed",
			"sql", sql,
			"recommendation_count", len(analyzer.recommendations),
			"duration", duration)
	}

	// Convert recommendations map to slice
	recommendations := make([]IndexRecommendation, 0, len(analyzer.recommendations))
	for _, rec := range analyzer.recommendations {
		recommendations = append(recommendations, *rec)
	}

	return sql, recommendations, nil
}

// analyzeVisit recursively walks the AST to track index-worthy patterns
func (a *analysisConverter) analyzeVisit(expr *exprpb.Expr) error {
	if expr == nil {
		return nil
	}

	// Track patterns based on expression type
	switch expr.ExprKind.(type) {
	case *exprpb.Expr_CallExpr:
		if err := a.analyzeCall(expr); err != nil {
			return err
		}
		// Recursively analyze call arguments
		c := expr.GetCallExpr()
		if c.GetTarget() != nil {
			if err := a.analyzeVisit(c.GetTarget()); err != nil {
				return err
			}
		}
		for _, arg := range c.GetArgs() {
			if err := a.analyzeVisit(arg); err != nil {
				return err
			}
		}

	case *exprpb.Expr_ComprehensionExpr:
		if err := a.analyzeComprehension(expr); err != nil {
			return err
		}
		// Recursively analyze comprehension parts
		comp := expr.GetComprehensionExpr()
		if err := a.analyzeVisit(comp.GetIterRange()); err != nil {
			return err
		}
		if err := a.analyzeVisit(comp.GetAccuInit()); err != nil {
			return err
		}
		if err := a.analyzeVisit(comp.GetLoopCondition()); err != nil {
			return err
		}
		if err := a.analyzeVisit(comp.GetLoopStep()); err != nil {
			return err
		}
		if err := a.analyzeVisit(comp.GetResult()); err != nil {
			return err
		}

	case *exprpb.Expr_SelectExpr:
		if err := a.analyzeSelect(expr); err != nil {
			return err
		}
		// Recursively analyze operand
		sel := expr.GetSelectExpr()
		if err := a.analyzeVisit(sel.GetOperand()); err != nil {
			return err
		}

	case *exprpb.Expr_ListExpr:
		// Recursively analyze list elements
		list := expr.GetListExpr()
		for _, elem := range list.GetElements() {
			if err := a.analyzeVisit(elem); err != nil {
				return err
			}
		}

	case *exprpb.Expr_StructExpr:
		// Recursively analyze struct entries
		st := expr.GetStructExpr()
		for _, entry := range st.GetEntries() {
			if entry.GetMapKey() != nil {
				if err := a.analyzeVisit(entry.GetMapKey()); err != nil {
					return err
				}
			}
			if err := a.analyzeVisit(entry.GetValue()); err != nil {
				return err
			}
		}
	}

	return nil
}

// analyzeCall analyzes function calls for index recommendations
func (a *analysisConverter) analyzeCall(expr *exprpb.Expr) error {
	c := expr.GetCallExpr()
	fun := c.GetFunction()

	switch fun {
	case overloads.Matches:
		// Regex matching benefits from GIN index with pg_trgm extension
		if err := a.recommendRegexIndex(expr); err != nil {
			return err
		}

	case operators.Equals, operators.NotEquals,
		operators.Greater, operators.GreaterEquals,
		operators.Less, operators.LessEquals:
		// Comparison operations benefit from B-tree indexes
		if err := a.recommendComparisonIndex(expr); err != nil {
			return err
		}

	case operators.In:
		// IN operations on arrays benefit from GIN indexes
		if err := a.recommendArrayIndex(expr); err != nil {
			return err
		}
	}

	return nil
}

// analyzeComprehension analyzes comprehension expressions for index recommendations
func (a *analysisConverter) analyzeComprehension(expr *exprpb.Expr) error {
	comprehension := expr.GetComprehensionExpr()
	if comprehension == nil {
		return nil
	}

	iterRange := comprehension.GetIterRange()

	// Check if this is a JSON array comprehension
	if a.isJSONArrayField(iterRange) {
		// Extract the column name from the iter range
		if column := a.extractColumnName(iterRange); column != "" {
			a.addRecommendation(column, &IndexRecommendation{
				Column:    column,
				IndexType: IndexTypeGIN,
				Expression: fmt.Sprintf("CREATE INDEX idx_%s_gin ON table_name USING GIN (%s);",
					sanitizeIndexName(column), column),
				Reason: fmt.Sprintf("JSONB array comprehension on '%s' benefits from GIN index for efficient array element access", column),
			})
		}
	} else {
		// Regular array comprehension
		if column := a.extractColumnName(iterRange); column != "" {
			a.addRecommendation(column, &IndexRecommendation{
				Column:    column,
				IndexType: IndexTypeGIN,
				Expression: fmt.Sprintf("CREATE INDEX idx_%s_gin ON table_name USING GIN (%s);",
					sanitizeIndexName(column), column),
				Reason: fmt.Sprintf("Array comprehension on '%s' benefits from GIN index for efficient array operations", column),
			})
		}
	}

	return nil
}

// analyzeSelect analyzes field selection for index recommendations
func (a *analysisConverter) analyzeSelect(expr *exprpb.Expr) error {
	sel := expr.GetSelectExpr()
	fieldName := sel.GetField()
	operand := sel.GetOperand()

	// Check if the operand is a JSON field (e.g., person.metadata in person.metadata.name)
	if operandSel := operand.GetSelectExpr(); operandSel != nil {
		operandField := operandSel.GetField()
		if identExpr := operandSel.GetOperand().GetIdentExpr(); identExpr != nil {
			tableName := identExpr.GetName()
			// Check if the parent field is JSON
			if a.isFieldJSON(tableName, operandField) {
				column := tableName + "." + operandField
				a.addRecommendation(column, &IndexRecommendation{
					Column:    column,
					IndexType: IndexTypeGIN,
					Expression: fmt.Sprintf("CREATE INDEX idx_%s_gin ON table_name USING GIN (%s);",
						sanitizeIndexName(column), column),
					Reason: fmt.Sprintf("JSON path operations on '%s' benefit from GIN index for efficient nested field access", column),
				})
			}
		}
	}

	// Also check if this field itself is JSON
	if identExpr := operand.GetIdentExpr(); identExpr != nil {
		tableName := identExpr.GetName()
		if a.isFieldJSON(tableName, fieldName) {
			column := tableName + "." + fieldName
			a.addRecommendation(column, &IndexRecommendation{
				Column:    column,
				IndexType: IndexTypeGIN,
				Expression: fmt.Sprintf("CREATE INDEX idx_%s_gin ON table_name USING GIN (%s);",
					sanitizeIndexName(column), column),
				Reason: fmt.Sprintf("JSON field '%s' benefits from GIN index for efficient access", column),
			})
		}
		// Track column access for potential B-tree indexes
		fullColumn := tableName + "." + fieldName
		a.visitedColumns[fullColumn] = true
	}

	return nil
}

// recommendRegexIndex recommends a GIN index with pg_trgm for regex operations
func (a *analysisConverter) recommendRegexIndex(expr *exprpb.Expr) error {
	c := expr.GetCallExpr()
	target := c.GetTarget()

	if target == nil && len(c.GetArgs()) >= 2 {
		target = c.GetArgs()[0]
	}

	if target != nil {
		if column := a.extractColumnName(target); column != "" {
			a.addRecommendation(column, &IndexRecommendation{
				Column:    column,
				IndexType: IndexTypeGIN,
				Expression: fmt.Sprintf("CREATE INDEX idx_%s_gin_trgm ON table_name USING GIN (%s gin_trgm_ops);",
					sanitizeIndexName(column), column),
				Reason: fmt.Sprintf("Regex matching on '%s' benefits from GIN index with pg_trgm extension for pattern matching", column),
			})
		}
	}

	return nil
}

// recommendComparisonIndex recommends a B-tree index for comparison operations
func (a *analysisConverter) recommendComparisonIndex(expr *exprpb.Expr) error {
	c := expr.GetCallExpr()
	args := c.GetArgs()

	if len(args) < 2 {
		return nil
	}

	lhs := args[0]

	// Extract column from left-hand side
	if column := a.extractColumnName(lhs); column != "" {
		// Check if this is a JSON field (skip B-tree recommendation for JSON)
		if !a.isJSONField(lhs) {
			a.addRecommendation(column, &IndexRecommendation{
				Column:    column,
				IndexType: IndexTypeBTree,
				Expression: fmt.Sprintf("CREATE INDEX idx_%s_btree ON table_name (%s);",
					sanitizeIndexName(column), column),
				Reason: fmt.Sprintf("Comparison operations on '%s' benefit from B-tree index for efficient range queries and equality checks", column),
			})
		}
	}

	return nil
}

// recommendArrayIndex recommends a GIN index for array containment operations
func (a *analysisConverter) recommendArrayIndex(expr *exprpb.Expr) error {
	c := expr.GetCallExpr()
	args := c.GetArgs()

	if len(args) < 2 {
		return nil
	}

	rhs := args[1]

	// Check if the right-hand side is an array field
	if a.isFieldArray(a.extractTableName(rhs), a.extractFieldName(rhs)) {
		if column := a.extractColumnName(rhs); column != "" {
			a.addRecommendation(column, &IndexRecommendation{
				Column:    column,
				IndexType: IndexTypeGIN,
				Expression: fmt.Sprintf("CREATE INDEX idx_%s_gin ON table_name USING GIN (%s);",
					sanitizeIndexName(column), column),
				Reason: fmt.Sprintf("Array membership tests on '%s' benefit from GIN index for efficient element lookups", column),
			})
		}
	}

	return nil
}

// extractColumnName extracts the full column name (table.column) from an expression
func (a *analysisConverter) extractColumnName(expr *exprpb.Expr) string {
	if sel := expr.GetSelectExpr(); sel != nil {
		fieldName := sel.GetField()
		if identExpr := sel.GetOperand().GetIdentExpr(); identExpr != nil {
			tableName := identExpr.GetName()
			return tableName + "." + fieldName
		}
	}

	if identExpr := expr.GetIdentExpr(); identExpr != nil {
		return identExpr.GetName()
	}

	return ""
}

// extractTableName extracts the table name from an expression
func (a *analysisConverter) extractTableName(expr *exprpb.Expr) string {
	if sel := expr.GetSelectExpr(); sel != nil {
		if identExpr := sel.GetOperand().GetIdentExpr(); identExpr != nil {
			return identExpr.GetName()
		}
	}
	return ""
}

// extractFieldName extracts the field name from an expression
func (a *analysisConverter) extractFieldName(expr *exprpb.Expr) string {
	if sel := expr.GetSelectExpr(); sel != nil {
		return sel.GetField()
	}
	return ""
}

// isJSONField checks if an expression refers to a JSON/JSONB field
func (a *analysisConverter) isJSONField(expr *exprpb.Expr) bool {
	if sel := expr.GetSelectExpr(); sel != nil {
		operand := sel.GetOperand()
		fieldName := sel.GetField()
		if identExpr := operand.GetIdentExpr(); identExpr != nil {
			tableName := identExpr.GetName()
			return a.isFieldJSON(tableName, fieldName)
		}
	}
	return false
}

// addRecommendation adds or updates an index recommendation
func (a *analysisConverter) addRecommendation(column string, rec *IndexRecommendation) {
	// Only add if we don't already have a recommendation for this column
	// or if the new recommendation is more specific (e.g., GIN over BTREE)
	existing, exists := a.recommendations[column]
	if !exists {
		a.recommendations[column] = rec
		return
	}

	// GIN indexes are more versatile than BTREE for JSON/array operations
	// If we already have a BTREE recommendation and we're suggesting GIN, upgrade it
	if existing.IndexType == IndexTypeBTree && rec.IndexType == IndexTypeGIN {
		a.recommendations[column] = rec
	}
}

// sanitizeIndexName creates a safe index name from a column name
func sanitizeIndexName(column string) string {
	// Replace dots and special characters with underscores
	sanitized := strings.ReplaceAll(column, ".", "_")
	sanitized = strings.ReplaceAll(sanitized, " ", "_")
	sanitized = strings.ReplaceAll(sanitized, "-", "_")

	// PostgreSQL index names are limited to 63 characters
	if len(sanitized) > 50 {
		sanitized = sanitized[:50]
	}

	return sanitized
}
