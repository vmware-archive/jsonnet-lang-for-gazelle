package jsonnet

import (
	"sort"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/repo"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

// Implement the Resolver interface
func (*jsonnetLang) Embeds(r *rule.Rule, from label.Label) []label.Label { return nil }
func (*jsonnetLang) Imports(c *config.Config, r *rule.Rule, f *rule.File) []resolve.ImportSpec {
	return nil
}
func (*jsonnetLang) Name() string { return languageName }
func (*jsonnetLang) Resolve(c *config.Config, ix *resolve.RuleIndex, rc *repo.RemoteCache, r *rule.Rule, imports interface{}, from label.Label) {
	if imports == nil {
		return
	}

	conf := getJsonnetConfig(c)
	deps := []string{}
	for _, fpath := range imports.(map[string]FilePath) {
		label := fpath.newLabel(conf)
		deps = append(deps, label.String())
	}

	r.DelAttr("deps")
	if len(deps) > 0 {
		sort.Strings(deps)
		r.SetAttr("deps", deps)
	}
}
