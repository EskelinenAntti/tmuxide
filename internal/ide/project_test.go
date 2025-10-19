package ide

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/shell"
	"github.com/eskelinenantti/tmuxide/internal/tmux"
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

var testEditor string = "editor"
var errNotInRepository = errors.New("not in repository")

func shellMock(repository string) shell.Shell {
	var err error
	if repository == "" {
		err = errNotInRepository
	}

	return shell.Shell{
		Git: GitMock{
			repository: repository,
			err:        err,
		},
		Tmux: nil,
		Path: PathMock{},
	}
}

func wantWindows(target string, lazygit bool) []tmux.WindowCommand {
	windows := []tmux.WindowCommand{
		{
			Cmd:  testEditor,
			Args: []string{target},
		}}
	if lazygit {
		windows = append(windows, tmux.WindowCommand{
			Cmd:  "lazygit",
			Args: nil,
		})
	}
	return windows
}

func TestProjectForDirectory(t *testing.T) {
	t.Setenv("EDITOR", testEditor)
	var dir = t.TempDir()

	project, err := ProjectFor(dir, shellMock(""))

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	if got, want := project.Root, dir; got != want {
		t.Fatalf("got=%v, want=%v", got, want)
	}

	if got, want := project.Name, filepath.Base(dir)+"-"; !strings.HasPrefix(got, want) {
		t.Fatalf("%v did not start with %v", got, want)
	}

	if got, want := project.Windows, wantWindows(dir, false); !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}

}

func TestProjectForFile(t *testing.T) {
	t.Setenv("EDITOR", testEditor)

	dir := t.TempDir()
	file := dir + "/file.txt"
	os.WriteFile(file, []byte{}, 0644)

	project, err := ProjectFor(file, shellMock(""))

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	if got, want := project.Root, dir; got != want {
		t.Fatalf("got=%v, want=%v", got, want)
	}

	if got, want := project.Name, "file_txt-"; !strings.HasPrefix(got, want) {
		t.Fatalf("%v did not start with %v", got, want)
	}

	if got, want := project.Windows, wantWindows(file, false); !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestProjectForDirectoryInRepository(t *testing.T) {
	t.Setenv("EDITOR", testEditor)

	repository := t.TempDir()
	dir := filepath.Join(repository, "path/to/dir/in/repository")

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		t.Fatalf("err=%v", err)
	}

	project, err := ProjectFor(dir, shellMock(repository))

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	if got, want := project.Root, repository; got != want {
		t.Fatalf("got=%v, want=%v", got, want)
	}

	if got, want := project.Name, filepath.Base(dir)+"-"; !strings.HasPrefix(got, want) {
		t.Fatalf("%v did not start with %v", got, want)
	}

	if got, want := project.Windows, wantWindows(dir, true); !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestProjectForFileInRepository(t *testing.T) {
	t.Setenv("EDITOR", testEditor)

	repository := t.TempDir()
	file := repository + "/file.txt"
	os.WriteFile(file, []byte{}, 0644)

	project, err := ProjectFor(file, shellMock(repository))

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	if got, want := project.Root, repository; got != want {
		t.Fatalf("got=%v, want=%v", got, want)
	}

	if got, want := project.Name, "file_txt-"; !strings.HasPrefix(got, want) {
		t.Fatalf("%v did not start with %v", got, want)
	}

	if got, want := project.Windows, wantWindows(file, true); !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestProjectForInvalidFile(t *testing.T) {
	t.Setenv("EDITOR", testEditor)

	dir := t.TempDir()
	file := dir + "/does-not-exist.txt"

	_, err := ProjectFor(file, shellMock(""))

	var pathError *os.PathError
	if got, want := err, &pathError; !errors.As(got, want) {
		t.Fatalf("got=%T, want=%T", got, want)
	}

	if got, want := err.Error(), fmt.Sprintf("stat %s: no such file or directory", file); got != want {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}
