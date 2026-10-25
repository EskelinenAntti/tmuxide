package multiplexer

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/shell"
	"github.com/eskelinenantti/tmuxide/internal/shell/herdr"
	"github.com/eskelinenantti/tmuxide/internal/shell/path"
	"github.com/eskelinenantti/tmuxide/internal/shell/tmux"
	"github.com/eskelinenantti/tmuxide/internal/test/mock"
	"github.com/eskelinenantti/tmuxide/internal/test/spy"
)

func TestSelect(t *testing.T) {
	for _, tt := range []struct {
		name, mux, env, tmux string
		herdr                bool
	}{
		{"default", "auto", "", "", false},
		{"inside tmux", "auto", "", "tmux", false},
		{"inside herdr", "auto", "1", "", true},
		{"herdr inside tmux", "auto", "1", "outer-tmux", true},
		{"explicit tmux", "tmux", "1", "", false},
		{"explicit herdr", "herdr", "", "tmux", true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("HERDR_ENV", tt.env)
			t.Setenv("TMUX", tt.tmux)
			mux, err := Select(tt.mux, mock.Path{}, &spy.SpyRunner{})
			if err != nil {
				t.Fatal(err)
			}
			if tt.herdr {
				_, ok := mux.(herdr.Cmd)
				if !ok {
					t.Fatalf("got %#v", mux)
				}
			} else if _, ok := mux.(tmux.Cmd); !ok {
				t.Fatalf("got %#v", mux)
			}
		})
	}
}

func TestInvalidOrMissingMultiplexer(t *testing.T) {
	for _, name := range []string{"tmux", "herdr"} {
		_, err := Select(name, mock.Path{Missing: []string{name}}, &spy.SpyRunner{})
		var missing shell.NotInstalledError
		if !errors.As(err, &missing) || missing.Cmd != name {
			t.Fatal(err)
		}
	}
	if _, err := Select("invalid", mock.Path{}, &spy.SpyRunner{}); err == nil {
		t.Fatal("expected invalid multiplexer error")
	}
}

func TestSelectHerdrUsesPATH(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "herdr"), []byte("#!/bin/sh\nexit 0\n"), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir)
	t.Setenv("HERDR_ENV", "1")
	for _, name := range []string{"herdr", "auto"} {
		mux, err := Select(name, path.Path{}, &spy.SpyRunner{})
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := mux.(herdr.Cmd); !ok {
			t.Fatalf("got %#v", mux)
		}
	}
}

func TestSelectHerdrReportsMissingCommand(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	_, err := Select("herdr", path.Path{}, &spy.SpyRunner{})
	var missing shell.NotInstalledError
	if !errors.As(err, &missing) || missing.Cmd != "herdr" {
		t.Fatalf("got %v, want missing herdr command", err)
	}
}
