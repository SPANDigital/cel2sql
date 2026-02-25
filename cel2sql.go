// Package cel2sql converts CEL (Common Expression Language) expressions to SQL conditions.
// It supports multiple SQL dialects through the dialect interface, with PostgreSQL as the default.
package cel2sql

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/common/overloads"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"

	"github.com/spandigital/cel2sql/v3/dialect"
	"github.com/spandigital/cel2sql/v3/dialect/postgres"
	"github.com/spandigital/cel2sql/v3/schema"
)

// Implementations based on `google/cel-go`'s unparser
// https://github.com/google/cel-go/blob/master/parser/unparser.go

// Resource limit constants.
const (
	// defaultMaxRecursionDepth is the default maximum recursion depth for visit()
	// to prevent stack overflow from deeply nested expressions (CWE-674: Uncontrolled Recursion).
	defaultMaxRecursionDepth = 100

	// maxComprehensionDepth is the maximum nesting depth for CEL comprehensions
	// to prevent resource exhaustion from deeply nested UNNEST/subquery operations (CWE-400).
	maxComprehensionDepth = 3

	// maxByteArrayLength is the maximum allowed length for byte arrays in non-parameterized mode
	// to prevent memory exhaustion from large hex-encoded SQL strings (CWE-400).
	// Each byte expands to ~4 characters in hex format (e.g., \xDE).
	// 10,000 bytes → ~40KB SQL output.
	maxByteArrayLength = 10000

	// defaultMaxSQLOutputLength is the default maximum length of generated SQL output
	// to prevent resource exhaustion from extremely large SQL queries (CWE-400).
	defaultMaxSQLOutputLength = 50000
)

// ConvertOption is a functional option for configuring the Convert function.
type ConvertOption func(*convertOptions)

// convertOptions holds configuration options for the Convert function.
type convertOptions struct {
	schemas      map[string]schema.Schema
	ctx          context.Context
	logger       *slog.Logger
	maxDepth     int             // Maximum recursion depth (0 = use default)
	maxOutputLen int             // Maximum SQL output length (0 = use default)
	dialect      dialect.Dialect // SQL dialect (nil = PostgreSQL default)
}

// WithDialect sets the SQL dialect for conversion.
// If not provided, PostgreSQL is used as the default dialect.
//
// Example:
//
//	import "github.com/spandigital/cel2sql/v3/dialect/mysql"
//
//	sql, err := cel2sql.Convert(ast, cel2sql.WithDialect(mysql.New()))
func WithDialect(d dialect.Dialect) ConvertOption {
	return func(o *convertOptions) {
		o.dialect = d
	}
}

// WithSchemas provides schema information for proper JSON/JSONB field handling.
// This option is required for correct SQL generation when using JSON/JSONB fields.
//
// Example:
//
//	schemas := provider.GetSchemas()
//	sql, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
func WithSchemas(schemas map[string]schema.Schema) ConvertOption {
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
	SQL        string // The generated SQL WHERE clause with placeholders
	Parameters []any  // Parameter values in order ($1, $2, etc.)
}

// Convert converts a CEL AST to a SQL WHERE clause condition.
// By default, PostgreSQL SQL is generated. Use WithDialect to select a different dialect.
//
// Example without options (PostgreSQL):
//
//	sql, err := cel2sql.Convert(ast)
//
// Example with schema information for JSON/JSONB support:
//
//	sql, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
//
// Example with a different dialect:
//
//	sql, err := cel2sql.Convert(ast, cel2sql.WithDialect(mysql.New()))
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

	// Default to PostgreSQL dialect if none specified
	if options.dialect == nil {
		options.dialect = postgres.New()
	}

	options.logger.Debug("starting CEL to SQL conversion")

	checkedExpr, err := cel.AstToCheckedExpr(ast)
	if err != nil {
		options.logger.Error("AST to CheckedExpr conversion failed", slog.Any("error", err))
		return "", err
	}

	un := &converter{
		typeMap:      checkedExpr.TypeMap,
		schemas:      options.schemas,
		ctx:          options.ctx,
		logger:       options.logger,
		dialect:      options.dialect,
		maxDepth:     options.maxDepth,
		maxOutputLen: options.maxOutputLen,
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
		slog.String("dialect", string(options.dialect.Name())),
		slog.Duration("duration", duration),
	)

	return result, nil
}

// ConvertParameterized converts a CEL AST to a parameterized SQL WHERE clause.
// Returns both the SQL string with placeholders and the parameter values.
// By default uses PostgreSQL ($1, $2). Use WithDialect for other placeholder styles.
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

	// Default to PostgreSQL dialect if none specified
	if options.dialect == nil {
		options.dialect = postgres.New()
	}

	options.logger.Debug("starting parameterized CEL to SQL conversion")

	checkedExpr, err := cel.AstToCheckedExpr(ast)
	if err != nil {
		options.logger.Error("AST to CheckedExpr conversion failed", slog.Any("error", err))
		return nil, err
	}

	un := &converter{
		typeMap:      checkedExpr.TypeMap,
		schemas:      options.schemas,
		ctx:          options.ctx,
		logger:       options.logger,
		dialect:      options.dialect,
		maxDepth:     options.maxDepth,
		maxOutputLen: options.maxOutputLen,
		parameterize: true, // Enable parameterization
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
	str                strings.Builder
	typeMap            map[int64]*exprpb.Type
	schemas            map[string]schema.Schema
	ctx                context.Context
	logger             *slog.Logger
	dialect            dialect.Dialect
	depth              int   // Current recursion depth
	maxDepth           int   // Maximum allowed recursion depth
	maxOutputLen       int   // Maximum allowed SQL output length
	comprehensionDepth int   // Current comprehension nesting depth
	parameterize       bool  // Enable parameterized output
	parameters         []any // Collected parameters for parameterized queries
	paramCount         int   // Parameter counter for placeholders
}

