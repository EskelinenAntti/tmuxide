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
	history.RecordQuery("late", filepath.Join(home, "project/late"))
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
		{"preferred match", "rebind", query, "pos(2)+track-current+unbind(result)\n"},
		{"normal ranking", "rebind", "unknown", "unbind(result)\n"},
		{"empty query", "rebind", "", "unbind(result)\n"},
		{"wait for preferred match", "rebind", "late", ""},
		{"manual navigation", "down", query, "unbind(result)\n"},
		{"already restored", "position", query, "unbind(result)\n"},
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

func TestFocusCommandUsesLivePaths(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_STATE_HOME", t.TempDir())
	t.Setenv("FZF_ACTION", "rebind")
	history.RecordQuery("a", filepath.Join(home, "a"))
	matches := filepath.Join(t.TempDir(), "matches")
	for _, tt := range []struct {
		name    string
		matches string
		want    string
	}{
		{"a renamed to b", "b\nc\n", ""},
		{"c renamed to a inherits preference", "b\na\n", "pos(2)+track-current+unbind(result)\n"},
		{"a deleted", "b\n", ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if err := os.WriteFile(matches, []byte(tt.matches), 0600); err != nil {
				t.Fatal(err)
			}
			cmd := newFocusCmd()
			var output bytes.Buffer
			cmd.SetOut(&output)
			cmd.SetArgs([]string{"--", "a", matches})
			if err := cmd.Execute(); err != nil {
				t.Fatal(err)
			}
			if got := output.String(); got != tt.want {
				t.Fatalf("focus action = %q, want %q", got, tt.want)
			}
		})
	}
}
