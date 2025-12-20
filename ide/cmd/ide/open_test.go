package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/ide"
	"github.com/eskelinenantti/tmuxide/internal/project"
	"github.com/eskelinenantti/tmuxide/internal/shell"
	"github.com/eskelinenantti/tmuxide/internal/test/mock"
	"github.com/eskelinenantti/tmuxide/internal/test/spy"
)

func TestOpen(t *testing.T) {
	os.Unsetenv("TMUX")

	dir := t.TempDir()
	t.Chdir(dir)

	tmux := &spy.Tmux{
		Errors: []string{"has-session"},
	}

	shellEnv := ShellEnv{
		Git:  mock.Git{},
		Tmux: tmux,
		Path: mock.Path{},
	}

	err := Open([]string{}, shellEnv)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	session := project.Name(dir)

	expectedCalls := []spy.Call{
		{Name: "has-session", Args: shell.Args{TargetSession: session}},
		{Name: "new-session", Args: shell.Args{SessionName: session, Detach: true, WorkingDir: dir}},
		{Name: "attach", Args: shell.Args{TargetSession: session}},
	}

	if got, want := tmux.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestOpenDirInsideRepository(t *testing.T) {
	os.Unsetenv("TMUX")

	repository := t.TempDir()
	dir := filepath.Join(repository, "path/to/dir/in/repository")

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		t.Fatalf("err=%v", err)
	}

	tmux := &spy.Tmux{
		Errors: []string{"has-session"},
	}

	shellEnv := ShellEnv{
		Git:  mock.Git{Repository: repository},
		Tmux: tmux,
		Path: mock.Path{},
	}

	err := Open([]string{dir}, shellEnv)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	session := project.Name(dir)

	expectedCalls := []spy.Call{
		{Name: "has-session", Args: shell.Args{TargetSession: session}},
		{Name: "new-session", Args: shell.Args{SessionName: session, Detach: true, WorkingDir: dir}},
		{Name: "attach", Args: shell.Args{TargetSession: session}},
	}

	if got, want := tmux.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestOpenDirWithProgram(t *testing.T) {
	os.Unsetenv("TMUX")

	dir := t.TempDir()
	tmux := &spy.Tmux{
		Errors: []string{"has-session", "has-session"},
	}

	shellEnv := ShellEnv{
		Git:  mock.Git{},
		Tmux: tmux,
		Path: mock.Path{},
	}

	err := Open([]string{dir, program}, shellEnv)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	session := project.Name(dir)

	expectedCalls := []spy.Call{
		{Name: "has-session", Args: shell.Args{TargetSession: session, TargetWindow: program}},
		{Name: "has-session", Args: shell.Args{TargetSession: session}},
		{Name: "new-session", Args: shell.Args{SessionName: session, Detach: true, WorkingDir: dir, Command: []string{program}}},
		{Name: "attach", Args: shell.Args{TargetSession: session}},
	}

	if got, want := tmux.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestOpenWithExistingSession(t *testing.T) {
	t.Setenv("TMUX", "test")

	dir := t.TempDir()
	session := project.Name(dir)

	tmux := &spy.Tmux{
		Errors: []string{"has-session"},
	}

	shellEnv := ShellEnv{
		Git:  mock.Git{},
		Tmux: tmux,
		Path: mock.Path{},
	}

	err := Open([]string{dir}, shellEnv)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	expectedCalls := []spy.Call{
		{Name: "has-session", Args: shell.Args{TargetSession: session}},
		{Name: "new-session", Args: shell.Args{SessionName: session, Detach: true, WorkingDir: dir}},
		{Name: "switch-client", Args: shell.Args{TargetSession: session}},
	}

	if got, want := tmux.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestOpenWithoutTmux(t *testing.T) {
	os.Unsetenv("TMUX")

	tmuxSpy := &spy.Tmux{}

	shellEnv := ShellEnv{
		Git:  mock.Git{},
		Tmux: tmuxSpy,
		Path: mock.Path{Missing: []string{"tmux"}},
	}

	err := Open([]string{}, shellEnv)

	if got, want := err, ide.ErrTmuxNotInstalled; !errors.Is(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
	var expectedCalls []spy.Call
	if got, want := tmuxSpy.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestOpenFile(t *testing.T) {
	os.Unsetenv("TMUX")

	dir := t.TempDir()
	file := dir + "/file.txt"
	os.WriteFile(file, []byte{}, 0644)

	tmux := &spy.Tmux{}

	shellEnv := ShellEnv{
		Git:  mock.Git{},
		Tmux: tmux,
		Path: mock.Path{},
	}

	err := Open([]string{file}, shellEnv)

	if got, want := err, project.ErrNotADirectory; !errors.Is(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}

	var expectedCalls []spy.Call
	if got, want := tmux.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}
