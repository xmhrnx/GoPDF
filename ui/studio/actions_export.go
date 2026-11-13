package studio

import (
	"bytes"
	"fmt"
	"image/png"
	"strings"

	"gopdf/internal/editor"
	"gopdf/internal/pdf"
)

// exportDPI is the resolution of exported page images.
const exportDPI = 200

// exportPageImage saves the current page, markup included, as a PNG.
func (s *Studio) exportPageImage() {
	if !s.hasDoc() {
		return
	}
	page, snap := s.view.current, s.doc.Snapshot()
	name := fmt.Sprintf("%s-page-%d.png", baseName(s.doc.Name), page+1)
	s.exportFile("Exporting image", name, func() ([]byte, error) {
		data, err := editor.ExportPages(snap, []int{page})
		if err != nil {
			return nil, err
		}
		r, err := pdf.NewRenderer(data)
		if err != nil {
			return nil, err
		}
		defer r.Close()
		img, err := r.Render(0, exportDPI)
		if err != nil {
			return nil, err
		}
		var buf bytes.Buffer
		err = png.Encode(&buf, img)
		return buf.Bytes(), err
	})
}

// exportText saves the text of every page as a plain text file.
func (s *Studio) exportText() {
	if !s.hasDoc() {
		return
	}
	data := s.doc.Doc().Data
	n := s.doc.Doc().PageCount()
	s.exportFile("Exporting text", baseName(s.doc.Name)+".txt", func() ([]byte, error) {
		r, err := pdf.NewRenderer(data)
		if err != nil {
			return nil, err
		}
		defer r.Close()
		var b strings.Builder
		for p := range n {
			t, err := r.Text(p)
			if err != nil {
				return nil, err
			}
			fmt.Fprintf(&b, "── Page %d ──\n\n%s\n\n", p+1, strings.TrimSpace(t))
		}
		return []byte(b.String()), nil
	})
}
