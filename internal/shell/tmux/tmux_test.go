package tmux

import (
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/project"
	"github.com/eskelinenantti/tmuxide/internal/test/spy"
	"github.com/google/go-cmp/cmp"
)

func TestEachEditorRequestCreatesUnnamedWindow(t *testing.T) {
	t.Setenv("TMUX", "test")
	runner := &spy.SpyRunner{}
	mux := Cmd{Runner: runner}
	p := project.Project{Name: "project", WorkingDir: "/project"}
	var want [][]string
	for range 2 {
		if err := mux.Open(p, []string{"nvim", "file"}); err != nil {
			t.Fatal(err)
		}
		want = append(want,
			[]string{"tmux", "has-session", "-t", "project:"},
			[]string{"tmux", "new-window", "-t", "project:", "-c", "/project", "nvim", "file"},
			[]string{"tmux", "switch-client", "-t", "project:"},
		)
	}
	if diff := cmp.Diff(want, runner.Calls); diff != "" {
		t.Fatal(diff)
	}
}
