package picker

import (
	"io"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/shell/fd"
	"github.com/eskelinenantti/tmuxide/internal/shell/fzf"
	"github.com/eskelinenantti/tmuxide/internal/test/spy"
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
		})
	}
}
