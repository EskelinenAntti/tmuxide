package fd

import (
	"io"
	"os"
	"os/exec"

	"github.com/eskelinenantti/tmuxide/internal/shell/runner"
)

type Cmd struct {
	runner.Runner
}

func (f Cmd) Fd(output io.Writer, roots ...string) error {
	args := []string{
		"--follow",
		"--hidden",
		"--exclude", "{.git,node_modules,Library}",
		".",
		"--base-directory",
		os.Getenv("HOME"),
	}
	args = append(args, "--")
	args = append(args, roots...)

	fdCmd := exec.Command("fd", args...)
	fdCmd.Stdout = output

	return f.Run(fdCmd)
}
