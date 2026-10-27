package fzf

import (
	"io"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/test/spy"
	"github.com/google/go-cmp/cmp"
)

func TestFzfPreservesRankedInputOrder(t *testing.T) {
	r := &spy.SpyRunner{}
	input, err := (Cmd{Runner: r}).Fzf(io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	if err := input.Close(); err != nil {
		t.Fatal(err)
	}
	want := [][]string{{"fzf", "--reverse", "--height", "70%", "--tmux", "70%", "--no-sort"}}
	if diff := cmp.Diff(want, r.Calls); diff != "" {
		t.Fatal(diff)
	}
}
