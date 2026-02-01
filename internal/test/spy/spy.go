package spy

import (
	"os/exec"
	"reflect"
	"slices"
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
