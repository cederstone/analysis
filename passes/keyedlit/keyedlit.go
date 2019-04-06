// Package keyedlit defines an analysis pass that checks that impomrtant keyed
// literals fields are explicitly set.
//
// It currently checks that any field named
// Timeout or KeepAlive is explicitly set instead of relying on default values.
//
// This pass guards against users trusting the default timeout value of 0 which
// usually indicates an infinite value. Timeouts and KeepAlives should never be
// set by default and certainly never default to infinity. This pass helps
// ensure that these values were thoroughly considered in your project.
package keyedlit

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const Doc = `check for important, unspecified fields in keyed literals.

This checker reports unset Timeout / KeepAlive fields in keyed literals. These
are often overlooked (e.g., when preparing a net/http.Client) and lead to
production issues due to the default value of infinity for Timeout and no
KeepAlives.`

var Analyzer = &analysis.Analyzer{
	Name:             "keyedlit",
	Doc:              Doc,
	Requires:         []*analysis.Analyzer{inspect.Analyzer},
	Run:              run,
	RunDespiteErrors: true,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CompositeLit)(nil),
	}
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		// Skip tests.
		if strings.HasSuffix(pass.Fset.Position(n.Pos()).Filename, "_test.go") {
			return
		}
		lit := n.(*ast.CompositeLit)
		// Get the type being created.
		t := pass.TypesInfo.TypeOf(lit).Underlying()
		// Ignore it unless the type is that of a struct.
		s, ok := t.(*types.Struct)
		if !ok {
			return
		}
		var typeName string
		switch x := lit.Type.(type) {
		case *ast.SelectorExpr:
			typeName = x.Sel.String()
		case *ast.Ident:
			typeName = x.String()
		default:
			return
		}
		// Ignore unless this is a keyed composite literal.
		isKeyedLiteral := false
		for _, e := range lit.Elts {
			_, ok := e.(*ast.KeyValueExpr)
			if !ok {
				// If any of the elements are not keyed, none
				// of them are.
				return
			}
			isKeyedLiteral = true
		}
		if !isKeyedLiteral {
			return
		}
		// Loop through its fields, looking for ones that contain the
		// substrings 'Timeout' or 'KeepAlive'.
		for ii := 0; ii < s.NumFields(); ii++ {
			field := s.Field(ii)
			if mustBeSpecified(field) {
				fieldIsSpecified := false
				for _, el := range lit.Elts {
					kve := el.(*ast.KeyValueExpr)
					keyIdent, ok := kve.Key.(*ast.Ident)
					if !ok {
						// No idea what this is, skip it.
						continue
					}
					if keyIdent.Name == field.Name() {
						fieldIsSpecified = true
						break
					}
				}
				if !fieldIsSpecified {
					pass.Reportf(lit.Pos(), "unspecified field %s of %s", field.Name(), typeName)
				}
			}
		}
		return
	})

	return nil, nil
}

func mustBeSpecified(field *types.Var) bool {
	if strings.Contains(field.Name(), "KeepAlive") &&
		field.Type().String() == "time.Duration" {
		return true
	}
	if strings.Contains(field.Name(), "Timeout") &&
		field.Type().String() == "time.Duration" {
		return true
	}
	return false
}
