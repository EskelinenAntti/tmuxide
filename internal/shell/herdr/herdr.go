package herdr

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/eskelinenantti/tmuxide/internal/project"
	"github.com/eskelinenantti/tmuxide/internal/shell/runner"
)

type Cmd struct {
	runner.Runner
	Binary string
}

func (h Cmd) existingWorkspace(project project.Project) (string, error) {
	workspaces, err := h.listWorkspaces()
	if err != nil {
		return "", err
	}
	id, err := workspaceID(workspaces, func(w workspace) bool {
		return w.Label == project.Name
	})
	if id != "" || err != nil {
		return id, err
	}
	if len(workspaces) == 0 || project.WorkingDir == "" {
		return "", nil
	}
	panes, err := h.listPanes()
	if err != nil {
		return "", err
	}
	workingDir := filepath.Clean(project.WorkingDir)
	matchingWorkspaces := make(map[string]bool)
	for _, pane := range panes {
		if pane.Cwd != "" && filepath.Clean(pane.Cwd) == workingDir {
			matchingWorkspaces[pane.WorkspaceID] = true
		}
	}
	return workspaceID(workspaces, func(w workspace) bool {
		return matchingWorkspaces[w.ID]
	})
}

func workspaceID(workspaces []workspace, matches func(workspace) bool) (string, error) {
	for _, workspace := range workspaces {
		if !matches(workspace) {
			continue
		}
		if workspace.ID == "" {
			return "", fmt.Errorf("herdr workspace has no id")
		}
		return workspace.ID, nil
	}
	return "", nil
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
