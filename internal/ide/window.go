package ide

import (
	"errors"
	"fmt"
	"os"

	"github.com/eskelinenantti/tmuxide/internal/input"
	"github.com/eskelinenantti/tmuxide/internal/project"
)

type ShellPath interface {
	Contains(path string) bool
}

type Window input.Args

var ErrEditorNotSet = errors.New(
	"No editor was configured. Specify the editor you would like to use by setting the $EDITOR variable.\n\n" +
		"For example, to use Vim as your editor, add the following line to your ~/.zshrc:\n" +
		"export EDITOR=vim\n",
)

var ErrTmuxNotInstalled = errors.New(
	"Did not find tmux, which is a required dependency for ide command.\n\n" +

		"You can install tmux e.g. via homebrew by running\n" +
		"brew install tmux\n",
)

func Windows(project project.Project, path ShellPath) ([]Window, error) {
	if !path.Contains("tmux") {
		return nil, ErrTmuxNotInstalled
	}

	editor, err := editor(project, path)
	if err != nil {
		return nil, err
	}
	windows := []Window{editor}

	lazygit, err := lazygit(project, path)
	if err == nil {
		windows = append(windows, lazygit)
	}
	return windows, nil
}

func lazygit(project project.Project, path ShellPath) (Window, error) {
	if !path.Contains("lazygit") {
		return Window{}, errors.New("Lazygit is not installed")
	}

	if !project.IsGitRepo {
		return Window{}, errors.New("Not insGit repository")
	}

	return Window{"lazygit"}, nil
}

func editor(project project.Project, path ShellPath) (Window, error) {
	editorCmd, hasEditor := os.LookupEnv("EDITOR")
	if !hasEditor || editorCmd == "" {
		return Window{}, ErrEditorNotSet
	}

	if !path.Contains(editorCmd) {
		return Window{}, fmt.Errorf(
			"Editor %s is not installed", editorCmd,
		)
	}

	return Window{editorCmd, project.TargetPath}, nil
}
