package analyzer

import (
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"strings"
)

var availableLoggers = []string{
	"log",
	"log/slog",
	"go.uber.org/zap",
}
var disabledRules = map[string]bool{}

func run(pass *analysis.Pass) (interface{}, error) {
	if pass.TypesInfo == nil {
		return nil, nil
	}
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)

		if !isLogger(pass, call) {
			return
		}

		stringsGot := extractAllStrings(call)
		stringsTrimmed := make([]string, 0, len(stringsGot))
		for _, str := range stringsGot {
			stringsTrimmed = append(stringsTrimmed, strings.Trim(str, `"`))
		}
		msg := strings.Join(stringsTrimmed, " ")

		if !disabledRules["low_first_letter"] {
			if errorMsg := LowFirstLetter(msg); errorMsg != "" {
				pass.Report(analysis.Diagnostic{
					Pos:     call.Pos(),
					Message: errorMsg,
				})
			}
		}
		if !disabledRules["english-only"] {
			if errorMsg := EnglishOnly(msg); errorMsg != "" {
				pass.Report(analysis.Diagnostic{
					Pos:     call.Pos(),
					Message: errorMsg,
				})
			}
		}
		if !disabledRules["special-symbols"] {
			if errorMsg := SpecialSymbols(msg); errorMsg != "" {
				pass.Report(analysis.Diagnostic{
					Pos:     call.Pos(),
					Message: errorMsg,
				})
			}
		}
		if !disabledRules["sensitive-words"] {
			if errorMsg := SensitiveWords(call); errorMsg != "" {
				pass.Report(analysis.Diagnostic{
					Pos:     call.Pos(),
					Message: errorMsg,
				})
			}
		}
	})
	return nil, nil
}

func isLogMeth(name string) bool {
	switch name {
	case
		"Debug", "Info", "Warn", "Error", "DPanic", "Fatal",
		"Debugf", "Infof", "Warnf", "Errorf", "DPanicf", "Panicf", "Fatalf",
		"Debugw", "Infow", "Warnw", "Errorw", "DPanicw", "Panicw", "Fatalw",
		"Print", "Printf", "Println":
		return true
	}
	return false
}

func isLogger(pass *analysis.Pass, call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	methodName := sel.Sel.Name
	if !isLogMeth(methodName) {
		return false
	}

	// проверяет по пакету метода
	selObj := pass.TypesInfo.ObjectOf(sel.Sel)
	if selObj != nil && selObj.Pkg() != nil {
		pkgPath := selObj.Pkg().Path()
		for _, knownLogger := range availableLoggers {
			if strings.HasPrefix(pkgPath, knownLogger+"/") || pkgPath == knownLogger {
				return true
			}
		}
	}

	return false
}

func extractAllStrings(node ast.Node) []string {
	var res []string
	var visit func(n ast.Node)
	visit = func(n ast.Node) {
		switch x := n.(type) {
		case *ast.BasicLit:
			if x.Kind == token.STRING {
				res = append(res, x.Value)
			}
		case *ast.CallExpr:
			for _, arg := range x.Args {
				visit(arg)
			}
		case *ast.BinaryExpr:
			visit(x.X)
			visit(x.Y)

		}
	}
	if call, ok := node.(*ast.CallExpr); ok {
		for _, arg := range call.Args {
			visit(arg)
		}
	}
	return res
}

func SetDisabled(rules map[string]bool) {
	disabledRules = rules
}
