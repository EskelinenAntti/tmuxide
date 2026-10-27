package history

import (
	"cmp"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
)

func filePath() string {
	stateHome := os.Getenv("XDG_STATE_HOME")
	if !filepath.IsAbs(stateHome) {
		stateHome = filepath.Join(os.Getenv("HOME"), ".local", "state")
	}
	return filepath.Join(stateHome, "tmuxide", "history.json")
}

func load() (map[string]int, error) {
	data, err := os.ReadFile(filePath())
	if errors.Is(err, os.ErrNotExist) {
		return make(map[string]int), nil
	}
	if err != nil {
		return nil, err
	}
	var counts map[string]int
	err = json.Unmarshal(data, &counts)
	if counts == nil {
		counts = make(map[string]int)
	}
	return counts, err
}

func Ranked() ([]string, error) {
	counts, err := load()
	if err != nil {
		return nil, err
	}
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
	return paths, nil
}

func Record(path string) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	counts, err := load()
	if err != nil {
		return err
	}
	counts[path]++
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
