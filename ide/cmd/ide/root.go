package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/eskelinenantti/tmuxide/internal/multiplexer"
	"github.com/eskelinenantti/tmuxide/internal/picker"
	"github.com/eskelinenantti/tmuxide/internal/project"
	"github.com/eskelinenantti/tmuxide/internal/shell"
	"github.com/eskelinenantti/tmuxide/internal/shell/path"
	"github.com/eskelinenantti/tmuxide/internal/shell/runner"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ide [file|folder]",
	Short: "tmuxide creates or switches to tmux sessions based on files and folders.",
	Args:  cobra.MaximumNArgs(1),
	Long: `tmuxide creates or switches to tmux sessions based on files and folders.

Run it without arguments to pick a location from a fuzzy finder. A tmux
session will be created (or reused) for that location automatically.

When a file is selected or passed as an argument, tmuxide opens it in
$EDITOR and creates the session for the repository root, or the file's
directory if it is not inside a git repository.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		mux, err := cmd.Flags().GetString("mux")
		if err != nil {
			return err
		}
		return Run(args, mux, runner.CmdRunner{}, path.Path{})
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

func Run(args []string, muxFlag string, runner shell.Runner, path path.ShellPath) error {
	mux, err := multiplexer.Select(muxFlag, path, runner)
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
		target, err = picker.Prompt(shell.Fd, shell.Fzf)
	} else {
		target = args[0]
	}

	if target == "" || err != nil {
		return err
	}

	resolved, err := project.ForPath(target, shell.Git)
	if err != nil {
		if errors.Is(err, project.ErrInvalidPath) {
			return err
		}
		return fmt.Errorf("could not edit %s: %w", target, err)
	}

	var command []string
	if resolved.File != "" {
		command = append(editorCmd, resolved.File)
	}

	return mux.Open(resolved.Project, command)
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
	rootCmd.Flags().String("mux", "auto", "Terminal host override")
}
