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
func (conf *jsonnetConfig) setNativeImports(extensions string) error {
	for _, extension := range strings.Split(extensions, ",") {
		// Ensure the extension is prefixed with a dot "."
		conf.NativeImports["."+strings.TrimPrefix(extension, ".")] = true
	}
	return nil
}
func (conf *jsonnetConfig) isNativeImport(extension string) bool {
	return conf.NativeImports[extension]
}

// setIgnoreFolders implements the stringFlag type so it can be used
// to register flags
func (conf *jsonnetConfig) setIgnoreFolders(folders string) error {
	for _, folder := range strings.Split(folders, ",") {
		conf.IgnoreFolders[folder] = true
	}
	return nil
}
func (conf *jsonnetConfig) shouldIgnoreFolder(folder string) bool {
	return conf.IgnoreFolders[folder]
}
func (conf *jsonnetConfig) registerIgnoreFoldersFlag(fs *flag.FlagSet) {
	fs.Var(
		stringFlag(conf.setIgnoreFolders),
		ignoreFoldersDirective,
		"comma-separated list of folders that should not be processed. If not specified, Gazelle will process all the folders.")
}
