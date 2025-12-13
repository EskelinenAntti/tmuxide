package main

import (
	"fmt"
	"os"
	"slices"

	"github.com/eskelinenantti/tmuxide/internal/ide"
	"github.com/eskelinenantti/tmuxide/internal/input"
	"github.com/eskelinenantti/tmuxide/internal/project"
	"github.com/eskelinenantti/tmuxide/internal/shell"
)

const helpMsgTemplate string = `Usage: %s [path]

Arguments
	path (optional) - Path to project root directory or file.`

type shellEnv struct {
	Git  project.Git
	Tmux ide.Tmux
	Path ide.ShellPath
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

func run(args []string, shell shellEnv) error {
	if containHelp(args) {
		return fmt.Errorf(helpMsgTemplate, args[0])
	}

	parsedArgs, err := input.Parse(args)
	if err != nil {
		return err
	}

	project, err := project.New(parsedArgs, shell.Git)
	if err != nil {
		return err
	}

	return ide.Start(parsedArgs, project, shell.Tmux, shell.Path)
}

func containHelp(args []string) bool {
	return slices.Contains(args, "-h") || slices.Contains(args, "--help")
}
