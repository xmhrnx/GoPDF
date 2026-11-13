package annot

import "gopdf/internal/geom"

// Style is the look new annotations get from the active tool.
type Style struct {
	Color    Color
	Width    float64
	FontSize float64
	Filled   bool
}

// DefaultStyle is the tool style the editor starts with.
var DefaultStyle = Style{Color: Swatches[6], Width: 2, FontSize: 16}

// New returns an annotation of kind k started at p with style s.
func New(id int64, page int, k Kind, p geom.Point, s Style) Annotation {
	a := Annotation{
		ID: id, Page: page, Kind: k,
		Color: s.Color, Width: s.Width, Opacity: 1,
		FontSize: s.FontSize, Filled: s.Filled,
		Points: []geom.Point{p},
	}
	switch k {
	case Marker:
		a.Opacity, a.Width = 0.4, max(s.Width*5, 10)
	case Highlight:
		a.Opacity, a.Width, a.Filled = 0.4, 0, true
	case Whiteout:
		a.Color, a.Width, a.Filled = White, 0, true
	case Text:
		a.Text = "Text"
	}
	if !k.Freehand() && k != Text {
		a.Points = append(a.Points, p)
	}
	return a
}

// Extend updates the annotation being drawn as the pointer moves to p.
func (a *Annotation) Extend(p geom.Point) {
	if a.Kind.Freehand() {
		if last := a.Points[len(a.Points)-1]; last.Dist(p) >= 0.75 {
			a.Points = append(a.Points, p)
		}
		return
	}
	if len(a.Points) > 1 {
		a.Points[len(a.Points)-1] = p
	}
}

// Degenerate reports whether a freshly drawn annotation is too small to
// keep (a click without a drag).
func (a Annotation) Degenerate() bool {
	switch {
	case a.Kind == Text:
		return false
	case a.Kind.Freehand():
		return len(a.Points) < 2
	default:
		return a.Points[0].Dist(a.Points[len(a.Points)-1]) < 2
	}
}
