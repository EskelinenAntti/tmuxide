package fd

import (
	"io"
	"testing"

	"github.com/eskelinenantti/zipln/internal/test/mock"
	"github.com/eskelinenantti/zipln/internal/test/spy"
	"github.com/google/go-cmp/cmp"
)

func TestFdStreamsShallowPathsBeforeDeepPaths(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	shallow := "project/\nproject/file.txt\n"
	deep := "project/nested/file.txt\n"
	want := shallow + deep
	runner := &spy.SpyRunner{Responses: []spy.Response{
		{OnRun: mock.WriteToStdout(shallow)},
		{OnRun: mock.WriteToStdout(deep)},
	}}
	input, err := (Cmd{Runner: runner}).Fd()
	if err != nil {
		t.Fatal(err)
	}
	defer input.Close()
	output, err := io.ReadAll(input)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(want, string(output)); diff != "" {
		t.Fatal(diff)
	}
	wantCalls := [][]string{
		{"fd", "--max-depth", "2", "--follow", "--hidden", "--exclude", "{.git,node_modules,Library}", ".", "--base-directory", home},
		{"fd", "--min-depth", "3", "--follow", "--hidden", "--exclude", "{.git,node_modules,Library}", ".", "--base-directory", home},
	}
	if diff := cmp.Diff(wantCalls, runner.Calls); diff != "" {
		t.Fatal(diff)
	}
}
