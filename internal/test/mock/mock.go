package mock

import (
	"errors"
	"os/exec"
	"slices"
)

func SimulateError(cmd *exec.Cmd) error { return errors.New("mock error") }
func WriteToStdout(value string) func(*exec.Cmd) error {
	return func(cmd *exec.Cmd) error {
		cmd.Stdout.Write([]byte(value))
		return nil
	}
}

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
