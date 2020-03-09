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
	Dir      string // Relative directory to the root of the code base; it works as pkg ref
	Ext      string // File name extension
	Filename string // File name
	Name     string // File name, without extension
	Path     string // File path, relative to the root of the code base
}

// FileInfo contains metadata extracted from a file
type FileInfo struct {
	Path        FilePath            // File path information
	Imports     map[string]FilePath // Jsonnet imports, from import
	DataImports map[string]FilePath // Data imports, from importstr
}

// RuleName computes a rule name for a given file path
func (fpath FilePath) RuleName(prefix string) string {
	// Replace non [a-zA-Z0-9_] characters with "_"
	str := ruleRe.ReplaceAllString(strings.ToLower(fpath.Name), "_")
	return str + "_" + prefix
}

// NewLabel computes a label for a given file path
func (fpath FilePath) NewLabel(prefix string) label.Label {
	return label.New("", fpath.Dir, fpath.RuleName(prefix))
}

// NewDataRef returns a ref for the given data file path
func (fpath FilePath) NewDataRef() string {
	return fmt.Sprintf("//:%s", fpath.Path)
}

// NewDataLabel returns a label for the given data file path
func (fpath FilePath) NewDataLabel() string {
	return fmt.Sprintf("//%s:%s", fpath.Dir, fpath.Filename)
}

// NewFilePath returns a file path from a given `dir` and `file` inside `dir`
func NewFilePath(dir string, file string) FilePath {
	ext := filepath.Ext(file)
	name := strings.TrimSuffix(file, ext)
	fp := FilePath{
		Dir:      dir,
		Ext:      ext,
		Filename: file,
		Name:     name,
		Path:     filepath.Join(dir, file),
	}
	return fp
}
