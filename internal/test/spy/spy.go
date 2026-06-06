package spy

import (
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
}

type FakeWriteCloser struct{}

func (f FakeWriteCloser) Close() error {
	return nil
}

func (f FakeWriteCloser) Write(p []byte) (n int, err error) {
	return 0, nil
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
	return FakeWriteCloser{}, t.Run(cmd)
}
