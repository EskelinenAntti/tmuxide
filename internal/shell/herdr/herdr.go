package herdr

import (
	"os"
	"os/exec"

	"github.com/eskelinenantti/tmuxide/internal/project"
	"github.com/eskelinenantti/tmuxide/internal/shell/runner"
)

type Cmd struct {
	runner.Runner
	Binary string
}

func (h Cmd) Open(project project.Project, command []string) error {
	id, err := h.existingWorkspace(project)
	if err != nil {
		return err
	}
	var tab commandTab
	if id == "" {
		id, tab, err = h.createWorkspace(project)
		if err != nil {
			return err
		}
	} else if len(command) > 0 {
		tab, err = h.createTab(id, project.WorkingDir)
		if err != nil {
			return err
		}
	}
	if len(command) > 0 {
		if err := h.runInTab(tab, command); err != nil {
			return err
		}
	}
	if _, err := h.output("workspace", "focus", id); err != nil {
		return err
	}
	if InSession() {
		return nil
	}
	cmd := exec.Command(h.Binary)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return h.Run(cmd)
}
