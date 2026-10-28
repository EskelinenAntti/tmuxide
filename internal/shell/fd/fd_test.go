package fd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/history"
	"github.com/eskelinenantti/tmuxide/internal/test/mock"
	"github.com/eskelinenantti/tmuxide/internal/test/spy"
	"github.com/google/go-cmp/cmp"
)

func TestFdWritesRankedPathsThenOneScan(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_STATE_HOME", t.TempDir())
	dir := filepath.Join(home, "project")
	if err := os.Mkdir(dir, 0755); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(file, nil, 0644); err != nil {
		t.Fatal(err)
	}
	history.Record(dir)
	history.Record(file)
	history.Record(file)
	history.Record(filepath.Join(home, "missing"))

	runner := &spy.SpyRunner{Responses: []spy.Response{
		{OnRun: mock.WriteToStdout("./project/\nproject/file.txt\nproject/child.txt\nother.txt\n")},
	}}
	var output bytes.Buffer
	if err := (Cmd{Runner: runner}).Fd(&output); err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff("project/file.txt\nproject/\nproject/child.txt\nother.txt\n", output.String()); diff != "" {
		t.Fatal(diff)
	}
	wantCalls := [][]string{{"fd", "--follow", "--hidden", "--exclude", "{.git,node_modules,Library}", ".", "--base-directory", home, "--"}}
	if diff := cmp.Diff(wantCalls, runner.Calls); diff != "" {
		t.Fatal(diff)
	}
}

func TestWriteSkipsDuplicates(t *testing.T) {
	var output bytes.Buffer
	seen := make(map[string]bool)
	for _, path := range []string{"project/", "project/file.txt", "./project", "./project/file.txt", "other file.txt"} {
		if err := write(&output, path, seen); err != nil {
			t.Fatal(err)
		}
	}
	want := "project/\nproject/file.txt\nother file.txt\n"
	if diff := cmp.Diff(want, output.String()); diff != "" {
		t.Fatal(diff)
	}
}
