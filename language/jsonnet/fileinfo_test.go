package jsonnet_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bitnami/jsonnet-gazelle/language/jsonnet"
	"github.com/bitnami/jsonnet-gazelle/language/jsonnet/fileinfo"
)

// FilePath includes a Root field that prevents comparing the FileInfo object
// properly. This function removes the Root field from all imports and the main
// FilePath.
func normalizeFileInfo(info *fileinfo.FileInfo) {
	info.Path.Root = ""
	for path, value := range info.Imports {
		value.Root = ""
		info.Imports[path] = value
	}
	for path, value := range info.DataImports {
		value.Root = ""
		info.DataImports[path] = value
	}
}

func TestJsonnetFileInfo(t *testing.T) {
	testCases := []struct {
		desc, dir, rel, name, content string
		want                          fileinfo.FileInfo
	}{
		{
			desc:    "empty",
			rel:     "pkg/foo",
			name:    "bar.jsonnet",
			content: "",
			want: fileinfo.FileInfo{
				Path:        fileinfo.FilePath{Package: "pkg/foo", Ext: ".jsonnet", Filename: "bar.jsonnet", Name: "bar", Path: "pkg/foo/bar.jsonnet"},
				Imports:     map[string]fileinfo.FilePath{},
				DataImports: map[string]fileinfo.FilePath{},
			},
		}, {
			desc:    "differnt quotes imports",
			rel:     "pkg/foo",
			name:    "bar.jsonnet",
			content: `(import 'singlequotes.jsonnet') + (import "doublequotes.jsonnet")`,
			want: fileinfo.FileInfo{
				Path: fileinfo.FilePath{Package: "pkg/foo", Ext: ".jsonnet", Filename: "bar.jsonnet", Name: "bar", Path: "pkg/foo/bar.jsonnet"},
				Imports: map[string]fileinfo.FilePath{
					"pkg/foo/singlequotes.jsonnet": {Package: "pkg/foo", Ext: ".jsonnet", Filename: "singlequotes.jsonnet", Name: "singlequotes", Path: "pkg/foo/singlequotes.jsonnet"},
					"pkg/foo/doublequotes.jsonnet": {Package: "pkg/foo", Ext: ".jsonnet", Filename: "doublequotes.jsonnet", Name: "doublequotes", Path: "pkg/foo/doublequotes.jsonnet"},
				},
				DataImports: map[string]fileinfo.FilePath{},
			},
		}, {
			desc:    "libsonnet import",
			rel:     "pkg/foo",
			name:    "bar.jsonnet",
			content: "import 'demo.libsonnet'",
			want: fileinfo.FileInfo{
				Path: fileinfo.FilePath{Package: "pkg/foo", Ext: ".jsonnet", Filename: "bar.jsonnet", Name: "bar", Path: "pkg/foo/bar.jsonnet"},
				Imports: map[string]fileinfo.FilePath{
					"pkg/foo/demo.libsonnet": {Package: "pkg/foo", Ext: ".libsonnet", Filename: "demo.libsonnet", Name: "demo", Path: "pkg/foo/demo.libsonnet"},
				},
				DataImports: map[string]fileinfo.FilePath{},
			},
		}, {
			desc:    "different folder imports",
			rel:     "pkg/foo",
			name:    "bar.jsonnet",
			content: "(import '../pkg.libsonnet') + (import '../../root.jsonnet')",
			want: fileinfo.FileInfo{
				Path: fileinfo.FilePath{Package: "pkg/foo", Ext: ".jsonnet", Filename: "bar.jsonnet", Name: "bar", Path: "pkg/foo/bar.jsonnet"},
				Imports: map[string]fileinfo.FilePath{
					"pkg/pkg.libsonnet": {Package: "pkg", Ext: ".libsonnet", Filename: "pkg.libsonnet", Name: "pkg", Path: "pkg/pkg.libsonnet"},
					"root.jsonnet":      {Package: "", Ext: ".jsonnet", Filename: "root.jsonnet", Name: "root", Path: "root.jsonnet"},
				},
				DataImports: map[string]fileinfo.FilePath{},
			},
		}, {
			desc:    "data import",
			rel:     "pkg/foo",
			name:    "bar.jsonnet",
			content: "importstr 'data/db.json'",
			want: fileinfo.FileInfo{
				Path:    fileinfo.FilePath{Package: "pkg/foo", Ext: ".jsonnet", Filename: "bar.jsonnet", Name: "bar", Path: "pkg/foo/bar.jsonnet"},
				Imports: map[string]fileinfo.FilePath{},
				DataImports: map[string]fileinfo.FilePath{
					"pkg/foo/data/db.json": {Package: "pkg/foo/data", Ext: ".json", Filename: "db.json", Name: "db", Path: "pkg/foo/data/db.json"},
				},
			},
		}, {
			desc:    "mixed data and jsonnet imports",
			rel:     "pkg/foo",
			name:    "bar.jsonnet",
			content: "import 'demo.libsonnet' { db: importstr 'data/db.json' }",
			want: fileinfo.FileInfo{
				Path: fileinfo.FilePath{Package: "pkg/foo", Ext: ".jsonnet", Filename: "bar.jsonnet", Name: "bar", Path: "pkg/foo/bar.jsonnet"},
				Imports: map[string]fileinfo.FilePath{
					"pkg/foo/demo.libsonnet": {Package: "pkg/foo", Ext: ".libsonnet", Filename: "demo.libsonnet", Name: "demo", Path: "pkg/foo/demo.libsonnet"},
				},
				DataImports: map[string]fileinfo.FilePath{
					"pkg/foo/data/db.json": {Package: "pkg/foo/data", Ext: ".json", Filename: "db.json", Name: "db", Path: "pkg/foo/data/db.json"},
				},
			},
		}, {
			desc:    "json-like import",
			rel:     "pkg/foo",
			name:    "bar.jsonnet",
			content: "import 'data/db.json'",
			want: fileinfo.FileInfo{
				Path:    fileinfo.FilePath{Package: "pkg/foo", Ext: ".jsonnet", Filename: "bar.jsonnet", Name: "bar", Path: "pkg/foo/bar.jsonnet"},
				Imports: map[string]fileinfo.FilePath{},
				DataImports: map[string]fileinfo.FilePath{
					"pkg/foo/data/db.json": {Package: "pkg/foo/data", Ext: ".json", Filename: "db.json", Name: "db", Path: "pkg/foo/data/db.json"},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			dir, err := ioutil.TempDir("", "test")
			if err != nil {
				t.Fatal(err)
			}
			dir = filepath.Join(dir, tc.rel)
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(dir)

			if err := ioutil.WriteFile(filepath.Join(dir, tc.name), []byte(tc.content), 0600); err != nil {
				t.Fatal(err)
			}

			got := jsonnet.NewFileInfo(&config.Config{}, dir, tc.rel, tc.name)
			normalizeFileInfo(&got)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %#v; want %#v", got, tc.want)
			}
		})
	}
}
