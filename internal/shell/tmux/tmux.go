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
	var create *exec.Cmd
	sessionExists := t.Run(sessionCommand("has-session", project.Name)) == nil
	if !sessionExists {
		create = exec.Command("tmux", "new-session")
	} else if len(command) > 0 {
		create = sessionCommand("new-window", project.Name)
	}
	if create != nil {
		if project.WorkingDir != "" {
			create.Args = append(create.Args, "-c", project.WorkingDir)
		}
		if !sessionExists {
			create.Args = append(create.Args, "-d")
			if project.Name != "" {
				create.Args = append(create.Args, "-s", project.Name)
			}
		}
		create.Args = append(create.Args, command...)
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

func sessionCommand(subCommand, session string) *exec.Cmd {
	cmd := exec.Command("tmux", subCommand)
	if session != "" {
		cmd.Args = append(cmd.Args, "-t", session+":")
	}
	return cmd
}
