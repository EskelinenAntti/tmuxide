package mock

import (
	"errors"
	"slices"
)

type Path struct {
	Missing []string
}

func (path Path) Contains(program string) bool {
	return !slices.Contains(path.Missing, program)
}

type Git struct {
	Repository string
}

func (git Git) RevParse(cwd string) (string, error) {

	if git.Repository == "" {
		return "", errors.New("not inside git repo")
	}

	return git.Repository, nil
}
