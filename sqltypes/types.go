package sqltypes

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

var (
	Date     = decls.NewAbstractType("DATE")
	Time     = decls.NewAbstractType("TIME")
	DateTime = decls.NewAbstractType("DATETIME")
	Interval = decls.NewAbstractType("INTERVAL")
	DatePart = decls.NewAbstractType("date_part")
)

func newConstantString(str string) *expr.Constant {
	return &expr.Constant{ConstantKind: &expr.Constant_StringValue{StringValue: str}}
}

var SQLTypeDeclarations = cel.Types(
	// Custom abstract types
	Date, Time, DateTime, Interval, DatePart,
)
