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

// Target describes the project to open and the file to edit, if any.
type Target struct {
	Project Project
	File    string
}

// ForPath keeps directory targets at the selected directory and resolves file
// targets to their repository root when available.
func ForPath(path string, git Git) (Target, error) {
	info, err := os.Stat(path)
	if err != nil {
		return Target{}, ErrInvalidPath
	}
	if info.IsDir() {
		proj, err := ForDir(path)
		return Target{Project: proj}, err
	}
	proj, err := forWorkingDir(filepath.Dir(path), git)
	if err != nil {
		return Target{}, err
	}
	return Target{Project: proj, File: path}, nil
}

func ForFile(file string, git Git) (Project, error) {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return Project{}, ErrInvalidPath
	}
	workingDir := file
	if !fileInfo.IsDir() {
		workingDir = filepath.Dir(file)
	}
	return forWorkingDir(workingDir, git)
}

func forWorkingDir(workingDir string, git Git) (Project, error) {
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
	digest := sha1.Sum([]byte(path))
	return fmt.Sprintf("%s-%x", sessionPrefix, digest[:2])
}
