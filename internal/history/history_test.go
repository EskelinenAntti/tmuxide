package history

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFilePath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	for _, stateHome := range []string{"", "relative/state", t.TempDir()} {
		t.Run(stateHome, func(t *testing.T) {
			t.Setenv("XDG_STATE_HOME", stateHome)
			base := stateHome
			if !filepath.IsAbs(base) {
				base = filepath.Join(home, ".local", "state")
			}
			want := filepath.Join(base, "tmuxide", "history.json")
			if got := filePath(); got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
		})
	}
}

func TestRecordAndRank(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	t.Setenv("XDG_STATE_HOME", filepath.Join(dir, "state"))
	paths, err := Ranked()
	if err != nil || len(paths) != 0 {
		t.Fatalf("empty history = %v, %v", paths, err)
	}
	for _, path := range []string{"b", "./a", "c", "./b/../c", "c"} {
		if err := Record(path); err != nil {
			t.Fatal(err)
		}
	}
	paths, err = Ranked()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		filepath.Join(dir, "c"),
		filepath.Join(dir, "a"),
		filepath.Join(dir, "b"),
	}
	if diff := cmp.Diff(want, paths); diff != "" {
		t.Fatal(diff)
	}
	counts, err := load()
	if err != nil {
		t.Fatal(err)
	}
	if got := counts[filepath.Join(dir, "c")]; got != 3 {
		t.Fatalf("selection count = %d, want 3", got)
	}
}

func TestInvalidHistoryIsNotOverwritten(t *testing.T) {
	t.Setenv("XDG_STATE_HOME", t.TempDir())
	if err := os.MkdirAll(filepath.Dir(filePath()), 0700); err != nil {
		t.Fatal(err)
	}
	const invalid = "not json"
	if err := os.WriteFile(filePath(), []byte(invalid), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := Ranked(); err == nil {
		t.Fatal("expected invalid history error")
	}
	if err := Record("project"); err == nil {
		t.Fatal("expected invalid history error")
	}
	data, err := os.ReadFile(filePath())
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != invalid {
		t.Fatalf("history was overwritten: %q", data)
	}
}
