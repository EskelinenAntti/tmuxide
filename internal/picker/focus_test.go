package picker

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/iotest"
)

func TestPreferredPosition(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	preferred := filepath.Join(home, "project", "child")
	for _, tt := range []struct {
		name    string
		matches string
		want    int
	}{
		{"preferred match", "project/\nproject/other/\nproject/child/\n", 3},
		{"different ranking", "project/child/\nproject/\n", 1},
		{"absolute path", "project/\n" + preferred + "/\n", 2},
		{"normalized path", "project/\n./project/child/\n", 2},
		{"preferred path has not arrived", "project/\nother/\n", 0},
		{"no matches yet", "", 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PreferredPosition(preferred, strings.NewReader(tt.matches))
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
	got, err := PreferredPosition(filepath.Join(os.Getenv("HOME"), "project "), strings.NewReader("project/\nproject /\n"))
	if err != nil || got != 2 {
		t.Fatalf("PreferredPosition = %d, %v; want 2, nil", got, err)
	}
}

func TestPreferredPositionReportsReadErrors(t *testing.T) {
	want := errors.New("read failed")
	_, err := PreferredPosition(filepath.Join(t.TempDir(), "project"), io.MultiReader(strings.NewReader("other/\n"), iotest.ErrReader(want)))
	if !errors.Is(err, want) {
		t.Fatalf("PreferredPosition error = %v, want %v", err, want)
	}
}
