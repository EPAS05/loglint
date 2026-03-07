package main

import (
    "github.com/EPAS05/loglint/logcheck"
    "golang.org/x/tools/go/analysis"
)

func New(conf any) ([]*analysis.Analyzer, error) {
    return []*analysis.Analyzer{logcheck.Analyzer}, nil
}