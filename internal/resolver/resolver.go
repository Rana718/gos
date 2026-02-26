package resolver

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ResolvePath(nameOrPath string, loadedPaths map[string]string) (string, error) {
	if nameOrPath == "." {
		return os.Getwd()
	}

	// Check stored name with backslash sub-path (e.g. "backend\src")
	parts := strings.Split(nameOrPath, "\\")
	if basePath, exists := loadedPaths[parts[0]]; exists {
		if len(parts) > 1 {
			return filepath.Join(basePath, filepath.Join(parts[1:]...)), nil
		}
		return basePath, nil
	}

	// Check stored name with forward-slash sub-path (e.g. "backend/src")
	parts = strings.Split(nameOrPath, "/")
	if basePath, exists := loadedPaths[parts[0]]; exists {
		if len(parts) > 1 {
			return filepath.Join(basePath, filepath.Join(parts[1:]...)), nil
		}
		return basePath, nil
	}

	// Treat as literal filesystem path
	if _, err := os.Stat(nameOrPath); err == nil {
		return nameOrPath, nil
	}

	return "", fmt.Errorf("'%s' is not a saved name or valid path", nameOrPath)
}
