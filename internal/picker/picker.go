package picker

import (
	"bytes"
	"cmp"
	"errors"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/eskelinenantti/herdr-zipln/internal/history"
	"github.com/eskelinenantti/herdr-zipln/internal/shell/fd"
	"github.com/eskelinenantti/herdr-zipln/internal/shell/fzf"
	"github.com/eskelinenantti/herdr-zipln/internal/shell/pipe"
)

func Prompt(fdCmd fd.Cmd, fzfCmd fzf.Cmd) (string, error) {
	var buffer bytes.Buffer
	fdStdout, err := fdCmd.Fd()
	if err != nil {
		return "", err
	}
	defer fdStdout.Close()

	// Cache a single scan for reloads; typing must not walk the filesystem again.
	data, err := io.ReadAll(pipe.DeduplicatePaths(fdStdout))
	if err != nil {
		return "", err
	}
	paths := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
	slices.SortFunc(paths, func(a, b string) int {
		if order := cmp.Compare(len(a), len(b)); order != 0 {
			return order
		}
		return cmp.Compare(a, b)
	})
	data = []byte(strings.Join(paths, "\n") + "\n")
	candidates, err := os.CreateTemp("", "zipln-candidates-*")
	if err != nil {
		return "", err
	}
	defer os.Remove(candidates.Name())
	if _, err := candidates.Write(data); err != nil {
		candidates.Close()
		return "", err
	}
	if err := candidates.Close(); err != nil {
		return "", err
	}
	if err := fzfCmd.Fzf(bytes.NewReader(data), &buffer, candidates.Name()); err != nil {
		if errors.Is(err, fzf.ErrCancelled) {
			return "", nil
		}
		return "", err
	}

	query, selection, ok := strings.Cut(strings.TrimSuffix(buffer.String(), "\n"), "\n")
	if !ok || selection == "" {
		return "", nil
	}
	if !filepath.IsAbs(selection) {
		selection = filepath.Join(os.Getenv("HOME"), selection)
	}
	history.RecordQuery(query, selection)
	return selection, nil
}
