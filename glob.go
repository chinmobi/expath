package expath

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
)

// doGlob is the main entrance of the glob routine.
//
func doGlob(pattern, root string, helper pathHelper, matches matchesHandler) (err error) {

	pattern, root = normalizePath(pattern, root)

	err = matches.setRoot(root)
	if err != nil {
		return
	}

	var segs []patternSeg

	segs, err = scanSegments(pattern)
	if err != nil || len(segs) == 0 {
		return
	}

	var matchedPath string

	if segs[0].pattern[0] == '/' {
		matchedPath = "/"
	}

	return segsGlob(segs, 0, root, matchedPath, helper, matches)
}

// segsGlob does the glob match of current pattern segment.
//
func segsGlob(segs []patternSeg, curr int,
	dir, matchedPath string,
	helper pathHelper, matches matchesHandler) (err error) {

	segsLen := len(segs)

	if (segsLen - curr) > 1 {
		if segs[curr].dirs >= 0 {

			var mh matchedSet
			err = normalGlob(dir, matchedPath, trimPath(segs[curr].pattern), helper, &mh)
			if err != nil {
				return
			}

			if (segsLen - curr) > 2 {
				for _, mp := range mh.matches {
					err = exglob(segs, curr+2,
						appendDirPath(dir, matchedPath, mp), mp,
						len(mp), 0,
						helper, matches)
					if err != nil {
						return
					}
				}
			} else {
				for _, mp := range mh.matches {
					err = anyDirsGlob(appendDirPath(dir, matchedPath, mp), mp, helper, matches)
					if err != nil {
						return
					}
				}
			}

		} else {
			err = exglob(segs, curr+1,
				dir, matchedPath,
				len(matchedPath), 0,
				helper, matches)
		}

	} else {
		if segs[curr].dirs >= 0 {
			err = normalGlob(dir, matchedPath, segs[curr].pattern, helper, matches)
		} else {
			err = anyDirsGlob(dir, matchedPath, helper, matches)
		}
	}

	return
}

func appendDirPath(dir, appendedPrefix, path string) string {
	nLen := len(appendedPrefix)
	// assert nLen < len(path)

	path = path[nLen:]

	if path[0] != '/' {
		nLen = len(dir)
		if nLen > 0 && !isDirSeparator(dir, nLen-1) {
			dir += string(os.PathSeparator)
		}
	}

	dir += filepath.FromSlash(path)

	return dir
}

