package studio

import (
	"github.com/xmhrnx/giode"

	"gopdf/internal/pdf"
)

// syncRenderer replaces the renderer when the document bytes changed,
// so rendered pages and thumbnails always show the current version.
func (s *Studio) syncRenderer() {
	if s.doc == nil {
		s.dropRenderer()
		return
	}
	if s.renderer != nil && s.rendered == s.doc.Version() {
		return
	}
	r, err := pdf.NewRenderer(s.doc.Doc().Data)
	if err != nil {
		s.toasts.Show(giode.ToastError, "Can't display this document: "+err.Error())
		return
	}
	old := s.renderer
	s.renderer, s.rendered = r, s.doc.Version()
	s.pages.SetRenderer(r, s.rendered)
	s.thumbs.SetRenderer(r, s.rendered)
	s.outline = r.Outline()
	s.ui.search.reset()
	if old != nil {
		old.Close()
	}
	s.clampCurrent()
}

func (s *Studio) dropRenderer() {
	if s.renderer == nil {
		return
	}
	s.pages.SetRenderer(nil, -1)
	s.thumbs.SetRenderer(nil, -1)
	s.renderer.Close()
	s.renderer, s.outline = nil, nil
}

func (s *Studio) clampCurrent() {
	if n := s.doc.Doc().PageCount(); s.view.current >= n {
		s.view.current = n - 1
	}
}
