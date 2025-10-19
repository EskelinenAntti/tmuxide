package git

import (
	"os/exec"
	"strings"
)

type Command interface {
	RevParse(cwd string) (string, error)
}

type ShellGit struct{}

func (ShellGit) RevParse(cwd string) (string, error) {
	cmd := exec.Command("git", "-C", cwd, "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
