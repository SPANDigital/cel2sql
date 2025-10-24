// Package cel2sql converts CEL (Common Expression Language) expressions to PostgreSQL SQL conditions.
package cel2sql

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/common/overloads"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"

	"github.com/spandigital/cel2sql/v2/pg"
)

// Implementations based on `google/cel-go`'s unparser
// https://github.com/google/cel-go/blob/master/parser/unparser.go

// Regex pattern complexity limits to prevent ReDoS attacks (CWE-1333).
const (
	// maxRegexPatternLength is the maximum allowed length for regex patterns
	// to prevent processing extremely long patterns that could cause DoS.
	maxRegexPatternLength = 500

	// maxRegexGroups is the maximum number of capture groups allowed in a pattern
	// to prevent memory exhaustion and slow matching.
	maxRegexGroups = 20

	// maxRegexNestingDepth is the maximum nesting depth for groups and quantifiers
	// to prevent catastrophic backtracking.
	maxRegexNestingDepth = 10

	// defaultMaxRecursionDepth is the default maximum recursion depth for visit()
	// to prevent stack overflow from deeply nested expressions (CWE-674: Uncontrolled Recursion).
	defaultMaxRecursionDepth = 100

	// maxComprehensionDepth is the maximum nesting depth for CEL comprehensions
	// to prevent resource exhaustion from deeply nested UNNEST/subquery operations (CWE-400).
	maxComprehensionDepth = 3

	// defaultMaxSQLOutputLength is the default maximum length of generated SQL output
	// to prevent resource exhaustion from extremely large SQL queries (CWE-400).
	defaultMaxSQLOutputLength = 50000
)

// ConvertOption is a functional option for configuring the Convert function.
type ConvertOption func(*convertOptions)

// convertOptions holds configuration options for the Convert function.
type convertOptions struct {
	schemas       map[string]pg.Schema
	ctx           context.Context
	logger        *slog.Logger
	maxDepth      int // Maximum recursion depth (0 = use default)
	maxOutputLen  int // Maximum SQL output length (0 = use default)
}

// WithSchemas provides schema information for proper JSON/JSONB field handling.
// This option is required for correct SQL generation when using JSON/JSONB fields.
//
// Example:
//
//	schemas := provider.GetSchemas()
//	sql, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
func WithSchemas(schemas map[string]pg.Schema) ConvertOption {
	return func(o *convertOptions) {
		o.schemas = schemas
	}
}

// WithContext provides a context for cancellation and timeout support.
// If not provided, operations will run without cancellation checks.
// This allows long-running conversions to be cancelled and enables timeout protection.
//
// Example with timeout:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	sql, err := cel2sql.Convert(ast, cel2sql.WithContext(ctx))
//
// Example with cancellation:
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//	sql, err := cel2sql.Convert(ast, cel2sql.WithContext(ctx), cel2sql.WithSchemas(schemas))
func WithContext(ctx context.Context) ConvertOption {
	return func(o *convertOptions) {
		o.ctx = ctx
	}
}

// WithLogger provides a logger for observability and debugging.
// If not provided, logging is disabled with zero overhead using slog.DiscardHandler.
//
// The logger enables visibility into:
//   - JSON path detection decisions (table, field, operator selection)
//   - Comprehension type identification (all, exists, filter, map)
//   - Schema lookups (hits/misses, field types)
//   - Performance metrics (conversion duration)
//   - Regex pattern transformations (RE2 to POSIX)
//   - Operator mapping decisions
//   - Error contexts with full details
//
// Example with JSON output:
//
//	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
//	sql, err := cel2sql.Convert(ast, cel2sql.WithLogger(logger))
//
// Example with text output:
//
//	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
//	sql, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas), cel2sql.WithLogger(logger))
func WithLogger(logger *slog.Logger) ConvertOption {
	return func(o *convertOptions) {
		o.logger = logger
	}
}

// WithMaxDepth sets the maximum recursion depth for expression traversal.
// If not provided, defaultMaxRecursionDepth (100) is used.
// This protects against stack overflow from deeply nested expressions (CWE-674).
//
// Example with custom depth:
//
//	sql, err := cel2sql.Convert(ast, cel2sql.WithMaxDepth(150))
//
// Example with multiple options:
//
//	sql, err := cel2sql.Convert(ast,
//	    cel2sql.WithMaxDepth(50),
//	    cel2sql.WithContext(ctx),
//	    cel2sql.WithSchemas(schemas))
func WithMaxDepth(maxDepth int) ConvertOption {
	return func(o *convertOptions) {
		o.maxDepth = maxDepth
	}
}

// WithMaxOutputLength sets the maximum length of generated SQL output.
// If not provided, defaultMaxSQLOutputLength (50000) is used.
// This protects against resource exhaustion from extremely large SQL queries (CWE-400).
//
// Example with custom output length limit:
//
//	sql, err := cel2sql.Convert(ast, cel2sql.WithMaxOutputLength(100000))
//
// Example with multiple options:
//
//	sql, err := cel2sql.Convert(ast,
//	    cel2sql.WithMaxOutputLength(25000),
//	    cel2sql.WithMaxDepth(50),
//	    cel2sql.WithContext(ctx))
func WithMaxOutputLength(maxLength int) ConvertOption {
	return func(o *convertOptions) {
		o.maxOutputLen = maxLength
	}
}

// Result represents the output of a CEL to SQL conversion with parameterized queries.
// It contains the SQL string with placeholders ($1, $2, etc.) and the corresponding parameter values.
type Result struct {
	SQL        string        // The generated SQL WHERE clause with placeholders
	Parameters []interface{} // Parameter values in order ($1, $2, etc.)
}

// Convert converts a CEL AST to a PostgreSQL SQL WHERE clause condition.
// Options can be provided to configure the conversion behavior.
//
// Example without options:
//
//	sql, err := cel2sql.Convert(ast)
//
// Example with schema information for JSON/JSONB support:
//
//	sql, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
func Convert(ast *cel.Ast, opts ...ConvertOption) (string, error) {
	start := time.Now()

	options := &convertOptions{
		logger:       slog.New(slog.DiscardHandler), // Default: no-op logger with zero overhead
		maxDepth:     defaultMaxRecursionDepth,      // Default: 100 recursion depth limit
		maxOutputLen: defaultMaxSQLOutputLength,     // Default: 50000 character output limit
	}
	for _, opt := range opts {
		opt(options)
	}

	options.logger.Debug("starting CEL to SQL conversion")

	checkedExpr, err := cel.AstToCheckedExpr(ast)
	if err != nil {
		options.logger.Error("AST to CheckedExpr conversion failed", slog.Any("error", err))
		return "", err
	}

	un := &converter{
		typeMap:       checkedExpr.TypeMap,
		schemas:       options.schemas,
		ctx:           options.ctx,
		logger:        options.logger,
		maxDepth:      options.maxDepth,
		maxOutputLen:  options.maxOutputLen,
	}

	if err := un.visit(checkedExpr.Expr); err != nil {
		options.logger.Error("conversion failed", slog.Any("error", err))
		return "", err
	}

	result := un.str.String()
	duration := time.Since(start)

	options.logger.LogAttrs(context.Background(), slog.LevelDebug,
		"conversion completed",
		slog.String("sql", result),
		slog.Duration("duration", duration),
	)

	return result, nil
}

// ConvertParameterized converts a CEL AST to a parameterized PostgreSQL SQL WHERE clause.
// Returns both the SQL string with placeholders ($1, $2, etc.) and the parameter values.
// This enables query plan caching and provides additional SQL injection protection.
//
// Constants that are parameterized:
//   - String literals: 'John' → $1
//   - Numeric literals: 42, 3.14 → $1, $2
//   - Byte literals: b"data" → $1
//
// Constants kept inline (for query plan optimization):
//   - TRUE, FALSE (boolean constants)
//   - NULL
//
// Example:
//
//	result, err := cel2sql.ConvertParameterized(ast,
//	    cel2sql.WithSchemas(schemas),
//	    cel2sql.WithContext(ctx))
//	// result.SQL: "user.age = $1 AND user.name = $2"
//	// result.Parameters: []interface{}{18, "John"}
//
//	// Execute with database/sql
//	rows, err := db.Query("SELECT * FROM users WHERE "+result.SQL, result.Parameters...)
func ConvertParameterized(ast *cel.Ast, opts ...ConvertOption) (*Result, error) {
	start := time.Now()

	options := &convertOptions{
		logger:       slog.New(slog.DiscardHandler), // Default: no-op logger with zero overhead
		maxDepth:     defaultMaxRecursionDepth,      // Default: 100 recursion depth limit
		maxOutputLen: defaultMaxSQLOutputLength,     // Default: 50000 character output limit
	}
	for _, opt := range opts {
		opt(options)
	}

	options.logger.Debug("starting parameterized CEL to SQL conversion")

	checkedExpr, err := cel.AstToCheckedExpr(ast)
	if err != nil {
		options.logger.Error("AST to CheckedExpr conversion failed", slog.Any("error", err))
		return nil, err
	}

	un := &converter{
		typeMap:       checkedExpr.TypeMap,
		schemas:       options.schemas,
		ctx:           options.ctx,
		logger:        options.logger,
		maxDepth:      options.maxDepth,
		maxOutputLen:  options.maxOutputLen,
		parameterize:  true, // Enable parameterization
	}

	if err := un.visit(checkedExpr.Expr); err != nil {
		options.logger.Error("conversion failed", slog.Any("error", err))
		return nil, err
	}

	sql := un.str.String()
	duration := time.Since(start)

	options.logger.LogAttrs(context.Background(), slog.LevelDebug,
		"parameterized conversion completed",
		slog.String("sql", sql),
		slog.Int("param_count", len(un.parameters)),
		slog.Duration("duration", duration),
	)

	return &Result{
		SQL:        sql,
		Parameters: un.parameters,
	}, nil
}

