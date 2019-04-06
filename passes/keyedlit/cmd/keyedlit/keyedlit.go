package main

import (
	"github.com/cederstone/analysis/passes/keyedlit"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(keyedlit.Analyzer) }
