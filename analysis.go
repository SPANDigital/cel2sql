// Package cel2sql provides query analysis and index recommendations
package cel2sql

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/common/overloads"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"

	"github.com/spandigital/cel2sql/v3/dialect"
	"github.com/spandigital/cel2sql/v3/dialect/postgres"
)

// Index type constants for recommendations (kept for backward compatibility).
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

	// IndexType specifies the index type (e.g., "BTREE", "GIN", "ART", "CLUSTERING")
	IndexType string

	// Expression is the complete DDL statement that can be executed directly
	Expression string

	// Reason explains why this index is recommended and what query patterns it optimizes
	Reason string
}

// analysisConverter extends converter with index pattern tracking capabilities
type analysisConverter struct {
	*converter
	recommendations map[string]*IndexRecommendation // Key: column name, Value: recommendation
	visitedColumns  map[string]bool                 // Track which columns have been accessed
	advisor         dialect.IndexAdvisor            // Dialect-specific index advisor
}

// AnalyzeQuery converts a CEL AST to SQL and provides dialect-specific index recommendations.
// It analyzes the query patterns to suggest indexes that would optimize performance.
//
// The function detects patterns that benefit from specific index types:
//   - JSON/JSONB path operations → GIN indexes (PostgreSQL), functional indexes (MySQL), search indexes (BigQuery)
//   - Array operations → GIN indexes (PostgreSQL), ART indexes (DuckDB)
//   - Regex matching → GIN indexes with pg_trgm (PostgreSQL), FULLTEXT indexes (MySQL)
//   - Comparison operations → B-tree indexes (PostgreSQL/MySQL/SQLite), ART (DuckDB), clustering (BigQuery)
//
// Use WithDialect() to get dialect-specific index recommendations. Defaults to PostgreSQL.
//
// Example:
//
//	sql, recommendations, err := cel2sql.AnalyzeQuery(ast,
//	    cel2sql.WithSchemas(schemas),
//	    cel2sql.WithDialect(mysql.New()))
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

	// Default to PostgreSQL dialect if none specified
	d := options.dialect
	if d == nil {
		d = postgres.New()
	}

	// Get the IndexAdvisor for the dialect (all built-in dialects implement it)
	advisor, hasAdvisor := dialect.GetIndexAdvisor(d)
	if !hasAdvisor {
		// Fallback: use PostgreSQL advisor for backward compatibility
		advisor = postgres.New()
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
		advisor:         advisor,
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
			"dialect", d.Name(),
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
		// Regex matching benefits from dialect-specific indexes
		a.recommendRegexIndex(expr)

	case operators.Equals, operators.NotEquals,
		operators.Greater, operators.GreaterEquals,
		operators.Less, operators.LessEquals:
		// Comparison operations benefit from indexes
		a.recommendComparisonIndex(expr)

	case operators.In:
		// IN operations on arrays benefit from indexes
		a.recommendArrayIndex(expr)
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
		if column := a.extractColumnName(iterRange); column != "" {
			a.recommendForPattern(column, dialect.PatternJSONArrayComprehension)
		}
	} else {
		// Regular array comprehension
		if column := a.extractColumnName(iterRange); column != "" {
			a.recommendForPattern(column, dialect.PatternArrayComprehension)
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
				a.recommendForPattern(column, dialect.PatternJSONAccess)
			}
		}
	}

	// Also check if this field itself is JSON
	if identExpr := operand.GetIdentExpr(); identExpr != nil {
		tableName := identExpr.GetName()
		if a.isFieldJSON(tableName, fieldName) {
			column := tableName + "." + fieldName
			a.recommendForPattern(column, dialect.PatternJSONAccess)
		}
		// Track column access for potential indexes
		fullColumn := tableName + "." + fieldName
		a.visitedColumns[fullColumn] = true
	}

	return nil
}

// recommendRegexIndex recommends an index for regex operations
func (a *analysisConverter) recommendRegexIndex(expr *exprpb.Expr) {
	c := expr.GetCallExpr()
	target := c.GetTarget()

	if target == nil && len(c.GetArgs()) >= 2 {
		target = c.GetArgs()[0]
	}

	if target != nil {
		if column := a.extractColumnName(target); column != "" {
			a.recommendForPattern(column, dialect.PatternRegexMatch)
		}
	}
}

// recommendComparisonIndex recommends an index for comparison operations
func (a *analysisConverter) recommendComparisonIndex(expr *exprpb.Expr) {
	c := expr.GetCallExpr()
	args := c.GetArgs()

	if len(args) < 2 {
		return
	}

	lhs := args[0]

	// Extract column from left-hand side
	if column := a.extractColumnName(lhs); column != "" {
		// Check if this is a JSON field (skip comparison recommendation for JSON)
		if !a.isJSONField(lhs) {
			a.recommendForPattern(column, dialect.PatternComparison)
		}
	}
}

// recommendArrayIndex recommends an index for array containment operations
func (a *analysisConverter) recommendArrayIndex(expr *exprpb.Expr) {
	c := expr.GetCallExpr()
	args := c.GetArgs()

	if len(args) < 2 {
		return
	}

	rhs := args[1]

	// Check if the right-hand side is an array field
	if a.isFieldArray(a.extractTableName(rhs), a.extractFieldName(rhs)) {
		if column := a.extractColumnName(rhs); column != "" {
			a.recommendForPattern(column, dialect.PatternArrayMembership)
		}
	}
}

// recommendForPattern asks the dialect's IndexAdvisor for a recommendation and stores it.
func (a *analysisConverter) recommendForPattern(column string, pattern dialect.PatternType) {
	rec := a.advisor.RecommendIndex(dialect.IndexPattern{
		Column:  column,
		Pattern: pattern,
	})
	if rec == nil {
		return
	}
	a.addRecommendation(column, &IndexRecommendation{
		Column:     rec.Column,
		IndexType:  rec.IndexType,
		Expression: rec.Expression,
		Reason:     rec.Reason,
	})
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

// addRecommendation adds or updates an index recommendation.
// When a more specialized recommendation exists for a column, it takes priority.
func (a *analysisConverter) addRecommendation(column string, rec *IndexRecommendation) {
	// Only add if we don't already have a recommendation for this column
	// or if the new recommendation is more specific
	existing, exists := a.recommendations[column]
	if !exists {
		a.recommendations[column] = rec
		return
	}

	// More specialized index types take priority over basic B-tree/comparison indexes
	if isBasicIndexType(existing.IndexType) && !isBasicIndexType(rec.IndexType) {
		a.recommendations[column] = rec
	}
}

// isBasicIndexType returns true if the index type is a basic comparison index
// that should be upgraded when a more specialized recommendation is available.
func isBasicIndexType(indexType string) bool {
	return indexType == IndexTypeBTree || indexType == "ART" || indexType == "CLUSTERING"
}
