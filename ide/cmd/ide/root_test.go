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

const editor string = "editor"

func TestSelectFolderFromPrompt(t *testing.T) {
	os.Unsetenv("TMUX")
	t.Setenv("EDITOR", editor)
	folder := "session"
	spyRunner := &spy.SpyRunner{
		Mocks: []spy.Mock{{
			Args: []string{
				"fzf", "--reverse", "--height", "70%", "--tmux", "70%",
			},
			OnRun: mock.WriteToStdout(folder),
		}},
	}
	err := Ide([]string{}, spyRunner, mock.Path{})
	if err != nil {
		t.Errorf("err=%v", err)
	}

	selectedPath := filepath.Join(os.Getenv("HOME"), folder)
	expectedCalls := [][]string{
		{"fzf", "--reverse", "--height", "70%", "--tmux", "70%"},
		{"tmux", "list-sessions", "-F", "Session: #S"},
		{"fd", "--follow", "--hidden", "--exclude", "{.git,node_modules,target,build,Library}", ".", "--base-directory", os.Getenv("HOME")},
		{"tmux", "has-session", "-t", selectedPath + ":"},
		{"tmux", "has-session", "-t", selectedPath + ":"},
		{"tmux", "attach", "-t", selectedPath + ":"},
	}

	if !cmp.Equal(expectedCalls, spyRunner.Calls) {
		t.Error(cmp.Diff(expectedCalls, spyRunner.Calls))
	}
}

func TestSelectSessionFromPrompt(t *testing.T) {
	os.Unsetenv("TMUX")
	t.Setenv("EDITOR", editor)
	session := "test-session"
	selection := "Session: " + session

	spyRunner := &spy.SpyRunner{
		Mocks: []spy.Mock{{
			Args: []string{
				"fzf", "--reverse", "--height", "70%", "--tmux", "70%",
			},
			OnRun: mock.WriteToStdout(selection),
		}},
	}
	err := Ide([]string{}, spyRunner, mock.Path{})
	if err != nil {
		t.Errorf("err=%v", err)
	}

	selectedPath := session
	expectedCalls := [][]string{
		{"fzf", "--reverse", "--height", "70%", "--tmux", "70%"},
		{"tmux", "list-sessions", "-F", "Session: #S"},
		{"fd", "--follow", "--hidden", "--exclude", "{.git,node_modules,target,build,Library}", ".", "--base-directory", os.Getenv("HOME")},
		{"tmux", "has-session", "-t", selectedPath + ":"},
		{"tmux", "has-session", "-t", selectedPath + ":"},
		{"tmux", "attach", "-t", selectedPath + ":"},
	}

	if !cmp.Equal(expectedCalls, spyRunner.Calls) {
		t.Error(cmp.Diff(expectedCalls, spyRunner.Calls))
	}
}

func TestSelectFolderFromPromptWhenAttachedToSession(t *testing.T) {
	t.Setenv("TMUX", "test")
	t.Setenv("EDITOR", editor)
	session := "session"
	spyRunner := &spy.SpyRunner{
		Mocks: []spy.Mock{{
			Args: []string{
				"fzf", "--reverse", "--height", "70%", "--tmux", "70%",
			},
			OnRun: mock.WriteToStdout(session),
		}},
	}

	err := Ide([]string{}, spyRunner, mock.Path{})
	if err != nil {
		t.Errorf("err=%v", err)
	}

	selectedPath := filepath.Join(os.Getenv("HOME"), session)
	expectedCalls := [][]string{
		{"fzf", "--reverse", "--height", "70%", "--tmux", "70%"},
		{"tmux", "list-sessions", "-F", "Session: #S"},
		{"fd", "--follow", "--hidden", "--exclude", "{.git,node_modules,target,build,Library}", ".", "--base-directory", os.Getenv("HOME")},
		{"tmux", "has-session", "-t", selectedPath + ":"},
		{"tmux", "has-session", "-t", selectedPath + ":"},
		{"tmux", "switch-client", "-t", selectedPath + ":"},
	}

	if !cmp.Equal(expectedCalls, spyRunner.Calls) {
		t.Error(cmp.Diff(expectedCalls, spyRunner.Calls))
	}
}

