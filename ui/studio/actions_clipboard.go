package studio

import (
	"gopdf/internal/annot"
	"gopdf/internal/editor"
	"gopdf/internal/geom"
)

// copySelection puts the selected markup on the clipboard.
func (s *Studio) copySelection() {
	a, ok := s.selectedAnnot()
	if !ok {
		return
	}
	c := a.Clone()
	s.clipboard.item, s.clipboard.pastes = &c, 0
}

// cutSelection copies the selection, then deletes it.
func (s *Studio) cutSelection() {
	if _, ok := s.selectedAnnot(); ok && s.canEdit() {
		s.copySelection()
		s.deleteSelected()
	}
}

// paste adds a copy of the clipboard to the current page, offset from
// the original, and selects it.
func (s *Studio) paste() {
	src := s.clipboard.item
	if src == nil || !s.canEdit() {
		return
	}
	s.clipboard.pastes++
	a := src.Clone()
	a.ID, a.Page = editor.NewID(), s.view.current
	d := pasteOffset * float64(s.clipboard.pastes)
	a.Translate(geom.Point{X: d, Y: d})
	s.keepOnPage(&a)
	s.doc.AddAnnotation(a)
	s.tool = ToolSelect
	s.selectAnnotation(a.ID)
}

// duplicateSelection pastes a copy of the selection next to it, leaving
// the clipboard as it was.
func (s *Studio) duplicateSelection() {
	a, ok := s.selectedAnnot()
	if !ok || !s.canEdit() {
		return
	}
	c := a.Clone()
	c.ID = editor.NewID()
	c.Translate(geom.Point{X: pasteOffset, Y: pasteOffset})
	s.keepOnPage(&c)
	s.doc.AddAnnotation(c)
	s.selectAnnotation(c.ID)
}

// keepOnPage moves a back inside its page when it sticks out past the
// right or bottom edge (a paste onto a smaller page).
func (s *Studio) keepOnPage(a *annot.Annotation) {
	w, h := s.doc.Doc().Pages[a.Page].Size()
	b := a.Bounds()
	d := geom.Point{X: min(0, w-b.Max.X), Y: min(0, h-b.Max.Y)}
	d.X = max(d.X, -b.Min.X)
	d.Y = max(d.Y, -b.Min.Y)
	a.Translate(d)
}

// duplicate copies the selected markup, or the current page when
// nothing is selected.
func (s *Studio) duplicate() {
	if _, ok := s.selectedAnnot(); ok {
		s.duplicateSelection()
		return
	}
	if s.hasDoc() {
		s.duplicatePage(s.view.current)
	}
}
