package picker

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/eskelinenantti/tmuxide/internal/shell/fd"
	"github.com/eskelinenantti/tmuxide/internal/shell/fzf"
	"github.com/eskelinenantti/tmuxide/internal/shell/tmux"
)

func Prompt(tmux tmux.Tmux, fd fd.Fd, fzf fzf.Fzf) (string, error) {
	var input bytes.Buffer

	if out, err := tmux.ListSessions(); err == nil {
		input.Write(out)
	}

	out, err := fd.Execute()
	if err != nil {
		return "", fmt.Errorf("failed to run fd %w: %s", err, string(out))
	}
	input.Write(out)

	out, err = fzf.Execute(&input)
	if err != nil {
		// Hide the error as most likely user just cancelled the operation
		return "", nil
	}
	selection := strings.TrimSpace(string(out))
	return selection, nil
}
