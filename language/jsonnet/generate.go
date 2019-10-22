package jsonnet

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

const (
	dataImpPrivateAttr    = "_jsonnet_data_imports"
	jsonnetImpPrivateAttr = "_jsonnet_imports"
	libraryKey            = "_jsonnet_library"
)

var (
	ruleRe = regexp.MustCompile(`[^\w]+`)
)

func (*jsonnetLang) GenerateRules(args language.GenerateArgs) language.GenerateResult {
	var res language.GenerateResult

	conf := getJsonnetConfig(args.Config)
	if conf.shouldIgnoreFolder(args.Rel) {
		return res
	}

	for _, name := range args.RegularFiles {
		if !conf.isNativeImport(filepath.Ext(name)) {
			continue
		}
		finfo := jsonnetFileInfo(args, name)
		res.Gen = append(res.Gen, finfo.newLibraryRule(finfo.Path.ruleName(conf)))
	}

	sort.SliceStable(res.Gen, func(i, j int) bool {
		return res.Gen[i].Name() < res.Gen[j].Name()
	})

	res.Imports = make([]interface{}, len(res.Gen))
	for i, r := range res.Gen {
		// The rule contains a private attribute with the imports
		// that we want to resolve back to deps in Resolve.
		res.Imports[i] = r.PrivateAttr(jsonnetImpPrivateAttr)
	}

	return res
}

func (finfo FileInfo) newLibraryRule(name string) *rule.Rule {
	r := rule.NewRule("jsonnet_library", name)
	r.SetAttr("srcs", []string{finfo.Path.Filename})
	r.SetAttr("visibility", []string{"//visibility:public"})

	// Mark jsonnet imports
	imports := make(map[string]FilePath, len(finfo.Imports))
	for imp, fpath := range finfo.Imports {
		imports[imp] = fpath
	}
	r.SetPrivateAttr(jsonnetImpPrivateAttr, imports)

	// Mark data imports
	dataImports := make(map[string]FilePath, len(finfo.DataImports))
	for imp, fpath := range finfo.DataImports {
		dataImports[imp] = fpath
	}
	r.SetPrivateAttr(dataImpPrivateAttr, dataImports)

	return r
}

func (fpath FilePath) ruleName(conf *jsonnetConfig) string {
	str := fpath.Name
	str = strings.ToLower(str)
	// Replace non [a-zA-Z0-9_] characters with "_"
	str = ruleRe.ReplaceAllString(str, "_")
	return str + libraryKey
}

func (fpath FilePath) newLabel(conf *jsonnetConfig) label.Label {
	return label.New("", fpath.Dir, fpath.ruleName(conf))
}

func (fpath FilePath) newDataRef() string {
	return fmt.Sprintf("//:%s", fpath.Path)
}

func (fpath FilePath) newDataLabel() string {
	return fmt.Sprintf("//%s:%s", fpath.Dir, fpath.Filename)
}
