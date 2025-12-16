package cmd

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/eskelinenantti/tmuxide/internal/ide"
	"github.com/eskelinenantti/tmuxide/internal/project"
	"github.com/eskelinenantti/tmuxide/internal/shell"
	"github.com/spf13/cobra"
)

var ErrEditorNotInstalled = errors.New("Editor not installed")
var ErrEditorNotSet = errors.New(
	"No editor was configured. Specify the editor you would like to use by setting the $EDITOR variable.\n\n" +
		"For example, to use Vim as your editor, add the following line to your ~/.zshrc:\n" +
		"export EDITOR=vim\n",
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Open editor within a new tmux session",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return Edit(args, shell.ShellEnv{
			Git:  shell.Git{},
			Tmux: shell.Tmux{},
			Path: shell.Path{},
		})
	},
}

func Edit(args []string, shell shell.ShellEnv) error {
	var editor = os.Getenv("EDITOR")

	if editor == "" {
		return ErrEditorNotSet
	}

	if !shell.Path.Contains(editor) {
		return ErrEditorNotInstalled
	}

	var editorPath string
	var err error
	switch len(args) {
	case 0:
		editorPath, err = os.Getwd()
	case 1:
		editorPath, err = filepath.Abs(args[0])
	}
	command := []string{editor, editorPath}

	project, err := project.New(editorPath, shell.Git)
	if err != nil {
		return err
	}

	return ide.Start(command, project, shell.Tmux, shell.Path)
}

func init() {
	rootCmd.AddCommand(editCmd)
}
