// Package signature turns hand-drawn strokes into reusable signatures
// and keeps them in the user's configuration folder.
package signature

import (
	"crypto/rand"
	"encoding/hex"
	"math"
	"time"

	"gopdf/internal/geom"
)

// Normalization constants.
const (
	minInk    = 2.0  // smallest drawing, in input units, that counts as a signature
	padding   = 0.04 // margin around the ink, as a fraction of the larger side
	minAspect = 0.2  // narrowest width/height kept (a tall scribble)
	maxPaths  = 4096 // a sanity cap on stored strokes and outlines
)

// Signature is a saved signature in unit coordinates of a box whose
// width/height ratio is Aspect: pen Strokes (drawn ones) and filled
// glyph outlines, Fills (typed ones, closed polygons filled with the
// non-zero rule).
type Signature struct {
	ID      string         `json:"id"`
	Strokes [][]geom.Point `json:"strokes,omitempty"`
	Fills   [][]geom.Point `json:"fills,omitempty"`
	Aspect  float64        `json:"aspect"`
	// DigitalID is the fingerprint of the digital ID that signs when this
	// signature is placed; empty for a mark without one.
	DigitalID string    `json:"digitalId,omitempty"`
	Created   time.Time `json:"created"`
}

// FromStrokes builds a signature from strokes drawn in any units (dp on
// a pad). It reports false when there is too little ink.
func FromStrokes(raw [][]geom.Point) (Signature, bool) {
	return normalize(raw, nil)
}

// normalize fits strokes and fills, in any units, into a padded unit
// box. It reports false when there is too little ink.
func normalize(strokes, fills [][]geom.Point) (Signature, bool) {
	strokes, fills = capPaths(strokes), capPaths(fills)
	box, ok := bounds(append(append([][]geom.Point(nil), strokes...), fills...))
	if !ok || max(box.Dx(), box.Dy()) < minInk {
		return Signature{}, false
	}
	pad := padding * max(box.Dx(), box.Dy())
	box = box.Inset(pad)
	w, h := box.Dx(), box.Dy()
	if w/h < minAspect {
		grow := (h*minAspect - w) / 2
		box.Min.X, box.Max.X, w = box.Min.X-grow, box.Max.X+grow, h*minAspect
	}
	toUnit := func(paths [][]geom.Point) [][]geom.Point {
		var out [][]geom.Point
		for _, path := range paths {
			unit := make([]geom.Point, len(path))
			for i, p := range path {
				unit[i] = geom.Point{X: round((p.X - box.Min.X) / w), Y: round((p.Y - box.Min.Y) / h)}
			}
			out = append(out, unit)
		}
		return out
	}
	return Signature{
		ID: newID(), Aspect: w / h, Created: time.Now(),
		Strokes: toUnit(strokes), Fills: toUnit(fills),
	}, true
}

// capPaths drops empty paths and keeps at most maxPaths.
func capPaths(paths [][]geom.Point) [][]geom.Point {
	var kept [][]geom.Point
	for _, p := range paths {
		if len(p) > 0 && len(kept) < maxPaths {
			kept = append(kept, p)
		}
	}
	return kept
}

// bounds returns the box around every point of paths.
func bounds(paths [][]geom.Point) (geom.Rect, bool) {
	var (
		box  geom.Rect
		seen bool
	)
	for _, path := range paths {
		for _, p := range path {
			if !seen {
				box, seen = geom.Rect{Min: p, Max: p}, true
			}
			box = box.Union(geom.Rect{Min: p, Max: p})
		}
	}
	return box, seen
}

// round keeps four decimals, plenty for ink and small on disk.
func round(v float64) float64 { return math.Round(v*1e4) / 1e4 }

func newID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return time.Now().Format("20060102150405.000000000")
	}
	return hex.EncodeToString(b)
}
