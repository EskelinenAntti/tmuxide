package fzf

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/shell/runner"
	"github.com/eskelinenantti/tmuxide/internal/test/spy"
	"github.com/google/go-cmp/cmp"
)

func TestIsCancelled(t *testing.T) {
	for _, code := range []int{0, 1, 2, 130} {
		t.Run(strconv.Itoa(code), func(t *testing.T) {
			err := fzfExitError(t, code)
			want := code == 130
			if got := isCancelled(err); got != want {
				t.Fatalf("IsCancelled(%v) = %v, want %v", err, got, want)
			}
			if err != nil && isCancelled(fmt.Errorf("command failed: %w", err)) != want {
				t.Fatal("wrapping changed cancellation detection")
			}
		})
	}
	if isCancelled(errors.New("exit status 130")) {
		t.Fatal("an ordinary error is not cancellation")
	}
}

func fzfExitError(t *testing.T, code int) error {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^TestFzfExitProcess$")
	cmd.Env = append(os.Environ(), "TMUXIDE_FZF_EXIT_CODE="+strconv.Itoa(code))
	err := cmd.Run()
	if code == 0 {
		if err != nil {
			t.Fatal(err)
		}
		return nil
	}
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || exitErr.ExitCode() != code {
		t.Fatalf("got %v, want exit code %d", err, code)
	}
	return err
}

func TestFzfExitProcess(t *testing.T) {
	if value := os.Getenv("TMUXIDE_FZF_EXIT_CODE"); value != "" {
		code, err := strconv.Atoi(value)
		if err != nil {
			t.Fatal(err)
		}
		os.Exit(code)
	}
}

func TestFzfNeedsOnlyStart(t *testing.T) {
	for _, tt := range []struct {
		name string
		err  error
	}{
		{"success", nil},
		{"failure", errors.New("start failed")},
	} {
		t.Run(tt.name, func(t *testing.T) {
			process := &spy.SpyRunner{Responses: []spy.RunFunc{func(cmd *exec.Cmd) error {
				if cmd.Stdout == nil || cmd.Stderr != os.Stderr {
					t.Error("fzf must collect output and use terminal stderr")
				}
				return tt.err
			}}}
			startOnly := struct{ runner.Starter }{process}
			_, err := (Cmd{Starter: startOnly}).Select(func(io.Writer) error { return nil })
			if !errors.Is(err, tt.err) {
				t.Fatalf("got %v, want %v", err, tt.err)
			}
			want := [][]string{{"fzf", "--reverse", "--height", "70%", "--tmux", "70%"}}
			if diff := cmp.Diff(want, process.Calls); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

type selectStarter struct {
	start func(*exec.Cmd) (io.WriteCloser, error)
}

func (s selectStarter) Start(cmd *exec.Cmd) (io.WriteCloser, error) {
	return s.start(cmd)
}

type selectionInput struct {
	strings.Builder
	closed bool
	err    error
}

func (i *selectionInput) Close() error {
	i.closed = true
	return i.err
}

func TestSelect(t *testing.T) {
	failure := errors.New("failed")
	cancelled := fzfExitError(t, 130)
	for _, tt := range []struct {
		name        string
		startErr    error
		producerErr error
		closeErr    error
		want        string
		wantErr     error
	}{
		{name: "selection", want: "folder/file"},
		{name: "start failure", startErr: failure, wantErr: failure},
		{name: "candidate failure", producerErr: failure, wantErr: failure},
		{name: "wait failure", closeErr: failure, wantErr: failure},
		{name: "cancelled", closeErr: fmt.Errorf("wrapped: %w", cancelled)},
	} {
		t.Run(tt.name, func(t *testing.T) {
			input := &selectionInput{err: tt.closeErr}
			starter := selectStarter{start: func(cmd *exec.Cmd) (io.WriteCloser, error) {
				if tt.startErr != nil {
					return nil, tt.startErr
				}
				_, err := io.WriteString(cmd.Stdout, "  folder/file\n")
				if err != nil {
					t.Fatal(err)
				}
				return input, nil
			}}
			produced := false
			got, err := (Cmd{Starter: starter}).Select(func(output io.Writer) error {
				produced = true
				if input.closed {
					t.Fatal("input closed before candidates were produced")
				}
				if _, err := io.WriteString(output, "folder/file\n"); err != nil {
					t.Fatal(err)
				}
				return tt.producerErr
			})
			if got != tt.want || !errors.Is(err, tt.wantErr) {
				t.Fatalf("got %q, %v; want %q, %v", got, err, tt.want, tt.wantErr)
			}
			if produced != (tt.startErr == nil) {
				t.Fatal("candidates must only be produced after a successful start")
			}
			if input.closed != (tt.startErr == nil && tt.producerErr == nil) {
				t.Fatal("must wait after producing candidates successfully")
			}
			if produced && input.String() != "folder/file\n" {
				t.Fatalf("candidate stream lost: %q", input.String())
			}
		})
	}
}
