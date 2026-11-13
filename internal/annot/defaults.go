package annot

import "gopdf/internal/geom"

// Style is the look new annotations get from the active tool.
type Style struct {
	Color    Color
	Width    float64 // stroke width
	FontSize float64 // text size
	Filled   bool
	Size     float64 // width of a placed signature
}

// Size limits, in points.
const (
	MinFontSize      = 4
	MaxFontSize      = 300
	MinSignatureSize = 24
	MaxSignatureSize = 600
)

// DefaultStyle is the tool style the editor starts with.
var DefaultStyle = Style{Color: Swatches[6], Width: 2, FontSize: 16, Size: 180}

// DefaultStyleFor returns the style kind starts with.
func DefaultStyleFor(k Kind) Style {
	s := DefaultStyle
	switch k {
	case Marker, Highlight:
		s.Color = Swatches[3]
	case Signature:
		s.Color, s.Width = Swatches[0], 2
	}
	return s
}

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
	if !k.Freehand() && !k.Placed() {
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
	case a.Kind.Placed():
		return false
	case a.Kind.Freehand():
		return len(a.Points) < 2
	default:
		return a.Points[0].Dist(a.Points[len(a.Points)-1]) < 2
	}
}

// NewSignature returns a signature of strokes (unit coordinates, aspect
// width/height) centered at p, s.Size points wide.
func NewSignature(id int64, page int, p geom.Point, strokes [][]geom.Point, aspect float64, s Style) Annotation {
	w := s.Size
	h := w / aspect
	a := Annotation{
		ID: id, Page: page, Kind: Signature, Color: s.Color, Width: s.Width, Opacity: 1,
		Points: []geom.Point{{X: p.X - w/2, Y: p.Y - h/2}, {X: p.X + w/2, Y: p.Y + h/2}},
	}
	a.Strokes = strokes
	return a.Clone()
}

// Resize sets a signature's width to w, keeping its aspect and center.
func (a *Annotation) Resize(w float64) {
	b := a.Box()
	if b.Dx() == 0 {
		return
	}
	h := w * b.Dy() / b.Dx()
	cx, cy := (b.Min.X+b.Max.X)/2, (b.Min.Y+b.Max.Y)/2
	a.Points = []geom.Point{{X: cx - w/2, Y: cy - h/2}, {X: cx + w/2, Y: cy + h/2}}
}
