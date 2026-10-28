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

type Process interface {
	io.WriteCloser
	Wait() error
	Kill() error
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
	cmd  *exec.Cmd
	wait func() error
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
	wait := sync.OnceValue(func() error {
		if err := cmd.Wait(); err != nil {
			return fmt.Errorf("'%v' %w: %w", cmd.Args, ErrCommandFailed, err)
		}
		return nil
	})
	return CmdReadCloser{ReadCloser: stdout, cmd: cmd, wait: wait}, nil
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

// Kill also closes stdin and reaps the command, unless Wait already did so.
func (c CmdProcess) Kill() error {
	c.stdin.Close()
	if c.cmd.ProcessState != nil {
		return nil
	}
	err := c.cmd.Process.Kill()
	c.cmd.Wait()
	if errors.Is(err, os.ErrProcessDone) {
		return nil
	}
	return err
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
