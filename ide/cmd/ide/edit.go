package cmd

import (
	"os"

	"github.com/eskelinenantti/tmuxide/internal/ide"
	"github.com/eskelinenantti/tmuxide/internal/project"
	"github.com/eskelinenantti/tmuxide/internal/shell"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Open editor within a new tmux session",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return Edit(args, shell.ShellEnv{
			Git:  shell.Git{},
			Tmux: shell.Tmux{},
			Path: shell.Path{},
		})
	},
}

func Edit(args []string, shell shell.ShellEnv) error {
	var input = project.Input{Command: []string{os.Getenv("EDITOR")}}

	if len(args) == 1 {
		input.EditorPath = args[0]
		input.Command = append(input.Command, input.EditorPath)
	}

	project, err := project.New(input, shell.Git)
	if err != nil {
		return err
	}

	return ide.Start(input, project, shell.Tmux, shell.Path)
}
func init() {
	rootCmd.AddCommand(editCmd)
}
