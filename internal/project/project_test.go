package project

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

type gitStub struct {
	root  string
	err   error
	calls []string
}

func (g *gitStub) RevParse(cwd string) (string, error) {
	g.calls = append(g.calls, cwd)
	return g.root, g.err
}

func TestName(t *testing.T) {
	for _, tt := range []struct {
		path string
		want string
	}{
		{"/home/user/project", "project-b28d"},
		{"/other/project", "project-6a86"},
		{"/home/user/my.project", "my_project-4f89"},
		{"relative/project", "project-7bef"},
		{"/", "/-4209"},
		{"", "_-da39"},
		{"/home/user/プロジェクト", "プロジェクト-8b12"},
	} {
		t.Run(tt.path, func(t *testing.T) {
			if got := Name(tt.path); got != tt.want {
				t.Fatalf("Name(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestForPathInvalidPath(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "file.txt")
	if err := os.WriteFile(file, nil, 0644); err != nil {
		t.Fatal(err)
	}
	dangling := filepath.Join(root, "dangling")
	if err := os.Symlink(filepath.Join(root, "missing"), dangling); err != nil {
		t.Fatal(err)
	}
	for _, target := range []string{filepath.Join(root, "missing"), filepath.Join(file, "child"), dangling} {
		t.Run(filepath.Base(target), func(t *testing.T) {
			git := &gitStub{}
			got, err := ForPath(target, git)
			if !errors.Is(err, ErrInvalidPath) {
				t.Fatalf("got %v, want %v", err, ErrInvalidPath)
			}
			if got != (Target{}) {
				t.Fatalf("got %+v, want an empty target", got)
			}
			if len(git.calls) != 0 {
				t.Fatalf("git called for invalid path: %v", git.calls)
			}
		})
	}
}

func TestForPath(t *testing.T) {
	root := t.TempDir()
	directory := filepath.Join(root, "nested")
	if err := os.Mkdir(directory, 0755); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(directory, "file.txt")
	if err := os.WriteFile(file, nil, 0644); err != nil {
		t.Fatal(err)
	}
	t.Chdir(root)
	for _, tt := range []struct {
		name   string
		path   string
		isFile bool
	}{
		{"directory", directory, false},
		{"relative directory", "nested", false},
		{"file", file, true},
		{"relative file", filepath.Join("nested", "file.txt"), true},
	} {
		for _, inRepository := range []bool{true, false} {
			repository := "outside repository"
			if inRepository {
				repository = "inside repository"
			}
			t.Run(tt.name+"/"+repository, func(t *testing.T) {
				git := &gitStub{root: root}
				if !inRepository {
					git.err = errors.New("not a repository")
				}
				got, err := ForPath(tt.path, git)
				if err != nil {
					t.Fatal(err)
				}
				wantDir := directory
				wantFile := ""
				if tt.isFile {
					wantFile = tt.path
					if inRepository {
						wantDir = root
					}
					if len(git.calls) != 1 || git.calls[0] != filepath.Dir(tt.path) {
						t.Fatalf("git calls = %v, want [%s]", git.calls, filepath.Dir(tt.path))
					}
				} else if len(git.calls) != 0 {
					t.Fatalf("directory target consulted git: %v", git.calls)
				}
				want := Target{Project: Project{Name: Name(wantDir), WorkingDir: wantDir}, File: wantFile}
				if got != want {
					t.Fatalf("got %+v, want %+v", got, want)
				}
			})
		}
	}
}
