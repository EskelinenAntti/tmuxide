package herdr

import (
	"fmt"
	"path/filepath"

	"github.com/eskelinenantti/tmuxide/internal/project"
)

type workspace struct {
	ID    string `json:"workspace_id"`
	Label string `json:"label"`
}

type pane struct {
	WorkspaceID string `json:"workspace_id"`
	Cwd         string `json:"cwd"`
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

func (h Cmd) listWorkspaces() ([]workspace, error) {
	list, err := call[struct {
		Type       string      `json:"type"`
		Workspaces []workspace `json:"workspaces"`
	}](h, "workspace", "list")
	if err != nil {
		return nil, fmt.Errorf("cannot list workspaces; start a Herdr session first: %w", err)
	}
	if list.Type != "workspace_list" || list.Workspaces == nil {
		return nil, fmt.Errorf("invalid herdr workspace list response")
	}
	return list.Workspaces, nil
}

func (h Cmd) listPanes() ([]pane, error) {
	list, err := call[struct {
		Type  string `json:"type"`
		Panes []pane `json:"panes"`
	}](h, "pane", "list")
	if err != nil {
		return nil, err
	}
	if list.Type != "pane_list" || list.Panes == nil {
		return nil, fmt.Errorf("invalid herdr pane list response")
	}
	return list.Panes, nil
}
