package lintergo

import (
	"fmt"
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
	"log"

	"github.com/GameXost/LinterGo/internal/analyzer"
)

type plugin struct{}

func (p *plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		analyzer.Analyzer,
	}, nil
}

func (p *plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}

func New(settings any) (register.LinterPlugin, error) {
	var err = "123"
	log.Println("Feq", err)
	fmt.Println("settings:", settings)

	if m, ok := settings.(map[string]any); ok {
		if words, ok := m["extra-ban-words"].([]any); ok {
			extraBanWords := make([]string, 0, len(words))
			for _, word := range words {
				if wordString, ok := word.(string); ok {
					extraBanWords = append(extraBanWords, wordString)
				}
			}
			analyzer.AddBanWords(extraBanWords)
		}
		disabled := make(map[string]bool)
		if rules, ok := m["disable-flags"].([]any); ok {
			for _, f := range rules {
				if rule, ok := f.(string); ok {
					disabled[rule] = true
				}
			}
		}
		analyzer.SetDisabled(disabled)
	}

	return &plugin{}, nil
}
func init() {
	register.Plugin("loglinter", New)
}
