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
	minInk     = 2.0  // smallest drawing, in input units, that counts as a signature
	padding    = 0.04 // margin around the ink, as a fraction of the larger side
	minAspect  = 0.2  // narrowest width/height kept (a tall scribble)
	maxStrokes = 512  // a sanity cap on stored strokes
)

// Signature is a saved signature: strokes in unit coordinates of a box
// whose width/height ratio is Aspect.
type Signature struct {
	ID      string         `json:"id"`
	Strokes [][]geom.Point `json:"strokes"`
	Aspect  float64        `json:"aspect"`
	Created time.Time      `json:"created"`
}

// FromStrokes builds a signature from strokes drawn in any units (dp on
// a pad). It reports false when there is too little ink.
func FromStrokes(raw [][]geom.Point) (Signature, bool) {
	var (
		box  geom.Rect
		seen bool
		kept [][]geom.Point
	)
	for _, st := range raw {
		if len(st) == 0 || len(kept) == maxStrokes {
			continue
		}
		kept = append(kept, st)
		for _, p := range st {
			if !seen {
				box, seen = geom.Rect{Min: p, Max: p}, true
			}
			box = box.Union(geom.Rect{Min: p, Max: p})
		}
	}
	if !seen || max(box.Dx(), box.Dy()) < minInk {
		return Signature{}, false
	}
	pad := padding * max(box.Dx(), box.Dy())
	box = box.Inset(pad)
	w, h := box.Dx(), box.Dy()
	if w/h < minAspect {
		grow := (h*minAspect - w) / 2
		box.Min.X, box.Max.X, w = box.Min.X-grow, box.Max.X+grow, h*minAspect
	}
	sig := Signature{ID: newID(), Aspect: w / h, Created: time.Now()}
	for _, st := range kept {
		unit := make([]geom.Point, len(st))
		for i, p := range st {
			unit[i] = geom.Point{X: round((p.X - box.Min.X) / w), Y: round((p.Y - box.Min.Y) / h)}
		}
		sig.Strokes = append(sig.Strokes, unit)
	}
	return sig, true
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
