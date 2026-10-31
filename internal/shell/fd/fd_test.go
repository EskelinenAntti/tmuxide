package fd

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/eskelinenantti/herdr-zipln/internal/shell/runner"
	"github.com/eskelinenantti/herdr-zipln/internal/test/mock"
	"github.com/eskelinenantti/herdr-zipln/internal/test/spy"
	"github.com/google/go-cmp/cmp"
)

func TestFdStreamsShallowBeforeDeep(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	shallow := "./project/\nproject/file.txt\n"
	deep := "project/src/nested/file.txt\n"
	runner := &spy.SpyRunner{Responses: []spy.Response{
		{OnRun: mock.WriteToStdout(shallow)},
		{OnRun: mock.WriteToStdout(deep)},
	}}
	input, err := (Cmd{Runner: runner}).Fd()
	if err != nil {
		t.Fatal(err)
	}
	defer input.Close()
	first := make([]byte, len(shallow))
	if _, err := io.ReadFull(input, first); err != nil {
		t.Fatal(err)
	}
	if string(first) != shallow || len(runner.Calls) != 1 {
		t.Fatalf("shallow output = %q, calls = %v; deep scan started too early", first, runner.Calls)
	}
	output, err := io.ReadAll(input)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(deep, string(output)); diff != "" {
		t.Fatal(diff)
	}
	wantCalls := [][]string{
		{"fd", "--follow", "--hidden", "--exclude", "{.git,node_modules,Library}", ".", "--base-directory", home, "--max-depth", "3"},
		{"fd", "--follow", "--hidden", "--exclude", "{.git,node_modules,Library}", ".", "--base-directory", home, "--min-depth", "4"},
	}
	if diff := cmp.Diff(wantCalls, runner.Calls); diff != "" {
		t.Fatal(diff)
	}
}

func TestScanCloseStopsActivePass(t *testing.T) {
	for _, tt := range []struct {
		name            string
		afterTransition bool
	}{{"shallow", false}, {"deep", true}} {
		t.Run(tt.name, func(t *testing.T) {
			shallow := &trackedReader{Reader: strings.NewReader("")}
			deep := &trackedReader{Reader: strings.NewReader("deep\n")}
			startedDeep := false
			input := &scan{ReadCloser: shallow, next: func() (io.ReadCloser, error) {
				startedDeep = true
				if !shallow.closed {
					t.Fatal("shallow scan was not closed before starting the deep scan")
				}
				return deep, nil
			}}
			if tt.afterTransition {
				if _, err := io.ReadFull(input, make([]byte, 1)); err != nil {
					t.Fatal(err)
				}
			}
			if err := input.Close(); err != nil {
				t.Fatal(err)
			}
			if _, err := input.Read(make([]byte, 1)); !errors.Is(err, os.ErrClosed) {
				t.Fatalf("Read after Close = %v, want closed stream", err)
			}
			if startedDeep != tt.afterTransition || deep.closed != tt.afterTransition || !shallow.closed {
				t.Fatalf("started deep = %v, closed shallow = %v, closed deep = %v", startedDeep, shallow.closed, deep.closed)
			}
		})
	}
}

type trackedReader struct {
	io.Reader
	closed bool
}

func (r *trackedReader) Read(p []byte) (int, error) {
	if r.closed {
		return 0, os.ErrClosed
	}
	return r.Reader.Read(p)
}

func (r *trackedReader) Close() error {
	r.closed = true
	return nil
}

func TestScanReportsErrors(t *testing.T) {
	want := errors.New("scan failed")
	for _, tt := range []struct {
		name string
		read io.Reader
		next func() (io.ReadCloser, error)
	}{
		{"shallow read", iotest.ErrReader(want), func() (io.ReadCloser, error) {
			t.Fatal("deep scan started after shallow scan failed")
			return nil, nil
		}},
		{"deep start", strings.NewReader(""), func() (io.ReadCloser, error) { return nil, want }},
		{"deep read", strings.NewReader(""), func() (io.ReadCloser, error) { return io.NopCloser(iotest.ErrReader(want)), nil }},
	} {
		t.Run(tt.name, func(t *testing.T) {
			input := &scan{ReadCloser: io.NopCloser(tt.read), next: tt.next}
			defer input.Close()
			if _, err := io.ReadAll(input); !errors.Is(err, want) {
				t.Fatalf("ReadAll error = %v, want %v", err, want)
			}
		})
	}
}

func TestFdDepthPartition(t *testing.T) {
	if _, err := exec.LookPath("fd"); err != nil {
		t.Skip("fd is not installed")
	}
	home := t.TempDir()
	t.Setenv("HOME", home)
	for _, path := range []string{"Code/project/src", ".config", "node_modules/package", "Library/cache", "ignored"} {
		if err := os.MkdirAll(filepath.Join(home, path), 0755); err != nil {
			t.Fatal(err)
		}
	}
	for _, path := range []string{"shallow.txt", "Code/project/README.md", "Code/project/src/deep.go", ".config/settings", "node_modules/package/hidden.js", "Library/cache/hidden", "ignored/hidden"} {
		if err := os.WriteFile(filepath.Join(home, path), nil, 0644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(home, ".fdignore"), []byte("ignored/\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(home, "Code"), filepath.Join(home, "linked")); err != nil {
		t.Fatal(err)
	}
	input, err := (Cmd{Runner: runner.CmdRunner{}}).Fd()
	if err != nil {
		t.Fatal(err)
	}
	defer input.Close()
	output, err := io.ReadAll(input)
	if err != nil {
		t.Fatal(err)
	}
	paths := strings.Split(strings.TrimSuffix(string(output), "\n"), "\n")
	deepStarted := false
	for _, path := range paths {
		depth := len(strings.Split(strings.TrimSuffix(path, "/"), "/"))
		if depth > 3 {
			deepStarted = true
		} else if deepStarted {
			t.Fatalf("shallow path %q arrived after deep results: %v", path, paths)
		}
	}
	slices.Sort(paths)
	want := []string{
		".config/", ".config/settings", ".fdignore", "shallow.txt",
		"Code/", "Code/project/", "Code/project/README.md", "Code/project/src/", "Code/project/src/deep.go",
		"linked/", "linked/project/", "linked/project/README.md", "linked/project/src/", "linked/project/src/deep.go",
	}
	slices.Sort(want)
	if diff := cmp.Diff(want, paths); diff != "" {
		t.Fatal(diff)
	}
}
