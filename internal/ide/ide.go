package ide

import (
	"errors"
	"os"
)

type Session interface {
	Exists() bool
	New(command string, args ...string) error
	Attach() error
	Switch() error
}

func Start(target string, session Session) error {
	editor, err := editor()
	if err != nil {
		return err
	}

	if !session.Exists() {
		if err := session.New(editor, target); err != nil {
			return err
		}
	}

	if isAttachedToSession() {
		return session.Switch()
	}

	return session.Attach()
}

func isAttachedToSession() bool {
	_, alreadyInSession := os.LookupEnv("TMUX")
	return alreadyInSession
}

func editor() (string, error) {
	editor, hasEditor := os.LookupEnv("EDITOR")
	if !hasEditor {
		return "", errors.New(
			"No editor was configured. Specify the editor you would like to use by setting the $EDITOR variable.\n\n" +
				"For example, to use Vim as your editor, add the following line to your ~/.zshrc:\n" +
				"export EDITOR=vim\n",
		)
	}
	return editor, nil
}
