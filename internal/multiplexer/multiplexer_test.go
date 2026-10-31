package multiplexer

import (
	"errors"
	"testing"

	"github.com/eskelinenantti/zipln/internal/shell"
	"github.com/eskelinenantti/zipln/internal/shell/herdr"
	"github.com/eskelinenantti/zipln/internal/shell/tmux"
	"github.com/eskelinenantti/zipln/internal/test/mock"
	"github.com/eskelinenantti/zipln/internal/test/spy"
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
				if _, ok := mux.(herdr.Cmd); !ok {
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
