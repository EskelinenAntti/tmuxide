package cmd

import (
	"fmt"

	"github.com/eskelinenantti/tmuxide/internal/ide"
	"github.com/eskelinenantti/tmuxide/internal/picker"
	"github.com/eskelinenantti/tmuxide/internal/project"
	"github.com/eskelinenantti/tmuxide/internal/shell"
	"github.com/eskelinenantti/tmuxide/internal/shell/path"
	"github.com/eskelinenantti/tmuxide/internal/shell/runner"
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
		return Open(args, runner.CmdRunner{}, path.Path{})
	}}

func Open(args []string, runner runner.Runner, path path.ShellPath) error {

	shell, err := shell.Init(path, runner)
	if err != nil {
		return err
	}

	var workingDir string
	var command []string
	if len(args) == 0 {
		workingDir, err = picker.Prompt(true, shell.Tmux, shell.Fd, shell.Fzf)
	} else {
		workingDir = args[0]
		command = args[1:]
	}

	if err != nil || workingDir == "" {
		return err
	}

	project, err := project.ForDir(workingDir, shell.Tmux)
	if err != nil {
		return fmt.Errorf("could not open %s: %w", workingDir, err)
	}

	return ide.Start(command, project, shell.Tmux)
}

func init() {
	rootCmd.AddCommand(openCmd)
}
