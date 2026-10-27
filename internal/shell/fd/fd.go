package fd

import (
	"bufio"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/eskelinenantti/tmuxide/internal/history"
	"github.com/eskelinenantti/tmuxide/internal/shell/runner"
)

type Cmd struct {
	runner.Runner
}

func (f Cmd) Fd(ctx context.Context, output io.Writer) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	home := os.Getenv("HOME")
	seen := make(map[string]bool)
	var folders []string
	for _, path := range history.Ranked() {
		if err := ctx.Err(); err != nil {
			return err
		}
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
		if err := write(output, path, seen); err != nil {
			return err
		}
	}

	if !slices.Contains(folders, ".") {
		folders = append(folders, ".")
	}
	for _, folder := range folders {
		input, err := f.open(ctx, folder)
		if err != nil {
			return err
		}
		scanner := bufio.NewScanner(input)
		for scanner.Scan() {
			if err := write(output, scanner.Text(), seen); err != nil {
				cancel()
				input.Close()
				return err
			}
		}
		if scanner.Err() != nil {
			cancel()
		}
		if err := errors.Join(scanner.Err(), input.Close()); err != nil {
			return err
		}
	}
	return nil
}

func (f Cmd) open(ctx context.Context, folder string) (io.ReadCloser, error) {
	args := []string{
		"--follow",
		"--hidden",
		"--exclude", "{.git,node_modules,Library}",
		".",
		"--base-directory",
		os.Getenv("HOME"),
		"--",
	}
	if folder != "." {
		args = append(args, folder)
	}

	return f.Read(exec.CommandContext(ctx, "fd", args...))
}

func write(output io.Writer, path string, seen map[string]bool) error {
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