type converter struct {
	str                 strings.Builder
	typeMap             map[int64]*exprpb.Type
	schemas             map[string]pg.Schema
	ctx                 context.Context
	logger              *slog.Logger
	depth               int           // Current recursion depth
	maxDepth            int           // Maximum allowed recursion depth
	maxOutputLen        int           // Maximum allowed SQL output length
	comprehensionDepth  int           // Current comprehension nesting depth
	parameterize        bool          // Enable parameterized output
	parameters          []interface{} // Collected parameters for parameterized queries
	paramCount          int           // Parameter counter for placeholders ($1, $2, etc.)
}

// checkContext checks if the context has been cancelled or expired.
// Returns nil if no context was provided or if the context is still active.
// Returns an error if the context has been cancelled or its deadline has exceeded.
func (con *converter) checkContext() error {
	if con.ctx == nil {
		return nil
	}
	if err := con.ctx.Err(); err != nil {
		return fmt.Errorf("conversion cancelled: %w", err)
	}
	return nil
}

func (con *converter) visit(expr *exprpb.Expr) error {
	// Track recursion depth
	con.depth++
	defer func() { con.depth-- }()

	// Check depth limit before context check (fail fast)
	// Allow depths up to and including maxDepth
	if con.depth > con.maxDepth {
		return fmt.Errorf("expression exceeds maximum recursion depth of %d", con.maxDepth)
	}

	// Check for context cancellation at the main recursion entry point
	if err := con.checkContext(); err != nil {
		return err
	}

	// Check SQL output length limit to prevent resource exhaustion (CWE-400)
	if con.str.Len() > con.maxOutputLen {
		return fmt.Errorf("generated SQL exceeds maximum output length of %d", con.maxOutputLen)
	}

	switch expr.ExprKind.(type) {
	case *exprpb.Expr_CallExpr:
		return con.visitCall(expr)
	// Comprehensions are supported (all, exists, exists_one, filter, map).
	case *exprpb.Expr_ComprehensionExpr:
		return con.visitComprehension(expr)
	case *exprpb.Expr_ConstExpr:
		return con.visitConst(expr)
	case *exprpb.Expr_IdentExpr:
		return con.visitIdent(expr)
	case *exprpb.Expr_ListExpr:
		return con.visitList(expr)
	case *exprpb.Expr_SelectExpr:
		return con.visitSelect(expr)
	case *exprpb.Expr_StructExpr:
		return con.visitStruct(expr)
	}
	return newConversionErrorf(errMsgUnsupportedExpression, "expr type: %T, id: %d", expr.ExprKind, expr.Id)
}

// isFieldJSON checks if a field in a table is a JSON/JSONB type using schema information
func (con *converter) isFieldJSON(tableName, fieldName string) bool {
	if con.schemas == nil {
		con.logger.Debug("no schemas provided for JSON detection")
		return false
	}

	schema, ok := con.schemas[tableName]
	if !ok {
		con.logger.Debug("schema not found for table", slog.String("table", tableName))
		return false
	}

	for _, field := range schema {
		if field.Name == fieldName {
			con.logger.LogAttrs(context.Background(), slog.LevelDebug,
				"field type lookup",
				slog.String("table", tableName),
				slog.String("field", fieldName),
				slog.Bool("is_json", field.IsJSON),
			)
			return field.IsJSON
		}
	}

	con.logger.Debug("field not found in schema",
		slog.String("table", tableName),
		slog.String("field", fieldName))
	return false
}

// getTableAndFieldFromSelectChain extracts the table name and field name from a select expression chain
// For obj.metadata, it returns ("obj", "metadata")
func (con *converter) getTableAndFieldFromSelectChain(expr *exprpb.Expr) (string, string, bool) {
	selectExpr := expr.GetSelectExpr()
	if selectExpr == nil {
		return "", "", false
	}

	fieldName := selectExpr.GetField()
	operand := selectExpr.GetOperand()

	// Check if the operand is an identifier (table name)
	if identExpr := operand.GetIdentExpr(); identExpr != nil {
		tableName := identExpr.GetName()
		return tableName, fieldName, true
	}

	return "", "", false
}

// isFieldJSONB checks if a field in a table is specifically JSONB (vs JSON) using schema information
func (con *converter) isFieldJSONB(tableName, fieldName string) bool {
	if con.schemas == nil {
		return false
	}

	schema, ok := con.schemas[tableName]
	if !ok {
		return false
	}

	for _, field := range schema {
		if field.Name == fieldName {
			return field.IsJSONB
		}
	}

	return false
}

// isFieldArray checks if a field in a table is an array using schema information
func (con *converter) isFieldArray(tableName, fieldName string) bool {
	if con.schemas == nil {
		return false
	}

	schema, ok := con.schemas[tableName]
	if !ok {
		return false
	}

	for _, field := range schema {
		if field.Name == fieldName {
			return field.Repeated
		}
	}

	return false
}

// getFieldElementType returns the element type of an array field using schema information
func (con *converter) getFieldElementType(tableName, fieldName string) string {
	if con.schemas == nil {
		return ""
	}

	schema, ok := con.schemas[tableName]
	if !ok {
		return ""
	}

	for _, field := range schema {
		if field.Name == fieldName && field.Repeated {
			return field.ElementType
		}
	}

	return ""
}

func (con *converter) visitCall(expr *exprpb.Expr) error {
	// Check for context cancellation before processing function calls
	if err := con.checkContext(); err != nil {
		return err
	}

	c := expr.GetCallExpr()
	fun := c.GetFunction()
	switch fun {
	// ternary operator
	case operators.Conditional:
		return con.visitCallConditional(expr)
	// index operator
	case operators.Index:
		return con.visitCallIndex(expr)
	// unary operators
	case operators.LogicalNot, operators.Negate:
		return con.visitCallUnary(expr)
	// binary operators
	case operators.Add,
		operators.Divide,
		operators.Equals,
		operators.Greater,
		operators.GreaterEquals,
		operators.In,
		operators.Less,
		operators.LessEquals,
		operators.LogicalAnd,
		operators.LogicalOr,
		operators.Multiply,
		operators.NotEquals,
		operators.OldIn,
		operators.Subtract:
		return con.visitCallBinary(expr)
	// standard function calls.
	default:
		return con.visitCallFunc(expr)
	}
}

