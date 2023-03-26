// Static analysis service

package main

import (
	"encoding/json"
	"gophkeeper/cmd/staticlint/analyzer"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"

	goc "github.com/go-critic/go-critic/checkers/analyzer"
	"github.com/gostaticanalysis/nilerr"
)

const Config = `config/config.json`

type StaticTestConfig struct {
	path string
}
type ConfigData struct {
	StaticCheck []string
	StyleCheck  []string
}

func main() {
	data, err := os.ReadFile(filepath.Join(filepath.Dir("C:/repo/GO/gophkeeper/cmd/staticlint/"), Config))
	if err != nil {
		panic(err)
	}

	var cfg ConfigData
	if err = json.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}

	analyzers := []*analysis.Analyzer{
		analyzer.DupAnalyzer,
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		httpresponse.Analyzer,
		goc.Analyzer,
		nilerr.Analyzer,
	}

	for _, v := range staticcheck.Analyzers {
		for _, sc := range cfg.StaticCheck {
			if strings.HasPrefix(v.Analyzer.Name, sc) {
				analyzers = append(analyzers, v.Analyzer)
			}
		}
	}
	for _, v := range stylecheck.Analyzers {
		for _, sc := range cfg.StyleCheck {
			if strings.HasPrefix(v.Analyzer.Name, sc) {
				analyzers = append(analyzers, v.Analyzer)
			}
		}
	}

	multichecker.Main(analyzers...)
}
