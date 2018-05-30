package expath

import (
	"testing"
)

type MatchTest struct {
	pattern, s string
	matched    bool
	err        error
}

var matchTests = []MatchTest{
	{"a/b/c", "a/b/c", true, nil}, // --- 0
	{"/a/b/c", "a/b/c", false, nil},
	{"/a/b/c/", "a/b/c", false, nil},
	{"a/b/c/", "a/b/c", false, nil},
	{"**", "a/b/c", true, nil}, // --- 1
	{"/**", "a/b/c", false, nil},
	{"/**/", "a/b/c", false, nil},
	{"**/", "a/b/c", false, nil},
	{"abc/**", "abc", true, nil}, // --- 2
	{"/abc/**", "abc", false, nil},
	{"/abc/**/", "abc", false, nil},
	{"abc/**/", "abc", false, nil},
	{"abc/**", "abc/a", true, nil}, // ---
	{"/abc/**", "abc/a", false, nil},
	{"/abc/**/", "abc/a", false, nil},
	{"abc/**/", "abc/a", false, nil},
	{"**/abc", "abc", true, nil}, // --- 3
	{"/**/abc", "abc", false, nil},
	{"/**/abc/", "abc", false, nil},
	{"**/abc/", "abc", false, nil},
	{"**/abc", "a/abc", true, nil}, // ---
	{"/**/abc", "a/abc", false, nil},
	{"/**/abc/", "a/abc", false, nil},
	{"**/abc/", "a/abc", false, nil},
	{"**/abc/**", "abc", true, nil}, // --- 4
	{"/**/abc/**", "abc", false, nil},
	{"/**/abc/**/", "abc", false, nil},
	{"**/abc/**/", "abc", false, nil},
	{"**/abc/**", "a/abc/c", true, nil}, // ---
	{"/**/abc/**", "a/abc/c", false, nil},
	{"/**/abc/**/", "a/abc/c", false, nil},
	{"**/abc/**/", "a/abc/c", false, nil},

	{"a/b/c", "/a/b/c", false, nil}, // --- 0
	{"/a/b/c", "/a/b/c", true, nil},
	{"/a/b/c/", "/a/b/c", false, nil},
	{"a/b/c/", "/a/b/c", false, nil},
	{"**", "/a/b/c", true, nil}, // --- 1
	{"/**", "/a/b/c", true, nil},
	{"/**/", "/a/b/c", false, nil},
	{"**/", "/a/b/c", false, nil},
	{"abc/**", "/abc", false, nil}, // --- 2
	{"/abc/**", "/abc", true, nil},
	{"/abc/**/", "/abc", false, nil},
	{"abc/**/", "/abc", false, nil},
	{"abc/**", "/abc/a", false, nil}, // ---
	{"/abc/**", "/abc/a", true, nil},
	{"/abc/**/", "/abc/a", false, nil},
	{"abc/**/", "/abc/a", false, nil},
	{"**/abc", "/abc", true, nil}, // --- 3
	{"/**/abc", "/abc", true, nil},
	{"/**/abc/", "/abc", false, nil},
	{"**/abc/", "/abc", false, nil},
	{"**/abc", "/a/abc", true, nil}, // ---
	{"/**/abc", "/a/abc", true, nil},
	{"/**/abc/", "/a/abc", false, nil},
	{"**/abc/", "/a/abc", false, nil},
	{"**/abc/**", "/abc", true, nil}, // --- 4
	{"/**/abc/**", "/abc", true, nil},
	{"/**/abc/**/", "/abc", false, nil},
	{"**/abc/**/", "/abc", false, nil},
	{"**/abc/**", "/a/abc/c", true, nil}, // ---
	{"/**/abc/**", "/a/abc/c", true, nil},
	{"/**/abc/**/", "/a/abc/c", false, nil},
	{"**/abc/**/", "/a/abc/c", false, nil},

	{"a/b/c", "/a/b/c/", false, nil}, // --- 0
	{"/a/b/c", "/a/b/c/", false, nil},
	{"/a/b/c/", "/a/b/c/", true, nil},
	{"a/b/c/", "/a/b/c/", false, nil},
	{"**", "/a/b/c/", true, nil}, // --- 1
	{"/**", "/a/b/c/", true, nil},
	{"/**/", "/a/b/c/", true, nil},
	{"**/", "/a/b/c/", true, nil},
	{"abc/**", "/abc/", false, nil}, // --- 2
	{"/abc/**", "/abc/", true, nil},
	{"/abc/**/", "/abc/", true, nil},
	{"abc/**/", "/abc/", false, nil},
	{"abc/**", "/abc/a/", false, nil}, // ---
	{"/abc/**", "/abc/a/", true, nil},
	{"/abc/**/", "/abc/a/", true, nil},
	{"abc/**/", "/abc/a/", false, nil},
	{"**/abc", "/abc/", false, nil}, // --- 3
	{"/**/abc", "/abc/", false, nil},
	{"/**/abc/", "/abc/", true, nil},
	{"**/abc/", "/abc/", true, nil},
	{"**/abc", "/a/abc/", false, nil}, // ---
	{"/**/abc", "/a/abc/", false, nil},
	{"/**/abc/", "/a/abc/", true, nil},
	{"**/abc/", "/a/abc/", true, nil},
	{"**/abc/**", "/abc/", true, nil}, // --- 4
	{"/**/abc/**", "/abc/", true, nil},
	{"/**/abc/**/", "/abc/", true, nil},
	{"**/abc/**/", "/abc/", true, nil},
	{"**/abc/**", "/a/abc/c/", true, nil}, // ---
	{"/**/abc/**", "/a/abc/c/", true, nil},
	{"/**/abc/**/", "/a/abc/c/", true, nil},
	{"**/abc/**/", "/a/abc/c/", true, nil},

	{"a/b/c", "a/b/c/", false, nil}, // --- 0
	{"/a/b/c", "a/b/c/", false, nil},
	{"/a/b/c/", "a/b/c/", false, nil},
	{"a/b/c/", "a/b/c/", true, nil},
	{"**", "a/b/c/", true, nil}, // --- 1
	{"/**", "a/b/c/", false, nil},
	{"/**/", "a/b/c/", false, nil},
	{"**/", "a/b/c/", true, nil},
	{"abc/**", "abc/", true, nil}, // --- 2
	{"/abc/**", "abc/", false, nil},
	{"/abc/**/", "abc/", false, nil},
	{"abc/**/", "abc/", true, nil},
	{"abc/**", "abc/a/", true, nil}, // ---
	{"/abc/**", "abc/a/", false, nil},
	{"/abc/**/", "abc/a/", false, nil},
	{"abc/**/", "abc/a/", true, nil},
	{"**/abc", "abc/", false, nil}, // --- 3
	{"/**/abc", "abc/", false, nil},
	{"/**/abc/", "abc/", false, nil},
	{"**/abc/", "abc/", true, nil},
	{"**/abc", "a/abc/", false, nil}, // ---
	{"/**/abc", "a/abc/", false, nil},
	{"/**/abc/", "a/abc/", false, nil},
	{"**/abc/", "a/abc/", true, nil},
	{"**/abc/**", "abc/", true, nil}, // --- 4
	{"/**/abc/**", "abc/", false, nil},
	{"/**/abc/**/", "abc/", false, nil},
	{"**/abc/**/", "abc/", true, nil},
	{"**/abc/**", "a/abc/c/", true, nil}, // ---
	{"/**/abc/**", "a/abc/c/", false, nil},
	{"/**/abc/**/", "a/abc/c/", false, nil},
	{"**/abc/**/", "a/abc/c/", true, nil},

	{"**/abc", "a/b/abc", true, nil}, // ---
	{"/**/abc", "a/b/abc", false, nil},
	{"/**/abc/", "a/b/abc", false, nil},
	{"**/abc/", "a/b/abc", false, nil},
	{"**/abc", "a/b/c/abc", true, nil}, // ---
	{"/**/abc", "a/b/c/abc", false, nil},
	{"/**/abc/", "a/b/c/abc", false, nil},
	{"**/abc/", "a/b/c/abc", false, nil},

	{"**/abc/**/def/**", "a/abc/c/d/def/f", true, nil}, // ---
	{"/**/abc/**/def/**", "a/abc/c/d/def/f", false, nil},
	{"/**/abc/**/def/**/", "a/abc/c/d/def/f", false, nil},
	{"**/abc/**/def/**/", "a/abc/c/d/def/f", false, nil},
	{"**/abc/**/def/**", "/a/abc/c/d/def/f", true, nil}, // ---
	{"/**/abc/**/def/**", "/a/abc/c/d/def/f", true, nil},
	{"/**/abc/**/def/**/", "/a/abc/c/d/def/f", false, nil},
	{"**/abc/**/def/**/", "/a/abc/c/d/def/f", false, nil},
	{"**/abc/**/def/**", "/a/abc/c/d/def/f/", true, nil}, // ---
	{"/**/abc/**/def/**", "/a/abc/c/d/def/f/", true, nil},
	{"/**/abc/**/def/**/", "/a/abc/c/d/def/f/", true, nil},
	{"**/abc/**/def/**/", "/a/abc/c/d/def/f/", true, nil},
	{"**/abc/**/def/**", "a/abc/c/d/def/f/", true, nil}, // ---
	{"/**/abc/**/def/**", "a/abc/c/d/def/f/", false, nil},
	{"/**/abc/**/def/**/", "a/abc/c/d/def/f/", false, nil},
	{"**/abc/**/def/**/", "a/abc/c/d/def/f/", true, nil},
}

