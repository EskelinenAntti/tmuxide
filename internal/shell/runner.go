package shell

import (
	"errors"
	"fmt"
	"os/exec"
)

var ErrCommandFailed = errors.New("command failed")

type Runner interface {
	Run(cmd exec.Cmd) error
}

type Parser interface {
	Parse() []string
}

type CmdRunner struct{}

func (c CmdRunner) Run(cmd exec.Cmd) error {
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%s %w %v: %w", cmd.Path, ErrCommandFailed, cmd.Args, err)
	}
	return nil
}