func (con *converter) visitCallBinary(expr *exprpb.Expr) error {
	c := expr.GetCallExpr()
	fun := c.GetFunction()
	args := c.GetArgs()
	lhs := args[0]
	// add parens if the current operator is lower precedence than the lhs expr operator.
	lhsParen := isComplexOperatorWithRespectTo(fun, lhs)
	rhs := args[1]
	// add parens if the current operator is lower precedence than the rhs expr operator,
	// or the same precedence and the operator is left recursive.
	rhsParen := isComplexOperatorWithRespectTo(fun, rhs)
	lhsType := con.getType(lhs)
	rhsType := con.getType(rhs)
	if (isTimestampRelatedType(lhsType) && isDurationRelatedType(rhsType)) ||
		(isTimestampRelatedType(rhsType) && isDurationRelatedType(lhsType)) {
		return con.callTimestampOperation(fun, lhs, rhs)
	}
	if !rhsParen && isLeftRecursive(fun) {
		rhsParen = isSamePrecedence(fun, rhs)
	}

	// Check if we need numeric casting for JSON text extraction
	needsNumericCasting := false
	if con.isJSONTextExtraction(lhs) && isNumericComparison(fun) && isNumericType(rhsType) {
		needsNumericCasting = true
		con.str.WriteString("(")
	}

	if err := con.visitMaybeNested(lhs, lhsParen); err != nil {
		return err
	}

	if needsNumericCasting {
		con.str.WriteString(")::numeric")
	}
	var operator string
	if fun == operators.Add && (lhsType.GetPrimitive() == exprpb.Type_STRING && rhsType.GetPrimitive() == exprpb.Type_STRING) {
		operator = "||"
	} else if fun == operators.Add && (rhsType.GetPrimitive() == exprpb.Type_BYTES && lhsType.GetPrimitive() == exprpb.Type_BYTES) {
		operator = "||"
	} else if fun == operators.Add && (isListType(lhsType) && isListType(rhsType)) {
		operator = "||"
	} else if fun == operators.Add && (isStringLiteral(lhs) || isStringLiteral(rhs)) {
		// If either operand is a string literal, assume string concatenation
		operator = "||"
	} else if fun == operators.Equals && (isNullLiteral(rhs) || isBoolLiteral(rhs)) {
		operator = "IS"
	} else if fun == operators.NotEquals && (isNullLiteral(rhs) || isBoolLiteral(rhs)) {
		operator = "IS NOT"
	} else if fun == operators.In && isListType(rhsType) {
		operator = "="
	} else if fun == operators.In && isFieldAccessExpression(rhs) {
		// In PostgreSQL, field access expressions in IN clauses are likely array membership tests
		// For both JSON arrays and regular arrays, we use the same operator
		operator = "="
	} else if fun == operators.In {
		operator = "IN"
	} else if op, found := standardSQLBinaryOperators[fun]; found {
		operator = op
	} else if op, found := operators.FindReverseBinaryOperator(fun); found {
		operator = op
	} else {
		return newConversionErrorf(errMsgInvalidOperator, "binary operator: %s", fun)
	}

	con.logger.LogAttrs(context.Background(), slog.LevelDebug,
		"binary operator conversion",
		slog.String("cel_op", fun),
		slog.String("sql_op", operator),
	)

	con.str.WriteString(" ")
	con.str.WriteString(operator)
	con.str.WriteString(" ")
	if fun == operators.In && (isListType(rhsType) || isFieldAccessExpression(rhs)) {
		// Check if we're dealing with a JSON array
		if isFieldAccessExpression(rhs) && con.isJSONArrayField(rhs) {
			// For JSON arrays, use jsonb_array_elements with ANY
			jsonFunc := con.getJSONArrayFunction(rhs)
			con.str.WriteString("ANY(ARRAY(SELECT ")

			// For nested JSON access like settings.permissions, we need to handle differently
			if con.isNestedJSONAccess(rhs) {
				// Use text extraction for the array elements
				con.str.WriteString("jsonb_array_elements_text(")
				// Generate the JSON path with -> instead of ->> to preserve JSONB type
				if err := con.visitNestedJSONForArray(rhs); err != nil {
					return err
				}
				con.str.WriteString(")))")
				return nil
			}
			// For direct JSON array access
			con.str.WriteString(jsonFunc)
			con.str.WriteString("(")
			if err := con.visitMaybeNested(rhs, rhsParen); err != nil {
				return err
			}
			con.str.WriteString(")))")
			return nil
		}
		con.str.WriteString("ANY(")
	}
	if err := con.visitMaybeNested(rhs, rhsParen); err != nil {
		return err
	}
	if fun == operators.In && (isListType(rhsType) || isFieldAccessExpression(rhs)) {
		// Check if we're dealing with a JSON array - already handled above for JSON arrays
		if !isFieldAccessExpression(rhs) || !con.isJSONArrayField(rhs) {
			con.str.WriteString(")")
		}
	}
	return nil
}

func (con *converter) visitCallConditional(expr *exprpb.Expr) error {
	c := expr.GetCallExpr()
	args := c.GetArgs()
	con.str.WriteString("CASE WHEN ")
	if err := con.visit(args[0]); err != nil {
		return err
	}
	con.str.WriteString(" THEN ")
	if err := con.visit(args[1]); err != nil {
		return err
	}
	con.str.WriteString(" ELSE ")
	if err := con.visit(args[2]); err != nil {
		return err
	}
	con.str.WriteString(" END")
	return nil
}

func (con *converter) callContains(target *exprpb.Expr, args []*exprpb.Expr) error {
	// Check if the target is a JSON/JSONB field
	if target != nil && con.isJSONArrayField(target) {
		// For JSON/JSONB arrays, use the ? operator
		if err := con.visit(target); err != nil {
			return err
		}
		con.str.WriteString(" ? ")
		if len(args) > 0 {
			if err := con.visit(args[0]); err != nil {
				return err
			}
		}
		return nil
	}

	// For regular strings, use POSITION
	con.str.WriteString("POSITION(")
	for i, arg := range args {
		err := con.visit(arg)
		if err != nil {
			return err
		}
		if i < len(args)-1 {
			con.str.WriteString(" IN ")
		}
	}
	if target != nil {
		con.str.WriteString(" IN ")
		nested := isBinaryOrTernaryOperator(target)
		err := con.visitMaybeNested(target, nested)
		if err != nil {
			return err
		}
	}
	con.str.WriteString(") > 0")
	return nil
}

func (con *converter) callStartsWith(target *exprpb.Expr, args []*exprpb.Expr) error {
	// CEL startsWith function: string.startsWith(prefix)
	// Convert to PostgreSQL: string LIKE 'prefix%'
	// or for more robust handling: LEFT(string, LENGTH(prefix)) = prefix

	if target == nil || len(args) == 0 {
		return errors.New("startsWith function requires both string and prefix arguments")
	}

	// Visit the string expression
	nested := isBinaryOrTernaryOperator(target)
	if err := con.visitMaybeNested(target, nested); err != nil {
		return err
	}

	con.str.WriteString(" LIKE ")

	// Visit the prefix argument and append '%' for LIKE pattern
	// If it's a constant string, we can append % directly
	if constExpr := args[0].GetConstExpr(); constExpr != nil && constExpr.GetStringValue() != "" {
		prefix := constExpr.GetStringValue()
		// Reject patterns containing null bytes
		if strings.Contains(prefix, "\x00") {
			return errors.New("LIKE patterns cannot contain null bytes")
		}
		// Escape special LIKE characters: %, _, \
		escaped := escapeLikePattern(prefix)
		con.str.WriteString("'")
		con.str.WriteString(escaped)
		con.str.WriteString("%'")
	} else {
		// For non-literal patterns, concatenate with %
		if err := con.visit(args[0]); err != nil {
			return err
		}
		con.str.WriteString(" || '%'")
	}

	return nil
}

func (con *converter) callEndsWith(target *exprpb.Expr, args []*exprpb.Expr) error {
	// CEL endsWith function: string.endsWith(suffix)
	// Convert to PostgreSQL: string LIKE '%suffix'
	// or for more robust handling: RIGHT(string, LENGTH(suffix)) = suffix

	if target == nil || len(args) == 0 {
		return errors.New("endsWith function requires both string and suffix arguments")
	}

	// Visit the string expression
	nested := isBinaryOrTernaryOperator(target)
	if err := con.visitMaybeNested(target, nested); err != nil {
		return err
	}

	con.str.WriteString(" LIKE ")

	// Visit the suffix argument and prepend '%' for LIKE pattern
	// If it's a constant string, we can prepend % directly
	if constExpr := args[0].GetConstExpr(); constExpr != nil && constExpr.GetStringValue() != "" {
		suffix := constExpr.GetStringValue()
		// Reject patterns containing null bytes
		if strings.Contains(suffix, "\x00") {
			return errors.New("LIKE patterns cannot contain null bytes")
		}
		// Escape special LIKE characters: %, _, \
		escaped := escapeLikePattern(suffix)
		con.str.WriteString("'%")
		con.str.WriteString(escaped)
		con.str.WriteString("'")
	} else {
		// For non-literal patterns, concatenate with %
		con.str.WriteString("'%' || ")
		if err := con.visit(args[0]); err != nil {
			return err
		}
	}

	return nil
}

func (con *converter) callCasting(function string, _ *exprpb.Expr, args []*exprpb.Expr) error {
	arg := args[0]
	if function == overloads.TypeConvertInt && isTimestampType(con.getType(arg)) {
		con.str.WriteString("EXTRACT(EPOCH FROM ")
		if err := con.visit(arg); err != nil {
			return err
		}
		con.str.WriteString(")::bigint")
		return nil
	}
	con.str.WriteString("CAST(")
	if err := con.visit(arg); err != nil {
		return err
	}
	con.str.WriteString(" AS ")
	switch function {
	case overloads.TypeConvertBool:
		con.str.WriteString("BOOLEAN")
	case overloads.TypeConvertBytes:
		con.str.WriteString("BYTEA")
	case overloads.TypeConvertDouble:
		con.str.WriteString("DOUBLE PRECISION")
	case overloads.TypeConvertInt:
		con.str.WriteString("BIGINT")
	case overloads.TypeConvertString:
		con.str.WriteString("TEXT")
	case overloads.TypeConvertUint:
		con.str.WriteString("BIGINT")
	}
	con.str.WriteString(")")
	return nil
}

