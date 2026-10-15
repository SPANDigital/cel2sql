package opensearch

import (
	"strings"

	"github.com/spandigital/cel2sql/v3/dialect/internal/regexsafe"
)

// convertRE2ToOpenSearch converts an RE2 regex pattern to the form OpenSearch
// SQL's REGEXP operator accepts. The SQL plugin evaluates REGEXP with
// java.util.regex.Pattern, a superset of the safe RE2 subset cel2sql accepts:
// \d, \w, \s, \b, (?:...) and an inline (?i) all work unchanged, so the
// pattern passes through as-is. Shared ReDoS / unsupported-feature validation
// lives in regexsafe.Validate.
//
// Like Spark, the (?i) flag stays inline (the Java engine honours it), so
// caseInsensitive is always reported as false; the flag is only stripped for
// validation.
func convertRE2ToOpenSearch(re2Pattern string) (string, bool, error) {
	toValidate := strings.TrimPrefix(re2Pattern, "(?i)")

	if err := regexsafe.Validate(toValidate); err != nil {
		return "", false, err
	}

	return re2Pattern, false, nil
}
