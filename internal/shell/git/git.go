package git

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/eskelinenantti/tmuxide/internal/shell/runner"
)

type Cmd struct {
	runner.Runner
}

func (g Cmd) RevParse(cwd string) (string, error) {
	cmd := exec.Command("git", "-C", cwd, "rev-parse", "--show-toplevel")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := g.Run(cmd)
	return strings.TrimSpace(out.String()), err
}
