package analyzer

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// MainExitAnalyzer analyzer that prohibits the use of a direct call
// to os.Exit in the main function of the main package.
var MainExitAnalyzer = &analysis.Analyzer{
	Name: "mainexit",
	Doc:  "check for use exit in main",
	Run:  mainExitRun,
}

type mainLocation struct {
	Pos token.Pos
	end token.Pos
}

func mainExitRun(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if file.Name.Name != "main" {
			continue
		}

		var ml mainLocation
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.FuncDecl:
				if x.Name.Name == "main" {
					ml = mainLocation{
						Pos: x.Pos(),
						end: x.End(),
					}
				}
			case *ast.CallExpr:
				handleExpr(x, pass, ml)
			}

			return true
		})
	}

	return nil, nil //nolint:nilnil
}

func handleExpr(exp *ast.CallExpr, p *analysis.Pass, ml mainLocation) bool {
	selExpr, ok := exp.Fun.(*ast.SelectorExpr)
	if !ok {
		return true
	}

	if selExpr.Sel.Name != "Exit" {
		return true
	}

	ident, ok := selExpr.X.(*ast.Ident)
	if !ok {
		return true
	}
	if ident.Name != "os" {
		return true
	}

	if selExpr.Pos() < ml.Pos || selExpr.End() > ml.end {
		return true
	}

	p.Reportf(exp.Pos(), "found os.Exit function call in main package, in main function")

	return false
}
