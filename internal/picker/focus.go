package picker

import (
	"bufio"
	"io"
	"os"
	"path/filepath"

	"github.com/eskelinenantti/herdr-zipln/internal/history"
)

// PreferredPosition locates the query's preferred path in fzf's ranked matches.
// Without a matching preference, the best fuzzy result stays selected.
func PreferredPosition(query string, matches io.Reader) (int, error) {
	preferred := history.Preferred(query)
	if preferred == "" {
		return 1, nil
	}
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
	return 1, scanner.Err()
}
