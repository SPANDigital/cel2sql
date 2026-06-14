// Package regexsafe centralises the RE2 ReDoS / unsupported-feature validation
// shared by every SQL dialect's regex conversion.
//
// Each dialect previously reimplemented this validation in its own regex.go,
// and the copies had drifted (some compiled the pattern first, the
// nested-quantifier and nesting-depth loops were written differently, error
// wording diverged). A single implementation guarantees that a check tightened
// for one dialect protects all of them. Only the dialect-specific parts —
// case-insensitivity handling and the final character-class / non-capturing
// group transform — remain in each dialect's regex.go.
package regexsafe

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Pattern complexity limits to prevent ReDoS attacks (CWE-1333).
const (
	// MaxPatternLength caps the raw pattern length.
	MaxPatternLength = 500
	// MaxGroups caps the number of capture groups.
	MaxGroups = 20
	// MaxNestingDepth caps how deeply groups may nest.
	MaxNestingDepth = 10
)

// quantifiedAlternation matches an alternation group immediately followed by a
// quantifier, e.g. (a|a)*, a classic exponential-backtracking shape.
var quantifiedAlternation = regexp.MustCompile(`\([^)]*\|[^)]*\)[*+]`)

// Validate runs the dialect-agnostic security checks on an RE2 pattern.
//
// Callers should strip any leading (?i) case-insensitivity flag before calling
// Validate (the flag's handling is dialect-specific); Validate still rejects
// inline flags other than (?i). It returns a non-nil error describing the first
// problem found, or nil if the pattern is safe to convert.
func Validate(pattern string) error {
	// 1. Length cap to bound the cost of every subsequent scan.
	if len(pattern) > MaxPatternLength {
		return fmt.Errorf("pattern length %d exceeds limit of %d characters", len(pattern), MaxPatternLength)
	}

	// 2. Reject features RE2 forbids or that cel2sql does not translate.
	// This runs before the compile check below so these patterns get a
	// descriptive message rather than Go's generic "invalid Perl syntax".
	if strings.Contains(pattern, "(?=") || strings.Contains(pattern, "(?!") {
		return errors.New("lookahead assertions (?=...), (?!...) are not supported")
	}
	if strings.Contains(pattern, "(?<=") || strings.Contains(pattern, "(?<!") {
		return errors.New("lookbehind assertions (?<=...), (?<!...) are not supported")
	}
	if strings.Contains(pattern, "(?P<") {
		return errors.New("named capture groups (?P<name>...) are not supported")
	}
	if strings.Contains(pattern, "(?m") || strings.Contains(pattern, "(?s") || strings.Contains(pattern, "(?-") {
		return errors.New("inline flags other than (?i) are not supported")
	}

	// 3. Catastrophic adjacent quantifiers, e.g. a++ / a*+.
	if matched, _ := regexp.MatchString(`[*+][*+]`, pattern); matched {
		return errors.New("regex contains catastrophic nested quantifiers that could cause ReDoS")
	}

	// 4. A quantified group whose body is itself quantified, e.g. (a+)+.
	if err := checkNestedQuantifiers(pattern); err != nil {
		return err
	}

	// 5. Capture-group count.
	groupCount := strings.Count(pattern, "(") - strings.Count(pattern, `\(`)
	if groupCount > MaxGroups {
		return fmt.Errorf("regex contains %d capture groups, exceeds limit of %d", groupCount, MaxGroups)
	}

	// 6. Quantified alternation, e.g. (a|a)*b.
	if quantifiedAlternation.MatchString(pattern) {
		return errors.New("regex contains quantified alternation that could cause ReDoS")
	}

	// 7. Group nesting depth.
	if depth := maxGroupDepth(pattern); depth > MaxNestingDepth {
		return fmt.Errorf("nesting depth %d exceeds limit of %d", depth, MaxNestingDepth)
	}

	// 8. Final catch-all: the pattern must compile under RE2. Runs last so the
	// heuristic checks above can return their descriptive messages for the
	// patterns Go's regexp rejects with a generic parse error.
	if _, err := regexp.Compile(pattern); err != nil {
		return fmt.Errorf("invalid regex pattern: %w", err)
	}

	return nil
}

// checkNestedQuantifiers rejects a quantified group that itself contains a
// quantifier (the (a+)+ catastrophic-backtracking shape). It tracks, per open
// group, whether a quantifier has been seen inside it, and when a group closes
// with a trailing quantifier while already containing one, flags it.
func checkNestedQuantifiers(pattern string) error {
	depth := 0
	groupHasQuantifier := make([]bool, 0)

	for i := 0; i < len(pattern); i++ {
		char := pattern[i]

		// Skip escaped characters.
		if i > 0 && pattern[i-1] == '\\' {
			continue
		}

		switch char {
		case '(':
			depth++
			groupHasQuantifier = append(groupHasQuantifier, false)
		case ')':
			if depth > 0 {
				depth--
				if i+1 < len(pattern) {
					next := pattern[i+1]
					if next == '*' || next == '+' || next == '?' || next == '{' {
						if len(groupHasQuantifier) > 0 && groupHasQuantifier[len(groupHasQuantifier)-1] {
							return errors.New("regex contains catastrophic nested quantifiers that could cause ReDoS")
						}
					}
				}
				if len(groupHasQuantifier) > 0 {
					// Propagate "contains a quantifier" up to the enclosing group.
					if len(groupHasQuantifier) > 1 && groupHasQuantifier[len(groupHasQuantifier)-1] {
						groupHasQuantifier[len(groupHasQuantifier)-2] = true
					}
					groupHasQuantifier = groupHasQuantifier[:len(groupHasQuantifier)-1]
				}
			}
		case '*', '+', '?', '{':
			if len(groupHasQuantifier) > 0 {
				groupHasQuantifier[len(groupHasQuantifier)-1] = true
			}
		}
	}

	return nil
}

// maxGroupDepth returns the deepest level of nested (unescaped) groups.
func maxGroupDepth(pattern string) int {
	maxDepth := 0
	current := 0
	for i := 0; i < len(pattern); i++ {
		if i > 0 && pattern[i-1] == '\\' {
			continue
		}
		switch pattern[i] {
		case '(':
			current++
			if current > maxDepth {
				maxDepth = current
			}
		case ')':
			if current > 0 {
				current--
			}
		}
	}
	return maxDepth
}
