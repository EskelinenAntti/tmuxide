package spy

import (
	"errors"
	"github.com/eskelinenantti/tmuxide/internal/shell"
	"slices"
)

type Call struct {
	Name string
	Args shell.Parser
}

type Tmux struct {
	Calls  []Call
	Errors []string
}

func (t *Tmux) Run(name string, args shell.Parser) error {
	call := Call{Name: name, Args: args}
	t.Calls = append(t.Calls, call)

	for i, error := range t.Errors {
		if name == error {
			t.Errors = slices.Delete(t.Errors, i, i+1)
			return errors.New("error")
		}
	}
	return nil
}

func (t *Tmux) Attach(name string, args shell.Parser) error {
	return t.Run(name, args)
}
