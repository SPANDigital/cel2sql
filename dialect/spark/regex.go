package spark

import (
	"strings"

	"github.com/spandigital/cel2sql/v3/dialect/internal/regexsafe"
)

// convertRE2ToSpark converts an RE2 regex pattern to Spark/Java regex format.
// Spark uses java.util.regex.Pattern, a superset of the safe RE2 subset
// cel2sql accepts: it natively supports \d, \w, \s, \b, (?:...) and inline
// (?i), so patterns pass through unchanged. Shared ReDoS / unsupported-feature
// validation lives in regexsafe.Validate.
//
// Unlike the other dialects, Spark keeps the (?i) flag inline (its regex engine
// honours it), so caseInsensitive is always reported as false; the flag is only
// stripped for the purpose of validation counting.
func convertRE2ToSpark(re2Pattern string) (string, bool, error) {
	toValidate := strings.TrimPrefix(re2Pattern, "(?i)")

	if err := regexsafe.Validate(toValidate); err != nil {
		return "", false, err
	}

	return re2Pattern, false, nil
}
