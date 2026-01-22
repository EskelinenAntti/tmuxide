package cmd

import (
	"github.com/eskelinenantti/tmuxide/internal/ide"
	"github.com/eskelinenantti/tmuxide/internal/shell"
	"github.com/spf13/cobra"
)

var helpList = "Open session selector"

var listCmd = &cobra.Command{
	Use:   "list",
	Short: helpList,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return List(ShellEnv{
			Tmux: shell.SubCmdRunner{Command: "tmux"},
		})
	},
}

func List(shellEnv ShellEnv) error {
	return ide.List(shellEnv.Tmux)
}

func init() {
	rootCmd.AddCommand(listCmd)
}
