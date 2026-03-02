package cmd

import (
	"errors"
	"os"

	"path/filepath"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/project"
	"github.com/eskelinenantti/tmuxide/internal/shell"
	"github.com/eskelinenantti/tmuxide/internal/test/mock"
	"github.com/eskelinenantti/tmuxide/internal/test/spy"
	"github.com/google/go-cmp/cmp"
)

func TestOpen(t *testing.T) {
	os.Unsetenv("TMUX")
	session := "session"
	spyRunner := &spy.SpyRunner{
		Mocks: []spy.Mock{{
			Args: []string{
				"fzf", "--reverse", "--height", "30%",
			},
			OnRun: mock.WriteToStdout(session),
		}},
	}
	err := Open([]string{}, spyRunner, mock.Path{})
	if err != nil {
		t.Errorf("err=%v", err)
	}

	selectedPath := filepath.Join(os.Getenv("HOME"), session)
	expectedCalls := [][]string{
		{"fzf", "--reverse", "--height", "30%"},
		{"tmux", "list-sessions", "-F", "Session: #S"},
		{"fd", "--type", "dir", "--follow", "--hidden", "--exclude", "{.git,node_modules,target,build,Library}", ".", "--base-directory", os.Getenv("HOME")},
		{"tmux", "has-session", "-t", selectedPath + ":"},
		{"tmux", "has-session", "-t", selectedPath + ":"},
		{"tmux", "attach", "-t", selectedPath + ":"},
	}

	if !cmp.Equal(expectedCalls, spyRunner.Calls) {
		t.Error(cmp.Diff(expectedCalls, spyRunner.Calls))
	}
}
func TestOpenWhenSelectsSession(t *testing.T) {
	os.Unsetenv("TMUX")
	session := "test-session"
	selection := "Session: " + session

	spyRunner := &spy.SpyRunner{
		Mocks: []spy.Mock{{
			Args: []string{
				"fzf", "--reverse", "--height", "30%",
			},
			OnRun: mock.WriteToStdout(selection),
		}},
	}
	err := Open([]string{}, spyRunner, mock.Path{})
	if err != nil {
		t.Errorf("err=%v", err)
	}

	selectedPath := session
	expectedCalls := [][]string{
		{"fzf", "--reverse", "--height", "30%"},
		{"tmux", "list-sessions", "-F", "Session: #S"},
		{"fd", "--type", "dir", "--follow", "--hidden", "--exclude", "{.git,node_modules,target,build,Library}", ".", "--base-directory", os.Getenv("HOME")},
		{"tmux", "has-session", "-t", selectedPath + ":"},
		{"tmux", "has-session", "-t", selectedPath + ":"},
		{"tmux", "attach", "-t", selectedPath + ":"},
	}

	if !cmp.Equal(expectedCalls, spyRunner.Calls) {
		t.Error(cmp.Diff(expectedCalls, spyRunner.Calls))
	}
}

func TestOpenWhenAttached(t *testing.T) {
	t.Setenv("TMUX", "test")
	session := "session"
	spyRunner := &spy.SpyRunner{
		Mocks: []spy.Mock{{
			Args: []string{
				"fzf", "--reverse", "--height", "30%",
			},
			OnRun: mock.WriteToStdout(session),
		}},
	}

	err := Open([]string{}, spyRunner, mock.Path{})
	if err != nil {
		t.Errorf("err=%v", err)
	}

	selectedPath := filepath.Join(os.Getenv("HOME"), session)
	expectedCalls := [][]string{
		{"fzf", "--reverse", "--height", "30%"},
		{"tmux", "list-sessions", "-F", "Session: #S"},
		{"fd", "--type", "dir", "--follow", "--hidden", "--exclude", "{.git,node_modules,target,build,Library}", ".", "--base-directory", os.Getenv("HOME")},
		{"tmux", "has-session", "-t", selectedPath + ":"},
		{"tmux", "has-session", "-t", selectedPath + ":"},
		{"tmux", "switch-client", "-t", selectedPath + ":"},
	}

	if !cmp.Equal(expectedCalls, spyRunner.Calls) {
		t.Error(cmp.Diff(expectedCalls, spyRunner.Calls))
	}
}

