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
	workingDir := path
	file := ""
	if !info.IsDir() {
		file = path
		workingDir = filepath.Dir(path)
		if repository, err := git.RevParse(workingDir); err == nil {
			workingDir = repository
		}
	}
	proj, err := ForDir(workingDir)
	if err != nil {
		return Target{}, err
	}
	return Target{Project: proj, File: file}, nil
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
