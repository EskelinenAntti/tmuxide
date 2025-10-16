package cmd

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/eskelinenantti/tmuxide/internal/git"
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

	editor, hasEditor := os.LookupEnv("EDITOR")
	if !hasEditor {
		cmd.PrintErr(
			"No editor was configured. Specify the editor you would like to use by setting the $EDITOR variable.\n\n" +
				"For example, to use Vim as your editor, add the following line to your ~/.zshrc:\n" +
				"export EDITOR=vim\n",
		)
		return errors.New("$EDITOR not set")
	}

	if _, err := exec.LookPath("tmux"); err != nil {
		cmd.PrintErr(
			"Did not find tmux, which is a required dependency for ide command.\n\n" +

				"You can install tmux e.g. via homebrew by running\n" +
				"brew install tmux\n",
		)
		return errors.New("tmux is not installed")
	}

	// TODO
	// Check if lazygit is installed
	// tests for tmux magic

	root, err := project.Root(target, git.RepositoryResolver{})
	if err != nil {
		return err
	}

	something := &tmux.AttachedSession{
		Session:    session.Name(target),
		WorkingDir: root,
		Command:    editor,
		Target:     target,
	}
	return tmux.Start(something)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
