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
	fzfCmd.Stdin = input
	fzfCmd.Stdout = output
	fzfCmd.Stderr = os.Stderr
	err := f.Run(fzfCmd)
	if fzfCmd.ProcessState != nil && fzfCmd.ProcessState.ExitCode() == 130 {
		return ErrCancelled
	}
	return err
}
