package fzf

import (
	"errors"
	"io"
	"os"
	"os/exec"

	"github.com/eskelinenantti/herdr-zipln/internal/shell/runner"
)

var ErrCancelled = errors.New("fzf cancelled")

type Cmd struct {
	runner.Runner
}

func (f Cmd) Fzf(input io.Reader, output io.Writer, candidates string) error {
	executable, err := os.Executable()
	if err != nil {
		return err
	}
	// Filtering and ranking happen in _rank so an exact query preference can
	// precede the normal fuzzy results without disabling relevance for the rest.
	fzfCmd := exec.Command("fzf", "--reverse", "--height", "70%", "--tmux", "70%", "--disabled", "--no-sort", "--print-query", "--with-shell", "sh -c", "--bind", `start:reload("$ZIPLN_EXECUTABLE" _rank -- "$ZIPLN_CANDIDATES" {q}),change:reload("$ZIPLN_EXECUTABLE" _rank -- "$ZIPLN_CANDIDATES" {q})`)
	fzfCmd.Env = append(os.Environ(), "ZIPLN_EXECUTABLE="+executable, "ZIPLN_CANDIDATES="+candidates)
	fzfCmd.Stdin = input
	fzfCmd.Stdout = output
	fzfCmd.Stderr = os.Stderr
	err = f.Run(fzfCmd)
	if fzfCmd.ProcessState != nil && fzfCmd.ProcessState.ExitCode() == 130 {
		return ErrCancelled
	}
	return err
}
