package bigquery

import (
	"strings"

	"github.com/spandigital/cel2sql/v3/dialect/internal/regexsafe"
)

// convertRE2ToBigQuery converts an RE2 regex pattern to BigQuery-compatible format.
// BigQuery uses RE2 natively, so \d, \w, \s and \b pass through unchanged. Shared
// ReDoS / unsupported-feature validation lives in regexsafe.Validate.
// Returns the converted pattern, whether it's case-insensitive, and any error.
func convertRE2ToBigQuery(re2Pattern string) (string, bool, error) {
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
