package cel2sql_test

import (
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/spandigital/cel2sql/v3"
)

// TestRegexConversion_CaseInsensitive tests that (?i) flag generates ~* operator
func TestRegexConversion_CaseInsensitive(t *testing.T) {
	env, err := cel.NewEnv(
		cel.Variable("email", cel.StringType),
	)
	require.NoError(t, err)

	tests := []struct {
		name         string
		celExpr      string
		wantOperator string
		wantPattern  string
	}{
		{
			name:         "case_insensitive_simple",
			celExpr:      `email.matches(r"(?i)test@example\.com")`,
			wantOperator: "~*",
			wantPattern:  "test@example\\.com",
		},
		{
			name:         "case_insensitive_complex",
			celExpr:      `email.matches(r"(?i)[a-z]+@[a-z]+\.[a-z]+")`,
			wantOperator: "~*",
			wantPattern:  "[a-z]+@[a-z]+\\.[a-z]+",
		},
		{
			name:         "case_sensitive_default",
			celExpr:      `email.matches(r"test@example\.com")`,
			wantOperator: "~",
			wantPattern:  "test@example\\.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.celExpr)
			require.NoError(t, issues.Err())

			sql, err := cel2sql.Convert(ast)
			require.NoError(t, err)

			// Check that the correct operator is used
			assert.Contains(t, sql, tt.wantOperator, "SQL should contain operator %s", tt.wantOperator)

			// Check that (?i) is stripped from the pattern
			assert.NotContains(t, sql, "(?i)", "Pattern should not contain (?i) flag")

			// Check that the pattern is present
			assert.Contains(t, sql, tt.wantPattern, "SQL should contain pattern %s", tt.wantPattern)
		})
	}
}

// TestRegexConversion_NonCapturingGroups tests that (?:...) is converted to (...)
func TestRegexConversion_NonCapturingGroups(t *testing.T) {
	env, err := cel.NewEnv(
		cel.Variable("text", cel.StringType),
	)
	require.NoError(t, err)

	tests := []struct {
		name        string
		celExpr     string
		wantPattern string
	}{
		{
			name:        "simple_non_capturing_group",
			celExpr:     `text.matches(r"(?:abc)")`,
			wantPattern: "(abc)",
		},
		{
			name:        "non_capturing_with_alternation",
			celExpr:     `text.matches(r"(?:cat|dog)")`,
			wantPattern: "(cat|dog)",
		},
		{
			name:        "nested_non_capturing",
			celExpr:     `text.matches(r"(?:(?:a|b)c)")`,
			wantPattern: "((a|b)c)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.celExpr)
			require.NoError(t, issues.Err())

			sql, err := cel2sql.Convert(ast)
			require.NoError(t, err)

			// Check that (?:) is converted to (
			assert.NotContains(t, sql, "(?:", "Pattern should not contain (?:")
			assert.Contains(t, sql, tt.wantPattern, "SQL should contain converted pattern %s", tt.wantPattern)
		})
	}
}

// TestRegexConversion_UnsupportedFeatures tests that unsupported features return errors
func TestRegexConversion_UnsupportedFeatures(t *testing.T) {
	env, err := cel.NewEnv(
		cel.Variable("text", cel.StringType),
	)
	require.NoError(t, err)

	tests := []struct {
		name          string
		pattern       string
		wantErrorText string
	}{
		{
			name:          "lookahead_positive",
			pattern:       `test(?=ing)`,
			wantErrorText: "lookahead assertions",
		},
		{
			name:          "lookahead_negative",
			pattern:       `test(?!ing)`,
			wantErrorText: "lookahead assertions",
		},
		{
			name:          "lookbehind_positive",
			pattern:       `(?<=pre)test`,
			wantErrorText: "lookbehind assertions",
		},
		{
			name:          "lookbehind_negative",
			pattern:       `(?<!pre)test`,
			wantErrorText: "lookbehind assertions",
		},
		{
			name:          "named_group",
			pattern:       `(?P<name>[a-z]+)`,
			wantErrorText: "named capture groups",
		},
		{
			name:          "multiline_flag",
			pattern:       `(?m)^test$`,
			wantErrorText: "inline flags other than",
		},
		{
			name:          "dotall_flag",
			pattern:       `(?s)test.+`,
			wantErrorText: "inline flags other than",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			celExpr := `text.matches(r"` + tt.pattern + `")`
			ast, issues := env.Compile(celExpr)
			require.NoError(t, issues.Err(), "CEL should compile")

			// Conversion should fail with clear error
			sql, err := cel2sql.Convert(ast)
			require.Error(t, err, "Should reject unsupported pattern: %s", tt.pattern)
			assert.Empty(t, sql)
			assert.Contains(t, err.Error(), tt.wantErrorText, "Error should mention %s", tt.wantErrorText)
		})
	}
}

