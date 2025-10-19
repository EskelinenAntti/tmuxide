package cmd

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/shell"
	"github.com/eskelinenantti/tmuxide/internal/tmux"
)

type PathMock struct{}

func (PathMock) Contains(program string) bool {
	return true
}

type GitMock struct {
	repository string
	err        error
}

func (git GitMock) RevParse(cwd string) (string, error) {
	return git.repository, git.err
}

type TmuxSpy struct {
	Calls []string
}

func (t *TmuxSpy) Attach(session string) error {
	t.Calls = append(t.Calls, "Attach")
	return nil
}

func (t *TmuxSpy) HasSession(name string) bool {
	t.Calls = append(t.Calls, "HasSession")
	return false
}

func (t *TmuxSpy) New(session string, dir string, cmd tmux.WindowCommand) error {
	t.Calls = append(t.Calls, "New")
	return nil
}

func (t *TmuxSpy) NewWindow(session string, dir string, cmd tmux.WindowCommand) error {
	t.Calls = append(t.Calls, "NewWindow")
	return nil
}

func (t *TmuxSpy) Switch(session string) error {
	t.Calls = append(t.Calls, "Switch")
	return nil
}

var tmuxSpy = TmuxSpy{Calls: []string{}}
var testEditor string = "editor"
var errNotInRepository = errors.New("not in repository")

func shellMock(repository string) shell.Shell {
	var err error
	if repository == "" {
		err = errNotInRepository
	}

	return shell.Shell{
		Git: GitMock{
			repository: repository,
			err:        err,
		},
		Tmux: &tmuxSpy,
		Path: PathMock{},
	}
}

func TestRun(t *testing.T) {
	t.Setenv("EDITOR", testEditor)
	os.Unsetenv("TMUX")

	var dir = t.TempDir()
	var shell = shellMock("")

	err := run([]string{dir}, shell)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	expectedCalls := []string{
		"HasSession", "New", "Attach",
	}

	if got, want := tmuxSpy.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}
