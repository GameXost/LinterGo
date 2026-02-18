package lintergo

import (
	"fmt"
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

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
	}

	return &plugin{}, nil
}
func init() {
	fmt.Println("reeee")
	register.Plugin("loglinter", New)
}
