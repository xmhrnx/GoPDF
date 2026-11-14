package editor

import (
	"fmt"

	"gopdf/internal/annot"
	"gopdf/internal/digitalid"
	"gopdf/internal/pdf"
)

// Signers holds unlocked digital IDs by fingerprint.
type Signers map[string]*digitalid.Signer

// LockedError names a digital ID that must be unlocked before saving.
type LockedError struct{ ID string }

func (e LockedError) Error() string { return fmt.Sprintf("digital ID %s is locked", e.ID) }

// Security is the protection applied when the document is written.
type Security struct {
	UserPassword  string
	OwnerPassword string
}

// Enabled reports whether a password is set.
func (sec Security) Enabled() bool { return sec.UserPassword != "" }

// Export returns the bytes to write for s: the markup flattened into the
// pages, the file encrypted when sec is enabled, then each sealed
// signature applied as a digital signature with its ID from signers.
// With no markup and no password the original bytes are kept, so new
// signatures add to the document's existing ones without breaking them.
func Export(s Snapshot, sec Security, signers Signers) ([]byte, error) {
	plain, sealed := splitSealed(s.Annots)
	for _, a := range sealed {
		if signers[a.Seal.ID] == nil {
			return nil, LockedError{ID: a.Seal.ID}
		}
	}
	data, err := pdf.Flatten(s.Doc.Data, plain)
	if err != nil {
		return nil, err
	}
	if sec.Enabled() {
		if data, err = pdf.Encrypt(data, sec.UserPassword, sec.OwnerPassword); err != nil {
			return nil, err
		}
	}
	for _, a := range sealed {
		if data, err = pdf.Seal(data, a.Page, s.Doc.Pages[a.Page], a, signers[a.Seal.ID], sec.UserPassword); err != nil {
			return nil, err
		}
	}
	return data, nil
}

// splitSealed separates signatures to seal from the rest of the markup.
func splitSealed(list []annot.Annotation) (plain, sealed []annot.Annotation) {
	for _, a := range list {
		if a.Kind == annot.Signature && a.Seal != nil {
			sealed = append(sealed, a)
		} else {
			plain = append(plain, a)
		}
	}
	return plain, sealed
}

// SealIDs lists the digital IDs saving s needs, each once.
func SealIDs(s Snapshot) []string {
	var ids []string
	seen := map[string]bool{}
	_, sealed := splitSealed(s.Annots)
	for _, a := range sealed {
		if !seen[a.Seal.ID] {
			seen[a.Seal.ID] = true
			ids = append(ids, a.Seal.ID)
		}
	}
	return ids
}

// ExportPages returns a new document holding only pages (old indexes,
// in order), markup flattened.
func ExportPages(s Snapshot, pages []int) ([]byte, error) {
	sub, err := Arrange(s, pages)
	if err != nil {
		return nil, err
	}
	sub.Annots = unseal(sub.Annots)
	return Export(sub, Security{}, nil)
}

// Split cuts the document into parts of at most every pages each.
func Split(s Snapshot, every int) ([][]byte, error) {
	var parts [][]byte
	n := s.Doc.PageCount()
	for start := 0; start < n; start += every {
		var pages []int
		for p := start; p < min(start+every, n); p++ {
			pages = append(pages, p)
		}
		part, err := ExportPages(s, pages)
		if err != nil {
			return nil, err
		}
		parts = append(parts, part)
	}
	return parts, nil
}

// unseal turns sealed signatures into plain marks, for copies of pages
// that don't carry digital signatures.
func unseal(list []annot.Annotation) []annot.Annotation {
	out := make([]annot.Annotation, len(list))
	for i, a := range list {
		out[i] = a.Clone()
		out[i].Seal = nil
	}
	return out
}
