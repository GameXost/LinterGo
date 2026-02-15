package analyzer

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"strings"
)

var loggerTypes = map[string]struct{}{
	"log":             {},
	"log/slog":        {},
	"go.uber.org/zap": {},
}

func run(pass *analysis.Pass) (interface{}, error) {
	fmt.Println("Shi;jon;ojnt")
	if pass.TypesInfo == nil {
		return nil, nil
	}
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)

		//sel, ok := call.Fun.(*ast.SelectorExpr)
		//if !ok {
		//	return
		//}

		if !isLogger(pass, call) {
			return
		}

		stringsGot := extractAllStrings(call)
		stringsTrimmed := make([]string, 0, len(stringsGot))
		for _, str := range stringsGot {
			stringsTrimmed = append(stringsTrimmed, strings.Trim(str, `"`))
		}
		msg := strings.Join(stringsTrimmed, " ")
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

//func ExtractLogMsg(call *ast.CallExpr) string {
//	if len(call.Args) == 0 {
//		return ""
//	}
//	var res []string
//
//	for _, arg := range call.Args {
//		str := collectStr(arg)
//		res = append(res, str...)
//	}
//
//	result := strings.Join(res, " ")
//	return result
//}
//
//func collectStr(expr ast.Expr) []string {
//	var res []string
//
//	ast.Inspect(expr, func(n ast.Node) bool {
//		if lit, ok := n.(*ast.BasicLit); ok && lit.Kind == token.STRING {
//			str := strings.Trim(lit.Value, `"`)
//			res = append(res, str)
//		}
//		return true
//	})
//	return res
//}

//func isLoggerCall(expr ast.Expr, pass *analysis.Pass) bool {
//	ident, ok := expr.(*ast.Ident)
//	if ok {
//		obj := pass.TypesInfo.Uses[ident]
//		if pkgName, ok := obj.(*types.PkgName); ok {
//			pkgPath := pkgName.Imported().Path()
//			_, has := loggerTypes[pkgPath]
//			return has
//		}
//	}
//
//	t := pass.TypesInfo.TypeOf(expr)
//	if t == nil {
//		return false
//	}
//
//	// распаковочка, если указатель
//	for {
//		if ptr, ok := t.(*types.Pointer); ok {
//			t = ptr.Elem()
//		} else {
//			break
//		}
//	}
//
//	named, ok := t.(*types.Named)
//	if !ok {
//		return false
//	}
//	obj := named.Obj()
//	if obj == nil || obj.Pkg() == nil {
//		return false
//	}
//	pkgPath := obj.Pkg().Path()
//	//typeName := obj.Name()
//
//	_, has := loggerTypes[pkgPath]
//	return has
//
//}

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
	receiverType := pass.TypesInfo.TypeOf(sel.X)
	if receiverType == nil {
		return false
	}
	if ptr, ok := receiverType.Underlying().(*types.Pointer); ok {
		receiverType = ptr.Elem()
	}
	named, ok := receiverType.(*types.Named)
	if !ok {
		return false
	}
	if named.Obj() == nil || named.Obj().Pkg() == nil {
		return false
	}

	selectorObj := pass.TypesInfo.ObjectOf(sel.Sel)
	if selectorObj == nil {
		return false
	}

	pkg := selectorObj.Pkg()

	if pkg == nil {
		return false
	}

	//pkgPath := pkg.Path()
	pkgPath := named.Obj().Pkg().Path()
	availableLoggers := []string{
		"log",
		"log/slog",
		"go.uber.org/zap",
	}
	for _, knownLogger := range availableLoggers {
		if strings.Contains(pkgPath, knownLogger) {
			return true
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
