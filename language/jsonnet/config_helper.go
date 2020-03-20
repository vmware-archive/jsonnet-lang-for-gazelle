// Copyright 2019 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package jsonnet

import (
	"flag"
	"strings"
)

const (
	ignoreFoldersDirective = "jsonnet_skip_folders"
)

var (
	nativeImports = []string{".jsonnet", ".libsonnet"}
)

// stringFlag implements flags.Value for a string flag
type stringFlag func(string) error

func (f stringFlag) Set(str string) error { return f(str) }
func (f stringFlag) String() string       { return "" }

// setNativeImports implements the stringFlag type so it can be used
// to register flags
func (conf *Config) setNativeImports(extensions string) error {
	for _, extension := range strings.Split(extensions, ",") {
		// Ensure the extension is prefixed with a dot "."
		conf.NativeImports["."+strings.TrimPrefix(extension, ".")] = true
	}
	return nil
}

// IsNativeImport returns whether a given extension is a native import or not
func (conf *Config) IsNativeImport(extension string) bool {
	return conf.NativeImports[extension]
}

// setIgnoreFolders implements the stringFlag type so it can be used
// to register flags
func (conf *Config) setIgnoreFolders(folders string) error {
	for _, folder := range strings.Split(folders, ",") {
		conf.IgnoreFolders[folder] = true
	}
	return nil
}

// ShouldIgnoreFolder returns whether a given folder should be ignored or not
func (conf *Config) ShouldIgnoreFolder(folder string) bool {
	return conf.IgnoreFolders[folder]
}
func (conf *Config) registerIgnoreFoldersFlag(fs *flag.FlagSet) {
	fs.Var(
		stringFlag(conf.setIgnoreFolders),
		ignoreFoldersDirective,
		"comma-separated list of folders that should not be processed. If not specified, Gazelle will process all the folders.")
}
