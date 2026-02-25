package postgres

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Regex pattern complexity limits to prevent ReDoS attacks (CWE-1333).
const (
	maxRegexPatternLength = 500
	maxRegexGroups        = 20
	maxRegexNestingDepth  = 10
)

// convertRE2ToPOSIX converts an RE2 regex pattern to POSIX ERE format for PostgreSQL.
// It performs security validation to prevent ReDoS attacks (CWE-1333).
// Returns: (posixPattern, caseInsensitive, error)
func convertRE2ToPOSIX(re2Pattern string) (string, bool, error) {
	// 1. Check pattern length to prevent processing extremely long patterns
	if len(re2Pattern) > maxRegexPatternLength {
		return "", false, fmt.Errorf("pattern length %d exceeds limit of %d characters", len(re2Pattern), maxRegexPatternLength)
	}

	// 2. Extract case-insensitive flag if present
	caseInsensitive := false
	if strings.HasPrefix(re2Pattern, "(?i)") {
		caseInsensitive = true
		re2Pattern = strings.TrimPrefix(re2Pattern, "(?i)")
	}

	// 3. Detect unsupported RE2 features and return errors
	if strings.Contains(re2Pattern, "(?=") || strings.Contains(re2Pattern, "(?!") {
		return "", false, errors.New("lookahead assertions (?=...), (?!...) are not supported in PostgreSQL POSIX regex")
	}
	if strings.Contains(re2Pattern, "(?<=") || strings.Contains(re2Pattern, "(?<!") {
		return "", false, errors.New("lookbehind assertions (?<=...), (?<!...) are not supported in PostgreSQL POSIX regex")
	}
	if strings.Contains(re2Pattern, "(?P<") {
		return "", false, errors.New("named capture groups (?P<name>...) are not supported in PostgreSQL POSIX regex")
	}
	if strings.Contains(re2Pattern, "(?m") || strings.Contains(re2Pattern, "(?s") || strings.Contains(re2Pattern, "(?-") {
		return "", false, errors.New("inline flags other than (?i) are not supported in PostgreSQL POSIX regex")
	}

	// 4. Detect catastrophic nested quantifiers
	if matched, _ := regexp.MatchString(`[*+][*+]`, re2Pattern); matched {
		return "", false, errors.New("regex contains catastrophic nested quantifiers that could cause ReDoS")
	}

	// Check for groups that contain quantifiers and are themselves quantified
	depth := 0
	groupHasQuantifier := make([]bool, 0)

	for i := 0; i < len(re2Pattern); i++ {
		char := re2Pattern[i]

		// Skip escaped characters
		if i > 0 && re2Pattern[i-1] == '\\' {
			continue
		}

		switch char {
		case '(':
			depth++
			groupHasQuantifier = append(groupHasQuantifier, false)
		case ')':
			if depth > 0 {
				depth--
				if i+1 < len(re2Pattern) {
					nextChar := re2Pattern[i+1]
					if nextChar == '*' || nextChar == '+' || nextChar == '?' || nextChar == '{' {
						if len(groupHasQuantifier) > 0 && groupHasQuantifier[len(groupHasQuantifier)-1] {
							return "", false, errors.New("regex contains catastrophic nested quantifiers that could cause ReDoS")
						}
					}
				}
				if len(groupHasQuantifier) > 0 {
					if len(groupHasQuantifier) > 1 {
						if groupHasQuantifier[len(groupHasQuantifier)-1] {
							groupHasQuantifier[len(groupHasQuantifier)-2] = true
						}
					}
					groupHasQuantifier = groupHasQuantifier[:len(groupHasQuantifier)-1]
				}
			}
		case '*', '+', '?':
			if len(groupHasQuantifier) > 0 {
				groupHasQuantifier[len(groupHasQuantifier)-1] = true
			}
		case '{':
			if len(groupHasQuantifier) > 0 {
				groupHasQuantifier[len(groupHasQuantifier)-1] = true
			}
		}
	}

	// 5. Count and limit capture groups
	groupCount := strings.Count(re2Pattern, "(") - strings.Count(re2Pattern, `\(`)
	if groupCount > maxRegexGroups {
		return "", false, fmt.Errorf("regex contains %d capture groups, exceeds limit of %d", groupCount, maxRegexGroups)
	}

	// 6. Detect exponential alternation patterns
	alternationPattern := regexp.MustCompile(`\([^)]*\|[^)]*\)[*+]`)
	if alternationPattern.MatchString(re2Pattern) {
		return "", false, errors.New("regex contains quantified alternation that could cause ReDoS")
	}

	// 7. Check nesting depth
	maxDepthVal := 0
	currentDepth := 0
	for _, char := range re2Pattern {
		if char == '(' && !strings.HasSuffix(re2Pattern[:strings.LastIndex(re2Pattern, string(char))], `\`) {
			currentDepth++
			if currentDepth > maxDepthVal {
				maxDepthVal = currentDepth
			}
		} else if char == ')' && !strings.HasSuffix(re2Pattern[:strings.LastIndex(re2Pattern, string(char))], `\`) {
			currentDepth--
		}
	}
	if maxDepthVal > maxRegexNestingDepth {
		return "", false, fmt.Errorf("nesting depth %d exceeds limit of %d", maxDepthVal, maxRegexNestingDepth)
	}

	// Passed all security checks - proceed with conversion
	posixPattern := re2Pattern

	// Convert RE2 patterns to POSIX equivalents
	posixPattern = strings.ReplaceAll(posixPattern, `\b`, `\y`)
	posixPattern = strings.ReplaceAll(posixPattern, `\B`, `[^[:alnum:]_]`)
	posixPattern = strings.ReplaceAll(posixPattern, `\d`, `[[:digit:]]`)
	posixPattern = strings.ReplaceAll(posixPattern, `\D`, `[^[:digit:]]`)
	posixPattern = strings.ReplaceAll(posixPattern, `\w`, `[[:alnum:]_]`)
	posixPattern = strings.ReplaceAll(posixPattern, `\W`, `[^[:alnum:]_]`)
	posixPattern = strings.ReplaceAll(posixPattern, `\s`, `[[:space:]]`)
	posixPattern = strings.ReplaceAll(posixPattern, `\S`, `[^[:space:]]`)
	posixPattern = strings.ReplaceAll(posixPattern, `(?:`, `(`)

	return posixPattern, caseInsensitive, nil
}
