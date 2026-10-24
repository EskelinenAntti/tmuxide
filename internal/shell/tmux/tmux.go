package tmux

import (
	"os"
	"os/exec"

	"github.com/eskelinenantti/tmuxide/internal/project"
	"github.com/eskelinenantti/tmuxide/internal/shell/runner"
)

type Cmd struct {
	runner.Runner
}

func (t Cmd) Open(project project.Project, command []string) error {
	var err error
	if !t.hasSession(project.Name) {
		err = t.new(project.Name, project.WorkingDir, command)
	} else if len(command) > 0 {
		err = t.newWindow(project.Name, project.WorkingDir, command)
	}
	if err != nil {
		return err
	}
	if _, attached := os.LookupEnv("TMUX"); attached {
		return t.switchSession(project.Name)
	}
	return t.attach(project.Name)
}

func (t Cmd) hasSession(session string) bool {
	tmuxCmd := tmuxCommand("has-session", args{TargetSession: session})
	return t.Run(tmuxCmd) == nil
}

func (t Cmd) new(session string, dir string, cmd []string) error {
	tmuxCmd := tmuxCommand("new-session", args{SessionName: session, Detach: true, WorkingDir: dir, Command: cmd})
	return t.Run(tmuxCmd)
}

func (t Cmd) newWindow(session string, workingDir string, cmd []string) error {
	tmuxCmd := tmuxCommand("new-window", args{WorkingDir: workingDir, TargetSession: session, Command: cmd})
	return t.Run(tmuxCmd)
}

func (t Cmd) attach(session string) error {
	tmuxCmd := tmuxCommand("attach", args{TargetSession: session})
	tmuxCmd.Stdin = os.Stdin
	tmuxCmd.Stdout = os.Stdout
	tmuxCmd.Stderr = os.Stderr
	return t.Run(tmuxCmd)
}

func (t Cmd) switchSession(session string) error {
	tmuxCmd := tmuxCommand("switch-client", args{TargetSession: session})
	return t.Run(tmuxCmd)
}

type args struct {
	TargetSession string
	Detach        bool
	SessionName   string
	WorkingDir    string
	Command       []string
}

func (a args) parse() []string {
	args := []string{}

	if a.TargetSession != "" {
		args = append(args, "-t", a.TargetSession+":")
	}

	if a.WorkingDir != "" {
		args = append(args, "-c", a.WorkingDir)
	}

	if a.Detach {
		args = append(args, "-d")
	}

	if a.SessionName != "" {
		args = append(args, "-s", a.SessionName)
	}

	if len(a.Command) > 0 {
		args = append(args, a.Command...)
	}

	return args
}

func tmuxCommand(subCommand string, args args) *exec.Cmd {
	cmd := exec.Command("tmux", subCommand)
	cmd.Args = append(cmd.Args, args.parse()...)
	return cmd
}
