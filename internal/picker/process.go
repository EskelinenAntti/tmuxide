package picker

import (
	"context"
	"errors"
	"io"
	"os/exec"
	"syscall"

	"github.com/eskelinenantti/tmuxide/internal/shell/runner"
)

func streamUnseen(output io.Writer, input io.ReadCloser, seen map[string]bool, cancel context.CancelFunc) error {
	err := copyUnseen(output, input, seen)
	if err != nil {
		cancel()
	}
	closeErr := input.Close()
	if errors.Is(err, syscall.EPIPE) {
		// Selecting or cancelling in fzf closes its input before discovery finishes.
		return nil
	}
	return errors.Join(err, closeErr)
}

func cleanupPicker(process runner.Process) {
	// Best-effort cleanup also covers early returns before the explicit close and wait.
	process.Close()
	process.Wait()
}

func IsUserCancelledErr(err error) bool {
	var exitErr *exec.ExitError
	return errors.As(err, &exitErr) && exitErr.ExitCode() == 130
}
