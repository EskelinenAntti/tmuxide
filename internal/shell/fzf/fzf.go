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
	args := []string{
		"--reverse",
		"--height",
		"30%",
	}
	fzfCmd := exec.Command("fzf", args...)
	fzfCmd.Stdout = output
	fzfCmd.Stderr = os.Stderr
	waiter, err := f.Start(fzfCmd)
	return waiter, err

}
