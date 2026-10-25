package herdr

import (
	"errors"
	"os"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/shell"
	"github.com/eskelinenantti/tmuxide/internal/test/mock"
	"github.com/eskelinenantti/tmuxide/internal/test/spy"
)

func TestNew(t *testing.T) {
	for _, tt := range []struct {
		name    string
		binary  string
		missing bool
	}{
		{name: "configured", binary: "/custom/herdr with spaces"},
		{name: "unconfigured"},
		{name: "missing executable", binary: "/missing/herdr", missing: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("HERDR_BIN_PATH", tt.binary)
			path := mock.Path{}
			if tt.missing {
				path.Missing = []string{tt.binary}
			}
			runner := &spy.SpyRunner{}
			cmd, err := New(path, runner)
			switch {
			case tt.binary == "":
				if err == nil || err.Error() != "HERDR_BIN_PATH is required for the Herdr multiplexer" {
					t.Fatalf("got %v, want missing configuration error", err)
				}
			case tt.missing:
				var missing shell.NotInstalledError
				if !errors.As(err, &missing) || missing.Cmd != tt.binary {
					t.Fatalf("got %v, want missing executable %q", err, tt.binary)
				}
			default:
				if err != nil || cmd.Binary != tt.binary || cmd.Runner != runner {
					t.Fatalf("got %+v, %v, want configured command", cmd, err)
				}
			}
			if err != nil && (cmd.Binary != "" || cmd.Runner != nil) {
				t.Fatalf("got partially configured command on failure: %+v", cmd)
			}
		})
	}
}

func TestInSession(t *testing.T) {
	for _, tt := range []struct {
		name  string
		value string
		want  bool
	}{
		{name: "unset"},
		{name: "empty"},
		{name: "inactive", value: "0"},
		{name: "noncanonical", value: "true"},
		{name: "active", value: "1", want: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("HERDR_ENV", tt.value)
			if tt.name == "unset" {
				if err := os.Unsetenv("HERDR_ENV"); err != nil {
					t.Fatal(err)
				}
			}
			if got := InSession(); got != tt.want {
				t.Fatalf("InSession() = %v, want %v", got, tt.want)
			}
		})
	}
}
