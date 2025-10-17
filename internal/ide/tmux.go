package ide

import (
	"fmt"
	"os"
	"os/exec"
)

func (project *Project) Exists() bool {
	cmd := exec.Command("tmux", "has-session", "-t", project.Name)
	return cmd.Run() == nil
}

func (project *Project) New() error {
	window := project.Windows[0]
	args := []string{"new-session", "-ds", project.Name, "-c", project.Root}

	args = append(args, window.Cmd)
	args = append(args, window.Args...)

	cmd := exec.Command("tmux", args...)
	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("Failed to create session: %w", err)
	}
	return nil
}

func (project *Project) Attach() error {
	cmd := exec.Command("tmux", "attach", "-t", project.Name)
	err := attachAndRun(cmd)
	if err != nil {
		return fmt.Errorf("Failed to attach session: %w", err)
	}
	return nil
}

func (project *Project) Switch() error {
	cmd := exec.Command("tmux", "switch-client", "-t", project.Name)
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
