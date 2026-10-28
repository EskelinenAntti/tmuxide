package fd

import (
	"context"
	"io"
	"os"
	"os/exec"

	"github.com/eskelinenantti/tmuxide/internal/shell/runner"
)

type Cmd struct {
	runner.Runner
}

func (f Cmd) Fd(ctx context.Context) (io.ReadCloser, error) {
	args := []string{
		"--follow",
		"--hidden",
		"--exclude", "{.git,node_modules,Library}",
		".",
		"--base-directory",
		os.Getenv("HOME"),
		"--",
	}
	return f.Read(exec.CommandContext(ctx, "fd", args...))
}
