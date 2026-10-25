package spy

import (
	"io"
	"os/exec"
)

type RunFunc func(cmd *exec.Cmd) error

type SpyRunner struct {
	Calls     [][]string
	Responses []RunFunc
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
	if response == nil {
		return nil
	}
	return response(cmd)
}

func (t *SpyRunner) Start(cmd *exec.Cmd) (io.WriteCloser, error) {
	return FakeWriteCloser{}, t.Run(cmd)
}
