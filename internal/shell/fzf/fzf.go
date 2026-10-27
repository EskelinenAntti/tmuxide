package fzf

import (
	"io"
	"os"
	"os/exec"

	"github.com/eskelinenantti/tmuxide/internal/shell/runner"
)

type Cmd struct {
	runner.Runner
}

func (f Cmd) Fzf(output io.Writer) (runner.WriteCloser, error) {
	// Keep frequency-ranked input order even while filtering with a query.
	fzfCmd := exec.Command("fzf", "--reverse", "--height", "70%", "--tmux", "70%", "--no-sort")
	fzfCmd.Stdout = output
	fzfCmd.Stderr = os.Stderr
	waiter, err := f.Start(fzfCmd)
	return waiter, err

}
