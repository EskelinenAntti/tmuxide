package cmd

import (
	"fmt"
	"os"
	"slices"

	"path/filepath"

	"github.com/eskelinenantti/tmuxide/internal/git"
	"github.com/eskelinenantti/tmuxide/internal/ide"
	"github.com/eskelinenantti/tmuxide/internal/path"
	"github.com/eskelinenantti/tmuxide/internal/shell"
	"github.com/eskelinenantti/tmuxide/internal/tmux"
)

func Execute() {
	err := run(os.Args, shell.Shell{
		Git:  git.ShellGit{},
		Tmux: tmux.ShellTmux{},
		Path: path.ShellPath{},
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(args []string, shell shell.Shell) error {

	if err := help(args); err != nil {
		return err
	}

	var target string
	var err error

	switch len(args) {
	case 1:
		target, err = os.Getwd()
	case 2:
		target, err = filepath.Abs(args[1])
	default:
		return fmt.Errorf("Invalid number of arguments. See %s --help.", args[0])
	}

	if err != nil {
		return err
	}

	if err := tmux.EnsureInstalled(shell.Path); err != nil {
		return err
	}

	project, err := ide.ProjectFor(target, shell)
	if err != nil {
		return err
	}

	return ide.Start(project, shell.Tmux)
}

const helpMsgTemplate string = `Usage: %s [path]

Arguments
	path (optional) - Path to project root directory or file.`

func help(args []string) error {
	if slices.Contains(args, "-h") || slices.Contains(args, "--help") {
		return fmt.Errorf(helpMsgTemplate, args[0])
	}
	return nil
}
