package jsonnet

import (
	"flag"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

type jsonnetConfig struct {
	NativeImports map[string]bool
	IgnoreFolders map[string]bool
}

func newJsonnetConfig() *jsonnetConfig {
	conf := &jsonnetConfig{
		NativeImports: make(map[string]bool, len(nativeImports)),
		IgnoreFolders: make(map[string]bool),
	}
	conf.setNativeImports(strings.Join(nativeImports, ","))
	return conf
}

func getJsonnetConfig(c *config.Config) *jsonnetConfig {
	conf := c.Exts[languageName]
	if conf == nil {
		conf = newJsonnetConfig()
		return conf.(*jsonnetConfig)
	}
	return conf.(*jsonnetConfig)
}

func (*jsonnetLang) CheckFlags(fs *flag.FlagSet, c *config.Config) error { return nil }
func (*jsonnetLang) Configure(c *config.Config, rel string, f *rule.File) {
	conf := getJsonnetConfig(c)

	if f != nil {
		for _, d := range f.Directives {
			switch d.Key {
			case ignoreFoldersDirective:
				conf.setIgnoreFolders(d.Value)
			}
		}
	}
}
func (*jsonnetLang) KnownDirectives() []string {
	return []string{
		ignoreFoldersDirective,
	}
}
func (*jsonnetLang) RegisterFlags(fs *flag.FlagSet, cmd string, c *config.Config) {
	conf := getJsonnetConfig(c)
	switch cmd {
	case "fix", "update", "update-repos":
		conf.registerIgnoreFoldersFlag(fs)
	default:
	}
	c.Exts[languageName] = conf
}
