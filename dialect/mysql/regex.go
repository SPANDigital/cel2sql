package mysql

import (
	"strings"

	"github.com/spandigital/cel2sql/v3/dialect/internal/regexsafe"
)

// convertRE2ToMySQL converts an RE2 regex pattern to MySQL-compatible format.
// MySQL 8.0+ uses ICU regex, which supports \d, \w, \s and \b natively, so the
// pattern passes through almost unchanged. Shared ReDoS / unsupported-feature
// validation lives in regexsafe.Validate.
// Returns the converted pattern, whether it's case-insensitive, and any error.
func convertRE2ToMySQL(re2Pattern string) (string, bool, error) {
	caseInsensitive := false
	if strings.HasPrefix(re2Pattern, "(?i)") {
		caseInsensitive = true
		re2Pattern = re2Pattern[4:]
	}

	if err := regexsafe.Validate(re2Pattern); err != nil {
		return "", false, err
	}

	// Convert non-capturing groups (?:...) to regular groups (...).
	pattern := strings.ReplaceAll(re2Pattern, "(?:", "(")

	return pattern, caseInsensitive, nil
}
