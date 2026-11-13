package annot

import "gopdf/internal/geom"

// RotatePage turns the annotations on page a quarter turn clockwise (cw)
// or counter-clockwise, for a page whose display size before the turn is
// w x h. Text and signatures stay upright; their centers move with the
// page.
func RotatePage(list []Annotation, page int, cw bool, w, h float64) {
	for i := range list {
		a := &list[i]
		if a.Page != page {
			continue
		}
		if a.Kind.Placed() {
			rotateUpright(a, cw, w, h)
			continue
		}
		for j, p := range a.Points {
			a.Points[j] = rotatePoint(p, cw, w, h)
		}
	}
}

func rotatePoint(p geom.Point, cw bool, w, h float64) geom.Point {
	if cw {
		return geom.RotateCW(p, h)
	}
	return geom.RotateCCW(p, w)
}

func rotateUpright(a *Annotation, cw bool, w, h float64) {
	b := a.Bounds()
	c := geom.Point{X: (b.Min.X + b.Max.X) / 2, Y: (b.Min.Y + b.Max.Y) / 2}
	nc := rotatePoint(c, cw, w, h)
	d := geom.Point{X: nc.X - b.Dx()/2, Y: nc.Y - b.Dy()/2}.Sub(b.Min)
	a.Translate(d)
}
