package tmux

import (
	"fmt"
	"os"
	"os/exec"
)

type Session struct {
	Session    string
	WorkingDir string
}

func (tmux *Session) Exists() bool {
	cmd := exec.Command("tmux", "has-session", "-t", tmux.Session)
	return cmd.Run() == nil
}

func (tmux *Session) New(command string, args ...string) error {
	tmuxArgs := append([]string{
		"new-session", "-ds", tmux.Session, "-c", tmux.WorkingDir, command,
	}, args...)

	cmd := exec.Command("tmux", tmuxArgs...)
	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("Failed to create session: %w", err)
	}
	return nil
}

func (tmux *Session) Attach() error {
	cmd := exec.Command("tmux", "attach", "-t", tmux.Session)
	err := attachAndRun(cmd)
	if err != nil {
		return fmt.Errorf("Failed to attach session: %w", err)
	}
	return nil
}

func (tmux *Session) Switch() error {
	cmd := exec.Command("tmux", "switch-client", "-t", tmux.Session)
	err := attachAndRun(cmd)
	if err != nil {
		return fmt.Errorf("Failed to switch session: %w", err)
	}
	return nil
}

func attachAndRun(cmd *exec.Cmd) error {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
