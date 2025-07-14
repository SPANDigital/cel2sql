package cel2sql

import (
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

var fieldNameRegexp = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]{0,127}$`)

// validateFieldName validates that a field name follows PostgreSQL naming conventions
func validateFieldName(name string) error {
	if !fieldNameRegexp.MatchString(name) {
		return fmt.Errorf("invalid field name \"%s\"", name)
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

// Byte conversion utilities

// bytesToOctets converts byte sequences to a string using a three digit octal encoded value
// per byte.
func bytesToOctets(byteVal []byte) string {
	var b strings.Builder
	for _, c := range byteVal {
		_, _ = fmt.Fprintf(&b, "\\%03o", c)
	}
	return b.String()
}