// callMatches handles CEL matches() function with RE2 to POSIX regex conversion
func (con *converter) callMatches(target *exprpb.Expr, args []*exprpb.Expr) error {
	// CEL matches function: string.matches(pattern) or matches(string, pattern)
	// Convert to PostgreSQL: string ~ 'posix_pattern'

	// Get the string to match against
	var stringExpr *exprpb.Expr
	var patternExpr *exprpb.Expr

	if target != nil {
		// Method call: string.matches(pattern)
		stringExpr = target
		if len(args) > 0 {
			patternExpr = args[0]
		}
	} else if len(args) >= 2 {
		// Function call: matches(string, pattern)
		stringExpr = args[0]
		patternExpr = args[1]
	}

	if stringExpr == nil || patternExpr == nil {
		return errors.New("matches function requires both string and pattern arguments")
	}

	// Visit the string expression
	if err := con.visit(stringExpr); err != nil {
		return err
	}

	con.str.WriteString(" ~ ")

	// Visit the pattern expression and convert from RE2 to POSIX if it's a string literal
	if constExpr := patternExpr.GetConstExpr(); constExpr != nil && constExpr.GetStringValue() != "" {
		// Convert RE2 pattern to POSIX
		re2Pattern := constExpr.GetStringValue()
		// Reject patterns containing null bytes
		if strings.Contains(re2Pattern, "\x00") {
			return errors.New("regex patterns cannot contain null bytes")
		}

		// Convert RE2 to POSIX with security validation
		posixPattern, err := convertRE2ToPOSIX(re2Pattern)
		if err != nil {
			return fmt.Errorf("invalid regex pattern: %w", err)
		}

		// Determine case sensitivity
		caseInsensitive := strings.HasPrefix(re2Pattern, "(?i)")

		con.logger.LogAttrs(context.Background(), slog.LevelDebug,
			"regex pattern conversion",
			slog.String("original_pattern", re2Pattern),
			slog.String("converted_pattern", posixPattern),
			slog.Bool("case_insensitive", caseInsensitive),
		)

		// Write the converted pattern as a string literal
		escaped := strings.ReplaceAll(posixPattern, "'", "''")
		con.str.WriteString("'")
		con.str.WriteString(escaped)
		con.str.WriteString("'")
	} else {
		// For non-literal patterns, we can't convert at compile time
		// Just use the pattern as-is and hope it's POSIX compatible
		if err := con.visit(patternExpr); err != nil {
			return err
		}
	}

	return nil
}

func (con *converter) visitCallFunc(expr *exprpb.Expr) error {
	c := expr.GetCallExpr()
	fun := c.GetFunction()
	target := c.GetTarget()
	args := c.GetArgs()
	switch fun {
	case overloads.Contains:
		return con.callContains(target, args)
	case overloads.StartsWith:
		return con.callStartsWith(target, args)
	case overloads.EndsWith:
		return con.callEndsWith(target, args)
	case overloads.Matches:
		return con.callMatches(target, args)
	case overloads.TypeConvertDuration:
		return con.callDuration(target, args)
	case "interval":
		return con.callInterval(target, args)
	case "timestamp":
		return con.callTimestampFromString(target, args)
	case overloads.TimeGetFullYear,
		overloads.TimeGetMonth,
		overloads.TimeGetDate,
		overloads.TimeGetHours,
		overloads.TimeGetMinutes,
		overloads.TimeGetSeconds,
		overloads.TimeGetMilliseconds,
		overloads.TimeGetDayOfYear,
		overloads.TimeGetDayOfMonth,
		overloads.TimeGetDayOfWeek:
		return con.callExtractFromTimestamp(fun, target, args)
	case overloads.TypeConvertBool,
		overloads.TypeConvertBytes,
		overloads.TypeConvertDouble,
		overloads.TypeConvertInt,
		overloads.TypeConvertString,
		overloads.TypeConvertUint:
		return con.callCasting(fun, target, args)
	}
	sqlFun, ok := standardSQLFunctions[fun]
	if !ok {
		if fun == overloads.Size {
			argType := con.getType(args[0])
			switch {
			case argType.GetPrimitive() == exprpb.Type_STRING:
				sqlFun = "LENGTH"
			case argType.GetPrimitive() == exprpb.Type_BYTES:
				sqlFun = "LENGTH"
			case isListType(argType):
				// Check if this is a JSON array field
				if len(args) > 0 && con.isJSONArrayField(args[0]) {
					// For JSON arrays, use jsonb_array_length
					con.str.WriteString("jsonb_array_length(")
					err := con.visit(args[0])
					if err != nil {
						return err
					}
					con.str.WriteString(")")
					return nil
				}
				// For PostgreSQL, we need to specify the array dimension (1 for 1D arrays)
				con.str.WriteString("ARRAY_LENGTH(")
				if target != nil {
					nested := isBinaryOrTernaryOperator(target)
					err := con.visitMaybeNested(target, nested)
					if err != nil {
						return err
					}
					con.str.WriteString(", ")
				}
				for i, arg := range args {
					err := con.visit(arg)
					if err != nil {
						return err
					}
					if i < len(args)-1 {
						con.str.WriteString(", ")
					}
				}
				con.str.WriteString(", 1)")
				return nil
			default:
				return newConversionErrorf(errMsgUnsupportedType, "size() argument type: %s", argType.String())
			}
		} else {
			sqlFun = strings.ToUpper(fun)
		}
	}
	con.str.WriteString(sqlFun)
	con.str.WriteString("(")
	if target != nil {
		nested := isBinaryOrTernaryOperator(target)
		err := con.visitMaybeNested(target, nested)
		if err != nil {
			return err
		}
		con.str.WriteString(", ")
	}
	for i, arg := range args {
		err := con.visit(arg)
		if err != nil {
			return err
		}
		if i < len(args)-1 {
			con.str.WriteString(", ")
		}
	}
	con.str.WriteString(")")
	return nil
}

func (con *converter) visitCallIndex(expr *exprpb.Expr) error {
	if isMapType(con.getType(expr.GetCallExpr().GetArgs()[0])) {
		return con.visitCallMapIndex(expr)
	}
	return con.visitCallListIndex(expr)
}

func (con *converter) visitCallMapIndex(expr *exprpb.Expr) error {
	c := expr.GetCallExpr()
	args := c.GetArgs()
	m := args[0]
	nested := isBinaryOrTernaryOperator(m)
	if err := con.visitMaybeNested(m, nested); err != nil {
		return err
	}
	fieldName, err := extractFieldName(args[1])
	if err != nil {
		return err
	}
	con.str.WriteString(".")
	con.str.WriteString(fieldName)
	return nil
}

func (con *converter) visitCallListIndex(expr *exprpb.Expr) error {
	c := expr.GetCallExpr()
	args := c.GetArgs()
	l := args[0]
	nested := isBinaryOrTernaryOperator(l)
	if err := con.visitMaybeNested(l, nested); err != nil {
		return err
	}
	con.str.WriteString("[")
	index := args[1]
	// PostgreSQL arrays are 1-indexed, CEL is 0-indexed, so add 1
	if constExpr := index.GetConstExpr(); constExpr != nil {
		idx := constExpr.GetInt64Value()
		if idx == math.MaxInt64 {
			return errors.New("array index overflow: cannot convert math.MaxInt64 to 1-based indexing")
		}
		if idx < 0 {
			return fmt.Errorf("invalid negative array index: %d", idx)
		}
		con.str.WriteString(strconv.FormatInt(idx+1, 10))
	} else {
		if err := con.visit(index); err != nil {
			return err
		}
		con.str.WriteString(" + 1")
	}
	con.str.WriteString("]")
	return nil
}

func (con *converter) visitCallUnary(expr *exprpb.Expr) error {
	c := expr.GetCallExpr()
	fun := c.GetFunction()
	args := c.GetArgs()
	var operator string
	if op, found := standardSQLUnaryOperators[fun]; found {
		operator = op
	} else if op, found := operators.FindReverse(fun); found {
		operator = op
	} else {
		return newConversionErrorf(errMsgInvalidOperator, "unary operator: %s", fun)
	}
	con.str.WriteString(operator)
	nested := isComplexOperator(args[0])
	return con.visitMaybeNested(args[0], nested)
}

func (con *converter) visitComprehension(expr *exprpb.Expr) error {
	// Track comprehension nesting depth to prevent resource exhaustion (CWE-400)
	con.comprehensionDepth++
	defer func() { con.comprehensionDepth-- }()

	// Check comprehension depth limit before context check (fail fast)
	if con.comprehensionDepth > maxComprehensionDepth {
		return fmt.Errorf("comprehension nesting depth %d exceeds maximum of %d",
			con.comprehensionDepth, maxComprehensionDepth)
	}

	// Check for context cancellation before processing comprehensions (potentially expensive)
	if err := con.checkContext(); err != nil {
		return err
	}

	info, err := con.identifyComprehension(expr)
	if err != nil {
		return fmt.Errorf("failed to identify comprehension: %w", err)
	}

	switch info.Type {
	case ComprehensionAll:
		return con.visitAllComprehension(expr, info)
	case ComprehensionExists:
		return con.visitExistsComprehension(expr, info)
	case ComprehensionExistsOne:
		return con.visitExistsOneComprehension(expr, info)
	case ComprehensionMap:
		return con.visitMapComprehension(expr, info)
	case ComprehensionFilter:
		return con.visitFilterComprehension(expr, info)
	case ComprehensionTransformList:
		return con.visitTransformListComprehension(expr, info)
	case ComprehensionTransformMap:
		return con.visitTransformMapComprehension(expr, info)
	case ComprehensionTransformMapEntry:
		return con.visitTransformMapEntryComprehension(expr, info)
	default:
		return newConversionErrorf(errMsgUnsupportedComprehension, "comprehension type: %s", info.Type.String())
	}
}

