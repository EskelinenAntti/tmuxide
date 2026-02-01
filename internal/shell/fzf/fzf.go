package fzf

import (
	"io"

	"github.com/eskelinenantti/tmuxide/internal/shell"
)

type Fzf struct {
	shell.Runner
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
