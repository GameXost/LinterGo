package analyzer

import (
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"strings"
)

var loggers = map[string]struct{}{
	"slog": {},
	"log":  {},
	"zap":  {},
}

func run(pass *analysis.Pass) (interface{}, error) {
	// обход инспектором более эффективен, однако при множественных обходах
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}

		ident, ok := sel.X.(*ast.Ident)
		if !ok {
			return
		}

		pkgName := ident.Name
		//methodName := sel.Sel.Name
		if _, has := loggers[pkgName]; !has {
			return
		}

		msg := ExtractLogMsg(call)
		if msg == "" {
			return
		}

		// работа правил
		if errorMsg := LowFirstLetter(msg); errorMsg != "" {
			pass.Report(analysis.Diagnostic{
				Pos:     call.Pos(),
				Message: errorMsg,
			})
		}
		if errorMsg := EnglishOnly(msg); errorMsg != "" {
			pass.Report(analysis.Diagnostic{
				Pos:     call.Pos(),
				Message: errorMsg,
			})
		}
		if errorMsg := SpecialSymbols(msg); errorMsg != "" {
			pass.Report(analysis.Diagnostic{
				Pos:     call.Pos(),
				Message: errorMsg,
			})
		}
		if errorMsg := SensitiveWords(call); errorMsg != "" {
			pass.Report(analysis.Diagnostic{
				Pos:     call.Pos(),
				Message: errorMsg,
			})
		}
	})
	return nil, nil
}

func ExtractLogMsg(call *ast.CallExpr) string {
	if len(call.Args) == 0 {
		return ""
	}
	var res []string
	for _, arg := range call.Args {
		str := collectStr(arg)
		res = append(res, str...)
	}

	return strings.Join(res, " ")
}

func collectStr(expr ast.Expr) []string {
	var res []string

	ast.Inspect(expr, func(n ast.Node) bool {
		if lit, ok := n.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			res = append(res, strings.Trim(lit.Value, `"`))
		}
		return true
	})
	return res
}
