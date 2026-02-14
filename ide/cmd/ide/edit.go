package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/eskelinenantti/tmuxide/internal/ide"
	"github.com/eskelinenantti/tmuxide/internal/picker"
	"github.com/eskelinenantti/tmuxide/internal/project"
	"github.com/eskelinenantti/tmuxide/internal/shell"
	"github.com/eskelinenantti/tmuxide/internal/shell/tmux"
	"github.com/spf13/cobra"
)

var ErrEditorNotInstalled = errors.New("editor not installed")
var ErrEditorEnvNotSet = errors.New("editor not configured")

var helpEdit = `Open editor inside a tmux session.

The working directory and name of the session are deduced from the given path with the following heuristics:

1. If path is inside a Git repository, the working directory is the repository root.
2. If the path points to a file outside of repository, the working directory is the surrounding directory.
3. If the path points to a directory outside of repository, the working directory is the directory itself.

If a session for the working directory already exists, the editor will open in that session. Otherwise, a new session is created.`

var editCmd = &cobra.Command{
	Use:   "edit [path]",
	Short: "Open editor inside a tmux session.",
	Long:  helpEdit,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return Edit(args, ShellEnv{
			CmdRunner: shell.CmdRunner{},
			Path:      shell.Path{},
		})
	},
}

func Edit(args []string, shellEnv ShellEnv) error {
	tmux, err := tmux.InitTmux(shellEnv.Path, shellEnv.CmdRunner)
	if err != nil {
		return err
	}

	editorCmd := strings.Fields(os.Getenv("EDITOR"))

	if len(editorCmd) == 0 {
		return ErrEditorEnvNotSet
	}

	if !shellEnv.Path.Contains(editorCmd[0]) {
		return ErrEditorNotInstalled
	}

	var editArg string
	if len(args) == 0 {
		editArg, err = picker.Prompt(false, tmux, shell.FdCmd{Runner: shellEnv.CmdRunner}, shell.FzfCmd{Runner: shellEnv.CmdRunner})
	} else {
		editArg = args[0]
	}

	if editArg == "" || err != nil {
		return err
	}

	project, err := project.ForPath(editArg, shell.GitCmd{Runner: shellEnv.CmdRunner}, tmux)

	if err != nil {
		return fmt.Errorf("could not edit %s: %w", editArg, err)
	}

	var command []string
	if editArg != project.Name {
		command = append(editorCmd, editArg)
	}

	return ide.Start(command, project, tmux)
}

func init() {
	rootCmd.AddCommand(editCmd)
}
