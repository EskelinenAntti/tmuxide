package editor

import (
	"errors"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/test/mock"
	"github.com/google/go-cmp/cmp"
)

func TestCommand(t *testing.T) {
	for _, tt := range []struct {
		name    string
		editor  string
		missing []string
		want    []string
		err     error
	}{
		{name: "empty", err: ErrNotConfigured},
		{name: "whitespace", editor: " \t\n", err: ErrNotConfigured},
		{name: "executable", editor: "nvim", want: []string{"nvim"}},
		{name: "arguments", editor: " nvim\t-R  --clean ", want: []string{"nvim", "-R", "--clean"}},
		{name: "absolute path", editor: "/custom/editor -R", want: []string{"/custom/editor", "-R"}},
		{name: "missing executable", editor: "nvim -R", missing: []string{"nvim"}, err: ErrNotInstalled},
		{name: "only executable is checked", editor: "nvim -R", missing: []string{"-R"}, want: []string{"nvim", "-R"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("EDITOR", tt.editor)
			got, err := Command(mock.Path{Missing: tt.missing})
			if !errors.Is(err, tt.err) {
				t.Fatalf("got error %v, want %v", err, tt.err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
