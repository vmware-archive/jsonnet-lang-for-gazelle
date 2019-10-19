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
	Dir      string // Relative directory to the root of the code base
	Ext      string // File name extension
	Filename string // File name
	Name     string // File name, without extension
}

// FileInfo contains metadata extracted from a file
type FileInfo struct {
	Path       FilePath            // File path information
	Imports    map[string]FilePath // Pure .jsonnet imports, which will probably depend on others
	StrImports map[string]FilePath // Plan imports, which do not depend on others
}

func jsonnetFileInfo(args language.GenerateArgs, name string) FileInfo {
	conf := getJsonnetConfig(args.Config)

	dir := args.Dir
	rel := args.Rel
	root := filepath.Clean(strings.TrimSuffix(dir, rel))
	fp := FilePath{
		Dir:      rel,
		Name:     strings.TrimSuffix(name, filepath.Ext(name)),
		Filename: name,
		Ext:      filepath.Ext(name),
	}
	info := FileInfo{
		Path:       fp,
		Imports:    make(map[string]FilePath),
		StrImports: make(map[string]FilePath),
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

			if conf.isNativeImport(ext) {
				info.Imports[imp] = impFp
				continue
			}

			// Normally, jsonnets `import foo` will import .jsonnet files
			// but you can still rename a .jsonnet file. Let's use the
			// jsonnet_allowed_imports in that case.
			if conf.isAllowedImport(ext) {
				info.StrImports[imp] = impFp
				continue
			}

			log.Printf("%s: import is not allowed. Use the jsonnet_allowed_imports directive to allowed it.", impFp.Filename)
			return info

		case match[importstrSubexpIndex] != nil:
			imp := string(match[importstrSubexpIndex])
			impPath := filepath.Join(root, rel, imp)
			impFp := resolveFilePath(root, impPath)
			ext := filepath.Ext(imp)

			if conf.isAllowedImport(ext) {
				info.StrImports[imp] = impFp
				continue
			}

			log.Printf("%s: import is not allowed. Use the jsonnet_allowed_imports directive to allowed it.", impFp.Filename)
			return info

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
		Name:     name,
		Filename: filename,
		Ext:      ext,
	}
	return fp
}
