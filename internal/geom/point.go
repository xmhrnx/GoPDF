// Package geom holds the small geometry types shared by the document
// model and the editor: points and rectangles in PDF points (1/72 in).
package geom

import "math"

// Point is a position in page display space: PDF points measured from the
// top-left corner of the page as it is shown (after its /Rotate).
type Point struct {
	X, Y float64
}

// Add returns p translated by q.
func (p Point) Add(q Point) Point { return Point{p.X + q.X, p.Y + q.Y} }

// Sub returns the vector from q to p.
func (p Point) Sub(q Point) Point { return Point{p.X - q.X, p.Y - q.Y} }

// Scale returns p multiplied by s.
func (p Point) Scale(s float64) Point { return Point{p.X * s, p.Y * s} }

// Dist returns the distance between p and q.
func (p Point) Dist(q Point) float64 { return math.Hypot(p.X-q.X, p.Y-q.Y) }
