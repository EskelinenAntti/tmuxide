package cmd

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/project"
	"github.com/eskelinenantti/tmuxide/internal/shell/tmux"
	"github.com/eskelinenantti/tmuxide/internal/test/mock"
	"github.com/eskelinenantti/tmuxide/internal/test/spy"
)

const program string = "program"
const editor string = "editor"

func TestEdit(t *testing.T) {
	os.Unsetenv("TMUX")
	t.Setenv("EDITOR", editor)

	dir := t.TempDir()
	file := dir + "/file.txt"
	os.WriteFile(file, []byte{}, 0644)
	t.Chdir(dir)

	tmuxSpy := &spy.Tmux{
		Errors: []string{"has-session", "has-session"},
	}

	shellEnv := ShellEnv{
		Git:  mock.Git{},
		Tmux: tmuxSpy,
		Path: mock.Path{},
	}

	err := Edit([]string{}, shellEnv)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	session := project.Name(dir)
	expectedCalls := []spy.Call{
		{Name: "has-session", Args: tmux.Args{TargetSession: session, TargetWindow: editor}},
		{Name: "has-session", Args: tmux.Args{TargetSession: session, TargetWindow: ""}},
		{Name: "new-session", Args: tmux.Args{SessionName: session, Detach: true, WorkingDir: dir, Command: []string{editor, dir}}},
		{Name: "attach", Args: tmux.Args{TargetSession: session}},
	}

	if got, want := tmuxSpy.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

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
		Git:  mock.Git{},
		Tmux: tmuxSpy,
		Path: mock.Path{},
	}

	err := Edit([]string{file}, shellEnv)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	session := project.Name(dir)

	expectedCalls := []spy.Call{
		{Name: "has-session", Args: tmux.Args{TargetSession: session, TargetWindow: editor}},
		{Name: "has-session", Args: tmux.Args{TargetSession: session, TargetWindow: ""}},
		{Name: "new-session", Args: tmux.Args{SessionName: session, Detach: true, WorkingDir: dir, Command: []string{editor, file}}},
		{Name: "attach", Args: tmux.Args{TargetSession: session}},
	}

	if got, want := tmuxSpy.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestEditNonExistingFile(t *testing.T) {
	os.Unsetenv("TMUX")
	t.Setenv("EDITOR", editor)

	dir := t.TempDir()
	file := dir + "/file.txt"

	tmuxSpy := &spy.Tmux{}

	shellEnv := ShellEnv{
		Git:  mock.Git{},
		Tmux: tmuxSpy,
		Path: mock.Path{},
	}

	err := Edit([]string{file}, shellEnv)

	if got, want := err, project.ErrInvalidPath; !errors.Is(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
	var expectedCalls []spy.Call
	if got, want := tmuxSpy.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
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
		Git:  mock.Git{},
		Tmux: tmuxSpy,
		Path: mock.Path{},
	}

	err := Edit([]string{dir}, shellEnv)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	session := project.Name(dir)

	expectedCalls := []spy.Call{
		{Name: "has-session", Args: tmux.Args{TargetSession: session, TargetWindow: editor}},
		{Name: "has-session", Args: tmux.Args{TargetSession: session, TargetWindow: ""}},
		{Name: "new-session", Args: tmux.Args{SessionName: session, Detach: true, WorkingDir: dir, Command: []string{editor, dir}}},
		{Name: "attach", Args: tmux.Args{TargetSession: session}},
	}

	if got, want := tmuxSpy.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
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
		Git:  mock.Git{Repository: repository},
		Tmux: tmuxSpy,
		Path: mock.Path{},
	}

	err := Edit([]string{file}, shellEnv)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	session := project.Name(repository)

	expectedCalls := []spy.Call{
		{Name: "has-session", Args: tmux.Args{TargetSession: session, TargetWindow: editor}},
		{Name: "has-session", Args: tmux.Args{TargetSession: session, TargetWindow: ""}},
		{Name: "new-session", Args: tmux.Args{SessionName: session, Detach: true, WorkingDir: repository, Command: []string{editor, file}}},
		{Name: "attach", Args: tmux.Args{TargetSession: session}},
	}

	if got, want := tmuxSpy.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
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
		Git:  mock.Git{},
		Tmux: tmuxSpy,
		Path: mock.Path{},
	}

	err := Edit([]string{dir}, shellEnv)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	expectedCalls := []spy.Call{
		{Name: "has-session", Args: tmux.Args{TargetSession: session, TargetWindow: editor}},
		{Name: "has-session", Args: tmux.Args{TargetSession: session, TargetWindow: ""}},
		{Name: "new-session", Args: tmux.Args{SessionName: session, Detach: true, WorkingDir: dir, Command: []string{editor, dir}}},
		{Name: "switch-client", Args: tmux.Args{TargetSession: session}},
	}

	if got, want := tmuxSpy.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestEditWithExistingWindow(t *testing.T) {
	t.Setenv("TMUX", "test")
	t.Setenv("EDITOR", editor)

	dir := t.TempDir()
	session := project.Name(dir)

	tmuxSpy := &spy.Tmux{}

	shellEnv := ShellEnv{
		Git:  mock.Git{},
		Tmux: tmuxSpy,
		Path: mock.Path{},
	}

	err := Edit([]string{dir}, shellEnv)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	expectedCalls := []spy.Call{
		{Name: "has-session", Args: tmux.Args{TargetSession: session, TargetWindow: editor}},
		{Name: "new-window", Args: tmux.Args{TargetSession: session, TargetWindow: editor, WindowName: editor, Kill: true, WorkingDir: dir, Command: []string{editor, dir}}},
		{Name: "switch-client", Args: tmux.Args{TargetSession: session}},
	}

	if got, want := tmuxSpy.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestEditWithUnsetEditor(t *testing.T) {
	os.Unsetenv("EDITOR")
	os.Unsetenv("TMUX")

	dir := t.TempDir()

	tmuxSpy := &spy.Tmux{}

	shellEnv := ShellEnv{
		Git:  mock.Git{},
		Tmux: tmuxSpy,
		Path: mock.Path{},
	}

	err := Edit([]string{dir}, shellEnv)

	if got, want := err, ErrEditorEnvNotSet; !errors.Is(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
	var expectedCalls []spy.Call
	if got, want := tmuxSpy.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestEditWithEditorNotInstalled(t *testing.T) {
	t.Setenv("EDITOR", editor)
	os.Unsetenv("TMUX")

	dir := t.TempDir()

	tmuxSpy := &spy.Tmux{}

	shellEnv := ShellEnv{
		Git:  mock.Git{},
		Tmux: tmuxSpy,
		Path: mock.Path{Missing: []string{editor}},
	}

	err := Edit([]string{dir}, shellEnv)

	if got, want := err, ErrEditorNotInstalled; !errors.Is(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
	var expectedCalls []spy.Call
	if got, want := tmuxSpy.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}
