package tmux

import (
	"errors"
	"os"
	"os/exec"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/project"
	"github.com/eskelinenantti/tmuxide/internal/test/spy"
	"github.com/google/go-cmp/cmp"
)

func TestOpen(t *testing.T) {
	failure := errors.New("tmux failed")
	tests := []struct {
		name      string
		attached  bool
		command   []string
		responses []spy.Response
		want      [][]string
		wantErr   error
	}{
		{
			name:     "reuse session inside tmux",
			attached: true,
			want: [][]string{
				{"tmux", "has-session", "-t", "project:"},
				{"tmux", "switch-client", "-t", "project:"},
			},
		},
		{
			name: "attach outside tmux",
			want: [][]string{
				{"tmux", "has-session", "-t", "project:"},
				{"tmux", "attach", "-t", "project:"},
			},
		},
		{
			name:      "create session without command",
			attached:  true,
			responses: []spy.Response{{OnRun: func(*exec.Cmd) error { return failure }}},
			want: [][]string{
				{"tmux", "has-session", "-t", "project:"},
				{"tmux", "new-session", "-c", "/project", "-d", "-s", "project"},
				{"tmux", "switch-client", "-t", "project:"},
			},
		},
		{
			name:      "create session with command",
			command:   []string{"nvim", "a file"},
			responses: []spy.Response{{OnRun: func(*exec.Cmd) error { return failure }}},
			want: [][]string{
				{"tmux", "has-session", "-t", "project:"},
				{"tmux", "new-session", "-c", "/project", "-d", "-s", "project", "nvim", "a file"},
				{"tmux", "attach", "-t", "project:"},
			},
		},
		{
			name: "session creation fails",
			responses: []spy.Response{
				{OnRun: func(*exec.Cmd) error { return failure }},
				{OnRun: func(*exec.Cmd) error { return failure }},
			},
			want: [][]string{
				{"tmux", "has-session", "-t", "project:"},
				{"tmux", "new-session", "-c", "/project", "-d", "-s", "project"},
			},
			wantErr: failure,
		},
		{
			name:      "window creation fails",
			command:   []string{"nvim", "file"},
			responses: []spy.Response{{}, {OnRun: func(*exec.Cmd) error { return failure }}},
			want: [][]string{
				{"tmux", "has-session", "-t", "project:"},
				{"tmux", "new-window", "-t", "project:", "-c", "/project", "nvim", "file"},
			},
			wantErr: failure,
		},
		{
			name:      "attach fails",
			responses: []spy.Response{{}, {OnRun: func(*exec.Cmd) error { return failure }}},
			want: [][]string{
				{"tmux", "has-session", "-t", "project:"},
				{"tmux", "attach", "-t", "project:"},
			},
			wantErr: failure,
		},
		{
			name:      "switch fails",
			attached:  true,
			responses: []spy.Response{{}, {OnRun: func(*exec.Cmd) error { return failure }}},
			want: [][]string{
				{"tmux", "has-session", "-t", "project:"},
				{"tmux", "switch-client", "-t", "project:"},
			},
			wantErr: failure,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("TMUX", "")
			if !tt.attached {
				if err := os.Unsetenv("TMUX"); err != nil {
					t.Fatal(err)
				}
			}
			runner := &spy.SpyRunner{Responses: tt.responses}
			mux := Cmd{Runner: runner}
			err := mux.Open(project.Project{Name: "project", WorkingDir: "/project"}, tt.command)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("got error %v, want %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(tt.want, runner.Calls); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestCreationCommand(t *testing.T) {
	tests := []struct {
		name          string
		project       project.Project
		command       []string
		sessionExists bool
		want          []string
	}{
		{
			name:          "reuse session without creating a window",
			sessionExists: true,
		},
		{
			name:          "reuse named session without creating a window",
			project:       project.Project{Name: "project", WorkingDir: "/a directory"},
			sessionExists: true,
		},
		{
			name: "create default session without command",
			want: []string{"tmux", "new-session", "-d"},
		},
		{
			name:    "create default session with command",
			command: []string{"nvim", "-R", "a file"},
			want:    []string{"tmux", "new-session", "-d", "nvim", "-R", "a file"},
		},
		{
			name:    "create named session in directory",
			project: project.Project{Name: "project", WorkingDir: "/a directory"},
			command: []string{"nvim", "-R", "a file"},
			want:    []string{"tmux", "new-session", "-c", "/a directory", "-d", "-s", "project", "nvim", "-R", "a file"},
		},
		{
			name:    "create named session with default directory",
			project: project.Project{Name: "project"},
			want:    []string{"tmux", "new-session", "-d", "-s", "project"},
		},
		{
			name:    "create unnamed session in directory",
			project: project.Project{WorkingDir: "/project"},
			want:    []string{"tmux", "new-session", "-c", "/project", "-d"},
		},
		{
			name:    "create unnamed session in directory with command",
			project: project.Project{WorkingDir: "/a directory"},
			command: []string{"nvim", "-R", "a file"},
			want:    []string{"tmux", "new-session", "-c", "/a directory", "-d", "nvim", "-R", "a file"},
		},
		{
			name:          "create window in named session",
			project:       project.Project{Name: "project", WorkingDir: "/a directory"},
			command:       []string{"nvim", "-R", "a file"},
			sessionExists: true,
			want:          []string{"tmux", "new-window", "-t", "project:", "-c", "/a directory", "nvim", "-R", "a file"},
		},
		{
			name:    "session options precede option-like command arguments",
			project: project.Project{Name: "project", WorkingDir: "/a directory"},
			command: []string{"nvim", "-c", "set number", "-s", "a file"},
			want:    []string{"tmux", "new-session", "-c", "/a directory", "-d", "-s", "project", "nvim", "-c", "set number", "-s", "a file"},
		},
		{
			name:          "window options precede option-like command arguments",
			project:       project.Project{Name: "project", WorkingDir: "/a directory"},
			command:       []string{"nvim", "-c", "set number", "-s", "a file"},
			sessionExists: true,
			want:          []string{"tmux", "new-window", "-t", "project:", "-c", "/a directory", "nvim", "-c", "set number", "-s", "a file"},
		},
		{
			name:          "create window with default session and directory",
			command:       []string{"nvim"},
			sessionExists: true,
			want:          []string{"tmux", "new-window", "nvim"},
		},
		{
			name:          "create window in named session with default directory",
			project:       project.Project{Name: "project"},
			command:       []string{"nvim"},
			sessionExists: true,
			want:          []string{"tmux", "new-window", "-t", "project:", "nvim"},
		},
		{
			name:          "create window in default session with directory",
			project:       project.Project{WorkingDir: "/a directory"},
			command:       []string{"nvim", "a file"},
			sessionExists: true,
			want:          []string{"tmux", "new-window", "-c", "/a directory", "nvim", "a file"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := creationCommand(tt.project, tt.command, tt.sessionExists)
			var args []string
			if cmd != nil {
				args = cmd.Args
			}
			if diff := cmp.Diff(tt.want, args); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestAttachUsesTerminalStreams(t *testing.T) {
	t.Setenv("TMUX", "")
	if err := os.Unsetenv("TMUX"); err != nil {
		t.Fatal(err)
	}
	runner := &spy.SpyRunner{Responses: []spy.Response{{}, {OnRun: func(cmd *exec.Cmd) error {
		if cmd.Stdin != os.Stdin || cmd.Stdout != os.Stdout || cmd.Stderr != os.Stderr {
			t.Error("attach must use the terminal's standard streams")
		}
		return nil
	}}}}
	if err := (Cmd{Runner: runner}).Open(project.Project{Name: "project", WorkingDir: "/project"}, nil); err != nil {
		t.Fatal(err)
	}
}

func TestOpenOmitsEmptyProjectOptions(t *testing.T) {
	t.Setenv("TMUX", "test")
	for _, exists := range []bool{false, true} {
		runner := &spy.SpyRunner{}
		wantCreate := []string{"tmux", "new-window", "nvim"}
		if !exists {
			runner.Responses = []spy.Response{{OnRun: func(*exec.Cmd) error { return errors.New("no session") }}}
			wantCreate = []string{"tmux", "new-session", "-d", "nvim"}
		}
		if err := (Cmd{Runner: runner}).Open(project.Project{}, []string{"nvim"}); err != nil {
			t.Fatal(err)
		}
		want := [][]string{{"tmux", "has-session"}, wantCreate, {"tmux", "switch-client"}}
		if diff := cmp.Diff(want, runner.Calls); diff != "" {
			t.Fatal(diff)
		}
	}
}

func TestEachEditorRequestCreatesUnnamedWindow(t *testing.T) {
	t.Setenv("TMUX", "test")
	runner := &spy.SpyRunner{}
	mux := Cmd{Runner: runner}
	p := project.Project{Name: "project", WorkingDir: "/project"}
	var want [][]string
	for range 2 {
		if err := mux.Open(p, []string{"nvim", "file"}); err != nil {
			t.Fatal(err)
		}
		want = append(want,
			[]string{"tmux", "has-session", "-t", "project:"},
			[]string{"tmux", "new-window", "-t", "project:", "-c", "/project", "nvim", "file"},
			[]string{"tmux", "switch-client", "-t", "project:"},
		)
	}
	if diff := cmp.Diff(want, runner.Calls); diff != "" {
		t.Fatal(diff)
	}
}
