package herdr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/eskelinenantti/zipln/internal/project"
	"github.com/eskelinenantti/zipln/internal/shell/runner"
)

type Cmd struct {
	runner.Runner
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
	cmd := exec.Command("herdr", args...)
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

func (h Cmd) Open(project project.Project, command []string) error {
	list, err := h.call("workspace", "list")
	if err != nil {
		return fmt.Errorf("cannot list workspaces; start a Herdr session first: %w", err)
	}
	if list.Type != "workspace_list" || list.Workspaces == nil {
		return fmt.Errorf("invalid herdr workspace list response")
	}
	var id string
	for _, workspace := range list.Workspaces {
		if workspace.Label == project.Name {
			id = workspace.ID
			if id == "" {
				return fmt.Errorf("herdr workspace has no id")
			}
			break
		}
	}
	if id == "" && len(list.Workspaces) > 0 && project.WorkingDir != "" {
		panes, err := h.call("pane", "list")
		if err != nil {
			return err
		}
		if panes.Type != "pane_list" || panes.Panes == nil {
			return fmt.Errorf("invalid herdr pane list response")
		}
		for _, workspace := range list.Workspaces {
			for _, pane := range panes.Panes {
				if pane.WorkspaceID == workspace.ID && pane.Cwd != "" && filepath.Clean(pane.Cwd) == filepath.Clean(project.WorkingDir) {
					id = workspace.ID
					if id == "" {
						return fmt.Errorf("herdr workspace has no id")
					}
					break
				}
			}
			if id != "" {
				break
			}
		}
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
		// pane run submits one shell command atomically. Quote every argv element
		// so paths and editor flags cannot become shell syntax.
		quoted := make([]string, len(command))
		for i, arg := range command {
			quoted[i] = "'" + strings.ReplaceAll(arg, "'", "'\"'\"'") + "'"
		}
		// pane run is quiet on success, unlike the JSON-returning create APIs.
		if _, err = h.output("pane", "run", created.RootPane.ID, strings.Join(quoted, " ")); err != nil {
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
	cmd := exec.Command("herdr")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return h.Run(cmd)
}
