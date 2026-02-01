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

	// Sequentially copy tmux output
	err = tmux.ListSessions(pipe)
	if err != nil {
		return "", nil
	}

	err = fd.Fd(pipe)
	if err != nil {
		return "", nil
	}

	pipe.Close()
	err = fzfCmd.Wait()
	if err == nil {
		return "", nil
	}
	return strings.TrimSpace(buffer.String()), nil
}
