package pipe

import (
	"bufio"
	"io"
	"path/filepath"
)

// DeduplicatePaths keeps the first occurrence of each path without sorting the stream.
func DeduplicatePaths(stdin io.Reader) io.Reader {
	return &pathFilter{scanner: bufio.NewScanner(stdin), seen: make(map[string]bool)}
}

type pathFilter struct {
	scanner *bufio.Scanner
	seen    map[string]bool
	pending string
}

func (f *pathFilter) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	for f.pending == "" {
		if !f.scanner.Scan() {
			if err := f.scanner.Err(); err != nil {
				return 0, err
			}
			return 0, io.EOF
		}
		path := f.scanner.Text()
		key := filepath.Clean(path)
		if path == "" || f.seen[key] {
			continue
		}
		f.seen[key] = true
		f.pending = path + "\n"
	}
	n := copy(p, f.pending)
	f.pending = f.pending[n:]
	return n, nil
}
