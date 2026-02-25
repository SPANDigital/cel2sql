// Package testutil provides multi-dialect test runners and helpers.
package testutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/spandigital/cel2sql/v3"
	"github.com/spandigital/cel2sql/v3/dialect"
	"github.com/spandigital/cel2sql/v3/testcases"
)

// RunConvertTests runs a set of ConvertTestCase entries for a given dialect.
// envFactory returns an EnvResult (CEL env + convert options) for the given EnvSetup key.
// Additional opts are appended after any env-specific options.
func RunConvertTests(
	t *testing.T,
	dialectName dialect.Name,
	cases []testcases.ConvertTestCase,
	envFactory func(envSetup string) (*EnvResult, error),
	opts ...cel2sql.ConvertOption,
) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			// Check skip
			if reason := tc.ShouldSkip(dialectName); reason != "" {
				t.Skip(reason)
			}

			// Check if we have an expectation for this dialect
			wantSQL, hasExpected := tc.ForDialect(dialectName)
			wantErr := tc.ShouldError(dialectName)

			if !hasExpected && !wantErr {
				t.Skipf("no expected SQL for dialect %s", dialectName)
			}

			// Build CEL environment
			envResult, err := envFactory(tc.EnvSetup)
			require.NoError(t, err, "failed to create CEL environment")

			// Compile CEL expression
			ast, issues := envResult.Env.Compile(tc.CELExpr)
			if issues != nil && issues.Err() != nil {
				if wantErr {
					return // expected compile error
				}
				t.Fatalf("CEL compile failed: %v", issues.Err())
			}

			// Merge options: env-specific first, then caller-provided
			allOpts := make([]cel2sql.ConvertOption, 0, len(envResult.Opts)+len(opts))
			allOpts = append(allOpts, envResult.Opts...)
			allOpts = append(allOpts, opts...)

			// Convert
			got, err := cel2sql.Convert(ast, allOpts...)
			if wantErr {
				assert.Error(t, err, "expected error for dialect %s", dialectName)
				return
			}

			if assert.NoError(t, err) {
				assert.Equal(t, wantSQL, got, "SQL mismatch for dialect %s", dialectName)
			}
		})
	}
}

// RunParameterizedTests runs a set of ParameterizedTestCase entries for a given dialect.
func RunParameterizedTests(
	t *testing.T,
	dialectName dialect.Name,
	cases []testcases.ParameterizedTestCase,
	envFactory func(envSetup string) (*EnvResult, error),
	opts ...cel2sql.ConvertOption,
) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			// Check skip
			if tc.SkipDialect != nil {
				if reason, ok := tc.SkipDialect[dialectName]; ok && reason != "" {
					t.Skip(reason)
				}
			}

			wantSQL, hasExpected := tc.WantSQL[dialectName]
			wantErr := tc.WantErr[dialectName]

			if !hasExpected && !wantErr {
				t.Skipf("no expected SQL for dialect %s", dialectName)
			}

			// Build CEL environment
			envResult, err := envFactory(tc.EnvSetup)
			require.NoError(t, err, "failed to create CEL environment")

			// Compile CEL expression
			ast, issues := envResult.Env.Compile(tc.CELExpr)
			if issues != nil && issues.Err() != nil {
				if wantErr {
					return
				}
				t.Fatalf("CEL compile failed: %v", issues.Err())
			}

			// Merge options
			allOpts := make([]cel2sql.ConvertOption, 0, len(envResult.Opts)+len(opts))
			allOpts = append(allOpts, envResult.Opts...)
			allOpts = append(allOpts, opts...)

			// Convert
			result, err := cel2sql.ConvertParameterized(ast, allOpts...)
			if wantErr {
				assert.Error(t, err)
				return
			}

			if assert.NoError(t, err) {
				assert.Equal(t, wantSQL, result.SQL, "SQL mismatch for dialect %s", dialectName)

				if wantParams, ok := tc.WantParams[dialectName]; ok && len(wantParams) > 0 {
					assert.Equal(t, wantParams, result.Parameters, "params mismatch for dialect %s", dialectName)
				}
			}
		})
	}
}

// RunAllConvertTests runs all standard test suites for a given dialect.
func RunAllConvertTests(
	t *testing.T,
	dialectName dialect.Name,
	envFactory func(envSetup string) (*EnvResult, error),
	opts ...cel2sql.ConvertOption,
) {
	t.Run("basic", func(t *testing.T) {
		RunConvertTests(t, dialectName, testcases.BasicTests(), envFactory, opts...)
	})
	t.Run("operators", func(t *testing.T) {
		RunConvertTests(t, dialectName, testcases.OperatorTests(), envFactory, opts...)
	})
	t.Run("strings", func(t *testing.T) {
		RunConvertTests(t, dialectName, testcases.StringTests(), envFactory, opts...)
	})
	t.Run("regex", func(t *testing.T) {
		RunConvertTests(t, dialectName, testcases.RegexTests(), envFactory, opts...)
	})
	t.Run("casts", func(t *testing.T) {
		RunConvertTests(t, dialectName, testcases.CastTests(), envFactory, opts...)
	})
	t.Run("arrays", func(t *testing.T) {
		RunConvertTests(t, dialectName, testcases.ArrayTests(), envFactory, opts...)
	})
	t.Run("timestamps", func(t *testing.T) {
		RunConvertTests(t, dialectName, testcases.TimestampTests(), envFactory, opts...)
	})
	t.Run("comprehensions", func(t *testing.T) {
		RunConvertTests(t, dialectName, testcases.ComprehensionTests(), envFactory, opts...)
	})
	t.Run("json", func(t *testing.T) {
		RunConvertTests(t, dialectName, testcases.JSONTests(), envFactory, opts...)
	})
}

// RunAllParameterizedTests runs all parameterized test suites for a given dialect.
func RunAllParameterizedTests(
	t *testing.T,
	dialectName dialect.Name,
	envFactory func(envSetup string) (*EnvResult, error),
	opts ...cel2sql.ConvertOption,
) {
	t.Run("parameterized", func(t *testing.T) {
		RunParameterizedTests(t, dialectName, testcases.ParameterizedTests(), envFactory, opts...)
	})
}
