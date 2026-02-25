package mysql

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

// convertRE2ToMySQL converts an RE2 regex pattern to MySQL-compatible format.
// MySQL 8.0+ uses ICU regex which supports most RE2 features.
// Returns the converted pattern, whether it's case-insensitive, and any error.
func convertRE2ToMySQL(re2Pattern string) (string, bool, error) {
	// 1. Pattern length validation
	if len(re2Pattern) > maxRegexPatternLength {
		return "", false, fmt.Errorf("pattern length %d exceeds limit of %d characters", len(re2Pattern), maxRegexPatternLength)
	}

	// 2. Validate pattern compiles
	if _, err := regexp.Compile(re2Pattern); err != nil {
		return "", false, fmt.Errorf("invalid regex pattern: %w", err)
	}

	// 3. Detect unsupported RE2 features
	if strings.Contains(re2Pattern, "(?=") || strings.Contains(re2Pattern, "(?!") {
		return "", false, errors.New("lookahead assertions (?=...), (?!...) are not supported in MySQL regex")
	}
	if strings.Contains(re2Pattern, "(?<=") || strings.Contains(re2Pattern, "(?<!") {
		return "", false, errors.New("lookbehind assertions (?<=...), (?<!...) are not supported in MySQL regex")
	}
	if strings.Contains(re2Pattern, "(?P<") {
		return "", false, errors.New("named capture groups (?P<name>...) are not supported in MySQL regex")
	}

	// 4. Detect catastrophic nested quantifiers
	if matched, _ := regexp.MatchString(`[*+][*+]`, re2Pattern); matched {
		return "", false, errors.New("regex contains catastrophic nested quantifiers that could cause ReDoS")
	}

	// 5. Check for nested quantifiers in groups
	depth := 0
	groupHasQuantifier := make([]bool, 0)
	for i := 0; i < len(re2Pattern); i++ {
		char := re2Pattern[i]
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
					groupHasQuantifier = groupHasQuantifier[:len(groupHasQuantifier)-1]
				}
			}
		case '*', '+', '?', '{':
			for j := range groupHasQuantifier {
				groupHasQuantifier[j] = true
			}
		}
	}

	// 6. Check group count limit
	groupCount := strings.Count(re2Pattern, "(") - strings.Count(re2Pattern, "\\(")
	if groupCount > maxRegexGroups {
		return "", false, fmt.Errorf("regex contains %d capture groups, exceeds limit of %d", groupCount, maxRegexGroups)
	}

	// 7. Check for quantified alternation
	quantifiedAlternation := regexp.MustCompile(`\([^)]*\|[^)]*\)[*+]`)
	if quantifiedAlternation.MatchString(re2Pattern) {
		return "", false, errors.New("regex contains quantified alternation that could cause ReDoS")
	}

	// 8. Check nesting depth
	maxDepthVal := 0
	currentDepth := 0
	for i := 0; i < len(re2Pattern); i++ {
		if i > 0 && re2Pattern[i-1] == '\\' {
			continue
		}
		switch re2Pattern[i] {
		case '(':
			currentDepth++
			if currentDepth > maxDepthVal {
				maxDepthVal = currentDepth
			}
		case ')':
			if currentDepth > 0 {
				currentDepth--
			}
		}
	}
	if maxDepthVal > maxRegexNestingDepth {
		return "", false, fmt.Errorf("nesting depth %d exceeds limit of %d", maxDepthVal, maxRegexNestingDepth)
	}

	// Process pattern: extract case-insensitivity, convert features
	caseInsensitive := false
	pattern := re2Pattern

	// Handle (?i) flag
	if strings.HasPrefix(pattern, "(?i)") {
		caseInsensitive = true
		pattern = pattern[4:]
	}

	// Handle inline flags other than (?i) at start
	if strings.Contains(pattern, "(?m") || strings.Contains(pattern, "(?s") || strings.Contains(pattern, "(?-") {
		return "", false, errors.New("inline flags other than (?i) are not supported in MySQL regex")
	}

	// Convert non-capturing groups (?:...) to regular groups (...)
	pattern = strings.ReplaceAll(pattern, "(?:", "(")

	// MySQL ICU regex supports \d, \w, \s natively - no conversion needed
	// Convert \b word boundary to MySQL's \b (same syntax in ICU)
	// No conversion needed for MySQL 8.0+

	return pattern, caseInsensitive, nil
}
