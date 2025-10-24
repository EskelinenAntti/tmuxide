package windows

import (
	"errors"
	"fmt"
	"os"

	"github.com/eskelinenantti/tmuxide/internal/ide/window"
	"github.com/eskelinenantti/tmuxide/internal/project"
)

var ErrEditorNotSet = errors.New(
	"No editor was configured. Specify the editor you would like to use by setting the $EDITOR variable.\n\n" +
		"For example, to use Vim as your editor, add the following line to your ~/.zshrc:\n" +
		"export EDITOR=vim\n",
)

var ErrTmuxNotInPath = errors.New(
	"Did not find tmux, which is a required dependency for ide command.\n\n" +

		"You can install tmux e.g. via homebrew by running\n" +
		"brew install tmux\n",
)

func Get(project project.Project, path window.Path) ([]window.Window, error) {
	if !path.Contains("tmux") {
		return nil, ErrTmuxNotInPath
	}

	editor, err := editor(project, path)
	if err != nil {
		return nil, err
	}
	windows := []window.Window{editor}

	lazygit, err := lazygit(project, path)
	if err == nil {
		windows = append(windows, lazygit)
	}
	return windows, nil
}

func lazygit(project project.Project, path window.Path) (window.Window, error) {
	if !path.Contains("lazygit") {
		return window.Window{}, errors.New("Lazygit is not installed")
	}

	if !project.IsGitRepo {
		return window.Window{}, errors.New("Not insGit repository")
	}

	return window.Window{"lazygit"}, nil
}

func editor(project project.Project, path window.Path) (window.Window, error) {
	editorCmd, hasEditor := os.LookupEnv("EDITOR")
	if !hasEditor || editorCmd == "" {
		return window.Window{}, ErrEditorNotSet
	}

	if !path.Contains(editorCmd) {
		return window.Window{}, fmt.Errorf(
			"Editor %s is not installed", editorCmd,
		)
	}

	return window.Window{editorCmd, project.TargetPath}, nil
}