// checkContext checks if the context has been cancelled or expired.
// Returns nil if no context was provided or if the context is still active.
// Returns an error if the context has been cancelled or its deadline has exceeded.
func (con *converter) checkContext() error {
	if con.ctx == nil {
		return nil
	}
	if err := con.ctx.Err(); err != nil {
		return fmt.Errorf("%w: %w", ErrContextCanceled, err)
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
		return fmt.Errorf("%w: depth %d exceeds limit of %d", ErrMaxDepthExceeded, con.depth, con.maxDepth)
	}

	// Check for context cancellation at the main recursion entry point
	if err := con.checkContext(); err != nil {
		return err
	}

	// Check SQL output length limit to prevent resource exhaustion (CWE-400)
	if con.str.Len() > con.maxOutputLen {
		return fmt.Errorf("%w: %d bytes exceeds limit of %d", ErrMaxOutputLengthExceeded, con.str.Len(), con.maxOutputLen)
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

	for _, field := range schema.Fields() {
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

	for _, field := range schema.Fields() {
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

	for _, field := range schema.Fields() {
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

	for _, field := range schema.Fields() {
		if field.Name == fieldName && field.Repeated {
			return field.ElementType
		}
	}

	return ""
}

// getArrayDimension returns the number of array dimensions for a field expression.
// Returns 1 if no schema information is available (backward compatible default).
// For multi-dimensional arrays, returns the detected dimension count (2 for int[][], 3 for int[][][], etc.)
func (con *converter) getArrayDimension(expr *exprpb.Expr) int {
	// Default to 1D arrays if we can't determine from schema
	if con.schemas == nil {
		return 1
	}

	// Try to extract field name from the select expression
	selectExpr := expr.GetSelectExpr()
	if selectExpr == nil {
		return 1
	}

	fieldName := selectExpr.GetField()
	operand := selectExpr.GetOperand()

	// Get the type of the operand from the type map
	operandType := con.typeMap[operand.GetId()]
	if operandType == nil {
		return 1
	}

	// Extract the type name (e.g., "TestTable" from the object type)
	typeName := operandType.GetMessageType()
	if typeName == "" {
		return 1
	}

	// Look up the schema by type name
	schema, ok := con.schemas[typeName]
	if !ok {
		return 1
	}

	// Find the field in the schema
	field, found := schema.FindField(fieldName)
	if !found || !field.Repeated {
		return 1
	}

	// If dimensions is explicitly set and > 0, use it
	if field.Dimensions > 0 {
		return field.Dimensions
	}

	// Otherwise default to 1
	return 1
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

	// Handle string concatenation via dialect before writing LHS.
	// This allows MySQL to use CONCAT() instead of ||.
	if fun == operators.Add &&
		((lhsType.GetPrimitive() == exprpb.Type_STRING && rhsType.GetPrimitive() == exprpb.Type_STRING) ||
			(isStringLiteral(lhs) || isStringLiteral(rhs))) {
		return con.dialect.WriteStringConcat(&con.str,
			func() error { return con.visitMaybeNested(lhs, lhsParen) },
			func() error { return con.visitMaybeNested(rhs, rhsParen) },
		)
	}

	// Handle array membership (IN operator with list) via dialect before writing LHS.
	// This allows dialects like SQLite to use a fundamentally different pattern
	// (e.g., "elem IN (SELECT value FROM json_each(array))") instead of "elem = ANY(array)".
	if fun == operators.In && isListType(rhsType) {
		// Non-JSON list membership
		if !isFieldAccessExpression(rhs) || !con.isJSONArrayField(rhs) {
			return con.dialect.WriteArrayMembership(&con.str,
				func() error { return con.visitMaybeNested(lhs, lhsParen) },
				func() error { return con.visitMaybeNested(rhs, rhsParen) },
			)
		}
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
		con.str.WriteString(")")
		con.dialect.WriteCastToNumeric(&con.str)
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
			// For JSON arrays, use dialect-specific JSON array membership
			jsonFunc := con.getJSONArrayFunction(rhs)

			// For nested JSON access like settings.permissions, we need to handle differently
			if con.isNestedJSONAccess(rhs) {
				// Use dialect-specific nested JSON array membership
				if err := con.dialect.WriteNestedJSONArrayMembership(&con.str, func() error {
					return con.visitNestedJSONForArray(rhs)
				}); err != nil {
					return err
				}
				return nil
			}
			// For direct JSON array access
			if err := con.dialect.WriteJSONArrayMembership(&con.str, jsonFunc, func() error {
				return con.visitMaybeNested(rhs, rhsParen)
			}); err != nil {
				return err
			}
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

	// For regular strings, use dialect-specific contains
	return con.dialect.WriteContains(&con.str,
		func() error {
			if target != nil {
				nested := isBinaryOrTernaryOperator(target)
				return con.visitMaybeNested(target, nested)
			}
			return nil
		},
		func() error {
			for i, arg := range args {
				if err := con.visit(arg); err != nil {
					return err
				}
				if i < len(args)-1 {
					con.str.WriteString(", ")
				}
			}
			return nil
		},
	)
}

func (con *converter) callStartsWith(target *exprpb.Expr, args []*exprpb.Expr) error {
	// CEL startsWith function: string.startsWith(prefix)
	// Convert to PostgreSQL: string LIKE 'prefix%'
	// or for more robust handling: LEFT(string, LENGTH(prefix)) = prefix

	if target == nil || len(args) == 0 {
		return fmt.Errorf("%w: startsWith function requires both string and prefix arguments", ErrInvalidArguments)
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
			return fmt.Errorf("%w: LIKE patterns cannot contain null bytes", ErrInvalidArguments)
		}
		// Escape special LIKE characters: %, _, \
		escaped := escapeLikePattern(prefix)
		con.str.WriteString("'")
		con.str.WriteString(escaped)
		con.str.WriteString("%'")
		con.dialect.WriteLikeEscape(&con.str)
	} else {
		// For non-literal patterns, escape special characters at runtime and concatenate with %
		con.str.WriteString("REPLACE(REPLACE(REPLACE(")
		if err := con.visit(args[0]); err != nil {
			return err
		}
		con.str.WriteString(", '\\\\', '\\\\\\\\'), '%', '\\%'), '_', '\\_') || '%'")
		con.dialect.WriteLikeEscape(&con.str)
	}

	return nil
}

func (con *converter) callEndsWith(target *exprpb.Expr, args []*exprpb.Expr) error {
	// CEL endsWith function: string.endsWith(suffix)
	// Convert to SQL: string LIKE '%suffix'

	if target == nil || len(args) == 0 {
		return fmt.Errorf("%w: endsWith function requires both string and suffix arguments", ErrInvalidArguments)
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
			return fmt.Errorf("%w: LIKE patterns cannot contain null bytes", ErrInvalidArguments)
		}
		// Escape special LIKE characters: %, _, \
		escaped := escapeLikePattern(suffix)
		con.str.WriteString("'%")
		con.str.WriteString(escaped)
		con.str.WriteString("'")
		con.dialect.WriteLikeEscape(&con.str)
	} else {
		// For non-literal patterns, escape special characters at runtime and concatenate with %
		con.str.WriteString("'%' || REPLACE(REPLACE(REPLACE(")
		if err := con.visit(args[0]); err != nil {
			return err
		}
		con.str.WriteString(", '\\\\', '\\\\\\\\'), '%', '\\%'), '_', '\\_')")
		con.dialect.WriteLikeEscape(&con.str)
	}

	return nil
}

func (con *converter) callCasting(function string, _ *exprpb.Expr, args []*exprpb.Expr) error {
	if len(args) == 0 {
		return fmt.Errorf("%w: type conversion requires an argument", ErrInvalidArguments)
	}
	arg := args[0]
	if function == overloads.TypeConvertInt && isTimestampType(con.getType(arg)) {
		return con.dialect.WriteEpochExtract(&con.str, func() error {
			return con.visit(arg)
		})
	}
	con.str.WriteString("CAST(")
	if err := con.visit(arg); err != nil {
		return err
	}
	con.str.WriteString(" AS ")
	// Map CEL type conversion function to dialect-specific type name
	var celTypeName string
	switch function {
	case overloads.TypeConvertBool:
		celTypeName = "bool"
	case overloads.TypeConvertBytes:
		celTypeName = "bytes"
	case overloads.TypeConvertDouble:
		celTypeName = "double"
	case overloads.TypeConvertInt:
		celTypeName = "int"
	case overloads.TypeConvertString:
		celTypeName = "string"
	case overloads.TypeConvertUint:
		celTypeName = "uint"
	}
	con.dialect.WriteTypeName(&con.str, celTypeName)
	con.str.WriteString(")")
	return nil
}

// callMatches handles CEL matches() function with regex conversion
func (con *converter) callMatches(target *exprpb.Expr, args []*exprpb.Expr) error {
	// CEL matches function: string.matches(pattern) or matches(string, pattern)

	// Check if the dialect supports regex
	if !con.dialect.SupportsRegex() {
		return fmt.Errorf("%w: regex matching is not supported by %s dialect", ErrUnsupportedDialectFeature, con.dialect.Name())
	}

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
		return fmt.Errorf("%w: matches function requires both string and pattern arguments", ErrInvalidArguments)
	}

	// Visit the pattern expression and convert if it's a string literal
	if constExpr := patternExpr.GetConstExpr(); constExpr != nil && constExpr.GetStringValue() != "" {
		re2Pattern := constExpr.GetStringValue()
		// Reject patterns containing null bytes
		if strings.Contains(re2Pattern, "\x00") {
			return fmt.Errorf("%w: regex patterns cannot contain null bytes", ErrInvalidRegexPattern)
		}

		// Convert RE2 to dialect-native format with security validation
		convertedPattern, caseInsensitive, err := con.dialect.ConvertRegex(re2Pattern)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrInvalidRegexPattern, err)
		}

		con.logger.LogAttrs(context.Background(), slog.LevelDebug,
			"regex pattern conversion",
			slog.String("original_pattern", re2Pattern),
			slog.String("converted_pattern", convertedPattern),
			slog.Bool("case_insensitive", caseInsensitive),
			slog.String("dialect", string(con.dialect.Name())),
		)

		// Use dialect-specific regex match writing
		return con.dialect.WriteRegexMatch(&con.str, func() error {
			return con.visit(stringExpr)
		}, convertedPattern, caseInsensitive)
	}
	// For non-literal patterns, we can't convert at compile time
	// Visit the string, then write regex operator, then visit the pattern
	if err := con.visit(stringExpr); err != nil {
		return err
	}
	con.str.WriteString(" ~ ")
	return con.visit(patternExpr)
}

