package tmux

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/eskelinenantti/tmuxide/internal/ide"
)

type Session struct {
	ide.Project
}

func (session Session) Exists() bool {
	cmd := exec.Command("tmux", "has-session", "-t", session.Name)
	return cmd.Run() == nil
}

func (session Session) New() error {
	window := session.Windows[0]
	args := []string{"new-session", "-ds", session.Name, "-c", session.Root}

	args = append(args, window.Cmd)
	args = append(args, window.Args...)

	cmd := exec.Command("tmux", args...)
	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("Failed to create session: %w", err)
	}
	return nil
}

func (session Session) Attach() error {
	cmd := exec.Command("tmux", "attach", "-t", session.Name)
	err := attachAndRun(cmd)
	if err != nil {
		return fmt.Errorf("Failed to attach session: %w", err)
	}
	return nil
}

func (session Session) Switch() error {
	cmd := exec.Command("tmux", "switch-client", "-t", session.Name)
	err := attachAndRun(cmd)
	if err != nil {
		return fmt.Errorf("Failed to switch session: %w", err)
	}
	return nil
}

func SessionFor(project ide.Project) (Session, error) {
	if err := EnsureInstalled(); err != nil {
		return Session{}, err
	}
	return Session{Project: project}, nil
}

func EnsureInstalled() error {
	if _, err := exec.LookPath("tmux"); err != nil {
		return errors.New(
			"Did not find tmux, which is a required dependency for ide command.\n\n" +

				"You can install tmux e.g. via homebrew by running\n" +
				"brew install tmux\n",
		)
	}
	return nil
}

func attachAndRun(cmd *exec.Cmd) error {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
