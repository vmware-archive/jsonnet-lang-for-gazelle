package jsonnet

import (
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/rule"
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
		// Skip non expected files extensions
		ext := filepath.Ext(name)
		if !conf.isNativeImport(ext) && !conf.isAllowedImport(ext) {
			continue
		}

		info := jsonnetFileInfo(args, name)
		res.Gen = append(res.Gen, info.newRule(conf))
	}

	sort.SliceStable(res.Gen, func(i, j int) bool {
		return res.Gen[i].Name() < res.Gen[j].Name()
	})

	res.Imports = make([]interface{}, len(res.Gen))
	for i, r := range res.Gen {
		// Add attributes marked as private so we resolve them back
		res.Imports[i] = r.PrivateAttr(config.GazelleImportsKey)
	}

	return res
}

func (finfo FileInfo) newRule(conf *jsonnetConfig) *rule.Rule {
	// We are generating a jsonnet_library per native import.
	// It will include others jsonnet_library and filegroups imports.
	if conf.isNativeImport(finfo.Path.Ext) {
		name := finfo.Path.ruleName(conf)
		return finfo.newLibraryRule(name)
	}

	// We are generating a single filegroup per non-native import.
	//
	// Non-native imports can be enabled with the jsonnet_allowed_imports
	// directive. However, it might generate unnecessary filegroups.
	//
	// In that case, we can either handcraft these filegroups or combine
	// it with the jsonnet_skip_folders directive to skip certain directories.
	name := finfo.Path.ruleName(conf)
	return finfo.newFilegroupRule(name)
}

func (finfo FileInfo) newLibraryRule(name string) *rule.Rule {
	r := rule.NewRule("jsonnet_library", name)
	r.SetAttr("srcs", []string{finfo.Path.Filename})
	r.SetAttr("visibility", []string{"//visibility:public"})

	// Add imports as private attributes so we can process them in the Resolve function.
	imports := make(map[string]FilePath, len(finfo.Imports)+len(finfo.StrImports))
	for imp, fpath := range finfo.Imports {
		imports[imp] = fpath
	}
	for imp, fpath := range finfo.StrImports {
		imports[imp] = fpath
	}
	r.SetPrivateAttr(config.GazelleImportsKey, imports)
	return r
}

func (finfo FileInfo) newFilegroupRule(name string) *rule.Rule {
	r := rule.NewRule("filegroup", name)
	r.SetAttr("srcs", []string{finfo.Path.Filename})
	r.SetAttr("visibility", []string{"//visibility:public"})
	return r
}

func (fpath FilePath) ruleName(conf *jsonnetConfig) string {
	prune := func(str string) string {
		str = strings.ToLower(str)
		// Replace non [a-zA-Z0-9_] characters with "_"
		str = ruleRe.ReplaceAllString(str, "_")
		return str
	}

	if conf.isNativeImport(fpath.Ext) {
		return prune(fpath.Name) + "_jsonnet_library"
	}

	return prune(fpath.Name) + "_filegroup"
}

func (fpath FilePath) newLabel(conf *jsonnetConfig) label.Label {
	return label.New("", fpath.Dir, fpath.ruleName(conf))
}
