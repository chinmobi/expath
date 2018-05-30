package expath

import (
	"runtime"
)

// A patternSeg is a pattern segment of the whole pattern.
// Each of the segment is either the normal pattern (none any-dirs' term) or just the any-dirs' pattern.
// The dirs indicates how many dirs that the pattern has (seperated by Seperator), -1 for any-dirs' pattern.
//
type patternSeg struct {
	pattern string
	dirs    int
}

// scanSegments scans the whole pattern and separates it into segments by the any-dirs' term ('**').
//
func scanSegments(pattern string) (segs []patternSeg, err error) {
	len := len(pattern)
	if len == 0 {
		return
	}

	var i, from, to, dirs int
	var ok bool

	if pattern[0] == '*' {
		if to, ok = scanAnyDirsPattern(pattern, 0, len, false); ok {
			segs = append(segs, patternSeg{pattern[from:to], -1})
		}
		from, i = to, to
	}

	for i < len {
		switch pattern[i] {
		case '\\':
			if runtime.GOOS != "windows" {
				i++
				continue
			}
			fallthrough
		case '/':
			if i > 0 {
				dirs++
			}

			to, ok = scanAnyDirsPattern(pattern, i+1, len, false)
			if ok || to >= len {
				if from < i {
					segs = append(segs, patternSeg{pattern[from : i+1], dirs})
					from = i + 1
				}

				dirs = 0

				switch to - from {
				case 0: // Nothing to do
				case 1: // Just "/"
					segs = append(segs, patternSeg{pattern[from:to], 0})
				default:
					segs = append(segs, patternSeg{pattern[from:to], -1})
				}

				from, i = to, to
			} else {
				i++
			}

		default:
			i++
		}
	}

	if from < i {
		segs = append(segs, patternSeg{pattern[from:i], dirs + 1})
	}

	return
}

// scanAnyDirsPattern used to identify the any-dirs' term ('**').
//
func scanAnyDirsPattern(pattern string, from, len int, preOK bool) (int, bool) {
	i := from
	if i < len && pattern[i] == '*' {
		i++
		if i < len && pattern[i] == '*' {
			for i++; i < len; {
				switch pattern[i] {
				case '*': // more star?
					i++
				case '\\':
					if runtime.GOOS != "windows" {
						return from, preOK
					}
					fallthrough
				case '/':
					return scanAnyDirsPattern(pattern, i+1, len, true)
				default:
					return from, preOK
				}
			}
			return i, true
		}
	}
	return from, preOK
}
