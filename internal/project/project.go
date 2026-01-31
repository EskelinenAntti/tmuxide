package project

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/eskelinenantti/tmuxide/internal/shell/tmux"
)

var ErrInvalidPath = errors.New("invalid path")
var ErrNotADirectory = errors.New("not a directory")

type Project struct {
	Name       string
	WorkingDir string
}

type Input struct {
	EditorPath string
}

type Git interface {
	RevParse(cwd string) (string, error)
}

func ForPath(path string, git Git, tmux tmux.Tmux) (Project, error) {
	if tmux.HasSession(path, "") {
		return Project{
			Name: path,
		}, nil
	}
	workingDir, err := repository(path, git)
	if err != nil {
		if workingDir, err = dir(path); err != nil {
			return Project{}, err
		}
	}

	absolutePath, err := filepath.Abs(workingDir)
	if err != nil {
		return Project{}, err
	}

	name := Name(absolutePath)
	return Project{
		Name:       name,
		WorkingDir: workingDir,
	}, nil
}

func ForDir(directory string, tmux tmux.Tmux) (Project, error) {
	if tmux.HasSession(directory, "") {
		return Project{
			Name: directory,
		}, nil
	}

	fileInfo, err := os.Stat(directory)
	if err != nil {
		return Project{}, ErrInvalidPath
	}

	if !fileInfo.IsDir() {
		return Project{}, ErrNotADirectory
	}

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
