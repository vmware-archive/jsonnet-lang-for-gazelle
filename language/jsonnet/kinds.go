package jsonnet

import "github.com/bazelbuild/bazel-gazelle/rule"

var (
	jsonnetKinds = map[string]rule.KindInfo{
		"jsonnet_library": {
			NonEmptyAttrs:  map[string]bool{"srcs": true},
			MergeableAttrs: map[string]bool{"srcs": true},
			ResolveAttrs:   map[string]bool{"deps": true},
		},
	}
	jsonnetLoads = []rule.LoadInfo{
		{
			Name:    "@io_bazel_rules_jsonnet//jsonnet:jsonnet.bzl",
			Symbols: []string{"jsonnet_library"},
		},
	}
)

func (*jsonnetLang) Kinds() map[string]rule.KindInfo { return jsonnetKinds }
func (*jsonnetLang) Loads() []rule.LoadInfo          { return jsonnetLoads }
