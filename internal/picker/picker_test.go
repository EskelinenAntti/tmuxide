package picker

import (
	"io"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/shell/fd"
	"github.com/eskelinenantti/tmuxide/internal/shell/fzf"
	"github.com/eskelinenantti/tmuxide/internal/test/spy"
	"github.com/google/go-cmp/cmp"
)

func TestPrompt(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	for _, tt := range []struct {
		name      string
		selection string
		want      string
	}{
		{"selected path", "folder/file\n", filepath.Join(home, "folder/file")},
		{"no selection", "", ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			process := &spy.SpyRunner{Responses: []spy.Response{{OnRun: func(cmd *exec.Cmd) error {
				_, err := io.WriteString(cmd.Stdout, tt.selection)
				return err
			}}}}
			got, err := Prompt(fd.Cmd{Runner: process}, fzf.Cmd{Starter: process})
			if err != nil || got != tt.want {
				t.Fatalf("got %q, %v; want %q", got, err, tt.want)
			}
			if len(process.Calls) != 2 || process.Calls[0][0] != "fzf" || process.Calls[1][0] != "fd" {
				t.Fatalf("unexpected commands: %v", process.Calls)
			}
			wantArgs := []string{"fd", "--follow", "--hidden", "--exclude", "{.git,node_modules,Library}", ".", "--base-directory", home}
			if diff := cmp.Diff(wantArgs, process.Calls[1]); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestPromptKeepsSearchAndSelectionBaseDirectoryConsistent(t *testing.T) {
	home := t.TempDir()
	otherHome := t.TempDir()
	t.Setenv("HOME", home)
	process := &spy.SpyRunner{Responses: []spy.Response{{OnRun: func(cmd *exec.Cmd) error {
		t.Setenv("HOME", otherHome)
		_, err := io.WriteString(cmd.Stdout, "folder/file\n")
		return err
	}}}}
	got, err := Prompt(fd.Cmd{Runner: process}, fzf.Cmd{Starter: process})
	if want := filepath.Join(home, "folder/file"); err != nil || got != want {
		t.Fatalf("got %q, %v; want %q", got, err, want)
	}
	want := [][]string{
		{"fzf", "--reverse", "--height", "70%", "--tmux", "70%"},
		{"fd", "--follow", "--hidden", "--exclude", "{.git,node_modules,Library}", ".", "--base-directory", home},
	}
	if diff := cmp.Diff(want, process.Calls); diff != "" {
		t.Fatal(diff)
	}
}
