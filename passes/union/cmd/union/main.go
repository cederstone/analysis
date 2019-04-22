package main

import (
	"github.com/cederstone/analysis/passes/union"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(union.Analyzer) }
