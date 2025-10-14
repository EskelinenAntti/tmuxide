package project

import (
	"os"
	"testing"
)

func TestDirectory(t *testing.T) {
	var want = t.TempDir()
	got, err := Root(want)

	if err != nil {
		t.Fatalf("Root(%s): err = %v", want, err)
	}

	if got != want {
		t.Fatalf("Root(%s) = %s, want %s", want, got, want)
	}
}

func TestFile(t *testing.T) {
	want := t.TempDir()
	file := want + "/file.txt"
	os.WriteFile(file, []byte{}, 0644)

	got, err := Root(file)

	if err != nil {
		t.Fatalf("Root(%s): err = %v", want, err)
	}

	if got != want {
		t.Fatalf("Root(%s) = %s, want %s", file, got, want)
	}
}
