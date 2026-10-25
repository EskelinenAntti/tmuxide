package shell

import (
	"errors"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/test/mock"
	"github.com/eskelinenantti/tmuxide/internal/test/spy"
)

func TestInit(t *testing.T) {
	runner := &spy.SpyRunner{}
	shell, err := Init(mock.Path{}, runner)
	if err != nil {
		t.Fatal(err)
	}
	if shell.Fd.Runner != runner || shell.Fzf.Starter != runner || shell.Git.Runner != runner {
		t.Fatalf("commands do not share the supplied runner: %+v", shell)
	}
	if len(runner.Calls) != 0 {
		t.Fatalf("initialization ran commands: %v", runner.Calls)
	}
}

func TestInitMissingDependency(t *testing.T) {
	for _, tt := range []struct {
		name    string
		missing []string
		want    string
	}{
		{"fd", []string{"fd"}, "fd"},
		{"fzf", []string{"fzf"}, "fzf"},
		{"git", []string{"git"}, "git"},
		{"first missing dependency", []string{"git", "fzf", "fd"}, "fd"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			runner := &spy.SpyRunner{}
			shell, err := Init(mock.Path{Missing: tt.missing}, runner)
			var missing NotInstalledError
			if !errors.As(err, &missing) || missing.Cmd != tt.want {
				t.Fatalf("got %v, want missing dependency %q", err, tt.want)
			}
			if !errors.Is(err, ErrCommandNotInstalled) {
				t.Fatalf("missing dependency error lost its cause: %v", err)
			}
			if shell != (Shell{}) {
				t.Fatalf("got partially initialized shell: %+v", shell)
			}
			if len(runner.Calls) != 0 {
				t.Fatalf("initialization ran commands: %v", runner.Calls)
			}
		})
	}
}
