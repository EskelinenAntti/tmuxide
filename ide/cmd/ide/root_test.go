package cmd

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/project"
	"github.com/eskelinenantti/tmuxide/internal/shell"
	"github.com/eskelinenantti/tmuxide/internal/test/mock"
	"github.com/eskelinenantti/tmuxide/internal/test/spy"
	"github.com/google/go-cmp/cmp"
)

const editor string = "editor"

func TestHelpDoesNotAdvertiseExperimentalMultiplexer(t *testing.T) {
	var output bytes.Buffer
	previous := rootCmd.OutOrStdout()
	rootCmd.SetOut(&output)
	t.Cleanup(func() { rootCmd.SetOut(previous) })
	requireNoError(t, rootCmd.Help())
	for _, text := range []string{output.String(), rootCmd.UsageString(), rootCmd.Short} {
		if strings.Contains(strings.ToLower(text), "herdr") {
			t.Fatalf("help advertises Herdr: %s", text)
		}
	}
	flag := rootCmd.Flags().Lookup("mux")
	if flag == nil || flag.Hidden || flag.DefValue != "auto" {
		t.Fatal("multiplexer override must remain available without special visibility handling")
	}
}

func TestRunMuxOverride(t *testing.T) {
	t.Setenv("HERDR_ENV", "1")
	t.Setenv("HERDR_BIN_PATH", "herdr")
	t.Setenv("EDITOR", editor)
	dir := t.TempDir()
	runner := &spy.SpyRunner{}
	err := run([]string{dir}, "tmux", runner, mock.Path{Missing: []string{"tmux"}})
	var missing shell.NotInstalledError
	if !errors.As(err, &missing) || missing.Cmd != "tmux" {
		t.Fatalf("explicit multiplexer was not selected: %v", err)
	}
	requireCalls(t, nil, runner.Calls)
}

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
	t.Setenv("HERDR_ENV", "")
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
			err := run([]string{}, "auto", spyRunner, mock.Path{})
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
	t.Setenv("HERDR_ENV", "")
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

			err := run([]string{dir}, "auto", spyRunner, mock.Path{})
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
	t.Setenv("HERDR_ENV", "")
	unsetenv(t, "TMUX")
	t.Setenv("EDITOR", editor)

	spyRunner := &spy.SpyRunner{}

	err := run([]string{}, "auto", spyRunner, mock.Path{Missing: []string{"tmux"}})

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
	t.Setenv("HERDR_ENV", "")
	tests := []struct {
		name          string
		attached      bool
		sessionExists bool
	}{
		{name: "creates and attaches to session"},
		{name: "creates and switches to session", attached: true},
		{name: "opens new editor window and attaches to session", sessionExists: true},
		{name: "opens new editor window and switches to session", attached: true, sessionExists: true},
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
			if !tt.sessionExists {
				responses = append(responses, spy.Response{OnRun: mock.SimulateError})
			}
			spyRunner := &spy.SpyRunner{Responses: responses}

			err := run([]string{file}, "auto", spyRunner, mock.Path{})
			requireNoError(t, err)

			expectedCalls := [][]string{
				{"git", "-C", dir, "rev-parse", "--show-toplevel"},
				{"tmux", "has-session", "-t", session + ":"},
			}
			if tt.sessionExists {
				expectedCalls = append(expectedCalls,
					[]string{"tmux", "new-window", "-t", session + ":", "-c", dir, "-n", editor, editor, file},
				)
			} else {
				expectedCalls = append(expectedCalls,
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
	t.Setenv("HERDR_ENV", "")
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
		},
	}

	err := run([]string{fileName}, "auto", spyRunner, mock.Path{})
	requireNoError(t, err)

	expectedCalls := [][]string{
		{"git", "-C", ".", "rev-parse", "--show-toplevel"},
		{"tmux", "has-session", "-t", session + ":"},
		{"tmux", "new-session", "-c", dir, "-d", "-s", session, editor, fileName},
		{"tmux", "attach", "-t", session + ":"},
	}

	requireCalls(t, expectedCalls, spyRunner.Calls)
}

func TestFileDoesNotExist(t *testing.T) {
	t.Setenv("HERDR_ENV", "")
	unsetenv(t, "TMUX")
	t.Setenv("EDITOR", editor)

	dir := t.TempDir()
	file := filepath.Join(dir, "file.txt")

	spyRunner := &spy.SpyRunner{}

	err := run([]string{file}, "auto", spyRunner, mock.Path{})

	if !errors.Is(err, project.ErrInvalidPath) {
		t.Fatalf("got=%v, want=%v", err, project.ErrInvalidPath)
	}

	requireCalls(t, nil, spyRunner.Calls)
}

