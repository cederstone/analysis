package keyedlit_test

import (
        "testing"

        "golang.org/x/tools/go/analysis/analysistest"
	"github.com/cederstone/analysis/passes/keyedlit"
)

func TestNonStrict(t *testing.T) {
        keyedlit.Analyzer.Flags.Set("strict", "false")
	
        testdata := analysistest.TestData()
        analysistest.Run(t, testdata, keyedlit.Analyzer, "nonstrict") // loads testdata/src/a/a.go.
}

func TestStrict(t *testing.T) {
        keyedlit.Analyzer.Flags.Set("strict", "true")
	
        testdata := analysistest.TestData()
        analysistest.Run(t, testdata, keyedlit.Analyzer, "strict") // loads testdata/src/a/a.go.
}
