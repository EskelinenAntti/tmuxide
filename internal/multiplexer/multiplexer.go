// Package multiplexer selects a terminal host without exposing its lifecycle to
// project selection or editor configuration.
package multiplexer

import (
	"fmt"
	"os"

	"github.com/eskelinenantti/tmuxide/internal/ide"
	"github.com/eskelinenantti/tmuxide/internal/shell"
	"github.com/eskelinenantti/tmuxide/internal/shell/herdr"
	"github.com/eskelinenantti/tmuxide/internal/shell/path"
	"github.com/eskelinenantti/tmuxide/internal/shell/runner"
	"github.com/eskelinenantti/tmuxide/internal/shell/tmux"
)

func Select(name string, path path.ShellPath, runner runner.Runner) (ide.Multiplexer, error) {
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
	var mux ide.Multiplexer
	var binary string
	switch name {
	case "tmux":
		binary = "tmux"
		mux = tmux.Cmd{Runner: runner}
	case "herdr":
		binary = os.Getenv("HERDR_BIN_PATH")
		if binary == "" {
			binary = "herdr"
		}
		mux = herdr.Cmd{Runner: runner, Binary: binary}
	default:
		return nil, fmt.Errorf("unknown backend %q (use auto, tmux, or herdr)", name)
	}
	if !path.Contains(binary) {
		return nil, shell.NotInstalledError{Cmd: binary}
	}
	return mux, nil
}
