package main

import (
	"path/filepath"
	"strings"
)

// check if the path is valid, such as /../a.txt
func check(filePath string) bool {
	if !strings.HasPrefix(filePath, "/") {
		return false
	}
	if filePath == "/" {
		return true
	}
	path1, err := filepath.Abs("." + filePath)
	path2, err := filepath.Abs(".")
	path3 := filepath.Clean(filePath)
	if err != nil {
		return false
	}
	return path1 == path2+path3
}
