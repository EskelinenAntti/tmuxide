package ide

import (
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/eskelinenantti/tmuxide/internal/git"
)

type Project struct {
	Name string
	Root string
}

func Repository(target string, git git.Command) (string, error) {
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

func ProjectFor(target string, repository string) (Project, error) {
	root, err := root(target, repository)
	if err != nil {
		return Project{}, err
	}

	name := name(target)

	return Project{
		Name: name,
		Root: root,
	}, nil
}

func name(path string) string {
	basename := filepath.Base(path)
	sessionPrefix := strings.ReplaceAll(basename, ".", "_")
	return strings.Join([]string{sessionPrefix, hash(path)}, "-")
}

func root(target string, repository string) (string, error) {
	absolutePath, err := filepath.Abs(target)
	if err != nil {
		return "", err
	}

	if repository != "" {
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

func hash(path string) string {
	hash := sha1.New()
	hash.Write([]byte(path))
	hashByteSlice := hash.Sum(nil)
	return fmt.Sprintf("%x", hashByteSlice)[:4]
}
