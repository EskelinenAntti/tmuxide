package cmd

import (
	"fmt"

	"github.com/eskelinenantti/tmuxide/internal/ide"
	"github.com/eskelinenantti/tmuxide/internal/project"
	"github.com/eskelinenantti/tmuxide/internal/shell"
	"github.com/spf13/cobra"
)

var helpOpen = `Opens tmux session for given directory, or creates one if it didn't exist.

Optionally, you can specify a command to be executed in the session.`

var openCmd = &cobra.Command{
	Use:   "open [directory] [command]",
	Short: "Open or create a tmux session for given directory.",
	Long:  helpOpen,
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return Open(args, ShellEnv{
			Git:  shell.Git{},
			Tmux: shell.SubCmdRunner{Command: "tmux"},
			Path: shell.Path{},
		})
	}}

func Open(args []string, shell ShellEnv) error {
	if len(args) == 0 {
		return ide.List(shell.Tmux, shell.Path)
	}

	workingDir := args[0]
	command := args[1:]

	project, err := project.ForDir(workingDir)
	if err != nil {
		return fmt.Errorf("could not open %s: %w", workingDir, err)
	}

	return ide.Start(command, project, shell.Tmux, shell.Path)
}

func init() {
	rootCmd.AddCommand(openCmd)
}
