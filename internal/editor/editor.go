package editor

import (
	"errors"
	"os"
	"strings"

	"github.com/eskelinenantti/tmuxide/internal/shell/path"
)

var ErrNotInstalled = errors.New("editor not installed")
var ErrNotConfigured = errors.New("editor not configured")

// Command resolves $EDITOR to an installed executable and its arguments.
func Command(path path.ShellPath) ([]string, error) {
	command := strings.Fields(os.Getenv("EDITOR"))
	if len(command) == 0 {
		return nil, ErrNotConfigured
	}
	if !path.Contains(command[0]) {
		return nil, ErrNotInstalled
	}
	return command, nil
}
