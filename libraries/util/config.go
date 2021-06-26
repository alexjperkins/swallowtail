package util

import (
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/monzo/terrors"
)

// RootDir returns the root directory of the project.
// This is the grandparent directory of this package.
func RootDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", terrors.Augment(err, "Failed to get cwd", nil)
	}

	// A dirty hack; this is the case for when we are in a Docker container.
	if cwd == "/" {
		return "", nil
	}

	// This allows us to pull the correct config path if we are running tests
	// from any child directory from the root project.
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	dd := path.Dir(d)
	return filepath.Dir(dd), nil
}
