package picker

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
)

// PreferredPosition locates a preferred absolute path in fzf's ranked matches.
// It returns zero if the path has not arrived in the stream yet.
func PreferredPosition(preferred string, matches io.Reader) (int, error) {
	scanner := bufio.NewScanner(matches)
	for position := 1; scanner.Scan(); position++ {
		path := scanner.Text()
		if !filepath.IsAbs(path) {
			path = filepath.Join(os.Getenv("HOME"), path)
		}
		if filepath.Clean(path) == preferred {
			return position, nil
		}
	}
	return 0, scanner.Err()
}
