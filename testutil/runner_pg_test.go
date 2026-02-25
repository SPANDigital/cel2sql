package testutil_test

import (
	"testing"

	"github.com/spandigital/cel2sql/v3/dialect"
	"github.com/spandigital/cel2sql/v3/testutil"
)

func TestPostgreSQLSharedCases(t *testing.T) {
	testutil.RunAllConvertTests(t, dialect.PostgreSQL, testutil.PostgreSQLEnvFactory())
}

func TestPostgreSQLParameterizedSharedCases(t *testing.T) {
	testutil.RunAllParameterizedTests(t, dialect.PostgreSQL, testutil.PostgreSQLEnvFactory())
}
