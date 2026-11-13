package geom

import "math"

// Rect is an axis-aligned rectangle; Min is the top-left corner in display
// space, or the lower-left corner in PDF user space.
type Rect struct {
	Min, Max Point
}

// RectFromPoints returns the normalized rectangle spanned by a and b.
func RectFromPoints(a, b Point) Rect {
	return Rect{
		Min: Point{math.Min(a.X, b.X), math.Min(a.Y, b.Y)},
		Max: Point{math.Max(a.X, b.X), math.Max(a.Y, b.Y)},
	}
}

// Dx returns the width of r.
func (r Rect) Dx() float64 { return r.Max.X - r.Min.X }

// Dy returns the height of r.
func (r Rect) Dy() float64 { return r.Max.Y - r.Min.Y }

// Inset grows r by d on every side (shrinks it for negative d).
func (r Rect) Inset(d float64) Rect {
	return Rect{Min: Point{r.Min.X - d, r.Min.Y - d}, Max: Point{r.Max.X + d, r.Max.Y + d}}
}

// Contains reports whether p lies inside r, edges included.
func (r Rect) Contains(p Point) bool {
	return p.X >= r.Min.X && p.X <= r.Max.X && p.Y >= r.Min.Y && p.Y <= r.Max.Y
}

// Union returns the smallest rectangle containing r and s.
func (r Rect) Union(s Rect) Rect {
	return Rect{
		Min: Point{math.Min(r.Min.X, s.Min.X), math.Min(r.Min.Y, s.Min.Y)},
		Max: Point{math.Max(r.Max.X, s.Max.X), math.Max(r.Max.Y, s.Max.Y)},
	}
}
