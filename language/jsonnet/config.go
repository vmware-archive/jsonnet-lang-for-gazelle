// Copyright 2019 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package jsonnet

import (
	"flag"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

// Config states the jsonnet configuration
type Config struct {
	NativeImports map[string]bool
	IgnoreFolders map[string]bool
}

func newConfig() *Config {
	conf := &Config{
		NativeImports: make(map[string]bool, len(nativeImports)),
		IgnoreFolders: make(map[string]bool),
	}
	conf.setNativeImports(strings.Join(nativeImports, ","))
	return conf
}

// GetConfig returns a new Config within jsonnet-specs
func GetConfig(c *config.Config) *Config {
	conf := c.Exts[languageName]
	if conf == nil {
		conf = newConfig()
		return conf.(*Config)
	}
	return conf.(*Config)
}

func (*Lang) CheckFlags(fs *flag.FlagSet, c *config.Config) error { return nil }
func (*Lang) Configure(c *config.Config, rel string, f *rule.File) {
	conf := GetConfig(c)

	if f != nil {
		for _, d := range f.Directives {
			switch d.Key {
			case ignoreFoldersDirective:
				conf.setIgnoreFolders(d.Value)
			}
		}
	}
}
func (*Lang) KnownDirectives() []string {
	return []string{
		ignoreFoldersDirective,
	}
}
func (*Lang) RegisterFlags(fs *flag.FlagSet, cmd string, c *config.Config) {
	conf := GetConfig(c)
	switch cmd {
	case "fix", "update", "update-repos":
		conf.registerIgnoreFoldersFlag(fs)
	default:
	}
	c.Exts[languageName] = conf
}
