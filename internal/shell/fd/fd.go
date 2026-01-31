package fd

import (
	"os"

	"github.com/eskelinenantti/tmuxide/internal/shell/tmux"
)

type Fd struct {
	tmux.Runner
}

type Args struct{}

func (f Fd) Execute() ([]byte, error) {
	return f.Output("", Args{})
}

func (a Args) Parse() []string {
	args := []string{
		"--follow",
		"--hidden",
		"--exclude", "{.git,node_modules,target,build,Library}",
		".",
		os.Getenv("HOME"),
	}
	return args
}