// callLowerASCII handles CEL lowerAscii() string function
func (con *converter) callLowerASCII(target *exprpb.Expr, args []*exprpb.Expr) error {
	// CEL lowerAscii function: string.lowerAscii()
	// Convert to PostgreSQL: LOWER(string)

	var stringExpr *exprpb.Expr
	switch {
	case target != nil:
		// Method call: string.lowerAscii()
		stringExpr = target
	case len(args) > 0:
		// Function call: lowerAscii(string)
		stringExpr = args[0]
	default:
		return fmt.Errorf("%w: lowerAscii() requires a string argument", ErrInvalidArguments)
	}

	con.str.WriteString("LOWER(")
	nested := isBinaryOrTernaryOperator(stringExpr)
	if err := con.visitMaybeNested(stringExpr, nested); err != nil {
		return err
	}
	con.str.WriteString(")")
	return nil
}

// callUpperASCII handles CEL upperAscii() string function
func (con *converter) callUpperASCII(target *exprpb.Expr, args []*exprpb.Expr) error {
	// CEL upperAscii function: string.upperAscii()
	// Convert to PostgreSQL: UPPER(string)

	var stringExpr *exprpb.Expr
	switch {
	case target != nil:
		// Method call: string.upperAscii()
		stringExpr = target
	case len(args) > 0:
		// Function call: upperAscii(string)
		stringExpr = args[0]
	default:
		return fmt.Errorf("%w: upperAscii() requires a string argument", ErrInvalidArguments)
	}

	con.str.WriteString("UPPER(")
	nested := isBinaryOrTernaryOperator(stringExpr)
	if err := con.visitMaybeNested(stringExpr, nested); err != nil {
		return err
	}
	con.str.WriteString(")")
	return nil
}

// callTrim handles CEL trim() string function
func (con *converter) callTrim(target *exprpb.Expr, args []*exprpb.Expr) error {
	// CEL trim function: string.trim()
	// Convert to PostgreSQL: TRIM(string)

	var stringExpr *exprpb.Expr
	switch {
	case target != nil:
		// Method call: string.trim()
		stringExpr = target
	case len(args) > 0:
		// Function call: trim(string)
		stringExpr = args[0]
	default:
		return fmt.Errorf("%w: trim() requires a string argument", ErrInvalidArguments)
	}

	con.str.WriteString("TRIM(")
	nested := isBinaryOrTernaryOperator(stringExpr)
	if err := con.visitMaybeNested(stringExpr, nested); err != nil {
		return err
	}
	con.str.WriteString(")")
	return nil
}

// callCharAt handles CEL charAt() string function
func (con *converter) callCharAt(target *exprpb.Expr, args []*exprpb.Expr) error {
	// CEL charAt function: string.charAt(index)
	// Convert to PostgreSQL: SUBSTRING(string, index+1, 1)
	// Note: CEL is 0-indexed, PostgreSQL SUBSTRING is 1-indexed

	var stringExpr *exprpb.Expr
	var indexExpr *exprpb.Expr

	if target != nil {
		// Method call: string.charAt(index)
		stringExpr = target
		if len(args) > 0 {
			indexExpr = args[0]
		}
	} else if len(args) >= 2 {
		// Function call: charAt(string, index)
		stringExpr = args[0]
		indexExpr = args[1]
	}

	if stringExpr == nil || indexExpr == nil {
		return fmt.Errorf("%w: charAt() requires both string and index arguments", ErrInvalidArguments)
	}

	con.str.WriteString("SUBSTRING(")
	nested := isBinaryOrTernaryOperator(stringExpr)
	if err := con.visitMaybeNested(stringExpr, nested); err != nil {
		return err
	}
	con.str.WriteString(", ")

	// Convert 0-indexed to 1-indexed
	// If index is a constant, we can add 1 at compile time
	if constExpr := indexExpr.GetConstExpr(); constExpr != nil {
		idx := constExpr.GetInt64Value()
		con.str.WriteString(strconv.FormatInt(idx+1, 10))
	} else {
		// For dynamic index, add 1 at runtime
		if err := con.visit(indexExpr); err != nil {
			return err
		}
		con.str.WriteString(" + 1")
	}

	con.str.WriteString(", 1)")
	return nil
}

// callIndexOf handles CEL indexOf() string function
func (con *converter) callIndexOf(target *exprpb.Expr, args []*exprpb.Expr) error {
	// CEL indexOf function: string.indexOf(substring) or string.indexOf(substring, offset)
	// Convert to PostgreSQL: POSITION(substring IN string) - 1 (to convert to 0-indexed)
	// Note: PostgreSQL POSITION is 1-indexed and returns 0 for not found, CEL returns -1 for not found

	var stringExpr *exprpb.Expr
	var substringExpr *exprpb.Expr
	var offsetExpr *exprpb.Expr

	if target != nil {
		// Method call: string.indexOf(substring [, offset])
		stringExpr = target
		if len(args) > 0 {
			substringExpr = args[0]
		}
		if len(args) > 1 {
			offsetExpr = args[1]
		}
	} else if len(args) >= 2 {
		// Function call: indexOf(string, substring [, offset])
		stringExpr = args[0]
		substringExpr = args[1]
		if len(args) > 2 {
			offsetExpr = args[2]
		}
	}

	if stringExpr == nil || substringExpr == nil {
		return fmt.Errorf("%w: indexOf() requires both string and substring arguments", ErrInvalidArguments)
	}

	if offsetExpr != nil {
		// With offset: use SUBSTRING to search from offset, then adjust result
		con.str.WriteString("CASE WHEN POSITION(")
		if err := con.visit(substringExpr); err != nil {
			return err
		}
		con.str.WriteString(" IN SUBSTRING(")
		nested := isBinaryOrTernaryOperator(stringExpr)
		if err := con.visitMaybeNested(stringExpr, nested); err != nil {
			return err
		}
		con.str.WriteString(", ")
		// Convert 0-indexed offset to 1-indexed
		if constExpr := offsetExpr.GetConstExpr(); constExpr != nil {
			offset := constExpr.GetInt64Value()
			con.str.WriteString(strconv.FormatInt(offset+1, 10))
		} else {
			if err := con.visit(offsetExpr); err != nil {
				return err
			}
			con.str.WriteString(" + 1")
		}
		con.str.WriteString(")) > 0 THEN POSITION(")
		if err := con.visit(substringExpr); err != nil {
			return err
		}
		con.str.WriteString(" IN SUBSTRING(")
		nested = isBinaryOrTernaryOperator(stringExpr)
		if err := con.visitMaybeNested(stringExpr, nested); err != nil {
			return err
		}
		con.str.WriteString(", ")
		if constExpr := offsetExpr.GetConstExpr(); constExpr != nil {
			offset := constExpr.GetInt64Value()
			con.str.WriteString(strconv.FormatInt(offset+1, 10))
		} else {
			if err := con.visit(offsetExpr); err != nil {
				return err
			}
			con.str.WriteString(" + 1")
		}
		con.str.WriteString(")) + ")
		if err := con.visit(offsetExpr); err != nil {
			return err
		}
		con.str.WriteString(" - 1 ELSE -1 END")
	} else {
		// Without offset: POSITION(substring IN string) - 1, return -1 if not found
		con.str.WriteString("CASE WHEN POSITION(")
		if err := con.visit(substringExpr); err != nil {
			return err
		}
		con.str.WriteString(" IN ")
		nested := isBinaryOrTernaryOperator(stringExpr)
		if err := con.visitMaybeNested(stringExpr, nested); err != nil {
			return err
		}
		con.str.WriteString(") > 0 THEN POSITION(")
		if err := con.visit(substringExpr); err != nil {
			return err
		}
		con.str.WriteString(" IN ")
		nested = isBinaryOrTernaryOperator(stringExpr)
		if err := con.visitMaybeNested(stringExpr, nested); err != nil {
			return err
		}
		con.str.WriteString(") - 1 ELSE -1 END")
	}

	return nil
}

