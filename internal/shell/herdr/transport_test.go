package herdr

import (
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/test/spy"
	"github.com/google/go-cmp/cmp"
)

func TestOutput(t *testing.T) {
	failure := errors.New("command failed")
	for _, tt := range []struct {
		name string
		err  error
	}{
		{"success", nil},
		{"failure", failure},
	} {
		t.Run(tt.name, func(t *testing.T) {
			runner := &spy.SpyRunner{Responses: []spy.Response{{OnRun: func(cmd *exec.Cmd) error {
				_, _ = cmd.Stdout.Write([]byte("response\n"))
				_, _ = cmd.Stderr.Write([]byte("  diagnostic\n"))
				return tt.err
			}}}}
			output, err := (Cmd{Runner: runner, Binary: "/custom/herdr"}).output("workspace", "list")
			if tt.err == nil {
				if err != nil || string(output) != "response\n" {
					t.Fatalf("got %q, %v", output, err)
				}
			} else if !errors.Is(err, failure) || err.Error() != "herdr workspace list: command failed: diagnostic" || output != nil {
				t.Fatalf("got %q, %v", output, err)
			}
			if diff := cmp.Diff([][]string{{"/custom/herdr", "workspace", "list"}}, runner.Calls); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestCall(t *testing.T) {
	type result struct {
		Value string `json:"value"`
	}
	for _, tt := range []struct {
		name string
		json string
		want result
		err  string
	}{
		{"result", `{"result":{"value":"decoded"}}`, result{Value: "decoded"}, ""},
		{"unknown fields", `{"result":{"value":"decoded","unknown":true}}`, result{Value: "decoded"}, ""},
		{"empty result", `{"result":{}}`, result{}, ""},
		{"missing result", `{}`, result{}, "herdr response has no result"},
		{"null result", `{"result":null}`, result{}, "herdr response has no result"},
		{"API error takes priority", `{"result":{},"error":{"code":"failed","message":"reason"}}`, result{}, "herdr: failed: reason"},
		{"invalid JSON", `not json`, result{}, "decode herdr response:"},
		{"invalid result type", `{"result":{"value":123}}`, result{}, "decode herdr response:"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			runner := &spy.SpyRunner{Responses: responses(tt.json)}
			got, err := call[result](Cmd{Runner: runner, Binary: "herdr"}, "workspace", "create")
			if tt.err == "" {
				if err != nil {
					t.Fatal(err)
				}
			} else if err == nil || !strings.HasPrefix(err.Error(), tt.err) {
				t.Fatalf("got %v, want %q", err, tt.err)
			}
			if tt.name == "invalid JSON" {
				var syntaxError *json.SyntaxError
				if !errors.As(err, &syntaxError) {
					t.Fatalf("JSON error not preserved: %v", err)
				}
			}
			if tt.name == "invalid result type" {
				var typeError *json.UnmarshalTypeError
				if !errors.As(err, &typeError) {
					t.Fatalf("JSON type error not preserved: %v", err)
				}
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
