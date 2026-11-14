package pdf

import "gopdf/internal/geom"

// apply maps p by the PDF matrix m.
func apply(m [6]float64, p geom.Point) geom.Point {
	return geom.Point{X: m[0]*p.X + m[2]*p.Y + m[4], Y: m[1]*p.X + m[3]*p.Y + m[5]}
}

// invert returns the inverse of m, which must not be singular.
func invert(m [6]float64) [6]float64 {
	det := m[0]*m[3] - m[1]*m[2]
	a, b, c, d := m[3]/det, -m[1]/det, -m[2]/det, m[0]/det
	return [6]float64{a, b, c, d, -(a*m[4] + c*m[5]), -(b*m[4] + d*m[5])}
}

// mapRect returns the box around r's corners mapped by m.
func mapRect(m [6]float64, r geom.Rect) geom.Rect {
	out := geom.RectFromPoints(apply(m, r.Min), apply(m, r.Max))
	return out.Union(geom.RectFromPoints(
		apply(m, geom.Point{X: r.Min.X, Y: r.Max.Y}),
		apply(m, geom.Point{X: r.Max.X, Y: r.Min.Y}),
	))
}

// ToDisplay maps a rectangle in the page's user space to display space.
func (p PageInfo) ToDisplay(r geom.Rect) geom.Rect { return mapRect(invert(p.displayMatrix()), r) }

// ToUser maps a rectangle in display space to the page's user space.
func (p PageInfo) ToUser(r geom.Rect) geom.Rect { return mapRect(p.displayMatrix(), r) }
