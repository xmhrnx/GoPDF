package studio

import (
	"fmt"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/styles/properties"

	"gopdf/internal/annot"
	"gopdf/internal/geom"
	"gopdf/internal/pdf"
	"gopdf/ui/theme"
)

// selectionPad is the gap between markup and its selection frame, in dp.
const selectionPad = 4

// markupLayer draws the page's markup (all but text, which is laid out
// as text) over a w x h dp page shown at zoom z.
func markupLayer(s *Studio, page, w, h int, z float64) giode.Element {
	list := s.doc.Annots()
	drawing := s.drawing
	selected := s.selected
	return giode.Box(giode.Styles{Position: giode.PositionAbsolute, Top: giode.Flush, Left: giode.Flush},
		giode.Canvas(giode.Styles{Width: properties.Width(w), Height: properties.Height(h)}, func(p *giode.Painter) {
			for _, a := range list {
				if a.Page == page {
					paintAnnotation(p, a, z)
				}
			}
			if drawing != nil && drawing.Page == page {
				paintAnnotation(p, *drawing, z)
			}
			for _, a := range list {
				if a.ID == selected && a.Page == page {
					paintSelection(p, a.Bounds(), z)
				}
			}
		}))
}

// rgba returns the annotation's color with its opacity as CSS hex.
func rgba(a annot.Annotation) string {
	return fmt.Sprintf("%s%02x", a.Color.Hex(), uint8(a.Opacity*255))
}

func paintAnnotation(p *giode.Painter, a annot.Annotation, z float64) {
	col, sw := rgba(a), float32(a.Width*z)
	switch {
	case a.Kind == annot.Text:
		return
	case a.Kind.Freehand():
		path := p.Path()
		path.MoveTo(f(a.Points[0].X, z), f(a.Points[0].Y, z))
		for _, pt := range a.Points[1:] {
			path.LineTo(f(pt.X, z), f(pt.Y, z))
		}
		path.Stroke(sw, col)
	case a.Kind == annot.Ellipse:
		paintEllipse(p, a.Box(), z, a.Filled, sw, col)
	case a.Kind.Boxed():
		paintBox(p, a, z, sw, col)
	case a.Kind == annot.Arrow:
		from, to := a.Points[0], a.Points[len(a.Points)-1]
		segment(p, from, to, z, sw, col)
		for _, h := range pdf.ArrowHead(from, to, a.Width) {
			segment(p, h, to, z, sw, col)
		}
	default:
		segment(p, a.Points[0], a.Points[len(a.Points)-1], z, sw, col)
	}
}

func f(v, z float64) float32 { return float32(v * z) }

func segment(p *giode.Painter, a, b geom.Point, z float64, w float32, col string) {
	p.Line(f(a.X, z), f(a.Y, z), f(b.X, z), f(b.Y, z), w, col)
}

func paintBox(p *giode.Painter, a annot.Annotation, z float64, sw float32, col string) {
	r := a.Box()
	if a.Filled {
		p.Rect(f(r.Min.X, z), f(r.Min.Y, z), f(r.Dx(), z), f(r.Dy(), z), 0, col)
	}
	if sw > 0 && a.Kind == annot.Rect {
		outline(p, r, z, sw, col)
	}
}

func outline(p *giode.Painter, r geom.Rect, z float64, sw float32, col string) {
	path := p.Path()
	path.MoveTo(f(r.Min.X, z), f(r.Min.Y, z))
	path.LineTo(f(r.Max.X, z), f(r.Min.Y, z))
	path.LineTo(f(r.Max.X, z), f(r.Max.Y, z))
	path.LineTo(f(r.Min.X, z), f(r.Max.Y, z))
	path.Close()
	path.Stroke(sw, col)
}

func paintEllipse(p *giode.Painter, r geom.Rect, z float64, filled bool, sw float32, col string) {
	const k = 0.5522847498
	cx, cy := f((r.Min.X+r.Max.X)/2, z), f((r.Min.Y+r.Max.Y)/2, z)
	rx, ry := f(r.Dx()/2, z), f(r.Dy()/2, z)
	path := p.Path()
	path.MoveTo(cx+rx, cy)
	path.CubeTo(cx+rx, cy+ry*k, cx+rx*k, cy+ry, cx, cy+ry)
	path.CubeTo(cx-rx*k, cy+ry, cx-rx, cy+ry*k, cx-rx, cy)
	path.CubeTo(cx-rx, cy-ry*k, cx-rx*k, cy-ry, cx, cy-ry)
	path.CubeTo(cx+rx*k, cy-ry, cx+rx, cy-ry*k, cx+rx, cy)
	path.Close()
	if filled {
		path.Fill(col)
		return
	}
	path.Stroke(sw, col)
}

// paintSelection frames the selected markup with corner handles.
func paintSelection(p *giode.Painter, b geom.Rect, z float64) {
	pad := float32(selectionPad)
	x0, y0 := f(b.Min.X, z)-pad, f(b.Min.Y, z)-pad
	x1, y1 := f(b.Max.X, z)+pad, f(b.Max.Y, z)+pad
	path := p.Path()
	path.MoveTo(x0, y0)
	path.LineTo(x1, y0)
	path.LineTo(x1, y1)
	path.LineTo(x0, y1)
	path.Close()
	path.Stroke(1.5, theme.Selection)
	for _, c := range [][2]float32{{x0, y0}, {x1, y0}, {x1, y1}, {x0, y1}} {
		p.Rect(c[0]-4, c[1]-4, 8, 8, 2, theme.Selection)
		p.Rect(c[0]-2.5, c[1]-2.5, 5, 5, 1, "#ffffff")
	}
}
