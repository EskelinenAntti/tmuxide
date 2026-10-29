package history

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRecordRanksByFrequency(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_STATE_HOME", t.TempDir())
	t.Chdir(dir)
	for _, path := range []string{"b", "a", "c", "./c", "c"} {
		Record(path)
	}
	want := []string{filepath.Join(dir, "c"), filepath.Join(dir, "a"), filepath.Join(dir, "b")}
	if diff := cmp.Diff(want, Ranked()); diff != "" {
		t.Fatal(diff)
	}
}

func TestInvalidHistoryCanBeReplaced(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_STATE_HOME", dir)
	if err := os.Mkdir(filepath.Dir(filePath()), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filePath(), []byte("invalid"), 0600); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "project")
	Record(path)
	if diff := cmp.Diff([]string{path}, Ranked()); diff != "" {
		t.Fatal(diff)
	}
}
