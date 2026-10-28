package picker

import (
	"bufio"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/eskelinenantti/tmuxide/internal/history"
	"github.com/eskelinenantti/tmuxide/internal/shell/fd"
)

func writeCandidates(fd fd.Cmd, output io.Writer) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	home := os.Getenv("HOME")
	seen := make(map[string]bool)
	for _, path := range history.Ranked() {
		fileInfo, err := os.Stat(path)
		if err != nil || strings.ContainsAny(path, "\r\n") {
			continue
		}
		if relative, err := filepath.Rel(home, path); err == nil && filepath.IsLocal(relative) {
			path = relative
		}
		if fileInfo.IsDir() {
			path += string(filepath.Separator)
		}
		if err := writePath(output, path, seen); err != nil {
			return err
		}
	}

	input, err := fd.Fd(ctx)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		if err := writePath(output, scanner.Text(), seen); err != nil {
			cancel()
			input.Close()
			return err
		}
	}
	if scanner.Err() != nil {
		cancel()
	}
	return errors.Join(scanner.Err(), input.Close())
}

func writePath(output io.Writer, path string, seen map[string]bool) error {
	key := filepath.Clean(path)
	if path == "" || seen[key] {
		return nil
	}
	if _, err := io.WriteString(output, path+"\n"); err != nil {
		return err
	}
	seen[key] = true
	return nil
}
