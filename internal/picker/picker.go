package picker

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

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
	// Selecting an entry can close fzf's input before fd finishes streaming.
	if scanErr != nil && !errors.Is(scanErr, syscall.EPIPE) {
		return "", scanErr
	}

	selection := strings.TrimSuffix(buffer.String(), "\n")
	if selection == "" || filepath.IsAbs(selection) {
		return selection, nil
	}
	return filepath.Join(os.Getenv("HOME"), selection), nil
}

func search(fd fd.Cmd, output *candidates) error {
	home := os.Getenv("HOME")
	// History is only an optimization: a missing or unreadable history must
	// never prevent the normal home-directory search.
	entries, _ := history.Default().Ranked()
	var folders []string
	for _, entry := range entries {
		info, err := os.Stat(entry.Path)
		if err != nil || strings.ContainsAny(entry.Path, "\r\n") {
			continue
		}
		path := entry.Path
		if relative, err := filepath.Rel(home, path); err == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
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
	// order before falling back to the rest of home. Prefer output errors to
	// fd's resulting exit status so an early selection remains a broken pipe.
	for _, folder := range folders {
		if err := fd.Fd(output, folder); err != nil && output.err == nil {
			return err
		}
		if err := output.flush(); err != nil {
			return err
		}
	}
	if err := fd.Fd(output); err != nil && output.err == nil {
		return err
	}
	return output.flush()
}

func IsUserCancelledErr(err error) bool {
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		exitCode := exitErr.ExitCode()
		return exitCode == 130
	}
	return false
}
