package project

import (
	"os"
	"path/filepath"
)

type RepositoryResolver interface {
	Root(string) (string, error)
}

func Root(target string, resolver RepositoryResolver) (string, error) {
	absolutePath, err := filepath.Abs(target)
	if err != nil {
		return "", err
	}

	if repository, err := resolver.Root(absolutePath); err == nil {
		return repository, nil
	}

	fileInfo, err := os.Stat(target)
	if err != nil {
		return "", err
	}

	if !fileInfo.IsDir() {
		return filepath.Dir(absolutePath), nil
	}

	return absolutePath, nil
}
