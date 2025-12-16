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

	var workingDir string
	var command []string
	var err error
	switch len(args) {
	case 0:
		workingDir, err = os.Getwd()
	case 1:
		workingDir = args[0]
	default:
		workingDir = args[0]
		command = args[1:]
	}

	if err != nil {
		return err
	}

	project := project.Project{
		Name:       project.Name(workingDir),
		WorkingDir: workingDir,
	}

	return ide.Start(command, project, shell.Tmux, shell.Path)

}

func init() {
	rootCmd.AddCommand(openCmd)
}
