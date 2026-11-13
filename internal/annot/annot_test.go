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

func TestSignaturePlacement(t *testing.T) {
	st := DefaultStyleFor(Signature)
	strokes := [][]geom.Point{{{X: 0, Y: 0}, {X: 1, Y: 1}}}
	a := NewSignature(1, 0, geom.Point{X: 100, Y: 100}, strokes, 3, st)
	if b := a.Box(); b.Dx() != 180 || b.Dy() != 60 {
		t.Fatalf("box = %v", b)
	}
	if p := a.StrokePoint(geom.Point{X: 1, Y: 1}); p != (geom.Point{X: 190, Y: 130}) {
		t.Fatalf("stroke end = %v", p)
	}
	c := a.Clone()
	c.Strokes[0][0].X = 0.5
	if a.Strokes[0][0].X != 0 {
		t.Fatal("clone shares strokes")
	}
	a.Resize(90)
	if b := a.Box(); b.Dx() != 90 || b.Dy() != 30 || !a.Hit(geom.Point{X: 100, Y: 100}) {
		t.Fatalf("resized box = %v", b)
	}
	RotatePage([]Annotation{a}, 0, true, 600, 800)
}

func TestParseHex(t *testing.T) {
	if c, ok := ParseHex("#1E293b"); !ok || c != (Color{0x1e, 0x29, 0x3b}) {
		t.Fatalf("got %v %v", c, ok)
	}
	if _, ok := ParseHex("#12345"); ok {
		t.Fatal("short hex accepted")
	}
}
