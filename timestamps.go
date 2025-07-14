package cel2sql

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/common/overloads"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

// isTimestampRelatedType checks if a type is timestamp-related (DATE, TIME, DATETIME, TIMESTAMP)
func isTimestampRelatedType(typ *exprpb.Type) bool {
	abstractType := typ.GetAbstractType()
	if abstractType != nil {
		name := abstractType.GetName()
		return name == "DATE" || name == "TIME" || name == "DATETIME"
	}
	return typ.GetWellKnown() == exprpb.Type_TIMESTAMP
}

// isTimestampType checks if a type is specifically a TIMESTAMP
func isTimestampType(typ *exprpb.Type) bool {
	return typ.GetWellKnown() == exprpb.Type_TIMESTAMP
}

// isDurationRelatedType checks if a type is duration-related (INTERVAL, DURATION)
func isDurationRelatedType(typ *exprpb.Type) bool {
	abstractType := typ.GetAbstractType()
	if abstractType != nil {
		name := abstractType.GetName()
		return name == "INTERVAL"
	}
	return typ.GetWellKnown() == exprpb.Type_DURATION
}

// callTimestampOperation handles timestamp arithmetic (addition/subtraction with durations)
func (con *converter) callTimestampOperation(fun string, lhs *exprpb.Expr, rhs *exprpb.Expr) error {
	lhsParen := isComplexOperatorWithRespectTo(fun, lhs)
	rhsParen := isComplexOperatorWithRespectTo(fun, rhs)
	lhsType := con.getType(lhs)
	rhsType := con.getType(rhs)

	var timestamp, duration *exprpb.Expr
	var timestampParen, durationParen bool
	switch {
	case isTimestampRelatedType(lhsType):
		timestamp, duration = lhs, rhs
		timestampParen, durationParen = lhsParen, rhsParen
	case isTimestampRelatedType(rhsType):
		timestamp, duration = rhs, lhs
		timestampParen, durationParen = rhsParen, lhsParen
	default:
		panic("lhs or rhs must be timestamp related type")
	}

	// PostgreSQL uses simple + and - operators for date arithmetic
	var sqlOp string
	switch fun {
	case operators.Add:
		sqlOp = "+"
	case operators.Subtract:
		sqlOp = "-"
	default:
		return fmt.Errorf("unsupported operation (%s)", fun)
	}

	if err := con.visitMaybeNested(timestamp, timestampParen); err != nil {
		return err
	}
	con.str.WriteString(" ")
	con.str.WriteString(sqlOp)
	con.str.WriteString(" ")
	if err := con.visitMaybeNested(duration, durationParen); err != nil {
		return err
	}
	return nil
}

// callDuration converts CEL duration expressions to PostgreSQL INTERVAL
func (con *converter) callDuration(_ *exprpb.Expr, args []*exprpb.Expr) error {
	if len(args) != 1 {
		return errors.New("arguments must be single")
	}
	arg := args[0]
	var durationString string
	switch arg.ExprKind.(type) {
	case *exprpb.Expr_ConstExpr:
		switch arg.GetConstExpr().ConstantKind.(type) {
		case *exprpb.Constant_StringValue:
			durationString = arg.GetConstExpr().GetStringValue()
		default:
			return fmt.Errorf("unsupported constant kind %t", arg.GetConstExpr().ConstantKind)
		}
	default:
		return fmt.Errorf("unsupported kind %t", arg.ExprKind)
	}
	d, err := time.ParseDuration(durationString)
	if err != nil {
		return err
	}
	con.str.WriteString("INTERVAL ")
	switch d {
	case d.Round(time.Hour):
		con.str.WriteString(strconv.FormatFloat(d.Hours(), 'f', 0, 64))
		con.str.WriteString(" HOUR")
	case d.Round(time.Minute):
		con.str.WriteString(strconv.FormatFloat(d.Minutes(), 'f', 0, 64))
		con.str.WriteString(" MINUTE")
	case d.Round(time.Second):
		con.str.WriteString(strconv.FormatFloat(d.Seconds(), 'f', 0, 64))
		con.str.WriteString(" SECOND")
	case d.Round(time.Millisecond):
		con.str.WriteString(strconv.FormatInt(d.Milliseconds(), 10))
		con.str.WriteString(" MILLISECOND")
	default:
		con.str.WriteString(strconv.FormatInt(d.Truncate(time.Microsecond).Microseconds(), 10))
		con.str.WriteString(" MICROSECOND")
	}
	return nil
}

// callInterval creates PostgreSQL INTERVAL expressions
func (con *converter) callInterval(_ *exprpb.Expr, args []*exprpb.Expr) error {
	con.str.WriteString("INTERVAL ")
	if err := con.visit(args[0]); err != nil {
		return err
	}
	con.str.WriteString(" ")
	datePart := args[1]
	con.str.WriteString(datePart.GetIdentExpr().GetName())
	return nil
}

// callExtractFromTimestamp handles timestamp field extraction (YEAR, MONTH, DAY, etc.)
func (con *converter) callExtractFromTimestamp(function string, target *exprpb.Expr, args []*exprpb.Expr) error {
	con.str.WriteString("EXTRACT(")
	switch function {
	case overloads.TimeGetFullYear:
		con.str.WriteString("YEAR")
	case overloads.TimeGetMonth:
		con.str.WriteString("MONTH")
	case overloads.TimeGetDate:
		con.str.WriteString("DAY")
	case overloads.TimeGetHours:
		con.str.WriteString("HOUR")
	case overloads.TimeGetMinutes:
		con.str.WriteString("MINUTE")
	case overloads.TimeGetSeconds:
		con.str.WriteString("SECOND")
	case overloads.TimeGetMilliseconds:
		con.str.WriteString("MILLISECOND")
	case overloads.TimeGetDayOfYear:
		con.str.WriteString("DAYOFYEAR")
	case overloads.TimeGetDayOfMonth:
		con.str.WriteString("DAY")
	case overloads.TimeGetDayOfWeek:
		con.str.WriteString("DAYOFWEEK")
	}
	con.str.WriteString(" FROM ")
	if err := con.visit(target); err != nil {
		return err
	}
	if isTimestampType(con.getType(target)) && len(args) == 1 {
		con.str.WriteString(" AT ")
		if err := con.visit(args[0]); err != nil {
			return err
		}
	}
	con.str.WriteString(")")
	if function == overloads.TimeGetMonth || function == overloads.TimeGetDayOfYear || function == overloads.TimeGetDayOfMonth || function == overloads.TimeGetDayOfWeek {
		con.str.WriteString(" - 1")
	}
	return nil
}

// callTimestampFromString converts string literals to PostgreSQL timestamps
func (con *converter) callTimestampFromString(_ *exprpb.Expr, args []*exprpb.Expr) error {
	if len(args) == 1 {
		// For PostgreSQL, we need to cast the string to a timestamp
		con.str.WriteString("CAST(")
		err := con.visit(args[0])
		if err != nil {
			return err
		}
		con.str.WriteString(" AS TIMESTAMP WITH TIME ZONE)")
		return nil
	} else if len(args) == 2 {
		// Handle timestamp(datetime, timezone) format
		con.str.WriteString("TIMESTAMP(")
		err := con.visit(args[0])
		if err != nil {
			return err
		}
		con.str.WriteString(", ")
		err = con.visit(args[1])
		if err != nil {
			return err
		}
		con.str.WriteString(")")
		return nil
	}

	return fmt.Errorf("timestamp function expects 1 or 2 arguments, got %d", len(args))
}
