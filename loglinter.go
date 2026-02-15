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
	return &plugin{}, nil
}

func init() {
	fmt.Println("reigstered")
	register.Plugin("loglinter", New)
}
