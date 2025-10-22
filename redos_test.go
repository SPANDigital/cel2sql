package cel2sql_test

import (
	"strings"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/spandigital/cel2sql/v2"
)

// TestReDoSProtection_NestedQuantifiers tests protection against catastrophic nested quantifiers
func TestReDoSProtection_NestedQuantifiers(t *testing.T) {
	env, err := cel.NewEnv(
		cel.Variable("text", cel.StringType),
	)
	require.NoError(t, err)

	// Test cases with dangerous nested quantifier patterns
	tests := []struct {
		name        string
		pattern     string
		description string
	}{
		{
			name:        "nested_star_quantifier",
			pattern:     `(a*)*`,
			description: "Catastrophic backtracking with nested * quantifiers",
		},
		{
			name:        "nested_plus_quantifier",
			pattern:     `(a+)+`,
			description: "Catastrophic backtracking with nested + quantifiers",
		},
		{
			name:        "complex_nested_pattern",
			pattern:     `(x+x+)+y`,
			description: "Complex nested pattern causing exponential time",
		},
		{
			name:        "nested_with_brace_quantifier",
			pattern:     `(a*){2,}`,
			description: "Nested pattern with brace quantifier",
		},
		{
			name:        "double_star",
			pattern:     `a**`,
			description: "Double star quantifier",
		},
		{
			name:        "double_plus",
			pattern:     `a++`,
			description: "Double plus quantifier",
		},
		{
			name:        "mixed_nested",
			pattern:     `(a*b*)+`,
			description: "Mixed nested quantifiers",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			celExpr := `text.matches(r"` + tt.pattern + `")`
			ast, issues := env.Compile(celExpr)
			require.NoError(t, issues.Err(), "CEL should compile")

			// Conversion should fail with ReDoS protection
			sql, err := cel2sql.Convert(ast)
			require.Error(t, err, "Should reject dangerous pattern: %s", tt.description)
			assert.Empty(t, sql)
			assert.Contains(t, err.Error(), "nested quantifiers", "Error should mention nested quantifiers")
		})
	}
}

// TestReDoSProtection_QuantifiedAlternation tests protection against quantified alternation patterns
func TestReDoSProtection_QuantifiedAlternation(t *testing.T) {
	env, err := cel.NewEnv(
		cel.Variable("text", cel.StringType),
	)
	require.NoError(t, err)

	tests := []struct {
		name        string
		pattern     string
		description string
	}{
		{
			name:        "quantified_alternation_star",
			pattern:     `(a|a)*b`,
			description: "Quantified alternation with overlapping branches",
		},
		{
			name:        "quantified_alternation_plus",
			pattern:     `(a|ab)+`,
			description: "Quantified alternation with prefix overlap",
		},
		{
			name:        "complex_alternation",
			pattern:     `(foo|foobar)*`,
			description: "Complex alternation pattern",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			celExpr := `text.matches(r"` + tt.pattern + `")`
			ast, issues := env.Compile(celExpr)
			require.NoError(t, issues.Err())

			sql, err := cel2sql.Convert(ast)
			require.Error(t, err, "Should reject dangerous pattern: %s", tt.description)
			assert.Empty(t, sql)
			assert.Contains(t, err.Error(), "alternation", "Error should mention alternation")
		})
	}
}

// TestReDoSProtection_PatternLength tests protection against extremely long patterns
func TestReDoSProtection_PatternLength(t *testing.T) {
	env, err := cel.NewEnv(
		cel.Variable("text", cel.StringType),
	)
	require.NoError(t, err)

	// Create a pattern longer than maxRegexPatternLength (500)
	longPattern := strings.Repeat("a", 501)
	celExpr := `text.matches(r"` + longPattern + `")`

	ast, issues := env.Compile(celExpr)
	require.NoError(t, issues.Err())

	sql, err := cel2sql.Convert(ast)
	require.Error(t, err, "Should reject excessively long pattern")
	assert.Empty(t, sql)
	assert.Contains(t, err.Error(), "exceeds maximum length", "Error should mention length limit")
}

// TestReDoSProtection_GroupLimit tests protection against excessive capture groups
func TestReDoSProtection_GroupLimit(t *testing.T) {
	env, err := cel.NewEnv(
		cel.Variable("text", cel.StringType),
	)
	require.NoError(t, err)

	// Create a pattern with more than maxRegexGroups (20) groups
	groups := strings.Repeat("(a)", 21)
	celExpr := `text.matches(r"` + groups + `")`

	ast, issues := env.Compile(celExpr)
	require.NoError(t, issues.Err())

	sql, err := cel2sql.Convert(ast)
	require.Error(t, err, "Should reject pattern with too many groups")
	assert.Empty(t, sql)
	assert.Contains(t, err.Error(), "capture groups", "Error should mention capture groups")
}

// TestReDoSProtection_NestingDepth tests protection against deeply nested patterns
func TestReDoSProtection_NestingDepth(t *testing.T) {
	env, err := cel.NewEnv(
		cel.Variable("text", cel.StringType),
	)
	require.NoError(t, err)

	// Create a pattern with more than maxRegexNestingDepth (10) levels of nesting
	pattern := strings.Repeat("(", 11) + "a" + strings.Repeat(")", 11)
	celExpr := `text.matches(r"` + pattern + `")`

	ast, issues := env.Compile(celExpr)
	require.NoError(t, issues.Err())

	sql, err := cel2sql.Convert(ast)
	require.Error(t, err, "Should reject deeply nested pattern")
	assert.Empty(t, sql)
	assert.Contains(t, err.Error(), "nesting depth", "Error should mention nesting depth")
}

