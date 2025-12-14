package cmd

import (
	"os"

	"github.com/eskelinenantti/tmuxide/internal/shell"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ide",
	Short: "Turn tmux and your favourite editor into an ide",
	RunE: func(cmd *cobra.Command, args []string) error {
		return Open(args, shell.ShellEnv{
			Git:  shell.Git{},
			Tmux: shell.Tmux{},
			Path: shell.Path{},
		})
	}}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
