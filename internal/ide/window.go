package ide

import (
	"errors"
	"os"
	"os/exec"

	"github.com/eskelinenantti/tmuxide/internal/tmux"
)

func WindowsFor(target string, repository string) ([]tmux.Cmd, error) {
	editor, err := editor(target)
	if err != nil {
		return []tmux.Cmd{}, err
	}

	windows := []tmux.Cmd{editor}

	if lazygit, err := lazygit(repository); err == nil {
		windows = append(windows, lazygit)
	}

	return windows, nil
}

func lazygit(repository string) (tmux.Cmd, error) {
	if _, err := exec.LookPath("lazygit"); err != nil {
		return tmux.Cmd{}, errors.New("Lazygit is not installed")
	}

	if repository == "" {
		return tmux.Cmd{}, errors.New("Not inside Git repository")
	}

	return tmux.Cmd{Cmd: "lazygit"}, nil
}

func editor(target string) (tmux.Cmd, error) {
	editorCmd, hasEditor := os.LookupEnv("EDITOR")
	if !hasEditor {
		return tmux.Cmd{}, errors.New(
			"No editor was configured. Specify the editor you would like to use by setting the $EDITOR variable.\n\n" +
				"For example, to use Vim as your editor, add the following line to your ~/.zshrc:\n" +
				"export EDITOR=vim\n",
		)
	}

	return tmux.Cmd{Cmd: editorCmd, Args: []string{target}}, nil
}
