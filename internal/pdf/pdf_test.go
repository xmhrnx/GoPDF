package pdf

import (
	"errors"
	"strings"
	"testing"

	"gopdf/internal/annot"
	"gopdf/internal/geom"
)

func mustLoad(t *testing.T, data []byte) *Document {
	t.Helper()
	d, err := Load(data)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func TestArrangeRotateInsert(t *testing.T) {
	data := samplePDF(3)
	out, err := Arrange(data, []int{2, 0, 0})
	if err != nil {
		t.Fatal(err)
	}
	d := mustLoad(t, out)
	if d.PageCount() != 3 {
		t.Fatalf("pages = %d", d.PageCount())
	}
	r, err := NewRenderer(out)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	if txt, _ := r.Text(0); !strings.Contains(txt, "Page 3") {
		t.Fatalf("first page text %q", txt)
	}
	rot, err := Rotate(out, []int{1}, 90)
	if err != nil {
		t.Fatal(err)
	}
	if d := mustLoad(t, rot); d.Pages[1].Rotate != 90 {
		t.Fatalf("rotation = %d", d.Pages[1].Rotate)
	}
	ins, err := InsertBlank(rot, -1, 300, 400)
	if err != nil {
		t.Fatal(err)
	}
	d = mustLoad(t, ins)
	if w, h := d.Pages[0].Size(); d.PageCount() != 4 || w != 300 || h != 400 {
		t.Fatalf("after insert: %d pages, first %vx%v", d.PageCount(), w, h)
	}
}

func TestMergeWatermarkNumbers(t *testing.T) {
	out, err := Merge(samplePDF(2), samplePDF(1))
	if err != nil {
		t.Fatal(err)
	}
	if out, err = AddWatermark(out, Watermark{Text: "DRAFT", FontSize: 48, Rotation: 45, Opacity: 0.3, Color: "#ff0000"}); err != nil {
		t.Fatal(err)
	}
	if out, err = AddPageNumbers(out, "%p / %P", BottomCenter, 10); err != nil {
		t.Fatal(err)
	}
	if d := mustLoad(t, out); d.PageCount() != 3 {
		t.Fatalf("pages = %d", d.PageCount())
	}
}

func TestFlattenRendersMarkup(t *testing.T) {
	data, err := Rotate(samplePDF(1), []int{0}, 90)
	if err != nil {
		t.Fatal(err)
	}
	box := annot.New(1, 0, annot.Rect, geom.Point{X: 20, Y: 20}, annot.Style{Color: annot.Color{R: 255}, Width: 4, Filled: true})
	box.Extend(geom.Point{X: 120, Y: 120})
	text := annot.New(2, 0, annot.Text, geom.Point{X: 200, Y: 200}, annot.DefaultStyle)
	text.Text = "Hello (world)\nsecond line é"
	hl := annot.New(3, 0, annot.Highlight, geom.Point{X: 300, Y: 300}, annot.Style{Color: annot.Swatches[3]})
	hl.Extend(geom.Point{X: 400, Y: 330})
	out, err := Flatten(data, []annot.Annotation{box, text, hl})
	if err != nil {
		t.Fatal(err)
	}
	r, err := NewRenderer(out)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	img, err := r.Render(0, 72)
	if err != nil {
		t.Fatal(err)
	}
	// The display-space point (70, 70) is inside the filled red box.
	c := img.RGBAAt(70, 70)
	if c.R < 200 || c.G > 60 || c.B > 60 {
		t.Fatalf("pixel at box center = %v, want red", c)
	}
	if txt, _ := r.Text(0); !strings.Contains(txt, "Hello (world)") {
		t.Fatalf("flattened text missing: %q", txt)
	}
}

func TestEncryptOpen(t *testing.T) {
	out, err := Encrypt(samplePDF(1), "secret", "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Open(out, ""); !errors.Is(err, ErrPassword) {
		t.Fatalf("open without password: %v", err)
	}
	if _, err := Open(out, "wrong"); !errors.Is(err, ErrPassword) {
		t.Fatalf("open with wrong password: %v", err)
	}
	d, err := Open(out, "secret")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewRenderer(d.Data); err != nil {
		t.Fatalf("decrypted file does not render: %v", err)
	}
}

func TestMetadata(t *testing.T) {
	out, err := SetMetadata(samplePDF(1), Metadata{Title: "Quarterly", Author: "Ada"})
	if err != nil {
		t.Fatal(err)
	}
	if m := mustLoad(t, out).Meta; m.Title != "Quarterly" || m.Author != "Ada" {
		t.Fatalf("metadata = %+v", m)
	}
}

func TestFindInText(t *testing.T) {
	m := FindInText(4, "Alpha beta\ngamma BETA delta", "beta", 10)
	if len(m) != 2 || m[0].Page != 4 || !strings.Contains(m[1].Snippet, "gamma BETA") {
		t.Fatalf("matches = %+v", m)
	}
}

// A damaged file makes MuPDF throw, which on Windows exercises the
// setjmp/longjmp pair behind its fz_try/fz_catch.
func TestRendererRejectsDamagedFile(t *testing.T) {
	if r, err := NewRenderer([]byte("%PDF-1.7\n1 0 obj << /Type /Catalog >>\ngarbage")); err == nil {
		r.Close()
		t.Fatal("damaged file opened without error")
	}
}

func TestFlattenSignature(t *testing.T) {
	strokes := [][]geom.Point{{{X: 0, Y: 0.5}, {X: 1, Y: 0.5}}}
	st := annot.DefaultStyleFor(annot.Signature)
	st.Width = 6
	sig := annot.NewSignature(1, 0, geom.Point{X: 300, Y: 400}, strokes, 4, st)
	out, err := Flatten(samplePDF(1), []annot.Annotation{sig})
	if err != nil {
		t.Fatal(err)
	}
	r, err := NewRenderer(out)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	img, err := r.Render(0, 72)
	if err != nil {
		t.Fatal(err)
	}
	// The stroke crosses the box's middle: (300, 400) is inked.
	if c := img.RGBAAt(300, 400); c.R > 120 {
		t.Fatalf("pixel on the signature = %v, want dark ink", c)
	}
}
