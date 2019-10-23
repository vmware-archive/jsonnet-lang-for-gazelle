package jsonnet

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/config"
)

func TestJsonnetFileInfo(t *testing.T) {
	testCases := []struct {
		desc, dir, rel, name, content string
		want                          FileInfo
	}{
		{
			desc:    "empty",
			rel:     "pkg/foo",
			name:    "bar.jsonnet",
			content: "",
			want: FileInfo{
				Path:        FilePath{Dir: "pkg/foo", Ext: ".jsonnet", Filename: "bar.jsonnet", Name: "bar", Path: "pkg/foo/bar.jsonnet"},
				Imports:     map[string]FilePath{},
				DataImports: map[string]FilePath{},
			},
		}, {
			desc:    "differnt quotes imports",
			rel:     "pkg/foo",
			name:    "bar.jsonnet",
			content: `(import 'singlequotes.jsonnet') + (import "doublequotes.jsonnet")`,
			want: FileInfo{
				Path: FilePath{Dir: "pkg/foo", Ext: ".jsonnet", Filename: "bar.jsonnet", Name: "bar", Path: "pkg/foo/bar.jsonnet"},
				Imports: map[string]FilePath{
					"pkg/foo/singlequotes.jsonnet": {Dir: "pkg/foo", Ext: ".jsonnet", Filename: "singlequotes.jsonnet", Name: "singlequotes", Path: "pkg/foo/singlequotes.jsonnet"},
					"pkg/foo/doublequotes.jsonnet": {Dir: "pkg/foo", Ext: ".jsonnet", Filename: "doublequotes.jsonnet", Name: "doublequotes", Path: "pkg/foo/doublequotes.jsonnet"},
				},
				DataImports: map[string]FilePath{},
			},
		}, {
			desc:    "libsonnet import",
			rel:     "pkg/foo",
			name:    "bar.jsonnet",
			content: "import 'demo.libsonnet'",
			want: FileInfo{
				Path: FilePath{Dir: "pkg/foo", Ext: ".jsonnet", Filename: "bar.jsonnet", Name: "bar", Path: "pkg/foo/bar.jsonnet"},
				Imports: map[string]FilePath{
					"pkg/foo/demo.libsonnet": {Dir: "pkg/foo", Ext: ".libsonnet", Filename: "demo.libsonnet", Name: "demo", Path: "pkg/foo/demo.libsonnet"},
				},
				DataImports: map[string]FilePath{},
			},
		}, {
			desc:    "different folder imports",
			rel:     "pkg/foo",
			name:    "bar.jsonnet",
			content: "(import '../pkg.libsonnet') + (import '../../root.jsonnet')",
			want: FileInfo{
				Path: FilePath{Dir: "pkg/foo", Ext: ".jsonnet", Filename: "bar.jsonnet", Name: "bar", Path: "pkg/foo/bar.jsonnet"},
				Imports: map[string]FilePath{
					"pkg/pkg.libsonnet": {Dir: "pkg", Ext: ".libsonnet", Filename: "pkg.libsonnet", Name: "pkg", Path: "pkg/pkg.libsonnet"},
					"root.jsonnet":      {Dir: "", Ext: ".jsonnet", Filename: "root.jsonnet", Name: "root", Path: "root.jsonnet"},
				},
				DataImports: map[string]FilePath{},
			},
		}, {
			desc:    "data import",
			rel:     "pkg/foo",
			name:    "bar.jsonnet",
			content: "importstr 'data/db.json'",
			want: FileInfo{
				Path:    FilePath{Dir: "pkg/foo", Ext: ".jsonnet", Filename: "bar.jsonnet", Name: "bar", Path: "pkg/foo/bar.jsonnet"},
				Imports: map[string]FilePath{},
				DataImports: map[string]FilePath{
					"pkg/foo/data/db.json": {Dir: "pkg/foo/data", Ext: ".json", Filename: "db.json", Name: "db", Path: "pkg/foo/data/db.json"},
				},
			},
		}, {
			desc:    "mixed data and jsonnet imports",
			rel:     "pkg/foo",
			name:    "bar.jsonnet",
			content: "import 'demo.libsonnet' { db: importstr 'data/db.json' }",
			want: FileInfo{
				Path: FilePath{Dir: "pkg/foo", Ext: ".jsonnet", Filename: "bar.jsonnet", Name: "bar", Path: "pkg/foo/bar.jsonnet"},
				Imports: map[string]FilePath{
					"pkg/foo/demo.libsonnet": {Dir: "pkg/foo", Ext: ".libsonnet", Filename: "demo.libsonnet", Name: "demo", Path: "pkg/foo/demo.libsonnet"},
				},
				DataImports: map[string]FilePath{
					"pkg/foo/data/db.json": {Dir: "pkg/foo/data", Ext: ".json", Filename: "db.json", Name: "db", Path: "pkg/foo/data/db.json"},
				},
			},
		}, {
			desc:    "json-like import",
			rel:     "pkg/foo",
			name:    "bar.jsonnet",
			content: "import 'data/db.json'",
			want: FileInfo{
				Path:    FilePath{Dir: "pkg/foo", Ext: ".jsonnet", Filename: "bar.jsonnet", Name: "bar", Path: "pkg/foo/bar.jsonnet"},
				Imports: map[string]FilePath{},
				DataImports: map[string]FilePath{
					"pkg/foo/data/db.json": {Dir: "pkg/foo/data", Ext: ".json", Filename: "db.json", Name: "db", Path: "pkg/foo/data/db.json"},
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

			got := jsonnetFileInfo(&config.Config{}, dir, tc.rel, tc.name)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %#v; want %#v", got, tc.want)
			}
		})
	}
}
