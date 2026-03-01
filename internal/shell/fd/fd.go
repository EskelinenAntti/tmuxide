package fd

import (
	"errors"
	"io"
	"os"
	"os/exec"

	"github.com/eskelinenantti/tmuxide/internal/shell/runner"
)

type Cmd struct {
	runner.Runner
}

func (f Cmd) Fd(filterDir bool, output io.Writer) error {
	args := []string{}
	if filterDir {
		args = append(args, "--type", "dir")
	}

	args = append(args,
		"--follow",
		"--hidden",
		"--exclude", "{.git,node_modules,target,build,Library}",
		".",
		"--base-directory",
		os.Getenv("HOME"),
	)
	fdCmd := exec.Command("fd", args...)
	fdCmd.Stdout = output

	err := f.Run(fdCmd)

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		if exitErr.ExitCode() == 1 {
			// This error occurs if fzf closes the pipe before the command is completed
			return nil
		}
		return err
	}
	return nil
}
