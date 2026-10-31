package fd

import (
	"errors"
	"io"
	"os"
	"os/exec"

	"github.com/eskelinenantti/herdr-zipln/internal/shell/runner"
)

type Cmd struct {
	runner.Runner
}

func (f Cmd) Fd() (io.ReadCloser, error) {
	args := []string{
		"--follow",
		"--hidden",
		"--exclude", "{.git,node_modules,Library}",
		".",
		"--base-directory",
		os.Getenv("HOME"),
	}
	// Bound the priority scan so deep trees cannot delay shallow directories.
	shallowArgs := append([]string{"--type", "d", "--max-depth", "2"}, args...)
	shallow, err := f.Stdout(exec.Command("fd", shallowArgs...))
	if err != nil {
		return nil, err
	}
	all, err := f.Stdout(exec.Command("fd", args...))
	if err != nil {
		shallow.Close()
		return nil, err
	}
	return &scanOutput{Reader: io.MultiReader(shallow, all), shallow: shallow, all: all}, nil
}

// Both scans stream and are stopped when the picker closes its input.
type scanOutput struct {
	io.Reader
	shallow io.ReadCloser
	all     io.ReadCloser
}

func (s *scanOutput) Close() error {
	return errors.Join(s.shallow.Close(), s.all.Close())
}
