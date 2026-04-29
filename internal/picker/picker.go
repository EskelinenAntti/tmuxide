package picker

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/eskelinenantti/tmuxide/internal/shell/fd"
	"github.com/eskelinenantti/tmuxide/internal/shell/fzf"
	"github.com/eskelinenantti/tmuxide/internal/shell/tmux"
)

func Prompt(tmux tmux.Cmd, fd fd.Cmd, fzf fzf.Cmd) (string, error) {
	var buffer bytes.Buffer
	fzfStdin, err := fzf.Fzf(&buffer)
	if err != nil {
		return "", err
	}

	sessionPrefix := "Session: "
	tmux.ListSessions(fzfStdin, sessionPrefix)
	err = fd.Fd(fzfStdin)
	if err != nil {
		return "", err
	}

	err = fzfStdin.Close()
	if err != nil {
		if IsUserCancelledErr(err) {
			return "", nil
		}
		return "", err
	}

	selection := strings.TrimSpace(buffer.String())
	if sessionName, isSession := strings.CutPrefix(selection, sessionPrefix); isSession {
		return sessionName, nil
	}

	return filepath.Join(os.Getenv("HOME"), selection), nil
}

func IsUserCancelledErr(err error) bool {
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		exitCode := exitErr.ExitCode()
		return exitCode == 130
	}
	return false
}
