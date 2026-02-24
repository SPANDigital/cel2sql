package pg_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/spandigital/cel2sql/v3/pg"
)

// TestConnectionStringLengthValidation tests that connection strings exceeding
// the maximum length are rejected
func TestConnectionStringLengthValidation(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		connStr   string
		expectErr bool
		errMsg    string
	}{
		{ //nolint:gosec // test data with fake credentials
			name:      "valid_short_connection_string",
			connStr:   "postgresql://user:pass@localhost:5432/db",
			expectErr: false,
		},
		{
			name:      "valid_connection_string_at_limit",
			connStr:   "postgresql://user:pass@localhost:5432/db?" + strings.Repeat("a", 959), // Total = 1000
			expectErr: false,
		},
		{
			name:      "connection_string_exceeds_limit",
			connStr:   "postgresql://user:pass@localhost:5432/db?" + strings.Repeat("a", 960), // Total = 1001
			expectErr: true,
			errMsg:    "connection string exceeds maximum length",
		},
		{
			name:      "very_long_connection_string",
			connStr:   strings.Repeat("postgresql://user:pass@localhost:5432/db?", 100), // Way over limit
			expectErr: true,
			errMsg:    "connection string exceeds maximum length",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := pg.NewTypeProviderWithConnection(ctx, tt.connStr)

			if tt.expectErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, provider)
			} else {
				// Note: These will fail to connect since they're not real DBs,
				// but we're testing length validation which happens first
				if err != nil {
					// Should fail with connection error, not length error
					assert.NotContains(t, err.Error(), "exceeds maximum length")
				}
				if provider != nil {
					provider.Close()
				}
			}
		})
	}
}

// TestMalformedConnectionStringsNoCredentialExposure tests that malformed
// connection strings don't expose credentials in error messages
func TestMalformedConnectionStringsNoCredentialExposure(t *testing.T) {
	ctx := context.Background()

	// Test various truly malformed connection strings that will cause parsing errors
	testCases := []struct {
		name    string
		connStr string
		secrets []string // Strings that should NOT appear in error
	}{
		{
			name:    "invalid_syntax_with_credentials",
			connStr: "postgresql://admin:secret@::invalid::/database",
			secrets: []string{"secret", "admin"},
		},
		{
			name:    "malformed_uri_with_password",
			connStr: "postgressql://user:P@ss123!@[invalid", // Invalid URI syntax
			secrets: []string{"P@ss123!", "user"},
		},
		{ //nolint:gosec // test data with fake credentials
			name:    "broken_url_encoding",
			connStr: "postgresql://app:pass%@host/db",
			secrets: []string{"pass%", "app"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			provider, err := pg.NewTypeProviderWithConnection(ctx, tc.connStr)

			// These should fail during parsing
			if err != nil {
				// Critical security check: error message must NOT contain any secrets
				errorMsg := err.Error()
				t.Logf("Error message (should be generic): %s", errorMsg)

				for _, secret := range tc.secrets {
					assert.NotContains(t, errorMsg, secret,
						"Error message MUST NOT expose credential: %s", secret)
				}

				// Verify we get generic error message (either length or connection error)
				assert.True(t,
					strings.Contains(errorMsg, "failed to create connection pool") ||
						strings.Contains(errorMsg, "exceeds maximum length"),
					"Error should be generic, got: %s", errorMsg)
			}

			if provider != nil {
				provider.Close()
			}
		})
	}
}

// TestValidConnectionWithPostgreSQL17 tests that valid connection strings work
// correctly with PostgreSQL 17 and don't break existing functionality
func TestValidConnectionWithPostgreSQL17(t *testing.T) {
	ctx := context.Background()

	// Create a PostgreSQL 17 container
	container, err := postgres.Run(ctx,
		"postgres:17",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(time.Second*60),
		),
	)
	require.NoError(t, err)

	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Errorf("failed to terminate container: %v", err)
		}
	}()

	// Get connection string
	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)
	t.Logf("Connection string length: %d characters", len(connStr))

	// Verify connection string is within limit
	assert.LessOrEqual(t, len(connStr), pg.MaxConnectionStringLength,
		"Generated connection string should be within limit")

	// Test that connection succeeds with valid connection string
	provider, err := pg.NewTypeProviderWithConnection(ctx, connStr)
	require.NoError(t, err, "Valid connection string should not error")
	require.NotNil(t, provider, "Provider should be created")
	defer provider.Close()

	// Verify provider is functional by accessing schemas
	schemas := provider.GetSchemas()
	assert.NotNil(t, schemas)
}

