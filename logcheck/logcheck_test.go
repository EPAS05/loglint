package logcheck_test

import (
	"testing"

	"github.com/EPAS05/loglint/logcheck"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, logcheck.Analyzer, "testpkg")
}
