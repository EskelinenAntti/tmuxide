package herdr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/eskelinenantti/tmuxide/internal/project"
	"github.com/eskelinenantti/tmuxide/internal/shell/runner"
)

type Cmd struct {
	runner.Runner
	Binary string
}

type workspace struct {
	ID    string `json:"workspace_id"`
	Label string `json:"label"`
}

type result struct {
	Type       string      `json:"type"`
	Workspaces []workspace `json:"workspaces"`
	Workspace  workspace   `json:"workspace"`
	Panes      []struct {
		WorkspaceID string `json:"workspace_id"`
		Cwd         string `json:"cwd"`
	} `json:"panes"`
	Tab struct {
		ID string `json:"tab_id"`
	} `json:"tab"`
	RootPane struct {
		ID string `json:"pane_id"`
	} `json:"root_pane"`
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

func (h Cmd) existingWorkspace(project project.Project) (string, error) {
	list, err := h.call("workspace", "list")
	if err != nil {
		return "", fmt.Errorf("cannot list workspaces; start a Herdr session first: %w", err)
	}
	if list.Type != "workspace_list" || list.Workspaces == nil {
		return "", fmt.Errorf("invalid herdr workspace list response")
	}
	for _, workspace := range list.Workspaces {
		if workspace.Label == project.Name {
			if workspace.ID == "" {
				return "", fmt.Errorf("herdr workspace has no id")
			}
			return workspace.ID, nil
		}
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
	for _, workspace := range list.Workspaces {
		if !matchingWorkspaces[workspace.ID] {
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
		if created.RootPane.ID == "" || created.Tab.ID == "" {
			return fmt.Errorf("herdr create response has no root pane or tab id")
		}
		// pane run is quiet on success, unlike the JSON-returning create APIs.
		if _, err = h.output("pane", "run", created.RootPane.ID, shellCommand(command)); err != nil {
			return err
		}
		if _, err = h.output("tab", "focus", created.Tab.ID); err != nil {
			return err
		}
	}
	if _, err = h.output("workspace", "focus", id); err != nil {
		return err
	}
	if os.Getenv("HERDR_ENV") == "1" {
		return nil
	}
	cmd := exec.Command(h.Binary)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return h.Run(cmd)
}
