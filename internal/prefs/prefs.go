// Package prefs is what the app remembers of the user's choices between
// runs: the style of each kind of markup and the signature in use.
package prefs

import "gopdf/internal/annot"

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
