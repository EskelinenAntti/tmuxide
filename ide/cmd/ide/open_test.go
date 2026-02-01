package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/project"
	"github.com/eskelinenantti/tmuxide/internal/test/mock"
	"github.com/eskelinenantti/tmuxide/internal/test/spy"
	"github.com/google/go-cmp/cmp"
)

/*
	func TestOpen(t *testing.T) {
		os.Unsetenv("TMUX")
		spyRunner := &spy.SpyRunner{}
		shellEnv := ShellEnv{
			CmdRunner: spyRunner,
			Path:      mock.Path{},
		}
		err := Open([]string{}, shellEnv)
		if err != nil {
			t.Errorf("err=%v", err)
		}

		expectedCalls := []spy.Call{
			{"choose-session"},
			{Name: "attach", Args: tmux.Args{}},
		}

		if !cmp.Equal(expectedCalls, tmuxSpy.Calls) {
			t.Error(cmp.Diff(expectedCalls, tmuxSpy.Calls))
		}
	}

	func TestOpenWhenAttached(t *testing.T) {
		t.Setenv("TMUX", "test")
		tmuxSpy := &spy.Tmux{}
		shellEnv := ShellEnv{
			TmuxRunner: tmuxSpy,
			Path:       mock.Path{},
		}
		err := Open([]string{}, shellEnv)
		if err != nil {
			t.Errorf("err=%v", err)
		}

		expectedCalls := []spy.Call{
			{Name: "choose-session", Args: tmux.Args{}},
		}

		if !cmp.Equal(expectedCalls, tmuxSpy.Calls) {
			t.Error(cmp.Diff(expectedCalls, tmuxSpy.Calls))
		}
	}

	func TestOpenWhenNoSessionsFound(t *testing.T) {
		t.Setenv("TMUX", "test")
		tmuxSpy := &spy.Tmux{
			Errors: []string{"choose-session"},
		}
		shellEnv := ShellEnv{
			TmuxRunner: tmuxSpy,
			Path:       mock.Path{},
		}

		err := Open([]string{}, shellEnv)
		if !errors.Is(err, ide.ErrNoSessionsFound) {
			t.Errorf("got=%v, want=%v", err, ide.ErrNoSessionsFound)
		}

		expectedCalls := []spy.Call{
			{Name: "choose-session", Args: tmux.Args{}},
		}

		if !cmp.Equal(expectedCalls, tmuxSpy.Calls) {
			t.Error(cmp.Diff(expectedCalls, tmuxSpy.Calls))
		}
	}
*/

func TestOpenDirInsideRepository(t *testing.T) {
	os.Unsetenv("TMUX")

	repository := t.TempDir()
	dir := filepath.Join(repository, "path/to/dir/in/repository")

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		t.Errorf("err=%v", err)
	}

	session := project.Name(dir)

	spyRunner := &spy.SpyRunner{
		Errors: [][]string{
			{"tmux", "has-session", "-t", dir + ":"},
			{"tmux", "has-session", "-t", session + ":"},
		},
	}

	shellEnv := ShellEnv{
		CmdRunner: spyRunner,
		Path:      mock.Path{},
	}

	err := Open([]string{dir}, shellEnv)

	if err != nil {
		t.Errorf("err=%v", err)
	}

	expectedCalls := [][]string{
		{"tmux", "has-session", "-t", dir + ":"},
		{"tmux", "has-session", "-t", session + ":"},
		{"tmux", "new-session", "-c", dir, "-d", "-s", session},
		{"tmux", "attach", "-t", session + ":"},
	}

	if !cmp.Equal(expectedCalls, spyRunner.Calls) {
		t.Error(cmp.Diff(expectedCalls, spyRunner.Calls))
	}
}

/*
func TestOpenDirWithProgram(t *testing.T) {
	os.Unsetenv("TMUX")

	dir := t.TempDir()
	tmuxSpy := &spy.Tmux{
		Errors: []string{"has-session", "has-session", "has-session"},
	}

	shellEnv := ShellEnv{
		Git:        mock.Git{},
		TmuxRunner: tmuxSpy,
		Path:       mock.Path{},
	}

	err := Open([]string{dir, program}, shellEnv)

	if err != nil {
		t.Errorf("err=%v", err)
	}

	session := project.Name(dir)

	expectedCalls := []spy.Call{
		{Name: "has-session", Args: tmux.Args{TargetSession: dir}},
		{Name: "has-session", Args: tmux.Args{TargetSession: session, TargetWindow: program}},
		{Name: "has-session", Args: tmux.Args{TargetSession: session}},
		{Name: "new-session", Args: tmux.Args{SessionName: session, Detach: true, WorkingDir: dir, Command: []string{program}}},
		{Name: "attach", Args: tmux.Args{TargetSession: session}},
	}

	if !cmp.Equal(expectedCalls, tmuxSpy.Calls) {
		t.Error(cmp.Diff(expectedCalls, tmuxSpy.Calls))
	}
}

func TestOpenWithExistingSession(t *testing.T) {
	t.Setenv("TMUX", "test")

	dir := t.TempDir()
	session := project.Name(dir)

	tmuxSpy := &spy.Tmux{
		Errors: []string{"has-session"},
	}

	shellEnv := ShellEnv{
		Git:        mock.Git{},
		TmuxRunner: tmuxSpy,
		Path:       mock.Path{},
	}

	err := Open([]string{dir}, shellEnv)

	if err != nil {
		t.Errorf("err=%v", err)
	}

	expectedCalls := []spy.Call{
		{Name: "has-session", Args: tmux.Args{TargetSession: dir}},
		{Name: "has-session", Args: tmux.Args{TargetSession: session}},
		{Name: "switch-client", Args: tmux.Args{TargetSession: session}},
	}

	if !cmp.Equal(expectedCalls, tmuxSpy.Calls) {
		t.Error(cmp.Diff(expectedCalls, tmuxSpy.Calls))
	}
}

func TestOpenWithoutTmux(t *testing.T) {
	os.Unsetenv("TMUX")

	tmuxSpy := &spy.Tmux{}

	shellEnv := ShellEnv{
		Git:        mock.Git{},
		TmuxRunner: tmuxSpy,
		Path:       mock.Path{Missing: []string{"tmux"}},
	}

	err := Open([]string{}, shellEnv)

	if !errors.Is(err, tmux.ErrTmuxNotInstalled) {
		t.Errorf("got=%v, want=%v", err, tmux.ErrTmuxNotInstalled)
	}
	var expectedCalls []spy.Call
	if !cmp.Equal(expectedCalls, tmuxSpy.Calls) {
		t.Error(cmp.Diff(expectedCalls, tmuxSpy.Calls))
	}
}

func TestOpenFile(t *testing.T) {
	os.Unsetenv("TMUX")

	dir := t.TempDir()
	file := dir + "/file.txt"
	os.WriteFile(file, []byte{}, 0644)

	tmuxSpy := &spy.Tmux{
		Errors: []string{"has-session"},
	}

	shellEnv := ShellEnv{
		Git:        mock.Git{},
		TmuxRunner: tmuxSpy,
		Path:       mock.Path{},
	}

	err := Open([]string{file}, shellEnv)

	if !errors.Is(err, project.ErrNotADirectory) {
		t.Errorf("got=%v, want=%v", err, project.ErrNotADirectory)
	}

	expectedCalls := []spy.Call{
		{Name: "has-session", Args: tmux.Args{TargetSession: file}},
	}

	if !cmp.Equal(expectedCalls, tmuxSpy.Calls) {
		t.Error(cmp.Diff(expectedCalls, tmuxSpy.Calls))
	}
}
*/
