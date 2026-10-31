package multiplexer

import (
	"fmt"
	"os"

	"github.com/eskelinenantti/zipln/internal/project"
	"github.com/eskelinenantti/zipln/internal/shell"
	"github.com/eskelinenantti/zipln/internal/shell/herdr"
	"github.com/eskelinenantti/zipln/internal/shell/path"
	"github.com/eskelinenantti/zipln/internal/shell/runner"
	"github.com/eskelinenantti/zipln/internal/shell/tmux"
)

type Multiplexer interface {
	Open(project.Project, []string) error
}

func Select(name string, path path.ShellPath, runner runner.Runner) (Multiplexer, error) {
	if name == "auto" {
		switch {
		case os.Getenv("HERDR_ENV") == "1":
			name = "herdr"
		case os.Getenv("TMUX") != "":
			name = "tmux"
		default:
			name = "tmux"
		}
	}
	var mux Multiplexer
	switch name {
	case "tmux":
		mux = tmux.Cmd{Runner: runner}
	case "herdr":
		mux = herdr.Cmd{Runner: runner}
	default:
		return nil, fmt.Errorf("unknown multiplexer %q (use auto, tmux, or herdr)", name)
	}
	if !path.Contains(name) {
		return nil, shell.NotInstalledError{Cmd: name}
	}
	return mux, nil
}
