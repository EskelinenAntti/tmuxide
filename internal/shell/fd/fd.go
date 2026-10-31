package fd

import (
	"io"
	"os"
	"os/exec"

	"github.com/eskelinenantti/herdr-zipln/internal/shell/runner"
)

type Cmd struct {
	runner.Runner
}

func (f Cmd) Fd() (io.ReadCloser, error) {
	args := []string{
		"--follow",
		"--hidden",
		"--exclude", "{.git,node_modules,Library}",
		".",
		"--base-directory",
		os.Getenv("HOME"),
	}
	// Finish the shallow pass before traversing deeper trees. The disjoint depth
	// ranges keep useful paths early without buffering or duplicating results.
	shallow, err := f.Stdout(exec.Command("fd", append(args, "--max-depth", "3")...))
	if err != nil {
		return nil, err
	}
	return &scan{
		ReadCloser: shallow,
		next: func() (io.ReadCloser, error) {
			return f.Stdout(exec.Command("fd", append(args, "--min-depth", "4")...))
		},
	}, nil
}

type scan struct {
	io.ReadCloser
	next func() (io.ReadCloser, error)
}

func (s *scan) Read(p []byte) (int, error) {
	for {
		n, err := s.ReadCloser.Read(p)
		if err != io.EOF || s.next == nil {
			return n, err
		}
		if n > 0 {
			return n, nil
		}
		if err := s.ReadCloser.Close(); err != nil {
			return 0, err
		}
		next := s.next
		s.next = nil
		deeper, err := next()
		if err != nil {
			return 0, err
		}
		s.ReadCloser = deeper
	}
}

func (s *scan) Close() error {
	// Closing an early selection stops the active pass without starting another.
	s.next = nil
	return s.ReadCloser.Close()
}
