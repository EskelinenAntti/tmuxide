package cmd

import (
	"errors"
	"os"

	"path/filepath"

	"github.com/eskelinenantti/tmuxide/internal/git"
	"github.com/eskelinenantti/tmuxide/internal/ide"
	"github.com/eskelinenantti/tmuxide/internal/path"
	"github.com/eskelinenantti/tmuxide/internal/shell"
	"github.com/eskelinenantti/tmuxide/internal/tmux"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tmuxide",
	Short: "",
	Long:  ``,
	RunE:  runCmd,
}

func runCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	return run(args, shell.Shell{
		Git:  git.ShellGit{},
		Tmux: tmux.ShellTmux{},
		Path: path.ShellPath{},
	})
}

func run(args []string, shell shell.Shell) error {
	var target string
	var err error

	switch len(args) {
	case 0:
		target, err = os.Getwd()
	case 1:
		target, err = filepath.Abs(args[0])
	default:
		return errors.New("Invalid number of arguments.")
	}

	if err != nil {
		return err
	}

	if err := tmux.EnsureInstalled(shell.Path); err != nil {
		return err
	}

	project, err := ide.ProjectFor(target, shell)
	if err != nil {
		return err
	}

	return ide.Start(project, shell.Tmux)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
