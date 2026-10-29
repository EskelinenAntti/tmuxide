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
