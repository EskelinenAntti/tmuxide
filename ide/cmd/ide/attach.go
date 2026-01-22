package cmd

import (
	"github.com/eskelinenantti/tmuxide/internal/ide"
	"github.com/eskelinenantti/tmuxide/internal/shell"
	"github.com/spf13/cobra"
)

var helpAttach = "Alias to `tmux attach`"

var attachCmd = &cobra.Command{
	Use:   "attach",
	Short: "Attach to previously open session",
	Long:  helpAttach,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return Attach(ShellEnv{
			Tmux: shell.SubCmdRunner{Command: "tmux"},
		})
	}}

func Attach(shell ShellEnv) error {
	return ide.Attach(shell.Tmux)
}

func init() {
	rootCmd.AddCommand(attachCmd)
}
