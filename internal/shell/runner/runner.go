package runner

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
)

var ErrCommandFailed = errors.New("command failed")

type Process interface {
	io.WriteCloser
	Wait() error
}

type Runner interface {
	Run(cmd *exec.Cmd) error
	Start(cmd *exec.Cmd) (Process, error)
	Read(cmd *exec.Cmd) (io.ReadCloser, error)
}

type CmdRunner struct{}
type CmdProcess struct {
	cmd   *exec.Cmd
	stdin io.WriteCloser
}

type CmdReadCloser struct {
	io.ReadCloser
	cmd *exec.Cmd
}

func (c CmdRunner) Run(cmd *exec.Cmd) error {
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("'%v' %w: %w", cmd.Args, ErrCommandFailed, err)
	}
	return nil
}

func (c CmdRunner) Start(cmd *exec.Cmd) (Process, error) {
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	err = cmd.Start()
	if err != nil {
		stdin.Close()
		return nil, fmt.Errorf("'%v' %w: %w", cmd.Args, ErrCommandFailed, err)
	}

	return CmdProcess{cmd: cmd, stdin: stdin}, nil
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
	return CmdReadCloser{ReadCloser: stdout, cmd: cmd}, nil
}

func (c CmdProcess) Write(p []byte) (n int, err error) {
	return c.stdin.Write(p)
}

func (c CmdProcess) Close() error {
	return c.stdin.Close()
}

func (c CmdProcess) Wait() error {
	return c.cmd.Wait()
}

func (c CmdReadCloser) Close() error {
	err := errors.Join(c.ReadCloser.Close(), c.cmd.Wait())
	if err != nil {
		return fmt.Errorf("'%v' %w: %w", c.cmd.Args, ErrCommandFailed, err)
	}
	return nil
}
