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
var ErrNoSessionsFound = errors.New("no session found")

func Start(command []string, project project.Project, tmuxRunner tmux.Runner, path ShellPath) error {
	if !path.Contains("tmux") {
		return ErrTmuxNotInstalled
	}

	tmux := tmux.Tmux{Runner: tmuxRunner}
	var err error
	if len(command) == 0 {
		err = startWithoutCommand(tmux, project)
	} else {
		err = startWithCommand(tmux, project, command)
	}

	if err != nil {
		return err
	}

	if isAttached() {
		return tmux.Switch(project.Name)
	}

	return tmux.Attach(project.Name)
}

func List(tmuxRunner tmux.Runner) error {
	tmux := tmux.Tmux{Runner: tmuxRunner}

	err := tmux.ChooseSession()
	if err != nil {
		return ErrNoSessionsFound
	}

	if isAttached() {
		return err
	}

	return tmux.Attach("")
}

func Attach(tmuxRunner tmux.Runner) error {
	tmux := tmux.Tmux{Runner: tmuxRunner}
	if isAttached() {
		return nil
	}
	return tmux.Attach("")
}

func Quit(tmuxRunner tmux.Runner) error {
	tmux := tmux.Tmux{Runner: tmuxRunner}
	if !isAttached() {
		return nil
	}
	return tmux.KillSession()
}

func startWithCommand(tmux tmux.Tmux, project project.Project, command []string) error {
	windowName := command[0]

	if tmux.HasSession(project.Name, windowName) {
		return tmux.NewWindow(project.Name, windowName, project.WorkingDir, windowName, command)
	} else if tmux.HasSession(project.Name, "") {
		return tmux.NewWindow(project.Name, "", project.WorkingDir, windowName, command)
	} else {
		return tmux.New(project.Name, project.WorkingDir, command)
	}
}

func startWithoutCommand(tmux tmux.Tmux, project project.Project) error {
	if tmux.HasSession(project.Name, "") {
		// When no command was provided and session exists, don't create any new windows or sessions
		return nil
	} else {
		return tmux.New(project.Name, project.WorkingDir, nil)
	}
}

func isAttached() bool {
	_, isAttached := os.LookupEnv("TMUX")
	return isAttached
}
