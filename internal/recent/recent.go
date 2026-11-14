// Package recent keeps the list of files opened last.
package recent

import (
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

// Add puts e at the top of list, once, and caps its length.
func Add(list []Entry, e Entry) []Entry {
	list = slices.DeleteFunc(slices.Clone(list), func(o Entry) bool { return o.Path == e.Path })
	list = slices.Insert(list, 0, e)
	if len(list) > maxEntries {
		list = list[:maxEntries]
	}
	return list
}

// Existing returns the entries whose files are still there.
func Existing(list []Entry) []Entry {
	return slices.DeleteFunc(slices.Clone(list), func(e Entry) bool {
		_, err := os.Stat(e.Path)
		return err != nil
	})
}