func TestMatch(t *testing.T) {
	for _, tt := range matchTests {
		pattern := tt.pattern
		s := tt.s

		ok, err := Match(pattern, s)
		if ok != tt.matched || err != tt.err {
			t.Errorf("Match(%#q, %#q) = %v, %q want %v, %q", pattern, s, ok, errp(err), tt.matched, errp(tt.err))
		}
	}
}

var matchTests0 = []MatchTest{
	{"**/abc", "/", false, nil}, // ---
	{"/**/abc", "/", false, nil},
	{"/**/abc/", "/", false, nil},
	{"**/abc/", "/", false, nil},
	{"abc/**", "/", false, nil}, // ---
	{"/abc/**", "/", false, nil},
	{"/abc/**/", "/", false, nil},
	{"abc/**/", "/", false, nil},

	{"**/abc", "", false, nil}, // ---
	{"/**/abc", "", false, nil},
	{"/**/abc/", "", false, nil},
	{"**/abc/", "", false, nil},
	{"abc/**", "", false, nil}, // ---
	{"/abc/**", "", false, nil},
	{"/abc/**/", "", false, nil},
	{"abc/**/", "", false, nil},
}

func TestMatch0(t *testing.T) {
	for _, tt := range matchTests0 {
		pattern := tt.pattern
		s := tt.s

		ok, err := Match(pattern, s)
		if ok != tt.matched || err != tt.err {
			t.Errorf("Match(%#q, %#q) = %v, %q want %v, %q", pattern, s, ok, errp(err), tt.matched, errp(tt.err))
		}
	}
}

