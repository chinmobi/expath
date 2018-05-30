package expath

import (
	"os"
)

// matchesHandler used by the glob routine to handle the matched result.
// It also acts as a role to decouple handling matched result from the glob algorithm.
//
type matchesHandler interface {
	onMatched(matched string) error
	onError(path string, err error) error
	setRoot(root string) error
}

// matchedSet implements the matchesHandler interface to collect the matched results.
//
type matchedSet struct {
	root    string
	matches []string
}

func (m *matchedSet) onMatched(matched string) error {
	m.matches = append(m.matches, matched)
	return nil
}

func (m *matchedSet) onError(path string, err error) error {
	return err
}

func (m *matchedSet) setRoot(root string) error {
	m.root = root
	return nil
}


// matchedSet implements the matchesHandler interface to use the GlobFunc to handle the matched results.
//
type matchesFunc struct {
	root   string
	globFn GlobFunc
}

func (m *matchesFunc) onMatched(matched string) error {
	var info matchesInfo
	info.root = m.root
	info.path = matched

	return m.globFn(&info, nil)
}

func (m *matchesFunc) onError(path string, err error) error {
	var info matchesInfo
	info.root = m.root
	info.path = path

	return m.globFn(&info, err)
}

func (m *matchesFunc) setRoot(root string) error {
	m.root = root
	return nil
}

// matchesInfo implements the GlobInfo interface.
//
type matchesInfo struct {
	root string
	path string
}

func (m *matchesInfo) AtRoot() string {
	return m.root
}

func (m *matchesInfo) Path() string {
	return m.path
}

func (m *matchesInfo) FullName() string {
	return appendDirPath(m.root, "", m.path)
}

func (m *matchesInfo) FileInfo() (os.FileInfo, error) {
	fullName := m.FullName()
	return os.Lstat(fullName)
}
