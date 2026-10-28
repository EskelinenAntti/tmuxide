package pipe

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"syscall"
	"testing"
	"testing/iotest"

	"github.com/google/go-cmp/cmp"
)

func TestDeduplicatePathsPreservesFirstOccurrence(t *testing.T) {
	stdin := io.MultiReader(
		strings.NewReader("project/\nproject/file.txt\n"),
		strings.NewReader("./project\n./project/file.txt\n\nother file.txt\n"),
	)
	output, err := io.ReadAll(iotest.OneByteReader(DeduplicatePaths(stdin)))
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff("project/\nproject/file.txt\nother file.txt\n", string(output)); diff != "" {
		t.Fatal(diff)
	}
}

func TestDeduplicatePathsPropagatesReadErrors(t *testing.T) {
	want := errors.New("read failed")
	stdin := io.MultiReader(strings.NewReader("project/\n"), iotest.ErrReader(want))
	output, err := io.ReadAll(DeduplicatePaths(stdin))
	if !errors.Is(err, want) {
		t.Fatalf("ReadAll error = %v, want %v", err, want)
	}
	if diff := cmp.Diff("project/\n", string(output)); diff != "" {
		t.Fatal(diff)
	}
}

func TestDeduplicatePathsRejectsOversizedLines(t *testing.T) {
	_, err := io.ReadAll(DeduplicatePaths(strings.NewReader(strings.Repeat("x", 70000))))
	if err == nil {
		t.Fatal("expected scanner error")
	}
}

type failingWriter struct{ err error }

func (w failingWriter) Write([]byte) (int, error) { return 0, w.err }

func TestPipe(t *testing.T) {
	t.Run("copies stdin to stdout", func(t *testing.T) {
		var stdout bytes.Buffer
		if err := Pipe(strings.NewReader("project/\n"), &stdout); err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff("project/\n", stdout.String()); diff != "" {
			t.Fatal(diff)
		}
	})
	t.Run("closed downstream pipe is normal completion", func(t *testing.T) {
		if err := Pipe(strings.NewReader("project/\n"), failingWriter{syscall.EPIPE}); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("other write failures are returned", func(t *testing.T) {
		want := errors.New("write failed")
		if err := Pipe(strings.NewReader("project/\n"), failingWriter{want}); !errors.Is(err, want) {
			t.Fatalf("Pipe error = %v, want %v", err, want)
		}
	})
}
