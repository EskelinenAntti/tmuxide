package shell

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
)

var ErrCommandFailed = errors.New("command failed")

type WriteCloser interface {
	Close() error
	Write(p []byte) (n int, err error)
}

type Runner interface {
	Run(cmd *exec.Cmd) error
	Start(cmd *exec.Cmd) (WriteCloser, error)
}

type Parser interface {
	Parse() []string
}

type CmdRunner struct{}
type CmdWriteCloser struct {
	cmd   *exec.Cmd
	stdin io.WriteCloser
}

func (c CmdRunner) Run(cmd *exec.Cmd) error {
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("'%v' %w: %w", cmd.Args, ErrCommandFailed, err)
	}
	return nil
}

func (c CmdRunner) Start(cmd *exec.Cmd) (WriteCloser, error) {
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("'%v' %w: %w", cmd.Args, ErrCommandFailed, err)
	}

	return CmdWriteCloser{cmd: cmd, stdin: stdin}, nil
}

func (c CmdWriteCloser) Write(p []byte) (n int, err error) {
	return c.stdin.Write(p)
}

func (c CmdWriteCloser) Close() error {
	err := c.stdin.Close()
	if err != nil {
		return err
	}
	return c.cmd.Wait()
}
