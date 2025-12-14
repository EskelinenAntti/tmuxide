package shell

import "github.com/eskelinenantti/tmuxide/internal/project"
import "github.com/eskelinenantti/tmuxide/internal/ide"

type ShellEnv struct {
	Git  project.Git
	Tmux ide.Tmux
	Path ide.ShellPath
}
