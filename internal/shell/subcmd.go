package shell

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

var ErrSubCmd = errors.New("command failed")

type SubCmdRunner struct {
	Command string
}

func (c SubCmdRunner) Attach(name string, args Parser) error {
	cmd := c.createCmd(name, args)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return c.runCmd(cmd)
}

func (c SubCmdRunner) Run(name string, args Parser) error {
	cmd := c.createCmd(name, args)
	return c.runCmd(cmd)
}

func (c SubCmdRunner) runCmd(cmd *exec.Cmd) error {
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%s %w %v: %w", c.Command, ErrSubCmd, cmd.Args, err)
	}
	return nil
}

func (c SubCmdRunner) createCmd(name string, args Parser) *exec.Cmd {
	cmd := exec.Command("tmux", name)
	cmd.Args = append(cmd.Args, args.Parse()...)
	return cmd
}
