package picker

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/eskelinenantti/herdr-zipln/internal/history"
	"github.com/eskelinenantti/herdr-zipln/internal/shell/fd"
	"github.com/eskelinenantti/herdr-zipln/internal/shell/fzf"
	"github.com/eskelinenantti/herdr-zipln/internal/shell/runner"
	"github.com/eskelinenantti/herdr-zipln/internal/test/mock"
	"github.com/eskelinenantti/herdr-zipln/internal/test/spy"
	"github.com/google/go-cmp/cmp"
)

func TestPromptWritesRankedPathsThenOneScan(t *testing.T) {
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

	r := &spy.SpyRunner{Responses: []spy.Response{
		{OnRun: mock.WriteToStdout("./project/\nproject/file.txt\nproject/child.txt\nother.txt\n")},
		{OnRun: mock.WriteToStdout("project/\n")},
	}}
	selection, err := Prompt(fd.Cmd{Runner: r}, fzf.Cmd{Runner: r})
	if err != nil {
		t.Fatal(err)
	}
	if selection != dir {
		t.Fatalf("Prompt selection = %q, want %q", selection, dir)
	}
	if diff := cmp.Diff("project/file.txt\nproject/\nproject/child.txt\nother.txt\n", r.Stdin.String()); diff != "" {
		t.Fatal(diff)
	}
	wantCalls := [][]string{
		{"fd", "--follow", "--hidden", "--exclude", "{.git,node_modules,Library}", ".", "--base-directory", home},
		{"fzf", "--reverse", "--height", "70%", "--tmux", "70%", "--no-sort"},
	}
	if diff := cmp.Diff(wantCalls, r.Calls); diff != "" {
		t.Fatal(diff)
	}
}

func TestPrompt(t *testing.T) {
	tests := []struct {
		name      string
		fd        string
		fzf       string
		selection string
		wantError bool
	}{
		{
			name:      "closes input before waiting for selection",
			fd:        "printf 'project/\n'",
			fzf:       "while IFS= read -r line; do :; done; printf 'project/\n'",
			selection: "project",
		},
		{
			name:      "absolute selection",
			fd:        "printf 'project/\n'",
			fzf:       "while IFS= read -r line; do :; done; printf '/outside/project/\n'",
			selection: "/outside/project/",
		},
		{
			name:      "selection with trailing space",
			fd:        "printf 'project /\n'",
			fzf:       "while IFS= read -r line; do :; done; printf 'project \n'",
			selection: "project ",
		},
		{
			name: "empty selection",
			fd:   "printf 'project/\n'",
			fzf:  "while IFS= read -r line; do :; done",
		},
		{
			name: "user cancellation",
			fd:   "printf 'project/\n'",
			fzf:  "exit 130",
		},
		{
			name: "user cancellation discards output",
			fd:   "printf 'project/\n'",
			fzf:  "printf 'project/\n'; exit 130",
		},
		{
			name:      "selection before scan finishes",
			fd:        "i=0; while [ \"$i\" -lt 100000 ]; do printf 'file%s\n' \"$i\"; i=$((i+1)); done",
			fzf:       "printf 'project/\n'",
			selection: "project",
		},
		{
			name:      "scan failure without selection",
			fd:        "exit 7",
			fzf:       "while IFS= read -r line; do :; done",
			wantError: true,
		},
		{
			name:      "scan failure takes precedence over selection",
			fd:        "exit 7",
			fzf:       "while IFS= read -r line; do :; done; printf 'project/\n'",
			wantError: true,
		},
		{
			name:      "scan exit 130 is not fzf cancellation",
			fd:        "exit 130",
			fzf:       "while IFS= read -r line; do :; done",
			wantError: true,
		},
		{
			name:      "fzf failure",
			fd:        "printf 'project/\n'",
			fzf:       "exit 2",
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			home := t.TempDir()
			t.Setenv("HOME", home)
			t.Setenv("XDG_STATE_HOME", t.TempDir())
			bin := t.TempDir()
			t.Setenv("PATH", bin)
			for name, script := range map[string]string{"fd": tt.fd, "fzf": tt.fzf} {
				if err := os.WriteFile(filepath.Join(bin, name), []byte("#!/bin/sh\n"+script+"\n"), 0755); err != nil {
					t.Fatal(err)
				}
			}
			r := runner.CmdRunner{}
			selection, err := Prompt(fd.Cmd{Runner: r}, fzf.Cmd{Runner: r})
			if (err != nil) != tt.wantError {
				t.Fatalf("Prompt error = %v, want error = %v", err, tt.wantError)
			}
			want := ""
			if filepath.IsAbs(tt.selection) {
				want = tt.selection
			} else if tt.selection != "" {
				want = filepath.Join(home, tt.selection)
			}
			if selection != want {
				t.Fatalf("Prompt selection = %q, want %q", selection, want)
			}
		})
	}
}
