package analyser

import (
	"fmt"
	"strings"
	"unicode"
)

func LowFirstLetter(str string) string {
	if len(str) == 0 {
		return ""
	}
	firstLetter := []rune(str)[0]
	if unicode.IsLower(firstLetter) {
		return ""
	}
	return "the log message must start with lowercase letter "
}

func AllEnglishLetters(str string) string {
	for _, r := range str {
		if unicode.IsLetter(r) && (r < 'a' || r > 'z') {
			return "the log message must be in english"
		}
	}
	return ""
}

func SpecialSymbols(str string) string {
	for _, r := range str {
		if !unicode.IsLetter(r) {
			return "the log message must not contain any special symbols "
		}
	}
	return ""
}

var blackList = map[string]struct{}{
	"password": {},
	"pswd":     {},
	"token":    {},
}

func SensitiveWords(str string) string {
	// строка может быть заведомо некорректной типа Hello!!! password
	// или же могут быть разделительные знаки кроме пробелов
	// хотя и выведется сообщение по 2 и 3 правилам, 4 тоже стоит вывести
	// так что надо убрать лишние знаки, разделить на слова и уже только тогда проверять на банлист

	// формируем строку только из букв (в целом неважно, ру или англ)
	var correctStr strings.Builder
	for _, r := range str {
		if unicode.IsLetter(r) || unicode.IsSpace(r) {
			correctStr.WriteRune(unicode.ToLower(r))
		}
	}
	words := strings.Fields(correctStr.String())
	specWordsFound := make([]string, 0, 2)
	for _, word := range words {
		if _, exists := blackList[word]; exists == true {
			specWordsFound = append(specWordsFound, word)
		}
	}
	if len(specWordsFound) == 0 {
		return ""
	}
	return fmt.Sprintf("the log message must not contain anny sensitive data: %s", strings.Join(specWordsFound, ", "))
}
