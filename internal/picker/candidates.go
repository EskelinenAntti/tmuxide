package picker

import (
	"io"
	"path/filepath"
	"strings"
)

// candidates streams unique paths without buffering the full filesystem scan.
// fd may split a path across writes, or include several paths in a single write.
type candidates struct {
	output  io.Writer
	seen    map[string]bool
	pending string
}

func (c *candidates) add(path string) error {
	key := filepath.Clean(path)
	if path == "" || c.seen[key] {
		return nil
	}
	_, err := io.WriteString(c.output, path+"\n")
	if err == nil {
		c.seen[key] = true
	}
	return err
}

func (c *candidates) Write(p []byte) (int, error) {
	c.pending += string(p)
	for {
		path, rest, ok := strings.Cut(c.pending, "\n")
		if !ok {
			return len(p), nil
		}
		c.pending = rest
		if err := c.add(path); err != nil {
			return 0, err
		}
	}
}
