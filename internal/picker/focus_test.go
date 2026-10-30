package picker

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/eskelinenantti/herdr-zipln/internal/history"
)

func TestPreferredPosition(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_STATE_HOME", t.TempDir())
	preferred := filepath.Join(home, "project", "child")
	history.RecordQuery("proj", preferred)
	history.RecordQuery("proj", preferred)
	history.RecordQuery("proj", filepath.Join(home, "project"))
	history.RecordQuery("other", filepath.Join(home, "other"))
	for _, tt := range []struct {
		name    string
		query   string
		matches string
		want    int
	}{
		{"most frequent", "proj", "project/\nproject/other/\nproject/child/\n", 3},
		{"different ranking", "proj", "project/child/\nproject/\n", 1},
		{"absolute path", "proj", "project/\n" + preferred + "/\n", 2},
		{"normalized path", "proj", "project/\n./project/child/\n", 2},
		{"unrelated query", "project", "project/\nproject/child/\n", 1},
		{"different preference", "other", "project/\nother/\n", 2},
		{"missing preference", "proj", "project/\nother/\n", 1},
		{"no matches", "proj", "", 1},
		{"empty query", "", "project/\nproject/child/\n", 1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PreferredPosition(tt.query, strings.NewReader(tt.matches))
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Fatalf("PreferredPosition = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestPreferredPositionPreservesSpaces(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("XDG_STATE_HOME", t.TempDir())
	history.RecordQuery("space", filepath.Join(os.Getenv("HOME"), "project "))
	got, err := PreferredPosition("space", strings.NewReader("project/\nproject /\n"))
	if err != nil || got != 2 {
		t.Fatalf("PreferredPosition = %d, %v; want 2, nil", got, err)
	}
}

func TestPreferredPositionReportsReadErrors(t *testing.T) {
	t.Setenv("XDG_STATE_HOME", t.TempDir())
	history.RecordQuery("proj", filepath.Join(t.TempDir(), "project"))
	want := errors.New("read failed")
	_, err := PreferredPosition("proj", io.MultiReader(strings.NewReader("other/\n"), iotest.ErrReader(want)))
	if !errors.Is(err, want) {
		t.Fatalf("PreferredPosition error = %v, want %v", err, want)
	}
}
