package spy

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

type RunFunc func(cmd *exec.Cmd) error

type Response struct {
	OnRun RunFunc
}

type SpyRunner struct {
	Calls     [][]string
	Responses []Response
	Stdin     bytes.Buffer
}

func (t *SpyRunner) Run(cmd *exec.Cmd) error {
	t.Calls = append(t.Calls, cmd.Args)
	// Reserve this command's response before reading stdin, which can lazily
	// start another command while the current command is running.
	var response Response
	if len(t.Responses) > 0 {
		response = t.Responses[0]
		t.Responses = t.Responses[1:]
	}
	var inputErr error
	if _, file := cmd.Stdin.(*os.File); cmd.Stdin != nil && !file {
		_, inputErr = io.Copy(&t.Stdin, cmd.Stdin)
	}
	if response.OnRun != nil {
		if err := response.OnRun(cmd); err != nil {
			return err
		}
	}
	return inputErr
}

func (t *SpyRunner) Stdout(cmd *exec.Cmd) (io.ReadCloser, error) {
	var output bytes.Buffer
	cmd.Stdout = &output
	if err := t.Run(cmd); err != nil {
		return nil, err
	}
	return io.NopCloser(&output), nil
}
