package ranked

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/eskelinenantti/tmuxide/internal/history"
)

func History() io.Reader {
	home := os.Getenv("HOME")
	var stdout strings.Builder
	for _, path := range history.Ranked() {
		fileInfo, err := os.Stat(path)
		if err != nil || strings.ContainsAny(path, "\r\n") {
			continue
		}
		if relative, err := filepath.Rel(home, path); err == nil && filepath.IsLocal(relative) {
			path = relative
		}
		if fileInfo.IsDir() {
			path += string(filepath.Separator)
		}
		stdout.WriteString(path)
		stdout.WriteString("\n")
	}
	return strings.NewReader(stdout.String())
}
