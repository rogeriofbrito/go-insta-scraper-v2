package util

import (
	"os"
	"path/filepath"

	"github.com/palantir/stacktrace"
)

func CreateWorkingDir(path string) error {
	// Check if directory exists
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// Create if not exists
		err = os.MkdirAll(path, 0777)
		if err != nil {
			return stacktrace.Propagate(err, "failed to create working dir at path %s", path)
		}
	} else if err == nil {
		// Remove all contents but keep the directory itself
		entries, err := os.ReadDir(path)
		if err != nil {
			return stacktrace.Propagate(err, "failed to read working dir at path %s", path)
		}

		for _, entry := range entries {
			entryPath := filepath.Join(path, entry.Name())
			err = os.RemoveAll(entryPath)
			if err != nil {
				return stacktrace.Propagate(err, "failed to remove working dir item at path %s", entryPath)
			}
		}
	} else {
		// Unexpected error
		return stacktrace.Propagate(err, "failed to check working dir at path %s", path)
	}

	return nil
}
