package fd

import (
	"io"
	"os/exec"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/test/spy"
	"github.com/google/go-cmp/cmp"
)

func TestFdUsesProvidedBaseDirectory(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	baseDir := t.TempDir()
	var outputWasSet bool
	process := &spy.SpyRunner{Responses: []spy.Response{{OnRun: func(cmd *exec.Cmd) error {
		outputWasSet = cmd.Stdout == io.Discard
		return nil
	}}}}
	if err := (Cmd{Runner: process}).Fd(baseDir, io.Discard); err != nil {
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
