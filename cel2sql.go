// Package cel2sql converts CEL (Common Expression Language) expressions to PostgreSQL SQL conditions.
package cel2sql

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/common/overloads"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

// Implementations based on `google/cel-go`'s unparser
// https://github.com/google/cel-go/blob/master/parser/unparser.go

// Convert converts a CEL AST to a PostgreSQL SQL WHERE clause condition.
func Convert(ast *cel.Ast) (string, error) {
	checkedExpr, err := cel.AstToCheckedExpr(ast)
	if err != nil {
		return "", err
	}
	un := &converter{
		typeMap: checkedExpr.TypeMap,
	}
	if err := un.visit(checkedExpr.Expr); err != nil {
		return "", err
	}
	return un.str.String(), nil
}

type converter struct {
	str     strings.Builder
	typeMap map[int64]*exprpb.Type
}

func (con *converter) visit(expr *exprpb.Expr) error {
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
	return fmt.Errorf("unsupported expr: %v", expr)
}

func (con *converter) visitCall(expr *exprpb.Expr) error {
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
	if err := con.visitMaybeNested(lhs, lhsParen); err != nil {
		return err
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
		return fmt.Errorf("cannot unmangle operator: %s", fun)
	}
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
	con.str.WriteString("IF(")
	if err := con.visit(args[0]); err != nil {
		return err
	}
	con.str.WriteString(", ")
	if err := con.visit(args[1]); err != nil {
		return err
	}
	con.str.WriteString(", ")
	if err := con.visit(args[2]); err != nil {
		return nil
	}
	con.str.WriteString(")")
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

func (con *converter) callCasting(function string, _ *exprpb.Expr, args []*exprpb.Expr) error {
	arg := args[0]
	if function == overloads.TypeConvertInt && isTimestampType(con.getType(arg)) {
		con.str.WriteString("UNIX_SECONDS(")
		if err := con.visit(arg); err != nil {
			return err
		}
		con.str.WriteString(")")
		return nil
	}
	con.str.WriteString("CAST(")
	if err := con.visit(arg); err != nil {
		return err
	}
	con.str.WriteString(" AS ")
	switch function {
	case overloads.TypeConvertBool:
		con.str.WriteString("BOOL")
	case overloads.TypeConvertBytes:
		con.str.WriteString("BYTES")
	case overloads.TypeConvertDouble:
		con.str.WriteString("FLOAT64")
	case overloads.TypeConvertInt:
		con.str.WriteString("INT64")
	case overloads.TypeConvertString:
		con.str.WriteString("STRING")
	case overloads.TypeConvertUint:
		con.str.WriteString("INT64")
	}
	con.str.WriteString(")")
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
				return fmt.Errorf("unsupported type: %v", argType)
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
		con.str.WriteString(strconv.FormatInt(constExpr.GetInt64Value()+1, 10))
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
		return fmt.Errorf("cannot unmangle operator: %s", fun)
	}
	con.str.WriteString(operator)
	nested := isComplexOperator(args[0])
	return con.visitMaybeNested(args[0], nested)
}

func (con *converter) visitComprehension(expr *exprpb.Expr) error {
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
		return fmt.Errorf("unsupported comprehension type: %v", info.Type)
	}
}

// Comprehension visit functions - Phase 1 placeholder implementations

