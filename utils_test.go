package cel2sql

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateFieldName(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		wantErr   bool
		errContains string
	}{
		// Valid field names
		{
			name:      "simple lowercase",
			fieldName: "username",
			wantErr:   false,
		},
		{
			name:      "simple uppercase",
			fieldName: "USERNAME",
			wantErr:   false,
		},
		{
			name:      "mixed case",
			fieldName: "UserName",
			wantErr:   false,
		},
		{
			name:      "with underscore",
			fieldName: "user_name",
			wantErr:   false,
		},
		{
			name:      "starts with underscore",
			fieldName: "_private",
			wantErr:   false,
		},
		{
			name:      "with numbers",
			fieldName: "field123",
			wantErr:   false,
		},
		{
			name:      "underscore and numbers",
			fieldName: "field_123_test",
			wantErr:   false,
		},
		{
			name:      "max length (63 chars)",
			fieldName: strings.Repeat("a", 63),
			wantErr:   false,
		},

		// Invalid field names - SQL injection attempts
		{
			name:      "SQL injection with semicolon",
			fieldName: "name; DROP TABLE users",
			wantErr:   true,
			errContains: "invalid field name",
		},
		{
			name:      "SQL injection with comment",
			fieldName: "field-- comment",
			wantErr:   true,
			errContains: "invalid field name",
		},
		{
			name:      "SQL injection with quotes",
			fieldName: "field' OR '1'='1",
			wantErr:   true,
			errContains: "invalid field name",
		},
		{
			name:      "SQL injection with double dash",
			fieldName: "field--",
			wantErr:   true,
			errContains: "invalid field name",
		},

		// Invalid field names - format violations
		{
			name:      "starts with number",
			fieldName: "123field",
			wantErr:   true,
			errContains: "must start with a letter or underscore",
		},
		{
			name:      "contains space",
			fieldName: "user name",
			wantErr:   true,
			errContains: "invalid field name",
		},
		{
			name:      "contains hyphen",
			fieldName: "user-name",
			wantErr:   true,
			errContains: "invalid field name",
		},
		{
			name:      "contains dot",
			fieldName: "user.name",
			wantErr:   true,
			errContains: "invalid field name",
		},
		{
			name:      "contains special characters",
			fieldName: "user@name",
			wantErr:   true,
			errContains: "invalid field name",
		},
		{
			name:      "contains newline",
			fieldName: "user\nname",
			wantErr:   true,
			errContains: "invalid field name",
		},
		{
			name:      "contains tab",
			fieldName: "user\tname",
			wantErr:   true,
			errContains: "invalid field name",
		},
		{
			name:      "empty string",
			fieldName: "",
			wantErr:   true,
			errContains: "cannot be empty",
		},

		// Invalid field names - too long
		{
			name:      "exceeds max length (64 chars)",
			fieldName: strings.Repeat("a", 64),
			wantErr:   true,
			errContains: "exceeds PostgreSQL maximum identifier length",
		},
		{
			name:      "way too long (200 chars)",
			fieldName: strings.Repeat("a", 200),
			wantErr:   true,
			errContains: "exceeds PostgreSQL maximum identifier length",
		},

		// Invalid field names - reserved SQL keywords
		{
			name:      "reserved keyword: select",
			fieldName: "select",
			wantErr:   true,
			errContains: "reserved SQL keyword",
		},
		{
			name:      "reserved keyword: SELECT (uppercase)",
			fieldName: "SELECT",
			wantErr:   true,
			errContains: "reserved SQL keyword",
		},
		{
			name:      "reserved keyword: from",
			fieldName: "from",
			wantErr:   true,
			errContains: "reserved SQL keyword",
		},
		{
			name:      "reserved keyword: where",
			fieldName: "where",
			wantErr:   true,
			errContains: "reserved SQL keyword",
		},
		{
			name:      "reserved keyword: and",
			fieldName: "and",
			wantErr:   true,
			errContains: "reserved SQL keyword",
		},
		{
			name:      "reserved keyword: or",
			fieldName: "or",
			wantErr:   true,
			errContains: "reserved SQL keyword",
		},
		{
			name:      "reserved keyword: not",
			fieldName: "not",
			wantErr:   true,
			errContains: "reserved SQL keyword",
		},
		{
			name:      "reserved keyword: null",
			fieldName: "null",
			wantErr:   true,
			errContains: "reserved SQL keyword",
		},
		{
			name:      "reserved keyword: true",
			fieldName: "true",
			wantErr:   true,
			errContains: "reserved SQL keyword",
		},
		{
			name:      "reserved keyword: false",
			fieldName: "false",
			wantErr:   true,
			errContains: "reserved SQL keyword",
		},
		{
			name:      "reserved keyword: union",
			fieldName: "union",
			wantErr:   true,
			errContains: "reserved SQL keyword",
		},
		{
			name:      "reserved keyword: drop",
			fieldName: "drop",
			wantErr:   true,
			errContains: "reserved SQL keyword",
		},
		{
			name:      "reserved keyword: insert",
			fieldName: "insert",
			wantErr:   true,
			errContains: "reserved SQL keyword",
		},
		{
			name:      "reserved keyword: update",
			fieldName: "update",
			wantErr:   true,
			errContains: "reserved SQL keyword",
		},
		{
			name:      "reserved keyword: delete",
			fieldName: "delete",
			wantErr:   true,
			errContains: "reserved SQL keyword",
		},
		{
			name:      "reserved keyword: table",
			fieldName: "table",
			wantErr:   true,
			errContains: "reserved SQL keyword",
		},
		{
			name:      "reserved keyword: create",
			fieldName: "create",
			wantErr:   true,
			errContains: "reserved SQL keyword",
		},
		{
			name:      "reserved keyword: alter",
			fieldName: "alter",
			wantErr:   true,
			errContains: "reserved SQL keyword",
		},

		// Edge cases
		{
			name:      "single character",
			fieldName: "a",
			wantErr:   false,
		},
		{
			name:      "single underscore",
			fieldName: "_",
			wantErr:   false,
		},
		{
			name:      "all underscores",
			fieldName: "___",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFieldName(tt.fieldName)

			if tt.wantErr {
				require.Error(t, err, "validateFieldName() should return error for: %s", tt.fieldName)
				if tt.errContains != "" {
					require.Contains(t, err.Error(), tt.errContains, "error message should contain expected text")
				}
			} else {
				require.NoError(t, err, "validateFieldName() should not return error for valid field name: %s", tt.fieldName)
			}
		})
	}
}

