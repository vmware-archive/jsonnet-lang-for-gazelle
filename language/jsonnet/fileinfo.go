package jsonnet

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/language"
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

func jsonnetFileInfo(args language.GenerateArgs, name string) FileInfo {
	conf := getJsonnetConfig(args.Config)

	dir := args.Dir
	rel := args.Rel
	root := filepath.Clean(strings.TrimSuffix(dir, rel))
	fp := FilePath{
		Dir:      rel,
		Ext:      filepath.Ext(name),
		Filename: name,
		Name:     strings.TrimSuffix(name, filepath.Ext(name)),
		Path:     filepath.Join(rel, name),
	}
	info := FileInfo{
		Path:        fp,
		Imports:     make(map[string]FilePath),
		DataImports: make(map[string]FilePath),
	}

	if !conf.isNativeImport(fp.Ext) {
		return info
	}

	filename := filepath.Join(rel, name)
	content, err := ioutil.ReadFile(filename)
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
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)

	fp := FilePath{
		Dir:      dir,
		Ext:      ext,
		Filename: filename,
		Name:     name,
		Path:     filepath.Join(dir, filename),
	}
	return fp
}
