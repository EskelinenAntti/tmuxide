package path

import "os/exec"

type Path interface {
	Contains(path string) bool
}

type ShellPath struct{}

func (ShellPath) Contains(path string) bool {
	_, err := exec.LookPath(path)
	return err != nil
}
