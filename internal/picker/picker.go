package picker

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/eskelinenantti/tmuxide/internal/shell/fd"
	"github.com/eskelinenantti/tmuxide/internal/shell/fzf"
	"github.com/eskelinenantti/tmuxide/internal/shell/tmux"
)

func Prompt(filterDir bool, tmux tmux.Cmd, fd fd.Cmd, fzf fzf.Cmd) (string, error) {
	var buffer bytes.Buffer
	fzfStdin, err := fzf.Fzf(&buffer)
	if err != nil {
		return "", err
	}

	sessionPrefix := "Session: "
	tmux.ListSessions(fzfStdin, sessionPrefix)
	err = fd.Fd(filterDir, fzfStdin)
	if err != nil {
		return "", err
	}

	err = fzfStdin.Close()
	if err != nil {
		// As a workaround, silence errors from fzf to not show an error if user closed it.
		return "", nil
	}

	selection := strings.TrimSpace(buffer.String())
	if sessionName, isSession := strings.CutPrefix(selection, sessionPrefix); isSession {
		return sessionName, nil
	}

	return filepath.Join(os.Getenv("HOME"), selection), nil
}
