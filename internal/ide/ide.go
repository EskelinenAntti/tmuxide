package ide

import (
	"os"

	"github.com/eskelinenantti/tmuxide/internal/project"
	"github.com/eskelinenantti/tmuxide/internal/shell/tmux"
)

func Start(command []string, project project.Project, tmux tmux.Cmd) error {
	var err error
	sessionExists := tmux.HasSession(project.Name)

	if !sessionExists {
		err = tmux.New(project.Name, project.WorkingDir, command)
	} else if len(command) > 0 {
		err = tmux.NewWindow(project.Name, project.WorkingDir, command[0], command)
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
