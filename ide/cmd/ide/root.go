package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/eskelinenantti/tmuxide/internal/ide"
	"github.com/eskelinenantti/tmuxide/internal/picker"
	"github.com/eskelinenantti/tmuxide/internal/project"
	"github.com/eskelinenantti/tmuxide/internal/shell"
	"github.com/eskelinenantti/tmuxide/internal/shell/path"
	"github.com/eskelinenantti/tmuxide/internal/shell/runner"
	"github.com/eskelinenantti/tmuxide/internal/workspace"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ide [file|folder]",
	Short: "Open projects in tmux sessions or Herdr workspaces.",
	Args:  cobra.MaximumNArgs(1),
	Long: `tmuxide opens projects in tmux sessions or native Herdr workspaces.

Run it without arguments to pick a location from a fuzzy finder. A project
workspace will be created (or reused) for that location automatically.
Inside Herdr, the Herdr backend is selected automatically; otherwise tmux
is used. Use --backend to override this choice.

When a file is selected or passed as an argument, tmuxide opens it in
$EDITOR and creates the workspace for the repository root, or the file's
directory if it is not inside a git repository.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		backend, err := cmd.Flags().GetString("backend")
		if err != nil {
			return err
		}
		return run(args, backend, runner.CmdRunner{}, path.Path{})
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
	return run(args, "auto", runner, path)
}

func run(args []string, backend string, runner runner.Runner, path path.ShellPath) error {
	host, fullscreen, err := workspace.Select(backend, path, runner)
	if err != nil {
		return err
	}
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
		shell.Fzf.Fullscreen = fullscreen
		target, err = picker.Prompt(shell.Fd, shell.Fzf)
	} else {
		target = args[0]
	}

	if target == "" || err != nil {
		return err
	}

	isDir, err := isDir(target)
	if err != nil {
		return err
	}

	var proj project.Project
	if isDir {
		proj, err = project.ForDir(target)
	} else {
		proj, err = project.ForFile(target, shell.Git)
	}

	if err != nil {
		return fmt.Errorf("could not edit %s: %w", target, err)
	}

	var command []string
	if !isDir {
		// The editor starts in the project root, not the caller's directory.
		target, err = filepath.Abs(target)
		if err != nil {
			return err
		}
		command = append(editorCmd, target)
	}

	return ide.Start(command, proj, host)
}

func isDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, project.ErrInvalidPath
	}

	return fileInfo.IsDir(), nil
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
	rootCmd.Flags().String("backend", "auto", "Terminal host: auto, tmux, or herdr")
}