// exglob does the glob match of the normal pattern segment that following the any-dirs' pattern segment.
//
func exglob(segs []patternSeg, curr int,
	dir, matchedPath string,
	mark, pendingDirs int,
	helper pathHelper, matches matchesHandler) error {

	names, err := helper.getNames(dir)
	if err != nil {
		return matches.onError(matchedPath, err)
	}
	if len(names) == 0 {
		return nil
	}

	pendingDirs++

	if pendingDirs >= segs[curr].dirs {

		segsLen := len(segs)

		for _, name := range names {
			d, p := appendDir(dir, name, false), appendPath(matchedPath, name)

			matched, err := path.Match(trimPath(segs[curr].pattern), trimPath(p[mark:]))
			if err != nil {
				return err
			}

			if matched {
				if curr == segsLen-1 {
					err = matches.onMatched(p)
				} else {
					err = segsGlob(segs, curr+1, d, p, helper, matches)
				}
			} else {
				m, _ := scanDirs(p, mark, len(p), 1)
				err = exglob(segs, curr, d, p, m, pendingDirs-1, helper, matches)
			}

			if err != nil {
				return err
			}
		}

	} else {
		for _, name := range names {
			d, p := appendDir(dir, name, false), appendPath(matchedPath, name)
			err = exglob(segs, curr, d, p, mark, pendingDirs, helper, matches)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func trimPath(path string) string {
	nLen := len(path)
	if nLen > 0 {
		i, j := 0, nLen-1
		if path[0] == '/' {
			i++
		}
		if j > i && path[j] == '/' {
			j--
		}
		return path[i : j+1]
	}
	return path
}

// normalGlob globs the normal pattern that not following the any-dirs' pattern.
//
func normalGlob(dir, matchedPath, pattern string,
	helper pathHelper, matches matchesHandler) error {

	head, tail, hasMeta, slashes := headOfPath(pattern)
	if len(head) == 0 {
		return nil
	}

	morePattern := (len(tail) > 0)

	if !hasMeta {
		dir = appendDir(dir, head, slashes > 0)
		matchedPath = appendPath(matchedPath, head)

		exists, err := helper.isExist(dir)
		if err != nil {
			return matches.onError(matchedPath, err)
		}
		if !exists {
			return err
		}

		if morePattern {
			return normalGlob(dir, matchedPath, tail, helper, matches)
		} else {
			err = matches.onMatched(matchedPath)
			if err != nil {
				return err
			}
		}

	} else {
		names, err := helper.getNames(dir)
		if err != nil {
			return matches.onError(matchedPath, err)
		}

		for _, name := range names {
			matched, err := path.Match(head, name)
			if err != nil {
				return err
			}

			if !matched {
				continue
			}

			if morePattern {
				d, p := appendDir(dir, name, false), appendPath(matchedPath, name)
				err = normalGlob(d, p, tail, helper, matches)
				if err != nil {
					return err
				}
			} else {
				mp := appendPath(matchedPath, name)
				err = matches.onMatched(mp)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func headOfPath(path string) (head, tail string, hasMeta bool, slashes int) {
	nLen := len(path)

	var i int
	if nLen > 0 && path[0] == '/' {
		i++
	}
	from := i

	for mark := -1; i < nLen; i++ {
		switch path[i] {
		case '*', '?', '[':
			hasMeta = true

			if mark >= 0 {
				head, tail = path[from:mark], path[mark:]
				return
			}
		case '/':
			mark = i

			if hasMeta {
				head, tail = path[from:mark], path[mark:]
				return
			}

			slashes++
		}
	}

	head, tail = path[from:i], path[i:]
	return
}

// anyDirsGlob globs the lastest any-dirs' pattern.
//
func anyDirsGlob(dir, matchedPath string,
	helper pathHelper, matches matchesHandler) error {

	names, err := helper.getNames(dir)
	if err != nil {
		return matches.onError(matchedPath, err)
	}

	if len(names) == 0 && isValidMatched(matchedPath) {
		return matches.onMatched(matchedPath)
	}

	for _, name := range names {
		d, p := appendDir(dir, name, false), appendPath(matchedPath, name)
		err = anyDirsGlob(d, p, helper, matches)
		if err != nil {
			break
		}
	}

	return err
}

func isValidMatched(path string) bool {
	switch len(path) {
	case 0:
		return false
	case 1:
		return path[0] != '/'
	default:
		return true
	}
}

func appendDir(dir, name string, hasSlash bool) string {
	nLen := len(dir)
	if nLen > 0 && !isDirSeparator(dir, nLen-1) {
		dir += string(os.PathSeparator)
	}

	if hasSlash {
		dir += filepath.FromSlash(name)
	} else {
		dir += name
	}

	return dir
}

func appendPath(path, name string) string {
	nLen := len(path)
	if nLen > 0 && path[nLen-1] != '/' {
		path += "/"
	}
	path += name

	return path
}

// normalizePath normalizes the pattern and root.
//
func normalizePath(pattern, root string) (string, string) {

	nLen := len(filepath.VolumeName(pattern))
	if nLen > 0 { // Use the pattern's root to replace the given root
		if pattern[nLen-1] == ':' &&
			nLen < len(pattern) && isDirSeparator(pattern, nLen) {
			root = pattern[:nLen+1]
		} else {
			root = pattern[:nLen]
		}

		pattern = pattern[nLen:]
	}

	nLen = len(pattern)
	if nLen > 0 && isDirSeparator(pattern, nLen-1) {
		pattern += "**"
		nLen += 2
	}

	mark := skipDotsDir(pattern, nLen)
	if mark > 0 {
		pre := pattern[:mark]
		pattern = pattern[mark:]

		nLen = len(root)
		if isDirSeparator(pre, 0) {
			if nLen > 0 && isDirSeparator(root, nLen-1) {
				pre = pre[1:]
			}
		} else {
			if nLen > 0 &&
				root[nLen-1] != ':' && !isDirSeparator(root, nLen-1) {
				root += string(os.PathSeparator)
			}
		}

		root += pre
	}

	root = filepath.Clean(filepath.FromSlash(root))
	pattern = filepath.ToSlash(pattern)

	nLen = len(pattern)
	if nLen > 0 && !isDirSeparator(pattern, 0) {
		nLen = len(root)
		if nLen > 0 && !isDirSeparator(root, nLen-1) {
			root += string(os.PathSeparator)
		}
	}

	return pattern, root
}

func skipDotsDir(pattern string, nLen int) (mark int) {
	mark = -1
	i := 0
	for ; i < nLen; i++ {
		switch pattern[i] {
		case '\\':
			if runtime.GOOS != "windows" {
				return
			}
			fallthrough
		case '/':
			mark = i
		case '.': // Nothing to do, just keeping loop
		default:
			return
		}
	}

	if mark != i-1 {
		mark = i
	}
	return
}
