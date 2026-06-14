// Package identsafe centralises the SQL identifier (field name) validation
// shared by every dialect.
//
// Each dialect previously reimplemented the same validateFieldName skeleton —
// empty check, optional length cap, identifier-format regexp, reserved-keyword
// lookup — differing only in the per-dialect data (max length, keyword set).
// Sharing the skeleton keeps the validation logic in one place while each
// dialect retains its own identifier rules.
package identsafe

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ValidateFieldName checks that name is a safe unquoted SQL identifier for a
// dialect.
//
//   - dialectName appears in the length-limit error message (e.g. "PostgreSQL").
//   - maxLen is the maximum identifier length; pass 0 (or negative) to skip the
//     length check for dialects that impose no practical limit.
//   - nameRE is the dialect's identifier-format pattern.
//   - reserved is the dialect's set of reserved keywords (lowercased keys).
func ValidateFieldName(name, dialectName string, maxLen int, nameRE *regexp.Regexp, reserved map[string]bool) error {
	if len(name) == 0 {
		return errors.New("field name cannot be empty")
	}

	if maxLen > 0 && len(name) > maxLen {
		return fmt.Errorf("field name %q exceeds %s maximum identifier length of %d characters", name, dialectName, maxLen)
	}

	if !nameRE.MatchString(name) {
		return fmt.Errorf("field name %q must start with a letter or underscore and contain only alphanumeric characters and underscores", name)
	}

	if reserved[strings.ToLower(name)] {
		return fmt.Errorf("field name %q is a reserved SQL keyword and cannot be used without quoting", name)
	}

	return nil
}