// Comprehension visit functions - Phase 1 placeholder implementations

func (con *converter) visitAllComprehension(expr *exprpb.Expr, info *ComprehensionInfo) error {
	// Generate SQL for ALL comprehension: all elements must satisfy the predicate
	// Pattern: NOT EXISTS (SELECT 1 FROM UNNEST(array) AS item WHERE NOT predicate)
	// For JSON arrays: NOT EXISTS (SELECT 1 FROM jsonb_array_elements(json_field) AS item WHERE NOT predicate)

	comprehension := expr.GetComprehensionExpr()
	if comprehension == nil {
		return newConversionError(errMsgUnsupportedComprehension, "expression is not a comprehension (ALL)")
	}

	iterRange := comprehension.GetIterRange()
	isJSONArray := con.isJSONArrayField(iterRange)

	con.str.WriteString("NOT EXISTS (SELECT 1 FROM ")

	if isJSONArray {
		jsonFunc := con.getJSONArrayFunction(iterRange)
		con.str.WriteString(jsonFunc)
		con.str.WriteString("(")
		if err := con.visit(iterRange); err != nil {
			return wrapConversionError(err, "visiting iter range in ALL comprehension")
		}
		con.str.WriteString(")")
	} else {
		con.str.WriteString("UNNEST(")
		if err := con.visit(iterRange); err != nil {
			return wrapConversionError(err, "visiting iter range in ALL comprehension")
		}
		con.str.WriteString(")")
	}

	con.str.WriteString(" AS ")
	con.str.WriteString(info.IterVar)

	con.str.WriteString(" WHERE ")

	// Add null checks for JSON arrays
	if isJSONArray {
		if err := con.visit(iterRange); err != nil {
			return wrapConversionError(err,"visiting iter range for null check")
		}
		con.str.WriteString(" IS NOT NULL AND ")
		typeofFunc := con.getJSONTypeofFunction(iterRange)
		con.str.WriteString(typeofFunc)
		con.str.WriteString("(")
		if err := con.visit(iterRange); err != nil {
			return wrapConversionError(err,"visiting iter range for type check")
		}
		con.str.WriteString(") = 'array'")

		if info.Predicate != nil {
			con.str.WriteString(" AND ")
		}
	}

	if info.Predicate != nil {
		con.str.WriteString("NOT (")
		if err := con.visit(info.Predicate); err != nil {
			return wrapConversionError(err,"visiting predicate in ALL comprehension")
		}
		con.str.WriteString(")")
	}

	con.str.WriteString(")")
	return nil
}

func (con *converter) visitExistsComprehension(expr *exprpb.Expr, info *ComprehensionInfo) error {
	// Generate SQL for EXISTS comprehension: at least one element satisfies the predicate
	// Pattern: EXISTS (SELECT 1 FROM UNNEST(array) AS item WHERE predicate)
	// For JSON arrays: EXISTS (SELECT 1 FROM jsonb_array_elements(json_field) AS item WHERE predicate)

	comprehension := expr.GetComprehensionExpr()
	if comprehension == nil {
		return newConversionError(errMsgUnsupportedComprehension, "expression is not a comprehension (EXISTS)")
	}

	iterRange := comprehension.GetIterRange()
	isJSONArray := con.isJSONArrayField(iterRange)

	con.str.WriteString("EXISTS (SELECT 1 FROM ")

	if isJSONArray {
		jsonFunc := con.getJSONArrayFunction(iterRange)
		con.str.WriteString(jsonFunc)
		con.str.WriteString("(")
		if err := con.visit(iterRange); err != nil {
			return wrapConversionError(err,"visiting iter range in EXISTS comprehension")
		}
		con.str.WriteString(")")
	} else {
		con.str.WriteString("UNNEST(")
		if err := con.visit(iterRange); err != nil {
			return wrapConversionError(err,"visiting iter range in EXISTS comprehension")
		}
		con.str.WriteString(")")
	}

	con.str.WriteString(" AS ")
	con.str.WriteString(info.IterVar)

	con.str.WriteString(" WHERE ")

	// Add null checks for JSON arrays
	if isJSONArray {
		if err := con.visit(iterRange); err != nil {
			return wrapConversionError(err,"visiting iter range for null check")
		}
		con.str.WriteString(" IS NOT NULL AND ")
		typeofFunc := con.getJSONTypeofFunction(iterRange)
		con.str.WriteString(typeofFunc)
		con.str.WriteString("(")
		if err := con.visit(iterRange); err != nil {
			return wrapConversionError(err,"visiting iter range for type check")
		}
		con.str.WriteString(") = 'array'")

		if info.Predicate != nil {
			con.str.WriteString(" AND ")
		}
	}

	if info.Predicate != nil {
		if err := con.visit(info.Predicate); err != nil {
			return wrapConversionError(err,"visiting predicate in EXISTS comprehension")
		}
	}

	con.str.WriteString(")")
	return nil
}

func (con *converter) visitExistsOneComprehension(expr *exprpb.Expr, info *ComprehensionInfo) error {
	// Generate SQL for EXISTS_ONE comprehension: exactly one element satisfies the predicate
	// Pattern: (SELECT COUNT(*) FROM UNNEST(array) AS item WHERE predicate) = 1
	// For JSON arrays: (SELECT COUNT(*) FROM jsonb_array_elements(json_field) AS item WHERE predicate) = 1

	comprehension := expr.GetComprehensionExpr()
	if comprehension == nil {
		return newConversionError(errMsgUnsupportedComprehension, "expression is not a comprehension (EXISTS_ONE)")
	}

	iterRange := comprehension.GetIterRange()
	isJSONArray := con.isJSONArrayField(iterRange)

	con.str.WriteString("(SELECT COUNT(*) FROM ")

	if isJSONArray {
		jsonFunc := con.getJSONArrayFunction(iterRange)
		con.str.WriteString(jsonFunc)
		con.str.WriteString("(")
		if err := con.visit(iterRange); err != nil {
			return wrapConversionError(err,"visiting iter range in EXISTS_ONE comprehension")
		}
		con.str.WriteString(")")
	} else {
		con.str.WriteString("UNNEST(")
		if err := con.visit(iterRange); err != nil {
			return wrapConversionError(err,"visiting iter range in EXISTS_ONE comprehension")
		}
		con.str.WriteString(")")
	}

	con.str.WriteString(" AS ")
	con.str.WriteString(info.IterVar)

	con.str.WriteString(" WHERE ")

	// Add null checks for JSON arrays
	if isJSONArray {
		if err := con.visit(iterRange); err != nil {
			return wrapConversionError(err,"visiting iter range for null check")
		}
		con.str.WriteString(" IS NOT NULL AND ")
		typeofFunc := con.getJSONTypeofFunction(iterRange)
		con.str.WriteString(typeofFunc)
		con.str.WriteString("(")
		if err := con.visit(iterRange); err != nil {
			return wrapConversionError(err,"visiting iter range for type check")
		}
		con.str.WriteString(") = 'array'")

		if info.Predicate != nil {
			con.str.WriteString(" AND ")
		}
	}

	if info.Predicate != nil {
		if err := con.visit(info.Predicate); err != nil {
			return wrapConversionError(err,"visiting predicate in EXISTS_ONE comprehension")
		}
	}

	con.str.WriteString(") = 1")
	return nil
}

func (con *converter) visitMapComprehension(expr *exprpb.Expr, info *ComprehensionInfo) error {
	// Generate SQL for MAP comprehension: transform elements using the transform expression
	// Pattern: ARRAY(SELECT transform FROM UNNEST(array) AS item [WHERE filter])
	// For JSON arrays: ARRAY(SELECT transform FROM jsonb_array_elements(json_field) AS item [WHERE filter])

	comprehension := expr.GetComprehensionExpr()
	if comprehension == nil {
		return newConversionError(errMsgUnsupportedComprehension, "expression is not a comprehension (MAP)")
	}

	iterRange := comprehension.GetIterRange()
	isJSONArray := con.isJSONArrayField(iterRange)

	con.str.WriteString("ARRAY(SELECT ")

	// Visit the transform expression
	if info.Transform != nil {
		if err := con.visit(info.Transform); err != nil {
			return wrapConversionError(err,"visiting transform in MAP comprehension")
		}
	} else {
		// If no transform, just return the variable itself
		con.str.WriteString(info.IterVar)
	}

	con.str.WriteString(" FROM ")

	if isJSONArray {
		jsonFunc := con.getJSONArrayFunction(iterRange)
		con.str.WriteString(jsonFunc)
		con.str.WriteString("(")
		if err := con.visit(iterRange); err != nil {
			return wrapConversionError(err,"visiting iter range in MAP comprehension")
		}
		con.str.WriteString(")")
	} else {
		con.str.WriteString("UNNEST(")
		if err := con.visit(iterRange); err != nil {
			return wrapConversionError(err,"visiting iter range in MAP comprehension")
		}
		con.str.WriteString(")")
	}

	con.str.WriteString(" AS ")
	con.str.WriteString(info.IterVar)

	// Add filter condition if present (for map with filter)
	if info.Filter != nil {
		con.str.WriteString(" WHERE ")
		if err := con.visit(info.Filter); err != nil {
			return wrapConversionError(err,"visiting filter in MAP comprehension")
		}
	}

	con.str.WriteString(")")
	return nil
}

