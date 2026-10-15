package testutil_test

import (
	"testing"

	"github.com/spandigital/cel2sql/v3/dialect"
	"github.com/spandigital/cel2sql/v3/testutil"
)

func TestOpenSearchSharedCases(t *testing.T) {
	testutil.RunAllConvertTests(t, dialect.OpenSearch, testutil.OpenSearchEnvFactory())
}

func TestOpenSearchParameterizedSharedCases(t *testing.T) {
	testutil.RunAllParameterizedTests(t, dialect.OpenSearch, testutil.OpenSearchEnvFactory())
}
