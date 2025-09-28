package util_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rogeriofbrito/go-insta-scraper-v2/util"
)

// Table-driven tests for CreateWorkingDir covering main scenarios.
func TestCreateWorkingDir_Scenarios(t *testing.T) {
	baseTmp := filepath.Join(os.TempDir(), "go-insta-scraper-v2-test")
	// Clean up before and after
	_ = os.RemoveAll(baseTmp)
	defer os.RemoveAll(baseTmp)

	tests := []struct {
		name         string
		setup        func(dir string) error
		expectExists bool
		expectEmpty  bool
	}{
		{
			name: "path_does_not_exist_creates_dir",
			setup: func(dir string) error {
				// Ensure directory does not exist
				return os.RemoveAll(dir)
			},
			expectExists: true,
			expectEmpty:  true,
		},
		{
			name: "path_exists_empty_dir",
			setup: func(dir string) error {
				_ = os.RemoveAll(dir)
				return os.MkdirAll(dir, 0777)
			},
			expectExists: true,
			expectEmpty:  true,
		},
		{
			name: "path_exists_with_items_removes_them",
			setup: func(dir string) error {
				_ = os.RemoveAll(dir)
				if err := os.MkdirAll(dir, 0777); err != nil {
					return err
				}
				// Create a file and a subdir
				if err := os.WriteFile(filepath.Join(dir, "file.txt"), []byte("data"), 0666); err != nil {
					return err
				}
				subdir := filepath.Join(dir, "subdir")
				if err := os.Mkdir(subdir, 0777); err != nil {
					return err
				}
				return os.WriteFile(filepath.Join(subdir, "nested.txt"), []byte("nested"), 0666)
			},
			expectExists: true,
			expectEmpty:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testDir := filepath.Join(baseTmp, tc.name)
			if err := tc.setup(testDir); err != nil {
				t.Fatalf("setup failed: %v", err)
			}
			if err := util.CreateWorkingDir(testDir); err != nil {
				t.Fatalf("CreateWorkingDir failed: %v", err)
			}
			// Check directory exists
			info, err := os.Stat(testDir)
			if tc.expectExists && (err != nil || !info.IsDir()) {
				t.Fatalf("expected dir to exist: %v", err)
			}
			// Check directory is empty
			entries, err := os.ReadDir(testDir)
			if err != nil {
				t.Fatalf("failed to read dir: %v", err)
			}
			if tc.expectEmpty && len(entries) != 0 {
				t.Fatalf("expected dir to be empty, found %d items", len(entries))
			}
		})
	}
}
