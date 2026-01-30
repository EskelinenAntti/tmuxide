package tmux

import (
	"errors"
	"fmt"
	"strings"
)

var ErrTmuxNotInstalled = errors.New("tmux not installed")

type Parser interface {
	Parse() []string
}

type Runner interface {
	Run(name string, args Parser) error
	Attach(name string, args Parser) error
}

type Tmux struct {
	Runner
}

func (t Tmux) HasSession(targetSession string, targetWindow string) bool {
	return t.Run("has-session", Args{TargetSession: targetSession, TargetWindow: targetWindow}) == nil
}

func (t Tmux) New(session string, dir string, cmd []string) error {
	return t.Run("new-session", Args{SessionName: session, Detach: true, WorkingDir: dir, Command: cmd})
}

func (t Tmux) NewWindow(session string, window string, workingDir string, name string, cmd []string) error {
	return t.Run("new-window", Args{Kill: true, WindowName: name, WorkingDir: workingDir, TargetSession: session, TargetWindow: window, Command: cmd})
}

func (t Tmux) Attach(session string) error {
	return t.Runner.Attach("attach", Args{TargetSession: session})
}

func (t Tmux) Switch(session string) error {
	return t.Runner.Attach("switch-client", Args{TargetSession: session})
}

func (t Tmux) Kill(session string) error {
	return t.Run("kill-session", Args{TargetSession: session})
}

func (t Tmux) ChooseSession() error {
	return t.Run("choose-session", Args{})
}

func (t Tmux) KillSession() error {
	return t.Run("kill-session", Args{})
}

type Args struct {
	TargetSession string
	TargetWindow  string
	Detach        bool
	SessionName   string
	WindowName    string
	WorkingDir    string
	Command       []string
	Kill          bool
}

func (a Args) Parse() []string {
	args := []string{}

	if a.TargetSession != "" || a.TargetWindow != "" {
		args = append(args, "-t", fmt.Sprintf("%s:%s", a.TargetSession, a.TargetWindow))
	}

	if a.WorkingDir != "" {
		args = append(args, "-c", a.WorkingDir)
	}

	if a.Detach {
		args = append(args, "-d")
	}

	if a.Kill {
		args = append(args, "-k")
	}

	if a.SessionName != "" {
		args = append(args, "-s", a.SessionName)
	}

	if a.WindowName != "" {
		args = append(args, "-n", a.WindowName)
	}

	if len(a.Command) > 0 {
		args = append(args, a.Command...)
	}

	return args
}

func (a Args) String() string {
	return strings.Join(a.Parse(), " ")
}

type ShellPath interface {
	Contains(path string) bool
}

func InitTmux(path ShellPath, tmuxRunner Runner) (Tmux, error) {
	if !path.Contains("tmux") {
		return Tmux{}, ErrTmuxNotInstalled
	}

	return Tmux{Runner: tmuxRunner}, nil
}
