package history

import (
	"cmp"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"syscall"
)

type Store struct {
	path string
}

type Entry struct {
	Path  string
	Count uint64
}

func Default() Store {
	stateHome := os.Getenv("XDG_STATE_HOME")
	if !filepath.IsAbs(stateHome) {
		stateHome = filepath.Join(os.Getenv("HOME"), ".local", "state")
	}
	return Store{path: filepath.Join(stateHome, "tmuxide", "history.json")}
}

func (s Store) load() (map[string]uint64, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return make(map[string]uint64), nil
	}
	if err != nil {
		return nil, err
	}
	counts := make(map[string]uint64)
	if err := json.Unmarshal(data, &counts); err != nil {
		return nil, err
	}
	if counts == nil {
		counts = make(map[string]uint64)
	}
	return counts, nil
}

func (s Store) Ranked() ([]Entry, error) {
	counts, err := s.load()
	if err != nil {
		return nil, err
	}
	var entries []Entry
	for path, count := range counts {
		if count > 0 && filepath.IsAbs(path) {
			entries = append(entries, Entry{Path: path, Count: count})
		}
	}
	slices.SortFunc(entries, func(a, b Entry) int {
		if order := cmp.Compare(b.Count, a.Count); order != 0 {
			return order
		}
		return cmp.Compare(a.Path, b.Path)
	})
	return entries, nil
}

func (s Store) Record(path string) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0700); err != nil {
		return err
	}

	// Multiple pickers can be open at once. Lock the read-modify-write, and
	// replace the file atomically so readers never see a partial history.
	lock, err := os.OpenFile(s.path+".lock", os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer lock.Close()
	if err := syscall.Flock(int(lock.Fd()), syscall.LOCK_EX); err != nil {
		return err
	}
	defer syscall.Flock(int(lock.Fd()), syscall.LOCK_UN)

	counts, err := s.load()
	if err != nil {
		return err
	}
	counts[path]++
	data, err := json.Marshal(counts)
	if err != nil {
		return err
	}
	file, err := os.CreateTemp(filepath.Dir(s.path), ".history-*")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	if _, err := file.Write(data); err != nil {
		file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return os.Rename(file.Name(), s.path)
}
