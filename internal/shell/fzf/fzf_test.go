package fzf

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/shell/runner"
	"github.com/eskelinenantti/tmuxide/internal/test/spy"
	"github.com/google/go-cmp/cmp"
)

func TestFzfNeedsOnlyStart(t *testing.T) {
	for _, tt := range []struct {
		name string
		err  error
	}{
		{"success", nil},
		{"failure", errors.New("start failed")},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			process := &spy.SpyRunner{Responses: []spy.Response{{OnRun: func(cmd *exec.Cmd) error {
				if cmd.Stdout != &output || cmd.Stderr != os.Stderr {
					t.Error("fzf must use the supplied output and terminal stderr")
				}
				return tt.err
			}}}}
			startOnly := struct{ runner.Starter }{process}
			writer, err := (Cmd{Starter: startOnly}).Fzf(&output)
			if !errors.Is(err, tt.err) {
				t.Fatalf("got %v, want %v", err, tt.err)
			}
			if tt.err == nil && writer == nil {
				t.Fatal("successful start returned no writer")
			}
			want := [][]string{{"fzf", "--reverse", "--height", "70%", "--tmux", "70%"}}
			if diff := cmp.Diff(want, process.Calls); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
