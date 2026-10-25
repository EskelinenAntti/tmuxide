package runner

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
)

var ErrCommandFailed = errors.New("command failed")

type Runner interface {
	Run(cmd *exec.Cmd) error
}

type Starter interface {
	Start(cmd *exec.Cmd) (io.WriteCloser, error)
}

type CmdRunner struct{}
type CmdWriteCloser struct {
	io.WriteCloser
	cmd *exec.Cmd
}

func (c CmdRunner) Run(cmd *exec.Cmd) error {
	return commandError(cmd, cmd.Run())
}

func commandError(cmd *exec.Cmd, err error) error {
	if err != nil {
		return fmt.Errorf("'%v' %w: %w", cmd.Args, ErrCommandFailed, err)
	}
	return nil
}

func (c CmdRunner) Start(cmd *exec.Cmd) (io.WriteCloser, error) {
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	err = cmd.Start()
	if err != nil {
		return nil, commandError(cmd, err)
	}

	return CmdWriteCloser{cmd: cmd, WriteCloser: stdin}, nil
}

func (c CmdWriteCloser) Close() error {
	err := c.WriteCloser.Close()
	if err != nil {
		return err
	}
	return c.cmd.Wait()
}
