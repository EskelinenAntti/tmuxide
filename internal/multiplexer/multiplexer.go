package multiplexer

import (
	"fmt"

	"github.com/eskelinenantti/tmuxide/internal/project"
	"github.com/eskelinenantti/tmuxide/internal/shell"
	"github.com/eskelinenantti/tmuxide/internal/shell/herdr"
	"github.com/eskelinenantti/tmuxide/internal/shell/path"
	"github.com/eskelinenantti/tmuxide/internal/shell/runner"
	"github.com/eskelinenantti/tmuxide/internal/shell/tmux"
)

type Multiplexer interface {
	Open(project.Project, []string) error
}

func Select(name string, path path.ShellPath, runner runner.Runner) (Multiplexer, error) {
	if name == "auto" {
		name = "tmux"
		if herdr.InSession() {
			name = "herdr"
		}
	}
	switch name {
	case "tmux":
		if !path.Contains("tmux") {
			return nil, shell.NotInstalledError{Cmd: "tmux"}
		}
		return tmux.Cmd{Runner: runner}, nil
	case "herdr":
		cmd, err := herdr.New(path, runner)
		if err != nil {
			return nil, err
		}
		return cmd, nil
	default:
		return nil, fmt.Errorf("unknown multiplexer %q (use auto, tmux, or herdr)", name)
	}
}
