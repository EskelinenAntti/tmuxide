package picker

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

func TestCandidates(t *testing.T) {
	var buffer bytes.Buffer
	output := &candidates{output: &buffer, seen: make(map[string]bool)}
	if err := output.add("project/"); err != nil {
		t.Fatal(err)
	}
	for _, chunk := range []string{"./project\npro", "ject/file.txt\n\nother/\n", "project/file.txt\nlast file \n"} {
		if n, err := output.Write([]byte(chunk)); err != nil || n != len(chunk) {
			t.Fatalf("Write = %d, %v", n, err)
		}
	}
	want := "project/\nproject/file.txt\nother/\nlast file \n"
	if got := buffer.String(); got != want {
		t.Fatalf("candidates = %q, want %q", got, want)
	}
}

type brokenWriter struct{}

func (brokenWriter) Write([]byte) (int, error) { return 0, io.ErrClosedPipe }

func TestCandidatesWriteError(t *testing.T) {
	output := &candidates{output: brokenWriter{}, seen: make(map[string]bool)}
	if _, err := output.Write([]byte("file\n")); !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("Write = %v, want closed pipe", err)
	}
}
