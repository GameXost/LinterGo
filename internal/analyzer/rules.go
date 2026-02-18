package analyzer

import (
	"fmt"
	"go/ast"
	"strings"
	"unicode"
)

func LowFirstLetter(str string) string {
	str = strings.TrimLeft(str, " ")
	if len(str) == 0 {
		return ""
	}
	firstLetter := []rune(str)[0]
	if unicode.IsLetter(firstLetter) {
		if unicode.IsUpper(firstLetter) {
			return "the log message must start with lowercase letter"
		}
	}

	return ""
}

func EnglishOnly(str string) string {
	strLower := strings.ToLower(str)
	for _, r := range strLower {
		if unicode.IsLetter(r) {
			if r < 'a' || r > 'z' {
				return "the log message must be in english"
			}
		}
	}
	return ""
}

func SpecialSymbols(str string) string {
	for _, r := range str {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && !unicode.IsSpace(r) {
			return "the log message must not contain any special symbols"
		}
	}
	return ""
}

var blackList = map[string]struct{}{
	"password": {},
	"pswd":     {},
	"token":    {},
	"api":      {},
}

func SensitiveWords(call *ast.CallExpr) string {
	found := make(map[string]struct{})
	res := make([]string, 0)
	for _, arg := range call.Args {
		ast.Inspect(arg, func(n ast.Node) bool {
			if ident, ok := n.(*ast.Ident); ok {
				name := strings.ToLower(ident.Name)
				for k := range blackList {
					if strings.Contains(name, k) {
						if _, has := found[ident.Name]; !has {
							found[ident.Name] = struct{}{}
							res = append(res, ident.Name)
						}
					}
				}
			}
			return true
		})
	}
	if len(found) == 0 {
		return ""
	}
	resFound := make([]string, 0, len(found))
	for _, k := range res {
		resFound = append(resFound, k)
	}
	return fmt.Sprintf("the log message must not contain any sensitive data: %s", strings.Join(resFound, ", "))
}
