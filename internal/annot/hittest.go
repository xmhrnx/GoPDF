package annot

import "gopdf/internal/geom"

// minHitSlop is the smallest distance, in points, that still counts as a
// hit on a thin stroke.
const minHitSlop = 4.0

// Hit reports whether p (display space) touches the annotation.
func (a Annotation) Hit(p geom.Point) bool {
	slop := max(minHitSlop, a.Width)
	switch {
	case a.Kind.Placed(), a.Kind == Highlight, a.Kind == Whiteout, a.Filled:
		return a.Bounds().Inset(slop / 2).Contains(p)
	case a.Kind.Boxed():
		return hitBoxOutline(a.Box(), p, slop)
	default:
		return hitPolyline(a.Points, p, slop)
	}
}

func hitPolyline(pts []geom.Point, p geom.Point, slop float64) bool {
	if len(pts) == 1 {
		return pts[0].Dist(p) <= slop
	}
	for i := 1; i < len(pts); i++ {
		if geom.SegmentDist(p, pts[i-1], pts[i]) <= slop {
			return true
		}
	}
	return false
}

func hitBoxOutline(r geom.Rect, p geom.Point, slop float64) bool {
	return r.Inset(slop).Contains(p) && !r.Inset(-slop).Contains(p)
}

// TopmostAt returns the index of the last annotation on page that p
// touches, or -1.
func TopmostAt(list []Annotation, page int, p geom.Point) int {
	for i := len(list) - 1; i >= 0; i-- {
		if list[i].Page == page && list[i].Hit(p) {
			return i
		}
	}
	return -1
}
