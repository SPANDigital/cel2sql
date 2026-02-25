package testutil_test

import (
	"testing"

	"github.com/spandigital/cel2sql/v3/dialect"
	"github.com/spandigital/cel2sql/v3/testutil"
)

func TestMySQLSharedCases(t *testing.T) {
	testutil.RunAllConvertTests(t, dialect.MySQL, testutil.MySQLEnvFactory())
}
