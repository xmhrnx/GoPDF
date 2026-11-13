// Package recent remembers the files opened last, in the user's
// configuration folder.
package recent

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"time"
)

// maxEntries is how many files the list keeps.
const maxEntries = 8

// Entry is one recently opened file.
type Entry struct {
	Path   string    `json:"path"`
	Opened time.Time `json:"opened"`
	Pages  int       `json:"pages"`
}

// Name returns the file's base name.
func (e Entry) Name() string { return filepath.Base(e.Path) }

// Store reads and writes the list in a JSON file.
type Store struct {
	file string
}

// NewStore returns a store at the default location for app.
func NewStore(app string) (*Store, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	return &Store{file: filepath.Join(dir, app, "recent.json")}, nil
}

// Load returns the entries whose files still exist, newest first.
func (s *Store) Load() ([]Entry, error) {
	b, err := os.ReadFile(s.file)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var list []Entry
	if err := json.Unmarshal(b, &list); err != nil {
		return nil, err
	}
	return slices.DeleteFunc(list, func(e Entry) bool {
		_, err := os.Stat(e.Path)
		return err != nil
	}), nil
}

// Add records e at the top of the list and returns the new list.
func (s *Store) Add(list []Entry, e Entry) ([]Entry, error) {
	list = slices.DeleteFunc(slices.Clone(list), func(o Entry) bool { return o.Path == e.Path })
	list = slices.Insert(list, 0, e)
	if len(list) > maxEntries {
		list = list[:maxEntries]
	}
	return list, s.save(list)
}

func (s *Store) save(list []Entry) error {
	b, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(s.file), 0o755); err != nil {
		return err
	}
	return os.WriteFile(s.file, b, 0o644)
}
