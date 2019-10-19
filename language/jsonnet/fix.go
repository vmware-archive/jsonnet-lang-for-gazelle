package jsonnet

import (
	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

func (*jsonnetLang) Fix(c *config.Config, f *rule.File) {}
