package regexsafe

import (
	"strings"
	"testing"
)

func TestValidate_SafePatterns(t *testing.T) {
	safe := []string{
		"a+",
		"^[0-9]+$",
		`\btest\b`,
		`\d{3}-\d{4}`,
		`[a-z]+@[a-z]+\.[a-z]+`,
		"(abc)+",
		"(?:abc)",
		strings.Repeat("a", MaxPatternLength), // exactly at the length limit
	}
	for _, p := range safe {
		if err := Validate(p); err != nil {
			t.Errorf("Validate(%q) = %v, want nil", p, err)
		}
	}
}

func TestValidate_RejectsWithDescriptiveMessage(t *testing.T) {
	cases := []struct {
		name    string
		pattern string
		wantSub string
	}{
		{"lookahead", "(?=foo)", "lookahead assertions"},
		{"neg_lookahead", "(?!foo)", "lookahead assertions"},
		{"lookbehind", "(?<=foo)bar", "lookbehind assertions"},
		{"neg_lookbehind", "(?<!foo)bar", "lookbehind assertions"},
		{"named_group", "(?P<name>x)", "named capture groups"},
		{"inline_flag_m", "(?m)^x", "inline flags other than"},
		{"inline_flag_s", "(?s).x", "inline flags other than"},
		{"double_star", "a**", "nested quantifiers"},
		{"double_plus", "a++", "nested quantifiers"},
		{"nested_quantified_group", "(a+)+", "nested quantifiers"},
		{"too_long", strings.Repeat("a", MaxPatternLength+1), "pattern length"},
		{"too_many_groups", strings.Repeat("()", MaxGroups+1), "capture groups"},
		{"quantified_alternation", "(a|a)*b", "quantified alternation"},
		{"too_deep", strings.Repeat("(", MaxNestingDepth+1) + "a" + strings.Repeat(")", MaxNestingDepth+1), "nesting depth"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := Validate(tc.pattern)
			if err == nil {
				t.Fatalf("Validate(%q) = nil, want error containing %q", tc.pattern, tc.wantSub)
			}
			if !strings.Contains(err.Error(), tc.wantSub) {
				t.Errorf("Validate(%q) error = %q, want substring %q", tc.pattern, err.Error(), tc.wantSub)
			}
		})
	}
}

// TestValidate_InvalidPatternCaughtByCompile covers a malformed pattern that the
// heuristic checks do not flag but RE2 cannot compile.
func TestValidate_InvalidPatternCaughtByCompile(t *testing.T) {
	err := Validate("[unterminated")
	if err == nil || !strings.Contains(err.Error(), "invalid regex pattern") {
		t.Errorf("Validate(unterminated class) = %v, want 'invalid regex pattern'", err)
	}
}
