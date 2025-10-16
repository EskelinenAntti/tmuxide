package tmux

import (
	"errors"
	"os"
	"os/exec"
)

type Session interface {
	Exists() (bool, error)
	New() error
	Attach() error
	Switch() error
	IsActive() bool
}

type AttachedSession struct {
	Session    string
	WorkingDir string
	Command    string
	Target     string
}

func (tmux *AttachedSession) Exists() (bool, error) {
	cmd := exec.Command("tmux", "has-session", "-t", tmux.Session)

	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			return false, err
		}
		if exitErr.ExitCode() != 1 {
			return false, err
		}

		return false, nil
	}

	return true, nil
}

func (tmux *AttachedSession) New() error {
	cmd := exec.Command("tmux", "new-session", "-ds", tmux.Session, "-c", tmux.WorkingDir, tmux.Command, tmux.Target)
	return cmd.Run()
}

func (*AttachedSession) IsActive() bool {
	_, alreadyInSession := os.LookupEnv("TMUX")
	return alreadyInSession
}

func (tmux *AttachedSession) Attach() error {
	cmd := exec.Command("tmux", "attach", "-t", tmux.Session)
	return attachAndRun(cmd)
}

func (tmux *AttachedSession) Switch() error {
	cmd := exec.Command("tmux", "switch-client", "-t", tmux.Session)
	return attachAndRun(cmd)
}

func Start(tmux Session) error {
	exists, err := tmux.Exists()
	if err != nil {
		return err
	}

	if !exists {
		if err := tmux.New(); err != nil {
			return err
		}
	}

	if tmux.IsActive() {
		return tmux.Switch()
	}

	return tmux.Attach()
}

func attachAndRun(cmd *exec.Cmd) error {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
