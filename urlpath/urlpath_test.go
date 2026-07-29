package urlpath_test

import (
	"testing"

	"github.com/notwithering/shit/urlpath"
)

func TestClean(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		// separators
		{"", "/"},
		{"/", "/"},
		{"//", "/"},
		{"///", "/"},
		{"a//", "/a/"},
		{"//a", "/a"},
		{"//a//", "/a/"},

		// traversal
		{"a/./b", "/a/b"},
		{"a/../b", "/b"},
		{".", "/"},
		{"..", "/"},

		// files
		{"a", "/a"},
		{"/a", "/a"},
		{"a/b", "/a/b"},
		{"/a/b", "/a/b"},

		// dirs
		{"a/", "/a/"},
		{"a/b/", "/a/b/"},
		{"/a/", "/a/"},
		{"/a/b/", "/a/b/"},
	}

	for _, test := range tests {
		result := urlpath.Clean(test.path)
		if result != test.expected {
			t.Errorf("Clean(%q) = %q, expected %q", test.path, result, test.expected)
		}
	}
}

func TestJoin(t *testing.T) {
	tests := []struct {
		path     []string
		expected string
	}{
		{[]string{}, "/"},
		{[]string{"a"}, "/a"},
		{[]string{"a", "b"}, "/a/b"},
		{[]string{"a/b", "c"}, "/a/b/c"},
		{[]string{"a/b", "c/d"}, "/a/b/c/d"},
	}

	for _, test := range tests {
		result := urlpath.Join(test.path...)
		if result != test.expected {
			t.Errorf("Join(%v) = %q, expected %q", test.path, result, test.expected)
		}
	}
}
