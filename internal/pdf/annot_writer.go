package pdf

import (
	"math"

	"gopdf/internal/annot"
	"gopdf/internal/geom"
)

// Names of the resources the flattened markup uses.
const (
	fontResName = "PSHelv"
	gsResPrefix = "PSGs"
	// arrowHeadLen is the arrow head length in multiples of the stroke width.
	arrowHeadLen  = 4.0
	minArrowHead  = 8.0
	arrowHeadSpan = math.Pi / 7
	// helvAscent is Helvetica's ascent as a fraction of the font size.
	helvAscent = 0.718
)

// writeAnnotations emits list (all on one page) as content operators in
// display space, which m maps to user space. gs names the ExtGState for
// each opacity in use.
func writeAnnotations(list []annot.Annotation, m [6]float64, gs func(opacity float64) string) []byte {
	var cb contentBuilder
	cb.op("q")
	cb.op("%s %s %s %s %s %s cm", num(m[0]), num(m[1]), num(m[2]), num(m[3]), num(m[4]), num(m[5]))
	for _, a := range list {
		cb.op("q")
		if a.Opacity < 1 {
			cb.op("/%s gs", gs(a.Opacity))
		}
		writeAnnotation(&cb, a)
		cb.op("Q")
	}
	cb.op("Q")
	return cb.bytes()
}

func writeAnnotation(cb *contentBuilder, a annot.Annotation) {
	r, g, b := a.Color.Floats()
	cb.op("%s %s %s RG %s %s %s rg", num(r), num(g), num(b), num(r), num(g), num(b))
	cb.op("%s w 1 J 1 j", num(a.Width))
	switch {
	case a.Kind == annot.Text:
		writeText(cb, a)
	case a.Kind == annot.Signature:
		for _, st := range a.Strokes {
			if len(st) < 2 {
				continue
			}
			pts := make([]geom.Point, len(st))
			for i, u := range st {
				pts[i] = a.StrokePoint(u)
			}
			cb.polyline(pts)
			cb.op("S")
		}
	case a.Kind.Freehand():
		cb.polyline(a.Points)
		cb.op("S")
	case a.Kind == annot.Ellipse:
		cb.ellipse(a.Box())
		paint(cb, a)
	case a.Kind.Boxed():
		cb.rect(a.Box())
		paint(cb, a)
	case a.Kind == annot.Arrow:
		writeArrow(cb, a)
	default:
		cb.polyline(a.Points)
		cb.op("S")
	}
}

func paint(cb *contentBuilder, a annot.Annotation) {
	switch {
	case a.Filled && a.Width > 0 && a.Kind != annot.Highlight && a.Kind != annot.Whiteout:
		cb.op("B")
	case a.Filled:
		cb.op("f")
	default:
		cb.op("S")
	}
}

func writeArrow(cb *contentBuilder, a annot.Annotation) {
	from, to := a.Points[0], a.Points[len(a.Points)-1]
	cb.polyline([]geom.Point{from, to})
	cb.op("S")
	for _, p := range ArrowHead(from, to, a.Width) {
		cb.polyline([]geom.Point{p, to})
		cb.op("S")
	}
}

// ArrowHead returns the two outer points of an arrow head at to, for a
// shaft from from drawn with stroke width w.
func ArrowHead(from, to geom.Point, w float64) [2]geom.Point {
	l := max(minArrowHead, w*arrowHeadLen)
	ang := math.Atan2(to.Y-from.Y, to.X-from.X)
	return [2]geom.Point{
		{X: to.X - l*math.Cos(ang-arrowHeadSpan), Y: to.Y - l*math.Sin(ang-arrowHeadSpan)},
		{X: to.X - l*math.Cos(ang+arrowHeadSpan), Y: to.Y - l*math.Sin(ang+arrowHeadSpan)},
	}
}

func writeText(cb *contentBuilder, a annot.Annotation) {
	o := a.Points[0]
	cb.op("BT")
	cb.op("/%s %s Tf", fontResName, num(a.FontSize))
	// The text matrix flips y back up so glyphs read upright in the
	// y-down display space set by the outer cm.
	cb.op("1 0 0 -1 %s %s Tm", num(o.X), num(o.Y+a.FontSize*helvAscent+(a.LineHeight()-a.FontSize)/2))
	for i, line := range a.Lines() {
		if i > 0 {
			cb.op("0 %s Td", num(-a.LineHeight()))
		}
		cb.op("%s Tj", pdfString(line))
	}
	cb.op("ET")
}
