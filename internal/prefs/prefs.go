// Package prefs remembers the user's choices between runs: the style of
// each kind of markup and the signature in use.
package prefs

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"gopdf/internal/annot"
)

// Prefs are the remembered choices.
type Prefs struct {
	Styles        map[string]Style `json:"styles"`        // by annot.Kind name
	Signature     string           `json:"signature"`     // ID of the signature to place
	SignatureFont string           `json:"signatureFont"` // typeface last used for a typed signature
	DigitalID     string           `json:"digitalId"`     // fingerprint of the digital ID last signed with
}

// Style is an annot.Style as stored.
type Style struct {
	Color    string  `json:"color"`
	Width    float64 `json:"width"`
	FontSize float64 `json:"fontSize"`
	Filled   bool    `json:"filled"`
	Size     float64 `json:"size"`
}

// Store reads and writes Prefs in a JSON file.
type Store struct {
	file string
}

// NewStore returns a store in the configuration folder of app.
func NewStore(app string) (*Store, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	return &Store{file: filepath.Join(dir, app, "prefs.json")}, nil
}

// Load returns the saved prefs, or empty ones when there are none.
func (s *Store) Load() (Prefs, error) {
	var p Prefs
	b, err := os.ReadFile(s.file)
	if errors.Is(err, fs.ErrNotExist) {
		return p, nil
	}
	if err != nil {
		return p, err
	}
	return p, json.Unmarshal(b, &p)
}

// Save writes p.
func (s *Store) Save(p Prefs) error {
	b, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(s.file), 0o755); err != nil {
		return err
	}
	return os.WriteFile(s.file, b, 0o644)
}

// StyleFor returns the remembered style of k, or its default.
func (p Prefs) StyleFor(k annot.Kind) annot.Style {
	st := annot.DefaultStyleFor(k)
	saved, ok := p.Styles[k.String()]
	if !ok {
		return st
	}
	if c, ok := annot.ParseHex(saved.Color); ok {
		st.Color = c
	}
	if saved.Width > 0 {
		st.Width = saved.Width
	}
	if saved.FontSize > 0 {
		st.FontSize = saved.FontSize
	}
	if saved.Size > 0 {
		st.Size = saved.Size
	}
	st.Filled = saved.Filled
	return st
}

// SetStyle remembers st for k.
func (p *Prefs) SetStyle(k annot.Kind, st annot.Style) {
	if p.Styles == nil {
		p.Styles = map[string]Style{}
	}
	p.Styles[k.String()] = Style{Color: st.Color.Hex(), Width: st.Width, FontSize: st.FontSize, Filled: st.Filled, Size: st.Size}
}
