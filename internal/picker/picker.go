package picker

import (
	"bytes"
	"strings"

	"github.com/eskelinenantti/tmuxide/internal/shell"
	"github.com/eskelinenantti/tmuxide/internal/shell/tmux"
)

func Prompt(tmux tmux.Tmux, fd shell.FdCmd, fzf shell.FzfCmd) (string, error) {
	var buffer bytes.Buffer
	fzfStdin, err := fzf.Fzf(&buffer)
	if err != nil {
		return "", err
	}

	tmux.ListSessions(fzfStdin)
	err = fd.Fd(fzfStdin)
	if err != nil {
		return "", err
	}

	err = fzfStdin.Close()
	if err != nil {
		// As a workaround, silence errors from fzf to not show an error if user closed it.
		return "", nil
	}
	return strings.TrimSpace(buffer.String()), nil
}
