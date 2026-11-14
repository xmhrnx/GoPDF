package signature

import (
	"testing"

	"gopdf/internal/geom"
)

func TestFromStrokesNormalizes(t *testing.T) {
	sig, ok := FromStrokes([][]geom.Point{
		{{X: 100, Y: 50}, {X: 300, Y: 50}},
		{{X: 100, Y: 100}, {X: 300, Y: 100}},
	})
	if !ok {
		t.Fatal("signature rejected")
	}
	for _, st := range sig.Strokes {
		for _, p := range st {
			if p.X < 0 || p.X > 1 || p.Y < 0 || p.Y > 1 {
				t.Fatalf("point %v outside the unit box", p)
			}
		}
	}
	if sig.Aspect < 3 || sig.Aspect > 4.5 {
		t.Fatalf("aspect = %v", sig.Aspect)
	}
	if _, ok := FromStrokes([][]geom.Point{{{X: 1, Y: 1}}}); ok {
		t.Fatal("a dot became a signature")
	}
}

func TestWithout(t *testing.T) {
	sig, _ := FromStrokes([][]geom.Point{{{X: 0, Y: 0}, {X: 40, Y: 10}}})
	if len(Without([]Signature{sig}, sig.ID)) != 0 {
		t.Fatal("Without kept the signature")
	}
}
