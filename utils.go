package cel2sql

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/cel-go/common/operators"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

// Type checking utilities

// isMapType checks if a type is a map type
func isMapType(typ *exprpb.Type) bool {
	_, ok := typ.TypeKind.(*exprpb.Type_MapType_)
	return ok
}

// isListType checks if a type is a list/array type
func isListType(typ *exprpb.Type) bool {
	_, ok := typ.TypeKind.(*exprpb.Type_ListType_)
	return ok
}

// Expression type checking utilities

// isNullLiteral checks if an expression is a NULL literal
func isNullLiteral(node *exprpb.Expr) bool {
	_, isConst := node.ExprKind.(*exprpb.Expr_ConstExpr)
	if !isConst {
		return false
	}
	_, isNull := node.GetConstExpr().ConstantKind.(*exprpb.Constant_NullValue)
	return isNull
}

// isBoolLiteral checks if an expression is a boolean literal
func isBoolLiteral(node *exprpb.Expr) bool {
	_, isConst := node.ExprKind.(*exprpb.Expr_ConstExpr)
	if !isConst {
		return false
	}
	_, isBool := node.GetConstExpr().ConstantKind.(*exprpb.Constant_BoolValue)
	return isBool
}

// isStringLiteral checks if an expression is a string literal
func isStringLiteral(node *exprpb.Expr) bool {
	_, isConst := node.ExprKind.(*exprpb.Expr_ConstExpr)
	if !isConst {
		return false
	}
	_, isString := node.GetConstExpr().ConstantKind.(*exprpb.Constant_StringValue)
	return isString
}

// isFieldAccessExpression checks if an expression is a field access (like trigram.cell[0].value)
func isFieldAccessExpression(expr *exprpb.Expr) bool {
	switch expr.GetExprKind().(type) {
	case *exprpb.Expr_SelectExpr:
		return true
	case *exprpb.Expr_CallExpr:
		// Check if it's an array index access
		call := expr.GetCallExpr()
		if call.GetFunction() == operators.Index {
			return true
		}
	}
	return false
}

// Field name validation and extraction

const (
	// maxPostgreSQLIdentifierLength is the maximum length for PostgreSQL identifiers
	// PostgreSQL's NAMEDATALEN is 64 bytes (including null terminator), so max usable length is 63
	maxPostgreSQLIdentifierLength = 63
)

var (
	// fieldNameRegexp validates PostgreSQL identifier format
	fieldNameRegexp = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

	// reservedSQLKeywords contains SQL keywords that should not be used as unquoted identifiers
	// This is a subset of PostgreSQL reserved keywords that could cause issues
	reservedSQLKeywords = map[string]bool{
		"all": true, "analyse": true, "analyze": true, "and": true, "any": true,
		"array": true, "as": true, "asc": true, "asymmetric": true, "both": true,
		"case": true, "cast": true, "check": true, "collate": true, "column": true,
		"constraint": true, "create": true, "cross": true, "current_catalog": true,
		"current_date": true, "current_role": true, "current_time": true,
		"current_timestamp": true, "current_user": true, "default": true,
		"deferrable": true, "desc": true, "distinct": true, "do": true, "else": true,
		"end": true, "except": true, "false": true, "fetch": true, "for": true,
		"foreign": true, "from": true, "grant": true, "group": true, "having": true,
		"in": true, "initially": true, "inner": true, "intersect": true, "into": true,
		"is": true, "join": true, "leading": true, "left": true, "like": true,
		"limit": true, "localtime": true, "localtimestamp": true, "natural": true,
		"not": true, "null": true, "offset": true, "on": true, "only": true,
		"or": true, "order": true, "outer": true, "overlaps": true, "placing": true,
		"primary": true, "references": true, "returning": true, "right": true,
		"select": true, "session_user": true, "similar": true, "some": true,
		"symmetric": true, "table": true, "then": true, "to": true, "trailing": true,
		"true": true, "union": true, "unique": true, "user": true, "using": true,
		"variadic": true, "when": true, "where": true, "window": true, "with": true,
		// Additional keywords that commonly cause issues
		"alter": true, "delete": true, "drop": true, "insert": true, "update": true,
	}
)

// validateFieldName validates that a field name follows PostgreSQL naming conventions
// and is safe to use in SQL queries without quoting
func validateFieldName(name string) error {
	// Check length (PostgreSQL max identifier length is 63 characters)
	if len(name) > maxPostgreSQLIdentifierLength {
		return fmt.Errorf("field name \"%s\" exceeds PostgreSQL maximum identifier length of %d characters", name, maxPostgreSQLIdentifierLength)
	}

	// Check if empty
	if len(name) == 0 {
		return errors.New("field name cannot be empty")
	}

	// Check format (must start with letter or underscore, contain only alphanumeric and underscore)
	if !fieldNameRegexp.MatchString(name) {
		return fmt.Errorf("invalid field name \"%s\": must start with a letter or underscore and contain only alphanumeric characters and underscores", name)
	}

	// Check for reserved SQL keywords (case-insensitive)
	if reservedSQLKeywords[strings.ToLower(name)] {
		return fmt.Errorf("field name \"%s\" is a reserved SQL keyword and cannot be used without quoting", name)
	}

	return nil
}

// extractFieldName extracts a field name from a string literal expression
func extractFieldName(node *exprpb.Expr) (string, error) {
	if !isStringLiteral(node) {
		return "", fmt.Errorf("unsupported type: %v", node)
	}
	fieldName := node.GetConstExpr().GetStringValue()
	if err := validateFieldName(fieldName); err != nil {
		return "", err
	}
	return fieldName, nil
}

// Numeric comparison utilities

// isNumericComparison checks if an operator is a numeric comparison
func isNumericComparison(op string) bool {
	return op == operators.Greater || op == operators.GreaterEquals ||
		op == operators.Less || op == operators.LessEquals ||
		op == operators.Equals || op == operators.NotEquals
}

// isNumericType checks if a type represents a numeric value
func isNumericType(typ *exprpb.Type) bool {
	if typ == nil {
		return false
	}
	primitive := typ.GetPrimitive()
	return primitive == exprpb.Type_INT64 ||
		primitive == exprpb.Type_UINT64 ||
		primitive == exprpb.Type_DOUBLE
}

// String pattern utilities

// escapeLikePattern escapes special characters in a LIKE pattern
// PostgreSQL LIKE special characters are: % (matches any sequence), _ (matches any single char), \ (escape)
func escapeLikePattern(pattern string) string {
	// Escape backslashes first, then % and _
	escaped := strings.ReplaceAll(pattern, `\`, `\\`)
	escaped = strings.ReplaceAll(escaped, `%`, `\%`)
	escaped = strings.ReplaceAll(escaped, `_`, `\_`)
	// Escape single quotes by doubling them (for PostgreSQL string literals)
	escaped = strings.ReplaceAll(escaped, `'`, `''`)
	return escaped
}

// escapeJSONFieldName escapes single quotes in JSON field names for safe use in PostgreSQL JSON path operators
// In PostgreSQL, single quotes within string literals must be escaped by doubling them
func escapeJSONFieldName(fieldName string) string {
	return strings.ReplaceAll(fieldName, "'", "''")
}
