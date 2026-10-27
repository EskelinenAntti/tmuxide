package picker

import (
	"bytes"
	"errors"
	"syscall"
	"testing"
)

func TestCandidates(t *testing.T) {
	var buffer bytes.Buffer
	output := &candidates{output: &buffer, seen: make(map[string]bool)}
	if err := output.add("project/"); err != nil {
		t.Fatal(err)
	}
	for _, chunk := range []string{"./project\npro", "ject/file.txt\n\nother/\n", "project/file.txt\nlast file "} {
		if n, err := output.Write([]byte(chunk)); err != nil || n != len(chunk) {
			t.Fatalf("Write = %d, %v", n, err)
		}
	}
	if err := output.flush(); err != nil {
		t.Fatal(err)
	}
	want := "project/\nproject/file.txt\nother/\nlast file \n"
	if got := buffer.String(); got != want {
		t.Fatalf("candidates = %q, want %q", got, want)
	}
}

type brokenWriter struct{}

func (brokenWriter) Write([]byte) (int, error) { return 0, syscall.EPIPE }

func TestCandidatesStopOnClosedPipe(t *testing.T) {
	output := &candidates{output: brokenWriter{}, seen: make(map[string]bool)}
	if _, err := output.Write([]byte("file\n")); !errors.Is(err, syscall.EPIPE) {
		t.Fatalf("Write = %v, want broken pipe", err)
	}
	if err := output.add("another"); !errors.Is(err, syscall.EPIPE) {
		t.Fatalf("add = %v, want broken pipe", err)
	}
}
