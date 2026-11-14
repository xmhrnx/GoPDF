package editor

import (
	"testing"

	"gopdf/internal/annot"
	"gopdf/internal/geom"
	"gopdf/internal/pdf"
)

func newTestSession(t *testing.T, pages int) *Session {
	t.Helper()
	data := samplePDF(pages)
	doc, err := pdf.Load(data)
	if err != nil {
		t.Fatal(err)
	}
	return NewSession(doc, "/tmp/test.pdf", "")
}

func TestDuplicateCarriesMarkup(t *testing.T) {
	s := newTestSession(t, 2)
	a := annot.New(NewID(), 1, annot.Rect, geom.Point{X: 1, Y: 1}, annot.DefaultStyle)
	a.Extend(geom.Point{X: 50, Y: 50})
	s.AddAnnotation(a)
	next, err := Duplicate(s.Snapshot(), 1)
	if err != nil {
		t.Fatal(err)
	}
	s.Commit(next)
	if s.Doc().PageCount() != 3 || len(s.Annots()) != 2 {
		t.Fatalf("pages %d annots %d", s.Doc().PageCount(), len(s.Annots()))
	}
	if s.Annots()[0].ID == s.Annots()[1].ID || s.Annots()[1].Page != 2 {
		t.Fatalf("copy = %+v", s.Annots()[1])
	}
	next, err = Delete(s.Snapshot(), []int{0})
	if err != nil {
		t.Fatal(err)
	}
	s.Commit(next)
	if got := s.Annots()[0].Page; got != 0 {
		t.Fatalf("after delete page = %d", got)
	}
	s.Undo()
	s.Undo()
	if s.Doc().PageCount() != 2 || len(s.Annots()) != 1 {
		t.Fatal("undo did not restore the document")
	}
}

func TestDragUndoesInOneStep(t *testing.T) {
	s := newTestSession(t, 1)
	a := annot.New(NewID(), 0, annot.Line, geom.Point{}, annot.DefaultStyle)
	s.AddAnnotation(a)
	s.BeginChange()
	for i := range 5 {
		a.Translate(geom.Point{X: float64(i)})
		s.ReplaceAnnotation(a)
	}
	s.Undo()
	if got, _ := s.Annotation(a.ID); got.Points[0].X != 0 {
		t.Fatalf("after undo x = %v", got.Points[0].X)
	}
}

func TestParseRanges(t *testing.T) {
	got, err := ParseRanges("1-3, 5, 7-", 8)
	if err != nil {
		t.Fatal(err)
	}
	want := []int{0, 1, 2, 4, 6, 7}
	if len(got) != len(want) {
		t.Fatalf("got %v", got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
	if _, err := ParseRanges("4-2", 8); err == nil {
		t.Fatal("backwards range accepted")
	}
}

func TestSplitAndExport(t *testing.T) {
	s := newTestSession(t, 5)
	parts, err := Split(s.Snapshot(), 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(parts) != 3 {
		t.Fatalf("parts = %d", len(parts))
	}
	out, err := Export(s.Snapshot(), Security{UserPassword: "pw"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := pdf.Open(out, "pw"); err != nil {
		t.Fatal(err)
	}
}

func TestVersionsUniqueAcrossSessions(t *testing.T) {
	a, b := newTestSession(t, 1), newTestSession(t, 1)
	if a.Version() == b.Version() {
		t.Fatal("two sessions share a document version")
	}
	v := a.Version()
	a.Commit(a.Snapshot())
	if a.Version() != v {
		t.Fatal("markup-only change renewed the version")
	}
}
