package shell

import (
	"io"
	"os"
	"os/exec"
)

type FdCmd struct {
	Runner
}

func (f FdCmd) Fd(output io.Writer) error {
	args := []string{
		"--follow",
		"--hidden",
		"--exclude", "{.git,node_modules,target,build,Library}",
		".",
		os.Getenv("HOME"),
	}

	fdCmd := exec.Command("fd", args...)
	fdCmd.Stdout = output

	return f.Run(fdCmd)
}
