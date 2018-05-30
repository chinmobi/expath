package expath

import (
	"testing"
)

func errp(e error) string {
	if e == nil {
		return "<nil>"
	}
	return e.Error()
}

type ScanSegmentsTest struct {
	pattern string
	segs    []patternSeg
	err     error
}

func isSegsEquals(s1, s2 []patternSeg) bool {
	len, len2 := len(s1), len(s2)
	if len == len2 {
		if len != 0 {
			for i := 0; i < len; i++ {
				if s1[i] != s2[i] {
					return false
				}
			}
		}
		return true
	}
	return false
}

var scanTests = []ScanSegmentsTest{
	{"**", []patternSeg{{"**", -1}}, nil},
	{"**/", []patternSeg{{"**/", -1}}, nil},
	{"/**/", []patternSeg{{"/**/", -1}}, nil},
	{"/**", []patternSeg{{"/**", -1}}, nil},

	{"abc", []patternSeg{{"abc", 1}}, nil},
	{"/abc", []patternSeg{{"/abc", 1}}, nil},
	{"/abc/", []patternSeg{{"/abc/", 1}}, nil},

	{"a/b/c", []patternSeg{{"a/b/c", 3}}, nil},
	{"/a/b/c", []patternSeg{{"/a/b/c", 3}}, nil},
	{"/a/b/c/", []patternSeg{{"/a/b/c/", 3}}, nil},

	{"**/a/b/c", []patternSeg{{"**/", -1}, {"a/b/c", 3}}, nil},
	{"**/a/b/c/", []patternSeg{{"**/", -1}, {"a/b/c/", 3}}, nil},
	{"/**/a/b/c", []patternSeg{{"/**/", -1}, {"a/b/c", 3}}, nil},
	{"/**/a/b/c/", []patternSeg{{"/**/", -1}, {"a/b/c/", 3}}, nil},
	{"a/b/c/**", []patternSeg{{"a/b/c/", 3}, {"**", -1}}, nil},
	{"a/b/c/**/", []patternSeg{{"a/b/c/", 3}, {"**/", -1}}, nil},
	{"/a/b/c/**", []patternSeg{{"/a/b/c/", 3}, {"**", -1}}, nil},
	{"/a/b/c/**/", []patternSeg{{"/a/b/c/", 3}, {"**/", -1}}, nil},

	{"**/a/b/c/**", []patternSeg{{"**/", -1}, {"a/b/c/", 3}, {"**", -1}}, nil},
	{"**/a/b/c/**/", []patternSeg{{"**/", -1}, {"a/b/c/", 3}, {"**/", -1}}, nil},
	{"/**/a/b/c/**", []patternSeg{{"/**/", -1}, {"a/b/c/", 3}, {"**", -1}}, nil},
	{"/**/a/b/c/**/", []patternSeg{{"/**/", -1}, {"a/b/c/", 3}, {"**/", -1}}, nil},

	{"a/b/c/**/d/e", []patternSeg{{"a/b/c/", 3}, {"**/", -1}, {"d/e", 2}}, nil},
	{"a/b/c/**/d/e/", []patternSeg{{"a/b/c/", 3}, {"**/", -1}, {"d/e/", 2}}, nil},

	{"a/b/c/**/d/e/**", []patternSeg{{"a/b/c/", 3}, {"**/", -1}, {"d/e/", 2}, {"**", -1}}, nil},

	{"a/b/c/**/d/e/**/f", []patternSeg{{"a/b/c/", 3}, {"**/", -1}, {"d/e/", 2}, {"**/", -1}, {"f", 1}}, nil},

	{"a/b/c/*/**/*/d/e/**/f", []patternSeg{{"a/b/c/*/", 4}, {"**/", -1}, {"*/d/e/", 3}, {"**/", -1}, {"f", 1}}, nil},

	{"*", []patternSeg{{"*", 1}}, nil},
	{"/*", []patternSeg{{"/*", 1}}, nil},
	{"/*/", []patternSeg{{"/*/", 1}}, nil},
	{"*/", []patternSeg{{"*/", 1}}, nil},

	{"/*/*", []patternSeg{{"/*/*", 2}}, nil},
	{"/*/*/", []patternSeg{{"/*/*/", 2}}, nil},
	{"*/*/", []patternSeg{{"*/*/", 2}}, nil},

	{"/*/**/*", []patternSeg{{"/*/", 1}, {"**/", -1}, {"*", 1}}, nil},
	{"/*/**/*/", []patternSeg{{"/*/", 1}, {"**/", -1}, {"*/", 1}}, nil},

	{"/", []patternSeg{{"/", 0}}, nil},

	{"", nil, nil},
}

func TestScanSegments(t *testing.T) {
	for _, tt := range scanTests {
		pattern := tt.pattern

		segs, err := scanSegments(pattern)
		if err != nil || !isSegsEquals(tt.segs, segs) {
			t.Errorf("scanSegments(%#q) = %v, %q want %v, %q", pattern, segs, errp(err), tt.segs, errp(tt.err))
		}
	}
}
