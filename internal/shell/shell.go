package shell

import (
	"github.com/eskelinenantti/tmuxide/internal/git"
	"github.com/eskelinenantti/tmuxide/internal/path"
	"github.com/eskelinenantti/tmuxide/internal/tmux"
)

type Shell struct {
	Git  git.Command
	Tmux tmux.Command
	Path path.Path
}

func Get() (Shell, error) {
	tmux, err := tmux.Get()
	if err != nil {
		return Shell{}, err
	}

	git := git.ShellGit{}
	path := path.ShellPath{}

	return Shell{
		Git:  git,
		Tmux: tmux,
		Path: path,
	}, nil
}
