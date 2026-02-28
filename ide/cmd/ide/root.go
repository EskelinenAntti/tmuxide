package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/eskelinenantti/tmuxide/internal/shell"
	"github.com/eskelinenantti/tmuxide/internal/shell/path"
	"github.com/eskelinenantti/tmuxide/internal/shell/runner"
	"github.com/spf13/cobra"
)

type ShellEnv struct {
	Path      path.ShellPath
	CmdRunner runner.Runner
}

var rootCmd = &cobra.Command{
	Use:   "ide",
	Short: "Turn tmux and your favourite editor into an IDE with tmuxide.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return Open(args, runner.CmdRunner{}, path.Path{})
	}}

var helpNoEditorConfigured = `
No editor was configured. Specify the editor you would like to use by setting the $EDITOR variable.
For example, to use Vim as your editor, add the following line to your ~/.zshrc or ~/.bashrc:
				
export EDITOR=vim`

var helpCommandNotInstalledTemplate = `
Did not find %s, which is a required dependency for ide command.

You can install %[1]s e.g. via homebrew by running:
brew install %[1]s`

func Execute() {
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = false
	err := rootCmd.Execute()

	if err == nil {
		return
	}

	if errors.Is(err, ErrEditorEnvNotSet) {
		rootCmd.PrintErrln(helpNoEditorConfigured)
	}

	var commandNotInstalled shell.NotInstalledError
	if errors.As(err, &commandNotInstalled) {
		rootCmd.PrintErrln(fmt.Sprintf(helpCommandNotInstalledTemplate, commandNotInstalled.Cmd))
	}

	os.Exit(1)
}

func init() {
}
