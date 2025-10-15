package project

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

type RepositoryResolverStub struct {
	root string
	err  error
}

func (resolver RepositoryResolverStub) Root(path string) (string, error) {
	return resolver.root, resolver.err
}

var notInGitRepo = RepositoryResolverStub{
	root: "",
	err:  errors.New("not a git repo"),
}

func TestRootDirectory(t *testing.T) {
	var dir = t.TempDir()

	root, err := Root(dir, notInGitRepo)

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

	root, err := Root(file, notInGitRepo)

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

	_, err := Root(file, notInGitRepo)

	var pathError *os.PathError
	if got, want := err, &pathError; !errors.As(got, want) {
		t.Fatalf("got=%T, want=%T", got, want)
	}

	if got, want := err.Error(), fmt.Sprintf("stat %s: no such file or directory", file); got != want {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestRootRepository(t *testing.T) {
	repository := t.TempDir()
	rootResolver := RepositoryResolverStub{
		root: repository,
		err:  nil,
	}

	dir := filepath.Join(repository, "path/to/dir/in/repository")

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		t.Fatalf("err=%v", err)
	}

	root, err := Root(dir, rootResolver)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	if got, want := root, repository; got != want {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}
