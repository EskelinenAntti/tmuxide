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

func unsetenv(t *testing.T, key string) {
	t.Helper()
	t.Setenv(key, "")
	if err := os.Unsetenv(key); err != nil {
		t.Fatal(err)
	}
}

func requireNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func requireCalls(t *testing.T, want, got [][]string) {
	t.Helper()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatal(diff)
	}
}

func createFile(t *testing.T, dir, name string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, nil, 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestSelectFolderFromPrompt(t *testing.T) {
	tests := []struct {
		name     string
		attached bool
	}{
		{name: "attaches to session"},
		{name: "switches to session", attached: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("EDITOR", editor)
			if tt.attached {
				t.Setenv("TMUX", "test")
			} else {
				unsetenv(t, "TMUX")
			}

			home := t.TempDir()
			t.Setenv("HOME", home)
			folder := "session"
			if err := os.Mkdir(filepath.Join(home, folder), 0755); err != nil {
				t.Fatal(err)
			}

			spyRunner := &spy.SpyRunner{
				Responses: []spy.Response{
					{OnRun: mock.WriteToStdout(folder)},
				},
			}
			err := Ide([]string{}, spyRunner, mock.Path{})
			requireNoError(t, err)

			session := project.Name(filepath.Join(home, folder))
			expectedCalls := [][]string{
				{"fzf", "--reverse", "--height", "70%", "--tmux", "70%"},
				{"fd", "--follow", "--hidden", "--exclude", "{.git,node_modules,Library}", ".", "--base-directory", home},
				{"tmux", "has-session", "-t", session + ":"},
			}
			if tt.attached {
				expectedCalls = append(expectedCalls, []string{"tmux", "switch-client", "-t", session + ":"})
			} else {
				expectedCalls = append(expectedCalls, []string{"tmux", "attach", "-t", session + ":"})
			}

			requireCalls(t, expectedCalls, spyRunner.Calls)
		})
	}
}

func TestFolderSessionWorkflow(t *testing.T) {
	tests := []struct {
		name          string
		attached      bool
		sessionExists bool
		nested        bool
	}{
		{name: "creates and attaches to session", nested: true},
		{name: "switches to existing session", attached: true, sessionExists: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("EDITOR", editor)
			if tt.attached {
				t.Setenv("TMUX", "test")
			} else {
				unsetenv(t, "TMUX")
			}

			dir := t.TempDir()
			if tt.nested {
				dir = filepath.Join(dir, "path/to/dir/in/repository")
				if err := os.MkdirAll(dir, os.ModePerm); err != nil {
					t.Fatal(err)
				}
			}
			session := project.Name(dir)

			spyRunner := &spy.SpyRunner{}
			if !tt.sessionExists {
				spyRunner.Responses = []spy.Response{{OnRun: mock.SimulateError}}
			}

			err := Ide([]string{dir}, spyRunner, mock.Path{})
			requireNoError(t, err)

			expectedCalls := [][]string{
				{"tmux", "has-session", "-t", session + ":"},
			}
			if !tt.sessionExists {
				expectedCalls = append(expectedCalls, []string{"tmux", "new-session", "-c", dir, "-d", "-s", session})
			}
			if tt.attached {
				expectedCalls = append(expectedCalls, []string{"tmux", "switch-client", "-t", session + ":"})
			} else {
				expectedCalls = append(expectedCalls, []string{"tmux", "attach", "-t", session + ":"})
			}

			requireCalls(t, expectedCalls, spyRunner.Calls)
		})
	}
}

func TestTmuxNotInstalled(t *testing.T) {
	unsetenv(t, "TMUX")
	t.Setenv("EDITOR", editor)

	spyRunner := &spy.SpyRunner{}

	err := Ide([]string{}, spyRunner, mock.Path{Missing: []string{"tmux"}})

	expectedError := shell.NotInstalledError{Cmd: "tmux"}
	var cmdNotInstalledError shell.NotInstalledError
	if !errors.As(err, &cmdNotInstalledError) {
		t.Fatalf("got=%v, want=%v", err, expectedError)
	}
	if cmdNotInstalledError != expectedError {
		t.Fatalf("got=%v, want=%v", cmdNotInstalledError, expectedError)
	}

	requireCalls(t, nil, spyRunner.Calls)
}

