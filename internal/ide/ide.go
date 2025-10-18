package ide

import (
	"os"
)

type Session interface {
	Exists() bool
	New(window Window) error
	NewWindow(window Window) error
	Attach() error
	Switch() error
}

func Start(session Session, windows []Window) error {
	if !session.Exists() {
		if err := create(session, windows); err != nil {
			return err
		}
	}

	if isAttached() {
		return session.Switch()
	}

	return session.Attach()
}

func create(session Session, windows []Window) error {
	if err := session.New(windows[0]); err != nil {
		return err
	}

	for _, window := range windows[1:] {
		if err := session.NewWindow(window); err != nil {
			return err
		}
	}

	return nil
}

func isAttached() bool {
	_, isAttached := os.LookupEnv("TMUX")
	return isAttached
}
