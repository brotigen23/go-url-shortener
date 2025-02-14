package main

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

var MainExitAnalyzer = &analysis.Analyzer{
	Name: "MainExitAnalyzer",
	Doc:  "check for not using exit() from main package",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if pass.Pkg.Name() != "main" {
			continue
		}

		ast.Inspect(file, func(n ast.Node) bool {
			callExpr, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			pkgIdent, ok := selExpr.X.(*ast.Ident)
			if !ok || pkgIdent.Name != "os" || selExpr.Sel.Name != "Exit" {
				return true
			}

			pos := pass.Fset.Position(callExpr.Pos())
			pass.Reportf(callExpr.Pos(), "exit() in: %s:%d", pos.Filename, pos.Line)
			return false
		})
	}
	return nil, nil
}

func main() {
	var mychecks []*analysis.Analyzer

	// Default analyzers
	mychecks = append(mychecks, printf.Analyzer)
	mychecks = append(mychecks, shadow.Analyzer)
	mychecks = append(mychecks, structtag.Analyzer)

	// Main Exit() check
	mychecks = append(mychecks, MainExitAnalyzer)

	// Staticcheck
	for _, v := range staticcheck.Analyzers {
		mychecks = append(mychecks, v.Analyzer)
	}

	multichecker.Main(
		mychecks...,
	)
}
