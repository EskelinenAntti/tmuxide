package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/project"
	"github.com/eskelinenantti/tmuxide/internal/test/mock"
	"github.com/eskelinenantti/tmuxide/internal/test/spy"
	"github.com/google/go-cmp/cmp"
)

const program string = "program"
const editor string = "editor"

func TestEditFile(t *testing.T) {
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

	shellEnv := ShellEnv{
		CmdRunner: spyRunner,
		Path:      mock.Path{},
	}

	err := Edit([]string{file}, shellEnv)

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

func TestEditRelativeFile(t *testing.T) {
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

	shellEnv := ShellEnv{
		CmdRunner: spyRunner,
		Path:      mock.Path{},
	}

	err := Edit([]string{fileName}, shellEnv)

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

func TestEditNonExistingFile(t *testing.T) {
	os.Unsetenv("TMUX")
	t.Setenv("EDITOR", editor)

	dir := t.TempDir()
	file := filepath.Join(dir, "file.txt")

	spyRunner := &spy.SpyRunner{
		Mocks: []spy.Mock{
			{Args: []string{"tmux", "has-session", "-t", file + ":"}, OnRun: mock.SimulateError},
		},
	}

	shellEnv := ShellEnv{
		CmdRunner: spyRunner,
		Path:      mock.Path{},
	}

	err := Edit([]string{file}, shellEnv)

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

func TestEditDirectory(t *testing.T) {
	os.Unsetenv("TMUX")
	t.Setenv("EDITOR", editor)

	dir := t.TempDir()
	session := project.Name(dir)

	spyRunner := &spy.SpyRunner{
		Mocks: []spy.Mock{
			{Args: []string{"git", "-C", dir, "rev-parse", "--show-toplevel"}, OnRun: mock.SimulateError},
			{Args: []string{"tmux", "has-session", "-t", dir + ":"}, OnRun: mock.SimulateError},
			{Args: []string{"tmux", "has-session", "-t", session + ":" + editor}, OnRun: mock.SimulateError},
			{Args: []string{"tmux", "has-session", "-t", session + ":"}, OnRun: mock.SimulateError},
		},
	}

	shellEnv := ShellEnv{
		CmdRunner: spyRunner,
		Path:      mock.Path{},
	}

	err := Edit([]string{dir}, shellEnv)

	if err != nil {
		t.Errorf("err=%v", err)
	}

	expectedCalls := [][]string{
		{"tmux", "has-session", "-t", dir + ":"},
		{"git", "-C", dir, "rev-parse", "--show-toplevel"},
		{"tmux", "has-session", "-t", session + ":" + editor},
		{"tmux", "has-session", "-t", session + ":"},
		{"tmux", "new-session", "-c", dir, "-d", "-s", session, editor, dir},
		{"tmux", "attach", "-t", session + ":"},
	}

	if !cmp.Equal(expectedCalls, spyRunner.Calls) {
		t.Error(cmp.Diff(expectedCalls, spyRunner.Calls))
	}
}

func TestEditFileInRepository(t *testing.T) {
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

	shellEnv := ShellEnv{
		CmdRunner: spyRunner,
		Path:      mock.Path{},
	}

	err := Edit([]string{file}, shellEnv)

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

func TestEditFromAnotherSession(t *testing.T) {
	t.Setenv("TMUX", "test")
	t.Setenv("EDITOR", editor)

	dir := t.TempDir()
	session := project.Name(dir)

	spyRunner := &spy.SpyRunner{
		Mocks: []spy.Mock{
			{Args: []string{"git", "-C", dir, "rev-parse", "--show-toplevel"}, OnRun: mock.SimulateError},
			{Args: []string{"tmux", "has-session", "-t", dir + ":"}, OnRun: mock.SimulateError},
			{Args: []string{"tmux", "has-session", "-t", session + ":" + editor}, OnRun: mock.SimulateError},
			{Args: []string{"tmux", "has-session", "-t", session + ":"}, OnRun: mock.SimulateError},
		},
	}

	shellEnv := ShellEnv{
		CmdRunner: spyRunner,
		Path:      mock.Path{},
	}

	err := Edit([]string{dir}, shellEnv)

	if err != nil {
		t.Errorf("err=%v", err)
	}

	expectedCalls := [][]string{
		{"tmux", "has-session", "-t", dir + ":"},
		{"git", "-C", dir, "rev-parse", "--show-toplevel"},
		{"tmux", "has-session", "-t", session + ":" + editor},
		{"tmux", "has-session", "-t", session + ":"},
		{"tmux", "new-session", "-c", dir, "-d", "-s", session, editor, dir},
		{"tmux", "switch-client", "-t", session + ":"},
	}

	if !cmp.Equal(expectedCalls, spyRunner.Calls) {
		t.Error(cmp.Diff(expectedCalls, spyRunner.Calls))
	}
}

func TestEditWithExistingWindow(t *testing.T) {
	t.Setenv("TMUX", "test")
	t.Setenv("EDITOR", editor)

	dir := t.TempDir()
	session := project.Name(dir)

	spyRunner := &spy.SpyRunner{
		Mocks: []spy.Mock{
			{Args: []string{"git", "-C", dir, "rev-parse", "--show-toplevel"}, OnRun: mock.SimulateError},
			{Args: []string{"tmux", "has-session", "-t", dir + ":"}, OnRun: mock.SimulateError},
		},
	}

	shellEnv := ShellEnv{
		CmdRunner: spyRunner,
		Path:      mock.Path{},
	}

	err := Edit([]string{dir}, shellEnv)

	if err != nil {
		t.Errorf("err=%v", err)
	}

	expectedCalls := [][]string{
		{"tmux", "has-session", "-t", dir + ":"},
		{"git", "-C", dir, "rev-parse", "--show-toplevel"},
		{"tmux", "has-session", "-t", session + ":" + editor},
		{"tmux", "new-window", "-t", session + ":" + editor, "-c", dir, "-k", "-n", editor, "editor", dir},
		{"tmux", "switch-client", "-t", session + ":"},
	}

	if !cmp.Equal(expectedCalls, spyRunner.Calls) {
		t.Error(cmp.Diff(expectedCalls, spyRunner.Calls))
	}
}

func TestEditWithUnsetEditor(t *testing.T) {
	os.Unsetenv("EDITOR")
	os.Unsetenv("TMUX")

	dir := t.TempDir()

	spyRunner := &spy.SpyRunner{
		Mocks: []spy.Mock{
			{Args: []string{"git", "-C", dir, "rev-parse", "--show-toplevel"}, OnRun: mock.SimulateError},
		},
	}

	shellEnv := ShellEnv{
		CmdRunner: spyRunner,
		Path:      mock.Path{},
	}

	err := Edit([]string{dir}, shellEnv)

	if !errors.Is(err, ErrEditorEnvNotSet) {
		t.Errorf("got=%v, want=%v", err, ErrEditorEnvNotSet)
	}
	var expectedCalls [][]string
	if !cmp.Equal(expectedCalls, spyRunner.Calls) {
		t.Error(cmp.Diff(expectedCalls, spyRunner.Calls))
	}
}

func TestEditWithEditorNotInstalled(t *testing.T) {
	t.Setenv("EDITOR", editor)
	os.Unsetenv("TMUX")

	dir := t.TempDir()
	mockPath := mock.Path{Missing: []string{editor}}

	spyRunner := &spy.SpyRunner{
		Mocks: []spy.Mock{
			{Args: []string{"git", "-C", dir, "rev-parse", "--show-toplevel"}, OnRun: mock.SimulateError},
		},
	}

	shellEnv := ShellEnv{
		CmdRunner: spyRunner,
		Path:      &mockPath,
	}

	err := Edit([]string{dir}, shellEnv)
	if !errors.Is(err, ErrEditorNotInstalled) {
		t.Errorf("got=%v, want=%v", err, ErrEditorNotInstalled)
	}
	var expectedCalls [][]string
	if !cmp.Equal(expectedCalls, spyRunner.Calls) {
		t.Error(cmp.Diff(expectedCalls, spyRunner.Calls))
	}
}