func TestTargetDirInsideRepository(t *testing.T) {
	os.Unsetenv("TMUX")
	t.Setenv("EDITOR", editor)

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

	err := Ide([]string{dir}, spyRunner, mock.Path{})

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

func TestSessionExistsForTargetDir(t *testing.T) {
	t.Setenv("TMUX", "test")
	t.Setenv("EDITOR", editor)

	dir := t.TempDir()
	session := project.Name(dir)

	spyRunner := &spy.SpyRunner{
		Mocks: []spy.Mock{
			{Args: []string{"tmux", "has-session", "-t", dir + ":"}, OnRun: mock.SimulateError},
		},
	}

	err := Ide([]string{dir}, spyRunner, mock.Path{})

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

func TestTmuxNotInstalled(t *testing.T) {
	os.Unsetenv("TMUX")
	t.Setenv("EDITOR", editor)

	spyRunner := &spy.SpyRunner{}

	err := Ide([]string{}, spyRunner, mock.Path{Missing: []string{"tmux"}})

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

func TestFile(t *testing.T) {
	os.Unsetenv("TMUX")
	t.Setenv("EDITOR", editor)

	dir := t.TempDir()
	file := filepath.Join(dir, "file.txt")
	os.WriteFile(file, []byte{}, 0644)

	session := project.Name(dir)

	spyRunner := &spy.SpyRunner{
		Mocks: []spy.Mock{
			{Args: []string{"git", "-C", dir, "rev-parse", "--show-toplevel"}, OnRun: mock.SimulateError},
			{Args: []string{"tmux", "has-session", "-t", file + ":"}, OnRun: mock.SimulateError},
			{Args: []string{"tmux", "has-session", "-t", session + ":" + editor}, OnRun: mock.SimulateError},
			{Args: []string{"tmux", "has-session", "-t", session + ":"}, OnRun: mock.SimulateError},
		},
	}

	err := Ide([]string{file}, spyRunner, mock.Path{})

	if err != nil {
		t.Errorf("err=%v", err)
	}

	expectedCalls := [][]string{
		{"tmux", "has-session", "-t", file + ":"},
		{"git", "-C", dir, "rev-parse", "--show-toplevel"},
		{"tmux", "has-session", "-t", session + ":" + editor},
		{"tmux", "has-session", "-t", session + ":"},
		{"tmux", "new-session", "-c", dir, "-d", "-s", session, editor, file},
		{"tmux", "attach", "-t", session + ":"},
	}

	if !cmp.Equal(expectedCalls, spyRunner.Calls) {
		t.Error(cmp.Diff(expectedCalls, spyRunner.Calls))
	}
}

func TestRelativePathToFile(t *testing.T) {
	os.Unsetenv("TMUX")
	t.Setenv("EDITOR", editor)

	dir := t.TempDir()
	fileName := "file.txt"
	file := filepath.Join(dir, fileName)
	os.WriteFile(file, []byte{}, 0644)

	t.Chdir(dir)

	session := project.Name(dir)

	spyRunner := &spy.SpyRunner{
		Mocks: []spy.Mock{
			{Args: []string{"git", "-C", ".", "rev-parse", "--show-toplevel"}, OnRun: mock.SimulateError},
			{Args: []string{"tmux", "has-session", "-t", fileName + ":"}, OnRun: mock.SimulateError},
			{Args: []string{"tmux", "has-session", "-t", session + ":" + editor}, OnRun: mock.SimulateError},
			{Args: []string{"tmux", "has-session", "-t", session + ":"}, OnRun: mock.SimulateError},
		},
	}

	err := Ide([]string{fileName}, spyRunner, mock.Path{})

	if err != nil {
		t.Errorf("err=%v", err)
	}

	expectedCalls := [][]string{
		{"tmux", "has-session", "-t", fileName + ":"},
		{"git", "-C", ".", "rev-parse", "--show-toplevel"},
		{"tmux", "has-session", "-t", session + ":" + editor},
		{"tmux", "has-session", "-t", session + ":"},
		{"tmux", "new-session", "-c", ".", "-d", "-s", session, editor, fileName},
		{"tmux", "attach", "-t", session + ":"},
	}

	if !cmp.Equal(expectedCalls, spyRunner.Calls) {
		t.Error(cmp.Diff(expectedCalls, spyRunner.Calls))
	}
}

func TestFileDoesNotExist(t *testing.T) {
	os.Unsetenv("TMUX")
	t.Setenv("EDITOR", editor)

	dir := t.TempDir()
	file := filepath.Join(dir, "file.txt")

	spyRunner := &spy.SpyRunner{
		Mocks: []spy.Mock{
			{Args: []string{"tmux", "has-session", "-t", file + ":"}, OnRun: mock.SimulateError},
		},
	}

	err := Ide([]string{file}, spyRunner, mock.Path{})

	if !errors.Is(err, project.ErrInvalidPath) {
		t.Errorf("got=%v, want=%v", err, project.ErrInvalidPath)
	}

	expectedCalls := [][]string{
		{"tmux", "has-session", "-t", file + ":"},
	}

	if !cmp.Equal(expectedCalls, spyRunner.Calls) {
		t.Error(cmp.Diff(expectedCalls, spyRunner.Calls))
	}
}

func TestFileInRepository(t *testing.T) {
	os.Unsetenv("TMUX")
	t.Setenv("EDITOR", editor)

	repository := t.TempDir()
	file := filepath.Join(repository, "file.txt")
	os.WriteFile(file, []byte{}, 0644)

	session := project.Name(repository)

	spyRunner := &spy.SpyRunner{
		Mocks: []spy.Mock{
			{Args: []string{"git", "-C", repository, "rev-parse", "--show-toplevel"}, OnRun: mock.WriteToStdout(repository)},
			{Args: []string{"tmux", "has-session", "-t", file + ":"}, OnRun: mock.SimulateError},
			{Args: []string{"tmux", "has-session", "-t", session + ":" + editor}, OnRun: mock.SimulateError},
			{Args: []string{"tmux", "has-session", "-t", session + ":"}, OnRun: mock.SimulateError},
		},
	}

	err := Ide([]string{file}, spyRunner, mock.Path{})

	if err != nil {
		t.Errorf("err=%v", err)
	}

	expectedCalls := [][]string{
		{"tmux", "has-session", "-t", file + ":"},
		{"git", "-C", repository, "rev-parse", "--show-toplevel"},
		{"tmux", "has-session", "-t", session + ":" + editor},
		{"tmux", "has-session", "-t", session + ":"},
		{"tmux", "new-session", "-c", repository, "-d", "-s", session, editor, file},
		{"tmux", "attach", "-t", session + ":"},
	}

	if !cmp.Equal(expectedCalls, spyRunner.Calls) {
		t.Error(cmp.Diff(expectedCalls, spyRunner.Calls))
	}
}

func TestFileWhenInAnotherSession(t *testing.T) {
	t.Setenv("TMUX", "test")
	t.Setenv("EDITOR", editor)

	dir := t.TempDir()
	file := filepath.Join(dir, "file.txt")
	os.WriteFile(file, []byte{}, 0644)

	session := project.Name(dir)

	spyRunner := &spy.SpyRunner{
		Mocks: []spy.Mock{
			{Args: []string{"git", "-C", dir, "rev-parse", "--show-toplevel"}, OnRun: mock.SimulateError},
			{Args: []string{"tmux", "has-session", "-t", file + ":"}, OnRun: mock.SimulateError},
			{Args: []string{"tmux", "has-session", "-t", session + ":" + editor}, OnRun: mock.SimulateError},
			{Args: []string{"tmux", "has-session", "-t", session + ":"}, OnRun: mock.SimulateError},
		},
	}

	err := Ide([]string{file}, spyRunner, mock.Path{})

	if err != nil {
		t.Errorf("err=%v", err)
	}

	expectedCalls := [][]string{
		{"tmux", "has-session", "-t", file + ":"},
		{"git", "-C", dir, "rev-parse", "--show-toplevel"},
		{"tmux", "has-session", "-t", session + ":" + editor},
		{"tmux", "has-session", "-t", session + ":"},
		{"tmux", "new-session", "-c", dir, "-d", "-s", session, editor, file},
		{"tmux", "switch-client", "-t", session + ":"},
	}

	if !cmp.Equal(expectedCalls, spyRunner.Calls) {
		t.Error(cmp.Diff(expectedCalls, spyRunner.Calls))
	}
}

func TestFileWhenSessionWithEditorExist(t *testing.T) {
	t.Setenv("TMUX", "test")
	t.Setenv("EDITOR", editor)

	dir := t.TempDir()
	file := filepath.Join(dir, "file.txt")
	os.WriteFile(file, []byte{}, 0644)

	session := project.Name(dir)

	spyRunner := &spy.SpyRunner{
		Mocks: []spy.Mock{
			{Args: []string{"git", "-C", dir, "rev-parse", "--show-toplevel"}, OnRun: mock.SimulateError},
			{Args: []string{"tmux", "has-session", "-t", file + ":"}, OnRun: mock.SimulateError},
		},
	}

	err := Ide([]string{file}, spyRunner, mock.Path{})

	if err != nil {
		t.Errorf("err=%v", err)
	}

	expectedCalls := [][]string{
		{"tmux", "has-session", "-t", file + ":"},
		{"git", "-C", dir, "rev-parse", "--show-toplevel"},
		{"tmux", "has-session", "-t", session + ":" + editor},
		{"tmux", "new-window", "-t", session + ":" + editor, "-c", dir, "-k", "-n", editor, "editor", file},
		{"tmux", "switch-client", "-t", session + ":"},
	}

	if !cmp.Equal(expectedCalls, spyRunner.Calls) {
		t.Error(cmp.Diff(expectedCalls, spyRunner.Calls))
	}
}

func TestEditorNotSet(t *testing.T) {
	os.Unsetenv("EDITOR")
	os.Unsetenv("TMUX")

	dir := t.TempDir()

	spyRunner := &spy.SpyRunner{
		Mocks: []spy.Mock{
			{Args: []string{"git", "-C", dir, "rev-parse", "--show-toplevel"}, OnRun: mock.SimulateError},
		},
	}

	err := Ide([]string{dir}, spyRunner, mock.Path{})

	if !errors.Is(err, ErrEditorEnvNotSet) {
		t.Errorf("got=%v, want=%v", err, ErrEditorEnvNotSet)
	}
	var expectedCalls [][]string
	if !cmp.Equal(expectedCalls, spyRunner.Calls) {
		t.Error(cmp.Diff(expectedCalls, spyRunner.Calls))
	}
}

func TestEditorNotInstalled(t *testing.T) {
	t.Setenv("EDITOR", editor)
	os.Unsetenv("TMUX")

	dir := t.TempDir()
	mockPath := mock.Path{Missing: []string{editor}}

	spyRunner := &spy.SpyRunner{
		Mocks: []spy.Mock{
			{Args: []string{"git", "-C", dir, "rev-parse", "--show-toplevel"}, OnRun: mock.SimulateError},
		},
	}

	err := Ide([]string{dir}, spyRunner, &mockPath)
	if !errors.Is(err, ErrEditorNotInstalled) {
		t.Errorf("got=%v, want=%v", err, ErrEditorNotInstalled)
	}
	var expectedCalls [][]string
	if !cmp.Equal(expectedCalls, spyRunner.Calls) {
		t.Error(cmp.Diff(expectedCalls, spyRunner.Calls))
	}
}
