package postgres

import (
	"strings"

	"github.com/spandigital/cel2sql/v3/dialect/internal/regexsafe"
)

// convertRE2ToPOSIX converts an RE2 regex pattern to POSIX ERE format for PostgreSQL.
// Shared ReDoS / unsupported-feature validation lives in regexsafe.Validate;
// this function adds only PostgreSQL's case-insensitivity handling and the
// RE2 → POSIX character-class translation.
// Returns: (posixPattern, caseInsensitive, error)
func convertRE2ToPOSIX(re2Pattern string) (string, bool, error) {
	// Extract the case-insensitive flag; PostgreSQL signals it via the ~* operator.
	caseInsensitive := false
	if strings.HasPrefix(re2Pattern, "(?i)") {
		caseInsensitive = true
		re2Pattern = strings.TrimPrefix(re2Pattern, "(?i)")
	}

	if err := regexsafe.Validate(re2Pattern); err != nil {
		return "", false, err
	}

	// Convert RE2 patterns to POSIX equivalents.
	posixPattern := re2Pattern
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
