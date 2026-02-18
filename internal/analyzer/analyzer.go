package analyzer

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

var Analyzer = &analysis.Analyzer{
	Name:             "loglinter",
	Doc:              "Checks logs to be proper",
	Run:              run,
	RunDespiteErrors: false,
	Requires:         []*analysis.Analyzer{inspect.Analyzer},
}
