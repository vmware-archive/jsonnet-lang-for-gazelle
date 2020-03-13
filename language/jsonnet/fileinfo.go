package jsonnet

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bitnami/jsonnet-gazelle/language/jsonnet/fileinfo"
)

// NewFileInfo returns a FileInfo from a file path information
func NewFileInfo(c *config.Config, dir string, rel string, name string, importer *Importer) (*fileinfo.FileInfo, error) {
	conf := GetConfig(c)
	root := filepath.Clean(strings.TrimSuffix(dir, rel))
	path, err := fileinfo.NewFilePath(root, rel, name)
	if err != nil {
		return nil, err
	}
	info := &fileinfo.FileInfo{
		Path:        path,
		Imports:     make(map[string]fileinfo.FilePath),
		DataImports: make(map[string]fileinfo.FilePath),
	}

	if !conf.IsNativeImport(path.Ext) {
		return nil, nil
	}

	imports, err := ParseFileImports(path.Abs(), importer)
	if err != nil {
		return nil, fmt.Errorf("error parsing file %q: %v", path.Filename, err)
	}

	for _, filename := range imports {
		abs, err := NormalizeImport(path, filename)
		if err != nil {
			return nil, fmt.Errorf("error normalizing import %q: %v", filename, err)
		}
		importPath, err := fileinfo.NewFilePath(path.Root, abs)
		if err != nil {
			return nil, err
		}

		if conf.IsNativeImport(importPath.Ext) {
			info.Imports[importPath.Path] = importPath
			continue
		}

		info.DataImports[importPath.Path] = importPath
	}

	return info, nil
}

// NormalizeImport normalizes an import string to be absolute, in any case.
// E.g. import '../foo.jsonnet' => import '/abs/to/foo.jsonnet'
func NormalizeImport(path fileinfo.FilePath, importstr string) (string, error) {
	if filepath.IsAbs(importstr) {
		if !strings.HasSuffix(importstr, path.Root) {
			return "", fmt.Errorf("cannot normalize %q. It is out of the root of the project %q", importstr, path.Root)
		}
		return importstr, nil
	}
	return path.Join(importstr), nil
}
