package project

import (
	"os"
	"path/filepath"
)

func Root(inputPath string) (string, error) {
	fileInfo, err := os.Stat(inputPath)
	if err != nil {
		return "", err
	}
	absolutePath, err := filepath.Abs(inputPath)
	if err != nil {
		return "", err
	}
	if !fileInfo.IsDir() {
		return filepath.Dir(absolutePath), nil
	}
	return absolutePath, nil
}
