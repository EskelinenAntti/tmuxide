package fd

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPathWriter(t *testing.T) {
	var paths []string
	output := &pathWriter{onPath: func(path string) error {
		paths = append(paths, path)
		return nil
	}}
	for _, chunk := range []string{"project/\npro", "ject/file.txt\nother file", ".txt\n"} {
		if _, err := output.Write([]byte(chunk)); err != nil {
			t.Fatal(err)
		}
	}
	if diff := cmp.Diff([]string{"project/", "project/file.txt", "other file.txt"}, paths); diff != "" {
		t.Fatal(diff)
	}
}

func TestPathWriterReturnsError(t *testing.T) {
	failure := errors.New("output closed")
	output := &pathWriter{onPath: func(string) error { return failure }}
	if _, err := output.Write([]byte("file.txt\n")); !errors.Is(err, failure) {
		t.Fatalf("Write = %v, want %v", err, failure)
	}
}
