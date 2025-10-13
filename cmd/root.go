package cmd

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/eskelinenantti/tmuxide/internal/project"
	"github.com/eskelinenantti/tmuxide/internal/session"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tmuxide",
	Short: "",
	Long:  ``,
	RunE:  run,
}

func run(cmd *cobra.Command, args []string) error {
	var selected string
	var err error
	switch len(args) {
	case 0:
		selected, err = os.Getwd()
	case 1:
		selected, err = filepath.Abs(args[0])
	default:
		// We should never end up here, but handle the error nicely nevertheless.
		cmd.Help()
		return errors.New("Invalid number of arguments.")
	}

	if err != nil {
		return err
	}

	root, err := project.Root(selected)
	if err != nil {
		return err
	}

	tmuxCmd := exec.Command("tmux", "new", "-c", root, "-s", session.Name(selected))
	return attachAndRun(tmuxCmd)
}

func attachAndRun(cmd *exec.Cmd) error {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
