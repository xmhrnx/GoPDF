package signature

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
)

// Store keeps signatures in a JSON file.
type Store struct {
	file string
}

// NewStore returns a store in the configuration folder of app.
func NewStore(app string) (*Store, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	return &Store{file: filepath.Join(dir, app, "signatures.json")}, nil
}

// Load returns the saved signatures, oldest first.
func (s *Store) Load() ([]Signature, error) {
	b, err := os.ReadFile(s.file)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var list []Signature
	return list, json.Unmarshal(b, &list)
}

// Save replaces the saved signatures with list.
func (s *Store) Save(list []Signature) error {
	b, err := json.Marshal(list)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(s.file), 0o755); err != nil {
		return err
	}
	return os.WriteFile(s.file, b, 0o600)
}

// Without returns list minus the signature with id.
func Without(list []Signature, id string) []Signature {
	return slices.DeleteFunc(slices.Clone(list), func(sig Signature) bool { return sig.ID == id })
}
