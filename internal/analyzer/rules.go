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
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r)) {
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

// строка может быть заведомо некорректной типа Hello!!! password
// или же могут быть разделительные знаки кроме пробелов
// хотя и выведется сообщение по 2 и 3 правилам, 4 тоже стоит вывести
// так что надо убрать лишние знаки, разделить на слова и уже только тогда проверять на банлист

func SensitiveWords(call *ast.CallExpr) string {

	// опытным путем было выяснено, что в строке могу быть бан слова, а чувствительные данные априори будут передаваться через переменные
	// таким образом будем искать бан ворды в аргументах - переменных в вызове

	found := make(map[string]struct{})
	for _, arg := range call.Args {
		ast.Inspect(arg, func(n ast.Node) bool {
			if ident, ok := n.(*ast.Ident); ok {
				name := strings.ToLower(ident.Name)
				for k := range blackList {
					if strings.Contains(name, k) {
						if _, has := found[ident.Name]; !has {
							found[ident.Name] = struct{}{}
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
	for k := range found {
		resFound = append(resFound, k)
	}
	//password := "1234"
	//apiKey := "1111"
	//slog.Error("password " + password + "apikey " + apiKey)

	return fmt.Sprintf("the log message must not contain any sensitive data: %s", strings.Join(resFound, ", "))
}