// callLastIndexOf handles CEL lastIndexOf() string function
func (con *converter) callLastIndexOf(target *exprpb.Expr, args []*exprpb.Expr) error {
	// CEL lastIndexOf function: string.lastIndexOf(substring)
	// Convert to PostgreSQL: LENGTH(string) - POSITION(REVERSE(substring) IN REVERSE(string)) - LENGTH(substring) + 1
	// Returns -1 if not found (CEL convention)

	var stringExpr *exprpb.Expr
	var substringExpr *exprpb.Expr

	if target != nil {
		// Method call: string.lastIndexOf(substring)
		stringExpr = target
		if len(args) > 0 {
			substringExpr = args[0]
		}
	} else if len(args) >= 2 {
		// Function call: lastIndexOf(string, substring)
		stringExpr = args[0]
		substringExpr = args[1]
	}

	if stringExpr == nil || substringExpr == nil {
		return fmt.Errorf("%w: lastIndexOf() requires both string and substring arguments", ErrInvalidArguments)
	}

	// Return -1 if not found, otherwise calculate position
	con.str.WriteString("CASE WHEN POSITION(REVERSE(")
	if err := con.visit(substringExpr); err != nil {
		return err
	}
	con.str.WriteString(") IN REVERSE(")
	nested := isBinaryOrTernaryOperator(stringExpr)
	if err := con.visitMaybeNested(stringExpr, nested); err != nil {
		return err
	}
	con.str.WriteString(")) > 0 THEN LENGTH(")
	nested = isBinaryOrTernaryOperator(stringExpr)
	if err := con.visitMaybeNested(stringExpr, nested); err != nil {
		return err
	}
	con.str.WriteString(") - POSITION(REVERSE(")
	if err := con.visit(substringExpr); err != nil {
		return err
	}
	con.str.WriteString(") IN REVERSE(")
	nested = isBinaryOrTernaryOperator(stringExpr)
	if err := con.visitMaybeNested(stringExpr, nested); err != nil {
		return err
	}
	con.str.WriteString(")) - LENGTH(")
	if err := con.visit(substringExpr); err != nil {
		return err
	}
	con.str.WriteString(") + 1 ELSE -1 END")

	return nil
}

// callSubstring handles CEL substring() string function
func (con *converter) callSubstring(target *exprpb.Expr, args []*exprpb.Expr) error {
	// CEL substring function: string.substring(start) or string.substring(start, end)
	// Convert to PostgreSQL: SUBSTRING(string, start+1 [, end-start])
	// Note: CEL is 0-indexed and end is exclusive, PostgreSQL SUBSTRING is 1-indexed

	var stringExpr *exprpb.Expr
	var startExpr *exprpb.Expr
	var endExpr *exprpb.Expr

	if target != nil {
		// Method call: string.substring(start [, end])
		stringExpr = target
		if len(args) > 0 {
			startExpr = args[0]
		}
		if len(args) > 1 {
			endExpr = args[1]
		}
	} else if len(args) >= 2 {
		// Function call: substring(string, start [, end])
		stringExpr = args[0]
		startExpr = args[1]
		if len(args) > 2 {
			endExpr = args[2]
		}
	}

	if stringExpr == nil || startExpr == nil {
		return fmt.Errorf("%w: substring() requires string and start arguments", ErrInvalidArguments)
	}

	con.str.WriteString("SUBSTRING(")
	nested := isBinaryOrTernaryOperator(stringExpr)
	if err := con.visitMaybeNested(stringExpr, nested); err != nil {
		return err
	}
	con.str.WriteString(", ")

	// Convert 0-indexed start to 1-indexed
	if constExpr := startExpr.GetConstExpr(); constExpr != nil {
		start := constExpr.GetInt64Value()
		con.str.WriteString(strconv.FormatInt(start+1, 10))
	} else {
		if err := con.visit(startExpr); err != nil {
			return err
		}
		con.str.WriteString(" + 1")
	}

	// If end is provided, calculate length as (end - start)
	if endExpr != nil {
		con.str.WriteString(", ")
		// If both start and end are constants, calculate length at compile time
		if startConst := startExpr.GetConstExpr(); startConst != nil {
			if endConst := endExpr.GetConstExpr(); endConst != nil {
				start := startConst.GetInt64Value()
				end := endConst.GetInt64Value()
				length := max(end-start, 0)
				con.str.WriteString(strconv.FormatInt(length, 10))
			} else {
				// End is dynamic, start is constant
				if err := con.visit(endExpr); err != nil {
					return err
				}
				con.str.WriteString(" - ")
				start := startConst.GetInt64Value()
				con.str.WriteString(strconv.FormatInt(start, 10))
			}
		} else {
			// Start is dynamic
			if err := con.visit(endExpr); err != nil {
				return err
			}
			con.str.WriteString(" - (")
			if err := con.visit(startExpr); err != nil {
				return err
			}
			con.str.WriteString(")")
		}
	}

	con.str.WriteString(")")
	return nil
}

