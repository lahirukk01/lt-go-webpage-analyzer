package utils

import (
	"path/filepath"
	"runtime"
)

func GetProjectRoot() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}
	currentDir := filepath.Dir(filename)
	projectRoot := filepath.Join(currentDir, "..", "..")

	return projectRoot
}
