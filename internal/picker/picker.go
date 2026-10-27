package picker

import (
	"bytes"
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
	var buffer bytes.Buffer
	fzfStdin, err := fzf.Fzf(&buffer)
	if err != nil {
		return "", err
	}

	scanErr := search(fd, &candidates{output: fzfStdin, seen: make(map[string]bool)})

	err = fzfStdin.Close()
	if err != nil {
		if IsUserCancelledErr(err) {
			return "", nil
		}
		return "", err
	}
	selection := strings.TrimSuffix(buffer.String(), "\n")
	if selection == "" {
		return "", scanErr
	}
	// A selection is usable even if fzf closed its input before fd finished.
	if filepath.IsAbs(selection) {
		return selection, nil
	}
	return filepath.Join(os.Getenv("HOME"), selection), nil
}

func search(fd fd.Cmd, output *candidates) error {
	home := os.Getenv("HOME")
	// History is only an optimization: a missing or unreadable history must
	// never prevent the normal home-directory search.
	paths, _ := history.Ranked()
	var folders []string
	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil || strings.ContainsAny(path, "\r\n") {
			continue
		}
		if relative, err := filepath.Rel(home, path); err == nil && filepath.IsLocal(relative) {
			path = relative
		}
		if info.IsDir() {
			folders = append(folders, path)
			path += string(filepath.Separator)
		}
		if err := output.add(path); err != nil {
			return err
		}
	}

	// Seed all known paths first, then scan selected folders in frequency
	// order before falling back to the rest of home.
	for _, folder := range folders {
		if err := fd.Fd(output, folder); err != nil {
			return err
		}
	}
	return fd.Fd(output)
}

func IsUserCancelledErr(err error) bool {
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		exitCode := exitErr.ExitCode()
		return exitCode == 130
	}
	return false
}
