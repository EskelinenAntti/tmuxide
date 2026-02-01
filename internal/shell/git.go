package shell

import (
	"bytes"
	"os/exec"
	"strings"
)

type GitCmd struct {
	Runner
}

func (g GitCmd) RevParse(cwd string) (string, error) {
	cmd := exec.Command("git", "-C", cwd, "rev-parse", "--show-toplevel")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := g.Run(cmd)
	return strings.TrimSpace(out.String()), err
}
