package shell

import (
	"io"
	"os"
	"os/exec"
)

type FzfCmd struct {
	Runner
}

func (f FzfCmd) Fzf(input io.Reader, output io.Writer) Waitable {
	args := []string{
		"--reverse",
		"--height",
		"30%",
	}
	fzfCmd := exec.Command("fzf", args...)
	fzfCmd.Stdin = input
	fzfCmd.Stdout = output
	fzfCmd.Stderr = os.Stderr
	return fzfCmd
}
