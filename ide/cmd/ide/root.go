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
	"github.com/eskelinenantti/tmuxide/internal/shell/path"
	"github.com/eskelinenantti/tmuxide/internal/shell/runner"
	"github.com/eskelinenantti/tmuxide/internal/shell/tmux"
	"github.com/spf13/cobra"
)

type ShellEnv struct {
	Path      path.ShellPath
	CmdRunner runner.Runner
}

var rootCmd = &cobra.Command{
	Use:   "ide [session|file|folder]",
	Short: "Turn tmux and your favourite editor into an IDE with tmuxide.",
	Long: `tmuxide creates or switches to tmux sessions based on files and folders.

Run it without arguments to pick a location from a fuzzy finder. A tmux
session will be created (or reused) for that location automatically.

When a file is selected or passed as an argument, tmuxide opens it in
$EDITOR and creates the session for the repository root, or the file's
directory if it is not inside a git repository.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return Ide(args, runner.CmdRunner{}, path.Path{})
	},
}

var helpNoEditorConfigured = `
No editor was configured. Specify the editor you would like to use by setting the $EDITOR variable.
For example, to use Vim as your editor, add the following line to your ~/.zshrc or ~/.bashrc:
				
export EDITOR=vim`

var helpCommandNotInstalledTemplate = `
Did not find %s, which is a required dependency for ide command.

You can install %[1]s e.g. via homebrew by running:
brew install %[1]s`

var ErrEditorNotInstalled = errors.New("editor not installed")
var ErrEditorEnvNotSet = errors.New("editor not configured")

func Ide(args []string, runner runner.Runner, path path.ShellPath) error {
	shell, err := shell.Init(path, runner)
	if err != nil {
		return err
	}

	editorCmd, err := editorCmd(path)
	if err != nil {
		return err
	}

	var target string
	if len(args) == 0 {
		target, err = picker.Prompt(shell.Tmux, shell.Fd, shell.Fzf)
	} else {
		target = args[0]
	}

	if target == "" || err != nil {
		return err
	}

	isSession, isDir, err := targetDetails(target, shell.Tmux)
	if err != nil {
		return err
	}

	var proj project.Project
	if isSession {
		proj = project.ForSession(target)
	} else if isDir {
		proj, err = project.ForDir(target)
	} else {
		proj, err = project.ForFile(target, shell.Git, shell.Tmux)
	}

	if err != nil {
		return fmt.Errorf("could not edit %s: %w", target, err)
	}

	var command []string
	if !isSession && !isDir {
		command = append(editorCmd, target)
	}

	return ide.Start(command, proj, shell.Tmux)
}

func targetDetails(target string, tmux tmux.Cmd) (bool, bool, error) {
	isSession := tmux.HasSession(target, "")
	var isDir bool
	if !isSession {
		fileInfo, err := os.Stat(target)
		if err != nil {
			return false, false, project.ErrInvalidPath
		}

		isDir = fileInfo.IsDir()
	}
	return isSession, isDir, nil
}

func Execute() {
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = false
	err := rootCmd.Execute()

	if err == nil {
		return
	}

	if errors.Is(err, ErrEditorEnvNotSet) {
		rootCmd.PrintErrln(helpNoEditorConfigured)
	}

	var commandNotInstalled shell.NotInstalledError
	if errors.As(err, &commandNotInstalled) {
		rootCmd.PrintErrln(fmt.Sprintf(helpCommandNotInstalledTemplate, commandNotInstalled.Cmd))
	}

	os.Exit(1)
}

func editorCmd(path path.ShellPath) ([]string, error) {
	editorCmd := strings.Fields(os.Getenv("EDITOR"))

	if len(editorCmd) == 0 {
		return nil, ErrEditorEnvNotSet
	}

	if !path.Contains(editorCmd[0]) {
		return nil, ErrEditorNotInstalled
	}
	return editorCmd, nil
}

func init() {
}
