// Package analyzer search call os.Exit in main packages and report position.
// Implements analysis.Analyzer type interface for multi-check.

package analyzer

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"golang.org/x/tools/go/analysis"
	"strings"
)

var DupAnalyzer = &analysis.Analyzer{
	Name: "duplicate_killer",
	Doc:  "detects duplicated code that can be refactored into functions",
	Run:  run,
}

type snippetCounter struct {
	snippetCounts map[string]int
}

func run(pass *analysis.Pass) (interface{}, error) {
	snippetCounts := make(map[string]int)

	for _, file := range pass.Files {
		ast.Walk(&snippetCounter{snippetCounts: snippetCounts}, file)
	}

	for snippet, count := range snippetCounts {
		if count > 1 && count*4 <= countNodeLines(snippet) {
			pass.Reportf(token.NoPos, "snippet %s occurs %d times and could be refactored into a func", snippet, count)
		}
	}

	return nil, nil
}
func (v *snippetCounter) Visit(node ast.Node) ast.Visitor {
	expr, ok := node.(ast.Expr)
	if !ok {
		return v
	}

	snippet := nodeToString(expr)
	numLines := countNodeLines(snippet)

	if numLines >= 4 {
		v.snippetCounts[snippet]++
	}

	return v
}

func countNodeLines(snippet string) int {
	return strings.Count(snippet, "\n") + 1
}

func nodeToString(node ast.Node) string {
	fSet := token.NewFileSet()
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fSet, node); err != nil {
		return ""
	}

	return buf.String()
}
