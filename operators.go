package cel2sql

import (
	"cel.dev/cel-go/common/operators"
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
	operators.Modulo: "MOD",
	// Note: overloads.StartsWith and overloads.EndsWith are handled specially in visitCallFunc
	// Note: overloads.Matches is handled specially in visitCallFunc with RE2 to POSIX conversion
}
