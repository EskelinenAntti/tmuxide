package ide

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/shell"
)

type PathMock struct{}

func (PathMock) Contains(program string) bool {
	return true
}

type GitMock struct {
	repository string
	err        error
}

func (git GitMock) RevParse(cwd string) (string, error) {
	return git.repository, git.err
}

var shellMock = shell.Shell{
	Git: GitMock{
		repository: "",
		err:        errors.New("not inside git repository"),
	},
	Tmux: nil,
	Path: PathMock{},
}

func TestProjectForDirectory(t *testing.T) {
	var dir = t.TempDir()

	project, err := ProjectFor(dir, shellMock)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	if got, want := project.Root, dir; got != want {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestProjectForFile(t *testing.T) {
	dir := t.TempDir()
	file := dir + "/file.txt"
	os.WriteFile(file, []byte{}, 0644)

	project, err := ProjectFor(file, shellMock)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	if got, want := project.Root, dir; got != want {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestProjectForRepository(t *testing.T) {
	repository := t.TempDir()
	dir := filepath.Join(repository, "path/to/dir/in/repository")

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		t.Fatalf("err=%v", err)
	}

	shellMock := shell.Shell{
		Git: GitMock{
			repository: repository,
			err:        nil,
		},
		Path: PathMock{},
	}

	project, err := ProjectFor(dir, shellMock)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	if got, want := project.Root, repository; got != want {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestProjectForInvalidFile(t *testing.T) {
	dir := t.TempDir()
	file := dir + "/does-not-exist.txt"

	_, err := ProjectFor(file, shellMock)

	var pathError *os.PathError
	if got, want := err, &pathError; !errors.As(got, want) {
		t.Fatalf("got=%T, want=%T", got, want)
	}

	if got, want := err.Error(), fmt.Sprintf("stat %s: no such file or directory", file); got != want {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}
