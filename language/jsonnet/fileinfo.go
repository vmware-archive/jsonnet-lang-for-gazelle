package jsonnet

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
)

var (
	jsonnetRe = buildJsonnetRegexp()
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

func jsonnetFileInfo(c *config.Config, dir string, rel string, name string) FileInfo {
	conf := getJsonnetConfig(c)
	root := filepath.Clean(strings.TrimSuffix(dir, rel))
	info := FileInfo{
		Path:        newFilePath(rel, name),
		Imports:     make(map[string]FilePath),
		DataImports: make(map[string]FilePath),
	}

	if !conf.isNativeImport(info.Path.Ext) {
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

			if !conf.isNativeImport(ext) {
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

func resolveFilePath(root string, file string) FilePath {
	filedir := filepath.Dir(file)
	dir := strings.TrimPrefix(strings.TrimPrefix(filedir, root), "/")
	filename := filepath.Base(file)
	return newFilePath(dir, filename)
}

func newFilePath(dir string, file string) FilePath {
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
