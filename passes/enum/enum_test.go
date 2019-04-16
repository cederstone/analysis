package enum_test

import (
	"testing"

	"github.com/cederstone/analysis/passes/enum"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestEnum(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, enum.Analyzer, "a") // loads testdata/src/a/a.go.
}