func (con *converter) visitFilterComprehension(expr *exprpb.Expr, info *ComprehensionInfo) error {
	// Generate SQL for FILTER comprehension: return elements that satisfy the predicate
	// Pattern: ARRAY(SELECT item FROM UNNEST(array) AS item WHERE predicate)
	// For JSON arrays: ARRAY(SELECT item FROM jsonb_array_elements(json_field) AS item WHERE predicate)

	comprehension := expr.GetComprehensionExpr()
	if comprehension == nil {
		return newConversionError(errMsgUnsupportedComprehension, "expression is not a comprehension (FILTER)")
	}

	iterRange := comprehension.GetIterRange()
	isJSONArray := con.isJSONArrayField(iterRange)

	con.str.WriteString("ARRAY(SELECT ")
	con.str.WriteString(info.IterVar)
	con.str.WriteString(" FROM ")

	if isJSONArray {
		jsonFunc := con.getJSONArrayFunction(iterRange)
		con.str.WriteString(jsonFunc)
		con.str.WriteString("(")
		if err := con.visit(iterRange); err != nil {
			return wrapConversionError(err,"visiting iter range in FILTER comprehension")
		}
		con.str.WriteString(")")
	} else {
		con.str.WriteString("UNNEST(")
		if err := con.visit(iterRange); err != nil {
			return wrapConversionError(err,"visiting iter range in FILTER comprehension")
		}
		con.str.WriteString(")")
	}

	con.str.WriteString(" AS ")
	con.str.WriteString(info.IterVar)

	if info.Predicate != nil {
		con.str.WriteString(" WHERE ")
		if err := con.visit(info.Predicate); err != nil {
			return wrapConversionError(err,"visiting predicate in FILTER comprehension")
		}
	}

	con.str.WriteString(")")
	return nil
}

func (con *converter) visitTransformListComprehension(expr *exprpb.Expr, info *ComprehensionInfo) error {
	// Generate SQL for TRANSFORM_LIST comprehension: similar to MAP but may have different semantics
	// Pattern: ARRAY(SELECT transform FROM UNNEST(array) AS item [WHERE filter])

	comprehension := expr.GetComprehensionExpr()
	if comprehension == nil {
		return newConversionError(errMsgUnsupportedComprehension, "expression is not a comprehension (TRANSFORM_LIST)")
	}

	con.str.WriteString("ARRAY(SELECT ")

	// Visit the transform expression
	if info.Transform != nil {
		if err := con.visit(info.Transform); err != nil {
			return wrapConversionError(err,"visiting transform in TRANSFORM_LIST comprehension")
		}
	} else {
		// If no transform, just return the variable itself
		con.str.WriteString(info.IterVar)
	}

	con.str.WriteString(" FROM UNNEST(")

	// Visit the iterable range (the array/list being comprehended over)
	if err := con.visit(comprehension.GetIterRange()); err != nil {
		return wrapConversionError(err,"visiting iter range in TRANSFORM_LIST comprehension")
	}

	con.str.WriteString(") AS ")
	con.str.WriteString(info.IterVar)

	// Add filter condition if present
	if info.Filter != nil {
		con.str.WriteString(" WHERE ")
		if err := con.visit(info.Filter); err != nil {
			return wrapConversionError(err,"visiting filter in TRANSFORM_LIST comprehension")
		}
	}

	con.str.WriteString(")")
	return nil
}

func (con *converter) visitTransformMapComprehension(_ *exprpb.Expr, _ *ComprehensionInfo) error {
	// Generate SQL for TRANSFORM_MAP comprehension: work with map entries
	// This is complex for PostgreSQL - maps are typically represented as JSON or composite types
	// For now, return an error indicating this needs special handling
	return errors.New("TRANSFORM_MAP comprehension requires map/JSON support: not yet implemented")
}

func (con *converter) visitTransformMapEntryComprehension(_ *exprpb.Expr, _ *ComprehensionInfo) error {
	// Generate SQL for TRANSFORM_MAP_ENTRY comprehension: work with map key-value pairs
	// This is complex for PostgreSQL - maps are typically represented as JSON or composite types
	// For now, return an error indicating this needs special handling
	return errors.New("TRANSFORM_MAP_ENTRY comprehension requires map/JSON support: not yet implemented")
}

func (con *converter) visitConst(expr *exprpb.Expr) error {
	c := expr.GetConstExpr()
	switch c.ConstantKind.(type) {
	case *exprpb.Constant_BoolValue:
		// Always inline TRUE/FALSE for PostgreSQL query plan efficiency
		if c.GetBoolValue() {
			con.str.WriteString("TRUE")
		} else {
			con.str.WriteString("FALSE")
		}
	case *exprpb.Constant_NullValue:
		// Always inline NULL for PostgreSQL query plan efficiency
		con.str.WriteString("NULL")
	case *exprpb.Constant_Int64Value:
		if con.parameterize {
			con.paramCount++
			con.str.WriteString(fmt.Sprintf("$%d", con.paramCount))
			con.parameters = append(con.parameters, c.GetInt64Value())
		} else {
			i := strconv.FormatInt(c.GetInt64Value(), 10)
			con.str.WriteString(i)
		}
	case *exprpb.Constant_Uint64Value:
		if con.parameterize {
			con.paramCount++
			con.str.WriteString(fmt.Sprintf("$%d", con.paramCount))
			con.parameters = append(con.parameters, c.GetUint64Value())
		} else {
			ui := strconv.FormatUint(c.GetUint64Value(), 10)
			con.str.WriteString(ui)
		}
	case *exprpb.Constant_DoubleValue:
		if con.parameterize {
			con.paramCount++
			con.str.WriteString(fmt.Sprintf("$%d", con.paramCount))
			con.parameters = append(con.parameters, c.GetDoubleValue())
		} else {
			d := strconv.FormatFloat(c.GetDoubleValue(), 'g', -1, 64)
			con.str.WriteString(d)
		}
	case *exprpb.Constant_StringValue:
		str := c.GetStringValue()
		// Reject strings containing null bytes
		if strings.Contains(str, "\x00") {
			return errors.New("string literals cannot contain null bytes")
		}

		if con.parameterize {
			con.paramCount++
			con.str.WriteString(fmt.Sprintf("$%d", con.paramCount))
			con.parameters = append(con.parameters, str)
		} else {
			// Use single quotes for PostgreSQL string literals
			// Escape single quotes by doubling them
			escaped := strings.ReplaceAll(str, "'", "''")
			con.str.WriteString("'")
			con.str.WriteString(escaped)
			con.str.WriteString("'")
		}
	case *exprpb.Constant_BytesValue:
		b := c.GetBytesValue()

		if con.parameterize {
			con.paramCount++
			con.str.WriteString(fmt.Sprintf("$%d", con.paramCount))
			con.parameters = append(con.parameters, b)
		} else {
			con.str.WriteString("'\\x")
			con.str.WriteString(hex.EncodeToString(b))
			con.str.WriteString("'")
		}
	default:
		return newConversionErrorf(errMsgUnsupportedExpression, "constant type: %T", c.ConstantKind)
	}
	return nil
}

func (con *converter) visitIdent(expr *exprpb.Expr) error {
	identName := expr.GetIdentExpr().GetName()

	// Validate identifier name for security (prevent SQL injection)
	if err := validateFieldName(identName); err != nil {
		return fmt.Errorf("invalid identifier name: %w", err)
	}

	// Check if this identifier needs numeric casting for JSON comprehensions
	if con.needsNumericCasting(identName) {
		con.str.WriteString("(")
		con.str.WriteString(identName)
		con.str.WriteString(")::numeric")
	} else {
		con.str.WriteString(identName)
	}
	return nil
}

func (con *converter) visitList(expr *exprpb.Expr) error {
	l := expr.GetListExpr()
	elems := l.GetElements()
	con.str.WriteString("ARRAY[")
	for i, elem := range elems {
		err := con.visit(elem)
		if err != nil {
			return err
		}
		if i < len(elems)-1 {
			con.str.WriteString(", ")
		}
	}
	con.str.WriteString("]")
	return nil
}

