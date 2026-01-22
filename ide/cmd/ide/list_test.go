package cmd

import (
	"os"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/shell/tmux"
	"github.com/eskelinenantti/tmuxide/internal/test/spy"
	"github.com/google/go-cmp/cmp"
)

func TestList(t *testing.T) {
	os.Unsetenv("TMUX")
	tmuxSpy := &spy.Tmux{}
	shellEnv := ShellEnv{
		Tmux: tmuxSpy,
	}
	err := List(shellEnv)
	if err != nil {
		t.Errorf("err=%v", err)
	}

	expectedCalls := []spy.Call{
		{Name: "choose-session", Args: tmux.Args{}},
		{Name: "attach", Args: tmux.Args{}},
	}

	if !cmp.Equal(tmuxSpy.Calls, expectedCalls) {
		t.Error(cmp.Diff(tmuxSpy.Calls, expectedCalls))
	}
}

func TestListWhenAttached(t *testing.T) {
	t.Setenv("TMUX", "test")
	tmuxSpy := &spy.Tmux{}
	shellEnv := ShellEnv{
		Tmux: tmuxSpy,
	}
	err := List(shellEnv)
	if err != nil {
		t.Errorf("err=%v", err)
	}

	expectedCalls := []spy.Call{
		{Name: "choose-session", Args: tmux.Args{}},
	}

	if !cmp.Equal(tmuxSpy.Calls, expectedCalls) {
		t.Error(cmp.Diff(tmuxSpy.Calls, expectedCalls))
	}
}
