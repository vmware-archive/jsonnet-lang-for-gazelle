package jsonnet

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

import (
	"github.com/bazelbuild/bazel-gazelle/language"
)

const (
	languageName = "jsonnet"
)

type jsonnetLang struct{}

// NewLanguage implements the language.Language interface
func NewLanguage() language.Language { return &jsonnetLang{} }
