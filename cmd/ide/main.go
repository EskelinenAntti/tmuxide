package main

import (
	"fmt"
	"os"

	"github.com/eskelinenantti/tmuxide/internal/ide"
	"github.com/eskelinenantti/tmuxide/internal/ide/window"
	"github.com/eskelinenantti/tmuxide/internal/input"
	"github.com/eskelinenantti/tmuxide/internal/project"
	"github.com/eskelinenantti/tmuxide/internal/shell"
)

type shellEnv struct {
	Git  project.Git
	Tmux ide.Tmux
	Path window.Path
}

func main() {
	err := run(os.Args, shellEnv{
		Git:  shell.Git{},
		Tmux: shell.Tmux{},
		Path: shell.Path{},
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(args input.Args, shell shellEnv) error {
	if args.ContainHelp() {
		return fmt.Errorf(helpMsgTemplate, args.Command())
	}

	target, err := input.Path(args)
	if err != nil {
		return err
	}

	project, err := project.New(target, shell.Git)
	if err != nil {
		return err
	}

	return ide.Start(project, shell.Tmux, shell.Path)
}

const helpMsgTemplate string = `Usage: %s [path]

Arguments
	path (optional) - Path to project root directory or file.`
