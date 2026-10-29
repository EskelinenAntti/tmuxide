package runner

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
)

var ErrCommandFailed = errors.New("command failed")

type Runner interface {
	Run(cmd *exec.Cmd) error
	Read(cmd *exec.Cmd) (io.ReadCloser, error)
}

type CmdRunner struct{}
type CmdReadCloser struct {
	io.ReadCloser
	cmd  *exec.Cmd
	wait func() error
}

func (c CmdRunner) Run(cmd *exec.Cmd) error {
	var stdin *commandStdin
	if _, file := cmd.Stdin.(*os.File); cmd.Stdin != nil && !file {
		stdin = &commandStdin{Reader: cmd.Stdin, cmd: cmd}
		cmd.Stdin = stdin
	}
	err := cmd.Run()
	if stdin != nil && stdin.err != nil {
		err = stdin.err
	}
	if err != nil {
		return fmt.Errorf("'%v' %w: %w", cmd.Args, ErrCommandFailed, err)
	}
	return nil
}

func (c CmdRunner) Read(cmd *exec.Cmd) (io.ReadCloser, error) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		stdout.Close()
		return nil, fmt.Errorf("'%v' %w: %w", cmd.Args, ErrCommandFailed, err)
	}
	wait := sync.OnceValue(func() error {
		if err := cmd.Wait(); err != nil {
			return fmt.Errorf("'%v' %w: %w", cmd.Args, ErrCommandFailed, err)
		}
		return nil
	})
	return CmdReadCloser{ReadCloser: stdout, cmd: cmd, wait: wait}, nil
}

type commandStdin struct {
	io.Reader
	cmd *exec.Cmd
	err error
}

func (s *commandStdin) Read(p []byte) (int, error) {
	n, err := s.Reader.Read(p)
	if err != nil && err != io.EOF {
		// A failed upstream command must not leave the consumer waiting for user input.
		s.err = err
		s.cmd.Process.Kill()
	}
	return n, err
}

func (c CmdReadCloser) Read(p []byte) (int, error) {
	n, err := c.ReadCloser.Read(p)
	if err == io.EOF {
		// Exit status travels with stdout, rather than being lost in a deferred Close.
		if waitErr := c.wait(); waitErr != nil {
			return n, waitErr
		}
	}
	return n, err
}

func (c CmdReadCloser) Close() error {
	if c.cmd.ProcessState == nil {
		c.cmd.Process.Kill()
	}
	err := c.ReadCloser.Close()
	if errors.Is(err, os.ErrClosed) {
		err = nil
	}
	return errors.Join(err, c.wait())
}
