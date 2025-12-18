package ide

import (
	"errors"
	"os"

	"github.com/eskelinenantti/tmuxide/internal/project"
)

type Tmux interface {
	HasSession(name string) bool
	New(session string, dir string, cmd []string) error
	NewWindow(session string, dir string, cmd []string) error
	Attach(session string) error
	Switch(session string) error
	Kill(session string) error
}

type ShellPath interface {
	Contains(path string) bool
}

var ErrTmuxNotInstalled = errors.New("tmux not installed")

func Start(command []string, project project.Project, tmux Tmux, path ShellPath) error {
	if !path.Contains("tmux") {
		return ErrTmuxNotInstalled
	}

	var err error
	if !tmux.HasSession(project.Name) {
		err = tmux.New(project.Name, project.WorkingDir, command)
	} else {
		err = tmux.NewWindow(project.Name, project.WorkingDir, command)
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
