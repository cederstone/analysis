// cedercheck runs the analysis passes defined in
// github.com/cederstone/analysis.
package main

import (
	"github.com/cederstone/analysis/passes/keyedlit"
	"github.com/cederstone/analysis/passes/enum"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() { multichecker.Main(enum.Analyzer, keyedlit.Analyzer) }
