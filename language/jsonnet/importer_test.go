package jsonnet_test

import (
	"reflect"
	"testing"

	"github.com/bitnami/jsonnet-gazelle/language/jsonnet"
	gojsonnet "github.com/google/go-jsonnet"
)

func TestParseSnippetImports(t *testing.T) {
	testCases := []struct {
		desc    string
		snippet string
		want    []string
	}{
		{
			desc:    "empty",
			snippet: "{}",
			want:    nil,
		},
		{
			desc:    "simple import",
			snippet: "import 'a.jsonnet'",
			want:    []string{"a.jsonnet"},
		},
		{
			desc:    "arbitrary import",
			snippet: "(import 'a.f.o.o')",
			want:    []string{"a.f.o.o"},
		},
		{
			desc:    "consecutive import",
			snippet: "(import 'a.jsonnet') + (import 'b.jsonnet')",
			want:    []string{"a.jsonnet", "b.jsonnet"},
		},
		{
			desc:    "repeated import",
			snippet: "(import 'a.jsonnet') + (import 'a.jsonnet')",
			want:    []string{"a.jsonnet"},
		},
		{
			desc:    "parent import",
			snippet: "(import '../a.jsonnet')",
			want:    []string{"../a.jsonnet"},
		},
		{
			desc:    "subfolder import",
			snippet: "(import 'b/a.jsonnet')",
			want:    []string{"b/a.jsonnet"},
		},
		{
			desc:    "simple importstr",
			snippet: "(importstr 'a.json')",
			want:    []string{"a.json"},
		},
	}

	filename := "test.jsonnet"
	importer := &jsonnet.Importer{&gojsonnet.FileImporter{}}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := importer.ParseSnippetImports(filename, tc.snippet)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %q; want %q", got, tc.want)
			}
		})
	}
}
