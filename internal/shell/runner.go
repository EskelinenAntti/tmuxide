package shell

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

var ErrSubCmd = errors.New("command failed")

type CmdRunner struct {
	Command string
}

type Parser interface {
	Parse() []string
}

type Runner interface {
	Run(name string, args Parser) error
	Attach(name string, args Parser) error
	Output(name string, args Parser) ([]byte, error)
	Pipe(name string, args Parser, input io.Reader) ([]byte, error)
}

func (c CmdRunner) Attach(subCommand string, args Parser) error {
	cmd := c.createCmd(subCommand, args)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return c.run(cmd)
}

func (c CmdRunner) Run(subCommand string, args Parser) error {
	cmd := c.createCmd(subCommand, args)
	return c.run(cmd)
}

func (c CmdRunner) Output(subCommand string, args Parser) ([]byte, error) {
	cmd := c.createCmd(subCommand, args)
	return cmd.CombinedOutput()
}

func (c CmdRunner) Pipe(subCmd string, args Parser, input io.Reader) ([]byte, error) {
	cmd := c.createCmd(subCmd, args)
	cmd.Stdin = input
	cmd.Stderr = os.Stderr
	return cmd.Output()
}

func (c CmdRunner) run(cmd *exec.Cmd) error {
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%s %w %v: %w", c.Command, ErrSubCmd, cmd.Args, err)
	}
	return nil
}

func (c CmdRunner) createCmd(subCommand string, args Parser) *exec.Cmd {
	cmd := exec.Command(c.Command)
	if subCommand != "" {
		cmd.Args = append(cmd.Args, subCommand)
	}
	cmd.Args = append(cmd.Args, args.Parse()...)
	return cmd
}
