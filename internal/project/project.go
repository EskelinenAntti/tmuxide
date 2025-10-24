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
	TargetPath string
	WorkingDir string
	IsGitRepo  bool
}

type Git interface {
	RevParse(cwd string) (string, error)
}

func New(target string, git Git) (Project, error) {
	name := Name(target)

	var workingDir string

	workingDir, err := repository(target, git)
	isGitRepo := err == nil

	if !isGitRepo {
		if workingDir, err = dir(target); err != nil {
			return Project{}, err
		}
	}

	return Project{
		Name:       name,
		TargetPath: target,
		WorkingDir: workingDir,
		IsGitRepo:  isGitRepo,
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
