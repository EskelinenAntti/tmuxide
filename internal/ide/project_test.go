package ide

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

type RepositoryStub struct {
	root string
	err  error
}

func (resolver RepositoryStub) Root(path string) (string, error) {
	return resolver.root, resolver.err
}

var notInGitRepo = RepositoryStub{
	root: "",
	err:  errors.New("not a git repo"),
}

func TestRootDirectory(t *testing.T) {
	var dir = t.TempDir()

	project, err := ProjectFor(dir, notInGitRepo)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	if got, want := project.Root, dir; got != want {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestRootFile(t *testing.T) {
	dir := t.TempDir()
	file := dir + "/file.txt"
	os.WriteFile(file, []byte{}, 0644)

	project, err := ProjectFor(file, notInGitRepo)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	if got, want := project.Root, dir; got != want {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestRootRepository(t *testing.T) {
	repositoryPath := t.TempDir()
	repository := RepositoryStub{
		root: repositoryPath,
		err:  nil,
	}

	dir := filepath.Join(repositoryPath, "path/to/dir/in/repository")

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		t.Fatalf("err=%v", err)
	}

	project, err := ProjectFor(dir, repository)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	if got, want := project.Root, repository.root; got != want {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestRootInvalidFile(t *testing.T) {
	dir := t.TempDir()
	file := dir + "/does-not-exist.txt"

	_, err := ProjectFor(file, notInGitRepo)

	var pathError *os.PathError
	if got, want := err, &pathError; !errors.As(got, want) {
		t.Fatalf("got=%T, want=%T", got, want)
	}

	if got, want := err.Error(), fmt.Sprintf("stat %s: no such file or directory", file); got != want {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}
