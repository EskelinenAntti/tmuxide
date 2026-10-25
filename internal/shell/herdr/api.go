package herdr

import (
	"fmt"
	"path/filepath"

	"github.com/eskelinenantti/tmuxide/internal/project"
)

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

func (h Cmd) createWorkspace(project project.Project) (string, commandTab, error) {
	created, err := call[struct {
		commandTabResult
		Workspace workspace `json:"workspace"`
	}](h, "workspace", "create", "--cwd", project.WorkingDir, "--label", filepath.Base(project.WorkingDir), "--no-focus")
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
