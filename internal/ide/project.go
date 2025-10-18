package ide

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Window struct {
	Cmd  string
	Args []string
}

type Project struct {
	Name    string
	Root    string
	Windows []Window
}

type Repository interface {
	Root(path string) (string, error)
}

func ProjectFor(target string, repository Repository) (Project, error) {
	root, err := root(target, repository)
	if err != nil {
		return Project{}, err
	}

	name := name(target)

	editor, err := editor(target)
	if err != nil {
		return Project{}, err
	}

	windows := []Window{editor}

	if lazygit, err := lazygit(target, repository); err != nil {
		windows = append(windows, lazygit)
	}

	return Project{
		Name:    name,
		Root:    root,
		Windows: windows,
	}, nil
}

func lazygit(target string, repository Repository) (Window, error) {
	if _, err := exec.LookPath("lazygit"); err != nil {
		return Window{}, errors.New("Lazygit is not isntalled")
	}

	if _, err := repository.Root(target); err != nil {
		return Window{}, errors.New("Lazygit is not isntalled")
	}

	return Window{Cmd: "lazygit"}, nil
}

func name(path string) string {
	basename := filepath.Base(path)
	sessionPrefix := strings.ReplaceAll(basename, ".", "-")
	return strings.Join([]string{sessionPrefix, hash(path)}, "-")
}

func root(target string, repository Repository) (string, error) {
	absolutePath, err := filepath.Abs(target)
	if err != nil {
		return "", err
	}

	if repository, err := repository.Root(absolutePath); err == nil {
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

func editor(target string) (Window, error) {
	editorCmd, hasEditor := os.LookupEnv("EDITOR")
	if !hasEditor {
		return Window{}, errors.New(
			"No editor was configured. Specify the editor you would like to use by setting the $EDITOR variable.\n\n" +
				"For example, to use Vim as your editor, add the following line to your ~/.zshrc:\n" +
				"export EDITOR=vim\n",
		)
	}

	return Window{Cmd: editorCmd, Args: []string{target}}, nil
}
