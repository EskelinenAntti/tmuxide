package spy

import (
	"github.com/eskelinenantti/tmuxide/internal/ide"
)

type Tmux struct {
	Calls    [][]string
	Sessions string
}

func (t *Tmux) HasSession(name string) bool {
	args := []string{"HasSession", name}
	t.Calls = append(t.Calls, args)
	return t.Sessions == name
}

func (t *Tmux) New(session string, dir string, window ide.Window) error {
	args := []string{"New", session, dir}
	args = append(args, window...)
	t.Calls = append(t.Calls, args)
	return nil
}

func (t *Tmux) NewWindow(session string, dir string, window ide.Window) error {
	args := []string{"NewWindow", session, dir}
	args = append(args, window...)
	t.Calls = append(t.Calls, args)
	return nil
}

func (t *Tmux) Attach(session string) error {
	args := []string{"Attach", session}
	t.Calls = append(t.Calls, args)
	return nil
}

func (t *Tmux) Switch(session string) error {
	args := []string{"Switch", session}
	t.Calls = append(t.Calls, args)
	return nil
}
