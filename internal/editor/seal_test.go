package editor

import (
	"errors"
	"testing"

	"gopdf/internal/annot"
	"gopdf/internal/digitalid"
	"gopdf/internal/geom"
	"gopdf/internal/pdf"
)

func sealed(id string) annot.Annotation {
	a := annot.NewSignature(NewID(), 0, geom.Point{X: 200, Y: 200},
		[][]geom.Point{{{X: 0, Y: 0}, {X: 1, Y: 1}}}, nil, 3, annot.DefaultStyleFor(annot.Signature))
	a.Seal = &annot.Seal{ID: id, Reason: "Approved"}
	return a
}

func TestExportSealsAndKeepsEarlierSignatures(t *testing.T) {
	store := digitalid.NewStoreAt(t.TempDir())
	info, err := store.Create(digitalid.Identity{Name: "Jane Doe"}, "secret-pass")
	if err != nil {
		t.Fatal(err)
	}
	signer, err := store.Unlock(info.Fingerprint, "secret-pass")
	if err != nil {
		t.Fatal(err)
	}
	s := newTestSession(t, 1)
	s.AddAnnotation(sealed(info.Fingerprint))
	if _, err := Export(s.Snapshot(), Security{}, nil); !errors.As(err, &LockedError{}) {
		t.Fatalf("export without the ID unlocked: %v", err)
	}
	if ids := SealIDs(s.Snapshot()); len(ids) != 1 || ids[0] != info.Fingerprint {
		t.Fatalf("SealIDs = %v", ids)
	}
	signers := Signers{info.Fingerprint: signer}
	out, err := Export(s.Snapshot(), Security{}, signers)
	if err != nil {
		t.Fatal(err)
	}

	// Reopened, adding a second signature keeps the first; a plain edit
	// would not.
	doc, err := pdf.Load(out)
	if err != nil {
		t.Fatal(err)
	}
	again := NewSession(doc, "", "signed.pdf")
	if n := again.SignaturesAtRisk(); n != 0 {
		t.Fatalf("at risk before changes = %d", n)
	}
	again.AddAnnotation(sealed(info.Fingerprint))
	if n := again.SignaturesAtRisk(); n != 0 {
		t.Fatalf("at risk with only a new signature = %d", n)
	}
	twice, err := Export(again.Snapshot(), Security{}, signers)
	if err != nil {
		t.Fatal(err)
	}
	final, err := pdf.Load(twice)
	if err != nil {
		t.Fatal(err)
	}
	if len(final.Signatures) != 2 || !final.Signatures[0].Intact || !final.Signatures[1].Intact {
		t.Fatalf("signatures = %+v", final.Signatures)
	}
	again.AddAnnotation(annot.New(NewID(), 0, annot.Rect, geom.Point{}, annot.DefaultStyle))
	if n := again.SignaturesAtRisk(); n != 1 {
		t.Fatalf("at risk after a plain edit = %d", n)
	}
}

func TestSignedAndEncrypted(t *testing.T) {
	store := digitalid.NewStoreAt(t.TempDir())
	info, _ := store.Create(digitalid.Identity{Name: "Jane Doe"}, "secret-pass")
	signer, _ := store.Unlock(info.Fingerprint, "secret-pass")
	s := newTestSession(t, 1)
	s.AddAnnotation(sealed(info.Fingerprint))
	out, err := Export(s.Snapshot(), Security{UserPassword: "open-me"}, Signers{info.Fingerprint: signer})
	if err != nil {
		t.Fatal(err)
	}
	doc, err := pdf.Open(out, "open-me")
	if err != nil {
		t.Fatal(err)
	}
	fields, err := pdf.ReadSignatures(out, doc.Pages, "open-me")
	if err != nil {
		t.Fatal(err)
	}
	if len(fields) != 1 || !fields[0].Intact {
		t.Fatalf("signatures of the encrypted file = %+v", fields)
	}
}
