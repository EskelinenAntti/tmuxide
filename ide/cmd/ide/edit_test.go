package cmd

import (
	"errors"
	"os"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/project"
	"github.com/eskelinenantti/tmuxide/internal/shell/tmux"
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
	file := dir + "/file.txt"
	os.WriteFile(file, []byte{}, 0644)

	tmuxSpy := &spy.Tmux{
		Errors: []string{"has-session", "has-session"},
	}

	shellEnv := ShellEnv{
		Git:        mock.Git{},
		TmuxRunner: tmuxSpy,
		Path:       mock.Path{},
	}

	err := Edit([]string{file}, shellEnv)

	if err != nil {
		t.Errorf("err=%v", err)
	}

	session := project.Name(dir)

	expectedCalls := []spy.Call{
		{Name: "has-session", Args: tmux.Args{TargetSession: session, TargetWindow: editor}},
		{Name: "has-session", Args: tmux.Args{TargetSession: session, TargetWindow: ""}},
		{Name: "new-session", Args: tmux.Args{SessionName: session, Detach: true, WorkingDir: dir, Command: []string{editor, file}}},
		{Name: "attach", Args: tmux.Args{TargetSession: session}},
	}

	if !cmp.Equal(tmuxSpy.Calls, expectedCalls) {
		t.Error(cmp.Diff(tmuxSpy.Calls, expectedCalls))
	}
}

func TestEditNonExistingFile(t *testing.T) {
	os.Unsetenv("TMUX")
	t.Setenv("EDITOR", editor)

	dir := t.TempDir()
	file := dir + "/file.txt"

	tmuxSpy := &spy.Tmux{}

	shellEnv := ShellEnv{
		Git:        mock.Git{},
		TmuxRunner: tmuxSpy,
		Path:       mock.Path{},
	}

	err := Edit([]string{file}, shellEnv)

	if !errors.Is(err, project.ErrInvalidPath) {
		t.Errorf("got=%v, want=%v", err, project.ErrInvalidPath)
	}
	var expectedCalls []spy.Call
	if !cmp.Equal(tmuxSpy.Calls, expectedCalls) {
		t.Error(cmp.Diff(tmuxSpy.Calls, expectedCalls))
	}
}

func TestEditDirectory(t *testing.T) {
	os.Unsetenv("TMUX")
	t.Setenv("EDITOR", editor)

	dir := t.TempDir()
	tmuxSpy := &spy.Tmux{
		Errors: []string{"has-session", "has-session"},
	}

	shellEnv := ShellEnv{
		Git:        mock.Git{},
		TmuxRunner: tmuxSpy,
		Path:       mock.Path{},
	}

	err := Edit([]string{dir}, shellEnv)

	if err != nil {
		t.Errorf("err=%v", err)
	}

	session := project.Name(dir)

	expectedCalls := []spy.Call{
		{Name: "has-session", Args: tmux.Args{TargetSession: session, TargetWindow: editor}},
		{Name: "has-session", Args: tmux.Args{TargetSession: session, TargetWindow: ""}},
		{Name: "new-session", Args: tmux.Args{SessionName: session, Detach: true, WorkingDir: dir, Command: []string{editor, dir}}},
		{Name: "attach", Args: tmux.Args{TargetSession: session}},
	}

	if !cmp.Equal(tmuxSpy.Calls, expectedCalls) {
		t.Error(cmp.Diff(tmuxSpy.Calls, expectedCalls))
	}
}

func TestEditFileInRepository(t *testing.T) {
	os.Unsetenv("TMUX")
	t.Setenv("EDITOR", editor)

	repository := t.TempDir()
	file := repository + "/file.txt"
	os.WriteFile(file, []byte{}, 0644)

	tmuxSpy := &spy.Tmux{
		Errors: []string{"has-session", "has-session"},
	}

	shellEnv := ShellEnv{
		Git:        mock.Git{Repository: repository},
		TmuxRunner: tmuxSpy,
		Path:       mock.Path{},
	}

	err := Edit([]string{file}, shellEnv)

	if err != nil {
		t.Errorf("err=%v", err)
	}

	session := project.Name(repository)

	expectedCalls := []spy.Call{
		{Name: "has-session", Args: tmux.Args{TargetSession: session, TargetWindow: editor}},
		{Name: "has-session", Args: tmux.Args{TargetSession: session, TargetWindow: ""}},
		{Name: "new-session", Args: tmux.Args{SessionName: session, Detach: true, WorkingDir: repository, Command: []string{editor, file}}},
		{Name: "attach", Args: tmux.Args{TargetSession: session}},
	}

	if !cmp.Equal(tmuxSpy.Calls, expectedCalls) {
		t.Error(cmp.Diff(tmuxSpy.Calls, expectedCalls))
	}
}

