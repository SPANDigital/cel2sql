// Package sqltypes provides custom SQL type definitions for CEL (Date, Time, DateTime).
package sqltypes

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
)

var (
	// Date represents a SQL DATE type for CEL.
	Date = decls.NewAbstractType("DATE")
	// Time represents a SQL TIME type for CEL.
	Time = decls.NewAbstractType("TIME")
	// DateTime represents a SQL DATETIME type for CEL.
	DateTime = decls.NewAbstractType("DATETIME")
	// Interval represents a SQL INTERVAL type for CEL.
	Interval = decls.NewAbstractType("INTERVAL")
	// DatePart represents a SQL date_part function type for CEL.
	DatePart = decls.NewAbstractType("date_part")
)

// SQLTypeDeclarations provides CEL type declarations for custom SQL types.
var SQLTypeDeclarations = cel.Types(
	// Custom abstract types
	Date, Time, DateTime, Interval, DatePart,
)
