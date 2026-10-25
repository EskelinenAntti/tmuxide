package herdr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

func (h Cmd) output(args ...string) ([]byte, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(h.Binary, args...)
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := h.Run(cmd); err != nil {
		return nil, fmt.Errorf("herdr %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return stdout.Bytes(), nil
}

func call[T any](h Cmd, args ...string) (T, error) {
	var zero T
	stdout, err := h.output(args...)
	if err != nil {
		return zero, err
	}
	var response struct {
		Result *T `json:"result"`
		Error  *struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(stdout, &response); err != nil {
		return zero, fmt.Errorf("decode herdr response: %w", err)
	}
	if response.Error != nil {
		return zero, fmt.Errorf("herdr: %s: %s", response.Error.Code, response.Error.Message)
	}
	if response.Result == nil {
		return zero, fmt.Errorf("herdr response has no result")
	}
	return *response.Result, nil
}
