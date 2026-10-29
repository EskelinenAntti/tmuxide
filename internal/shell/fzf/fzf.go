package fzf

import (
	"errors"
	"io"
	"os"
	"os/exec"

	"github.com/eskelinenantti/tmuxide/internal/shell/runner"
)

var ErrCancelled = errors.New("fzf cancelled")

type Cmd struct {
	runner.Runner
}

func (f Cmd) Fzf(input io.Reader, output io.Writer) error {
	fzfCmd := exec.Command("fzf", "--reverse", "--height", "70%", "--tmux", "70%", "--no-sort")
	fzfCmd.Stderr = os.Stderr
	err := f.Pipe(input, fzfCmd, output)
	var inputErr *runner.InputError
	if errors.As(err, &inputErr) {
		return err
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == 130 {
		return ErrCancelled
	}
	return err
}
