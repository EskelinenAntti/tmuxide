package project

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var ErrInvalidPath = errors.New("invalid path")

type Project struct {
	Name       string
	WorkingDir string
}

type Git interface {
	RevParse(cwd string) (string, error)
}

func ForFile(file string, git Git) (Project, error) {
	workingDir, err := dir(file)
	if err != nil {
		return Project{}, err
	}
	if repository, err := git.RevParse(workingDir); err == nil {
		workingDir = repository
	}

	return ForDir(workingDir)
}

func ForDir(directory string) (Project, error) {
	absoluteDir, err := filepath.Abs(directory)
	if err != nil {
		return Project{}, err
	}

	return Project{
		Name:       Name(absoluteDir),
		WorkingDir: absoluteDir,
	}, nil
}

func Name(path string) string {
	basename := filepath.Base(path)
	sessionPrefix := strings.ReplaceAll(basename, ".", "_")
	return strings.Join([]string{sessionPrefix, hash(path)}, "-")
}

func dir(target string) (string, error) {
	fileInfo, err := os.Stat(target)
	if err != nil {
		return "", ErrInvalidPath
	}

	if !fileInfo.IsDir() {
		return filepath.Dir(target), nil
	}

	return target, nil
}

func hash(path string) string {
	hash := sha1.New()
	hash.Write([]byte(path))
	hashByteSlice := hash.Sum(nil)
	return fmt.Sprintf("%x", hashByteSlice)[:4]
}
