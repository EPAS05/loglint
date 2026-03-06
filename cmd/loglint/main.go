package main

import (
	"github.com/EPAS05/loglint/logcheck"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(logcheck.Analyzer)
}