func TestOpenDirInsideRepository(t *testing.T) {
	os.Unsetenv("TMUX")

	repository := t.TempDir()
	dir := filepath.Join(repository, "path/to/dir/in/repository")

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		t.Errorf("err=%v", err)
	}

	session := project.Name(dir)

	spyRunner := &spy.SpyRunner{
		Mocks: []spy.Mock{
			{Args: []string{"tmux", "has-session", "-t", dir + ":"}, OnRun: mock.SimulateError},
			{Args: []string{"tmux", "has-session", "-t", session + ":"}, OnRun: mock.SimulateError},
		},
	}

	err := Open([]string{dir}, spyRunner, mock.Path{})

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

func TestOpenDirWithProgram(t *testing.T) {
	os.Unsetenv("TMUX")

	dir := t.TempDir()
	session := project.Name(dir)
	spyRunner := &spy.SpyRunner{
		Mocks: []spy.Mock{
			{Args: []string{"tmux", "has-session", "-t", dir + ":"}, OnRun: mock.SimulateError},
			{Args: []string{"tmux", "has-session", "-t", session + ":" + program}, OnRun: mock.SimulateError},
			{Args: []string{"tmux", "has-session", "-t", session + ":"}, OnRun: mock.SimulateError},
		},
	}

	err := Open([]string{dir, program}, spyRunner, mock.Path{})

	if err != nil {
		t.Errorf("err=%v", err)
	}

	expectedCalls := [][]string{
		{"tmux", "has-session", "-t", dir + ":"},
		{"tmux", "has-session", "-t", session + ":" + program},
		{"tmux", "has-session", "-t", session + ":"},
		{"tmux", "new-session", "-c", dir, "-d", "-s", session, program},
		{"tmux", "attach", "-t", session + ":"},
	}

	if !cmp.Equal(expectedCalls, spyRunner.Calls) {
		t.Error(cmp.Diff(expectedCalls, spyRunner.Calls))
	}
}

func TestOpenWithExistingSession(t *testing.T) {
	t.Setenv("TMUX", "test")

	dir := t.TempDir()
	session := project.Name(dir)

	spyRunner := &spy.SpyRunner{
		Mocks: []spy.Mock{
			{Args: []string{"tmux", "has-session", "-t", dir + ":"}, OnRun: mock.SimulateError},
		},
	}

	err := Open([]string{dir}, spyRunner, mock.Path{})

	if err != nil {
		t.Errorf("err=%v", err)
	}

	expectedCalls := [][]string{
		{"tmux", "has-session", "-t", dir + ":"},
		{"tmux", "has-session", "-t", session + ":"},
		{"tmux", "switch-client", "-t", session + ":"},
	}

	if !cmp.Equal(expectedCalls, spyRunner.Calls) {
		t.Error(cmp.Diff(expectedCalls, spyRunner.Calls))
	}
}

func TestOpenWithoutTmux(t *testing.T) {
	os.Unsetenv("TMUX")

	spyRunner := &spy.SpyRunner{}

	err := Open([]string{}, spyRunner, mock.Path{Missing: []string{"tmux"}})

	expectedError := shell.NotInstalledError{Cmd: "tmux"}
	var cmdNotInstalledError shell.NotInstalledError
	if !errors.As(err, &cmdNotInstalledError) {
		t.Errorf("got=%v, want=%v", err, expectedError)
	}
	if !cmp.Equal(cmdNotInstalledError, expectedError) {
		t.Error(cmp.Diff(cmdNotInstalledError, expectedError))
	}

	var expectedCalls [][]string
	if !cmp.Equal(expectedCalls, spyRunner.Calls) {
		t.Error(cmp.Diff(expectedCalls, spyRunner.Calls))
	}
}

func TestOpenFile(t *testing.T) {
	os.Unsetenv("TMUX")

	dir := t.TempDir()
	file := dir + "/file.txt"
	os.WriteFile(file, []byte{}, 0644)

	spyRunner := &spy.SpyRunner{
		Mocks: []spy.Mock{
			{Args: []string{"tmux", "has-session", "-t", file + ":"}, OnRun: mock.SimulateError},
		},
	}

	err := Open([]string{file}, spyRunner, mock.Path{})

	if !errors.Is(err, project.ErrNotADirectory) {
		t.Errorf("got=%v, want=%v", err, project.ErrNotADirectory)
	}

	expectedCalls := [][]string{
		{"tmux", "has-session", "-t", file + ":"},
	}

	if !cmp.Equal(expectedCalls, spyRunner.Calls) {
		t.Error(cmp.Diff(expectedCalls, spyRunner.Calls))
	}
}
