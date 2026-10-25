package fzf

import (
	"errors"
	"io"
	"os"
	"os/exec"

	"github.com/eskelinenantti/tmuxide/internal/shell/runner"
)

type Cmd struct {
	runner.Starter
}

// IsCancelled reports whether fzf exited because the user cancelled selection.
func IsCancelled(err error) bool {
	var exitErr *exec.ExitError
	return errors.As(err, &exitErr) && exitErr.ExitCode() == 130
}

func (f Cmd) Fzf(output io.Writer) (runner.WriteCloser, error) {
	fzfCmd := exec.Command("fzf", "--reverse", "--height", "70%", "--tmux", "70%")
	fzfCmd.Stdout = output
	fzfCmd.Stderr = os.Stderr
	waiter, err := f.Start(fzfCmd)
	return waiter, err

}
