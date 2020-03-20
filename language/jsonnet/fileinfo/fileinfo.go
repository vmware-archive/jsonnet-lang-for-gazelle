// Copyright 2019 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package fileinfo

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/label"
)

var (
	ruleRe = regexp.MustCompile(`[^\w]+`)
)

// FilePath contains the path information for a file
type FilePath struct {
	Root     string // Absolute path to the root of the workspace
	Package  string // Workspace-relative path of the directory containing the file.
	Ext      string // File name extension
	Filename string // File name
	Name     string // File name, without extension
	Path     string // File path, relative to the root of the workspace
}

// FileInfo contains metadata extracted from a file
type FileInfo struct {
	Path        FilePath            // File path information
	Imports     map[string]FilePath // Jsonnet imports, from import
	DataImports map[string]FilePath // Data imports, from importstr
}

// Join filepath.Joins any number of path elements into a single path prepending
// the FilePath Root and Package paths to the path elements.
func (fp FilePath) Join(elem ...string) string {
	return filepath.Join(append([]string{fp.Root, fp.Package}, elem...)...)
}

// RuleName computes a rule name for a given file path
func (fp FilePath) RuleName(prefix string) string {
	// Replace non [a-zA-Z0-9_] characters with "_"
	str := ruleRe.ReplaceAllString(strings.ToLower(fp.Name), "_")
	return str + "_" + prefix
}

// NewLabel computes a label for a given file path
func (fp FilePath) NewLabel(prefix string) label.Label {
	return label.New("", fp.Package, fp.RuleName(prefix))
}

// NewDataRef returns a ref for the given data file path
func (fp FilePath) NewDataRef() string {
	return fmt.Sprintf("//:%s", fp.Path)
}

// NewDataLabel returns a label for the given data file path
func (fp FilePath) NewDataLabel() string {
	return fmt.Sprintf("//%s:%s", fp.Package, fp.Filename)
}

// Abs returns the absolute path to the file
func (fp FilePath) Abs() string {
	return fp.Join(fp.Filename)
}

// NewFilePath constructs a FilePath structure given a root directory and one or more path elements.
//
// The path elements are filepath.Join-ed together interpreted as relative to root.
//
// If the path elements forms an absolute path, NewFilePath returns a FilePath structure
// so Package is relative to the root directory. Therefore, Package will contain as many ".." symbols
// in order to "escape" from its absolute path and then "join" the root directory.
func NewFilePath(root string, elem ...string) (FilePath, error) {
	// We don't know the elements shape so let's join them and then split them
	// into dir and file
	path := filepath.Join(elem...)

	// Get rid of the root part in absolute paths
	if filepath.IsAbs(path) {
		p, err := filepath.Rel(root, path)
		if err != nil {
			return FilePath{}, err
		}
		path = p
	}

	dir, file := filepath.Split(path)
	ext := filepath.Ext(file)
	return FilePath{
		Root:     root,
		Package:  strings.TrimSuffix(dir, "/"),
		Ext:      ext,
		Filename: file,
		Name:     strings.TrimSuffix(file, ext),
		Path:     path,
	}, nil
}
