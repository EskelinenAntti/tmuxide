package project

import (
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Project struct {
	Name       string
	WorkingDir string
}

type Input struct {
	WorkingDir string
	Command    []string
	EditorPath string
}

type Git interface {
	RevParse(cwd string) (string, error)
}

func New(input Input, git Git) (Project, error) {
	if input.WorkingDir != "" {
		return Project{
			Name:       Name(input.WorkingDir),
			WorkingDir: input.WorkingDir,
		}, nil
	}

	target := input.EditorPath
	if input.EditorPath == "" {
		var err error
		if target, err = os.Getwd(); err != nil {
			return Project{}, err
		}
	}

	workingDir, err := repository(target, git)
	if err != nil {
		if workingDir, err = dir(target); err != nil {
			return Project{}, err
		}
	}

	name := Name(workingDir)
	return Project{
		Name:       name,
		WorkingDir: workingDir,
	}, nil
}

func Name(path string) string {
	basename := filepath.Base(path)
	sessionPrefix := strings.ReplaceAll(basename, ".", "_")
	return strings.Join([]string{sessionPrefix, hash(path)}, "-")
}

func dir(target string) (string, error) {
	absolutePath, err := filepath.Abs(target)
	if err != nil {
		return "", err
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

func repository(target string, git Git) (string, error) {
	fileInfo, err := os.Stat(target)
	if err != nil {
		return "", err
	}

	var cwd string
	if fileInfo.IsDir() {
		cwd = target
	} else {
		cwd = filepath.Dir(target)
	}

	return git.RevParse(cwd)
}

func hash(path string) string {
	hash := sha1.New()
	hash.Write([]byte(path))
	hashByteSlice := hash.Sum(nil)
	return fmt.Sprintf("%x", hashByteSlice)[:4]
}
