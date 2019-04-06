package main

import (
	"github.com/cederstone/analysis/passes/enum"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(enum.Analyzer) }
