package project

import (
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

func Name(path string) string {
	basename := filepath.Base(path)
	sessionPrefix := strings.ReplaceAll(basename, ".", "-")
	return strings.Join([]string{sessionPrefix, hash(path)}, "-")
}

func hash(path string) string {
	hash := sha1.New()
	hash.Write([]byte(path))
	hashByteSlice := hash.Sum(nil)
	return fmt.Sprintf("%x", hashByteSlice)[:4]
}
