package shell

import (
	"bytes"
	"os"
	"os/exec"
)

type FdCmd struct {
	Runner
}

func (f FdCmd) Fd() ([]byte, error) {
	args := []string{
		"--follow",
		"--hidden",
		"--exclude", "{.git,node_modules,target,build,Library}",
		".",
		os.Getenv("HOME"),
	}

	fdCmd := exec.Command("fd", args...)
	var out bytes.Buffer
	fdCmd.Stdout = &out

	err := f.Run(fdCmd)
	return out.Bytes(), err
}
