package fzf

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/eskelinenantti/tmuxide/internal/shell/runner"
)

type Cmd struct {
	runner.Starter
}

func isCancelled(err error) bool {
	var exitErr *exec.ExitError
	return errors.As(err, &exitErr) && exitErr.ExitCode() == 130
}

// Select streams candidates into fzf and returns the selected candidate.
// Cancelling selection returns an empty string without an error.
func (f Cmd) Select(candidates func(io.Writer) error) (string, error) {
	var output bytes.Buffer
	fzfCmd := exec.Command("fzf", "--reverse", "--height", "70%", "--tmux", "70%")
	fzfCmd.Stdout = &output
	fzfCmd.Stderr = os.Stderr
	waiter, err := f.Start(fzfCmd)
	if err != nil {
		return "", err
	}
	if err := candidates(waiter); err != nil {
		return "", err
	}
	if err := waiter.Close(); err != nil {
		if isCancelled(err) {
			return "", nil
		}
		return "", err
	}
	return strings.TrimSpace(output.String()), nil
}
