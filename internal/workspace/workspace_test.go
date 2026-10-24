package workspace

import (
	"errors"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/shell"
	"github.com/eskelinenantti/tmuxide/internal/shell/herdr"
	"github.com/eskelinenantti/tmuxide/internal/shell/tmux"
	"github.com/eskelinenantti/tmuxide/internal/test/mock"
	"github.com/eskelinenantti/tmuxide/internal/test/spy"
)

func TestSelect(t *testing.T) {
	for _, tt := range []struct {
		name, backend, env, binary string
		herdr                      bool
	}{
		{"default", "auto", "", "", false},
		{"inside herdr", "auto", "1", "", true},
		{"explicit tmux", "tmux", "1", "", false},
		{"explicit herdr", "herdr", "", "", true},
		{"host binary", "herdr", "1", "/opt/herdr", true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("HERDR_ENV", tt.env)
			t.Setenv("HERDR_BIN_PATH", tt.binary)
			t.Setenv("TMUX", "outer-tmux")
			host, fullscreen, err := Select(tt.backend, mock.Path{}, &spy.SpyRunner{})
			if err != nil {
				t.Fatal(err)
			}
			if fullscreen != tt.herdr {
				t.Fatal("incorrect picker mode")
			}
			if tt.herdr {
				cmd, ok := host.(herdr.Cmd)
				want := tt.binary
				if want == "" {
					want = "herdr"
				}
				if !ok || cmd.Binary != want {
					t.Fatalf("got %#v", host)
				}
			} else if _, ok := host.(tmux.Cmd); !ok {
				t.Fatalf("got %#v", host)
			}
		})
	}
}

func TestInvalidOrMissingBackend(t *testing.T) {
	t.Setenv("HERDR_BIN_PATH", "")
	for _, name := range []string{"tmux", "herdr"} {
		_, _, err := Select(name, mock.Path{Missing: []string{name}}, &spy.SpyRunner{})
		var missing shell.NotInstalledError
		if !errors.As(err, &missing) || missing.Cmd != name {
			t.Fatal(err)
		}
	}
	if _, _, err := Select("invalid", mock.Path{}, &spy.SpyRunner{}); err == nil {
		t.Fatal("expected invalid backend error")
	}
}