func (con *converter) visitSelect(expr *exprpb.Expr) error {
	sel := expr.GetSelectExpr()

	// Validate field name for security (prevent SQL injection)
	fieldName := sel.GetField()
	if err := validateFieldName(fieldName); err != nil {
		return fmt.Errorf("invalid field name in select expression: %w", err)
	}

	// Handle the case when the select expression was generated by the has() macro.
	if sel.GetTestOnly() {
		return con.visitHasFunction(expr)
	}

	// Check if we should use JSON path operators
	// We need to determine if the operand is a JSON/JSONB field
	useJSONPath := con.shouldUseJSONPath(sel.GetOperand(), fieldName)
	useJSONObjectAccess := con.isJSONObjectFieldAccess(expr)

	// Check if this is a nested JSON path that requires special handling
	if useJSONPath && !useJSONObjectAccess {
		// Use the specialized JSON path builder for nested access
		return con.buildJSONPath(expr)
	}

	nested := !sel.GetTestOnly() && isBinaryOrTernaryOperator(sel.GetOperand())

	if useJSONObjectAccess && con.isNumericJSONField(fieldName) {
		// For numeric JSON fields, wrap in parentheses for casting
		con.str.WriteString("(")
	}

	err := con.visitMaybeNested(sel.GetOperand(), nested)
	if err != nil {
		return err
	}

	switch {
	case useJSONPath:
		// Use ->> for text extraction
		con.str.WriteString("->>")
		con.str.WriteString("'")
		con.str.WriteString(escapeJSONFieldName(fieldName))
		con.str.WriteString("'")
	case useJSONObjectAccess:
		// Use -> for JSON object field access in comprehensions
		con.str.WriteString("->>'")
		con.str.WriteString(escapeJSONFieldName(fieldName))
		con.str.WriteString("'")
		if con.isNumericJSONField(fieldName) {
			// Close parentheses and add numeric cast
			con.str.WriteString(")::numeric")
		}
	default:
		// Regular field selection
		con.str.WriteString(".")
		con.str.WriteString(fieldName)
	}

	return nil
}

// visitHasFunction handles the has() macro for field existence checks
func (con *converter) visitHasFunction(expr *exprpb.Expr) error {
	sel := expr.GetSelectExpr()
	operand := sel.GetOperand()
	field := sel.GetField()

	// Check if this is a direct JSON field access (e.g., table.json_column.key)
	if con.isDirectJSONFieldAccess(operand, field) {
		// For direct JSON field access, use the appropriate existence operator
		err := con.visitMaybeNested(operand, isBinaryOrTernaryOperator(operand))
		if err != nil {
			return err
		}

		// Check if this is a JSONB field
		if con.isJSONBField(operand) {
			// Use JSONB's ? operator for existence check
			con.str.WriteString(" ? '")
			con.str.WriteString(escapeJSONFieldName(field))
			con.str.WriteString("'")
		} else {
			// For JSON fields, check if the field is not null
			con.str.WriteString("->'")
			con.str.WriteString(escapeJSONFieldName(field))
			con.str.WriteString("' IS NOT NULL")
		}
		return nil
	}

	// Check if this is a nested JSON path (e.g., table.json_column.key.subkey)
	if con.hasJSONFieldInChain(expr) {
		return con.visitNestedJSONHas(expr)
	}

	// For regular struct fields, check if the field is not null
	err := con.visitMaybeNested(operand, isBinaryOrTernaryOperator(operand))
	if err != nil {
		return err
	}
	con.str.WriteString(".")
	con.str.WriteString(field)
	con.str.WriteString(" IS NOT NULL")

	return nil
}

// isDirectJSONFieldAccess checks if this represents a direct JSON field access (table.json_column.key)
func (con *converter) isDirectJSONFieldAccess(operand *exprpb.Expr, _ string) bool {
	// Check if operand is a select expression that refers to a JSON column
	if selectExpr := operand.GetSelectExpr(); selectExpr != nil {
		parentField := selectExpr.GetField()

		// Check if the parent field is a known JSON column
		jsonFields := []string{"metadata", "properties", "content", "structure", "taxonomy", "analytics", "classification"}
		for _, jsonField := range jsonFields {
			if parentField == jsonField {
				return true
			}
		}
	}

	return false
}

// visitNestedJSONHas handles has() for deeply nested JSON paths
func (con *converter) visitNestedJSONHas(expr *exprpb.Expr) error {
	// For nested JSON paths, we use jsonb_extract_path_text and check for NOT NULL
	// This is more reliable than trying to use ? operator on nested paths
	con.str.WriteString("jsonb_extract_path_text(")

	// Get the root JSON column and remaining path segments
	rootColumn, pathSegments := con.getJSONRootAndPath(expr)

	// Visit the root column without adding JSON access operators
	if err := con.visitJSONColumnReference(rootColumn); err != nil {
		return err
	}

	// Add path segments as arguments
	for _, segment := range pathSegments {
		con.str.WriteString(", '")
		con.str.WriteString(escapeJSONFieldName(segment))
		con.str.WriteString("'")
	}

	con.str.WriteString(") IS NOT NULL")
	return nil
}

// visitJSONColumnReference visits a JSON column reference without adding JSON access operators
// This is used for jsonb_extract_path_text where we need the column reference as-is
func (con *converter) visitJSONColumnReference(expr *exprpb.Expr) error {
	if selectExpr := expr.GetSelectExpr(); selectExpr != nil {
		operand := selectExpr.GetOperand()
		field := selectExpr.GetField()

		// Visit the operand (table name)
		if err := con.visit(operand); err != nil {
			return err
		}

		// Add the field name with a simple dot notation
		con.str.WriteString(".")
		con.str.WriteString(field)
		return nil
	}

	// If it's not a SelectExpr, just visit it normally
	return con.visit(expr)
}

// getJSONRootAndPath extracts the root JSON column and path segments for a JSON expression
func (con *converter) getJSONRootAndPath(expr *exprpb.Expr) (*exprpb.Expr, []string) {
	var pathSegments []string
	current := expr

	// Walk up the chain to collect path segments until we reach a JSON column
	for {
		if sel := current.GetSelectExpr(); sel != nil {
			fieldName := sel.GetField()
			operand := sel.GetOperand()

			// Check if this field is a JSON column
			if con.isJSONColumn(operand, fieldName) {
				// We've found the JSON column boundary
				// Create the JSON column expression (table.json_column)
				jsonColumnExpr := &exprpb.Expr{
					ExprKind: &exprpb.Expr_SelectExpr{
						SelectExpr: &exprpb.Expr_Select{
							Operand: operand,
							Field:   fieldName,
						},
					},
				}
				return jsonColumnExpr, pathSegments
			}
			// This field is part of the path, continue up the chain
			pathSegments = append([]string{fieldName}, pathSegments...)
			current = operand
		} else {
			break
		}
	}

	// If we didn't find a clear JSON column boundary, return what we have
	return current, pathSegments
}

// isJSONColumn checks if the operand refers to a JSON column
func (con *converter) isJSONColumn(operand *exprpb.Expr, field string) bool {
	// Check if the field name is a known JSON column
	jsonColumns := []string{"metadata", "properties", "content", "structure", "taxonomy", "analytics", "classification"}
	for _, jsonCol := range jsonColumns {
		if field == jsonCol {
			// Additional check: make sure the operand is a table reference, not another JSON field
			if con.isTableReference(operand) {
				return true
			}
		}
	}
	return false
}

// isTableReference checks if an expression refers to a table (not a JSON field)
func (con *converter) isTableReference(expr *exprpb.Expr) bool {
	if identExpr := expr.GetIdentExpr(); identExpr != nil {
		// Direct table reference (e.g., "information_assets")
		return true
	}

	// For now, assume SelectExpr that doesn't have JSON field characteristics is also a table reference
	// This is a simplification but should work for our use cases
	return false
}

func (con *converter) visitStruct(expr *exprpb.Expr) error {
	s := expr.GetStructExpr()
	// If the message name is non-empty, then this should be treated as message construction.
	if s.GetMessageName() != "" {
		return con.visitStructMsg(expr)
	}
	// Otherwise, build a map.
	return con.visitStructMap(expr)
}

func (con *converter) visitStructMsg(expr *exprpb.Expr) error {
	m := expr.GetStructExpr()
	entries := m.GetEntries()
	con.str.WriteString(m.GetMessageName())
	con.str.WriteString("{")
	for i, entry := range entries {
		f := entry.GetFieldKey()
		con.str.WriteString(f)
		con.str.WriteString(": ")
		v := entry.GetValue()
		err := con.visit(v)
		if err != nil {
			return err
		}
		if i < len(entries)-1 {
			con.str.WriteString(", ")
		}
	}
	con.str.WriteString("}")
	return nil
}

func (con *converter) visitStructMap(expr *exprpb.Expr) error {
	m := expr.GetStructExpr()
	entries := m.GetEntries()
	con.str.WriteString("ROW(")
	for i, entry := range entries {
		v := entry.GetValue()
		if err := con.visit(v); err != nil {
			return err
		}
		if i < len(entries)-1 {
			con.str.WriteString(", ")
		}
	}
	con.str.WriteString(")")
	return nil
}

func (con *converter) visitMaybeNested(expr *exprpb.Expr, nested bool) error {
	if nested {
		con.str.WriteString("(")
	}
	err := con.visit(expr)
	if err != nil {
		return err
	}
	if nested {
		con.str.WriteString(")")
	}
	return nil
}

func (con *converter) getType(node *exprpb.Expr) *exprpb.Type {
	return con.typeMap[node.GetId()]
}