func TestEditFromAnotherSession(t *testing.T) {
	t.Setenv("TMUX", "test")
	t.Setenv("EDITOR", editor)

	dir := t.TempDir()
	session := project.Name(dir)

	tmuxSpy := &spy.Tmux{
		Errors: []string{"has-session", "has-session"},
	}

	shellEnv := ShellEnv{
		Git:        mock.Git{},
		TmuxRunner: tmuxSpy,
		Path:       mock.Path{},
	}

	err := Edit([]string{dir}, shellEnv)

	if err != nil {
		t.Errorf("err=%v", err)
	}

	expectedCalls := []spy.Call{
		{Name: "has-session", Args: tmux.Args{TargetSession: session, TargetWindow: editor}},
		{Name: "has-session", Args: tmux.Args{TargetSession: session, TargetWindow: ""}},
		{Name: "new-session", Args: tmux.Args{SessionName: session, Detach: true, WorkingDir: dir, Command: []string{editor, dir}}},
		{Name: "switch-client", Args: tmux.Args{TargetSession: session}},
	}

	if !cmp.Equal(tmuxSpy.Calls, expectedCalls) {
		t.Error(cmp.Diff(tmuxSpy.Calls, expectedCalls))
	}
}

func TestEditWithExistingWindow(t *testing.T) {
	t.Setenv("TMUX", "test")
	t.Setenv("EDITOR", editor)

	dir := t.TempDir()
	session := project.Name(dir)

	tmuxSpy := &spy.Tmux{}

	shellEnv := ShellEnv{
		Git:        mock.Git{},
		TmuxRunner: tmuxSpy,
		Path:       mock.Path{},
	}

	err := Edit([]string{dir}, shellEnv)

	if err != nil {
		t.Errorf("err=%v", err)
	}

	expectedCalls := []spy.Call{
		{Name: "has-session", Args: tmux.Args{TargetSession: session, TargetWindow: editor}},
		{Name: "new-window", Args: tmux.Args{TargetSession: session, TargetWindow: editor, WindowName: editor, Kill: true, WorkingDir: dir, Command: []string{editor, dir}}},
		{Name: "switch-client", Args: tmux.Args{TargetSession: session}},
	}

	if !cmp.Equal(tmuxSpy.Calls, expectedCalls) {
		t.Error(cmp.Diff(tmuxSpy.Calls, expectedCalls))
	}
}

func TestEditWithUnsetEditor(t *testing.T) {
	os.Unsetenv("EDITOR")
	os.Unsetenv("TMUX")

	dir := t.TempDir()

	tmuxSpy := &spy.Tmux{}

	shellEnv := ShellEnv{
		Git:        mock.Git{},
		TmuxRunner: tmuxSpy,
		Path:       mock.Path{},
	}

	err := Edit([]string{dir}, shellEnv)

	if !errors.Is(err, ErrEditorEnvNotSet) {
		t.Errorf("got=%v, want=%v", err, ErrEditorEnvNotSet)
	}
	var expectedCalls []spy.Call
	if !cmp.Equal(tmuxSpy.Calls, expectedCalls) {
		t.Error(cmp.Diff(tmuxSpy.Calls, expectedCalls))
	}
}

func TestEditWithEditorNotInstalled(t *testing.T) {
	t.Setenv("EDITOR", editor)
	os.Unsetenv("TMUX")

	dir := t.TempDir()

	tmuxSpy := &spy.Tmux{}

	shellEnv := ShellEnv{
		Git:        mock.Git{},
		TmuxRunner: tmuxSpy,
		Path:       mock.Path{Missing: []string{editor}},
	}

	err := Edit([]string{dir}, shellEnv)
	if !errors.Is(err, ErrEditorNotInstalled) {
		t.Errorf("got=%v, want=%v", err, ErrEditorNotInstalled)
	}
	var expectedCalls []spy.Call
	if !cmp.Equal(tmuxSpy.Calls, expectedCalls) {
		t.Error(cmp.Diff(tmuxSpy.Calls, expectedCalls))
	}
}
