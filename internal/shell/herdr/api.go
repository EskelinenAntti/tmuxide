package herdr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

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
