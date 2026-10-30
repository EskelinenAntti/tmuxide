package history

import (
	"path/filepath"
	"testing"
)

func TestQuerySelections(t *testing.T) {
	t.Setenv("XDG_STATE_HOME", t.TempDir())
	t.Chdir(t.TempDir())
	RecordQuery("proj", "first")
	RecordQuery("other", "second")
	RecordQuery("proj", "third")
	for query, path := range map[string]string{"proj": "third", "other": "second", "unknown": "", "": ""} {
		want := ""
		if path != "" {
			want, _ = filepath.Abs(path)
		}
		if got := Preferred(query); got != want {
			t.Fatalf("Preferred(%q) = %q, want %q", query, got, want)
		}
	}
	RecordQuery("", "ignored")
	RecordQuery("proj", "")
	if Preferred("") != "" || Preferred("proj") == "" {
		t.Fatal("empty queries and selections must not be recorded")
	}
}

func TestMostSelectedForEachQuery(t *testing.T) {
	t.Setenv("XDG_STATE_HOME", t.TempDir())
	dir := t.TempDir()
	first := filepath.Join(dir, "first")
	second := filepath.Join(dir, "second")
	for _, path := range []string{first, first, second} {
		RecordQuery("proj", path)
	}
	RecordQuery("other", second)
	if got := Preferred("proj"); got != first {
		t.Fatalf("most selected = %q, want %q", got, first)
	}
	if got := Preferred("other"); got != second {
		t.Fatalf("other query = %q, want %q", got, second)
	}
	RecordQuery("proj", second)
	if got := Preferred("proj"); got != second {
		t.Fatalf("most recent tied selection = %q, want %q", got, second)
	}
}
