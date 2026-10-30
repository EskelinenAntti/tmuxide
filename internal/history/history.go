package history

import (
	"cmp"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
)

type selection struct {
	Path  string `json:"path"`
	Count int    `json:"count"`
}

func filePath() string {
	stateHome := os.Getenv("XDG_STATE_HOME")
	if !filepath.IsAbs(stateHome) {
		stateHome = filepath.Join(os.Getenv("HOME"), ".local", "state")
	}
	return filepath.Join(stateHome, "zipln", "queries.json")
}

func load() map[string][]selection {
	data, _ := os.ReadFile(filePath())
	var selections map[string][]selection
	if json.Unmarshal(data, &selections) != nil || selections == nil {
		return make(map[string][]selection)
	}
	return selections
}

func Preferred(query string) string {
	if selections := load()[query]; len(selections) > 0 {
		return selections[0].Path
	}
	return ""
}

func RecordQuery(query, path string) {
	if query == "" || path == "" {
		return
	}
	path, err := filepath.Abs(path)
	if err != nil {
		return
	}
	selections := load()
	choices := selections[query]
	count := 1
	for i, choice := range choices {
		if choice.Path == path {
			count += choice.Count
			choices = slices.Delete(choices, i, i+1)
			break
		}
	}
	// Stable sorting keeps the most recent selection first when counts tie.
	choices = append([]selection{{Path: path, Count: count}}, choices...)
	slices.SortStableFunc(choices, func(a, b selection) int {
		return cmp.Compare(b.Count, a.Count)
	})
	selections[query] = choices
	data, err := json.Marshal(selections)
	if err != nil {
		return
	}
	filename := filePath()
	if os.MkdirAll(filepath.Dir(filename), 0700) != nil {
		return
	}
	temporary, err := os.CreateTemp(filepath.Dir(filename), "queries-*.tmp")
	if err != nil {
		return
	}
	defer os.Remove(temporary.Name())
	_, writeErr := temporary.Write(data)
	closeErr := temporary.Close()
	if writeErr == nil && closeErr == nil {
		_ = os.Rename(temporary.Name(), filename)
	}
}