// callReplace handles CEL replace() string function
func (con *converter) callReplace(target *exprpb.Expr, args []*exprpb.Expr) error {
	// CEL replace function: string.replace(old, new) or string.replace(old, new, limit)
	// Convert to PostgreSQL: REPLACE(string, old, new)
	// Note: PostgreSQL REPLACE replaces all occurrences, limit parameter not supported

	var stringExpr *exprpb.Expr
	var oldExpr *exprpb.Expr
	var newExpr *exprpb.Expr
	var limitExpr *exprpb.Expr

	if target != nil {
		// Method call: string.replace(old, new [, limit])
		stringExpr = target
		if len(args) > 0 {
			oldExpr = args[0]
		}
		if len(args) > 1 {
			newExpr = args[1]
		}
		if len(args) > 2 {
			limitExpr = args[2]
		}
	} else if len(args) >= 3 {
		// Function call: replace(string, old, new [, limit])
		stringExpr = args[0]
		oldExpr = args[1]
		newExpr = args[2]
		if len(args) > 3 {
			limitExpr = args[3]
		}
	}

	if stringExpr == nil || oldExpr == nil || newExpr == nil {
		return fmt.Errorf("%w: replace() requires string, old, and new arguments", ErrInvalidArguments)
	}

	// Check if limit is provided and is not -1 (replace all)
	if limitExpr != nil {
		if constExpr := limitExpr.GetConstExpr(); constExpr != nil {
			limit := constExpr.GetInt64Value()
			if limit != -1 {
				return fmt.Errorf("%w: replace() with limit != -1 is not supported in SQL conversion (PostgreSQL REPLACE replaces all occurrences)", ErrUnsupportedOperation)
			}
		} else {
			// Dynamic limit - we can't determine if it's -1 at compile time
			return fmt.Errorf("%w: replace() with dynamic limit is not supported in SQL conversion", ErrUnsupportedOperation)
		}
	}

	con.str.WriteString("REPLACE(")
	nested := isBinaryOrTernaryOperator(stringExpr)
	if err := con.visitMaybeNested(stringExpr, nested); err != nil {
		return err
	}
	con.str.WriteString(", ")
	if err := con.visit(oldExpr); err != nil {
		return err
	}
	con.str.WriteString(", ")
	if err := con.visit(newExpr); err != nil {
		return err
	}
	con.str.WriteString(")")
	return nil
}

// callReverse handles CEL reverse() string function
func (con *converter) callReverse(target *exprpb.Expr, args []*exprpb.Expr) error {
	// CEL reverse function: string.reverse()
	// Convert to PostgreSQL: REVERSE(string)

	var stringExpr *exprpb.Expr
	switch {
	case target != nil:
		// Method call: string.reverse()
		stringExpr = target
	case len(args) > 0:
		// Function call: reverse(string)
		stringExpr = args[0]
	default:
		return fmt.Errorf("%w: reverse() requires a string argument", ErrInvalidArguments)
	}

	con.str.WriteString("REVERSE(")
	nested := isBinaryOrTernaryOperator(stringExpr)
	if err := con.visitMaybeNested(stringExpr, nested); err != nil {
		return err
	}
	con.str.WriteString(")")
	return nil
}

// callSplit handles CEL split() string function
func (con *converter) callSplit(target *exprpb.Expr, args []*exprpb.Expr) error {
	// CEL split function: string.split(delimiter) or string.split(delimiter, limit)
	// Convert to PostgreSQL: STRING_TO_ARRAY(string, delimiter)
	// With limit support:
	//   limit = -1 or no limit: STRING_TO_ARRAY(string, delimiter) (unlimited)
	//   limit = 0: ARRAY[]::text[] (empty array)
	//   limit = 1: ARRAY[string] (no split)
	//   limit > 1: Complex SQL with REGEXP_SPLIT_TO_ARRAY and array slicing

	var stringExpr *exprpb.Expr
	var delimiterExpr *exprpb.Expr
	var limitExpr *exprpb.Expr

	if target != nil {
		// Method call: string.split(delimiter [, limit])
		stringExpr = target
		if len(args) > 0 {
			delimiterExpr = args[0]
		}
		if len(args) > 1 {
			limitExpr = args[1]
		}
	} else if len(args) >= 2 {
		// Function call: split(string, delimiter [, limit])
		stringExpr = args[0]
		delimiterExpr = args[1]
		if len(args) > 2 {
			limitExpr = args[2]
		}
	}

	if stringExpr == nil || delimiterExpr == nil {
		return fmt.Errorf("%w: split() requires string and delimiter arguments", ErrInvalidArguments)
	}

	// Validate delimiter for security (check for null bytes)
	if constExpr := delimiterExpr.GetConstExpr(); constExpr != nil {
		if strVal := constExpr.GetStringValue(); strings.ContainsRune(strVal, '\x00') {
			return fmt.Errorf("%w: split() delimiter cannot contain null bytes", ErrInvalidArguments)
		}
	}

	// Handle limit parameter
	var limit int64 = -1 // Default: unlimited splits
	if limitExpr != nil {
		if constExpr := limitExpr.GetConstExpr(); constExpr != nil {
			limit = constExpr.GetInt64Value()
		} else {
			return fmt.Errorf("%w: split() with dynamic limit is not supported in SQL conversion", ErrUnsupportedOperation)
		}
	}

	// Generate SQL based on limit value
	writeStr := func() error {
		nested := isBinaryOrTernaryOperator(stringExpr)
		return con.visitMaybeNested(stringExpr, nested)
	}
	writeDelim := func() error {
		return con.visit(delimiterExpr)
	}

	switch {
	case limit == 0:
		// Empty array
		con.dialect.WriteEmptyTypedArray(&con.str, "text")
		return nil

	case limit == 1:
		// Return original string as single-element array
		con.dialect.WriteArrayLiteralOpen(&con.str)
		if err := writeStr(); err != nil {
			return err
		}
		con.dialect.WriteArrayLiteralClose(&con.str)
		return nil

	case limit == -1:
		// Unlimited splits
		return con.dialect.WriteSplit(&con.str, writeStr, writeDelim)

	case limit > 1:
		// Positive limit - use dialect-specific split with limit
		return con.dialect.WriteSplitWithLimit(&con.str, writeStr, writeDelim, limit)

	default:
		// Negative limits other than -1 are not supported
		return fmt.Errorf("%w: split() with negative limit other than -1 is not supported", ErrUnsupportedOperation)
	}
}

// callJoin handles CEL join() function
func (con *converter) callJoin(target *exprpb.Expr, args []*exprpb.Expr) error {
	// CEL join function: array.join() or array.join(delimiter)
	// Convert to PostgreSQL: ARRAY_TO_STRING(array, delimiter, '')
	// Default delimiter is empty string if not provided

	var arrayExpr *exprpb.Expr
	var delimiterExpr *exprpb.Expr

	if target != nil {
		// Method call: array.join([delimiter])
		arrayExpr = target
		if len(args) > 0 {
			delimiterExpr = args[0]
		}
	} else if len(args) >= 1 {
		// Function call: join(array [, delimiter])
		arrayExpr = args[0]
		if len(args) > 1 {
			delimiterExpr = args[1]
		}
	}

	if arrayExpr == nil {
		return fmt.Errorf("%w: join() requires an array argument", ErrInvalidArguments)
	}

	// Validate delimiter for security (check for null bytes)
	if delimiterExpr != nil {
		if constExpr := delimiterExpr.GetConstExpr(); constExpr != nil {
			if strVal := constExpr.GetStringValue(); strings.ContainsRune(strVal, '\x00') {
				return fmt.Errorf("%w: join() delimiter cannot contain null bytes", ErrInvalidArguments)
			}
		}
	}

	// Generate SQL using dialect-specific join
	writeArray := func() error {
		nested := isBinaryOrTernaryOperator(arrayExpr)
		return con.visitMaybeNested(arrayExpr, nested)
	}
	var writeDelim func() error
	if delimiterExpr != nil {
		writeDelim = func() error {
			return con.visit(delimiterExpr)
		}
	}
	return con.dialect.WriteJoin(&con.str, writeArray, writeDelim)
}

