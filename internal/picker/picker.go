package picker

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/eskelinenantti/herdr-zipln/internal/history"
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

	query, selection, ok := strings.Cut(strings.TrimSuffix(buffer.String(), "\n"), "\n")
	if !ok || selection == "" {
		return "", nil
	}
	if !filepath.IsAbs(selection) {
		selection = filepath.Join(os.Getenv("HOME"), selection)
	}
	history.RecordQuery(query, selection)
	return selection, nil
}
