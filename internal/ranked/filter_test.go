package ranked

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/eskelinenantti/herdr-zipln/internal/history"
	"github.com/google/go-cmp/cmp"
)

func TestFilter(t *testing.T) {
	if _, err := exec.LookPath("fzf"); err != nil {
		t.Skip("requires fzf")
	}
	t.Setenv("HOME", t.TempDir())
	t.Setenv("XDG_STATE_HOME", t.TempDir())
	input := "project/\nproject/child/file.txt\nproject/other.txt\nelsewhere/project-notes.txt\nother/\n"
	filter := func(query string) string {
		t.Helper()
		var output bytes.Buffer
		if err := Filter(strings.NewReader(input), &output, query); err != nil {
			t.Fatal(err)
		}
		return output.String()
	}
	baseline := filter("proj")
	if first, _, _ := strings.Cut(baseline, "\n"); first != "project/" {
		t.Fatalf("top-level folder must rank above its children: %q", baseline)
	}
	target := "project/child/file.txt"
	history.RecordQuery("proj", filepath.Join(os.Getenv("HOME"), target))
	got := filter("proj")
	want := target + "\n" + strings.Replace(baseline, target+"\n", "", 1)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatal(diff)
	}
	// The preference must not leak into a different query matching the same paths.
	if got := filter("project"); strings.HasPrefix(got, target+"\n") {
		t.Fatalf("unrelated query was reordered: %q", got)
	}
	history.RecordQuery("proj", filepath.Join(os.Getenv("HOME"), "missing"))
	if diff := cmp.Diff(baseline, filter("proj")); diff != "" {
		t.Fatal(diff)
	}
	history.RecordQuery("proj", filepath.Join(os.Getenv("HOME"), "other"))
	if diff := cmp.Diff(baseline, filter("proj")); diff != "" {
		t.Fatal(diff)
	}
	if got := filter("no-such-match"); got != "" {
		t.Fatalf("no matches = %q", got)
	}
	t.Setenv("FZF_DEFAULT_OPTS", "--no-sort --print-query --reverse")
	if diff := cmp.Diff(baseline, filter("proj")); diff != "" {
		t.Fatal(diff)
	}
}
