package shell

import "os/exec"

type Path struct{}

func (Path) Contains(path string) bool {
	_, err := exec.LookPath(path)
	return err == nil
}
