package codegen

import (
	"fmt"
	"os"
	"path/filepath"
)

// Root finds the repository root by walking up from the working directory
// looking for go.mod, so the tool works from any subdirectory.
func Root() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("no go.mod found above %s", dir)
		}
		dir = parent
	}
}
