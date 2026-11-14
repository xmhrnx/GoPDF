package prefs

import (
	"encoding/json"
	"testing"

	"gopdf/internal/annot"
)

func TestStylesPerKindSurviveJSON(t *testing.T) {
	var p Prefs
	pen := annot.DefaultStyleFor(annot.Pen)
	pen.Color, pen.Width = annot.Swatches[1], 8
	p.SetStyle(annot.Pen, pen)
	b, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	var got Prefs
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}
	if st := got.StyleFor(annot.Pen); st != pen {
		t.Fatalf("pen = %+v, want %+v", st, pen)
	}
	if st := got.StyleFor(annot.Rect); st != annot.DefaultStyleFor(annot.Rect) {
		t.Fatalf("rect = %+v", st)
	}
}
