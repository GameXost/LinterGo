package linterGo

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"linterGo/internal/analyzer"
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

// фабрика
func New(settings any) (register.LinterPlugin, error) {
	return &plugin{}, nil
}

func init() {
	register.Plugin("loglinter", New)
}
