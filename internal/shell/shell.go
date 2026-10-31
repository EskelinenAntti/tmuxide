package shell

import (
	"errors"
	"fmt"
	"strings"

	"github.com/eskelinenantti/zipln/internal/shell/fd"
	"github.com/eskelinenantti/zipln/internal/shell/fzf"
	"github.com/eskelinenantti/zipln/internal/shell/git"
	"github.com/eskelinenantti/zipln/internal/shell/path"
	"github.com/eskelinenantti/zipln/internal/shell/runner"
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

type MissingToolsError struct {
	Tools []string
}

func (e MissingToolsError) Error() string {
	return "couldn't find " + strings.Join(e.Tools, ", ")
}

func (e MissingToolsError) Unwrap() []error {
	errs := make([]error, len(e.Tools))
	for i, command := range e.Tools {
		errs[i] = NotInstalledError{Cmd: command}
	}
	return errs
}

type Shell struct {
	Fd  fd.Cmd
	Fzf fzf.Cmd
	Git git.Cmd
}

func Init(path path.ShellPath, runner runner.Runner) (Shell, error) {
	var missing []string
	for _, dependency := range dependencies {
		if !path.Contains(dependency) {
			missing = append(missing, dependency)
		}
	}
	if len(missing) > 0 {
		return Shell{}, MissingToolsError{Tools: missing}
	}

	return Shell{
		Fd:  fd.Cmd{Runner: runner},
		Fzf: fzf.Cmd{Runner: runner},
		Git: git.Cmd{Runner: runner},
	}, nil
}
