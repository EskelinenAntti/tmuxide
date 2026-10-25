package shell

import (
	"errors"
	"fmt"

	"github.com/eskelinenantti/tmuxide/internal/shell/fd"
	"github.com/eskelinenantti/tmuxide/internal/shell/fzf"
	"github.com/eskelinenantti/tmuxide/internal/shell/git"
	"github.com/eskelinenantti/tmuxide/internal/shell/path"
	"github.com/eskelinenantti/tmuxide/internal/shell/runner"
)

var ErrCommandNotInstalled = errors.New("not installed")
var dependencies = []string{"fd", "fzf", "git"}

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
	Fd  fd.Cmd
	Fzf fzf.Cmd
	Git git.Cmd
}

// Runner provides both command capabilities needed to initialize the shell.
type Runner interface {
	runner.Runner
	runner.Starter
}

func Init(path path.ShellPath, runner Runner) (Shell, error) {
	for _, dependency := range dependencies {
		if err := assertInstalled(dependency, path); err != nil {
			return Shell{}, err
		}
	}

	return Shell{
		Fd:  fd.Cmd{Runner: runner},
		Fzf: fzf.Cmd{Starter: runner},
		Git: git.Cmd{Runner: runner},
	}, nil
}

func assertInstalled(command string, path path.ShellPath) error {
	if !path.Contains(command) {
		return NotInstalledError{Cmd: command}
	}
	return nil
}