func (con *converter) visitAllComprehension(expr *exprpb.Expr, info *ComprehensionInfo) error {
	// Generate SQL for ALL comprehension: all elements must satisfy the predicate
	// Pattern: NOT EXISTS (SELECT 1 FROM UNNEST(array) AS item WHERE NOT predicate)
	// For JSON arrays: NOT EXISTS (SELECT 1 FROM jsonb_array_elements(json_field) AS item WHERE NOT predicate)

	comprehension := expr.GetComprehensionExpr()
	if comprehension == nil {
		return errors.New("expression is not a comprehension")
	}

	iterRange := comprehension.GetIterRange()
	isJSONArray := con.isJSONArrayField(iterRange)

	con.str.WriteString("NOT EXISTS (SELECT 1 FROM ")

	if isJSONArray {
		jsonFunc := con.getJSONArrayFunction(iterRange)
		con.str.WriteString(jsonFunc)
		con.str.WriteString("(")
		if err := con.visit(iterRange); err != nil {
			return fmt.Errorf("failed to visit iter range in ALL comprehension: %w", err)
		}
		con.str.WriteString(")")
	} else {
		con.str.WriteString("UNNEST(")
		if err := con.visit(iterRange); err != nil {
			return fmt.Errorf("failed to visit iter range in ALL comprehension: %w", err)
		}
		con.str.WriteString(")")
	}

	con.str.WriteString(" AS ")
	con.str.WriteString(info.IterVar)

	con.str.WriteString(" WHERE ")

	// Add null checks for JSON arrays
	if isJSONArray {
		if err := con.visit(iterRange); err != nil {
			return fmt.Errorf("failed to visit iter range for null check: %w", err)
		}
		con.str.WriteString(" IS NOT NULL AND ")
		typeofFunc := con.getJSONTypeofFunction(iterRange)
		con.str.WriteString(typeofFunc)
		con.str.WriteString("(")
		if err := con.visit(iterRange); err != nil {
			return fmt.Errorf("failed to visit iter range for type check: %w", err)
		}
		con.str.WriteString(") = 'array'")

		if info.Predicate != nil {
			con.str.WriteString(" AND ")
		}
	}

	if info.Predicate != nil {
		con.str.WriteString("NOT (")
		if err := con.visit(info.Predicate); err != nil {
			return fmt.Errorf("failed to visit predicate in ALL comprehension: %w", err)
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
		return errors.New("expression is not a comprehension")
	}

	iterRange := comprehension.GetIterRange()
	isJSONArray := con.isJSONArrayField(iterRange)

	con.str.WriteString("EXISTS (SELECT 1 FROM ")

	if isJSONArray {
		jsonFunc := con.getJSONArrayFunction(iterRange)
		con.str.WriteString(jsonFunc)
		con.str.WriteString("(")
		if err := con.visit(iterRange); err != nil {
			return fmt.Errorf("failed to visit iter range in EXISTS comprehension: %w", err)
		}
		con.str.WriteString(")")
	} else {
		con.str.WriteString("UNNEST(")
		if err := con.visit(iterRange); err != nil {
			return fmt.Errorf("failed to visit iter range in EXISTS comprehension: %w", err)
		}
		con.str.WriteString(")")
	}

	con.str.WriteString(" AS ")
	con.str.WriteString(info.IterVar)

	con.str.WriteString(" WHERE ")

	// Add null checks for JSON arrays
	if isJSONArray {
		if err := con.visit(iterRange); err != nil {
			return fmt.Errorf("failed to visit iter range for null check: %w", err)
		}
		con.str.WriteString(" IS NOT NULL AND ")
		typeofFunc := con.getJSONTypeofFunction(iterRange)
		con.str.WriteString(typeofFunc)
		con.str.WriteString("(")
		if err := con.visit(iterRange); err != nil {
			return fmt.Errorf("failed to visit iter range for type check: %w", err)
		}
		con.str.WriteString(") = 'array'")

		if info.Predicate != nil {
			con.str.WriteString(" AND ")
		}
	}

	if info.Predicate != nil {
		if err := con.visit(info.Predicate); err != nil {
			return fmt.Errorf("failed to visit predicate in EXISTS comprehension: %w", err)
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
		return errors.New("expression is not a comprehension")
	}

	iterRange := comprehension.GetIterRange()
	isJSONArray := con.isJSONArrayField(iterRange)

	con.str.WriteString("(SELECT COUNT(*) FROM ")

	if isJSONArray {
		jsonFunc := con.getJSONArrayFunction(iterRange)
		con.str.WriteString(jsonFunc)
		con.str.WriteString("(")
		if err := con.visit(iterRange); err != nil {
			return fmt.Errorf("failed to visit iter range in EXISTS_ONE comprehension: %w", err)
		}
		con.str.WriteString(")")
	} else {
		con.str.WriteString("UNNEST(")
		if err := con.visit(iterRange); err != nil {
			return fmt.Errorf("failed to visit iter range in EXISTS_ONE comprehension: %w", err)
		}
		con.str.WriteString(")")
	}

	con.str.WriteString(" AS ")
	con.str.WriteString(info.IterVar)

	con.str.WriteString(" WHERE ")

	// Add null checks for JSON arrays
	if isJSONArray {
		if err := con.visit(iterRange); err != nil {
			return fmt.Errorf("failed to visit iter range for null check: %w", err)
		}
		con.str.WriteString(" IS NOT NULL AND ")
		typeofFunc := con.getJSONTypeofFunction(iterRange)
		con.str.WriteString(typeofFunc)
		con.str.WriteString("(")
		if err := con.visit(iterRange); err != nil {
			return fmt.Errorf("failed to visit iter range for type check: %w", err)
		}
		con.str.WriteString(") = 'array'")

		if info.Predicate != nil {
			con.str.WriteString(" AND ")
		}
	}

	if info.Predicate != nil {
		if err := con.visit(info.Predicate); err != nil {
			return fmt.Errorf("failed to visit predicate in EXISTS_ONE comprehension: %w", err)
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
		return errors.New("expression is not a comprehension")
	}

	iterRange := comprehension.GetIterRange()
	isJSONArray := con.isJSONArrayField(iterRange)

	con.str.WriteString("ARRAY(SELECT ")

	// Visit the transform expression
	if info.Transform != nil {
		if err := con.visit(info.Transform); err != nil {
			return fmt.Errorf("failed to visit transform in MAP comprehension: %w", err)
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
			return fmt.Errorf("failed to visit iter range in MAP comprehension: %w", err)
		}
		con.str.WriteString(")")
	} else {
		con.str.WriteString("UNNEST(")
		if err := con.visit(iterRange); err != nil {
			return fmt.Errorf("failed to visit iter range in MAP comprehension: %w", err)
		}
		con.str.WriteString(")")
	}

	con.str.WriteString(" AS ")
	con.str.WriteString(info.IterVar)

	// Add filter condition if present (for map with filter)
	if info.Filter != nil {
		con.str.WriteString(" WHERE ")
		if err := con.visit(info.Filter); err != nil {
			return fmt.Errorf("failed to visit filter in MAP comprehension: %w", err)
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
		return errors.New("expression is not a comprehension")
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
			return fmt.Errorf("failed to visit iter range in FILTER comprehension: %w", err)
		}
		con.str.WriteString(")")
	} else {
		con.str.WriteString("UNNEST(")
		if err := con.visit(iterRange); err != nil {
			return fmt.Errorf("failed to visit iter range in FILTER comprehension: %w", err)
		}
		con.str.WriteString(")")
	}

	con.str.WriteString(" AS ")
	con.str.WriteString(info.IterVar)

	if info.Predicate != nil {
		con.str.WriteString(" WHERE ")
		if err := con.visit(info.Predicate); err != nil {
			return fmt.Errorf("failed to visit predicate in FILTER comprehension: %w", err)
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
		return errors.New("expression is not a comprehension")
	}

	con.str.WriteString("ARRAY(SELECT ")

	// Visit the transform expression
	if info.Transform != nil {
		if err := con.visit(info.Transform); err != nil {
			return fmt.Errorf("failed to visit transform in TRANSFORM_LIST comprehension: %w", err)
		}
	} else {
		// If no transform, just return the variable itself
		con.str.WriteString(info.IterVar)
	}

	con.str.WriteString(" FROM UNNEST(")

	// Visit the iterable range (the array/list being comprehended over)
	if err := con.visit(comprehension.GetIterRange()); err != nil {
		return fmt.Errorf("failed to visit iter range in TRANSFORM_LIST comprehension: %w", err)
	}

	con.str.WriteString(") AS ")
	con.str.WriteString(info.IterVar)

	// Add filter condition if present
	if info.Filter != nil {
		con.str.WriteString(" WHERE ")
		if err := con.visit(info.Filter); err != nil {
			return fmt.Errorf("failed to visit filter in TRANSFORM_LIST comprehension: %w", err)
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
		if c.GetBoolValue() {
			con.str.WriteString("TRUE")
		} else {
			con.str.WriteString("FALSE")
		}
	case *exprpb.Constant_BytesValue:
		b := c.GetBytesValue()
		con.str.WriteString(`b"`)
		con.str.WriteString(bytesToOctets(b))
		con.str.WriteString(`"`)
	case *exprpb.Constant_DoubleValue:
		d := strconv.FormatFloat(c.GetDoubleValue(), 'g', -1, 64)
		con.str.WriteString(d)
	case *exprpb.Constant_Int64Value:
		i := strconv.FormatInt(c.GetInt64Value(), 10)
		con.str.WriteString(i)
	case *exprpb.Constant_NullValue:
		con.str.WriteString("NULL")
	case *exprpb.Constant_StringValue:
		// Use single quotes for PostgreSQL string literals
		str := c.GetStringValue()
		// Escape single quotes by doubling them
		escaped := strings.ReplaceAll(str, "'", "''")
		con.str.WriteString("'")
		con.str.WriteString(escaped)
		con.str.WriteString("'")
	case *exprpb.Constant_Uint64Value:
		ui := strconv.FormatUint(c.GetUint64Value(), 10)
		con.str.WriteString(ui)
	default:
		return fmt.Errorf("unimplemented : %v", expr)
	}
	return nil
}

func (con *converter) visitIdent(expr *exprpb.Expr) error {
	identName := expr.GetIdentExpr().GetName()

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
	// handle the case when the select expression was generated by the has() macro.
	if sel.GetTestOnly() {
		con.str.WriteString("has(")
	}

	// Check if we should use JSON path operators
	// We need to determine if the operand is a JSON/JSONB field
	useJSONPath := con.shouldUseJSONPath(sel.GetOperand(), sel.GetField())
	useJSONObjectAccess := con.isJSONObjectFieldAccess(expr)

	nested := !sel.GetTestOnly() && isBinaryOrTernaryOperator(sel.GetOperand())

	if useJSONObjectAccess && con.isNumericJSONField(sel.GetField()) {
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
		con.str.WriteString(sel.GetField())
		con.str.WriteString("'")
	case useJSONObjectAccess:
		// Use -> for JSON object field access in comprehensions
		fieldName := sel.GetField()
		con.str.WriteString("->>'")
		con.str.WriteString(fieldName)
		con.str.WriteString("'")
		if con.isNumericJSONField(fieldName) {
			// Close parentheses and add numeric cast
			con.str.WriteString(")::numeric")
		}
	default:
		// Regular field selection
		con.str.WriteString(".")
		con.str.WriteString(sel.GetField())
	}

	if sel.GetTestOnly() {
		con.str.WriteString(")")
	}
	return nil
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
	con.str.WriteString("STRUCT(")
	for i, entry := range entries {
		v := entry.GetValue()
		if err := con.visit(v); err != nil {
			return err
		}
		con.str.WriteString(" AS ")
		fieldName, err := extractFieldName(entry.GetMapKey())
		if err != nil {
			return err
		}
		con.str.WriteString(fieldName)
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
