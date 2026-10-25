package herdr

import (
	"fmt"
	"path/filepath"
)

type workspace struct {
	ID    string `json:"workspace_id"`
	Label string `json:"label"`
}

type pane struct {
	WorkspaceID string `json:"workspace_id"`
	Cwd         string `json:"cwd"`
}

type commandTabResult struct {
	Tab struct {
		ID string `json:"tab_id"`
	} `json:"tab"`
	RootPane struct {
		ID string `json:"pane_id"`
	} `json:"root_pane"`
}

type commandTab struct {
	tabID  string
	paneID string
}

func (r commandTabResult) commandTab() commandTab {
	return commandTab{tabID: r.Tab.ID, paneID: r.RootPane.ID}
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

func (h Cmd) createWorkspace(workingDir string) (string, commandTab, error) {
	created, err := call[struct {
		commandTabResult
		Workspace workspace `json:"workspace"`
	}](h, "workspace", "create", "--cwd", workingDir, "--label", filepath.Base(workingDir), "--no-focus")
	if err != nil {
		return "", commandTab{}, err
	}
	if created.Workspace.ID == "" {
		return "", commandTab{}, fmt.Errorf("herdr workspace create response has no workspace id")
	}
	return created.Workspace.ID, created.commandTab(), nil
}

func (h Cmd) createTab(workspaceID, workingDir string) (commandTab, error) {
	created, err := call[commandTabResult](h, "tab", "create", "--workspace", workspaceID, "--cwd", workingDir, "--no-focus")
	if err != nil {
		return commandTab{}, err
	}
	return created.commandTab(), nil
}

func (h Cmd) runInTab(tab commandTab, command []string) error {
	if tab.paneID == "" || tab.tabID == "" {
		return fmt.Errorf("herdr create response has no root pane or tab id")
	}
	// pane run is quiet on success, unlike the JSON-returning create APIs.
	if _, err := h.output("pane", "run", tab.paneID, shellCommand(command)); err != nil {
		return err
	}
	_, err := h.output("tab", "focus", tab.tabID)
	return err
}
