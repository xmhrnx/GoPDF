package editor

import "gopdf/internal/pdf"

// Security is the protection applied when the document is written.
type Security struct {
	UserPassword  string
	OwnerPassword string
}

// Enabled reports whether a password is set.
func (sec Security) Enabled() bool { return sec.UserPassword != "" }

// Export returns the bytes to write for s: the markup flattened into the
// pages and, when sec is enabled, the file encrypted.
func Export(s Snapshot, sec Security) ([]byte, error) {
	data, err := pdf.Flatten(s.Doc.Data, s.Annots)
	if err != nil {
		return nil, err
	}
	if sec.Enabled() {
		return pdf.Encrypt(data, sec.UserPassword, sec.OwnerPassword)
	}
	return data, nil
}

// ExportPages returns a new document holding only pages (old indexes,
// in order), markup flattened.
func ExportPages(s Snapshot, pages []int) ([]byte, error) {
	sub, err := Arrange(s, pages)
	if err != nil {
		return nil, err
	}
	return Export(sub, Security{})
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
