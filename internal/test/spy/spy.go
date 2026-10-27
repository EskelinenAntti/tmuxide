package spy

import (
	"bytes"
	"os/exec"

	"github.com/eskelinenantti/tmuxide/internal/shell/runner"
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

type FakeWriteCloser struct {
	*bytes.Buffer
}

func (f FakeWriteCloser) Close() error {
	return nil
}

func (t *SpyRunner) Run(cmd *exec.Cmd) error {
	t.Calls = append(t.Calls, cmd.Args)

	if len(t.Responses) == 0 {
		return nil
	}

	response := t.Responses[0]
	t.Responses = t.Responses[1:]
	if response.OnRun == nil {
		return nil
	}
	return response.OnRun(cmd)
}

func (t *SpyRunner) Start(cmd *exec.Cmd) (runner.WriteCloser, error) {
	return FakeWriteCloser{Buffer: &t.Stdin}, t.Run(cmd)
}
