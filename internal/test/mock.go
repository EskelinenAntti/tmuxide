package test

import (
	"errors"
	"slices"

	"github.com/eskelinenantti/tmuxide/internal/tmux"
)

type PathMock struct {
	Missing []string
}

func (pathMock PathMock) Contains(program string) bool {
	return !slices.Contains(pathMock.Missing, program)
}

type GitMock struct {
	Repository string
}

func (git GitMock) RevParse(cwd string) (string, error) {

	if git.Repository == "" {
		return "", errors.New("not inside git repo")
	}

	return git.Repository, nil
}

type TmuxSpy struct {
	Calls    [][]string
	Sessions string
}

func (t *TmuxSpy) Attach(session string) error {
	args := []string{"Attach", session}
	t.Calls = append(t.Calls, args)
	return nil
}

func (t *TmuxSpy) HasSession(name string) bool {
	args := []string{"HasSession", name}
	t.Calls = append(t.Calls, args)
	return t.Sessions == name
}

func (t *TmuxSpy) New(session string, dir string, cmd tmux.WindowCommand) error {
	args := []string{"New", session, dir}
	args = append(args, cmd.Cmd)
	args = append(args, cmd.Args...)
	t.Calls = append(t.Calls, args)
	return nil
}

func (t *TmuxSpy) NewWindow(session string, dir string, cmd tmux.WindowCommand) error {
	args := []string{"NewWindow", session, dir}
	args = append(args, cmd.Cmd)
	args = append(args, cmd.Args...)
	t.Calls = append(t.Calls, args)
	return nil
}

func (t *TmuxSpy) Switch(session string) error {
	args := []string{"Switch", session}
	t.Calls = append(t.Calls, args)
	return nil
}
