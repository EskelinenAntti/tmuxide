package cmd

import (
	"errors"
	"os"

	"path/filepath"

	"github.com/eskelinenantti/tmuxide/internal/git"
	"github.com/eskelinenantti/tmuxide/internal/ide"
	"github.com/eskelinenantti/tmuxide/internal/tmux"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tmuxide",
	Short: "",
	Long:  ``,
	RunE:  run,
}

func run(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	var target string
	var err error

	switch len(args) {
	case 0:
		target, err = os.Getwd()
	case 1:
		target, err = filepath.Abs(args[0])
	default:
		// We should never end up here, but handle the error nicely nevertheless.
		return errors.New("Invalid number of arguments.")
	}

	if err != nil {
		cmd.PrintErr(err)
	}

	repository, err := ide.Repository(target, git.ShellGit{})
	if err != nil {
		return err
	}

	project, err := ide.ProjectFor(target, repository)
	if err != nil {
		return nil
	}

	windows, err := ide.WindowsFor(target, repository)
	if err != nil {
		return err
	}

	tmux, err := tmux.Command()
	if err != nil {
		return err
	}

	return ide.Start(project, windows, tmux)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