// callFormat handles CEL format() function
func (con *converter) callFormat(target *exprpb.Expr, args []*exprpb.Expr) error {
	// CEL format function: format_string.format(args_list)
	// Convert to PostgreSQL: FORMAT(format_string, arg1, arg2, ...)
	// Supports: %s (string), %d (decimal/integer), %f (float)
	// Unsupported: %b (binary), %x (hex), etc.

	var formatExpr *exprpb.Expr
	var argsExpr *exprpb.Expr

	if target != nil {
		// Method call: format_string.format(args)
		formatExpr = target
		if len(args) > 0 {
			argsExpr = args[0]
		}
	} else if len(args) >= 2 {
		// Function call: format(format_string, args)
		formatExpr = args[0]
		argsExpr = args[1]
	}

	if formatExpr == nil || argsExpr == nil {
		return fmt.Errorf("%w: format() requires format string and arguments list", ErrInvalidArguments)
	}

	// Format string must be a constant
	constFormat := formatExpr.GetConstExpr()
	if constFormat == nil {
		return fmt.Errorf("%w: format() requires a constant format string", ErrUnsupportedOperation)
	}

	formatString := constFormat.GetStringValue()

	// Security: Check format string length limit (1000 chars)
	if len(formatString) > 1000 {
		return fmt.Errorf("%w: format() format string exceeds maximum length of 1000 characters", ErrInvalidArguments)
	}

	// Parse format string to extract specifiers and validate
	specifiers, err := parseFormatString(formatString)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidArguments, err)
	}

	// Arguments must be a constant list
	listExpr := argsExpr.GetListExpr()
	if listExpr == nil {
		return fmt.Errorf("%w: format() requires a constant list of arguments", ErrUnsupportedOperation)
	}

	argElements := listExpr.GetElements()

	// Validate argument count matches specifier count
	if len(argElements) != len(specifiers) {
		return fmt.Errorf("%w: format() argument count mismatch - format string has %d placeholders but got %d arguments", ErrInvalidArguments, len(specifiers), len(argElements))
	}

	// Convert CEL format specifiers to PostgreSQL format specifiers
	pgFormatString := convertFormatString(formatString, specifiers)

	// Generate SQL: FORMAT(format_string, arg1, arg2, ...)
	con.str.WriteString("FORMAT(")
	con.str.WriteString("'")
	con.str.WriteString(strings.ReplaceAll(pgFormatString, "'", "''")) // Escape single quotes
	con.str.WriteString("'")

	// Add each argument
	for _, argExpr := range argElements {
		con.str.WriteString(", ")
		if err := con.visit(argExpr); err != nil {
			return err
		}
	}

	con.str.WriteString(")")
	return nil
}

// parseFormatString extracts format specifiers from a CEL format string
// Returns list of specifiers (%s, %d, %f) and validates them
func parseFormatString(format string) ([]string, error) {
	var specifiers []string
	i := 0
	for i < len(format) {
		if format[i] == '%' {
			if i+1 >= len(format) {
				return nil, fmt.Errorf("%w: incomplete format specifier at end of string", ErrInvalidArguments)
			}

			nextChar := format[i+1]
			switch nextChar {
			case 's', 'd', 'f':
				// Supported specifiers
				specifiers = append(specifiers, "%"+string(nextChar))
				i += 2
			case '%':
				// Escaped percent sign %%
				i += 2
			case 'b', 'x', 'X', 'o', 'e', 'E', 'g', 'G':
				// Unsupported specifiers
				return nil, fmt.Errorf("unsupported format specifier %%%c - only %%s, %%d, and %%f are supported", nextChar)
			default:
				return nil, fmt.Errorf("invalid format specifier %%%c", nextChar)
			}
		} else {
			i++
		}
	}
	return specifiers, nil
}

// convertFormatString converts CEL format string to PostgreSQL FORMAT syntax
// PostgreSQL uses %s for all types, but we keep %d and %f for type hinting
func convertFormatString(format string, _ []string) string {
	// PostgreSQL FORMAT() uses %s for strings, %I for identifiers, %L for literals
	// For our purposes, we convert all to %s and let PostgreSQL handle type conversion
	result := format
	result = strings.ReplaceAll(result, "%d", "%s")
	result = strings.ReplaceAll(result, "%f", "%s")
	// %s stays as %s
	return result
}

