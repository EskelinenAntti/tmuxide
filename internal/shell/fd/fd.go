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
	shallowArgs := append([]string{"--max-depth", "2"}, args...)
	shallow, err := f.Stdout(exec.Command("fd", shallowArgs...))
	if err != nil {
		return nil, err
	}
	deepArgs := append([]string{"--min-depth", "3"}, args...)
	deep, err := f.Stdout(exec.Command("fd", deepArgs...))
	if err != nil {
		shallow.Close()
		return nil, err
	}
	return &scanOutput{Reader: io.MultiReader(shallow, deep), shallow: shallow, deep: deep}, nil
}

type scanOutput struct {
	io.Reader
	shallow io.ReadCloser
	deep    io.ReadCloser
}

func (s *scanOutput) Close() error {
	return errors.Join(s.shallow.Close(), s.deep.Close())
}
