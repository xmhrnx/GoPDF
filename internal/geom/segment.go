package geom

// SegmentDist returns the distance from p to the segment a-b.
func SegmentDist(p, a, b Point) float64 {
	ab := b.Sub(a)
	l2 := ab.X*ab.X + ab.Y*ab.Y
	if l2 == 0 {
		return p.Dist(a)
	}
	t := ((p.X-a.X)*ab.X + (p.Y-a.Y)*ab.Y) / l2
	switch {
	case t < 0:
		t = 0
	case t > 1:
		t = 1
	}
	return p.Dist(a.Add(ab.Scale(t)))
}
