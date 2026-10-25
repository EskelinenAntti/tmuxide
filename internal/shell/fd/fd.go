package fd

import (
	"errors"
	"io"
	"os/exec"

	"github.com/eskelinenantti/tmuxide/internal/shell/runner"
)

type Cmd struct {
	runner.Runner
}

// StreamCandidates writes file and directory paths relative to baseDir to output.
func (f Cmd) StreamCandidates(baseDir string, output io.Writer) error {
	args := []string{
		"--follow",
		"--hidden",
		"--exclude", "{.git,node_modules,Library}",
		".",
		"--base-directory",
		baseDir,
	}

	fdCmd := exec.Command("fd", args...)
	fdCmd.Stdout = output

	err := f.Run(fdCmd)
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		return nil
	}
	// This error occurs if fzf closes the pipe before the command is completed.
	if exitErr.ExitCode() == 1 {
		return nil
	}
	return err
}
