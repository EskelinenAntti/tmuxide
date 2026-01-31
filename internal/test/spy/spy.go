package spy

import (
	"errors"
	"slices"

	"github.com/eskelinenantti/tmuxide/internal/shell/tmux"
)

type Call struct {
	Name string
	Args tmux.Parser
}

type Tmux struct {
	Calls  []Call
	Errors []string
}

func (t *Tmux) Run(name string, args tmux.Parser) error {
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

func (t *Tmux) Attach(name string, args tmux.Parser) error {
	return t.Run(name, args)
}

func (t *Tmux) Output(name string, args tmux.Parser) ([]byte, error) {
	return t.Output(name, args)
}
