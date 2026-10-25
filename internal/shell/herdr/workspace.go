package herdr

import (
	"fmt"
	"path/filepath"

	"github.com/eskelinenantti/tmuxide/internal/project"
)

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
