package runner

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
	"testing/iotest"
	"time"
)

func TestStdoutStreamsAndCancels(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	cmd := exec.CommandContext(ctx, "sh", "-c", "printf 'project/\n'; exec sleep 60")
	input, err := (CmdRunner{}).Stdout(cmd)
	if err != nil {
		cancel()
		t.Fatal(err)
	}
	defer input.Close()
	defer cancel()
	line, err := bufio.NewReader(input).ReadString('\n')
	if err != nil || line != "project/\n" {
		t.Fatalf("Read = %q, %v", line, err)
	}
	cancel()
	if err := input.Close(); !errors.Is(err, ErrCommandFailed) || cmd.ProcessState == nil {
		t.Fatalf("command was not stopped and reaped: %v", err)
	}
}

func TestRunReturnsExitStatus(t *testing.T) {
	err := (CmdRunner{}).Run(exec.Command("sh", "-c", "exit 7"))
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || exitErr.ExitCode() != 7 {
		t.Fatalf("Run = %v, want exit status 7", err)
	}
}

func TestStdoutReturnsExitStatusAtEOF(t *testing.T) {
	cmd := exec.Command("sh", "-c", "printf 'project/\n'; exit 7")
	input, err := (CmdRunner{}).Stdout(cmd)
	if err != nil {
		t.Fatal(err)
	}
	defer input.Close()
	output, err := io.ReadAll(input)
	var exitErr *exec.ExitError
	if !errors.Is(err, ErrCommandFailed) || !errors.As(err, &exitErr) || exitErr.ExitCode() != 7 {
		t.Fatalf("ReadAll error = %v, want exit status 7", err)
	}
	if string(output) != "project/\n" || cmd.ProcessState == nil {
		t.Fatalf("output = %q, process state = %v", output, cmd.ProcessState)
	}
}

func TestStdoutClosesSuccessfullyAfterEOF(t *testing.T) {
	cmd := exec.Command("sh", "-c", "printf 'project/\n'")
	input, err := (CmdRunner{}).Stdout(cmd)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.ReadAll(input); err != nil {
		input.Close()
		t.Fatal(err)
	}
	for range 2 {
		if err := input.Close(); err != nil {
			t.Fatal(err)
		}
	}
}

func TestStdoutCloseStopsUnreadCommand(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "sh", "-c", "exec sleep 60")
	input, err := (CmdRunner{}).Stdout(cmd)
	if err != nil {
		t.Fatal(err)
	}
	if err := input.Close(); !errors.Is(err, ErrCommandFailed) || cmd.ProcessState == nil || ctx.Err() != nil {
		t.Fatalf("Close did not stop and reap command: %v, context error = %v", err, ctx.Err())
	}
}

func TestPipeStopsCommandOnInputFailure(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "sh", "-c", "exec sleep 60")
	want := errors.New("upstream failed")
	err := (CmdRunner{}).Pipe(iotest.ErrReader(want), cmd, io.Discard)
	var inputErr *InputError
	if !errors.As(err, &inputErr) || !errors.Is(err, want) {
		t.Fatalf("Pipe error = %v, want upstream InputError", err)
	}
	if cmd.ProcessState == nil || ctx.Err() != nil {
		t.Fatalf("Pipe did not stop and reap command: state = %v, context error = %v", cmd.ProcessState, ctx.Err())
	}
}

func TestRunLeavesFileStdinUnwrapped(t *testing.T) {
	stdin, err := os.CreateTemp(t.TempDir(), "stdin")
	if err != nil {
		t.Fatal(err)
	}
	defer stdin.Close()
	if _, err := stdin.WriteString("project\n"); err != nil {
		t.Fatal(err)
	}
	if _, err := stdin.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("sh", "-c", `read line; test "$line" = project`)
	cmd.Stdin = stdin
	if err := (CmdRunner{}).Run(cmd); err != nil {
		t.Fatal(err)
	}
	if cmd.Stdin != stdin {
		t.Fatal("file stdin was wrapped, which would break terminal handling")
	}
}

func TestRunLeavesReaderStdinUnwrapped(t *testing.T) {
	want := errors.New("input failed")
	stdin := iotest.ErrReader(want)
	cmd := exec.Command("sh", "-c", "read line; exit 0")
	cmd.Stdin = stdin
	err := (CmdRunner{}).Run(cmd)
	var inputErr *InputError
	if !errors.Is(err, want) || errors.As(err, &inputErr) {
		t.Fatalf("Run error = %v, want ordinary command error", err)
	}
	if cmd.Stdin != stdin {
		t.Fatal("Run replaced the command's stdin")
	}
}

func TestPipeConnectsStdinAndStdout(t *testing.T) {
	stdin := strings.NewReader("project\n")
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", `read line; printf '%s\n' "$line"`)
	if err := (CmdRunner{}).Pipe(stdin, cmd, &stdout); err != nil {
		t.Fatal(err)
	}
	if stdout.String() != "project\n" {
		t.Fatalf("stdout = %q", stdout.String())
	}
}

func TestPipeKeepsConsumerExitStatusSeparate(t *testing.T) {
	cmd := exec.Command("sh", "-c", "exit 130")
	err := (CmdRunner{}).Pipe(nil, cmd, io.Discard)
	var inputErr *InputError
	var exitErr *exec.ExitError
	if errors.As(err, &inputErr) || !errors.As(err, &exitErr) || exitErr.ExitCode() != 130 {
		t.Fatalf("Pipe error = %v, want consumer exit status 130", err)
	}
}

func TestPipeMarksProducerExitStatusAsInputFailure(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	runner := CmdRunner{}
	stdin, err := runner.Stdout(exec.CommandContext(ctx, "sh", "-c", "exit 130"))
	if err != nil {
		t.Fatal(err)
	}
	defer stdin.Close()
	cmd := exec.CommandContext(ctx, "sh", "-c", "exec sleep 60")
	err = runner.Pipe(stdin, cmd, io.Discard)
	var inputErr *InputError
	var exitErr *exec.ExitError
	if !errors.As(err, &inputErr) || !errors.As(inputErr, &exitErr) || exitErr.ExitCode() != 130 {
		t.Fatalf("Pipe error = %v, want upstream exit status 130 in InputError", err)
	}
	if cmd.ProcessState == nil || ctx.Err() != nil {
		t.Fatalf("consumer not stopped and reaped: state = %v, context error = %v", cmd.ProcessState, ctx.Err())
	}
}
