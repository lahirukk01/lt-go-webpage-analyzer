package utils

import "testing"

func TestGetProjectRoot(t *testing.T) {
	root := GetProjectRoot()
	if root == "" {
		t.Error("GetProjectRoot() returned an empty string, which indicates failure")
	}
}
