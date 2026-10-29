package pipe

import (
	"errors"
	"io"
	"strings"
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
