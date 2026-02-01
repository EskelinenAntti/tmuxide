package shell

import (
	"errors"
	"fmt"
	"os/exec"
)

var ErrCommandFailed = errors.New("command failed")

type Runner interface {
	Run(cmd *exec.Cmd) error
	Start(cmd *exec.Cmd) error
}

type Parser interface {
	Parse() []string
}

type CmdRunner struct{}

func (c CmdRunner) Run(cmd *exec.Cmd) error {
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("'%v' %w: %w", cmd.Args, ErrCommandFailed, err)
	}
	return nil
}

func (c CmdRunner) Start(cmd *exec.Cmd) error {
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("'%v' %w: %w", cmd.Args, ErrCommandFailed, err)
	}
	return nil
}
