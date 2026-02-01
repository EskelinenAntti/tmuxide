package tmux

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/eskelinenantti/tmuxide/internal/shell"
)

var ErrTmuxNotInstalled = errors.New("tmux not installed")

type Tmux struct {
	shell.Runner
}

func (t Tmux) HasSession(targetSession string, targetWindow string) bool {
	tmuxCmd := tmuxCommand("has-session", Args{TargetSession: targetSession, TargetWindow: targetWindow})
	return t.Run(*tmuxCmd) == nil
}

func (t Tmux) New(session string, dir string, cmd []string) error {
	tmuxCmd := tmuxCommand("new-session", Args{SessionName: session, Detach: true, WorkingDir: dir, Command: cmd})
	return t.Run(*tmuxCmd)
}

func (t Tmux) NewWindow(session string, window string, workingDir string, name string, cmd []string) error {
	tmuxCmd := tmuxCommand("new-window", Args{Kill: true, WindowName: name, WorkingDir: workingDir, TargetSession: session, TargetWindow: window, Command: cmd})
	return t.Run(*tmuxCmd)
}

func (t Tmux) Attach(session string) error {
	tmuxCmd := tmuxCommand("attach", Args{TargetSession: session})
	tmuxCmd.Stdin = os.Stdin
	tmuxCmd.Stdout = os.Stdout
	tmuxCmd.Stderr = os.Stderr
	return t.Run(*tmuxCmd)
}

func (t Tmux) Switch(session string) error {
	tmuxCmd := tmuxCommand("switch-client", Args{TargetSession: session})
	return t.Run(*tmuxCmd)
}

func (t Tmux) Kill(session string) error {
	tmuxCmd := tmuxCommand("kill-session", Args{TargetSession: session})
	return t.Run(*tmuxCmd)
}

func (t Tmux) ChooseSession() error {
	tmuxCmd := tmuxCommand("choose-session", Args{})
	return t.Run(*tmuxCmd)
}

func (t Tmux) ListSessions() ([]byte, error) {
	tmuxCmd := tmuxCommand("list-sessions", Args{Format: "#S"})
	var out bytes.Buffer
	tmuxCmd.Stdout = &out
	err := t.Run(*tmuxCmd)
	return out.Bytes(), err
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
	Format        string
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

	if a.Format != "" {
		args = append(args, "-F", a.Format)
	}

	if len(a.Command) > 0 {
		args = append(args, a.Command...)
	}

	return args
}

func (a Args) String() string {
	return strings.Join(a.Parse(), " ")
}

func InitTmux(path shell.ShellPath, tmuxRunner shell.Runner) (Tmux, error) {
	if !path.Contains("tmux") {
		return Tmux{}, ErrTmuxNotInstalled
	}

	return Tmux{Runner: tmuxRunner}, nil
}

func tmuxCommand(subCommand string, args Args) *exec.Cmd {
	cmd := exec.Command("tmux", subCommand)
	cmd.Args = append(cmd.Args, args.Parse()...)
	return cmd
}
