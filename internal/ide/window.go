package ide

import (
	"errors"
	"fmt"

	"github.com/eskelinenantti/tmuxide/internal/input"
)

var ErrTmuxNotInstalled = errors.New(
	"Did not find tmux, which is a required dependency for ide command.\n\n" +

		"You can install tmux e.g. via homebrew by running\n" +
		"brew install tmux\n",
)

var ErrUnknownProgram = errors.New("Unknown program")

type Window []string

type ShellPath interface {
	Contains(path string) bool
}

func Windows(args input.Args, path ShellPath) ([]Window, error) {
	if !path.Contains("tmux") {
		return nil, ErrTmuxNotInstalled
	}

	windows, err := windows(args, path)
	if err != nil {
		return nil, err
	}
	return windows, nil
}

func windows(args input.Args, path ShellPath) ([]Window, error) {
	if args.Program == "" {
		return []Window{}, nil
	}

	if !path.Contains(args.Program) {
		return []Window{}, fmt.Errorf("%w: %s", ErrUnknownProgram, args.Program)
	}

	return []Window{
		{args.Program, args.Path},
	}, nil
}
