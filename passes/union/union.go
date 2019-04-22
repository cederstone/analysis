// Package union performs totality checking for tagged unions.
//
// Go emulates closed tagged unions through the use of interfaces containing
// unexported methods. This pass checks that any type switch on such an
// interface value includes cases for all the types that satisfy the
// interface. The unexported 'tag' function name must not take any parameters
// nor return any values.

package union

import (
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const Doc = `find closed tagged unions.

Go emulates closed tagged unions through the use of interfaces containing
unexported methods. This pass checks that any type switch on such an interface
value includes cases for all the types that satisfy the
interface. The unexported 'tag' function name must not take any parameters nor
return any values.`

var Analyzer = &analysis.Analyzer{
	Name:             "union",
	Doc:              Doc,
	Requires:         []*analysis.Analyzer{inspect.Analyzer},
	Run:              run,
	RunDespiteErrors: false,
	FactTypes:        []analysis.Fact{new(union)},
}

type union struct {
	Interface *types.Named
	Members   []types.Type
}

func (*union) AFact() {}

func (u *union) String() string {
	return "Union"
}

func run(pass *analysis.Pass) (interface{}, error) {
	findTaggedUnions(pass)
	checkTaggedUnions(pass)
	return nil, nil
}

func findTaggedUnions(pass *analysis.Pass) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	unions := []*union{}
	// Find closed tagged unions. Consider using types.Type like guru instead
	// (See https://github.com/golang/tools/blob/master/cmd/guru/implements.go)
	nodeFilter := []ast.Node{
		(*ast.TypeSpec)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		typespec, ok := n.(*ast.TypeSpec)
		if !ok {
			return
		}
		typedef, ok := typespec.Type.(*ast.InterfaceType)
		if !ok {
			return
		}
		for _, m := range typedef.Methods.List {
			funcT, ok := m.Type.(*ast.FuncType)
			if !ok {
				continue
			}
			if funcT.Params.NumFields() != 0 || funcT.Results.NumFields() != 0 {
				// This field takes parameters or returns
				// values, this is not a union 'tag' field.
				continue
			}
			for _, name := range m.Names {
				if len(name.Name) == 0 {
					// Not sure if this is possible, but if
					// it is, this prevents an
					// out-of-bounds error below.
					continue
				}
				if name.Name[0] >= 'a' && name.Name[0] <= 'z' {
					// This type has a tag field! Record that fact.
					t := pass.TypesInfo.TypeOf(typespec.Name).(*types.Named)
					unions = append(unions, &union{Interface: t, Members: nil})
				}
			}
		}
	})
	for _, u := range unions {
		// Find types that form part of the closed tagged union.
		for _, typ := range pass.TypesInfo.Types {
			if types.IsInterface(typ.Type) {
				// Ignore interfaces as members of a closed tagged
				// union must be concrete types.
				continue
			}
			if types.AssignableTo(typ.Type, u.Interface) {
				alreadyAssigned := false
				for _, member := range u.Members {
					if types.AssignableTo(member, typ.Type) && types.AssignableTo(typ.Type, member) &&
						types.Identical(member, typ.Type) {
						alreadyAssigned = true
						break
					}
				}
				if !alreadyAssigned {
					u.Members = append(u.Members, typ.Type)
				}
			}
		}
		pass.ExportObjectFact(u.Interface.Obj(), &union{u.Interface, u.Members})
	}
}

func checkTaggedUnions(pass *analysis.Pass) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Find switch statements where the value is one of the enums and not
	// all values have case statements.
	nodeFilter := []ast.Node{
		(*ast.TypeSwitchStmt)(nil),
	}
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		stmt := n.(*ast.TypeSwitchStmt)
		var typeAssertIdent *ast.Ident
		switch x := stmt.Assign.(type) {
		case *ast.ExprStmt:
			texpr, ok := x.X.(*ast.TypeAssertExpr)
			if !ok {
				// We don't support this case.
				return
			}
			typeAssertIdent = texpr.X.(*ast.Ident)
		case *ast.AssignStmt:
			if len(x.Rhs) != 1 {
				// We don't support this case.
				return
			}
			texpr, ok := x.Rhs[0].(*ast.TypeAssertExpr)
			if !ok {
				// We don't support this case.
				return
			}
			typeAssertIdent = texpr.X.(*ast.Ident)
		default:
			// We don't support this case.
			return
		}
		t := pass.TypesInfo.TypeOf(typeAssertIdent)
		named, ok := t.(*types.Named)
		if !ok {
			return
		}

		u := new(union)
		ok = pass.ImportObjectFact(named.Obj(), u)
		if !ok {
			fmt.Printf("Cannot find object fact: %s\n", named.Obj())
			return
		}
		if !types.AssignableTo(t, u.Interface) || !types.AssignableTo(u.Interface, t) ||
			!types.Identical(t, u.Interface) {
			return
		}
		for _, member := range u.Members {
			had := false
			for _, stmt := range stmt.Body.List {
				caseClause, ok := stmt.(*ast.CaseClause)
				if !ok {
					continue
				}
				for _, caseEl := range caseClause.List {
					caseT := pass.TypesInfo.TypeOf(caseEl)
					if types.AssignableTo(member, caseT) && types.AssignableTo(caseT, member) &&
						types.Identical(member, caseT) {
						had = true
					}
				}
			}
			if !had {
				pass.Reportf(stmt.Pos(), fmt.Sprintf("non-total type switch over union: missing %s", member.String()))
			}
		}
	})
}
