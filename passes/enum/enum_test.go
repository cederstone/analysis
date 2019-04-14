package enum_test

import (
        "testing"

        "golang.org/x/tools/go/analysis/analysistest"
	"github.com/cederstone/analysis/passes/enum"
)

func TestNonStrict(t *testing.T) {
        testdata := analysistest.TestData()
        analysistest.Run(t, testdata, enum.Analyzer, "a") // loads testdata/src/a/a.go.
}
