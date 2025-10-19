package ide

import (
	"errors"
	"fmt"
	"os"

	"github.com/eskelinenantti/tmuxide/internal/path"
	"github.com/eskelinenantti/tmuxide/internal/tmux"
)

func WindowsFor(target string, repository string, pathLooker path.PathLooker) ([]tmux.Cmd, error) {
	editor, err := editor(target, pathLooker)
	if err != nil {
		return []tmux.Cmd{}, err
	}

	windows := []tmux.Cmd{editor}

	if lazygit, err := lazygit(repository, pathLooker); err == nil {
		windows = append(windows, lazygit)
	}

	return windows, nil
}

func lazygit(repository string, pathLooker path.PathLooker) (tmux.Cmd, error) {
	if _, err := pathLooker.LookPath("lazygit"); err != nil {
		return tmux.Cmd{}, errors.New("Lazygit is not installed")
	}

	if repository == "" {
		return tmux.Cmd{}, errors.New("Not inside Git repository")
	}

	return tmux.Cmd{Cmd: "lazygit"}, nil
}

func editor(target string, pathLooker path.PathLooker) (tmux.Cmd, error) {
	editorCmd, hasEditor := os.LookupEnv("EDITOR")
	if !hasEditor {
		return tmux.Cmd{}, errors.New(
			"No editor was configured. Specify the editor you would like to use by setting the $EDITOR variable.\n\n" +
				"For example, to use Vim as your editor, add the following line to your ~/.zshrc:\n" +
				"export EDITOR=vim\n",
		)
	}

	_, err := pathLooker.LookPath(editorCmd)
	if err != nil {
		return tmux.Cmd{}, fmt.Errorf(
			"Editor %s is not installed: %w", editorCmd, err,
		)
	}

	return tmux.Cmd{Cmd: editorCmd, Args: []string{target}}, nil
}
