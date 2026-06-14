package duckdb

import (
	"strings"

	"github.com/spandigital/cel2sql/v3/dialect/internal/regexsafe"
)

// convertRE2ToDuckDB converts an RE2 regex pattern to DuckDB-compatible format.
// DuckDB uses RE2 natively, so \d, \w, \s and \b pass through unchanged. Shared
// ReDoS / unsupported-feature validation lives in regexsafe.Validate.
// Returns the converted pattern, whether it's case-insensitive, and any error.
func convertRE2ToDuckDB(re2Pattern string) (string, bool, error) {
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
