package ide

import (
	"os"

	"github.com/eskelinenantti/tmuxide/internal/tmux"
)

func Start(project Project, tmux tmux.Command) error {
	if !tmux.HasSession(project.Name) {
		if err := create(project, project.Windows, tmux); err != nil {
			return err
		}
	}

	if isAttached() {
		return tmux.Switch(project.Name)
	}

	return tmux.Attach(project.Name)
}

func create(project Project, windows []tmux.WindowCommand, tmux tmux.Command) error {
	if err := tmux.New(project.Name, project.Root, windows[0]); err != nil {
		return err
	}

	for _, window := range windows[1:] {
		if err := tmux.NewWindow(project.Name, project.Root, window); err != nil {
			return err
		}
	}

	return nil
}

func isAttached() bool {
	_, isAttached := os.LookupEnv("TMUX")
	return isAttached
}