// TestRegexConversion_CharacterClasses tests that character classes are converted correctly
func TestRegexConversion_CharacterClasses(t *testing.T) {
	env, err := cel.NewEnv(
		cel.Variable("text", cel.StringType),
	)
	require.NoError(t, err)

	tests := []struct {
		name         string
		celPattern   string
		wantContains string
	}{
		{
			name:         "digit_class",
			celPattern:   `\d+`,
			wantContains: "[[:digit:]]",
		},
		{
			name:         "word_class",
			celPattern:   `\w+`,
			wantContains: "[[:alnum:]_]",
		},
		{
			name:         "whitespace_class",
			celPattern:   `\s+`,
			wantContains: "[[:space:]]",
		},
		{
			name:         "non_digit_class",
			celPattern:   `\D+`,
			wantContains: "[^[:digit:]]",
		},
		{
			name:         "non_word_class",
			celPattern:   `\W+`,
			wantContains: "[^[:alnum:]_]",
		},
		{
			name:         "non_whitespace_class",
			celPattern:   `\S+`,
			wantContains: "[^[:space:]]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			celExpr := `text.matches(r"` + tt.celPattern + `")`
			ast, issues := env.Compile(celExpr)
			require.NoError(t, issues.Err())

			sql, err := cel2sql.Convert(ast)
			require.NoError(t, err)
			assert.Contains(t, sql, tt.wantContains, "SQL should contain converted character class")
		})
	}
}

// TestRegexConversion_CombinedFeatures tests combining multiple conversion features
func TestRegexConversion_CombinedFeatures(t *testing.T) {
	env, err := cel.NewEnv(
		cel.Variable("email", cel.StringType),
	)
	require.NoError(t, err)

	tests := []struct {
		name         string
		celExpr      string
		wantOperator string
		wantContains []string
		notContains  []string
	}{
		{
			name:         "case_insensitive_with_non_capturing",
			celExpr:      `email.matches(r"(?i)(?:admin|user)@\w+\.com")`,
			wantOperator: "~*",
			wantContains: []string{"(admin|user)", "[[:alnum:]_]+"},
			notContains:  []string{"(?i)", "(?:"},
		},
		{
			name:         "case_insensitive_with_char_classes",
			celExpr:      `email.matches(r"(?i)\w+@\w+\.\w+")`,
			wantOperator: "~*",
			wantContains: []string{"[[:alnum:]_]+"},
			notContains:  []string{"(?i)", `\w`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.celExpr)
			require.NoError(t, issues.Err())

			sql, err := cel2sql.Convert(ast)
			require.NoError(t, err)

			assert.Contains(t, sql, tt.wantOperator, "SQL should use operator %s", tt.wantOperator)

			for _, want := range tt.wantContains {
				assert.Contains(t, sql, want, "SQL should contain %s", want)
			}

			for _, notWant := range tt.notContains {
				assert.NotContains(t, sql, notWant, "SQL should not contain %s", notWant)
			}
		})
	}
}

// TestRegexConversion_BackwardCompatibility tests that existing patterns still work
func TestRegexConversion_BackwardCompatibility(t *testing.T) {
	env, err := cel.NewEnv(
		cel.Variable("email", cel.StringType),
	)
	require.NoError(t, err)

	tests := []struct {
		name    string
		celExpr string
	}{
		{
			name:    "simple_literal",
			celExpr: `email.matches(r"test@example.com")`,
		},
		{
			name:    "with_anchors",
			celExpr: `email.matches(r"^[a-z]+@[a-z]+\.[a-z]+$")`,
		},
		{
			name:    "with_quantifiers",
			celExpr: `email.matches(r"[a-z]{3,10}@[a-z]+\.com")`,
		},
		{
			name:    "with_alternation",
			celExpr: `email.matches(r"(admin|user)@example\.com")`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.celExpr)
			require.NoError(t, issues.Err())

			sql, err := cel2sql.Convert(ast)
			require.NoError(t, err)
			assert.NotEmpty(t, sql, "Should generate valid SQL")
		})
	}
}