func TestMatchAnyDirs(t *testing.T) {
	matchTests := []MatchTest{
		{"**", "abc", true, nil},
		{"**", "abc/", true, nil},
		{"**", "/abc", true, nil},
		{"**", "/abc/", true, nil},
		{"**", "", true, nil},
		{"**", "/", true, nil},

		{"/**", "abc", false, nil},
		{"/**", "abc/", false, nil},
		{"/**", "/abc", true, nil},
		{"/**", "/abc/", true, nil},
		{"/**", "", false, nil},
		{"/**", "/", true, nil},

		{"/**/", "abc", false, nil},
		{"/**/", "abc/", false, nil},
		{"/**/", "/abc", false, nil},
		{"/**/", "/abc/", true, nil},
		{"/**/", "", false, nil},
		{"/**/", "/", true, nil},

		{"**/", "abc", false, nil},
		{"**/", "abc/", true, nil},
		{"**/", "/abc", false, nil},
		{"**/", "/abc/", true, nil},
		{"**/", "", false, nil},
		{"**/", "/", true, nil},
	}

	for _, tt := range matchTests {
		pattern := tt.pattern
		s := tt.s

		ok, err := matchAnyDirs(pattern, s)
		if ok != tt.matched || err != tt.err {
			t.Errorf("matchAnyDirs(%#q, %#q) = %v, %q want %v, %q", pattern, s, ok, errp(err), tt.matched, errp(tt.err))
		}
	}
}
