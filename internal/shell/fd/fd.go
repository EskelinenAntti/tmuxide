package fd

import (
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
	return f.Stdout(exec.Command("fd", args...))
}