// callQuote handles CEL quote() function - returns error (not in ext.Strings())
func (con *converter) callQuote(_ *exprpb.Expr, _ []*exprpb.Expr) error {
	// Note: quote() is not actually part of CEL's ext.Strings() standard extension
	// It may be part of CEL spec but is not commonly implemented
	return fmt.Errorf("%w: quote() is not part of CEL ext.Strings() standard extension", ErrUnsupportedOperation)
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
	// CEL string extension functions
	case "lowerAscii":
		return con.callLowerASCII(target, args)
	case "upperAscii":
		return con.callUpperASCII(target, args)
	case "trim":
		return con.callTrim(target, args)
	case "charAt":
		return con.callCharAt(target, args)
	case "indexOf":
		return con.callIndexOf(target, args)
	case "lastIndexOf":
		return con.callLastIndexOf(target, args)
	case "substring":
		return con.callSubstring(target, args)
	case "replace":
		return con.callReplace(target, args)
	case "reverse":
		return con.callReverse(target, args)
	// Unsupported string extension functions (return errors)
	case "split":
		return con.callSplit(target, args)
	case "join":
		return con.callJoin(target, args)
	case "format":
		return con.callFormat(target, args)
	case "quote":
		return con.callQuote(target, args)
	}
	sqlFun, ok := standardSQLFunctions[fun]
	if !ok {
		if fun == overloads.Size {
			// Handle both method calls (target != nil) and function calls (len(args) > 0)
			var argExpr *exprpb.Expr
			switch {
			case target != nil:
				// Method call: t.size()
				argExpr = target
			case len(args) > 0:
				// Function call: size(t) - though this is rare for size()
				argExpr = args[0]
			default:
				return fmt.Errorf("%w: size() requires a target or argument", ErrInvalidArguments)
			}

			argType := con.getType(argExpr)
			switch {
			case argType.GetPrimitive() == exprpb.Type_STRING, argType.GetPrimitive() == exprpb.Type_BYTES:
				// For strings and bytes, directly write LENGTH(arg) and return
				con.str.WriteString("LENGTH(")
				nested := isBinaryOrTernaryOperator(argExpr)
				err := con.visitMaybeNested(argExpr, nested)
				if err != nil {
					return err
				}
				con.str.WriteString(")")
				return nil
			case isListType(argType):
				// Check if this is a JSON array field
				if con.isJSONArrayField(argExpr) {
					// For JSON arrays, use dialect-specific JSON array length
					return con.dialect.WriteJSONArrayLength(&con.str, func() error {
						return con.visit(argExpr)
					})
				}
				// For native arrays, use dialect-specific array length
				dimension := con.getArrayDimension(argExpr)
				return con.dialect.WriteArrayLength(&con.str, dimension, func() error {
					nested := isBinaryOrTernaryOperator(argExpr)
					return con.visitMaybeNested(argExpr, nested)
				})
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
	args := expr.GetCallExpr().GetArgs()
	if len(args) == 0 {
		return fmt.Errorf("%w: index operator requires at least one argument", ErrInvalidArguments)
	}
	if isMapType(con.getType(args[0])) {
		return con.visitCallMapIndex(expr)
	}
	return con.visitCallListIndex(expr)
}

func (con *converter) visitCallMapIndex(expr *exprpb.Expr) error {
	c := expr.GetCallExpr()
	args := c.GetArgs()
	if len(args) < 2 {
		return fmt.Errorf("%w: map index operator requires map and key arguments", ErrInvalidArguments)
	}
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
	if len(args) < 2 {
		return fmt.Errorf("%w: list index operator requires list and index arguments", ErrInvalidArguments)
	}
	l := args[0]
	index := args[1]

	// Check for constant index
	if constExpr := index.GetConstExpr(); constExpr != nil {
		idx := constExpr.GetInt64Value()
		if idx == math.MaxInt64 {
			return fmt.Errorf("%w: array index overflow, cannot convert math.MaxInt64 to 1-based indexing", ErrInvalidArguments)
		}
		if idx < 0 {
			return fmt.Errorf("%w: negative array index %d is not supported", ErrInvalidArguments, idx)
		}
		return con.dialect.WriteListIndexConst(&con.str, func() error {
			nested := isBinaryOrTernaryOperator(l)
			return con.visitMaybeNested(l, nested)
		}, idx)
	}

	// Dynamic index
	return con.dialect.WriteListIndex(&con.str, func() error {
		nested := isBinaryOrTernaryOperator(l)
		return con.visitMaybeNested(l, nested)
	}, func() error {
		return con.visit(index)
	})
}

func (con *converter) visitCallUnary(expr *exprpb.Expr) error {
	c := expr.GetCallExpr()
	fun := c.GetFunction()
	args := c.GetArgs()
	if len(args) == 0 {
		return fmt.Errorf("%w: unary operator requires an argument", ErrInvalidArguments)
	}
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
		return fmt.Errorf("%w: depth %d exceeds limit of %d",
			ErrMaxComprehensionDepthExceeded, con.comprehensionDepth, maxComprehensionDepth)
	}

	// Check for context cancellation before processing comprehensions (potentially expensive)
	if err := con.checkContext(); err != nil {
		return err
	}

	info, err := con.identifyComprehension(expr)
	if err != nil {
		return fmt.Errorf("%w: failed to identify comprehension: %w", ErrInvalidComprehension, err)
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
	comprehension := expr.GetComprehensionExpr()
	if comprehension == nil {
		return newConversionError(errMsgUnsupportedComprehension, "expression is not a comprehension (ALL)")
	}

	iterRange := comprehension.GetIterRange()

	con.str.WriteString("NOT EXISTS (SELECT 1 FROM ")
	if err := con.writeComprehensionSource(iterRange); err != nil {
		return wrapConversionError(err, "visiting iter range in ALL comprehension")
	}
	con.str.WriteString(" AS ")
	con.str.WriteString(info.IterVar)

	if info.Predicate != nil {
		con.str.WriteString(" WHERE NOT (")
		if err := con.visit(info.Predicate); err != nil {
			return wrapConversionError(err, "visiting predicate in ALL comprehension")
		}
		con.str.WriteString(")")
	}

	con.str.WriteString(")")
	return nil
}

func (con *converter) visitExistsComprehension(expr *exprpb.Expr, info *ComprehensionInfo) error {
	comprehension := expr.GetComprehensionExpr()
	if comprehension == nil {
		return newConversionError(errMsgUnsupportedComprehension, "expression is not a comprehension (EXISTS)")
	}

	iterRange := comprehension.GetIterRange()

	con.str.WriteString("EXISTS (SELECT 1 FROM ")
	if err := con.writeComprehensionSource(iterRange); err != nil {
		return wrapConversionError(err, "visiting iter range in EXISTS comprehension")
	}
	con.str.WriteString(" AS ")
	con.str.WriteString(info.IterVar)

	if info.Predicate != nil {
		con.str.WriteString(" WHERE ")
		if err := con.visit(info.Predicate); err != nil {
			return wrapConversionError(err, "visiting predicate in EXISTS comprehension")
		}
	}

	con.str.WriteString(")")
	return nil
}

func (con *converter) visitExistsOneComprehension(expr *exprpb.Expr, info *ComprehensionInfo) error {
	comprehension := expr.GetComprehensionExpr()
	if comprehension == nil {
		return newConversionError(errMsgUnsupportedComprehension, "expression is not a comprehension (EXISTS_ONE)")
	}

	iterRange := comprehension.GetIterRange()

	con.str.WriteString("(SELECT COUNT(*) FROM ")
	if err := con.writeComprehensionSource(iterRange); err != nil {
		return wrapConversionError(err, "visiting iter range in EXISTS_ONE comprehension")
	}
	con.str.WriteString(" AS ")
	con.str.WriteString(info.IterVar)

	if info.Predicate != nil {
		con.str.WriteString(" WHERE ")
		if err := con.visit(info.Predicate); err != nil {
			return wrapConversionError(err, "visiting predicate in EXISTS_ONE comprehension")
		}
	}

	con.str.WriteString(") = 1")
	return nil
}

func (con *converter) visitMapComprehension(expr *exprpb.Expr, info *ComprehensionInfo) error {
	comprehension := expr.GetComprehensionExpr()
	if comprehension == nil {
		return newConversionError(errMsgUnsupportedComprehension, "expression is not a comprehension (MAP)")
	}

	iterRange := comprehension.GetIterRange()

	con.dialect.WriteArraySubqueryOpen(&con.str)
	if info.Transform != nil {
		if err := con.visit(info.Transform); err != nil {
			return wrapConversionError(err, "visiting transform in MAP comprehension")
		}
	} else {
		con.str.WriteString(info.IterVar)
	}
	con.dialect.WriteArraySubqueryExprClose(&con.str)
	con.str.WriteString(" FROM ")
	if err := con.writeComprehensionSource(iterRange); err != nil {
		return wrapConversionError(err, "visiting iter range in MAP comprehension")
	}
	con.str.WriteString(" AS ")
	con.str.WriteString(info.IterVar)

	if info.Filter != nil {
		con.str.WriteString(" WHERE ")
		if err := con.visit(info.Filter); err != nil {
			return wrapConversionError(err, "visiting filter in MAP comprehension")
		}
	}

	con.str.WriteString(")")
	return nil
}

func (con *converter) visitFilterComprehension(expr *exprpb.Expr, info *ComprehensionInfo) error {
	comprehension := expr.GetComprehensionExpr()
	if comprehension == nil {
		return newConversionError(errMsgUnsupportedComprehension, "expression is not a comprehension (FILTER)")
	}

	iterRange := comprehension.GetIterRange()

	con.dialect.WriteArraySubqueryOpen(&con.str)
	con.str.WriteString(info.IterVar)
	con.dialect.WriteArraySubqueryExprClose(&con.str)
	con.str.WriteString(" FROM ")
	if err := con.writeComprehensionSource(iterRange); err != nil {
		return wrapConversionError(err, "visiting iter range in FILTER comprehension")
	}
	con.str.WriteString(" AS ")
	con.str.WriteString(info.IterVar)

	if info.Predicate != nil {
		con.str.WriteString(" WHERE ")
		if err := con.visit(info.Predicate); err != nil {
			return wrapConversionError(err, "visiting predicate in FILTER comprehension")
		}
	}

	con.str.WriteString(")")
	return nil
}

func (con *converter) visitTransformListComprehension(expr *exprpb.Expr, info *ComprehensionInfo) error {
	comprehension := expr.GetComprehensionExpr()
	if comprehension == nil {
		return newConversionError(errMsgUnsupportedComprehension, "expression is not a comprehension (TRANSFORM_LIST)")
	}

	con.dialect.WriteArraySubqueryOpen(&con.str)
	if info.Transform != nil {
		if err := con.visit(info.Transform); err != nil {
			return wrapConversionError(err, "visiting transform in TRANSFORM_LIST comprehension")
		}
	} else {
		con.str.WriteString(info.IterVar)
	}
	con.dialect.WriteArraySubqueryExprClose(&con.str)
	con.str.WriteString(" FROM ")
	if err := con.writeComprehensionSource(comprehension.GetIterRange()); err != nil {
		return wrapConversionError(err, "visiting iter range in TRANSFORM_LIST comprehension")
	}
	con.str.WriteString(" AS ")
	con.str.WriteString(info.IterVar)

	if info.Filter != nil {
		con.str.WriteString(" WHERE ")
		if err := con.visit(info.Filter); err != nil {
			return wrapConversionError(err, "visiting filter in TRANSFORM_LIST comprehension")
		}
	}

	con.str.WriteString(")")
	return nil
}

func (con *converter) visitTransformMapComprehension(_ *exprpb.Expr, _ *ComprehensionInfo) error {
	// Generate SQL for TRANSFORM_MAP comprehension: work with map entries
	// This is complex for PostgreSQL - maps are typically represented as JSON or composite types
	// For now, return an error indicating this needs special handling
	return fmt.Errorf("%w: TRANSFORM_MAP comprehension requires map/JSON support (not yet implemented)", ErrInvalidComprehension)
}

func (con *converter) visitTransformMapEntryComprehension(_ *exprpb.Expr, _ *ComprehensionInfo) error {
	// Generate SQL for TRANSFORM_MAP_ENTRY comprehension: work with map key-value pairs
	// This is complex for PostgreSQL - maps are typically represented as JSON or composite types
	// For now, return an error indicating this needs special handling
	return fmt.Errorf("%w: TRANSFORM_MAP_ENTRY comprehension requires map/JSON support (not yet implemented)", ErrInvalidComprehension)
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
			con.dialect.WriteParamPlaceholder(&con.str, con.paramCount)
			con.parameters = append(con.parameters, c.GetInt64Value())
		} else {
			i := strconv.FormatInt(c.GetInt64Value(), 10)
			con.str.WriteString(i)
		}
	case *exprpb.Constant_Uint64Value:
		if con.parameterize {
			con.paramCount++
			con.dialect.WriteParamPlaceholder(&con.str, con.paramCount)
			con.parameters = append(con.parameters, c.GetUint64Value())
		} else {
			ui := strconv.FormatUint(c.GetUint64Value(), 10)
			con.str.WriteString(ui)
		}
	case *exprpb.Constant_DoubleValue:
		if con.parameterize {
			con.paramCount++
			con.dialect.WriteParamPlaceholder(&con.str, con.paramCount)
			con.parameters = append(con.parameters, c.GetDoubleValue())
		} else {
			d := strconv.FormatFloat(c.GetDoubleValue(), 'g', -1, 64)
			con.str.WriteString(d)
		}
	case *exprpb.Constant_StringValue:
		str := c.GetStringValue()
		// Reject strings containing null bytes
		if strings.Contains(str, "\x00") {
			return fmt.Errorf("%w: string literals cannot contain null bytes", ErrInvalidArguments)
		}

		if con.parameterize {
			con.paramCount++
			con.dialect.WriteParamPlaceholder(&con.str, con.paramCount)
			con.parameters = append(con.parameters, str)
		} else {
			con.dialect.WriteStringLiteral(&con.str, str)
		}
	case *exprpb.Constant_BytesValue:
		b := c.GetBytesValue()

		if con.parameterize {
			con.paramCount++
			con.dialect.WriteParamPlaceholder(&con.str, con.paramCount)
			con.parameters = append(con.parameters, b)
		} else {
			// Validate byte array length to prevent resource exhaustion (CWE-400)
			if len(b) > maxByteArrayLength {
				return fmt.Errorf("%w: %d bytes exceeds limit of %d bytes", ErrInvalidByteArrayLength, len(b), maxByteArrayLength)
			}
			if err := con.dialect.WriteBytesLiteral(&con.str, b); err != nil {
				return err
			}
		}
	default:
		return newConversionErrorf(errMsgUnsupportedExpression, "constant type: %T", c.ConstantKind)
	}
	return nil
}

