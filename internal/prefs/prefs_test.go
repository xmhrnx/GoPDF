package prefs

import (
	"path/filepath"
	"testing"

	"gopdf/internal/annot"
)

func TestStylesPerKind(t *testing.T) {
	s := &Store{file: filepath.Join(t.TempDir(), "a", "prefs.json")}
	var p Prefs
	pen := annot.DefaultStyleFor(annot.Pen)
	pen.Color, pen.Width = annot.Swatches[1], 8
	p.SetStyle(annot.Pen, pen)
	if err := s.Save(p); err != nil {
		t.Fatal(err)
	}
	got, err := s.Load()
	if err != nil {
		t.Fatal(err)
	}
	if st := got.StyleFor(annot.Pen); st != pen {
		t.Fatalf("pen = %+v, want %+v", st, pen)
	}
	if st := got.StyleFor(annot.Rect); st != annot.DefaultStyleFor(annot.Rect) {
		t.Fatalf("rect = %+v", st)
	}
}
