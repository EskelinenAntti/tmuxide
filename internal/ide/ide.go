package ide

import (
	"errors"
	"os"

	"github.com/eskelinenantti/tmuxide/internal/project"
	"github.com/eskelinenantti/tmuxide/internal/shell/tmux"
)

type ShellPath interface {
	Contains(path string) bool
}

var ErrTmuxNotInstalled = errors.New("tmux not installed")

func Start(command []string, project project.Project, tmuxRunner tmux.Runner, path ShellPath) error {
	if !path.Contains("tmux") {
		return ErrTmuxNotInstalled
	}

	windowName := ""
	if len(command) > 0 {
		windowName = command[0]
	}

	tmux := tmux.Tmux{Runner: tmuxRunner}
	var err error
	if len(command) == 0 {
		if !tmux.HasSession(project.Name, "") {
			err = tmux.New(project.Name, project.WorkingDir, command)
		}
		// When no command was provided and session exists, simply attach to the existing session.
	} else {
		if tmux.HasSession(project.Name, command[0]) {
			err = tmux.NewWindow(project.Name, command[0], project.WorkingDir, windowName, command)
		} else if tmux.HasSession(project.Name, "") {
			err = tmux.NewWindow(project.Name, "", project.WorkingDir, windowName, command)
		} else {
			err = tmux.New(project.Name, project.WorkingDir, command)
		}
	}

	if err != nil {
		return err
	}

	if isAttached() {
		return tmux.Switch(project.Name)
	}

	return tmux.Attach(project.Name)
}

func isAttached() bool {
	_, isAttached := os.LookupEnv("TMUX")
	return isAttached
}
