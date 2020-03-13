package jsonnet

import (
	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
	"github.com/google/go-jsonnet/toolutils"
)

// Importer hooks a jsonnet.Importer to parse an AST and obtain a list
// of the imports from a snippet.
type Importer struct {
	Importer jsonnet.Importer
}

func visit(n ast.Node, f func(ast.Node)) {
	f(n)
	for _, c := range toolutils.Children(n) {
		visit(c, f)
	}
}

// ParseFileImports returns the file names referenced by import and importstr
// expressions in a file.
func ParseFileImports(filename string, i *Importer) ([]string, error) {
	contents, _, err := i.Importer.Import("", filename)
	if err != nil {
		return nil, err
	}
	return i.ParseSnippetImports(filename, contents.String())
}

// ParseSnippetImports returns the file names referenced by import and importstr
// expressions in a snippet. It ensures uniqueness.
func (i *Importer) ParseSnippetImports(filename string, snippet string) ([]string, error) {
	node, err := jsonnet.SnippetToAST(filename, snippet)
	if err != nil {
		return nil, err
	}

	var files []string
	seen := map[string]struct{}{}
	collect := func(s string) {
		if _, found := seen[s]; !found {
			seen[s] = struct{}{}
			files = append(files, s)
		}
	}
	visit(node, func(n ast.Node) {
		switch i := n.(type) {
		case *ast.Import:
			collect(i.File.Value)
		case *ast.ImportStr:
			collect(i.File.Value)
		}
	})

	return files, nil
}
