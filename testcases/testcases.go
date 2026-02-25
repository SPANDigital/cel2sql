// Package testcases defines shared test case types and helpers for multi-dialect testing.
package testcases

import "github.com/spandigital/cel2sql/v3/dialect"

// Category classifies a test case for organization and selective running.
type Category string

// Test case categories.
const (
	CategoryBasic         Category = "basic"
	CategoryOperator      Category = "operator"
	CategoryString        Category = "string"
	CategoryRegex         Category = "regex"
	CategoryJSON          Category = "json"
	CategoryArray         Category = "array"
	CategoryComprehension Category = "comprehension"
	CategoryTimestamp     Category = "timestamp"
	CategoryParameterized Category = "parameterized"
	CategoryCast          Category = "cast"
	CategoryFieldAccess   Category = "field_access"
)

// ConvertTestCase defines a single CEL-to-SQL conversion test case
// with expected output per dialect.
type ConvertTestCase struct {
	// Name is the test case name (used for t.Run).
	Name string

	// CELExpr is the CEL expression source to compile and convert.
	CELExpr string

	// Category classifies the test case.
	Category Category

	// EnvSetup identifies which CEL environment setup to use.
	// Empty string means "default" (basic types, no schema).
	EnvSetup string

	// WantSQL maps dialect name to expected SQL output.
	// If a dialect is absent, the test is skipped for that dialect.
	WantSQL map[dialect.Name]string

	// WantErr maps dialect name to whether an error is expected.
	// If a dialect is absent, no error is expected.
	WantErr map[dialect.Name]bool

	// SkipDialect maps dialect name to a skip reason.
	// If a dialect is present, the test is skipped with the given message.
	SkipDialect map[dialect.Name]string
}

// ForDialect returns the expected SQL for a given dialect, and whether the
// test case has an expectation for that dialect.
func (tc *ConvertTestCase) ForDialect(d dialect.Name) (sql string, hasExpected bool) {
	sql, hasExpected = tc.WantSQL[d]
	return
}

// ShouldError returns whether an error is expected for the given dialect.
func (tc *ConvertTestCase) ShouldError(d dialect.Name) bool {
	return tc.WantErr[d]
}

// ShouldSkip returns the skip reason for a dialect, or empty string if not skipped.
func (tc *ConvertTestCase) ShouldSkip(d dialect.Name) string {
	if tc.SkipDialect == nil {
		return ""
	}
	return tc.SkipDialect[d]
}

// ParameterizedTestCase defines a test case for parameterized SQL conversion.
type ParameterizedTestCase struct {
	// Name is the test case name.
	Name string

	// CELExpr is the CEL expression source.
	CELExpr string

	// Category classifies the test case.
	Category Category

	// EnvSetup identifies which CEL environment setup to use.
	EnvSetup string

	// WantSQL maps dialect name to expected parameterized SQL output.
	WantSQL map[dialect.Name]string

	// WantParams maps dialect name to expected parameter values.
	WantParams map[dialect.Name][]any

	// WantErr maps dialect name to whether an error is expected.
	WantErr map[dialect.Name]bool

	// SkipDialect maps dialect name to a skip reason.
	SkipDialect map[dialect.Name]string
}