func TestFileInRepository(t *testing.T) {
	t.Setenv("HERDR_ENV", "")
	unsetenv(t, "TMUX")
	t.Setenv("EDITOR", editor)

	repository := t.TempDir()
	file := createFile(t, repository, "file.txt")

	session := project.Name(repository)

	spyRunner := &spy.SpyRunner{
		Responses: []spy.Response{
			{OnRun: mock.WriteToStdout(repository)},
			{OnRun: mock.SimulateError},
		},
	}

	err := run([]string{file}, "auto", spyRunner, mock.Path{})
	requireNoError(t, err)

	expectedCalls := [][]string{
		{"git", "-C", repository, "rev-parse", "--show-toplevel"},
		{"tmux", "has-session", "-t", session + ":"},
		{"tmux", "new-session", "-c", repository, "-d", "-s", session, editor, file},
		{"tmux", "attach", "-t", session + ":"},
	}

	requireCalls(t, expectedCalls, spyRunner.Calls)
}

func TestEditorNotSet(t *testing.T) {
	t.Setenv("HERDR_ENV", "")
	unsetenv(t, "EDITOR")
	unsetenv(t, "TMUX")

	dir := t.TempDir()

	spyRunner := &spy.SpyRunner{}

	err := run([]string{dir}, "auto", spyRunner, mock.Path{})

	if !errors.Is(err, ErrEditorEnvNotSet) {
		t.Fatalf("got=%v, want=%v", err, ErrEditorEnvNotSet)
	}
	requireCalls(t, nil, spyRunner.Calls)
}

func TestEditorNotInstalled(t *testing.T) {
	t.Setenv("HERDR_ENV", "")
	t.Setenv("EDITOR", editor)
	unsetenv(t, "TMUX")

	dir := t.TempDir()
	mockPath := mock.Path{Missing: []string{editor}}

	spyRunner := &spy.SpyRunner{}

	err := run([]string{dir}, "auto", spyRunner, &mockPath)
	if !errors.Is(err, ErrEditorNotInstalled) {
		t.Fatalf("got=%v, want=%v", err, ErrEditorNotInstalled)
	}
	requireCalls(t, nil, spyRunner.Calls)
}

func TestHerdrFolderFromPrompt(t *testing.T) {
	t.Setenv("HERDR_ENV", "1")
	t.Setenv("HERDR_BIN_PATH", "/opt/herdr")
	t.Setenv("TMUX", "outer-tmux")
	t.Setenv("EDITOR", editor)
	home := t.TempDir()
	t.Setenv("HOME", home)
	dir := filepath.Join(home, "project")
	requireNoError(t, os.Mkdir(dir, 0755))
	runner := &spy.SpyRunner{Responses: []spy.Response{
		{OnRun: mock.WriteToStdout("project")}, {},
		{OnRun: mock.WriteToStdout(`{"result":{"type":"workspace_list","workspaces":[]}}`)},
		{OnRun: mock.WriteToStdout(`{"result":{"workspace":{"workspace_id":"w1"}}}`)},
		{OnRun: mock.WriteToStdout(`{"result":{"type":"ok"}}`)},
	}}
	err := run(nil, "auto", runner, mock.Path{Missing: []string{"tmux", "herdr"}})
	requireNoError(t, err)
	requireCalls(t, [][]string{
		{"fzf", "--reverse", "--height", "70%", "--tmux", "70%"},
		{"fd", "--follow", "--hidden", "--exclude", "{.git,node_modules,Library}", ".", "--base-directory", home},
		{"/opt/herdr", "workspace", "list"},
		{"/opt/herdr", "workspace", "create", "--cwd", dir, "--label", project.Name(dir), "--no-focus"},
		{"/opt/herdr", "workspace", "focus", "w1"},
	}, runner.Calls)
}

func TestHerdrRelativeFileInRepository(t *testing.T) {
	t.Setenv("HERDR_ENV", "1")
	t.Setenv("HERDR_BIN_PATH", "herdr")
	t.Setenv("EDITOR", "editor -R")
	dir := t.TempDir()
	subdir := filepath.Join(dir, "sub")
	requireNoError(t, os.Mkdir(subdir, 0755))
	createFile(t, subdir, "file.txt")
	t.Chdir(dir)
	runner := &spy.SpyRunner{Responses: []spy.Response{
		{OnRun: mock.WriteToStdout(dir)},
		{OnRun: mock.WriteToStdout(`{"result":{"type":"workspace_list","workspaces":[]}}`)},
		{OnRun: mock.WriteToStdout(`{"result":{"workspace":{"workspace_id":"w1"},"tab":{"tab_id":"w1:t1"},"root_pane":{"pane_id":"w1:p1"}}}`)},
		{OnRun: mock.WriteToStdout(`{"result":{"type":"ok"}}`)},
		{OnRun: mock.WriteToStdout(`{"result":{"type":"ok"}}`)},
		{OnRun: mock.WriteToStdout(`{"result":{"type":"ok"}}`)},
	}}
	requireNoError(t, run([]string{"sub/file.txt"}, "auto", runner, mock.Path{Missing: []string{"tmux"}}))
	requireCalls(t, [][]string{
		{"git", "-C", "sub", "rev-parse", "--show-toplevel"},
		{"herdr", "workspace", "list"},
		{"herdr", "workspace", "create", "--cwd", dir, "--label", project.Name(dir), "--no-focus"},
		{"herdr", "pane", "run", "w1:p1", "'editor' '-R' 'sub/file.txt'"},
		{"herdr", "tab", "focus", "w1:t1"},
		{"herdr", "workspace", "focus", "w1"},
	}, runner.Calls)
}
