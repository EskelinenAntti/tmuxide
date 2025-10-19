package cmd

import (
	"errors"
	"os"

	"path/filepath"

	"github.com/eskelinenantti/tmuxide/internal/ide"
	"github.com/eskelinenantti/tmuxide/internal/shell"
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

	shell, err := shell.Get()
	if err != nil {
		return err
	}

	return run(args, shell)
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

	project, err := ide.ProjectFor(target, shell)
	if err != nil {
		return nil
	}

	return ide.Start(project, shell.Tmux)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
