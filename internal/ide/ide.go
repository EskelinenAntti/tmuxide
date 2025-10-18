package ide

import (
	"os"
)

type Session interface {
	Exists() bool
	New() error
	Attach() error
	Switch() error
}

func Start(session Session) error {
	if !session.Exists() {
		if err := session.New(); err != nil {
			return err
		}
	}

	if isAttached() {
		return session.Switch()
	}

	return session.Attach()
}

func isAttached() bool {
	_, isAttached := os.LookupEnv("TMUX")
	return isAttached
}
