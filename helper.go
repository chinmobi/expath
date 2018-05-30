package expath

import (
	"os"
)

// pathHelper helps the glob routine to acquire the path information.
// It also acts as a role to decouple retrieving path information from the glob algorithm.
//
type pathHelper interface {
	getNames(dir string) (names []string, err error)
	isExist(dir string) (bool, error)
}

// filePathHelper implements the pathHelper interface by retrieving os file information.
//
type filePathHelper struct{}

func (filePathHelper) getNames(dir string) ([]string, error) {
	fi, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	if !fi.IsDir() {
		return nil, nil
	}

	d, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if names == nil {
		names = []string{}
	}
	return names, err
}

func (filePathHelper) isExist(dir string) (bool, error) {
	if _, err := os.Lstat(dir); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}
