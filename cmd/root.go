package cmd

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/eskelinenantti/tmuxide/internal/git"
	"github.com/eskelinenantti/tmuxide/internal/ide"
	"github.com/eskelinenantti/tmuxide/internal/project"
	"github.com/eskelinenantti/tmuxide/internal/session"
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

	if _, err := exec.LookPath("tmux"); err != nil {
		cmd.PrintErr(
			"Did not find tmux, which is a required dependency for ide command.\n\n" +

				"You can install tmux e.g. via homebrew by running\n" +
				"brew install tmux\n",
		)
		return errors.New("tmux is not installed")
	}

	root, err := project.Root(target, git.RepositoryResolver{})
	if err != nil {
		return err
	}

	return ide.Start(target, &tmux.Session{
		Session:    session.Name(target),
		WorkingDir: root,
	})
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
