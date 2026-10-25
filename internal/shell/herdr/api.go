package herdr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/eskelinenantti/tmuxide/internal/project"
)

type result struct {
	Type       string      `json:"type"`
	Workspaces []workspace `json:"workspaces"`
	Workspace  workspace   `json:"workspace"`
	Panes      []pane      `json:"panes"`
	Tab        struct {
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

func (r result) commandTab() commandTab {
	return commandTab{tabID: r.Tab.ID, paneID: r.RootPane.ID}
}

func (h Cmd) createWorkspace(project project.Project) (string, commandTab, error) {
	created, err := h.call("workspace", "create", "--cwd", project.WorkingDir, "--label", filepath.Base(project.WorkingDir), "--no-focus")
	if err != nil {
		return "", commandTab{}, err
	}
	if created.Workspace.ID == "" {
		return "", commandTab{}, fmt.Errorf("herdr workspace create response has no workspace id")
	}
	return created.Workspace.ID, created.commandTab(), nil
}

func (h Cmd) createTab(workspaceID, workingDir string) (commandTab, error) {
	created, err := h.call("tab", "create", "--workspace", workspaceID, "--cwd", workingDir, "--no-focus")
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

func (h Cmd) output(args ...string) ([]byte, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(h.Binary, args...)
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := h.Run(cmd); err != nil {
		return nil, fmt.Errorf("herdr %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return stdout.Bytes(), nil
}

func (h Cmd) call(args ...string) (result, error) {
	stdout, err := h.output(args...)
	if err != nil {
		return result{}, err
	}
	var response struct {
		Result *result `json:"result"`
		Error  *struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(stdout, &response); err != nil {
		return result{}, fmt.Errorf("decode herdr response: %w", err)
	}
	if response.Error != nil {
		return result{}, fmt.Errorf("herdr: %s: %s", response.Error.Code, response.Error.Message)
	}
	if response.Result == nil {
		return result{}, fmt.Errorf("herdr response has no result")
	}
	return *response.Result, nil
}
