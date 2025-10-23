package cmd

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/ide"
	"github.com/eskelinenantti/tmuxide/internal/shell"
	"github.com/eskelinenantti/tmuxide/internal/test"
	"github.com/eskelinenantti/tmuxide/internal/tmux"
)

var testEditor string = "editor"

func TestRunWithDirectoryArgument(t *testing.T) {
	t.Setenv("EDITOR", testEditor)
	os.Unsetenv("TMUX")

	dir := t.TempDir()
	tmux := &test.TmuxSpy{}

	shell := shell.Shell{
		Git:  test.GitMock{},
		Tmux: tmux,
		Path: test.PathMock{},
	}

	err := run([]string{dir}, shell)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	session := ide.Name(dir)

	expectedCalls := [][]string{
		{"HasSession", session},
		{"New", session, dir, testEditor, dir},
		{"Attach", session},
	}

	if got, want := tmux.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestRunWithoutArguments(t *testing.T) {
	t.Setenv("EDITOR", testEditor)
	os.Unsetenv("TMUX")

	dir := t.TempDir()
	t.Chdir(dir)

	tmux := &test.TmuxSpy{}

	shell := shell.Shell{
		Git:  test.GitMock{},
		Tmux: tmux,
		Path: test.PathMock{},
	}

	err := run([]string{dir}, shell)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	session := ide.Name(dir)

	expectedCalls := [][]string{
		{"HasSession", session},
		{"New", session, dir, testEditor, dir},
		{"Attach", session},
	}

	if got, want := tmux.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestRunWithoutEditor(t *testing.T) {
	os.Unsetenv("EDITOR")
	os.Unsetenv("TMUX")

	tmux := &test.TmuxSpy{}
	dir := t.TempDir()

	shell := shell.Shell{
		Git:  test.GitMock{},
		Tmux: tmux,
		Path: test.PathMock{},
	}

	err := run([]string{dir}, shell)

	if got, want := err, ide.ErrEditorNotSet; !errors.Is(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}

	var expectedCalls [][]string
	if got, want := tmux.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestRunWithTmuxSessionExists(t *testing.T) {
	t.Setenv("EDITOR", testEditor)
	t.Setenv("TMUX", "test")

	dir := t.TempDir()
	session := ide.Name(dir)

	tmux := &test.TmuxSpy{
		Sessions: session,
	}

	shell := shell.Shell{
		Git:  test.GitMock{},
		Tmux: tmux,
		Path: test.PathMock{},
	}

	err := run([]string{dir}, shell)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	expectedCalls := [][]string{
		{"HasSession", session},
		{"Switch", session},
	}

	if got, want := tmux.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestRunInsideTmux(t *testing.T) {
	t.Setenv("EDITOR", testEditor)
	t.Setenv("TMUX", "test")

	dir := t.TempDir()
	tmux := &test.TmuxSpy{}

	shell := shell.Shell{
		Git:  test.GitMock{},
		Tmux: tmux,
		Path: test.PathMock{},
	}

	err := run([]string{dir}, shell)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	session := ide.Name(dir)

	expectedCalls := [][]string{
		{"HasSession", session},
		{"New", session, dir, testEditor, dir},
		{"Switch", session},
	}

	if got, want := tmux.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestRunWithoutTmux(t *testing.T) {
	t.Setenv("EDITOR", testEditor)
	os.Unsetenv("TMUX")

	tmuxSpy := &test.TmuxSpy{}
	dir := t.TempDir()

	shell := shell.Shell{
		Git:  test.GitMock{},
		Tmux: tmuxSpy,
		Path: test.PathMock{Missing: []string{"tmux"}},
	}

	err := run([]string{dir}, shell)

	if got, want := err, tmux.ErrTmuxNotInPath; !errors.Is(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
	var expectedCalls [][]string
	if got, want := tmuxSpy.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}
