package path

import "os/exec"

type PathLooker interface {
	LookPath(path string) (string, error)
}

type ExecPathLooker struct{}

func (ExecPathLooker) LookPath(path string) (string, error) {
	return exec.LookPath(path)
}
