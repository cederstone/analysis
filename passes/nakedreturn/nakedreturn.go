// Package nakedreturn defines an analysis pass that checks that there are no
// naked returns.
//
// A naked return is a return statement that does not explicitly specify the
// names of the variables being returned and instead rely on named return
// values for that purpose.
//
// Naked returns increase cognitive load by requiring the reader to remember
// the names of the return values while they are reading a function.
//
// With explicit returns, the reader can start reading at any given return
// statement and step backwards through the function definition. Since
// reversing through a function definition normally doesn't include branches,
// this strategy for code comprehension is often simpler. Reading functions in
// reverse sometimes allows the reader to stop once they have collected enough
// context.
//
// Sometimes, for very short functions, naked returns can be quite
// elegant. However, it is my belief that they are mostly abused and not worth
// it.
//
// Being explicit is usually better.
package nakedreturn

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const Doc = `check for naked returns.

This checker reports naked returns. These are return statements whose arguments
are not specified that is defined in a function that has named return values in
its signature.`

var Analyzer = &analysis.Analyzer{
	Name:             "nakedreturn",
	Doc:              Doc,
	Requires:         []*analysis.Analyzer{inspect.Analyzer},
	Run:              run,
	RunDespiteErrors: true,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.ReturnStmt)(nil),
	}
	inspect.WithStack(nodeFilter, func(n ast.Node, descending bool, stack []ast.Node) (prune bool) {
		if descending {
			return true
		}
		ret := n.(*ast.ReturnStmt)
	traverseAncestors:
		for ii := len(stack) - 2; ii > 0; ii-- {
			ancestor := stack[ii]
			var funcType *ast.FuncType
			switch fn := ancestor.(type) {
			case *ast.FuncDecl:
				funcType = fn.Type
			case *ast.FuncLit:
				funcType = fn.Type
			default:
				continue traverseAncestors
			}
			if funcType.Results == nil {
				// This function has no return values, skip.
				return true
			}
			if len(funcType.Results.List) == 0 {
				// This function has no return values, skip.
				return true
			}
			for _, field := range funcType.Results.List {
				if len(field.Names) == 0 {
					// This function does not use named return
					// variables so naked returns are impossible.
					return true
				}
			}
			// This function uses named return values. Look for naked
			// returns in this function.
			if len(ret.Results) == 0 {
				// Naked return!
				pass.Reportf(ret.Pos(), "return values not explicitly specified")
			}
			return true
		}
		return true
	})

	return nil, nil
}
