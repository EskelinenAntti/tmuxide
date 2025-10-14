package project

import (
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v6"
)

func Root(inputPath string) (string, error) {
	absolutePath, err := filepath.Abs(inputPath)
	if err != nil {
		return "", err
	}

	if repository, err := getGitRoot(absolutePath); err == nil {
		return repository, nil
	}

	fileInfo, err := os.Stat(inputPath)
	if err != nil {
		return "", err
	}

	if !fileInfo.IsDir() {
		return filepath.Dir(absolutePath), nil
	}

	return absolutePath, nil
}

func getGitRoot(path string) (string, error) {
	repo, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{
		DetectDotGit: true, // detect parent .git directories
	})
	if err != nil {
		return "", err
	}
	wt, err := repo.Worktree()
	if err != nil {
		return "", err
	}
	return wt.Filesystem.Root(), nil
}
