package spy

import (
	"errors"
	"os/exec"
	"reflect"
	"slices"
)

type SpyRunner struct {
	Calls  [][]string
	Errors [][]string
}

func (t *SpyRunner) Run(cmd exec.Cmd) error {
	call := cmd.Args
	t.Calls = append(t.Calls, call)

	for i, error := range t.Errors {
		if reflect.DeepEqual(error, call) {
			t.Errors = slices.Delete(t.Errors, i, i+1)
			return errors.New("error")
		}
	}
	return nil
}
