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

func (f Cmd) Fzf(input io.Reader, output io.Writer) error {
	executable, err := os.Executable()
	if err != nil {
		return err
	}
	// Stream the scan immediately. Track a restored selection as more results
	// arrive, and release it on query changes so normal fuzzy ranking resumes.
	fzfCmd := exec.Command("fzf",
		"--reverse", "--height", "70%", "--tmux", "70%",
		"--sort", "--scheme=default", "--tiebreak=length,index", "--print-query",
		"--with-shell", "sh -c",
		"--bind", `start,change:untrack-current+first+rebind(result),result:transform("$ZIPLN_EXECUTABLE" _focus -- {q} {*f})`,
	)
	fzfCmd.Env = append(os.Environ(), "ZIPLN_EXECUTABLE="+executable)
	fzfCmd.Stdin = input
	fzfCmd.Stdout = output
	fzfCmd.Stderr = os.Stderr
	err = f.Run(fzfCmd)
	if fzfCmd.ProcessState != nil && fzfCmd.ProcessState.ExitCode() == 130 {
		return ErrCancelled
	}
	return err
}
