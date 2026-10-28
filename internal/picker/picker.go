package picker

import (
	"bytes"
	"context"

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
	defer cleanupPicker(fzfStdin)

	seen, err := writeRanked(fzfStdin, history.Ranked())
	if err != nil {
		return "", err
	}

	results, err := fd.Fd(ctx)
	if err != nil {
		return "", err
	}

	if err := streamUnseen(fzfStdin, results, seen, cancel); err != nil {
		return "", err
	}

	if err := fzfStdin.Close(); err != nil {
		return "", err
	}

	if err := fzfStdin.Wait(); err != nil {
		if IsUserCancelledErr(err) {
			return "", nil
		}
		return "", err
	}

	return selectedPath(buffer.String()), nil
}
