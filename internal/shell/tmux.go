package shell

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

var ErrTmuxCommand = errors.New("command tmux")

type Tmux struct{}

func (Tmux) HasSession(session string, window string) bool {
	cmd := exec.Command("tmux", "has-session", "-t", fmt.Sprintf("%s:%s", session, window))
	return cmd.Run() == nil
}

func (Tmux) HasWindow(session string, window string) bool {
	cmd := exec.Command("tmux", "has-session", "-t", session)
	return cmd.Run() == nil
}

func (Tmux) New(session string, dir string, cmd []string) error {
	args := []string{"new-session", "-ds", session, "-c", dir}
	args = append(args, cmd...)
	return run(exec.Command("tmux", args...))
}

func (Tmux) NewWindow(session string, window string, workingDir string, name string, cmd []string) error {
	args := []string{"new-window", "-k", "-c", workingDir, "-t", fmt.Sprintf("%s:%s", session, window)}

	if name != "" {
		args = append(args, "-n", name)
	}

	args = append(args, cmd...)
	return run(exec.Command("tmux", args...))
}

func (Tmux) Attach(session string) error {
	tmuxCmd := exec.Command("tmux", "attach", "-t", session)
	return attach(tmuxCmd)
}

func (Tmux) Switch(session string) error {
	cmd := exec.Command("tmux", "switch-client", "-t", session)
	return attach(cmd)
}

func (Tmux) Kill(session string) error {
	cmd := exec.Command("tmux", "kill-session", "-t", session)
	return run(cmd)
}

func attach(cmd *exec.Cmd) error {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return run(cmd)
}

func run(cmd *exec.Cmd) error {
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%w %v: %w", ErrTmuxCommand, cmd.Args, err)
	}
	return nil
}
