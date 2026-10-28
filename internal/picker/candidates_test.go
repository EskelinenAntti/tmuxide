package picker

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWriteRankedFormatsExistingPaths(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	dir := filepath.Join(home, "project")
	if err := os.Mkdir(dir, 0755); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(file, nil, 0644); err != nil {
		t.Fatal(err)
	}
	external := filepath.Join(t.TempDir(), "external.txt")
	if err := os.WriteFile(external, nil, 0644); err != nil {
		t.Fatal(err)
	}
	newline := filepath.Join(home, "new\nline")
	if err := os.WriteFile(newline, nil, 0644); err != nil {
		t.Fatal(err)
	}
	paths := []string{file, dir, filepath.Join(home, "missing"), newline, external}
	var output bytes.Buffer
	if _, err := writeRanked(&output, paths); err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff("project/file.txt\nproject/\n"+external+"\n", output.String()); diff != "" {
		t.Fatal(diff)
	}
}

func TestWritePathSkipsDuplicates(t *testing.T) {
	var output bytes.Buffer
	seen := make(map[string]bool)
	for _, path := range []string{"project/", "project/file.txt", "./project", "./project/file.txt", "other file.txt"} {
		if err := writePath(&output, path, seen); err != nil {
			t.Fatal(err)
		}
	}
	want := "project/\nproject/file.txt\nother file.txt\n"
	if diff := cmp.Diff(want, output.String()); diff != "" {
		t.Fatal(diff)
	}
}
