package picker

import (
	"os"
	"path/filepath"

	"github.com/eskelinenantti/tmuxide/internal/shell/fd"
	"github.com/eskelinenantti/tmuxide/internal/shell/fzf"
)

func Prompt(fd fd.Cmd, finder fzf.Cmd) (string, error) {
	selection, err := finder.Select(fd.Fd)
	if err != nil || selection == "" {
		return "", err
	}

	return filepath.Join(os.Getenv("HOME"), selection), nil
}
