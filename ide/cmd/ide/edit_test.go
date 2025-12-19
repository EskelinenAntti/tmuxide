package cmd

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/project"
	"github.com/eskelinenantti/tmuxide/internal/shell"
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

	tmux := &spy.Tmux{}

	shell := shell.ShellEnv{
		Git:  mock.Git{},
		Tmux: tmux,
		Path: mock.Path{},
	}

	err := Edit([]string{}, shell)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	session := project.Name(dir)

	expectedCalls := [][]string{
		{"HasSession", session, editor},
		{"HasSession", session, ""},
		{"New", session, dir, editor, dir},
		{"Attach", session},
	}

	if got, want := tmux.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}
func TestEditFile(t *testing.T) {
	os.Unsetenv("TMUX")
	t.Setenv("EDITOR", editor)

	dir := t.TempDir()
	file := dir + "/file.txt"
	os.WriteFile(file, []byte{}, 0644)

	tmux := &spy.Tmux{}

	shell := shell.ShellEnv{
		Git:  mock.Git{},
		Tmux: tmux,
		Path: mock.Path{},
	}

	err := Edit([]string{file}, shell)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	session := project.Name(dir)

	expectedCalls := [][]string{
		{"HasSession", session, editor},
		{"HasSession", session, ""},
		{"New", session, dir, editor, file},
		{"Attach", session},
	}

	if got, want := tmux.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestEditNonExistingFile(t *testing.T) {
	os.Unsetenv("TMUX")
	t.Setenv("EDITOR", editor)

	dir := t.TempDir()
	file := dir + "/file.txt"

	tmux := &spy.Tmux{}

	shell := shell.ShellEnv{
		Git:  mock.Git{},
		Tmux: tmux,
		Path: mock.Path{},
	}

	err := Edit([]string{file}, shell)

	if got, want := err, project.ErrInvalidPath; !errors.Is(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
	var expectedCalls [][]string
	if got, want := tmux.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestEditDirectory(t *testing.T) {
	os.Unsetenv("TMUX")
	t.Setenv("EDITOR", editor)

	dir := t.TempDir()
	tmux := &spy.Tmux{}

	shell := shell.ShellEnv{
		Git:  mock.Git{},
		Tmux: tmux,
		Path: mock.Path{},
	}

	err := Edit([]string{dir}, shell)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	session := project.Name(dir)

	expectedCalls := [][]string{
		{"HasSession", session, editor},
		{"HasSession", session, ""},
		{"New", session, dir, editor, dir},
		{"Attach", session},
	}

	if got, want := tmux.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestEditFileInRepository(t *testing.T) {
	os.Unsetenv("TMUX")
	t.Setenv("EDITOR", editor)

	repository := t.TempDir()
	file := repository + "/file.txt"
	os.WriteFile(file, []byte{}, 0644)

	tmux := &spy.Tmux{}

	shell := shell.ShellEnv{
		Git:  mock.Git{Repository: repository},
		Tmux: tmux,
		Path: mock.Path{},
	}

	err := Edit([]string{file}, shell)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	session := project.Name(repository)

	expectedCalls := [][]string{
		{"HasSession", session, editor},
		{"HasSession", session, ""},
		{"New", session, repository, editor, file},
		{"Attach", session},
	}

	if got, want := tmux.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestEditWithExistingSession(t *testing.T) {
	t.Setenv("TMUX", "test")
	t.Setenv("EDITOR", editor)

	dir := t.TempDir()
	session := project.Name(dir)

	tmux := &spy.Tmux{
		Session: session,
		Window:  "",
	}

	shell := shell.ShellEnv{
		Git:  mock.Git{},
		Tmux: tmux,
		Path: mock.Path{},
	}

	err := Edit([]string{dir}, shell)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	expectedCalls := [][]string{
		{"HasSession", session, editor},
		{"HasSession", session, ""},
		{"NewWindow", session, "", dir, editor, editor, dir},
		{"Switch", session},
	}

	if got, want := tmux.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestEditWithExistingWindow(t *testing.T) {
	t.Setenv("TMUX", "test")
	t.Setenv("EDITOR", editor)

	dir := t.TempDir()
	session := project.Name(dir)

	tmux := &spy.Tmux{
		Session: session,
		Window:  editor,
	}

	shell := shell.ShellEnv{
		Git:  mock.Git{},
		Tmux: tmux,
		Path: mock.Path{},
	}

	err := Edit([]string{dir}, shell)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	expectedCalls := [][]string{
		{"HasSession", session, editor},
		{"NewWindow", session, editor, dir, editor, editor, dir},
		{"Switch", session},
	}

	if got, want := tmux.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestEditWithUnsetEditor(t *testing.T) {
	os.Unsetenv("EDITOR")
	os.Unsetenv("TMUX")

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

	err := Edit([]string{dir}, shell)

	if got, want := err, ErrEditorEnvNotSet; !errors.Is(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
	var expectedCalls [][]string
	if got, want := tmux.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestEditWithEditorNotInstalled(t *testing.T) {
	t.Setenv("EDITOR", editor)
	os.Unsetenv("TMUX")

	dir := t.TempDir()
	session := project.Name(dir)

	tmux := &spy.Tmux{
		Session: session,
	}

	shell := shell.ShellEnv{
		Git:  mock.Git{},
		Tmux: tmux,
		Path: mock.Path{Missing: []string{editor}},
	}

	err := Edit([]string{dir}, shell)

	if got, want := err, ErrEditorNotInstalled; !errors.Is(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
	var expectedCalls [][]string
	if got, want := tmux.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}
