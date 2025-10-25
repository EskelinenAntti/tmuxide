package test

import (
	"errors"
	"slices"
)

type Path struct {
	Missing []string
}

type Git struct {
	Repository string
}

func (pathMock Path) Contains(program string) bool {
	return !slices.Contains(pathMock.Missing, program)
}

func (git Git) RevParse(cwd string) (string, error) {

	if git.Repository == "" {
		return "", errors.New("not inside git repo")
	}

	return git.Repository, nil
}
