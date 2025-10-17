package git

import (
	"fmt"
	"os/exec"
	"strings"
)

type Repository struct{}

func (Repository) Root(target string) (string, error) {
	cmd := exec.Command("git", "-C", target, "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("%s is not a valid Git repository", target)
	}
	return strings.TrimSpace(string(out)), nil
}
