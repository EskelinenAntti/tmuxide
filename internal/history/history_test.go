package history

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDefault(t *testing.T) {
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
			if got := Default().path; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
		})
	}
}

func TestRecordAndRank(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	store := Store{path: filepath.Join(dir, "state", "history.json")}
	entries, err := store.Ranked()
	if err != nil || len(entries) != 0 {
		t.Fatalf("empty history = %v, %v", entries, err)
	}
	for _, path := range []string{"b", "./a", "c", "./b/../c", "c"} {
		if err := store.Record(path); err != nil {
			t.Fatal(err)
		}
	}
	// Reopen the store to verify that ranking is persisted across runs.
	entries, err = (Store{path: store.path}).Ranked()
	if err != nil {
		t.Fatal(err)
	}
	want := []Entry{
		{Path: filepath.Join(dir, "c"), Count: 3},
		{Path: filepath.Join(dir, "a"), Count: 1},
		{Path: filepath.Join(dir, "b"), Count: 1},
	}
	if diff := cmp.Diff(want, entries); diff != "" {
		t.Fatal(diff)
	}
	for path, mode := range map[string]os.FileMode{
		filepath.Dir(store.path): 0700,
		store.path:               0600,
		store.path + ".lock":     0600,
	} {
		info, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}
		if got := info.Mode().Perm(); got != mode {
			t.Fatalf("permissions of %s = %o, want %o", path, got, mode)
		}
	}
}

func TestConcurrentRecords(t *testing.T) {
	store := Store{path: filepath.Join(t.TempDir(), "history.json")}
	path := filepath.Join(t.TempDir(), "project")
	const selections = 20
	var wg sync.WaitGroup
	for range selections {
		wg.Go(func() {
			if err := store.Record(path); err != nil {
				t.Error(err)
			}
		})
	}
	wg.Wait()
	entries, err := store.Ranked()
	if err != nil {
		t.Fatal(err)
	}
	want := []Entry{{Path: path, Count: selections}}
	if diff := cmp.Diff(want, entries); diff != "" {
		t.Fatal(diff)
	}
}

func TestInvalidHistoryIsNotOverwritten(t *testing.T) {
	store := Store{path: filepath.Join(t.TempDir(), "history.json")}
	const invalid = "not json"
	if err := os.WriteFile(store.path, []byte(invalid), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Ranked(); err == nil {
		t.Fatal("expected invalid history error")
	}
	if err := store.Record("project"); err == nil {
		t.Fatal("expected invalid history error")
	}
	data, err := os.ReadFile(store.path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != invalid {
		t.Fatalf("history was overwritten: %q", data)
	}
}
