package ranked

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/eskelinenantti/herdr-zipln/internal/history"
)

// Filter preserves fzf's relevance order except for the last selection made for
// this exact query. Only promote it if it still occurs in the matching results.
func Filter(input io.Reader, output io.Writer, query string) error {
	cmd := exec.Command("fzf", "--filter", query, "--scheme=default", "--tiebreak=length,index", "--no-print-query")
	// Interactive user options (bindings, --no-sort, etc.) must not alter the
	// noninteractive ranking protocol.
	for _, variable := range os.Environ() {
		if !strings.HasPrefix(variable, "FZF_DEFAULT_OPTS=") && !strings.HasPrefix(variable, "FZF_DEFAULT_OPTS_FILE=") {
			cmd.Env = append(cmd.Env, variable)
		}
	}
	cmd.Stdin = input
	var matches bytes.Buffer
	cmd.Stdout = &matches
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		var exit *exec.ExitError
		if !errors.As(err, &exit) || exit.ExitCode() != 1 {
			return err
		}
	}
	preferred := history.Preferred(query)
	paths := strings.Split(strings.TrimSuffix(matches.String(), "\n"), "\n")
	for i, path := range paths {
		absolute := path
		if !filepath.IsAbs(absolute) {
			absolute = filepath.Join(os.Getenv("HOME"), absolute)
		}
		if preferred != "" && filepath.Clean(absolute) == preferred {
			copy(paths[1:i+1], paths[:i])
			paths[0] = path
			break
		}
	}
	if matches.Len() == 0 {
		return nil
	}
	_, err := io.WriteString(output, strings.Join(paths, "\n")+"\n")
	return err
}
