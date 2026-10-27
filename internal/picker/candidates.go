package picker

import (
	"bytes"
	"io"
	"path/filepath"
)

// candidates streams unique paths without buffering the full filesystem scan.
// fd may split a path across writes, or include several paths in a single write.
type candidates struct {
	output  io.Writer
	seen    map[string]bool
	pending []byte
	err     error
}

func (c *candidates) add(path string) error {
	if c.err != nil {
		return c.err
	}
	key := filepath.Clean(path)
	if path == "" || c.seen[key] {
		return nil
	}
	_, c.err = io.WriteString(c.output, path+"\n")
	if c.err == nil {
		c.seen[key] = true
	}
	return c.err
}

func (c *candidates) Write(p []byte) (int, error) {
	c.pending = append(c.pending, p...)
	for {
		end := bytes.IndexByte(c.pending, '\n')
		if end < 0 {
			break
		}
		if err := c.add(string(c.pending[:end])); err != nil {
			return 0, err
		}
		c.pending = c.pending[end+1:]
	}
	return len(p), nil
}

func (c *candidates) flush() error {
	if len(c.pending) == 0 {
		return c.err
	}
	path := string(c.pending)
	c.pending = nil
	return c.add(path)
}
