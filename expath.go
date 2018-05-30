// Package expath extends the standard library path/filepath's Match and Glob routines
// for supporting any-directories' pattern (** matches zero or more directories in a path).

package expath

import (
	"os"
)

// Match reports whether name matches the shell file name pattern.
// The function extends the standard library path/filepath's Match function.
// Its pattern syntax is compatible with the standard Match's pattern syntax,
// and supporting the new features:
//
//	term:
//		'**'         matches zero or more directories in a path
//
// As the standard Match function, the only possible returned error is filepath.ErrBadPattern, when pattern
// is malformed.
//
func Match(pattern, name string) (matched bool, err error) {
	segs, err := scanSegments(pattern)
	if err != nil {
		return false, err
	}

	switch len(segs) {
	case 1:
		return matchASeg(segs[0], name)
	case 0:
		return matchASeg(patternSeg{"", 0}, name)
	default:
		return matchSegs(segs, name)
	}
}

// Glob returns the names of all files matching pattern or nil
// if there is no matching file. The syntax of patterns is the same
// as in Match.
//
// If the pattern ends with Separator ('/'), this is the same as '/**'.
//
// Unlike the standard library path/filepath's Glob function, this Glob function has an extra root argument.
// The root argument indicates that the pattern path based on the root (empty root means the current direction).
//
func Glob(pattern, root string) (matches []string, atRoot string, err error) {
	var helper filePathHelper
	var mh matchedSet

	err = doGlob(pattern, root, &helper, &mh)

	atRoot = mh.root
	matches = mh.matches
	return
}

// GlobInfo used for the GlobFunc callback function to supply the glob information.
//
type GlobInfo interface {
	AtRoot() string
	Path() string
	FullName() string
	FileInfo() (os.FileInfo, error)
}

// GlobFunc used by GlobFn, called for each matched file name or encountered file error.
//
// err != nil means encountering a file error, returning nil could ignore the error.
//
type GlobFunc func(info GlobInfo, err error) error

// GlobFn uses the GlobFunc callback function to handle each matched file name or encountered file error.
//
func GlobFn(pattern, root string, globFn GlobFunc) error {
	var helper filePathHelper
	var mf matchesFunc

	return doGlob(pattern, root, &helper, &mf)
}
