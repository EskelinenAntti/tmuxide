package picker

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/eskelinenantti/tmuxide/internal/ranked"
	"github.com/eskelinenantti/tmuxide/internal/shell/fd"
	"github.com/eskelinenantti/tmuxide/internal/shell/fzf"
	"github.com/eskelinenantti/tmuxide/internal/shell/pipe"
)

func Prompt(ranked ranked.History, fdCmd fd.Cmd, fzfCmd fzf.Cmd) (string, error) {
	var buffer bytes.Buffer
	rankedStdout := ranked.Ranked()
	fdStdout, err := fdCmd.Fd()
	if err != nil {
		return "", err
	}
	defer fdStdout.Close()

	stdin := io.MultiReader(rankedStdout, fdStdout)
	if err := fzfCmd.Fzf(pipe.DeduplicatePaths(stdin), &buffer); err != nil {
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
