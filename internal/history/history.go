package history

import (
	"cmp"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
)

func filePath() string {
	stateHome := os.Getenv("XDG_STATE_HOME")
	if !filepath.IsAbs(stateHome) {
		stateHome = filepath.Join(os.Getenv("HOME"), ".local", "state")
	}
	return filepath.Join(stateHome, "zipln", "history.json")
}

func load() map[string]int {
	data, _ := os.ReadFile(filePath())
	var counts map[string]int
	if err := json.Unmarshal(data, &counts); err != nil || counts == nil {
		return make(map[string]int)
	}
	return counts
}

func save(counts map[string]int) error {
	data, err := json.Marshal(counts)
	if err != nil {
		return err
	}
	filename := filePath()
	if err := os.MkdirAll(filepath.Dir(filename), 0700); err != nil {
		return err
	}
	temporary := filename + ".tmp"
	defer os.Remove(temporary)
	if err := os.WriteFile(temporary, data, 0600); err != nil {
		return err
	}
	return os.Rename(temporary, filename)
}

func Ranked() []string {
	counts := load()
	paths := make([]string, 0, len(counts))
	for path := range counts {
		paths = append(paths, path)
	}
	slices.SortFunc(paths, func(a, b string) int {
		if order := cmp.Compare(counts[b], counts[a]); order != 0 {
			return order
		}
		return cmp.Compare(a, b)
	})
	return paths
}

func Record(path string) {
	path, err := filepath.Abs(path)
	if err != nil {
		return
	}
	counts := load()
	counts[path]++
	_ = save(counts)
}
