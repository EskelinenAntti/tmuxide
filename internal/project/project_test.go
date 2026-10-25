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

func TestForFile(t *testing.T) {
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
		target string
		cwd    string
	}{
		{"file", file, directory},
		{"directory", directory, directory},
		{"relative file", filepath.Join("nested", "file.txt"), "nested"},
		{"relative directory", "nested", "nested"},
	} {
		for _, inRepository := range []bool{true, false} {
			name := "outside repository"
			if inRepository {
				name = "inside repository"
			}
			t.Run(tt.name+"/"+name, func(t *testing.T) {
				git := &gitStub{root: root}
				wantDir := root
				if !inRepository {
					git.err = errors.New("not a repository")
					wantDir = directory
				}
				got, err := ForFile(tt.target, git)
				if err != nil {
					t.Fatal(err)
				}
				want := Project{Name: Name(wantDir), WorkingDir: wantDir}
				if got != want {
					t.Fatalf("got %+v, want %+v", got, want)
				}
				if len(git.calls) != 1 || git.calls[0] != tt.cwd {
					t.Fatalf("git calls = %v, want [%s]", git.calls, tt.cwd)
				}
			})
		}
	}
}

func TestForFileInvalidPath(t *testing.T) {
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
			got, err := ForFile(target, git)
			if !errors.Is(err, ErrInvalidPath) {
				t.Fatalf("got %v, want %v", err, ErrInvalidPath)
			}
			if got != (Project{}) {
				t.Fatalf("got %+v, want an empty project", got)
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

func TestForPathInvalidPath(t *testing.T) {
	git := &gitStub{}
	got, err := ForPath(filepath.Join(t.TempDir(), "missing"), git)
	if !errors.Is(err, ErrInvalidPath) || got != (Target{}) {
		t.Fatalf("got %+v, %v; want empty target and %v", got, err, ErrInvalidPath)
	}
	if len(git.calls) != 0 {
		t.Fatalf("invalid target consulted git: %v", git.calls)
	}
}
