package history

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRecordRanksByFrequency(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_STATE_HOME", t.TempDir())
	t.Chdir(dir)
	for _, path := range []string{"b", "a", "c", "./c", "c"} {
		if err := Record(path); err != nil {
			t.Fatal(err)
		}
	}
	paths, err := Ranked()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{filepath.Join(dir, "c"), filepath.Join(dir, "a"), filepath.Join(dir, "b")}
	if diff := cmp.Diff(want, paths); diff != "" {
		t.Fatal(diff)
	}
}
