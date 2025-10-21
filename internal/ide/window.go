package ide

import (
	"errors"
	"fmt"
	"os"

	"github.com/eskelinenantti/tmuxide/internal/path"
	"github.com/eskelinenantti/tmuxide/internal/shell"
	"github.com/eskelinenantti/tmuxide/internal/tmux"
)

var ErrEditorNotSet = errors.New(
	"No editor was configured. Specify the editor you would like to use by setting the $EDITOR variable.\n\n" +
		"For example, to use Vim as your editor, add the following line to your ~/.zshrc:\n" +
		"export EDITOR=vim\n",
)

func windowsFor(target string, repository string, shell shell.Shell) ([]tmux.WindowCommand, error) {
	editor, err := editor(target, shell.Path)
	if err != nil {
		return []tmux.WindowCommand{}, err
	}

	windows := []tmux.WindowCommand{editor}

	if lazygit, err := lazygit(repository, shell.Path); err == nil {
		windows = append(windows, lazygit)
	}

	return windows, nil
}

func lazygit(repository string, path path.Path) (tmux.WindowCommand, error) {
	if !path.Contains("lazygit") {
		return tmux.WindowCommand{}, errors.New("Lazygit is not installed")
	}

	if repository == "" {
		return tmux.WindowCommand{}, errors.New("Not inside Git repository")
	}

	return tmux.WindowCommand{Cmd: "lazygit"}, nil
}

func editor(target string, path path.Path) (tmux.WindowCommand, error) {
	editorCmd, hasEditor := os.LookupEnv("EDITOR")
	if !hasEditor || editorCmd == "" {
		return tmux.WindowCommand{}, ErrEditorNotSet
	}

	if !path.Contains(editorCmd) {
		return tmux.WindowCommand{}, fmt.Errorf(
			"Editor %s is not installed", editorCmd,
		)
	}

	return tmux.WindowCommand{Cmd: editorCmd, Args: []string{target}}, nil
}
