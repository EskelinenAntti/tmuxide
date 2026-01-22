package cmd

import (
	"github.com/eskelinenantti/tmuxide/internal/ide"
	"github.com/eskelinenantti/tmuxide/internal/shell"
	"github.com/spf13/cobra"
)

var helpQuit = "Close session"

var quitCmd = &cobra.Command{
	Use:   "quit",
	Short: helpQuit,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return Quit(ShellEnv{
			Tmux: shell.SubCmdRunner{Command: "tmux"},
		})
	},
}

func Quit(shellEnv ShellEnv) error {
	return ide.Quit(shellEnv.Tmux)
}

func init() {
	rootCmd.AddCommand(quitCmd)
}
