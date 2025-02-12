package main

import (
	"strings"

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
	return nil, nil
}

func main() {
	var mychecks []*analysis.Analyzer

	// Default analyzers
	mychecks = append(mychecks, printf.Analyzer)
	mychecks = append(mychecks, shadow.Analyzer)
	mychecks = append(mychecks, structtag.Analyzer)

	// Staticcheck
	for _, v := range staticcheck.Analyzers {
		if strings.HasPrefix(v.Analyzer.Name, "SA90") {
			mychecks = append(mychecks, v.Analyzer)
		}
		// Check for infinite loops
		if v.Analyzer.Name == "S1006" {
			mychecks = append(mychecks, v.Analyzer)
		}
		// Poorly chosen name for error variable
		if v.Analyzer.Name == "ST1012" {
			mychecks = append(mychecks, v.Analyzer)
		}
		// Apply De Morgan’s law
		if v.Analyzer.Name == "QF1001" {
			mychecks = append(mychecks, v.Analyzer)
		}
	}

	multichecker.Main(
		mychecks...,
	)
}
