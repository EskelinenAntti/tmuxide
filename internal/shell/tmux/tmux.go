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
	if !t.HasSession(project.Name) {
		err = t.New(project.Name, project.WorkingDir, command)
	} else if len(command) > 0 {
		err = t.NewWindow(project.Name, project.WorkingDir, command)
	}
	if err != nil {
		return err
	}
	if _, attached := os.LookupEnv("TMUX"); attached {
		return t.Switch(project.Name)
	}
	return t.Attach(project.Name)
}

func (t Cmd) HasSession(session string) bool {
	tmuxCmd := tmuxCommand("has-session", Args{TargetSession: session})
	return t.Run(tmuxCmd) == nil
}

func (t Cmd) New(session string, dir string, cmd []string) error {
	tmuxCmd := tmuxCommand("new-session", Args{SessionName: session, Detach: true, WorkingDir: dir, Command: cmd})
	return t.Run(tmuxCmd)
}

func (t Cmd) NewWindow(session string, workingDir string, cmd []string) error {
	tmuxCmd := tmuxCommand("new-window", Args{WorkingDir: workingDir, TargetSession: session, Command: cmd})
	return t.Run(tmuxCmd)
}

func (t Cmd) Attach(session string) error {
	tmuxCmd := tmuxCommand("attach", Args{TargetSession: session})
	tmuxCmd.Stdin = os.Stdin
	tmuxCmd.Stdout = os.Stdout
	tmuxCmd.Stderr = os.Stderr
	return t.Run(tmuxCmd)
}

func (t Cmd) Switch(session string) error {
	tmuxCmd := tmuxCommand("switch-client", Args{TargetSession: session})
	return t.Run(tmuxCmd)
}

type Args struct {
	TargetSession string
	Detach        bool
	SessionName   string
	WorkingDir    string
	Command       []string
}

func (a Args) Parse() []string {
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

func tmuxCommand(subCommand string, args Args) *exec.Cmd {
	cmd := exec.Command("tmux", subCommand)
	cmd.Args = append(cmd.Args, args.Parse()...)
	return cmd
}
