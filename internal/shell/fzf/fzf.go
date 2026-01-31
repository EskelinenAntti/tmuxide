package fzf

import (
	"io"

	"github.com/eskelinenantti/tmuxide/internal/shell/tmux"
)

type Fzf struct {
	tmux.Runner
}

type Args struct{}

func (f Fzf) Execute(input io.Reader) ([]byte, error) {
	return f.Pipe("", Args{}, input)
}

func (a Args) Parse() []string {
	args := []string{
		"--reverse",
		"--height",
		"30%",
	}
	return args
}
