package fd

import (
	"errors"
	"io"
	"os/exec"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/test/spy"
	"github.com/google/go-cmp/cmp"
)

func TestFdSearchRoots(t *testing.T) {
	t.Setenv("HOME", "/home/test")
	r := &spy.SpyRunner{}
	if err := (Cmd{Runner: r}).Fd(io.Discard, "-project", "/outside/with spaces"); err != nil {
		t.Fatal(err)
	}
	want := [][]string{{"fd", "--follow", "--hidden", "--exclude", "{.git,node_modules,Library}", ".", "--base-directory", "/home/test", "--", "-project", "/outside/with spaces"}}
	if diff := cmp.Diff(want, r.Calls); diff != "" {
		t.Fatal(diff)
	}
}

func TestFdErrors(t *testing.T) {
	failure := errors.New("could not start fd")
	exitErr := exec.Command("sh", "-c", "exit 2").Run()
	for _, tt := range []struct {
		name string
		err  error
		want error
	}{
		{name: "success"},
		{name: "start failure", err: failure, want: failure},
		{name: "closed fzf pipe", err: exec.Command("sh", "-c", "exit 1").Run()},
		{name: "unexpected exit", err: exitErr, want: exitErr},
	} {
		t.Run(tt.name, func(t *testing.T) {
			r := &spy.SpyRunner{Responses: []spy.Response{{OnRun: func(*exec.Cmd) error { return tt.err }}}}
			if err := (Cmd{Runner: r}).Fd(io.Discard); !errors.Is(err, tt.want) {
				t.Fatalf("Fd = %v, want %v", err, tt.want)
			}
		})
	}
}
