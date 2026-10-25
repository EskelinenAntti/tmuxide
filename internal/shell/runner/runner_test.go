package runner

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestCommandFailures(t *testing.T) {
	for _, start := range []bool{false, true} {
		name := "Run"
		if start {
			name = "Start"
		}
		t.Run(name, func(t *testing.T) {
			cmd := exec.Command(filepath.Join(t.TempDir(), "missing"), "argument")
			runner := CmdRunner{}
			var err error
			if start {
				var writer io.WriteCloser
				writer, err = runner.Start(cmd)
				if writer != nil {
					t.Fatal("failed start returned a writer")
				}
			} else {
				err = runner.Run(cmd)
			}
			if !errors.Is(err, ErrCommandFailed) {
				t.Fatalf("got %v, want command failure", err)
			}
			var pathErr *os.PathError
			if !errors.As(err, &pathErr) {
				t.Fatalf("underlying path error was lost: %v", err)
			}
			want := fmt.Sprintf("'%v' command failed: %s", cmd.Args, pathErr)
			if err.Error() != want {
				t.Fatalf("got %q, want %q", err, want)
			}
		})
	}
}

func TestCommandSuccess(t *testing.T) {
	for _, start := range []bool{false, true} {
		name := "Run"
		if start {
			name = "Start"
		}
		t.Run(name, func(t *testing.T) {
			cmd := exec.Command(os.Args[0], "-test.run=^TestRunnerProcess$")
			cmd.Env = append(os.Environ(), "TMUXIDE_RUNNER_PROCESS=1")
			runner := CmdRunner{}
			if start {
				writer, err := runner.Start(cmd)
				if err != nil {
					t.Fatal(err)
				}
				if err := writer.Close(); err != nil {
					t.Fatal(err)
				}
			} else if err := runner.Run(cmd); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestRunnerProcess(t *testing.T) {
	switch os.Getenv("TMUXIDE_RUNNER_PROCESS") {
	case "1":
		os.Exit(0)
	case "failure":
		os.Exit(7)
	case "echo":
		if _, err := io.Copy(os.Stdout, os.Stdin); err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}
}

func TestStartWriteClose(t *testing.T) {
	cmd := exec.Command(os.Args[0], "-test.run=^TestRunnerProcess$")
	cmd.Env = append(os.Environ(), "TMUXIDE_RUNNER_PROCESS=echo")
	var output bytes.Buffer
	cmd.Stdout = &output
	var starter Starter = CmdRunner{}
	writer, err := starter.Start(cmd)
	if err != nil {
		t.Fatal(err)
	}
	const input = "candidate one\ncandidate two\n"
	n, writeErr := io.WriteString(writer, input)
	closeErr := writer.Close()
	if writeErr != nil || n != len(input) {
		t.Fatalf("wrote %d bytes: %v", n, writeErr)
	}
	if closeErr != nil {
		t.Fatal(closeErr)
	}
	if output.String() != input {
		t.Fatalf("got %q, want %q", output.String(), input)
	}
}

func TestRunExitFailure(t *testing.T) {
	cmd := exec.Command(os.Args[0], "-test.run=^TestRunnerProcess$")
	cmd.Env = append(os.Environ(), "TMUXIDE_RUNNER_PROCESS=failure")
	err := (CmdRunner{}).Run(cmd)
	var exitErr *exec.ExitError
	if !errors.Is(err, ErrCommandFailed) || !errors.As(err, &exitErr) {
		t.Fatalf("got %v, want a wrapped exit error", err)
	}
	if exitErr.ExitCode() != 7 {
		t.Fatalf("got exit code %d, want 7", exitErr.ExitCode())
	}
}

func TestStartCloseExitFailure(t *testing.T) {
	cmd := exec.Command(os.Args[0], "-test.run=^TestRunnerProcess$")
	cmd.Env = append(os.Environ(), "TMUXIDE_RUNNER_PROCESS=failure")
	writer, err := (CmdRunner{}).Start(cmd)
	if err != nil {
		t.Fatal(err)
	}
	err = writer.Close()
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || exitErr.ExitCode() != 7 {
		t.Fatalf("got %v, want exit code 7", err)
	}
}

func TestWriteAfterClose(t *testing.T) {
	cmd := exec.Command(os.Args[0], "-test.run=^TestRunnerProcess$")
	cmd.Env = append(os.Environ(), "TMUXIDE_RUNNER_PROCESS=echo")
	writer, err := (CmdRunner{}).Start(cmd)
	if err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	n, err := writer.Write([]byte("candidate"))
	if n != 0 || !errors.Is(err, os.ErrClosed) {
		t.Fatalf("wrote %d bytes: %v, want a closed pipe error", n, err)
	}
}
