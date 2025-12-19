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

	tmux := &spy.Tmux{}

	shell := shell.ShellEnv{
		Git:  mock.Git{},
		Tmux: tmux,
		Path: mock.Path{},
	}

	err := Open([]string{}, shell)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	session := project.Name(dir)

	expectedCalls := [][]string{
		{"HasSession", session, ""},
		{"New", session, dir},
		{"Attach", session},
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

	tmux := &spy.Tmux{}

	shell := shell.ShellEnv{
		Git:  mock.Git{Repository: repository},
		Tmux: tmux,
		Path: mock.Path{},
	}

	err := Open([]string{dir}, shell)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	session := project.Name(dir)

	expectedCalls := [][]string{
		{"HasSession", session, ""},
		{"New", session, dir},
		{"Attach", session},
	}

	if got, want := tmux.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestOpenDirWithProgram(t *testing.T) {
	os.Unsetenv("TMUX")

	dir := t.TempDir()
	tmux := &spy.Tmux{}

	shell := shell.ShellEnv{
		Git:  mock.Git{},
		Tmux: tmux,
		Path: mock.Path{},
	}

	err := Open([]string{dir, program}, shell)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	session := project.Name(dir)

	expectedCalls := [][]string{
		{"HasSession", session, program},
		{"HasSession", session, ""},
		{"New", session, dir, program},
		{"Attach", session},
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
		Session: session,
	}

	shell := shell.ShellEnv{
		Git:  mock.Git{},
		Tmux: tmux,
		Path: mock.Path{},
	}

	err := Open([]string{dir}, shell)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	expectedCalls := [][]string{
		{"HasSession", session, ""},
		{"NewWindow", session, "", dir},
		{"Switch", session},
	}

	if got, want := tmux.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestOpenWithoutTmux(t *testing.T) {
	os.Unsetenv("TMUX")

	tmuxSpy := &spy.Tmux{}

	shell := shell.ShellEnv{
		Git:  mock.Git{},
		Tmux: tmuxSpy,
		Path: mock.Path{Missing: []string{"tmux"}},
	}

	err := Open([]string{}, shell)

	if got, want := err, ide.ErrTmuxNotInstalled; !errors.Is(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
	var expectedCalls [][]string
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

	shell := shell.ShellEnv{
		Git:  mock.Git{},
		Tmux: tmux,
		Path: mock.Path{},
	}

	err := Open([]string{file}, shell)

	if got, want := err, project.ErrNotADirectory; !errors.Is(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}

	var expectedCalls [][]string
	if got, want := tmux.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}
