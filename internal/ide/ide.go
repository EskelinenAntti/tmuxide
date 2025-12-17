package ide

import (
	"errors"
	"os"

	"github.com/eskelinenantti/tmuxide/internal/project"
)

type Tmux interface {
	HasSession(name string) bool
	New(session string, dir string, args Window) error
	NewWindow(session string, dir string, args Window) error
	Attach(session string) error
	Switch(session string) error
	Kill(session string) error
}

type Window []string

type ShellPath interface {
	Contains(path string) bool
}

var ErrTmuxNotInstalled = errors.New("tmux not installed")

func Start(command []string, project project.Project, tmux Tmux, path ShellPath) error {
	if !path.Contains("tmux") {
		return ErrTmuxNotInstalled
	}

	windows := []Window{}
	if len(command) > 0 {
		windows = []Window{command}
	}

	if tmux.HasSession(project.Name) {
		if err := tmux.Kill(project.Name); err != nil {
			return err
		}
	}

	if err := create(project, windows, tmux); err != nil {
		return err
	}

	if isAttached() {
		return tmux.Switch(project.Name)
	}

	return tmux.Attach(project.Name)
}

func create(project project.Project, windows []Window, tmux Tmux) error {
	var mainWindow = Window{}
	var otherWindows = []Window{}
	if len(windows) > 0 {
		mainWindow = windows[0]
		otherWindows = windows[1:]
	}

	if err := tmux.New(project.Name, project.WorkingDir, mainWindow); err != nil {
		return err
	}

	for _, window := range otherWindows {
		if err := tmux.NewWindow(project.Name, project.WorkingDir, window); err != nil {
			return err
		}
	}

	return nil
}

func isAttached() bool {
	_, isAttached := os.LookupEnv("TMUX")
	return isAttached
}
