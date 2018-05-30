package expath

import (
	"path/filepath"
	"runtime"
)

// matchASeg used to optimize the match routine if there is only one segment of the whole pattern,
// the whole pattern is either normal pattern or just any-dirs' pattern.
//
func matchASeg(seg patternSeg, name string) (bool, error) {
	if seg.dirs >= 0 {
		return filepath.Match(seg.pattern, name)
	}
	return matchAnyDirs(seg.pattern, name)
}

func matchAnyDirs(pattern, name string) (matched bool, err error) {
	a, z := pattern[0], pattern[len(pattern)-1]

	if a == '*' && z == '*' { // "**"
		return true, nil
	}

	nLen := len(name)
	if nLen <= 0 {
		return false, nil
	}

	if a == '*' { // "**/"
		if isDirSeparator(name, nLen-1) { // ".../"
			return true, nil
		}
	} else if z == '*' { // "/**"
		if isDirSeparator(name, 0) { // "/..."
			return true, nil
		}
	} else { // "/**/"
		if isDirSeparator(name, 0) &&
			isDirSeparator(name, nLen-1) { // "/.../"
			return true, nil
		}
	}

	return false, nil
}

func isDirSeparator(name string, i int) bool {
	switch name[i] {
	case '\\':
		if runtime.GOOS == "windows" {
			return true
		}
	case '/':
		return true
	}
	return false
}

// matchSegs is main routine to match each pattern segment.
//
func matchSegs(segs []patternSeg, name string) (matched bool, err error) {
	nLen := len(name)
	if nLen <= 0 {
		return
	}

	var i, from, to int
	var mark int

	segsLen := len(segs)
	// assert segsLen > 1

	seg := segs[0]
	pattern := seg.pattern

	if seg.dirs < 0 {
		if pattern[0] == '*' {
			if isDirSeparator(name, 0) {
				from++
			}
		} else {
			if isDirSeparator(name, 0) {
				from++
			} else {
				return
			}
		}
	} else {
		if pattern[0] != name[0] {
			return
		}

		if to, mark = scanDirs(name, from, nLen, seg.dirs); mark < 0 {
			return
		}

		if mark == 0 && isSegLastAndAny(segs, 1, segsLen) {
			pattern = pattern[:len(pattern)-1]
		}
		if matched, err = filepath.Match(pattern, name[from:to]); !matched {
			return
		}

		from = to

		seg = segs[1]
		// assert seg.dirs < 0

		i++
	}

	i++
	for {
		if i < segsLen {
			to, matched, err = searchMatched(segs, i, segsLen, name, from, nLen)
			if !matched {
				return
			}

			from = to
			i++

			if i >= segsLen {
				if from >= nLen {
					return true, nil
				}
				break
			}

			seg = segs[i]
			// assert seg.dirs < 0

			i++
		} else {
			if from >= nLen {
				from-- // To check the last dir separator
			}
			return matchAnyDirs(seg.pattern, name[from:])
		}
	}

	return
}

func searchMatched(segs []patternSeg, i, segsLen int, name string, from, nLen int) (to int, matched bool, err error) {
	seg := segs[i]
	// assert seg.dirs > 0

	origin := from

	var mark int

	if to, mark = scanDirs(name, from, nLen, seg.dirs); mark < 0 {
		return
	}

	pattern := seg.pattern

	for {
		if mark == 0 && isSegLastAndAny(segs, i+1, segsLen) {
			pattern = pattern[:len(pattern)-1]
		}
		matched, err = filepath.Match(pattern, name[from:to])

		if matched {
			return
		} else if err != nil {
			break
		}

		// move next
		from, _ = scanDirs(name, from, nLen, 1)

		if to, mark = scanDirs(name, to, nLen, 1); mark < 0 {
			break
		}
	}

	to = origin
	return
}

func scanDirs(name string, from, len, dirs int) (int, int) {
	i := from
	for ; i < len; i++ {
		switch name[i] {
		case '\\':
			if runtime.GOOS != "windows" {
				continue
			}
			fallthrough
		case '/':
			if i > from {
				dirs--
				if dirs <= 0 {
					return i + 1, 1
				}
			}
		}
	}

	if i > from {
		dirs--
		if dirs <= 0 {
			return i, 0
		}
	}

	return from, -1
}

func isSegLastAndAny(segs []patternSeg, i, segsLen int) bool {
	if i == segsLen-1 {
		// assert segs[i].dirs < 0

		pattern := segs[i].pattern
		return pattern[len(pattern)-1] == '*'
	}

	return false
}
