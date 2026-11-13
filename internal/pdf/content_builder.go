package pdf

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"

	"gopdf/internal/geom"
)

// kappa places Bézier control points to approximate a quarter ellipse.
const kappa = 0.5522847498

// contentBuilder writes PDF content-stream operators.
type contentBuilder struct {
	b strings.Builder
}

func num(v float64) string { return strconv.FormatFloat(v, 'f', -1, 64) }

func (cb *contentBuilder) op(format string, args ...any) {
	fmt.Fprintf(&cb.b, format, args...)
	cb.b.WriteByte('\n')
}

func (cb *contentBuilder) moveTo(p geom.Point) { cb.op("%s %s m", num(p.X), num(p.Y)) }
func (cb *contentBuilder) lineTo(p geom.Point) { cb.op("%s %s l", num(p.X), num(p.Y)) }

func (cb *contentBuilder) curveTo(c1, c2, p geom.Point) {
	cb.op("%s %s %s %s %s %s c", num(c1.X), num(c1.Y), num(c2.X), num(c2.Y), num(p.X), num(p.Y))
}

func (cb *contentBuilder) rect(r geom.Rect) {
	cb.op("%s %s %s %s re", num(r.Min.X), num(r.Min.Y), num(r.Dx()), num(r.Dy()))
}

func (cb *contentBuilder) polyline(pts []geom.Point) {
	cb.moveTo(pts[0])
	for _, p := range pts[1:] {
		cb.lineTo(p)
	}
}

func (cb *contentBuilder) ellipse(r geom.Rect) {
	cx, cy := (r.Min.X+r.Max.X)/2, (r.Min.Y+r.Max.Y)/2
	rx, ry := r.Dx()/2, r.Dy()/2
	kx, ky := rx*kappa, ry*kappa
	cb.moveTo(geom.Point{X: cx + rx, Y: cy})
	cb.curveTo(geom.Point{X: cx + rx, Y: cy + ky}, geom.Point{X: cx + kx, Y: cy + ry}, geom.Point{X: cx, Y: cy + ry})
	cb.curveTo(geom.Point{X: cx - kx, Y: cy + ry}, geom.Point{X: cx - rx, Y: cy + ky}, geom.Point{X: cx - rx, Y: cy})
	cb.curveTo(geom.Point{X: cx - rx, Y: cy - ky}, geom.Point{X: cx - kx, Y: cy - ry}, geom.Point{X: cx, Y: cy - ry})
	cb.curveTo(geom.Point{X: cx + kx, Y: cy - ry}, geom.Point{X: cx + rx, Y: cy - ky}, geom.Point{X: cx + rx, Y: cy})
	cb.op("h")
}

// pdfString encodes s for a WinAnsi Type 1 font as a literal string.
func pdfString(s string) string {
	enc := charmap.Windows1252.NewEncoder()
	var b strings.Builder
	b.WriteByte('(')
	for _, r := range s {
		bs, err := enc.Bytes([]byte(string(r)))
		if err != nil || len(bs) != 1 {
			bs = []byte{'?'}
		}
		switch c := bs[0]; c {
		case '(', ')', '\\':
			b.WriteByte('\\')
			b.WriteByte(c)
		default:
			if c < 32 || c > 126 {
				fmt.Fprintf(&b, "\\%03o", c)
			} else {
				b.WriteByte(c)
			}
		}
	}
	b.WriteByte(')')
	return b.String()
}

func (cb *contentBuilder) bytes() []byte { return []byte(cb.b.String()) }
