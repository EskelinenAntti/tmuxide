package picker

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/eskelinenantti/tmuxide/internal/shell/fd"
	"github.com/eskelinenantti/tmuxide/internal/shell/fzf"
)

func Prompt(fd fd.Cmd, fzf fzf.Cmd) (string, error) {
	var buffer bytes.Buffer
	fzfStdin, err := fzf.Fzf(&buffer)
	if err != nil {
		return "", err
	}

	fdErr := fd.Fd(fzfStdin)
	fdErr = errors.Join(fdErr, fzfStdin.Close())
	err = fzfStdin.Wait()
	if err != nil {
		if IsUserCancelledErr(err) {
			return "", nil
		}
		return "", err
	}

	selection := strings.TrimSuffix(buffer.String(), "\n")
	if selection == "" {
		return "", fdErr
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
