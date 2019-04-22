package union_test

import (
	"testing"

	"github.com/cederstone/analysis/passes/union"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestUnion(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, union.Analyzer, "a", "b") // loads testdata/src/
}
