package picker

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/history"
	"github.com/eskelinenantti/tmuxide/internal/shell/fd"
	"github.com/eskelinenantti/tmuxide/internal/shell/fzf"
	"github.com/eskelinenantti/tmuxide/internal/shell/runner"
	"github.com/eskelinenantti/tmuxide/internal/test/mock"
	"github.com/eskelinenantti/tmuxide/internal/test/spy"
	"github.com/google/go-cmp/cmp"
)

func TestSearchRanksHistoryBeforeScanningFoldersAndHome(t *testing.T) {
	dir := t.TempDir()
	home := filepath.Join(dir, "home")
	t.Setenv("HOME", home)
	t.Setenv("XDG_STATE_HOME", t.TempDir())
	for _, dir := range []string{"popular", "other folder"} {
		if err := os.MkdirAll(filepath.Join(home, dir), 0755); err != nil {
			t.Fatal(err)
		}
	}
	file := filepath.Join(home, "popular", "file.txt")
	if err := os.WriteFile(file, nil, 0644); err != nil {
		t.Fatal(err)
	}
	outside := filepath.Join(dir, "outside.txt")
	if err := os.WriteFile(outside, nil, 0644); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{
		file, file, file,
		filepath.Join(home, "popular"), filepath.Join(home, "popular"),
		filepath.Join(home, "other folder"), outside, filepath.Join(home, "missing"),
	} {
		if err := history.Record(path); err != nil {
			t.Fatal(err)
		}
	}
	r := &spy.SpyRunner{Responses: []spy.Response{
		{OnRun: mock.WriteToStdout("popular/file.txt\npopular/child.txt\n")},
		{OnRun: mock.WriteToStdout("other folder/child.txt\n")},
		{OnRun: mock.WriteToStdout("popular/\npopular/file.txt\npopular/child.txt\nother folder/\nnew.txt\n")},
	}}
	var buffer bytes.Buffer
	if err := search(fd.Cmd{Runner: r}, &candidates{output: &buffer, seen: make(map[string]bool)}); err != nil {
		t.Fatal(err)
	}
	wantCalls := [][]string{
		{"fd", "--follow", "--hidden", "--exclude", "{.git,node_modules,Library}", ".", "--base-directory", home, "--", "popular"},
		{"fd", "--follow", "--hidden", "--exclude", "{.git,node_modules,Library}", ".", "--base-directory", home, "--", "other folder"},
		{"fd", "--follow", "--hidden", "--exclude", "{.git,node_modules,Library}", ".", "--base-directory", home, "--"},
	}
	if diff := cmp.Diff(wantCalls, r.Calls); diff != "" {
		t.Fatal(diff)
	}
	want := "popular/file.txt\npopular/\nother folder/\n" + outside + "\npopular/child.txt\nother folder/child.txt\nnew.txt\n"
	if got := buffer.String(); got != want {
		t.Fatalf("results = %q, want %q", got, want)
	}
}

func TestSearchWithoutUsableHistory(t *testing.T) {
	for _, data := range []string{"", "not json"} {
		t.Run(data, func(t *testing.T) {
			home := t.TempDir()
			t.Setenv("HOME", home)
			t.Setenv("XDG_STATE_HOME", t.TempDir())
			if data != "" {
				path := filepath.Join(os.Getenv("XDG_STATE_HOME"), "tmuxide", "history.json")
				if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, []byte(data), 0600); err != nil {
					t.Fatal(err)
				}
			}
			r := &spy.SpyRunner{Responses: []spy.Response{{OnRun: mock.WriteToStdout("file\nfolder/\n")}}}
			var buffer bytes.Buffer
			if err := search(fd.Cmd{Runner: r}, &candidates{output: &buffer, seen: make(map[string]bool)}); err != nil {
				t.Fatal(err)
			}
			if got := buffer.String(); got != "file\nfolder/\n" {
				t.Fatalf("results = %q", got)
			}
			if len(r.Calls) != 1 {
				t.Fatalf("expected one home scan, got %v", r.Calls)
			}
		})
	}
}

type promptRunner struct {
	selected string
	writeErr error
	closeErr error
	scanErr  error
	closed   bool
}

func (r *promptRunner) Run(cmd *exec.Cmd) error {
	if r.scanErr != nil {
		return r.scanErr
	}
	_, err := io.WriteString(cmd.Stdout, "file\n")
	return err
}

func (r *promptRunner) Start(cmd *exec.Cmd) (runner.WriteCloser, error) {
	if _, err := io.WriteString(cmd.Stdout, r.selected); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *promptRunner) Write(p []byte) (int, error) {
	if r.writeErr != nil {
		return 0, r.writeErr
	}
	return len(p), nil
}

func (r *promptRunner) Close() error {
	r.closed = true
	return r.closeErr
}

func TestPrompt(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_STATE_HOME", t.TempDir())
	failure := errors.New("scan failed")
	for _, tt := range []struct {
		name   string
		runner promptRunner
		want   string
		err    error
	}{
		{name: "relative path", runner: promptRunner{selected: "folder/file\n"}, want: filepath.Join(home, "folder/file")},
		{name: "absolute path", runner: promptRunner{selected: "/outside/file\n"}, want: "/outside/file"},
		{name: "spaces preserved", runner: promptRunner{selected: " file \n"}, want: filepath.Join(home, " file ")},
		{name: "empty selection"},
		{name: "early selection", runner: promptRunner{selected: "file\n", writeErr: io.ErrClosedPipe}, want: filepath.Join(home, "file")},
		{name: "selection despite scan error", runner: promptRunner{selected: "file\n", scanErr: failure}, want: filepath.Join(home, "file")},
		{name: "scan failure", runner: promptRunner{scanErr: failure}, err: failure},
		{name: "cancelled", runner: promptRunner{closeErr: exec.Command("sh", "-c", "exit 130").Run()}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			r := &tt.runner
			got, err := Prompt(fd.Cmd{Runner: r}, fzf.Cmd{Runner: r})
			if got != tt.want || !errors.Is(err, tt.err) {
				t.Fatalf("Prompt = %q, %v; want %q, %v", got, err, tt.want, tt.err)
			}
			if !r.closed {
				t.Fatal("fzf input was not closed")
			}
		})
	}
}
