package pdf

import (
	"testing"

	"gopdf/internal/annot"
	"gopdf/internal/digitalid"
	"gopdf/internal/geom"
)

func testSigner(t *testing.T) *digitalid.Signer {
	t.Helper()
	store := digitalid.NewStoreAt(t.TempDir())
	info, err := store.Create(digitalid.Identity{Name: "Jane Doe", Email: "jane@example.com", Organization: "Acme"}, "secret-pass")
	if err != nil {
		t.Fatal(err)
	}
	signer, err := store.Unlock(info.Fingerprint, "secret-pass")
	if err != nil {
		t.Fatal(err)
	}
	return signer
}

func sealMark(page int) annot.Annotation {
	st := annot.DefaultStyleFor(annot.Signature)
	a := annot.NewSignature(1, page, geom.Point{X: 200, Y: 300}, [][]geom.Point{{{X: 0, Y: 0.5}, {X: 1, Y: 0.5}}}, nil, 3, st)
	a.Seal = &annot.Seal{Reason: "I approve", Location: "Berlin"}
	return a
}

func TestSealAndRead(t *testing.T) {
	signer := testSigner(t)
	data, err := Rotate(samplePDF(2), []int{1}, 90)
	if err != nil {
		t.Fatal(err)
	}
	doc := mustLoad(t, data)
	mark := sealMark(1)
	signed, err := Seal(data, 1, doc.Pages[1], mark, signer, "")
	if err != nil {
		t.Fatal(err)
	}
	fields, err := ReadSignatures(signed, doc.Pages, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(fields) != 1 {
		t.Fatalf("fields = %+v", fields)
	}
	f := fields[0]
	if !f.Signed || !f.Intact || !f.Whole || f.Trusted || !f.SelfSig {
		t.Fatalf("field = %+v", f)
	}
	if f.Name != "Jane Doe" || f.Email != "jane@example.com" || f.Reason != "I approve" || f.Place != "Berlin" || f.Page != 1 {
		t.Fatalf("details = %+v", f)
	}
	want := mark.Box()
	if d := f.Rect.Min.Dist(want.Min) + f.Rect.Max.Dist(want.Max); d > 1 {
		t.Fatalf("rect = %v, want %v", f.Rect, want)
	}

	// A second signature appended later leaves the first intact.
	second, err := Seal(signed, 0, doc.Pages[0], sealMark(0), signer, "")
	if err != nil {
		t.Fatal(err)
	}
	fields, err = ReadSignatures(second, doc.Pages, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(fields) != 2 {
		t.Fatalf("after countersigning: %d fields", len(fields))
	}
	for _, f := range fields {
		if !f.Intact {
			t.Fatalf("signature broken by countersigning: %+v", f)
		}
	}
}

func TestSealDetectsTampering(t *testing.T) {
	signer := testSigner(t)
	data := samplePDF(1)
	doc := mustLoad(t, data)
	signed, err := Seal(data, 0, doc.Pages[0], sealMark(0), signer, "")
	if err != nil {
		t.Fatal(err)
	}
	// Change "Page 1" in the signed content to "Page 2".
	tampered := append([]byte(nil), signed...)
	i := indexOf(tampered, []byte("(Page 1)"))
	if i < 0 {
		t.Fatal("content not found")
	}
	tampered[i+6] = '2'
	fields, err := ReadSignatures(tampered, doc.Pages, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(fields) != 1 || fields[0].Intact {
		t.Fatalf("tampered file reads as intact: %+v", fields)
	}
}

func indexOf(b, sub []byte) int {
	for i := 0; i+len(sub) <= len(b); i++ {
		if string(b[i:i+len(sub)]) == string(sub) {
			return i
		}
	}
	return -1
}

func TestSignedFieldIsDisplayed(t *testing.T) {
	signer := testSigner(t)
	data := samplePDF(1)
	doc := mustLoad(t, data)
	mark := sealMark(0) // a black stroke across the middle of its box
	signed, err := Seal(data, 0, doc.Pages[0], mark, signer, "")
	if err != nil {
		t.Fatal(err)
	}
	loaded := mustLoad(t, signed)
	if len(loaded.Signatures) != 1 || !loaded.Signatures[0].Intact {
		t.Fatalf("signatures = %+v", loaded.Signatures)
	}
	r, err := NewRenderer(loaded.Display)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	img, err := r.Render(0, 72)
	if err != nil {
		t.Fatal(err)
	}
	b := mark.Box()
	cx, cy := int((b.Min.X+b.Max.X)/2), int((b.Min.Y+b.Max.Y)/2)
	if c := img.RGBAAt(cx, cy); c.R > 120 {
		t.Fatalf("pixel on the signature stroke = %v, want ink", c)
	}
}
