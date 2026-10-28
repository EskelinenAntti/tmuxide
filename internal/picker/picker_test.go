package picker

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/shell/fd"
	"github.com/eskelinenantti/tmuxide/internal/shell/fzf"
	"github.com/eskelinenantti/tmuxide/internal/shell/runner"
)

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
			name: "user cancellation",
			fd:   "printf 'project/\n'",
			fzf:  "exit 130",
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
			if tt.selection != "" {
				want = filepath.Join(home, tt.selection)
			}
			if selection != want {
				t.Fatalf("Prompt selection = %q, want %q", selection, want)
			}
		})
	}
}
