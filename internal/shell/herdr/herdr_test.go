package herdr

import (
	"strings"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/project"
	"github.com/eskelinenantti/tmuxide/internal/test/mock"
	"github.com/eskelinenantti/tmuxide/internal/test/spy"
	"github.com/google/go-cmp/cmp"
)

const empty = `{"result":{"type":"workspace_list","workspaces":[]}}`
const existing = `{"result":{"type":"workspace_list","workspaces":[{"workspace_id":"w1","label":"project"}]}}`
const created = `{"result":{"type":"workspace_created","workspace":{"workspace_id":"w1"},"tab":{"tab_id":"w1:t1"},"root_pane":{"pane_id":"w1:p1"}}}`
const tab = `{"result":{"type":"tab_created","tab":{"tab_id":"w1:t2"},"root_pane":{"pane_id":"w1:p2"}}}`
const ok = `{"result":{"type":"ok"}}`

func responses(json ...string) []spy.Response {
	var result []spy.Response
	for _, text := range json {
		result = append(result, spy.Response{OnRun: mock.WriteToStdout(text)})
	}
	return result
}

func TestOpen(t *testing.T) {
	for _, tt := range []struct {
		name    string
		inside  bool
		command []string
		json    []string
		calls   [][]string
	}{
		{"create folder", true, nil, []string{empty, created, ok}, [][]string{
			{"herdr", "workspace", "list"},
			{"herdr", "workspace", "create", "--cwd", "/project", "--label", "project", "--no-focus"},
			{"herdr", "workspace", "focus", "w1"},
		}},
		{"reuse folder", true, nil, []string{existing, ok}, [][]string{
			{"herdr", "workspace", "list"}, {"herdr", "workspace", "focus", "w1"},
		}},
		{"attach outside", false, nil, []string{existing, ok}, [][]string{
			{"herdr", "workspace", "list"}, {"herdr", "workspace", "focus", "w1"}, {"herdr"},
		}},
		{"file in new workspace", true, []string{"nvim", "/project/file"}, []string{empty, created, ok, ok, ok}, [][]string{
			{"herdr", "workspace", "list"},
			{"herdr", "workspace", "create", "--cwd", "/project", "--label", "project", "--no-focus"},
			{"herdr", "pane", "run", "w1:p1", "'nvim' '/project/file'"},
			{"herdr", "tab", "focus", "w1:t1"}, {"herdr", "workspace", "focus", "w1"},
		}},
		{"file in existing workspace", true, []string{"nvim", "-R", "/project/a'b; $(touch bad)"}, []string{existing, tab, ok, ok, ok}, [][]string{
			{"herdr", "workspace", "list"},
			{"herdr", "tab", "create", "--workspace", "w1", "--cwd", "/project", "--label", "nvim", "--no-focus"},
			{"herdr", "pane", "run", "w1:p2", "'nvim' '-R' '/project/a'\"'\"'b; $(touch bad)'"},
			{"herdr", "tab", "focus", "w1:t2"}, {"herdr", "workspace", "focus", "w1"},
		}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("HERDR_ENV", "")
			if tt.inside {
				t.Setenv("HERDR_ENV", "1")
			}
			runner := &spy.SpyRunner{Responses: responses(tt.json...)}
			err := (Cmd{Runner: runner, Binary: "herdr"}).Open(project.Project{Name: "project", WorkingDir: "/project"}, tt.command)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.calls, runner.Calls); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestFailuresStopWorkflow(t *testing.T) {
	t.Setenv("HERDR_ENV", "1")
	for _, tt := range []struct {
		name string
		json []string
		want string
	}{
		{"invalid JSON", []string{"not json"}, "decode herdr response"},
		{"API error", []string{`{"error":{"code":"not_found","message":"missing"}}`}, "not_found: missing"},
		{"missing result", []string{`{}`}, "no result"},
		{"invalid list", []string{ok}, "invalid herdr workspace list"},
		{"missing workspace id", []string{empty, ok}, "no workspace id"},
		{"missing pane", []string{existing, ok}, "no root pane or tab id"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			runner := &spy.SpyRunner{Responses: responses(tt.json...)}
			err := (Cmd{Runner: runner, Binary: "herdr"}).Open(project.Project{Name: "project"}, []string{"editor", "file"})
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("got %v, want %s", err, tt.want)
			}
			if len(runner.Calls) != len(tt.json) {
				t.Fatal(runner.Calls)
			}
		})
	}
	for i := 0; i < 5; i++ {
		runner := &spy.SpyRunner{Responses: responses(empty, created, ok, ok, ok)}
		runner.Responses[i] = spy.Response{OnRun: mock.SimulateError}
		err := (Cmd{Runner: runner, Binary: "herdr"}).Open(project.Project{Name: "project"}, []string{"editor", "file"})
		if err == nil || len(runner.Calls) != i+1 {
			t.Fatalf("failure %d: %v, %v", i, err, runner.Calls)
		}
	}
}

func TestMissingServerDoesNotCreateWorkspace(t *testing.T) {
	runner := &spy.SpyRunner{Responses: []spy.Response{{OnRun: mock.SimulateError}}}
	err := (Cmd{Runner: runner, Binary: "herdr"}).Open(project.Project{}, nil)
	if err == nil || !strings.Contains(err.Error(), "start a Herdr session first") || len(runner.Calls) != 1 {
		t.Fatal(err, runner.Calls)
	}
}
