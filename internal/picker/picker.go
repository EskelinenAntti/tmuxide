package picker

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/eskelinenantti/herdr-zipln/internal/shell/fd"
	"github.com/eskelinenantti/herdr-zipln/internal/shell/fzf"
)

func Prompt(fdCmd fd.Cmd, fzfCmd fzf.Cmd) (string, error) {
	var buffer bytes.Buffer
	fdStdout, err := fdCmd.Fd()
	if err != nil {
		return "", err
	}
	defer fdStdout.Close()

	if err := fzfCmd.Fzf(fdStdout, &buffer); err != nil {
		if errors.Is(err, fzf.ErrCancelled) {
			return "", nil
		}
		return "", err
	}

	selection := strings.TrimSuffix(buffer.String(), "\n")
	if selection == "" || filepath.IsAbs(selection) {
		return selection, nil
	}
	return filepath.Join(os.Getenv("HOME"), selection), nil
}
