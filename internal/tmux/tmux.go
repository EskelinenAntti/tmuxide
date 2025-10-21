package tmux

import (
	"errors"
	"os"
	"os/exec"

	"github.com/eskelinenantti/tmuxide/internal/path"
)

type WindowCommand struct {
	Cmd  string
	Args []string
}

type Command interface {
	HasSession(name string) bool
	New(session string, dir string, cmd WindowCommand) error
	NewWindow(session string, dir string, cmd WindowCommand) error
	Attach(session string) error
	Switch(session string) error
}

type ShellTmux struct{}

var ErrTmuxNotInPath = errors.New(
	"Did not find tmux, which is a required dependency for ide command.\n\n" +

		"You can install tmux e.g. via homebrew by running\n" +
		"brew install tmux\n",
)

func (tmux ShellTmux) HasSession(session string) bool {
	cmd := exec.Command("tmux", "has-session", "-t", session)
	return cmd.Run() == nil
}

func (ShellTmux) New(session string, dir string, cmd WindowCommand) error {
	tmuxCmd := tmuxCommand([]string{"new-session", "-ds", session, "-c", dir}, cmd)
	return tmuxCmd.Run()
}

func (ShellTmux) NewWindow(session string, dir string, cmd WindowCommand) error {
	tmuxCmd := tmuxCommand([]string{"new-window", "-d", "-t", session, "-c", dir}, cmd)
	return tmuxCmd.Run()
}

func (ShellTmux) Attach(session string) error {
	tmuxCmd := exec.Command("tmux", "attach", "-t", session)
	return attachAndRun(tmuxCmd)
}

func (ShellTmux) Switch(session string) error {
	cmd := exec.Command("tmux", "switch-client", "-t", session)
	return attachAndRun(cmd)
}

func EnsureInstalled(path path.Path) error {
	if !path.Contains("tmux") {
		return ErrTmuxNotInPath
	}
	return nil
}

func attachAndRun(cmd *exec.Cmd) error {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func tmuxCommand(tmuxArgs []string, cmd WindowCommand) *exec.Cmd {
	tmuxArgs = append(tmuxArgs, cmd.Cmd)
	tmuxArgs = append(tmuxArgs, cmd.Args...)

	return exec.Command("tmux", tmuxArgs...)
}
