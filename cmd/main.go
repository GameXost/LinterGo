package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"
	"linterGo/internal/analyzer"
)

func main() {
	singlechecker.Main(analyser.Analyzer)
}