// isLeftRecursive indicates whether the parser resolves the call in a left-recursive manner as
// this can have an effect of how parentheses affect the order of operations in the AST.
func isLeftRecursive(op string) bool {
	return op != operators.LogicalAnd && op != operators.LogicalOr
}

// isSamePrecedence indicates whether the precedence of the input operator is the same as the
// precedence of the (possible) operation represented in the input Expr.
//
// If the expr is not a Call, the result is false.
func isSamePrecedence(op string, expr *exprpb.Expr) bool {
	if expr.GetCallExpr() == nil {
		return false
	}
	c := expr.GetCallExpr()
	other := c.GetFunction()
	return operators.Precedence(op) == operators.Precedence(other)
}

// isLowerPrecedence indicates whether the precedence of the input operator is lower precedence
// than the (possible) operation represented in the input Expr.
//
// If the expr is not a Call, the result is false.
func isLowerPrecedence(op string, expr *exprpb.Expr) bool {
	if expr.GetCallExpr() == nil {
		return false
	}
	c := expr.GetCallExpr()
	other := c.GetFunction()
	return operators.Precedence(op) < operators.Precedence(other)
}

// Indicates whether the expr is a complex operator, i.e., a call expression
// with 2 or more arguments.
func isComplexOperator(expr *exprpb.Expr) bool {
	if expr.GetCallExpr() != nil && len(expr.GetCallExpr().GetArgs()) >= 2 {
		return true
	}
	return false
}

// Indicates whether it is a complex operation compared to another.
// expr is *not* considered complex if it is not a call expression or has
// less than two arguments, or if it has a higher precedence than op.
func isComplexOperatorWithRespectTo(op string, expr *exprpb.Expr) bool {
	if expr.GetCallExpr() == nil || len(expr.GetCallExpr().GetArgs()) < 2 {
		return false
	}
	return isLowerPrecedence(op, expr)
}

// Indicate whether this is a binary or ternary operator.
func isBinaryOrTernaryOperator(expr *exprpb.Expr) bool {
	if expr.GetCallExpr() == nil || len(expr.GetCallExpr().GetArgs()) < 2 {
		return false
	}
	_, isBinaryOp := operators.FindReverseBinaryOperator(expr.GetCallExpr().GetFunction())
	return isBinaryOp || isSamePrecedence(operators.Conditional, expr)
}

// convertRE2ToPOSIX converts an RE2 regex pattern to POSIX ERE format for PostgreSQL.
// It performs security validation to prevent ReDoS attacks (CWE-1333).
// Note: This is a basic conversion for common patterns. Full RE2 to POSIX conversion is complex.
func convertRE2ToPOSIX(re2Pattern string) (string, error) {
	// 1. Check pattern length to prevent processing extremely long patterns
	if len(re2Pattern) > maxRegexPatternLength {
		return "", fmt.Errorf("regex pattern exceeds maximum length of %d characters", maxRegexPatternLength)
	}

	// 2. Detect catastrophic nested quantifiers that cause exponential backtracking
	// Patterns like (a+)+, (a*)*,  (x+x+)+, ((a)+b)+, etc. are extremely dangerous

	// Check for doubled quantifiers
	if matched, _ := regexp.MatchString(`[*+][*+]`, re2Pattern); matched {
		return "", errors.New("regex contains catastrophic nested quantifiers that could cause ReDoS")
	}

	// Check for groups that contain quantifiers and are themselves quantified
	// This catches patterns like (a+)+, ((a)+b)+, (a*b*)*, etc.
	// We need to check if any opening paren eventually leads to a closing paren followed by a quantifier,
	// and if there are quantifiers between those parens.
	depth := 0
	groupHasQuantifier := make([]bool, 0)

	for i := 0; i < len(re2Pattern); i++ {
		char := re2Pattern[i]

		// Skip escaped characters
		if i > 0 && re2Pattern[i-1] == '\\' {
			continue
		}

		switch char {
		case '(':
			depth++
			groupHasQuantifier = append(groupHasQuantifier, false)
		case ')':
			if depth > 0 {
				depth--
				// Check if the closing paren is followed by a quantifier
				if i+1 < len(re2Pattern) {
					nextChar := re2Pattern[i+1]
					if nextChar == '*' || nextChar == '+' || nextChar == '?' || nextChar == '{' {
						// This group is quantified. Check if it contains quantifiers
						if len(groupHasQuantifier) > 0 && groupHasQuantifier[len(groupHasQuantifier)-1] {
							return "", errors.New("regex contains catastrophic nested quantifiers that could cause ReDoS")
						}
					}
				}
				if len(groupHasQuantifier) > 0 {
					// Pop the last group
					if len(groupHasQuantifier) > 1 {
						// If inner group had quantifier, mark outer group as having quantifier too
						if groupHasQuantifier[len(groupHasQuantifier)-1] {
							groupHasQuantifier[len(groupHasQuantifier)-2] = true
						}
					}
					groupHasQuantifier = groupHasQuantifier[:len(groupHasQuantifier)-1]
				}
			}
		case '*', '+', '?':
			// Mark that current group contains a quantifier
			if len(groupHasQuantifier) > 0 {
				groupHasQuantifier[len(groupHasQuantifier)-1] = true
			}
		case '{':
			// Brace quantifier {n,m}
			if len(groupHasQuantifier) > 0 {
				groupHasQuantifier[len(groupHasQuantifier)-1] = true
			}
		}
	}

	// 3. Count and limit capture groups to prevent memory exhaustion
	groupCount := strings.Count(re2Pattern, "(") - strings.Count(re2Pattern, `\(`)
	if groupCount > maxRegexGroups {
		return "", fmt.Errorf("regex contains %d capture groups, exceeds maximum of %d", groupCount, maxRegexGroups)
	}

	// 4. Detect exponential alternation patterns like (a|a)*b or (a|ab)*
	alternationPattern := regexp.MustCompile(`\([^)]*\|[^)]*\)[*+]`)
	if alternationPattern.MatchString(re2Pattern) {
		// Check if alternation has overlapping branches (more dangerous)
		// This is a simple heuristic - full analysis would be more complex
		return "", errors.New("regex contains quantified alternation that could cause ReDoS")
	}

	// 5. Check nesting depth to prevent deeply nested patterns
	maxDepth := 0
	currentDepth := 0
	for _, char := range re2Pattern {
		if char == '(' && !strings.HasSuffix(re2Pattern[:strings.LastIndex(re2Pattern, string(char))], `\`) {
			currentDepth++
			if currentDepth > maxDepth {
				maxDepth = currentDepth
			}
		} else if char == ')' && !strings.HasSuffix(re2Pattern[:strings.LastIndex(re2Pattern, string(char))], `\`) {
			currentDepth--
		}
	}
	if maxDepth > maxRegexNestingDepth {
		return "", fmt.Errorf("regex nesting depth %d exceeds maximum of %d", maxDepth, maxRegexNestingDepth)
	}

	// Passed all security checks - proceed with conversion
	posixPattern := re2Pattern

	// Basic conversions for common differences between RE2 and POSIX:

	// 1. Word boundaries: \b -> [[:<:]] and [[:<:]] (PostgreSQL extension)
	//    Note: PostgreSQL supports \y for word boundaries in some contexts
	posixPattern = strings.ReplaceAll(posixPattern, `\b`, `\y`)

	// 2. Non-word boundaries: \B -> [^[:alnum:]_] (approximate)
	//    This is a simplification; exact conversion is complex
	posixPattern = strings.ReplaceAll(posixPattern, `\B`, `[^[:alnum:]_]`)

	// 3. Digit shortcuts: \d -> [[:digit:]] or [0-9]
	posixPattern = strings.ReplaceAll(posixPattern, `\d`, `[[:digit:]]`)

	// 4. Non-digit shortcuts: \D -> [^[:digit:]] or [^0-9]
	posixPattern = strings.ReplaceAll(posixPattern, `\D`, `[^[:digit:]]`)

	// 5. Word character shortcuts: \w -> [[:alnum:]_]
	posixPattern = strings.ReplaceAll(posixPattern, `\w`, `[[:alnum:]_]`)

	// 6. Non-word character shortcuts: \W -> [^[:alnum:]_]
	posixPattern = strings.ReplaceAll(posixPattern, `\W`, `[^[:alnum:]_]`)

	// 7. Whitespace shortcuts: \s -> [[:space:]]
	posixPattern = strings.ReplaceAll(posixPattern, `\s`, `[[:space:]]`)

	// 8. Non-whitespace shortcuts: \S -> [^[:space:]]
	posixPattern = strings.ReplaceAll(posixPattern, `\S`, `[^[:space:]]`)

	// Note: Many RE2 features are not directly convertible to POSIX ERE:
	// - Lookahead/lookbehind assertions (?=...), (?!...), (?<=...), (?<!...)
	// - Non-capturing groups (?:...)
	// - Named groups (?P<name>...)
	// - Case-insensitive flags (?i)
	// - Multiline flags (?m)
	// - Unicode character classes
	//
	// For these cases, the pattern is returned as-is, which may cause PostgreSQL errors
	// if the pattern uses unsupported RE2 features.

	return posixPattern, nil
}
