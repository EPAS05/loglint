package main

import (
	"github.com/EPAS05/loglint/logcheck"
	"golang.org/x/tools/go/analysis"
)

type analyzerPlugin struct{}

func (p *analyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		logcheck.Analyzer,
	}
}

var AnalyzerPlugin analyzerPlugin

func main() {}