func TestEditorSessionWorkflow(t *testing.T) {
	tests := []struct {
		name                string
		attached            bool
		editorSessionExists bool
	}{
		{name: "creates and attaches to session"},
		{name: "creates and switches to session", attached: true},
		{name: "reuses editor and switches to session", editorSessionExists: true},
		{name: "reuses editor session", attached: true, editorSessionExists: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("EDITOR", editor)
			if tt.attached {
				t.Setenv("TMUX", "test")
			} else {
				unsetenv(t, "TMUX")
			}

			dir := t.TempDir()
			file := createFile(t, dir, "file.txt")
			session := project.Name(dir)

			responses := []spy.Response{{OnRun: mock.SimulateError}}
			if !tt.editorSessionExists {
				responses = append(responses,
					spy.Response{OnRun: mock.SimulateError},
					spy.Response{OnRun: mock.SimulateError},
				)
			}
			spyRunner := &spy.SpyRunner{Responses: responses}

			err := Ide([]string{file}, spyRunner, mock.Path{})
			requireNoError(t, err)

			expectedCalls := [][]string{
				{"git", "-C", dir, "rev-parse", "--show-toplevel"},
				{"tmux", "has-session", "-t", session + ":" + editor},
			}
			if tt.editorSessionExists {
				expectedCalls = append(expectedCalls,
					[]string{"tmux", "new-window", "-t", session + ":" + editor, "-c", dir, "-k", "-n", editor, editor, file},
				)
			} else {
				expectedCalls = append(expectedCalls,
					[]string{"tmux", "has-session", "-t", session + ":"},
					[]string{"tmux", "new-session", "-c", dir, "-d", "-s", session, editor, file},
				)
			}
			if tt.attached {
				expectedCalls = append(expectedCalls, []string{"tmux", "switch-client", "-t", session + ":"})
			} else {
				expectedCalls = append(expectedCalls, []string{"tmux", "attach", "-t", session + ":"})
			}

			requireCalls(t, expectedCalls, spyRunner.Calls)
		})
	}
}

func TestRelativePathToFile(t *testing.T) {
	unsetenv(t, "TMUX")
	t.Setenv("EDITOR", editor)

	dir := t.TempDir()
	fileName := "file.txt"
	createFile(t, dir, fileName)

	t.Chdir(dir)

	session := project.Name(dir)

	spyRunner := &spy.SpyRunner{
		Responses: []spy.Response{
			{OnRun: mock.SimulateError},
			{OnRun: mock.SimulateError},
			{OnRun: mock.SimulateError},
		},
	}

	err := Ide([]string{fileName}, spyRunner, mock.Path{})
	requireNoError(t, err)

	expectedCalls := [][]string{
		{"git", "-C", ".", "rev-parse", "--show-toplevel"},
		{"tmux", "has-session", "-t", session + ":" + editor},
		{"tmux", "has-session", "-t", session + ":"},
		{"tmux", "new-session", "-c", ".", "-d", "-s", session, editor, fileName},
		{"tmux", "attach", "-t", session + ":"},
	}

	requireCalls(t, expectedCalls, spyRunner.Calls)
}

func TestFileDoesNotExist(t *testing.T) {
	unsetenv(t, "TMUX")
	t.Setenv("EDITOR", editor)

	dir := t.TempDir()
	file := filepath.Join(dir, "file.txt")

	spyRunner := &spy.SpyRunner{}

	err := Ide([]string{file}, spyRunner, mock.Path{})

	if !errors.Is(err, project.ErrInvalidPath) {
		t.Fatalf("got=%v, want=%v", err, project.ErrInvalidPath)
	}

	requireCalls(t, nil, spyRunner.Calls)
}

func TestFileInRepository(t *testing.T) {
	unsetenv(t, "TMUX")
	t.Setenv("EDITOR", editor)

	repository := t.TempDir()
	file := createFile(t, repository, "file.txt")

	session := project.Name(repository)

	spyRunner := &spy.SpyRunner{
		Responses: []spy.Response{
			{OnRun: mock.WriteToStdout(repository)},
			{OnRun: mock.SimulateError},
			{OnRun: mock.SimulateError},
		},
	}

	err := Ide([]string{file}, spyRunner, mock.Path{})
	requireNoError(t, err)

	expectedCalls := [][]string{
		{"git", "-C", repository, "rev-parse", "--show-toplevel"},
		{"tmux", "has-session", "-t", session + ":" + editor},
		{"tmux", "has-session", "-t", session + ":"},
		{"tmux", "new-session", "-c", repository, "-d", "-s", session, editor, file},
		{"tmux", "attach", "-t", session + ":"},
	}

	requireCalls(t, expectedCalls, spyRunner.Calls)
}

func TestEditorNotSet(t *testing.T) {
	unsetenv(t, "EDITOR")
	unsetenv(t, "TMUX")

	dir := t.TempDir()

	spyRunner := &spy.SpyRunner{}

	err := Ide([]string{dir}, spyRunner, mock.Path{})

	if !errors.Is(err, ErrEditorEnvNotSet) {
		t.Fatalf("got=%v, want=%v", err, ErrEditorEnvNotSet)
	}
	requireCalls(t, nil, spyRunner.Calls)
}

func TestEditorNotInstalled(t *testing.T) {
	t.Setenv("EDITOR", editor)
	unsetenv(t, "TMUX")

	dir := t.TempDir()
	mockPath := mock.Path{Missing: []string{editor}}

	spyRunner := &spy.SpyRunner{}

	err := Ide([]string{dir}, spyRunner, &mockPath)
	if !errors.Is(err, ErrEditorNotInstalled) {
		t.Fatalf("got=%v, want=%v", err, ErrEditorNotInstalled)
	}
	requireCalls(t, nil, spyRunner.Calls)
}
