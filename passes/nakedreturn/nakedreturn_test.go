package nakedreturn_test

import (
	"testing"

	"github.com/cederstone/analysis/passes/nakedreturn"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestNakedReturn(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, nakedreturn.Analyzer, "test1") // loads testdata/src/test1/test1.go.
}
