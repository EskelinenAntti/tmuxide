package spy

import (
	"os/exec"
	"reflect"
	"slices"

	"github.com/eskelinenantti/tmuxide/internal/shell"
)

type MockFunc func(cmd *exec.Cmd) error

type Mock struct {
	Args  []string
	OnRun MockFunc
}

type SpyRunner struct {
	Calls [][]string
	Mocks []Mock
}

type FakeWriteCloser struct{}

func (f FakeWriteCloser) Close() error {
	return nil
}

func (f FakeWriteCloser) Write(p []byte) (n int, err error) {
	return 0, nil
}

func (t *SpyRunner) Run(cmd *exec.Cmd) error {
	call := cmd.Args
	t.Calls = append(t.Calls, call)

	for i, mock := range t.Mocks {
		if reflect.DeepEqual(mock.Args, call) {
			t.Mocks = slices.Delete(t.Mocks, i, i+1)
			return mock.OnRun(cmd)
		}
	}

	return nil
}

func (t *SpyRunner) Start(cmd *exec.Cmd) (shell.WriteCloser, error) {
	return FakeWriteCloser{}, t.Run(cmd)
}
