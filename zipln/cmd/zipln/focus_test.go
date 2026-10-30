package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/eskelinenantti/herdr-zipln/internal/history"
)

func TestFocusCommand(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_STATE_HOME", t.TempDir())
	query := "-proj 'quoted'"
	history.RecordQuery(query, filepath.Join(home, "project/child"))
	matches := filepath.Join(t.TempDir(), "matches")
	if err := os.WriteFile(matches, []byte("project/\nproject/child/\n"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		name   string
		action string
		query  string
		want   string
	}{
		{"preferred match", "rebind", query, "pos(2)\n"},
		{"normal ranking", "rebind", "unknown", "pos(1)\n"},
		{"manual navigation", "down", query, ""},
		{"already restored", "position", query, ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("FZF_ACTION", tt.action)
			cmd := newFocusCmd()
			var output bytes.Buffer
			cmd.SetOut(&output)
			cmd.SetArgs([]string{"--", tt.query, matches})
			if err := cmd.Execute(); err != nil {
				t.Fatal(err)
			}
			if got := output.String(); got != tt.want {
				t.Fatalf("focus action = %q, want %q", got, tt.want)
			}
		})
	}
}
