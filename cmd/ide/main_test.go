package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/ide"
	"github.com/eskelinenantti/tmuxide/internal/project"
	"github.com/eskelinenantti/tmuxide/internal/test/mock"
	"github.com/eskelinenantti/tmuxide/internal/test/spy"
)

const command string = "ide"
const testEditor string = "editor"

func TestRunWithFile(t *testing.T) {
	t.Setenv("EDITOR", testEditor)
	os.Unsetenv("TMUX")

	dir := t.TempDir()
	file := dir + "/file.txt"
	os.WriteFile(file, []byte{}, 0644)

	tmux := &spy.Tmux{}

	shell := shellEnv{
		Git:  mock.Git{},
		Tmux: tmux,
		Path: mock.Path{},
	}

	err := run([]string{command, file}, shell)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	session := project.Name(dir)

	expectedCalls := [][]string{
		{"HasSession", session},
		{"New", session, dir, testEditor, file},
		{"Attach", session},
	}

	if got, want := tmux.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestRunWithDirectory(t *testing.T) {
	t.Setenv("EDITOR", testEditor)
	os.Unsetenv("TMUX")

	dir := t.TempDir()
	tmux := &spy.Tmux{}

	shell := shellEnv{
		Git:  mock.Git{},
		Tmux: tmux,
		Path: mock.Path{},
	}

	err := run([]string{command, dir}, shell)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	session := project.Name(dir)

	expectedCalls := [][]string{
		{"HasSession", session},
		{"New", session, dir, testEditor, dir},
		{"Attach", session},
	}

	if got, want := tmux.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestRunWithFileInRepository(t *testing.T) {
	t.Setenv("EDITOR", testEditor)
	os.Unsetenv("TMUX")

	repository := t.TempDir()
	file := repository + "/file.txt"
	os.WriteFile(file, []byte{}, 0644)

	tmux := &spy.Tmux{}

	shell := shellEnv{
		Git:  mock.Git{Repository: repository},
		Tmux: tmux,
		Path: mock.Path{},
	}

	err := run([]string{command, file}, shell)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	session := project.Name(repository)

	expectedCalls := [][]string{
		{"HasSession", session},
		{"New", session, repository, testEditor, file},
		{"Attach", session},
	}

	if got, want := tmux.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestRunWithDirectoryInRepository(t *testing.T) {
	t.Setenv("EDITOR", testEditor)
	os.Unsetenv("TMUX")

	repository := t.TempDir()
	dir := filepath.Join(repository, "path/to/dir/in/repository")

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		t.Fatalf("err=%v", err)
	}

	tmux := &spy.Tmux{}

	shell := shellEnv{
		Git:  mock.Git{Repository: repository},
		Tmux: tmux,
		Path: mock.Path{},
	}

	err := run([]string{command, dir}, shell)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	session := project.Name(repository)

	expectedCalls := [][]string{
		{"HasSession", session},
		{"New", session, repository, testEditor, dir},
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

	tmux := &spy.Tmux{}

	shell := shellEnv{
		Git:  mock.Git{},
		Tmux: tmux,
		Path: mock.Path{},
	}

	err := run([]string{command}, shell)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	session := project.Name(dir)

	expectedCalls := [][]string{
		{"HasSession", session},
		{"New", session, dir, testEditor, dir},
		{"Attach", session},
	}

	if got, want := tmux.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestRunHelp(t *testing.T) {
	os.Unsetenv("EDITOR")
	os.Unsetenv("TMUX")

	tmux := &spy.Tmux{}
	dir := t.TempDir()

	shell := shellEnv{
		Git:  mock.Git{},
		Tmux: tmux,
		Path: mock.Path{},
	}

	err := run([]string{command, dir, "-h"}, shell)

	if got, want := err.Error(), fmt.Sprintf(helpMsgTemplate, command); got != want {
		t.Fatalf("got=%v, want=%v", got, want)
	}

	var expectedCalls [][]string
	if got, want := tmux.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestRunWithoutEditor(t *testing.T) {
	os.Unsetenv("EDITOR")
	os.Unsetenv("TMUX")

	tmux := &spy.Tmux{}
	dir := t.TempDir()

	shell := shellEnv{
		Git:  mock.Git{},
		Tmux: tmux,
		Path: mock.Path{},
	}

	err := run([]string{command, dir}, shell)

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
	session := project.Name(dir)

	tmux := &spy.Tmux{
		Sessions: session,
	}

	shell := shellEnv{
		Git:  mock.Git{},
		Tmux: tmux,
		Path: mock.Path{},
	}

	err := run([]string{command, dir}, shell)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	expectedCalls := [][]string{
		{"HasSession", session},
		{"Kill", session},
		{"New", session, dir, testEditor, dir},
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
	tmux := &spy.Tmux{}

	shell := shellEnv{
		Git:  mock.Git{},
		Tmux: tmux,
		Path: mock.Path{},
	}

	err := run([]string{command, dir}, shell)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	session := project.Name(dir)

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

	tmuxSpy := &spy.Tmux{}
	dir := t.TempDir()

	shell := shellEnv{
		Git:  mock.Git{},
		Tmux: tmuxSpy,
		Path: mock.Path{Missing: []string{"tmux"}},
	}

	err := run([]string{command, dir}, shell)

	if got, want := err, ide.ErrTmuxNotInstalled; !errors.Is(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
	var expectedCalls [][]string
	if got, want := tmuxSpy.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}
