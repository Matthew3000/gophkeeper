// Static analysis service

package main

import (
	"encoding/json"
	"flag"
	"github.com/caarlos0/env/v6"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"gophkeeper/cmd/staticlint/analyzer"
	"log"
	"os"
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

type envConfig struct {
	Path string `env:"CONFIG_PATH" envDefault:"C:/repo/GO/gophkeeper/cmd/staticlint/config/config.json"`
}
type testConfig struct {
	StaticCheck []string
	StyleCheck  []string
}

func main() {
	var pathCfg envConfig
	if err := env.Parse(&pathCfg); err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&pathCfg.Path, "p", pathCfg.Path, "Config path")
	flag.Parse()

	data, err := os.ReadFile(pathCfg.Path)
	if err != nil {
		panic(err)
	}

	var cfg testConfig
	if err = json.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}

	analyzers := []*analysis.Analyzer{
		analyzer.DuplicateAnalyzer,
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		httpresponse.Analyzer,
		goc.Analyzer,
		nilerr.Analyzer,
		unusedresult.Analyzer,
		nilfunc.Analyzer,
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
