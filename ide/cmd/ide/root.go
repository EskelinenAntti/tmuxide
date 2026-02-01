package cmd

import (
	"errors"
	"os"

	"github.com/eskelinenantti/tmuxide/internal/project"
	"github.com/eskelinenantti/tmuxide/internal/shell"
	"github.com/eskelinenantti/tmuxide/internal/shell/tmux"
	"github.com/spf13/cobra"
)

type ShellEnv struct {
	Git        project.Git
	TmuxRunner shell.Runner
	Path       shell.ShellPath
	FdRunner   shell.Runner
	FzfRunner  shell.Runner
}

var rootCmd = &cobra.Command{
	Use:   "ide",
	Short: "Turn tmux and your favourite editor into an IDE with tmuxide.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return Open(args, ShellEnv{
			Git:        shell.Git{},
			TmuxRunner: shell.CmdRunner{Command: "tmux"},
			Path:       shell.Path{},
			FdRunner:   shell.CmdRunner{Command: "fd"},
			FzfRunner:  shell.CmdRunner{Command: "fzf"},
		})
	}}

var helpNoEditorConfigured = `
No editor was configured. Specify the editor you would like to use by setting the $EDITOR variable.
For example, to use Vim as your editor, add the following line to your ~/.zshrc or ~/.bashrc:
				
export EDITOR=vim`

var helpTmuxNotInstalled = `
Did not find tmux, which is a required dependency for ide command.

You can install tmux e.g. via homebrew by running:
brew install tmux`

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

	if errors.Is(err, tmux.ErrTmuxNotInstalled) {
		rootCmd.PrintErrln(helpTmuxNotInstalled)
	}

	os.Exit(1)
}

func init() {
}
