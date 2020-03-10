package jsonnet

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bitnami/jsonnet-gazelle/language/jsonnet/fileinfo"
	"github.com/juju/errors"
)

var (
	jsonnetRe = buildJsonnetRegexp()
)

// NewFileInfo returns a FileInfo from a file path information
func NewFileInfo(c *config.Config, dir string, rel string, name string) fileinfo.FileInfo {
	conf := GetConfig(c)
	root := filepath.Clean(strings.TrimSuffix(dir, rel))
	path, err := fileinfo.NewFilePath(root, rel, name)
	if err != nil {
		log.Println(err)
		return fileinfo.FileInfo{}
	}
	info := fileinfo.FileInfo{
		Path:        path,
		Imports:     make(map[string]fileinfo.FilePath),
		DataImports: make(map[string]fileinfo.FilePath),
	}

	if err := addImports(conf, &info, path); err != nil {
		log.Println(err)
		return info
	}

	return info
}

// addImports add the imports found in the provided `path` to `info`.
func addImports(conf *Config, info *fileinfo.FileInfo, path fileinfo.FilePath) error {
	if !conf.IsNativeImport(path.Ext) {
		return nil
	}

	content, err := ioutil.ReadFile(path.Abs())
	if err != nil {
		return fmt.Errorf("error reading file %q: %w", path.Filename, err)
	}

	for _, match := range jsonnetRe.FindAllSubmatch(content, -1) {
		switch {
		case match[importSubexpIndex] != nil:
			importAbsPath, err := NormalizeImport(path, string(match[importSubexpIndex]))
			if err != nil {
				return errors.Trace(err)
			}
			importPath, err := fileinfo.NewFilePath(path.Root, importAbsPath)
			if err != nil {
				return errors.Trace(err)
			}

			if !conf.IsNativeImport(importPath.Ext) {
				// Raw JSON can be imported this way too.
				if importPath.Ext == ".json" {
					// We should handle this import as a data import, though
					info.DataImports[importPath.Path] = importPath
					continue
				}

				return errors.Errorf("%s: unknown %s extension for the `import` construct.", importPath.Filename, importPath.Ext)
			}

			info.Imports[importPath.Path] = importPath

		case match[importstrSubexpIndex] != nil:
			importAbsPath, err := NormalizeImport(path, string(match[importstrSubexpIndex]))
			if err != nil {
				return errors.Trace(err)
			}
			importPath, err := fileinfo.NewFilePath(path.Root, importAbsPath)
			if err != nil {
				return errors.Trace(err)
			}

			info.DataImports[importPath.Path] = importPath

		default:
			// Nothing to extract.
		}
	}

	return nil
}

// NormalizeImport normalizes an import string to be absolute, in any case.
// E.g. import '../foo.jsonnet' => import '/abs/to/foo.jsonnet'
func NormalizeImport(path fileinfo.FilePath, importstr string) (string, error) {
	if filepath.IsAbs(importstr) {
		if !strings.HasSuffix(importstr, path.Root) {
			return "", errors.Errorf("%q cannot be normalized. It is out of the root of the project %q", importstr, path.Root)
		}
		return importstr, nil
	}
	return path.Join(importstr), nil
}

const (
	importSubexpIndex    = 1
	importstrSubexpIndex = 2
)

// Based on https://jsonnet.org/ref/spec.html
func buildJsonnetRegexp() *regexp.Regexp {
	imp := `["']+([^"']+)["']+`
	importStmt := `import\s+` + imp
	importstrStmt := `importstr\s+` + imp
	jsonnetReSrc := strings.Join([]string{importStmt, importstrStmt}, "|")
	return regexp.MustCompile(jsonnetReSrc)
}
