package cmd

import (
	"os"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/shell/tmux"
	"github.com/eskelinenantti/tmuxide/internal/test/spy"
	"github.com/google/go-cmp/cmp"
)

func TestQuit(t *testing.T) {
	t.Setenv("TMUX", "test")
	tmuxSpy := &spy.Tmux{}
	shellEnv := ShellEnv{
		Tmux: tmuxSpy,
	}

	err := Quit(shellEnv)
	if err != nil {
		t.Errorf("err=%v", err)
	}

	expectedCalls := []spy.Call{
		{Name: "kill-session", Args: tmux.Args{}},
	}

	if !cmp.Equal(tmuxSpy.Calls, expectedCalls) {
		t.Error(cmp.Diff(tmuxSpy.Calls, expectedCalls))
	}
}

func TestQuitWhenNotAttached(t *testing.T) {
	os.Unsetenv("TMUX")
	tmuxSpy := &spy.Tmux{}
	shellEnv := ShellEnv{
		Tmux: tmuxSpy,
	}

	err := Quit(shellEnv)
	if err != nil {
		t.Errorf("err=%v", err)
	}

	var expectedCalls []spy.Call

	if !cmp.Equal(tmuxSpy.Calls, expectedCalls) {
		t.Error(cmp.Diff(tmuxSpy.Calls, expectedCalls))
	}
}
