package ide

import (
	"errors"
	"os"
	"os/exec"
)

func WindowsFor(target string, repository string) ([]Window, error) {
	editor, err := editor(target)
	if err != nil {
		return []Window{}, err
	}

	windows := []Window{editor}

	if lazygit, err := lazygit(repository); err != nil {
		windows = append(windows, lazygit)
	}

	return windows, nil
}

func lazygit(repository string) (Window, error) {
	if _, err := exec.LookPath("lazygit"); err != nil {
		return Window{}, errors.New("Lazygit is not isntalled")
	}

	if repository == "" {
		return Window{}, errors.New("Not inside Git repository")
	}

	return Window{Cmd: "lazygit"}, nil
}

func editor(target string) (Window, error) {
	editorCmd, hasEditor := os.LookupEnv("EDITOR")
	if !hasEditor {
		return Window{}, errors.New(
			"No editor was configured. Specify the editor you would like to use by setting the $EDITOR variable.\n\n" +
				"For example, to use Vim as your editor, add the following line to your ~/.zshrc:\n" +
				"export EDITOR=vim\n",
		)
	}

	return Window{Cmd: editorCmd, Args: []string{target}}, nil
}
