package expath

import (
	"errors"
	"path/filepath"
	"runtime"
	"testing"
)

type testPathHelper struct {
	testPath string
}

func (t testPathHelper) getNames(dir string) (names []string, err error) {
	testPath := trimPath(t.testPath)
	dir = trimDir(filepath.ToSlash(dir))

	nLen := len(dir)
	if nLen < len(testPath) {
		if nLen == 0 {
			// Nothing to do.
		} else if testPath[nLen] == '/' && dir == testPath[:nLen] {
			testPath = testPath[nLen+1:]
		} else {
			return
		}

		to, _ := scanDirs(testPath, 0, len(testPath), 1)
		if to > 0 {
			name := trimPath(testPath[:to])
			names = append(names, name)
		}
	}

	return
}

func (t testPathHelper) isExist(dir string) (bool, error) {
	testPath := trimPath(t.testPath)
	dir = trimDir(filepath.ToSlash(dir))

	nLen := len(dir)
	if nLen <= len(testPath) {
		if dir == testPath[:nLen] {
			if nLen == len(testPath) || testPath[nLen] == '/' {
				return true, nil
			}
		}
	}

	return false, nil
}

func trimDir(dir string) string {
	mark := skipDotsDir(dir, len(dir))
	if mark > 0 {
		dir = dir[mark:]
	}
	return trimPath(dir)
}

func wrapGlob(pattern, path string) (bool, error) {
	var helper testPathHelper
	helper.testPath = path

	var mh matchedSet

	err := doGlob(pattern, "", &helper, &mh)
	if err != nil {
		return false, err
	}

	for _, mp := range mh.matches {
		if mp == path {
			return true, nil
		} else {
			if trimPath(mp) == trimPath(path) {
				return true, nil
			}

			return false, errors.New(mp)
		}
	}

	return false, nil
}

func TestGlob(t *testing.T) {
	for _, tt := range matchTests {
		pattern := tt.pattern
		s := tt.s

		ok, err := wrapGlob(pattern, s)
		if !ok || err != tt.err {
			t.Errorf("Match(%#q, %#q) = %v, %q want %v, %q", pattern, s, ok, errp(err), true, errp(tt.err))
		}
	}
}

func TestGlob0(t *testing.T) {
	for _, tt := range matchTests0 {
		pattern := tt.pattern
		s := tt.s

		ok, err := wrapGlob(pattern, s)
		if ok != tt.matched || err != tt.err {
			t.Errorf("Match(%#q, %#q) = %v, %q want %v, %q", pattern, s, ok, errp(err), tt.matched, errp(tt.err))
		}
	}
}

type normalizePathTest struct {
	pattern, root     string
	toPattern, toRoot string
}

func TestNormalizePath(t *testing.T) {
	tests := []normalizePathTest{
		{"/", "", "/**", "."},
		{"./", "", "/**", "."},
		{"../", "", "/**", ".."},

		{"abc", "", "abc", "./"},
		{"abc/", "", "abc/**", "./"},
		{"/abc/", "", "/abc/**", "."},
		{"/abc", "", "/abc", "."},

		{"./abc", "", "/abc", "."},
		{"../abc", "", "/abc", ".."},
		{"./abc/", "", "/abc/**", "."},
		{"../abc/", "", "/abc/**", ".."},

		{"./abc", "d", "/abc", "d"},
		{"../abc", "d", "/abc", "."},
		{"../abc", "d/e", "/abc", "d"},
		{"./abc/", "d", "/abc/**", "d"},
		{"../abc/", "d", "/abc/**", "."},

		{"abc", "d", "abc", "d/"},
		{"abc/", "d", "abc/**", "d/"},
		{"/abc/", "d", "/abc/**", "d"},
		{"/abc", "d", "/abc", "d"},

		{"./abc", "/d", "/abc", "/d"},
		{"../abc", "/d", "/abc", "/"},
		{"../abc", "/d/e", "/abc", "/d"},
		{"./abc/", "/d", "/abc/**", "/d"},
		{"../abc/", "/d", "/abc/**", "/"},

		{"abc", "/d", "abc", "/d/"},
		{"abc/", "/d", "abc/**", "/d/"},
		{"/abc/", "/d", "/abc/**", "/d"},
		{"/abc", "/d", "/abc", "/d"},

		{"", "abc", "", "abc"},
		{"", "/abc", "", "/abc"},
		{"", "/abc/", "", "/abc"},
		{"", "abc/", "", "abc"},

		{"", "/", "", "/"},
		{"", "./", "", "."},
		{"", "../", "", ".."},

		{"", "", "", "."},
	}

	wintests := []normalizePathTest{
		{"abc", "c:", "abc", `c:.\`},
		{"/abc", "c:", "/abc", `c:.`},
		{"abc", `c:\`, "abc", `c:\`},
		{"/abc", `c:\`, "/abc", `c:\`},

		{"abc", `c:\d`, "abc", `c:\d\`},
		{"/abc", `c:\d`, "/abc", `c:\d`},

		{"./abc", "c:", "/abc", `c:.`},
		{"../abc", "c:", "/abc", `c:..`},
		{"./abc/", "c:", "/abc/**", `c:.`},
		{"../abc/", "c:", "/abc/**", `c:..`},

		{"./abc", `c:\`, "/abc", `c:\`},
		{"../abc", `c:\`, "/abc", `c:\`},
		{"./abc/", `c:\`, "/abc/**", `c:\`},
		{"../abc/", `c:\`, "/abc/**", `c:\`},

		{"./abc", `c:\d`, "/abc", `c:\d`},
		{"../abc", `c:\d`, "/abc", `c:\`},
		{"./abc/", `c:\d`, "/abc/**", `c:\d`},
		{"../abc/", `c:\d`, "/abc/**", `c:\`},

		{"c:", "", "", `c:.`},
		{`c:\`, "", "/**", `c:\`},
		{`c:\abc`, "", "/abc", `c:\`},
		{`c:\abc\`, "", "/abc/**", `c:\`},

		{"", "c:", "", `c:.`},
		{"", `c:\`, "", `c:\`},
		{"", `c:\abc`, "", `c:\abc`},
		{"", `c:\abc\`, "", `c:\abc`},
	}

	if runtime.GOOS == "windows" {
		for i := range tests {
			tests[i].toRoot = filepath.FromSlash(tests[i].toRoot)
		}

		tests = append(tests, wintests...)
	}

	for _, tt := range tests {
		toPattern, toRoot := normalizePath(tt.pattern, tt.root)
		if toPattern != tt.toPattern || toRoot != tt.toRoot {
			t.Errorf("normalizePath(%#q, %#q) = (%#q, %#q) want (%#q, %#q)", tt.pattern, tt.root, toPattern, toRoot, tt.toPattern, tt.toRoot)
		}
	}

}
