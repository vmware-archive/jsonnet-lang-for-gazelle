package jsonnet

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bitnami/jsonnet-gazelle/language/jsonnet/fileinfo"
)

var (
	jsonnetRe = buildJsonnetRegexp()
)

// NewFileInfo returns a FileInfo from a file path information
func NewFileInfo(c *config.Config, dir string, rel string, name string) fileinfo.FileInfo {
	conf := GetConfig(c)
	root := filepath.Clean(strings.TrimSuffix(dir, rel))
	info := fileinfo.FileInfo{
		Path:        fileinfo.NewFilePath(rel, name),
		Imports:     make(map[string]fileinfo.FilePath),
		DataImports: make(map[string]fileinfo.FilePath),
	}

	if !conf.IsNativeImport(info.Path.Ext) {
		return info
	}

	filename := filepath.Join(rel, name)
	content, err := ioutil.ReadFile(filepath.Join(dir, name))
	if err != nil {
		log.Printf("%s: error reading file: %+v\n", filename, err)
		return info
	}

	for _, match := range jsonnetRe.FindAllSubmatch(content, -1) {
		switch {
		case match[importSubexpIndex] != nil:
			imp := string(match[importSubexpIndex])
			impPath := filepath.Join(root, rel, imp)
			impFp := resolveFilePath(root, impPath)
			ext := filepath.Ext(imp)

			if !conf.IsNativeImport(ext) {
				// Raw JSON can be imported this way too.
				if ext == ".json" {
					// We should handle this import as a data import, though
					info.DataImports[impFp.Path] = impFp
					continue
				}

				log.Printf("%s: unknown %s extension for the `import` construct.", impFp.Filename, ext)
				return info
			}

			info.Imports[impFp.Path] = impFp

		case match[importstrSubexpIndex] != nil:
			imp := string(match[importstrSubexpIndex])
			impPath := filepath.Join(root, rel, imp)
			impFp := resolveFilePath(root, impPath)

			info.DataImports[impFp.Path] = impFp

		default:
			// Nothing to extract.
		}
	}

	return info
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

func resolveFilePath(root string, file string) fileinfo.FilePath {
	filedir := filepath.Dir(file)
	dir := strings.TrimPrefix(strings.TrimPrefix(filedir, root), "/")
	filename := filepath.Base(file)
	return fileinfo.NewFilePath(dir, filename)
}
