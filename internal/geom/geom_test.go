package geom

import "testing"

func TestRotateRoundTrip(t *testing.T) {
	const w, h = 200.0, 100.0
	p := Point{30, 20}
	q := RotateCW(p, h) // the rotated page is h wide and w tall
	if got := RotateCCW(q, h); got != p {
		t.Fatalf("round trip = %v, want %v", got, p)
	}
	if q.X < 0 || q.X > h || q.Y < 0 || q.Y > w {
		t.Fatalf("rotated point %v outside %vx%v", q, h, w)
	}
}

func TestSegmentDist(t *testing.T) {
	if d := SegmentDist(Point{5, 5}, Point{0, 0}, Point{10, 0}); d != 5 {
		t.Fatalf("dist = %v, want 5", d)
	}
	if d := SegmentDist(Point{-3, 4}, Point{0, 0}, Point{10, 0}); d != 5 {
		t.Fatalf("dist past end = %v, want 5", d)
	}
}
