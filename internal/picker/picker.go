package picker

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/eskelinenantti/tmuxide/internal/history"
	"github.com/eskelinenantti/tmuxide/internal/shell/fd"
	"github.com/eskelinenantti/tmuxide/internal/shell/fzf"
)

func Prompt(fd fd.Cmd, fzf fzf.Cmd) (string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var buffer bytes.Buffer
	fzfStdin, err := fzf.Fzf(&buffer)
	if err != nil {
		return "", err
	}

	seen := make(map[string]bool)
	candidateErr := writeRanked(fzfStdin, history.Ranked(), seen)
	if candidateErr == nil {
		input, err := fd.Fd(ctx)
		if err != nil {
			candidateErr = err
		} else {
			candidateErr = copyUnseen(fzfStdin, input, seen)
			if candidateErr != nil {
				cancel()
			}
			candidateErr = errors.Join(candidateErr, input.Close())
		}
	}
	candidateErr = errors.Join(candidateErr, fzfStdin.Close())
	err = fzfStdin.Wait()
	if err != nil {
		if IsUserCancelledErr(err) {
			return "", nil
		}
		return "", err
	}

	selection := strings.TrimSuffix(buffer.String(), "\n")
	if selection == "" {
		return "", candidateErr
	}
	if filepath.IsAbs(selection) {
		return selection, nil
	}
	return filepath.Join(os.Getenv("HOME"), selection), nil
}

func IsUserCancelledErr(err error) bool {
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		exitCode := exitErr.ExitCode()
		return exitCode == 130
	}
	return false
}
