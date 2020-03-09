// Package jsonnet provides support for jsonnet rules.
// It generates jsonnet_library rules.
//
// Configuration
//
// Configuration is largely controlled by Mode:
//
// - disable: jsonnet_library are left alone (neither
//            generated nor deleted).
// - default: jsonnet_library rules are emitted.
//
// The jsonnet mode may be set with the -jsonnet command line flag or the
// "# gazelle:jsonnet" directive.
//
// Rule generation
//
//
// Dependency resolution
//
package jsonnet

import (
	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bitnami/jsonnet-gazelle/language/jsonnet/fileinfo"
)

const (
	languageName = "jsonnet"
)

// Lang implements language.Language
type Lang struct {
	FileInfoFunc func(c *config.Config, dir string, rel string, name string) fileinfo.FileInfo
}

// NewLanguage implements the language.Language interface
func NewLanguage() language.Language { return &Lang{FileInfoFunc: NewFileInfo} }
