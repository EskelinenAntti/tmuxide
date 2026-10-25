package multiplexer

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/shell"
	"github.com/eskelinenantti/tmuxide/internal/shell/herdr"
	"github.com/eskelinenantti/tmuxide/internal/shell/path"
	"github.com/eskelinenantti/tmuxide/internal/shell/tmux"
	"github.com/eskelinenantti/tmuxide/internal/test/mock"
	"github.com/eskelinenantti/tmuxide/internal/test/spy"
)

func TestSelect(t *testing.T) {
	t.Setenv("HERDR_BIN_PATH", "/injected/herdr")
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
				cmd, ok := mux.(herdr.Cmd)
				if !ok || cmd.Binary != "/injected/herdr" {
					t.Fatalf("got %#v", mux)
				}
			} else if _, ok := mux.(tmux.Cmd); !ok {
				t.Fatalf("got %#v", mux)
			}
		})
	}
}

func TestInvalidOrMissingMultiplexer(t *testing.T) {
	t.Setenv("HERDR_BIN_PATH", "/injected/herdr")
	for _, name := range []string{"tmux", "herdr"} {
		_, err := Select(name, mock.Path{Missing: []string{name, "/injected/herdr"}}, &spy.SpyRunner{})
		var missing shell.NotInstalledError
		if !errors.As(err, &missing) || missing.Cmd != map[string]string{"tmux": "tmux", "herdr": "/injected/herdr"}[name] {
			t.Fatal(err)
		}
	}
	if _, err := Select("invalid", mock.Path{}, &spy.SpyRunner{}); err == nil {
		t.Fatal("expected invalid multiplexer error")
	}
}

func TestSelectHerdrUsesInjectedBinaryWithoutPATH(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "herdr with spaces")
	if err := os.WriteFile(binary, []byte("#!/bin/sh\nexit 0\n"), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", "")
	t.Setenv("HERDR_BIN_PATH", binary)
	t.Setenv("HERDR_ENV", "1")
	for _, name := range []string{"herdr", "auto"} {
		mux, err := Select(name, path.Path{}, &spy.SpyRunner{})
		if err != nil {
			t.Fatal(err)
		}
		cmd, ok := mux.(herdr.Cmd)
		if !ok || cmd.Binary != binary {
			t.Fatalf("got %#v, want injected binary %q", mux, binary)
		}
	}
}

func TestSelectHerdrReportsMissingInjectedBinary(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "missing-herdr")
	t.Setenv("HERDR_BIN_PATH", binary)
	_, err := Select("herdr", path.Path{}, &spy.SpyRunner{})
	var missing shell.NotInstalledError
	if !errors.As(err, &missing) || missing.Cmd != binary {
		t.Fatalf("got %v, want missing injected binary %q", err, binary)
	}
}

func TestSelectHerdrRequiresInjectedBinary(t *testing.T) {
	t.Setenv("HERDR_BIN_PATH", "")
	_, err := Select("herdr", mock.Path{}, &spy.SpyRunner{})
	if err == nil || !strings.Contains(err.Error(), "HERDR_BIN_PATH is required") {
		t.Fatalf("got %v, want missing HERDR_BIN_PATH error", err)
	}
}
