package ranked

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/eskelinenantti/herdr-zipln/internal/history"
	"github.com/google/go-cmp/cmp"
)

func TestRankedFormatsExistingPaths(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_STATE_HOME", t.TempDir())
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
	for _, path := range []string{file, file, file, dir, dir, filepath.Join(home, "missing"), newline, external} {
		history.Record(path)
	}
	output, err := io.ReadAll(History())
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff("project/file.txt\nproject/\n"+external+"\n", string(output)); diff != "" {
		t.Fatal(diff)
	}
}

func TestRankedWithoutHistory(t *testing.T) {
	t.Setenv("XDG_STATE_HOME", t.TempDir())
	output, err := io.ReadAll(History())
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff("", string(output)); diff != "" {
		t.Fatal(diff)
	}
}
