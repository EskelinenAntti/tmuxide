package picker

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/eskelinenantti/tmuxide/internal/shell"
	"github.com/eskelinenantti/tmuxide/internal/shell/tmux"
)

func Prompt(tmux tmux.Tmux, fd shell.FdCmd, fzf shell.FzfCmd) (string, error) {
	args := []string{
		"--reverse",
		"--height",
		"30%",
	}
	fzfCmd := exec.Command("fzf", args...)

	var buffer bytes.Buffer
	fzfCmd.Stdout = &buffer
	fzfCmd.Stderr = os.Stderr

	pipe, err := fzfCmd.StdinPipe()
	if err != nil {
		return "", err
	}

	if err := fzfCmd.Start(); err != nil {
		return "", err
	}

	tmux.ListSessions(pipe)

	err = fd.Fd(pipe)
	if err != nil {
		return "", err
	}

	pipe.Close()
	err = fzfCmd.Wait()
	if err != nil {
		// As a workaround, silence errors from fzf to not show an error if user closed it.
		return "", nil
	}
	return strings.TrimSpace(buffer.String()), nil
}
