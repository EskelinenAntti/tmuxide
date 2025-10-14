package project

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

func TestRootDirectory(t *testing.T) {
	var dir = t.TempDir()
	root, err := Root(dir)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	if got, want := root, dir; got != want {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestRootFile(t *testing.T) {
	dir := t.TempDir()
	file := dir + "/file.txt"
	os.WriteFile(file, []byte{}, 0644)

	root, err := Root(file)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	if got, want := root, dir; got != want {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestRootInvalidFile(t *testing.T) {
	dir := t.TempDir()
	file := dir + "/does-not-exist.txt"

	_, err := Root(file)

	var pathError *os.PathError
	if got, want := err, &pathError; !errors.As(got, want) {
		t.Fatalf("got=%T, want=%T", got, want)
	}

	if got, want := err.Error(), fmt.Sprintf("stat %s: no such file or directory", file); got != want {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}
