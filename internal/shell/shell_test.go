package shell

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/eskelinenantti/zipln/internal/test/mock"
)

func TestInitMissingDependencies(t *testing.T) {
	for _, missing := range [][]string{nil, {"fd"}, {"fzf"}, {"fd", "fzf"}, {"fd", "fzf", "git"}} {
		t.Run(strings.Join(missing, ","), func(t *testing.T) {
			_, err := Init(mock.Path{Missing: missing}, nil)
			if len(missing) == 0 {
				if err != nil {
					t.Fatal(err)
				}
				return
			}
			var deps MissingToolsError
			if !errors.As(err, &deps) || !reflect.DeepEqual(deps.Tools, missing) {
				t.Fatalf("Init error = %v, want missing tools %v", err, missing)
			}
			if want := "couldn't find " + strings.Join(missing, ", "); err.Error() != want {
				t.Fatalf("error = %q, want %q", err, want)
			}
			var single NotInstalledError
			if !errors.Is(err, ErrCommandNotInstalled) || !errors.As(err, &single) || single.Cmd != missing[0] {
				t.Fatalf("missing dependency error no longer unwraps: %v", err)
			}
		})
	}
}
