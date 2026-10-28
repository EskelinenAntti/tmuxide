package picker

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

func writeRanked(output io.Writer, paths []string) (map[string]bool, error) {
	home := os.Getenv("HOME")
	seen := make(map[string]bool)
	for _, path := range paths {
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
			if errors.Is(err, syscall.EPIPE) {
				return seen, nil
			}
			return nil, err
		}
	}
	return seen, nil
}

func copyUnseen(output io.Writer, input io.Reader, seen map[string]bool) error {
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		if err := writePath(output, scanner.Text(), seen); err != nil {
			return err
		}
	}
	return scanner.Err()
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

func selectedPath(output string) string {
	selection := strings.TrimSuffix(output, "\n")
	if selection == "" || filepath.IsAbs(selection) {
		return selection
	}
	return filepath.Join(os.Getenv("HOME"), selection)
}
