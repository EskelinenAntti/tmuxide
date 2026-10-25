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
	"github.com/eskelinenantti/tmuxide/internal/shell/runner"
	"github.com/eskelinenantti/tmuxide/internal/shell/tmux"
	"github.com/eskelinenantti/tmuxide/internal/test/mock"
	"github.com/eskelinenantti/tmuxide/internal/test/spy"
)

func TestSelect(t *testing.T) {
	t.Setenv("HERDR_BIN_PATH", "/path/to/herdr")
	for _, tt := range []struct {
		name, mux, env, tmux string
		herdr                bool
	}{
		{"default", "auto", "", "", false},
		{"inside tmux", "auto", "", "tmux", false},
		{"inactive herdr", "auto", "0", "", false},
		{"inactive herdr inside tmux", "auto", "0", "tmux", false},
		{"noncanonical herdr flag", "auto", "true", "", false},
		{"inside herdr", "auto", "1", "", true},
		{"herdr inside tmux", "auto", "1", "outer-tmux", true},
		{"explicit tmux", "tmux", "1", "", false},
		{"explicit tmux inside nested hosts", "tmux", "1", "outer-tmux", false},
		{"explicit herdr", "herdr", "", "tmux", true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("HERDR_ENV", tt.env)
			t.Setenv("TMUX", tt.tmux)
			runOnly := struct{ runner.Runner }{&spy.SpyRunner{}}
			mux, err := Select(tt.mux, mock.Path{}, runOnly)
			if err != nil {
				t.Fatal(err)
			}
			if tt.herdr {
				cmd, ok := mux.(herdr.Cmd)
				if !ok || cmd.Binary != "/path/to/herdr" {
					t.Fatalf("got %#v", mux)
				}
			} else if _, ok := mux.(tmux.Cmd); !ok {
				t.Fatalf("got %#v", mux)
			}
		})
	}
}

func TestInvalidOrMissingMultiplexer(t *testing.T) {
	t.Setenv("HERDR_BIN_PATH", "/path/to/herdr")
	for _, name := range []string{"tmux", "herdr"} {
		_, err := Select(name, mock.Path{Missing: []string{name, "/path/to/herdr"}}, &spy.SpyRunner{})
		var missing shell.NotInstalledError
		if !errors.As(err, &missing) || missing.Cmd != map[string]string{"tmux": "tmux", "herdr": "/path/to/herdr"}[name] {
			t.Fatal(err)
		}
	}
	if _, err := Select("invalid", mock.Path{}, &spy.SpyRunner{}); err == nil {
		t.Fatal("expected invalid multiplexer error")
	}
}

func TestSelectHerdrUsesConfiguredBinaryWithoutPATH(t *testing.T) {
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
			t.Fatalf("got %#v, want configured binary %q", mux, binary)
		}
	}
}

func TestSelectHerdrReportsMissingConfiguredBinary(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "missing-herdr")
	t.Setenv("HERDR_BIN_PATH", binary)
	_, err := Select("herdr", path.Path{}, &spy.SpyRunner{})
	var missing shell.NotInstalledError
	if !errors.As(err, &missing) || missing.Cmd != binary {
		t.Fatalf("got %v, want missing configured binary %q", err, binary)
	}
}

func TestSelectHerdrRequiresConfiguredBinary(t *testing.T) {
	t.Setenv("HERDR_BIN_PATH", "")
	_, err := Select("herdr", mock.Path{}, &spy.SpyRunner{})
	if err == nil || !strings.Contains(err.Error(), "HERDR_BIN_PATH is required") {
		t.Fatalf("got %v, want missing HERDR_BIN_PATH error", err)
	}
}