// TestReDoSProtection_SafePatterns tests that safe patterns are still allowed
func TestReDoSProtection_SafePatterns(t *testing.T) {
	env, err := cel.NewEnv(
		cel.Variable("text", cel.StringType),
	)
	require.NoError(t, err)

	tests := []struct {
		name        string
		pattern     string
		description string
	}{
		{
			name:        "simple_literal",
			pattern:     `hello`,
			description: "Simple literal pattern",
		},
		{
			name:        "character_class",
			pattern:     `[a-z]+`,
			description: "Character class with quantifier",
		},
		{
			name:        "anchored_pattern",
			pattern:     `^test$`,
			description: "Anchored pattern",
		},
		{
			name:        "digit_pattern",
			pattern:     `\d{3}-\d{4}`,
			description: "Digit pattern with braces",
		},
		{
			name:        "word_boundary",
			pattern:     `\bword\b`,
			description: "Word boundary pattern",
		},
		{
			name:        "safe_alternation",
			pattern:     `cat|dog`,
			description: "Simple alternation without quantifier",
		},
		{
			name:        "safe_group",
			pattern:     `(abc)+`,
			description: "Simple grouped pattern with quantifier",
		},
		{
			name:        "email_pattern",
			pattern:     `^[a-z]+@[a-z]+\.[a-z]+$`,
			description: "Email validation pattern",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			celExpr := `text.matches(r"` + tt.pattern + `")`
			ast, issues := env.Compile(celExpr)
			require.NoError(t, issues.Err())

			sql, err := cel2sql.Convert(ast)
			require.NoError(t, err, "Safe pattern should be allowed: %s", tt.description)
			assert.NotEmpty(t, sql)
			assert.Contains(t, sql, "~", "Should contain regex operator")
		})
	}
}

// TestReDoSProtection_RealWorldAttackPatterns tests actual ReDoS attack patterns from CVE databases
func TestReDoSProtection_RealWorldAttackPatterns(t *testing.T) {
	env, err := cel.NewEnv(
		cel.Variable("text", cel.StringType),
	)
	require.NoError(t, err)

	// Real-world catastrophic backtracking patterns
	tests := []struct {
		name        string
		pattern     string
		cve         string
		description string
	}{
		{
			name:        "email_redos",
			pattern:     `(a+)+@example.com`,
			cve:         "Similar to CVE-2019-5021",
			description: "Email validation ReDoS pattern",
		},
		{
			name:        "json_redos",
			pattern:     `(.*,)*`,
			cve:         "Similar to various JSON parsers",
			description: "JSON parsing ReDoS pattern",
		},
		{
			name:        "url_redos",
			pattern:     `^(([a-z])+.)+[A-Z]([a-z])+$`,
			cve:         "Similar to URL validation ReDoS",
			description: "URL validation backtracking",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			celExpr := `text.matches(r"` + tt.pattern + `")`
			ast, issues := env.Compile(celExpr)
			require.NoError(t, issues.Err())

			sql, err := cel2sql.Convert(ast)
			require.Error(t, err, "Should block real-world attack pattern: %s (%s)", tt.description, tt.cve)
			assert.Empty(t, sql)
		})
	}
}

// TestReDoSProtection_EdgeCases tests edge cases in pattern validation
func TestReDoSProtection_EdgeCases(t *testing.T) {
	env, err := cel.NewEnv(
		cel.Variable("text", cel.StringType),
	)
	require.NoError(t, err)

	tests := []struct {
		name        string
		pattern     string
		shouldPass  bool
		description string
	}{
		{
			name:        "escaped_parentheses",
			pattern:     `\(test\)`,
			shouldPass:  true,
			description: "Escaped parentheses should not count as groups",
		},
		{
			name:        "exactly_max_groups",
			pattern:     strings.Repeat("(a)", 20), // Exactly maxRegexGroups
			shouldPass:  true,
			description: "Exactly at the limit should pass",
		},
		{
			name:        "exactly_max_depth",
			pattern:     strings.Repeat("(", 10) + "a" + strings.Repeat(")", 10), // Exactly maxRegexNestingDepth
			shouldPass:  true,
			description: "Exactly at nesting depth limit should pass",
		},
		{
			name:        "empty_pattern",
			pattern:     ``,
			shouldPass:  true,
			description: "Empty pattern should be allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			celExpr := `text.matches(r"` + tt.pattern + `")`
			ast, issues := env.Compile(celExpr)
			require.NoError(t, issues.Err())

			sql, err := cel2sql.Convert(ast)
			if tt.shouldPass {
				require.NoError(t, err, "Pattern should pass: %s", tt.description)
				assert.NotEmpty(t, sql)
			} else {
				require.Error(t, err, "Pattern should fail: %s", tt.description)
				assert.Empty(t, sql)
			}
		})
	}
}

// TestReDoSProtection_FunctionStyleMatches tests matches() function style as well
func TestReDoSProtection_FunctionStyleMatches(t *testing.T) {
	env, err := cel.NewEnv(
		cel.Variable("text", cel.StringType),
	)
	require.NoError(t, err)

	// Test with function-style matches(str, pattern)
	celExpr := `matches(text, "(a+)+")`
	ast, issues := env.Compile(celExpr)
	require.NoError(t, issues.Err())

	sql, err := cel2sql.Convert(ast)
	require.Error(t, err, "Function-style matches should also be protected")
	assert.Empty(t, sql)
	assert.Contains(t, err.Error(), "nested quantifiers")
}
