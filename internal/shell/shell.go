package shell

import (
	"errors"
	"fmt"

	"github.com/eskelinenantti/tmuxide/internal/shell/fd"
	"github.com/eskelinenantti/tmuxide/internal/shell/fzf"
	"github.com/eskelinenantti/tmuxide/internal/shell/git"
	"github.com/eskelinenantti/tmuxide/internal/shell/path"
	"github.com/eskelinenantti/tmuxide/internal/shell/runner"
	"github.com/eskelinenantti/tmuxide/internal/shell/tmux"
)

var ErrCommandNotInstalled = errors.New("not installed")

type NotInstalledError struct {
	Cmd string
}

func (e NotInstalledError) Error() string {
	return fmt.Sprintf("%s %v", e.Cmd, ErrCommandNotInstalled)
}

func (e NotInstalledError) Unwrap() error {
	return ErrCommandNotInstalled
}

type Shell struct {
	Tmux tmux.Cmd
	Fd   fd.Cmd
	Fzf  fzf.Cmd
	Git  git.Cmd
}

func Init(path path.ShellPath, runner runner.Runner) (Shell, error) {
	if err := assertInstalled("tmux", path); err != nil {
		return Shell{}, err
	}
	if err := assertInstalled("fd", path); err != nil {
		return Shell{}, err
	}
	if err := assertInstalled("fzf", path); err != nil {
		return Shell{}, err
	}
	if err := assertInstalled("git", path); err != nil {
		return Shell{}, err
	}

	return Shell{
		Tmux: tmux.Cmd{Runner: runner},
		Fd:   fd.Cmd{Runner: runner},
		Fzf:  fzf.Cmd{Runner: runner},
		Git:  git.Cmd{Runner: runner},
	}, nil
}

func assertInstalled(command string, path path.ShellPath) error {
	if !path.Contains(command) {
		return NotInstalledError{Cmd: command}
	}
	return nil
}
