package history

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func filePath() string {
	stateHome := os.Getenv("XDG_STATE_HOME")
	if !filepath.IsAbs(stateHome) {
		stateHome = filepath.Join(os.Getenv("HOME"), ".local", "state")
	}
	return filepath.Join(stateHome, "zipln", "queries.json")
}

func load() map[string]string {
	data, _ := os.ReadFile(filePath())
	var selections map[string]string
	if json.Unmarshal(data, &selections) != nil || selections == nil {
		return make(map[string]string)
	}
	return selections
}

func Preferred(query string) string {
	return load()[query]
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
	selections[query] = path
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
