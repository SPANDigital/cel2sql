package cel2sql

import (
	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/common/overloads"
)

// standardSQLBinaryOperators maps CEL binary operators to PostgreSQL SQL operators
var standardSQLBinaryOperators = map[string]string{
	operators.LogicalAnd: "AND",
	operators.LogicalOr:  "OR",
	operators.Equals:     "=",
}

// standardSQLUnaryOperators maps CEL unary operators to PostgreSQL SQL operators
var standardSQLUnaryOperators = map[string]string{
	operators.LogicalNot: "NOT ",
}

// standardSQLFunctions maps CEL function names to PostgreSQL function names
var standardSQLFunctions = map[string]string{
	operators.Modulo:     "MOD",
	overloads.StartsWith: "STARTS_WITH",
	overloads.EndsWith:   "ENDS_WITH",
	overloads.Matches:    "REGEXP_CONTAINS",
}
