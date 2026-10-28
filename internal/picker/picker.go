package picker

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/eskelinenantti/tmuxide/internal/shell/fd"
	"github.com/eskelinenantti/tmuxide/internal/shell/fzf"
	"github.com/eskelinenantti/tmuxide/internal/shell/pipe"
	"github.com/eskelinenantti/tmuxide/internal/shell/ranked"
)

func Prompt(ranked ranked.Cmd, fd fd.Cmd, fzf fzf.Cmd) (string, error) {
	var buffer bytes.Buffer
	fzfStdin, err := fzf.Fzf(&buffer)
	if err != nil {
		return "", err
	}
	defer fzfStdin.Kill()

	rankedStdout := ranked.Ranked()
	fdStdout, err := fd.Fd(context.Background())
	if err != nil {
		return "", err
	}
	defer fdStdout.Close()

	stdin := io.MultiReader(rankedStdout, fdStdout)
	if err := pipe.Pipe(pipe.DeduplicatePaths(stdin), fzfStdin); err != nil {
		return "", err
	}

	if err := fzfStdin.Close(); err != nil {
		return "", err
	}

	if err := fzfStdin.Wait(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 130 {
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
