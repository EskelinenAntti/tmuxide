package shell

import "os/exec"

type Path struct{}

type ShellPath interface {
	Contains(path string) bool
}

func (Path) Contains(path string) bool {
	_, err := exec.LookPath(path)
	return err == nil
}
