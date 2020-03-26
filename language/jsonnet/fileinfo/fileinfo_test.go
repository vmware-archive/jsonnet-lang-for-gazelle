// Copyright 2019 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package fileinfo_test

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/vmware/jsonnet-lang-for-gazelle/language/jsonnet/fileinfo"
)

func TestNewFilePath(t *testing.T) {
	testCases := []struct {
		root string
		elem []string
		want fileinfo.FilePath
	}{
		{"/a", []string{"b/c", "d.txt"}, fileinfo.FilePath{Root: "/a", Package: "b/c", Ext: ".txt", Filename: "d.txt", Name: "d", Path: "b/c/d.txt"}},
		{"/a", []string{"b/c/d.txt"}, fileinfo.FilePath{Root: "/a", Package: "b/c", Ext: ".txt", Filename: "d.txt", Name: "d", Path: "b/c/d.txt"}},
		{"/a", []string{"/a/b/c/d.txt"}, fileinfo.FilePath{Root: "/a", Package: "b/c", Ext: ".txt", Filename: "d.txt", Name: "d", Path: "b/c/d.txt"}},
		// This is a valid case for the method, but we should not allow imports outside the root code base
		{"/a", []string{"/b/c/d.txt"}, fileinfo.FilePath{Root: "/a", Package: "../b/c", Ext: ".txt", Filename: "d.txt", Name: "d", Path: "../b/c/d.txt"}},
	}

	for _, tc := range testCases {
		t.Run(filepath.Join(tc.elem...), func(t *testing.T) {
			got, _ := fileinfo.NewFilePath(tc.root, tc.elem...)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got: %#v; want: %#v", got, tc.want)
			}
		})
	}
}

func TestRuleName(t *testing.T) {
	testCases := []struct {
		path   fileinfo.FilePath // we only need the Name attribute
		prefix string
		want   string
	}{
		{fileinfo.FilePath{Name: "foo"}, "rule_test", "foo_rule_test"},
		{fileinfo.FilePath{Name: "foo1234"}, "rule_test", "foo1234_rule_test"},
		{fileinfo.FilePath{Name: "FOO1234"}, "rule_test", "foo1234_rule_test"},
		{fileinfo.FilePath{Name: "foo_1234"}, "rule_test", "foo_1234_rule_test"},
		{fileinfo.FilePath{Name: "foo.12.34"}, "rule_test", "foo_12_34_rule_test"},
	}

	for _, tc := range testCases {
		t.Run(tc.want, func(t *testing.T) {
			if got := tc.path.RuleName(tc.prefix); got != tc.want {
				t.Errorf("got: %q; want: %q", got, tc.want)
			}
		})
	}
}

func TestNewDataRef(t *testing.T) {
	testCases := []struct {
		path fileinfo.FilePath // we only need the Path attribute
		want string
	}{
		{fileinfo.FilePath{Path: "a/b/foo.ext"}, "//:a/b/foo.ext"},
		{fileinfo.FilePath{Path: "a/b/foo1234.ext"}, "//:a/b/foo1234.ext"},
		{fileinfo.FilePath{Path: "a/b/FOO1234.ext"}, "//:a/b/FOO1234.ext"},
		{fileinfo.FilePath{Path: "a/b/foo_1234.ext"}, "//:a/b/foo_1234.ext"},
		{fileinfo.FilePath{Path: "a/b/foo.12.34.ext"}, "//:a/b/foo.12.34.ext"},
	}

	for _, tc := range testCases {
		t.Run(tc.want, func(t *testing.T) {
			if got := tc.path.NewDataRef(); got != tc.want {
				t.Errorf("got: %q; want: %q", got, tc.want)
			}
		})
	}
}

func TestNewLabel(t *testing.T) {
	testCases := []struct {
		path   fileinfo.FilePath // we only need the Package attribute
		prefix string
		want   label.Label
	}{
		{fileinfo.FilePath{Package: "a/b"}, "library", label.Label{Repo: "", Pkg: "a/b", Name: "_library", Relative: false}},
	}

	for _, tc := range testCases {
		t.Run(tc.prefix, func(t *testing.T) {
			got := tc.path.NewLabel(tc.prefix)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got: %#v; want: %#v", got, tc.want)
			}
		})
	}
}

func TestNewDataLabel(t *testing.T) {
	testCases := []struct {
		path fileinfo.FilePath // we only need the Package and Filename attributes
		want string
	}{
		{fileinfo.FilePath{Package: "a/b", Filename: "foo.ext"}, "//a/b:foo.ext"},
		{fileinfo.FilePath{Package: "a/b", Filename: "foo1234.ext"}, "//a/b:foo1234.ext"},
		{fileinfo.FilePath{Package: "a/b", Filename: "FOO1234.ext"}, "//a/b:FOO1234.ext"},
		{fileinfo.FilePath{Package: "a/b", Filename: "foo_1234.ext"}, "//a/b:foo_1234.ext"},
		{fileinfo.FilePath{Package: "a/b", Filename: "foo.12.34.ext"}, "//a/b:foo.12.34.ext"},
	}

	for _, tc := range testCases {
		t.Run(tc.want, func(t *testing.T) {
			if got := tc.path.NewDataLabel(); got != tc.want {
				t.Errorf("got: %q; want: %q", got, tc.want)
			}
		})
	}
}
