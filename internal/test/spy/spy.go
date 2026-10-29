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

func (t *SpyRunner) Stdout(cmd *exec.Cmd) (io.ReadCloser, error) {
	var output bytes.Buffer
	cmd.Stdout = &output
	if err := t.Run(cmd); err != nil {
		return nil, err
	}
	return io.NopCloser(&output), nil
}

func (t *SpyRunner) Pipe(stdin io.Reader, cmd *exec.Cmd, stdout io.Writer) error {
	cmd.Stdin, cmd.Stdout = stdin, stdout
	err := t.Run(cmd)
	if stdin != nil {
		if _, err := io.Copy(&t.Stdin, stdin); err != nil {
			return &runner.InputError{Err: err}
		}
	}
	return err
}
