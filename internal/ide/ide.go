package ide

import (
	"os"

	"github.com/eskelinenantti/tmuxide/internal/project"
)

type Tmux interface {
	HasSession(name string) bool
	New(session string, dir string, args Window) error
	NewWindow(session string, dir string, args Window) error
	Attach(session string) error
	Switch(session string) error
}

func Start(project project.Project, tmux Tmux, path ShellPath) error {
	windows, err := Windows(project, path)
	if err != nil {
		return err
	}

	if !tmux.HasSession(project.Name) {
		if err := create(project, windows, tmux); err != nil {
			return err
		}
	}

	if isAttached() {
		return tmux.Switch(project.Name)
	}

	return tmux.Attach(project.Name)
}

func create(project project.Project, windows []Window, tmux Tmux) error {
	if err := tmux.New(project.Name, project.WorkingDir, windows[0]); err != nil {
		return err
	}

	for _, window := range windows[1:] {
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
