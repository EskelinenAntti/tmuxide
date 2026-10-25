package fzf

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/shell/runner"
	"github.com/eskelinenantti/tmuxide/internal/test/spy"
	"github.com/google/go-cmp/cmp"
)

func TestIsCancelled(t *testing.T) {
	for _, code := range []int{0, 1, 2, 130} {
		t.Run(strconv.Itoa(code), func(t *testing.T) {
			cmd := exec.Command(os.Args[0], "-test.run=^TestFzfExitProcess$")
			cmd.Env = append(os.Environ(), "TMUXIDE_FZF_EXIT_CODE="+strconv.Itoa(code))
			err := cmd.Run()
			if code != 0 {
				var exitErr *exec.ExitError
				if !errors.As(err, &exitErr) || exitErr.ExitCode() != code {
					t.Fatalf("got %v, want exit code %d", err, code)
				}
			} else if err != nil {
				t.Fatal(err)
			}
			want := code == 130
			if got := IsCancelled(err); got != want {
				t.Fatalf("IsCancelled(%v) = %v, want %v", err, got, want)
			}
			if err != nil && IsCancelled(fmt.Errorf("command failed: %w", err)) != want {
				t.Fatal("wrapping changed cancellation detection")
			}
		})
	}
	if IsCancelled(errors.New("exit status 130")) {
		t.Fatal("an ordinary error is not cancellation")
	}
}

func TestFzfExitProcess(t *testing.T) {
	if value := os.Getenv("TMUXIDE_FZF_EXIT_CODE"); value != "" {
		code, err := strconv.Atoi(value)
		if err != nil {
			t.Fatal(err)
		}
		os.Exit(code)
	}
}

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
