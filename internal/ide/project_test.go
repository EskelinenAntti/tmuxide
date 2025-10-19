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

var shellMock = shell.Shell{
	Git: GitMock{
		repository: "",
		err:        errors.New("not inside git repository"),
	},
	Tmux: nil,
	Path: PathMock{},
}

var testEditor string = "editor"

func wantWindows(target string) []tmux.WindowCommand {
	return []tmux.WindowCommand{
		{
			Cmd:  testEditor,
			Args: []string{target},
		}}
}

func wantWindowsInRepository(target string) []tmux.WindowCommand {
	return []tmux.WindowCommand{
		{
			Cmd:  testEditor,
			Args: []string{target},
		},
		{
			Cmd:  "lazygit",
			Args: nil,
		},
	}
}

func TestProjectForDirectory(t *testing.T) {
	t.Setenv("EDITOR", testEditor)
	var dir = t.TempDir()

	project, err := ProjectFor(dir, shellMock)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	if got, want := project.Root, dir; got != want {
		t.Fatalf("got=%v, want=%v", got, want)
	}

	if got, want := project.Name, filepath.Base(dir)+"-"; !strings.HasPrefix(got, want) {
		t.Fatalf("%v did not start with %v", got, want)
	}

	if got, want := project.Windows, wantWindows(dir); !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}

}

func TestProjectForFile(t *testing.T) {
	t.Setenv("EDITOR", testEditor)

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

	if got, want := project.Name, "file_txt-"; !strings.HasPrefix(got, want) {
		t.Fatalf("%v did not start with %v", got, want)
	}

	if got, want := project.Windows, wantWindows(file); !reflect.DeepEqual(got, want) {
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

	if got, want := project.Name, filepath.Base(dir)+"-"; !strings.HasPrefix(got, want) {
		t.Fatalf("%v did not start with %v", got, want)
	}

	if got, want := project.Windows, wantWindowsInRepository(dir); !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestProjectForFileInRepository(t *testing.T) {
	t.Setenv("EDITOR", testEditor)

	repository := t.TempDir()
	file := repository + "/file.txt"
	os.WriteFile(file, []byte{}, 0644)

	shellMock := shell.Shell{
		Git: GitMock{
			repository: repository,
			err:        nil,
		},
		Path: PathMock{},
	}

	project, err := ProjectFor(file, shellMock)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	if got, want := project.Root, repository; got != want {
		t.Fatalf("got=%v, want=%v", got, want)
	}

	if got, want := project.Name, "file_txt-"; !strings.HasPrefix(got, want) {
		t.Fatalf("%v did not start with %v", got, want)
	}

	if got, want := project.Windows, wantWindowsInRepository(file); !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestProjectForInvalidFile(t *testing.T) {
	t.Setenv("EDITOR", testEditor)

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
