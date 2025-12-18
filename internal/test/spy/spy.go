package spy

type Tmux struct {
	Calls   [][]string
	Session string
	Window  string
}

func (t *Tmux) HasSession(session string, window string) bool {
	args := []string{"HasSession", session, window}
	t.Calls = append(t.Calls, args)
	return t.Session == session && t.Window == window
}

func (t *Tmux) New(session string, dir string, cmd []string) error {
	args := []string{"New", session, dir}
	args = append(args, cmd...)
	t.Calls = append(t.Calls, args)
	return nil
}

func (t *Tmux) NewWindow(session string, window string, dir string, cmd []string) error {
	args := []string{"NewWindow", session, window, dir}
	args = append(args, cmd...)
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

func (t *Tmux) Kill(session string) error {
	args := []string{"Kill", session}
	t.Calls = append(t.Calls, args)
	return nil
}