func (con *converter) visitIdent(expr *exprpb.Expr) error {
	identName := expr.GetIdentExpr().GetName()

	// Validate identifier name for security (prevent SQL injection)
	if err := con.dialect.ValidateFieldName(identName); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidFieldName, err)
	}

	// Check if this identifier needs numeric casting for JSON comprehensions
	if con.needsNumericCasting(identName) {
		con.str.WriteString("(")
		con.str.WriteString(identName)
		con.str.WriteString(")")
		con.dialect.WriteCastToNumeric(&con.str)
	} else {
		con.str.WriteString(identName)
	}
	return nil
}

func (con *converter) visitList(expr *exprpb.Expr) error {
	l := expr.GetListExpr()
	elems := l.GetElements()
	con.dialect.WriteArrayLiteralOpen(&con.str)
	for i, elem := range elems {
		err := con.visit(elem)
		if err != nil {
			return err
		}
		if i < len(elems)-1 {
			con.str.WriteString(", ")
		}
	}
	con.dialect.WriteArrayLiteralClose(&con.str)
	return nil
}

func (con *converter) visitSelect(expr *exprpb.Expr) error {
	sel := expr.GetSelectExpr()

	// Validate field name for security (prevent SQL injection)
	fieldName := sel.GetField()
	if err := con.dialect.ValidateFieldName(fieldName); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidFieldName, err)
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

	writeBase := func() error {
		return con.visitMaybeNested(sel.GetOperand(), nested)
	}

	switch {
	case useJSONPath:
		// Use dialect-specific JSON field access (text extraction)
		if err := con.dialect.WriteJSONFieldAccess(&con.str, writeBase, fieldName, true); err != nil {
			return err
		}
	case useJSONObjectAccess:
		// Use dialect-specific JSON object field access in comprehensions
		isNumeric := con.isNumericJSONField(fieldName)
		if isNumeric {
			con.str.WriteString("(")
		}
		if err := con.dialect.WriteJSONFieldAccess(&con.str, writeBase, fieldName, true); err != nil {
			return err
		}
		if isNumeric {
			// Close parentheses and add numeric cast
			con.str.WriteString(")")
			con.dialect.WriteCastToNumeric(&con.str)
		}
	default:
		// Regular field selection
		if err := writeBase(); err != nil {
			return err
		}
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
		isJSONB := con.isJSONBField(operand)
		return con.dialect.WriteJSONExistence(&con.str, isJSONB, field, func() error {
			return con.visitMaybeNested(operand, isBinaryOrTernaryOperator(operand))
		})
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
		if slices.Contains(jsonFields, parentField) {
			return true
		}
	}

	return false
}

// visitNestedJSONHas handles has() for deeply nested JSON paths
func (con *converter) visitNestedJSONHas(expr *exprpb.Expr) error {
	// Get the root JSON column and remaining path segments
	rootColumn, pathSegments := con.getJSONRootAndPath(expr)

	return con.dialect.WriteJSONExtractPath(&con.str, pathSegments, func() error {
		return con.visitJSONColumnReference(rootColumn)
	})
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
	con.dialect.WriteStructOpen(&con.str)
	for i, entry := range entries {
		v := entry.GetValue()
		if err := con.visit(v); err != nil {
			return err
		}
		if i < len(entries)-1 {
			con.str.WriteString(", ")
		}
	}
	con.dialect.WriteStructClose(&con.str)
	return nil
}

// writeComprehensionSource writes the source expression for a comprehension (UNNEST or JSON function).
func (con *converter) writeComprehensionSource(iterRange *exprpb.Expr) error {
	isJSONArray := con.isJSONArrayField(iterRange)
	if isJSONArray {
		jsonFunc := con.getJSONArrayFunction(iterRange)
		isJSONB := con.isJSONBField(iterRange)
		// Determine if we need text extraction or object extraction
		asText := strings.HasSuffix(jsonFunc, "_text")
		return con.dialect.WriteJSONArrayElements(&con.str, isJSONB, asText, func() error {
			return con.visit(iterRange)
		})
	}
	return con.dialect.WriteUnnest(&con.str, func() error {
		return con.visit(iterRange)
	})
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
