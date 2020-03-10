package jsonnet

import (
	"sort"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/repo"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/bitnami/jsonnet-gazelle/language/jsonnet/fileinfo"
)

// Implement the Resolver interface
func (*Lang) Embeds(r *rule.Rule, from label.Label) []label.Label { return nil }

// Imports returns a list of ImportSpecs that can be used to import the rule r.
// This is used to populate RuleIndex for all the current existing rules.
func (*Lang) Imports(c *config.Config, r *rule.Rule, f *rule.File) []resolve.ImportSpec {
	// We need to resolve whether to set a label or a ref for a data file dependency.
	//
	// We can face several cases:
	//
	// 1. The data is isolated; we can include the data as a ref: deps = ["//:data/a.json"]
	//
	//	/BUILD.bazel
	//	/a.jsonnet
	//	/data/a.json
	//
	//	Then, we can generate a jsonnet_library using //:data/a.json
	//
	// 2. The data is within a pkg; we have to include the data as label: deps = ["//data:a.json"]
	//
	//	/BUILD.bazel
	//	/a.jsonnet
	//	/data/BUILD.bazel
	//	/data/a.jsonnet
	//	/data/a.json
	//
	// Therefore, we want to identify each rule by its pkg inside the workspace.

	return []resolve.ImportSpec{
		resolve.ImportSpec{Lang: "any", Imp: f.Pkg},
	}
}
func (*Lang) Name() string { return languageName }
func (*Lang) Resolve(c *config.Config, ix *resolve.RuleIndex, rc *repo.RemoteCache, r *rule.Rule, imports interface{}, from label.Label) {
	if imports == nil {
		return
	}

	var resolveFunc func(ix *resolve.RuleIndex, rc *repo.RemoteCache, r *rule.Rule, imports interface{}, from label.Label)
	switch r.Kind() {
	case libraryRule:
		resolveFunc = resolveLibraryRule
	case toJSONRule:
		resolveFunc = resolveToJSONRule
	}

	if resolveFunc != nil {
		resolveFunc(ix, rc, r, imports, from)
	}
}

func resolveLibraryRule(ix *resolve.RuleIndex, rc *repo.RemoteCache, r *rule.Rule, imports interface{}, from label.Label) {
	// Data imports will be added either as labels or refs, depending on whether its
	// directory is also a pkg or not.
	srcs := []string{}
	for _, fpath := range r.PrivateAttr(dataImpPrivateAttr).(map[string]fileinfo.FilePath) {
		// If the rules index responds to the relative directory of the data dependency,
		// it means that there is at least a rule belonging to a BUILD file in that
		// directory. In that case, we should refer to the data dependency by its
		// label.
		spec := resolve.ImportSpec{Lang: "any", Imp: fpath.Package}
		if matches := ix.FindRulesByImport(spec, "jsonnet"); len(matches) > 0 {
			srcs = append(srcs, fpath.NewDataLabel())
			continue
		}
		// Otherwise, we can refer to it by a plain ref.
		srcs = append(srcs, fpath.NewDataRef())
	}

	if len(srcs) > 0 {
		sort.Strings(srcs)
		// Leave self-import at the top
		r.SetAttr("srcs", append(r.AttrStrings("srcs"), srcs...))
	}

	// Jsonnet imports will be added as labels, as they will certainly be part of a pkg
	deps := []string{}
	for _, fpath := range imports.(map[string]fileinfo.FilePath) {
		deps = append(deps, fpath.NewLabel(libraryRulePrefix).String())
	}

	r.DelAttr("deps")
	if len(deps) > 0 {
		sort.Strings(deps)
		r.SetAttr("deps", deps)
	}
}

func resolveToJSONRule(ix *resolve.RuleIndex, rc *repo.RemoteCache, r *rule.Rule, imports interface{}, from label.Label) {
	deps := []string{}
	for _, fpath := range imports.(map[string]fileinfo.FilePath) {
		deps = append(deps, fpath.NewLabel(libraryRulePrefix).String())
	}

	r.DelAttr("deps")
	if len(deps) > 0 {
		sort.Strings(deps)
		r.SetAttr("deps", deps)
	}
}
