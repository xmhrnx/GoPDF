package annot

import (
	"testing"

	"gopdf/internal/geom"
)

func TestDrawAndHitRect(t *testing.T) {
	a := New(1, 0, Rect, geom.Point{X: 10, Y: 10}, DefaultStyle)
	a.Extend(geom.Point{X: 110, Y: 60})
	if a.Degenerate() {
		t.Fatal("dragged rectangle reported degenerate")
	}
	if !a.Hit(geom.Point{X: 10, Y: 30}) {
		t.Fatal("edge of outline not hit")
	}
	if a.Hit(geom.Point{X: 60, Y: 35}) {
		t.Fatal("inside of unfilled outline hit")
	}
}

func TestTopmostAt(t *testing.T) {
	list := []Annotation{
		New(1, 0, Highlight, geom.Point{X: 0, Y: 0}, DefaultStyle),
		New(2, 0, Highlight, geom.Point{X: 0, Y: 0}, DefaultStyle),
	}
	for i := range list {
		list[i].Extend(geom.Point{X: 50, Y: 50})
	}
	if got := TopmostAt(list, 0, geom.Point{X: 25, Y: 25}); got != 1 {
		t.Fatalf("TopmostAt = %d, want 1", got)
	}
	if got := TopmostAt(list, 1, geom.Point{X: 25, Y: 25}); got != -1 {
		t.Fatalf("other page = %d, want -1", got)
	}
}

func TestRotatePage(t *testing.T) {
	list := []Annotation{New(1, 0, Line, geom.Point{X: 0, Y: 0}, DefaultStyle)}
	list[0].Extend(geom.Point{X: 100, Y: 0})
	RotatePage(list, 0, true, 200, 300)
	want := []geom.Point{{X: 300, Y: 0}, {X: 300, Y: 100}}
	for i, p := range list[0].Points {
		if p != want[i] {
			t.Fatalf("point %d = %v, want %v", i, p, want[i])
		}
	}
}
