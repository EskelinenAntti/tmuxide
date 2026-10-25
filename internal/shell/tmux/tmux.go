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
	sessionExists := t.Run(sessionCommand("has-session", project.Name)) == nil
	if create := creationCommand(project, command, sessionExists); create != nil {
		if err := t.Run(create); err != nil {
			return err
		}
	}
	if _, attached := os.LookupEnv("TMUX"); attached {
		return t.Run(sessionCommand("switch-client", project.Name))
	}
	tmuxCmd := sessionCommand("attach", project.Name)
	tmuxCmd.Stdin = os.Stdin
	tmuxCmd.Stdout = os.Stdout
	tmuxCmd.Stderr = os.Stderr
	return t.Run(tmuxCmd)
}

func creationCommand(project project.Project, command []string, sessionExists bool) *exec.Cmd {
	if sessionExists && len(command) == 0 {
		return nil
	}

	var directoryArgs []string
	if project.WorkingDir != "" {
		directoryArgs = []string{"-c", project.WorkingDir}
	}

	if sessionExists {
		cmd := sessionCommand("new-window", project.Name)
		cmd.Args = append(cmd.Args, directoryArgs...)
		cmd.Args = append(cmd.Args, command...)
		return cmd
	}

	cmd := exec.Command("tmux", "new-session")
	cmd.Args = append(cmd.Args, directoryArgs...)
	cmd.Args = append(cmd.Args, "-d")
	if project.Name != "" {
		cmd.Args = append(cmd.Args, "-s", project.Name)
	}
	cmd.Args = append(cmd.Args, command...)
	return cmd
}

func sessionCommand(subCommand, session string) *exec.Cmd {
	cmd := exec.Command("tmux", subCommand)
	if session != "" {
		cmd.Args = append(cmd.Args, "-t", session+":")
	}
	return cmd
}
