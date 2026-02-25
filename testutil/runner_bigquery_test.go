package testutil_test

import (
	"testing"

	"github.com/spandigital/cel2sql/v3/dialect"
	"github.com/spandigital/cel2sql/v3/testutil"
)

func TestBigQuerySharedCases(t *testing.T) {
	testutil.RunAllConvertTests(t, dialect.BigQuery, testutil.BigQueryEnvFactory())
}

func TestBigQueryParameterizedSharedCases(t *testing.T) {
	testutil.RunAllParameterizedTests(t, dialect.BigQuery, testutil.BigQueryEnvFactory())
}
