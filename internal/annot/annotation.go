package annot

import (
	"strings"

	"gopdf/internal/geom"
)

// Approximate Helvetica metrics used for text bounds before the text is
// typeset by the PDF writer.
const (
	avgGlyphWidth = 0.52
	lineSpacing   = 1.2
)

// Annotation is one piece of markup on a page. Coordinates are in page
// display space (points from the top-left of the page as shown).
type Annotation struct {
	ID       int64
	Page     int
	Kind     Kind
	Color    Color
	Width    float64 // stroke width in points
	Opacity  float64 // 0 < Opacity <= 1
	Points   []geom.Point
	Text     string
	FontSize float64
	Filled   bool
}

// Clone returns a deep copy of a.
func (a Annotation) Clone() Annotation {
	a.Points = append([]geom.Point(nil), a.Points...)
	return a
}

// Lines returns the text split into lines.
func (a Annotation) Lines() []string { return strings.Split(a.Text, "\n") }

// LineHeight returns the distance between text baselines.
func (a Annotation) LineHeight() float64 { return a.FontSize * lineSpacing }

// Bounds returns the area the annotation covers, stroke included.
func (a Annotation) Bounds() geom.Rect {
	if len(a.Points) == 0 {
		return geom.Rect{}
	}
	if a.Kind == Text {
		return a.textBounds()
	}
	r := geom.Rect{Min: a.Points[0], Max: a.Points[0]}
	for _, p := range a.Points[1:] {
		r = r.Union(geom.Rect{Min: p, Max: p})
	}
	return r.Inset(a.Width / 2)
}

func (a Annotation) textBounds() geom.Rect {
	widest := 1
	lines := a.Lines()
	for _, l := range lines {
		if n := len([]rune(l)); n > widest {
			widest = n
		}
	}
	w := float64(widest) * a.FontSize * avgGlyphWidth
	h := float64(len(lines)) * a.LineHeight()
	o := a.Points[0]
	return geom.Rect{Min: o, Max: geom.Point{X: o.X + w, Y: o.Y + h}}
}

// Translate moves every point by d.
func (a *Annotation) Translate(d geom.Point) {
	for i := range a.Points {
		a.Points[i] = a.Points[i].Add(d)
	}
}

// Box returns the normalized rectangle of a boxed annotation.
func (a Annotation) Box() geom.Rect {
	if len(a.Points) < 2 {
		return geom.Rect{}
	}
	return geom.RectFromPoints(a.Points[0], a.Points[len(a.Points)-1])
}
