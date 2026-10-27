package picker

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/eskelinenantti/tmuxide/internal/history"
	fdcmd "github.com/eskelinenantti/tmuxide/internal/shell/fd"
	"github.com/eskelinenantti/tmuxide/internal/shell/fzf"
)

func Prompt(fdCmd fdcmd.Cmd, fzf fzf.Cmd) (string, error) {
	var buffer bytes.Buffer
	fzfStdin, err := fzf.Fzf(&buffer)
	if err != nil {
		return "", err
	}

	fdErr := fd(fdCmd, fzfStdin)

	err = fzfStdin.Close()
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

func fd(cmd fdcmd.Cmd, output io.Writer) error {
	shownPaths := make(map[string]bool)
	writePath := func(path string) error {
		key := filepath.Clean(path)
		if path == "" || shownPaths[key] {
			return nil
		}
		if _, err := io.WriteString(output, path+"\n"); err != nil {
			return err
		}
		shownPaths[key] = true
		return nil
	}

	home := os.Getenv("HOME")
	paths, _ := history.Ranked()
	var folders []string
	for _, path := range paths {
		fileInfo, err := os.Stat(path)
		if err != nil || strings.ContainsAny(path, "\r\n") {
			continue
		}
		if relative, err := filepath.Rel(home, path); err == nil && filepath.IsLocal(relative) {
			path = relative
		}
		if fileInfo.IsDir() {
			folders = append(folders, path)
			path += string(filepath.Separator)
		}
		if err := writePath(path); err != nil {
			return err
		}
	}

	for _, folder := range folders {
		if err := cmd.Fd(writePath, folder); err != nil {
			return err
		}
	}
	return cmd.Fd(writePath)
}

func IsUserCancelledErr(err error) bool {
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		exitCode := exitErr.ExitCode()
		return exitCode == 130
	}
	return false
}
