package spy

import (
	"bytes"
	"io"
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

type FakeProcess struct {
	*bytes.Buffer
	done chan struct{}
}

func (f FakeProcess) Close() error {
	close(f.done)
	return nil
}

func (f FakeProcess) Wait() error {
	<-f.done
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

func (t *SpyRunner) Start(cmd *exec.Cmd) (runner.Process, error) {
	return FakeProcess{Buffer: &t.Stdin, done: make(chan struct{})}, t.Run(cmd)
}

func (t *SpyRunner) Read(cmd *exec.Cmd) (io.ReadCloser, error) {
	var output bytes.Buffer
	cmd.Stdout = &output
	if err := t.Run(cmd); err != nil {
		return nil, err
	}
	return io.NopCloser(&output), nil
}
