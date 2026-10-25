package herdr

import (
	"fmt"
	"os"

	"github.com/eskelinenantti/tmuxide/internal/shell"
	"github.com/eskelinenantti/tmuxide/internal/shell/path"
	"github.com/eskelinenantti/tmuxide/internal/shell/runner"
)

// New configures a Herdr command using the executable supplied by the environment.
func New(path path.ShellPath, runner runner.Runner) (Cmd, error) {
	binary := os.Getenv("HERDR_BIN_PATH")
	if binary == "" {
		return Cmd{}, fmt.Errorf("HERDR_BIN_PATH is required for the Herdr multiplexer")
	}
	if !path.Contains(binary) {
		return Cmd{}, shell.NotInstalledError{Cmd: binary}
	}
	return Cmd{Runner: runner, Binary: binary}, nil
}

// InSession reports whether the current process is running inside Herdr.
func InSession() bool {
	return os.Getenv("HERDR_ENV") == "1"
}
