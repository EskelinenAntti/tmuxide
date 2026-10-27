package fd

import (
	"os"
	"os/exec"
	"strings"

	"github.com/eskelinenantti/tmuxide/internal/shell/runner"
)

type Cmd struct {
	runner.Runner
}

func (f Cmd) Fd(onPath func(string) error, roots ...string) error {
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
	fdCmd.Stdout = &pathWriter{onPath: onPath}

	return f.Run(fdCmd)
}

type pathWriter struct {
	onPath  func(string) error
	pending string
}

func (w *pathWriter) Write(p []byte) (int, error) {
	w.pending += string(p)
	for {
		path, rest, ok := strings.Cut(w.pending, "\n")
		if !ok {
			return len(p), nil
		}
		w.pending = rest
		if err := w.onPath(path); err != nil {
			return 0, err
		}
	}
}
