package fd

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/test/spy"
	"github.com/google/go-cmp/cmp"
)

func TestStreamCandidatesErrors(t *testing.T) {
	exitError := func(code string) error {
		t.Helper()
		err := exec.Command("sh", "-c", "exit "+code).Run()
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			t.Fatalf("got %v, want exit error", err)
		}
		return err
	}
	pipeClosed := exitError("1")
	failure := exitError("2")
	for _, tt := range []struct {
		name string
		err  error
		want error
	}{
		{name: "success"},
		{name: "closed candidate pipe", err: pipeClosed},
		{name: "wrapped closed candidate pipe", err: fmt.Errorf("command: %w", pipeClosed)},
		{name: "other exit failure", err: failure, want: failure},
		{name: "wrapped exit failure", err: fmt.Errorf("command: %w", failure), want: failure},
		{name: "non-exit error", err: errors.New("cannot start command")},
	} {
		t.Run(tt.name, func(t *testing.T) {
			process := &spy.SpyRunner{Responses: []spy.Response{{OnRun: func(*exec.Cmd) error {
				return tt.err
			}}}}
			err := (Cmd{Runner: process}).StreamCandidates("/project", io.Discard)
			if !errors.Is(err, tt.want) {
				t.Fatalf("got %v, want %v", err, tt.want)
			}
			if tt.want != nil && err != tt.err {
				t.Fatalf("got %v, want original error %v", err, tt.err)
			}
		})
	}
}

func TestStreamCandidatesUsesProvidedBaseDirectory(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	baseDir := t.TempDir()
	var outputWasSet bool
	process := &spy.SpyRunner{Responses: []spy.Response{{OnRun: func(cmd *exec.Cmd) error {
		outputWasSet = cmd.Stdout == io.Discard
		return nil
	}}}}
	if err := (Cmd{Runner: process}).StreamCandidates(baseDir, io.Discard); err != nil {
		t.Fatal(err)
	}
	want := [][]string{{"fd", "--follow", "--hidden", "--exclude", "{.git,node_modules,Library}", ".", "--base-directory", baseDir}}
	if diff := cmp.Diff(want, process.Calls); diff != "" {
		t.Fatal(diff)
	}
	if !outputWasSet {
		t.Fatal("fd output was not connected to the provided writer")
	}
}