func TestValidateFieldName_AllReservedKeywords(t *testing.T) {
	// Test a sample of the reserved keywords to ensure they're all properly rejected
	sampleKeywords := []string{
		"all", "and", "any", "array", "as", "case", "cast", "check",
		"create", "cross", "current_date", "default", "delete", "distinct",
		"drop", "else", "end", "except", "false", "for", "from", "group",
		"having", "in", "inner", "insert", "intersect", "into", "is", "join",
		"left", "like", "limit", "natural", "not", "null", "offset", "on",
		"or", "order", "outer", "primary", "references", "right", "select",
		"table", "then", "to", "true", "union", "unique", "update", "user",
		"using", "when", "where", "with",
	}

	for _, keyword := range sampleKeywords {
		t.Run("keyword_"+keyword, func(t *testing.T) {
			// Test lowercase
			err := validateFieldName(keyword)
			require.Error(t, err, "Should reject reserved keyword: %s", keyword)
			require.Contains(t, err.Error(), "reserved SQL keyword")

			// Test uppercase
			err = validateFieldName(strings.ToUpper(keyword))
			require.Error(t, err, "Should reject reserved keyword (uppercase): %s", strings.ToUpper(keyword))
			require.Contains(t, err.Error(), "reserved SQL keyword")

			// Test mixed case (capitalize first letter)
			if len(keyword) > 0 {
				mixedCase := strings.ToUpper(keyword[:1]) + keyword[1:]
				err = validateFieldName(mixedCase)
				require.Error(t, err, "Should reject reserved keyword (mixed case): %s", mixedCase)
				require.Contains(t, err.Error(), "reserved SQL keyword")
			}
		})
	}
}
