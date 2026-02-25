package cel2sql

import (
	"fmt"
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
		return newConversionError(errMsgInvalidTimestampOp, "timestamp operation requires at least one timestamp operand")
	}

	var sqlOp string
	switch fun {
	case operators.Add:
		sqlOp = "+"
	case operators.Subtract:
		sqlOp = "-"
	default:
		return newConversionError(errMsgInvalidTimestampOp, "unsupported timestamp operation")
	}

	return con.dialect.WriteTimestampArithmetic(&con.str, sqlOp,
		func() error { return con.visitMaybeNested(timestamp, timestampParen) },
		func() error { return con.visitMaybeNested(duration, durationParen) },
	)
}

// callDuration converts CEL duration expressions to PostgreSQL INTERVAL
func (con *converter) callDuration(_ *exprpb.Expr, args []*exprpb.Expr) error {
	if len(args) != 1 {
		return fmt.Errorf("%w: duration function requires exactly 1 argument, got %d", ErrInvalidArguments, len(args))
	}
	arg := args[0]
	var durationString string
	switch arg.ExprKind.(type) {
	case *exprpb.Expr_ConstExpr:
		switch arg.GetConstExpr().ConstantKind.(type) {
		case *exprpb.Constant_StringValue:
			durationString = arg.GetConstExpr().GetStringValue()
		default:
			return newConversionError(errMsgInvalidDuration, "unsupported constant type for duration")
		}
	default:
		return newConversionError(errMsgInvalidDuration, "unsupported expression type for duration")
	}
	d, err := time.ParseDuration(durationString)
	if err != nil {
		return err
	}
	var value int64
	var unit string
	switch d {
	case d.Round(time.Hour):
		value = int64(d.Hours())
		unit = "HOUR"
	case d.Round(time.Minute):
		value = int64(d.Minutes())
		unit = "MINUTE"
	case d.Round(time.Second):
		value = int64(d.Seconds())
		unit = "SECOND"
	case d.Round(time.Millisecond):
		value = d.Milliseconds()
		unit = "MILLISECOND"
	default:
		value = d.Truncate(time.Microsecond).Microseconds()
		unit = "MICROSECOND"
	}
	con.dialect.WriteDuration(&con.str, value, unit)
	return nil
}

// callInterval creates INTERVAL expressions using the dialect
func (con *converter) callInterval(_ *exprpb.Expr, args []*exprpb.Expr) error {
	datePart := args[1]
	unit := datePart.GetIdentExpr().GetName()
	return con.dialect.WriteInterval(&con.str, func() error {
		return con.visit(args[0])
	}, unit)
}

// callExtractFromTimestamp handles timestamp field extraction (YEAR, MONTH, DAY, etc.)
func (con *converter) callExtractFromTimestamp(function string, target *exprpb.Expr, args []*exprpb.Expr) error {
	var part string
	switch function {
	case overloads.TimeGetFullYear:
		part = "YEAR"
	case overloads.TimeGetMonth:
		part = "MONTH"
	case overloads.TimeGetDate:
		part = "DAY"
	case overloads.TimeGetHours:
		part = "HOUR"
	case overloads.TimeGetMinutes:
		part = "MINUTE"
	case overloads.TimeGetSeconds:
		part = "SECOND"
	case overloads.TimeGetMilliseconds:
		part = "MILLISECONDS"
	case overloads.TimeGetDayOfYear:
		part = "DOY"
	case overloads.TimeGetDayOfMonth:
		part = "DAY"
	case overloads.TimeGetDayOfWeek:
		part = "DOW"
	}

	writeExpr := func() error {
		return con.visit(target)
	}

	var writeTZ func() error
	if isTimestampType(con.getType(target)) && len(args) == 1 {
		writeTZ = func() error {
			return con.visit(args[0])
		}
	}

	if err := con.dialect.WriteExtract(&con.str, part, writeExpr, writeTZ); err != nil {
		return err
	}

	// Apply CEL-specific adjustments (these are universal, not dialect-specific)
	switch function {
	case overloads.TimeGetMonth, overloads.TimeGetDayOfYear, overloads.TimeGetDayOfMonth:
		con.str.WriteString(" - 1")
	}
	return nil
}

// callTimestampFromString converts string literals to timestamps using the dialect
func (con *converter) callTimestampFromString(_ *exprpb.Expr, args []*exprpb.Expr) error {
	if len(args) == 1 {
		return con.dialect.WriteTimestampCast(&con.str, func() error {
			return con.visit(args[0])
		})
	} else if len(args) == 2 {
		// Handle timestamp(datetime, timezone) format
		// For most dialects: datetime AT TIME ZONE timezone
		err := con.visit(args[0])
		if err != nil {
			return err
		}
		con.str.WriteString(" AT TIME ZONE ")
		err = con.visit(args[1])
		if err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("%w: timestamp function expects 1 or 2 arguments, got %d", ErrInvalidArguments, len(args))
}
