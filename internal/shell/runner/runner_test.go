package runner

import (
	"bufio"
	"context"
	"errors"
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
