package shell

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

type FzfCmd struct {
	Runner
}

func (f FzfCmd) Fzf(input io.Reader) ([]byte, error) {
	args := []string{
		"--reverse",
		"--height",
		"30%",
	}
	fzfCmd := exec.Command("fzf", args...)
	var out bytes.Buffer
	fzfCmd.Stderr = os.Stderr
	fzfCmd.Stdin = input
	fzfCmd.Stdout = &out

	err := f.Run(fzfCmd)
	return out.Bytes(), err
}
