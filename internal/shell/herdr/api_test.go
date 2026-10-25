package herdr

import (
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/project"
	"github.com/eskelinenantti/tmuxide/internal/test/spy"
	"github.com/google/go-cmp/cmp"
)

func TestCreateWorkspace(t *testing.T) {
	for _, tt := range []struct {
		name string
		json string
		id   string
		tab  commandTab
		err  string
	}{
		{"created", created, "w1", commandTab{tabID: "w1:t1", paneID: "w1:p1"}, ""},
		{"no command tab", `{"result":{"workspace":{"workspace_id":"w1"}}}`, "w1", commandTab{}, ""},
		{"missing workspace id", ok, "", commandTab{}, "herdr workspace create response has no workspace id"},
		{"API error", `{"error":{"code":"failed","message":"reason"}}`, "", commandTab{}, "herdr: failed: reason"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			runner := &spy.SpyRunner{Responses: responses(tt.json)}
			id, tab, err := (Cmd{Runner: runner, Binary: "/custom/herdr"}).createWorkspace(project.Project{Name: "hashed-name", WorkingDir: "/projects/project"})
			if (err != nil && err.Error() != tt.err) || (err == nil && tt.err != "") {
				t.Fatalf("got error %v, want %q", err, tt.err)
			}
			if id != tt.id || tab != tt.tab {
				t.Fatalf("got %q, %+v; want %q, %+v", id, tab, tt.id, tt.tab)
			}
			want := [][]string{{"/custom/herdr", "workspace", "create", "--cwd", "/projects/project", "--label", "project", "--no-focus"}}
			if diff := cmp.Diff(want, runner.Calls); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestCreateTab(t *testing.T) {
	for _, tt := range []struct {
		name string
		json string
		tab  commandTab
		err  string
	}{
		{"created", tab, commandTab{tabID: "w1:t2", paneID: "w1:p2"}, ""},
		{"missing identifiers", ok, commandTab{}, ""},
		{"API error", `{"error":{"code":"failed","message":"reason"}}`, commandTab{}, "herdr: failed: reason"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			runner := &spy.SpyRunner{Responses: responses(tt.json)}
			tab, err := (Cmd{Runner: runner, Binary: "/custom/herdr"}).createTab("w1", "/project")
			if (err != nil && err.Error() != tt.err) || (err == nil && tt.err != "") {
				t.Fatalf("got error %v, want %q", err, tt.err)
			}
			if tab != tt.tab {
				t.Fatalf("got %+v, want %+v", tab, tt.tab)
			}
			want := [][]string{{"/custom/herdr", "tab", "create", "--workspace", "w1", "--cwd", "/project", "--no-focus"}}
			if diff := cmp.Diff(want, runner.Calls); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestRunInTab(t *testing.T) {
	failure := errors.New("command failed")
	for _, tt := range []struct {
		name   string
		tab    commandTab
		failAt int
		calls  int
	}{
		{"run then focus", commandTab{tabID: "w1:t2", paneID: "w1:p2"}, -1, 2},
		{"missing pane", commandTab{tabID: "w1:t2"}, -1, 0},
		{"missing tab", commandTab{paneID: "w1:p2"}, -1, 0},
		{"missing both identifiers", commandTab{}, -1, 0},
		{"run failure skips focus", commandTab{tabID: "w1:t2", paneID: "w1:p2"}, 0, 1},
		{"focus failure", commandTab{tabID: "w1:t2", paneID: "w1:p2"}, 1, 2},
	} {
		t.Run(tt.name, func(t *testing.T) {
			runner := &spy.SpyRunner{}
			if tt.failAt >= 0 {
				runner.Responses = make([]spy.Response, tt.failAt+1)
				runner.Responses[tt.failAt].OnRun = func(*exec.Cmd) error { return failure }
			}
			err := (Cmd{Runner: runner, Binary: "/custom/herdr"}).runInTab(tt.tab, []string{"nvim", "a'b"})
			if tt.calls == 0 {
				if err == nil || err.Error() != "herdr create response has no root pane or tab id" {
					t.Fatalf("got error %v, want missing identifiers error", err)
				}
			} else if tt.failAt >= 0 {
				if !errors.Is(err, failure) {
					t.Fatalf("got error %v, want %v", err, failure)
				}
			} else if err != nil {
				t.Fatal(err)
			}
			want := [][]string{
				{"/custom/herdr", "pane", "run", "w1:p2", "'nvim' 'a'\"'\"'b'"},
				{"/custom/herdr", "tab", "focus", "w1:t2"},
			}
			if len(runner.Calls) != tt.calls {
				t.Fatalf("got calls %v, want %d calls", runner.Calls, tt.calls)
			}
			if tt.calls > 0 {
				if diff := cmp.Diff(want[:tt.calls], runner.Calls); diff != "" {
					t.Fatal(diff)
				}
			}
		})
	}
}

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

func TestListResponses(t *testing.T) {
	for _, tt := range []struct {
		name string
		kind string
		json string
		want any
		err  string
	}{
		{"workspaces", "workspace", existing, []workspace{{ID: "w1", Label: "project"}}, ""},
		{"empty workspaces", "workspace", empty, []workspace{}, ""},
		{"wrong workspace type", "workspace", `{"result":{"type":"ok","workspaces":[]}}`, []workspace(nil), "invalid herdr workspace list response"},
		{"missing workspaces", "workspace", `{"result":{"type":"workspace_list"}}`, []workspace(nil), "invalid herdr workspace list response"},
		{"null workspaces", "workspace", `{"result":{"type":"workspace_list","workspaces":null}}`, []workspace(nil), "invalid herdr workspace list response"},
		{"workspace API error", "workspace", `{"error":{"code":"failed","message":"reason"}}`, []workspace(nil), "cannot list workspaces; start a Herdr session first: herdr: failed: reason"},
		{"panes", "pane", `{"result":{"type":"pane_list","panes":[{"workspace_id":"w1","cwd":"/project"}]}}`, []pane{{WorkspaceID: "w1", Cwd: "/project"}}, ""},
		{"empty panes", "pane", `{"result":{"type":"pane_list","panes":[]}}`, []pane{}, ""},
		{"wrong pane type", "pane", `{"result":{"type":"ok","panes":[]}}`, []pane(nil), "invalid herdr pane list response"},
		{"missing panes", "pane", `{"result":{"type":"pane_list"}}`, []pane(nil), "invalid herdr pane list response"},
		{"null panes", "pane", `{"result":{"type":"pane_list","panes":null}}`, []pane(nil), "invalid herdr pane list response"},
		{"pane API error", "pane", `{"error":{"code":"failed","message":"reason"}}`, []pane(nil), "herdr: failed: reason"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			runner := &spy.SpyRunner{Responses: responses(tt.json)}
			cmd := Cmd{Runner: runner, Binary: "herdr"}
			var got any
			var err error
			if tt.kind == "workspace" {
				got, err = cmd.listWorkspaces()
			} else {
				got, err = cmd.listPanes()
			}
			if tt.err == "" {
				if err != nil {
					t.Fatal(err)
				}
			} else if err == nil || err.Error() != tt.err {
				t.Fatalf("got error %v, want %q", err, tt.err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff([][]string{{"herdr", tt.kind, "list"}}, runner.Calls); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