// TestInvalidConnectionFormats tests various invalid connection string formats
func TestInvalidConnectionFormats(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name    string
		connStr string
	}{
		{ //nolint:gosec // test data with fake credentials
			name:    "invalid_scheme",
			connStr: "mysql://user:pass@localhost/db",
		},
		{
			name:    "invalid_port",
			connStr: "postgresql://user:pass@localhost:abc/db",
		},
		{
			name:    "garbage_string",
			connStr: "this is not a connection string",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			provider, err := pg.NewTypeProviderWithConnection(ctx, tc.connStr)

			// Verify error is generic if one occurs
			if err != nil {
				errorMsg := err.Error()
				// Should be either length error or generic connection error
				assert.True(t,
					strings.Contains(errorMsg, "failed to create connection pool") ||
						strings.Contains(errorMsg, "exceeds maximum length"),
					"Error should be generic, got: %s", errorMsg)

				// Verify no connection details leaked
				assert.NotContains(t, errorMsg, tc.connStr)
			}

			if provider != nil {
				provider.Close()
			}
		})
	}
}

// TestConnectionStringSecurityProperties verifies security properties of connection handling
func TestConnectionStringSecurityProperties(t *testing.T) {
	ctx := context.Background()

	t.Run("error_does_not_leak_connection_string", func(t *testing.T) {
		// Use a malformed connection string that will actually cause a parsing error
		connStr := "postgresql://UNIQUE_USER_12345:UNIQUE_PASS_67890@[::invalid/UNIQUE_DB_ABC"

		provider, err := pg.NewTypeProviderWithConnection(ctx, connStr)

		if err != nil {
			errorMsg := err.Error()
			t.Logf("Error message: %s", errorMsg)

			// Verify none of the unique markers appear in error
			assert.NotContains(t, errorMsg, "UNIQUE_USER_12345")
			assert.NotContains(t, errorMsg, "UNIQUE_PASS_67890")
			assert.NotContains(t, errorMsg, "UNIQUE_DB_ABC")

			// Should be our generic message
			assert.Contains(t, errorMsg, "failed to create connection pool")
		}

		if provider != nil {
			provider.Close()
		}
	})

	t.Run("length_validation_happens_before_parsing", func(t *testing.T) {
		// Create a very long invalid connection string
		longConnStr := "postgresql://user:password@host/db?" + strings.Repeat("x", 2000)

		provider, err := pg.NewTypeProviderWithConnection(ctx, longConnStr)
		require.Error(t, err)
		require.Nil(t, provider)

		// Should fail with length error, not parsing error
		assert.Contains(t, err.Error(), "exceeds maximum length")
		assert.NotContains(t, err.Error(), "user")
		assert.NotContains(t, err.Error(), "password")
	})

	t.Run("no_connection_details_in_error_type", func(t *testing.T) {
		// Use malformed string that will cause parse error
		connStr := "postgresql://secret_user:secret_pass@[invalid"

		provider, err := pg.NewTypeProviderWithConnection(ctx, connStr)

		if err != nil {
			// Error should not expose details
			errorMsg := err.Error()
			t.Logf("Error message: %s", errorMsg)

			assert.Contains(t, errorMsg, "failed to create connection pool")

			// Verify it's a plain error without wrapped details
			assert.NotContains(t, errorMsg, "secret_user")
			assert.NotContains(t, errorMsg, "secret_pass")
		}

		if provider != nil {
			provider.Close()
		}
	})
}

// TestConnectionStringLengthConstant verifies the constant value is reasonable
func TestConnectionStringLengthConstant(t *testing.T) {
	// Verify constant is set to expected value
	assert.Equal(t, 1000, pg.MaxConnectionStringLength,
		"MaxConnectionStringLength should be 1000 (aligns with ODBC standard of 1024)")

	// Verify typical connection strings fit well within limit
	typicalConnStr := "postgresql://username:password@hostname.example.com:5432/database_name?sslmode=require&connect_timeout=10" //nolint:gosec // test data with fake credentials
	assert.Less(t, len(typicalConnStr), pg.MaxConnectionStringLength,
		"Typical connection string should fit within limit")

	// Verify even complex connection strings fit
	complexConnStr := "postgresql://my_application_user:My$ecureP@ssw0rd!@prod-db-primary.us-east-1.company.internal:5432/application_database?sslmode=verify-full&sslrootcert=/etc/ssl/certs/ca-bundle.crt&connect_timeout=30&application_name=MyApp"
	assert.Less(t, len(complexConnStr), pg.MaxConnectionStringLength,
		"Complex connection string should fit within limit")
}
