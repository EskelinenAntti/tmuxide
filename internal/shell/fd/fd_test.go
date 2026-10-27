package fd

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWriteSkipsDuplicates(t *testing.T) {
	var output bytes.Buffer
	seen := make(map[string]bool)
	for _, path := range []string{"project/", "project/file.txt", "./project", "./project/file.txt", "other file.txt"} {
		if err := write(&output, path, seen); err != nil {
			t.Fatal(err)
		}
	}
	want := "project/\nproject/file.txt\nother file.txt\n"
	if diff := cmp.Diff(want, output.String()); diff != "" {
		t.Fatal(diff)
	}
}
