package spy

import (
	"errors"
	"os/exec"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRunConsumesCallbacksInOrder(t *testing.T) {
	failure := errors.New("command failed")
	first := exec.Command("first", "argument")
	third := exec.Command("third")
	var callbacks []*exec.Cmd
	runner := &SpyRunner{Responses: []RunFunc{
		func(cmd *exec.Cmd) error {
			callbacks = append(callbacks, cmd)
			return failure
		},
		nil,
		func(cmd *exec.Cmd) error {
			callbacks = append(callbacks, cmd)
			return nil
		},
	}}
	for i, cmd := range []*exec.Cmd{first, exec.Command("second"), third, exec.Command("fourth")} {
		var wantErr error
		if i == 0 {
			wantErr = failure
		}
		if err := runner.Run(cmd); err != wantErr {
			t.Fatalf("call %d: got error %v, want %v", i, err, wantErr)
		}
		if len(runner.Calls) != i+1 {
			t.Fatalf("call %d: got %d recorded calls", i, len(runner.Calls))
		}
	}
	if len(callbacks) != 2 || callbacks[0] != first || callbacks[1] != third {
		t.Fatal("callbacks did not receive their commands in order")
	}
	if len(runner.Responses) != 0 {
		t.Fatal("responses were not consumed")
	}
	want := [][]string{{"first", "argument"}, {"second"}, {"third"}, {"fourth"}}
	if diff := cmp.Diff(want, runner.Calls); diff != "" {
		t.Fatal(diff)
	}
}

func TestStartUsesQueuedCallback(t *testing.T) {
	failure := errors.New("start failed")
	for _, wantErr := range []error{nil, failure} {
		cmd := exec.Command("command")
		runner := &SpyRunner{Responses: []RunFunc{func(got *exec.Cmd) error {
			if got != cmd {
				t.Fatal("callback received a different command")
			}
			return wantErr
		}}}
		writer, err := runner.Start(cmd)
		if err != wantErr {
			t.Fatalf("got error %v, want %v", err, wantErr)
		}
		if writer == nil {
			t.Fatal("missing fake input writer")
		}
		if diff := cmp.Diff([][]string{{"command"}}, runner.Calls); diff != "" {
			t.Fatal(diff)
		}
		if len(runner.Responses) != 0 {
			t.Fatal("response was not consumed")
		}
	}
}
