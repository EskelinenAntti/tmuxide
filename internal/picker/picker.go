package picker

import (
	"io"
	"os"
	"path/filepath"

	"github.com/eskelinenantti/tmuxide/internal/shell/fd"
	"github.com/eskelinenantti/tmuxide/internal/shell/fzf"
)

func Prompt(fd fd.Cmd, finder fzf.Cmd) (string, error) {
	baseDir := os.Getenv("HOME")
	selection, err := finder.Select(func(output io.Writer) error {
		return fd.Fd(baseDir, output)
	})
	if err != nil || selection == "" {
		return "", err
	}

	return filepath.Join(baseDir, selection), nil
}
