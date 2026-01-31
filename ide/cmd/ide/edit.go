package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/eskelinenantti/tmuxide/internal/ide"
	"github.com/eskelinenantti/tmuxide/internal/project"
	"github.com/eskelinenantti/tmuxide/internal/shell"
	"github.com/eskelinenantti/tmuxide/internal/shell/fd"
	"github.com/eskelinenantti/tmuxide/internal/shell/fzf"
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
			Git:        shell.Git{},
			TmuxRunner: shell.CmdRunner{Command: "tmux"},
			FzfRunner:  shell.CmdRunner{Command: "fzf"},
			FdRunner:   shell.CmdRunner{Command: "fd"},
			Path:       shell.Path{},
		})
	},
}

func Edit(args []string, shell ShellEnv) error {
	tmux, err := tmux.InitTmux(shell.Path, shell.TmuxRunner)
	if err != nil {
		return err
	}

	editorCmd := strings.Fields(os.Getenv("EDITOR"))

	if len(editorCmd) == 0 {
		return ErrEditorEnvNotSet
	}

	if !shell.Path.Contains(editorCmd[0]) {
		return ErrEditorNotInstalled
	}

	var editorPath string
	if len(args) > 0 {
		editorPath = args[0]
	} else if editorPath, err = promptPath(tmux, fd.Fd{Runner: shell.FdRunner}, fzf.Fzf{Runner: shell.FzfRunner}); editorPath == "" || err != nil {
		return err
	}

	project, err := project.ForPath(editorPath, shell.Git, tmux)
	if err != nil {
		return fmt.Errorf("could not edit %s: %w", editorPath, err)
	}

	command := append(editorCmd, editorPath)
	return ide.Start(command, project, tmux)
}

func init() {
	rootCmd.AddCommand(editCmd)
}

func promptPath(tmux tmux.Tmux, fd fd.Fd, fzf fzf.Fzf) (string, error) {
	var input bytes.Buffer

	if out, err := tmux.ListSessions(); err == nil {
		input.Write(out)
	}

	out, err := fd.Execute()
	if err != nil {
		return "", fmt.Errorf("failed to run fd %w: %s", err, string(out))
	}
	input.Write(out)

	out, err = fzf.Execute(&input)
	if err != nil {
		// Hide the error as most likely user just cancelled the operation
		return "", nil
	}
	selection := strings.TrimSpace(string(out))
	return selection, nil
}
