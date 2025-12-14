package cmd

import (
	"os"

	"github.com/eskelinenantti/tmuxide/internal/ide"
	"github.com/eskelinenantti/tmuxide/internal/project"
	"github.com/eskelinenantti/tmuxide/internal/shell"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Open new session within a folder",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return Open(args, shell.ShellEnv{
			Git:  shell.Git{},
			Tmux: shell.Tmux{},
			Path: shell.Path{},
		})
	}}

func Open(args []string, shell shell.ShellEnv) error {
	var input = project.Input{}

	var err error
	switch len(args) {
	case 0:
		input.WorkingDir, err = os.Getwd()
	case 1:
		input.WorkingDir = args[0]
	default:
		input.WorkingDir = args[0]
		input.Command = args[1:]
	}

	if err != nil {
		return err
	}

	project, err := project.New(input, shell.Git)
	if err != nil {
		return err
	}

	return ide.Start(input, project, shell.Tmux, shell.Path)

}

func init() {
	rootCmd.AddCommand(openCmd)
}
