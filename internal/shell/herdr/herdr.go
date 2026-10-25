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
	list, err := h.call("workspace", "list")
	if err != nil {
		return "", fmt.Errorf("cannot list workspaces; start a Herdr session first: %w", err)
	}
	if list.Type != "workspace_list" || list.Workspaces == nil {
		return "", fmt.Errorf("invalid herdr workspace list response")
	}
	id, err := workspaceID(list.Workspaces, func(w workspace) bool {
		return w.Label == project.Name
	})
	if id != "" || err != nil {
		return id, err
	}
	if len(list.Workspaces) == 0 || project.WorkingDir == "" {
		return "", nil
	}
	panes, err := h.call("pane", "list")
	if err != nil {
		return "", err
	}
	if panes.Type != "pane_list" || panes.Panes == nil {
		return "", fmt.Errorf("invalid herdr pane list response")
	}
	workingDir := filepath.Clean(project.WorkingDir)
	matchingWorkspaces := make(map[string]bool)
	for _, pane := range panes.Panes {
		if pane.Cwd != "" && filepath.Clean(pane.Cwd) == workingDir {
			matchingWorkspaces[pane.WorkspaceID] = true
		}
	}
	return workspaceID(list.Workspaces, func(w workspace) bool {
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
	var created result
	if id == "" {
		created, err = h.call("workspace", "create", "--cwd", project.WorkingDir, "--label", filepath.Base(project.WorkingDir), "--no-focus")
		if err != nil {
			return err
		}
		id = created.Workspace.ID
		if id == "" {
			return fmt.Errorf("herdr workspace create response has no workspace id")
		}
	} else if len(command) > 0 {
		created, err = h.call("tab", "create", "--workspace", id, "--cwd", project.WorkingDir, "--no-focus")
		if err != nil {
			return err
		}
	}
	if len(command) > 0 {
		if err := h.runInTab(created, command); err != nil {
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

func (h Cmd) runInTab(created result, command []string) error {
	if created.RootPane.ID == "" || created.Tab.ID == "" {
		return fmt.Errorf("herdr create response has no root pane or tab id")
	}
	// pane run is quiet on success, unlike the JSON-returning create APIs.
	if _, err := h.output("pane", "run", created.RootPane.ID, shellCommand(command)); err != nil {
		return err
	}
	_, err := h.output("tab", "focus", created.Tab.ID)
	return err
}
