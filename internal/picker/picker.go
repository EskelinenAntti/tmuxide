package picker

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/eskelinenantti/tmuxide/internal/shell/fd"
	"github.com/eskelinenantti/tmuxide/internal/shell/fzf"
)

func Prompt(fd fd.Cmd, finder fzf.Cmd) (string, error) {
	var buffer bytes.Buffer
	fzfStdin, err := finder.Fzf(&buffer)
	if err != nil {
		return "", err
	}

	err = fd.Fd(fzfStdin)
	if err != nil {
		return "", err
	}

	err = fzfStdin.Close()
	if err != nil {
		if fzf.IsCancelled(err) {
			return "", nil
		}
		return "", err
	}

	selection := strings.TrimSpace(buffer.String())
	return filepath.Join(os.Getenv("HOME"), selection), nil
}
