package runner

import (
	"bufio"
	"context"
	"errors"
	"io"
	"os/exec"
	"testing"
	"time"
)

func TestReadStreamsAndCancels(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	cmd := exec.CommandContext(ctx, "sh", "-c", "printf 'project/\n'; exec sleep 60")
	input, err := (CmdRunner{}).Read(cmd)
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

func TestWaitReturnsExitStatus(t *testing.T) {
	process, err := (CmdRunner{}).Start(exec.Command("sh", "-c", "exit 7"))
	if err != nil {
		t.Fatal(err)
	}
	defer process.Close()
	err = process.Wait()
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || exitErr.ExitCode() != 7 {
		t.Fatalf("Wait = %v, want exit status 7", err)
	}
}

func TestReadReturnsExitStatusAtEOF(t *testing.T) {
	cmd := exec.Command("sh", "-c", "printf 'project/\n'; exit 7")
	input, err := (CmdRunner{}).Read(cmd)
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

func TestReadClosesSuccessfullyAfterEOF(t *testing.T) {
	cmd := exec.Command("sh", "-c", "printf 'project/\n'")
	input, err := (CmdRunner{}).Read(cmd)
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

func TestReadCloseStopsUnreadCommand(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "sh", "-c", "exec sleep 60")
	input, err := (CmdRunner{}).Read(cmd)
	if err != nil {
		t.Fatal(err)
	}
	if err := input.Close(); !errors.Is(err, ErrCommandFailed) || cmd.ProcessState == nil || ctx.Err() != nil {
		t.Fatalf("Close did not stop and reap command: %v, context error = %v", err, ctx.Err())
	}
}

func TestKillReapsCommand(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "sh", "-c", "exec sleep 60")
	process, err := (CmdRunner{}).Start(cmd)
	if err != nil {
		t.Fatal(err)
	}
	for range 2 {
		if err := process.Kill(); err != nil {
			t.Fatal(err)
		}
	}
	if cmd.ProcessState == nil || ctx.Err() != nil {
		t.Fatalf("Kill did not reap command: state = %v, context error = %v", cmd.ProcessState, ctx.Err())
	}
}

func TestKillAfterWait(t *testing.T) {
	process, err := (CmdRunner{}).Start(exec.Command("sh", "-c", "exit 0"))
	if err != nil {
		t.Fatal(err)
	}
	if err := process.Wait(); err != nil {
		t.Fatal(err)
	}
	if err := process.Kill(); err != nil {
		t.Fatal(err)
	}
}
