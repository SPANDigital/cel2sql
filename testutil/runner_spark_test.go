package testutil_test

import (
	"testing"

	"github.com/spandigital/cel2sql/v3/dialect"
	"github.com/spandigital/cel2sql/v3/testutil"
)

func TestSparkSharedCases(t *testing.T) {
	testutil.RunAllConvertTests(t, dialect.Spark, testutil.SparkEnvFactory())
}

func TestSparkParameterizedSharedCases(t *testing.T) {
	testutil.RunAllParameterizedTests(t, dialect.Spark, testutil.SparkEnvFactory())
}
