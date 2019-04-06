// Package enum performs totality checking for switches over enum.
//
// Go emulates enums using a combination of int-type consts defined using iota
// for the first element. This makes it difficult to know when a switch/case
// statement covers an entire enum. This pass ensures that any switch over an
// enum explicitly lists all members.
package enum

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const Doc = `check that switch/case statements explicitly check all enum values.

Go emulates enums using a combination of int-type consts defined using iota for
the first element. This makes it difficult to know when a switch/case statement
covers an entire enum. This pass ensures that any switch over an enum
explicitly lists all members.`

var Analyzer = &analysis.Analyzer{
	Name:             "enum",
	Doc:              Doc,
	Requires:         []*analysis.Analyzer{inspect.Analyzer},
	Run:              run,
	RunDespiteErrors: true,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Find int types, since other types aren't candidates for being enums.
	candidates := map[types.Type]struct{}{}
	nodeFilter := []ast.Node{
		(*ast.TypeSpec)(nil),
	}
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		typedef := n.(*ast.TypeSpec)
		// Ignore types that aren't int aliases.
		id, ok := typedef.Type.(*ast.Ident)
		if !ok {
			return
		}
		if id.Name != "int" {
			return
		}
		t := pass.TypesInfo.TypeOf(typedef.Name)
		candidates[t] = struct{}{}
	})
	// Drop types that have a member declared directly, without using the
	// iota pattern.
	nodeFilter = []ast.Node{
		(*ast.GenDecl)(nil),
	}
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		constdecl := n.(*ast.GenDecl)
		// Ignore declarations that are not const.
		if constdecl.Tok != token.CONST {
			return
		}
		// If the RHS exists and is not iota, drop the candidate type.
		for _, spec := range constdecl.Specs {
			valspec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			if len(valspec.Names) != 1 {
				// Enums are defined one value per line. If one
				// of the values in this expression is of a
				// candidate type, delete the candidate.
				for _, name := range valspec.Names {
					t := pass.TypesInfo.TypeOf(name)
					if _, ok := candidates[t]; ok {
						delete(candidates, t)
					}
				}
				continue
			}
			t := pass.TypesInfo.TypeOf(valspec.Names[0])
			if _, ok := candidates[t]; !ok {
				continue
			}
			if len(valspec.Values) == 0 {
				continue
			}
			switch x := valspec.Values[0].(type) {
			case *ast.Ident:
				tok := x
				if tok.Name != "iota" {
					delete(candidates, t)
				}
			case *ast.CallExpr:
				call := x
				if len(call.Args) != 1 {
					delete(candidates, t)
				}
				tok, ok := call.Args[0].(*ast.Ident)
				if !ok {
					delete(candidates, t)
				}
				if tok.Name != "iota" {
					delete(candidates, t)
				}
			default:
				delete(candidates, t)
			}
		}
	})
	// Calculate the enums
	enums := map[types.Type][]types.Object{}
	for id, v := range pass.TypesInfo.Defs {
		if v == nil {
			continue
		}
		if _, ok := v.(*types.Const); !ok {
			continue
		}
		if _, ok := candidates[v.Type()]; ok {
			enums[v.Type()] = append(enums[v.Type()], pass.TypesInfo.Defs[id])
		}
	}
	// Remove unpopulated enum candidates.
	for t, members := range enums {
		if len(members) == 0 {
			delete(enums, t)
		}
	}
	// Find switch statements where the value is one of the enums and not
	// all values have case statements.
	nodeFilter = []ast.Node{
		(*ast.SwitchStmt)(nil),
	}
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		stmt := n.(*ast.SwitchStmt)
		t := pass.TypesInfo.TypeOf(stmt.Tag)
		members, ok := enums[t]
		if !ok {
			// Ignore switch statements over types that aren't
			// enums.
			return
		}
		expect := map[types.Object]struct{}{}
		for ii := range members {
			expect[members[ii]] = struct{}{}
		}
		for _, lstmt := range stmt.Body.List {
			cc, ok := lstmt.(*ast.CaseClause)
			if !ok {
				continue
			}
			for _, lexpr := range cc.List {
				lid, ok := lexpr.(*ast.Ident)
				if !ok {
					continue
				}
				obj := pass.TypesInfo.ObjectOf(lid)
				delete(expect, obj)
			}
		}
		if len(expect) != 0 {
			pass.Reportf(stmt.Pos(), "non-total switch over enum")
		}
	})
	return nil, nil
}
