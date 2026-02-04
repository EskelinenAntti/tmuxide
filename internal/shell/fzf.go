package shell

import (
	"io"
	"os"
	"os/exec"
)

type FzfCmd struct {
	Runner
}

func (f FzfCmd) Fzf(output io.Writer) (WriteCloser, error) {
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
