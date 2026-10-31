package fd

import (
	"io"
	"testing"

	"github.com/eskelinenantti/herdr-zipln/internal/test/mock"
	"github.com/eskelinenantti/herdr-zipln/internal/test/spy"
	"github.com/google/go-cmp/cmp"
)

func TestFdReturnsUnmodifiedOutput(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	want := "./project/\nproject/file.txt\nproject/file.txt\n"
	runner := &spy.SpyRunner{Responses: []spy.Response{
		{},
		{OnRun: mock.WriteToStdout(want)},
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
	wantCalls := [][]string{{"fd", "--type", "d", "--max-depth", "2", "--follow", "--hidden", "--exclude", "{.git,node_modules,Library}", ".", "--base-directory", home}, {"fd", "--follow", "--hidden", "--exclude", "{.git,node_modules,Library}", ".", "--base-directory", home}}
	if diff := cmp.Diff(wantCalls, runner.Calls); diff != "" {
		t.Fatal(diff)
	}
}
