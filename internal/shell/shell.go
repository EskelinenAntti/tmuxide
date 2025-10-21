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
	return Shell{
		Git:  git.ShellGit{},
		Tmux: tmux.ShellTmux{},
		Path: path.ShellPath{},
	}, nil
}
