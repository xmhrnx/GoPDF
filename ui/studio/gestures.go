package studio

import (
	"github.com/xmhrnx/giode"

	"gopdf/internal/annot"
	"gopdf/internal/editor"
	"gopdf/internal/geom"
)

// moveGesture is markup being dragged with the select tool.
type moveGesture struct {
	page   int
	start  geom.Point
	origin annot.Annotation
}

// cursor is the pointer shape over pages for the active tool.
func (s *Studio) cursor() giode.Cursor {
	switch {
	case s.tool == ToolSelect:
		return giode.CursorDefault
	case s.tool.Kind() == annot.Text:
		return giode.CursorText
	case s.tool.Kind() == annot.Signature:
		return giode.CursorPointer
	default:
		return giode.CursorCrosshair
	}
}

// clickPage handles a click (no drag) at p on page.
func (s *Studio) clickPage(page int, p geom.Point) {
	s.view.current = page
	s.ui.viewer.blurFields()
	if !s.canEdit() {
		return
	}
	switch {
	case s.tool == ToolSelect:
		s.selectAnnotation(s.hitTest(page, p))
	case s.tool.Kind() == annot.Text:
		s.placeText(page, p)
	case s.tool.Kind() == annot.Signature:
		s.placeSignature(page, p)
	}
}

func (s *Studio) hitTest(page int, p geom.Point) int64 {
	list := s.doc.Annots()
	if i := annot.TopmostAt(list, page, p); i >= 0 {
		return list[i].ID
	}
	return 0
}

// placeText adds a text box at p and opens it for typing.
func (s *Studio) placeText(page int, p geom.Point) {
	a := annot.New(editor.NewID(), page, annot.Text, p, s.memory.style(annot.Text))
	s.doc.AddAnnotation(a)
	s.tool = ToolSelect
	s.selectAnnotation(a.ID)
	s.view.inspector = true
	s.ui.inspector.editText()
}

// placeSignature puts the active signature centered at p, or asks for
// one when none is saved yet.
func (s *Studio) placeSignature(page int, p geom.Point) {
	sig, ok := s.signatures.active(s)
	if !ok {
		s.ui.signaturePad.show()
		return
	}
	a := annot.NewSignature(editor.NewID(), page, p, sig.Strokes, sig.Aspect, s.memory.style(annot.Signature))
	s.keepOnPage(&a)
	s.doc.AddAnnotation(a)
	s.tool = ToolSelect
	s.selectAnnotation(a.ID)
}

// startGesture begins drawing or moving at p on page.
func (s *Studio) startGesture(page int, p geom.Point) {
	s.view.current = page
	s.ui.viewer.blurFields()
	if !s.canEdit() {
		return
	}
	switch {
	case s.tool == ToolSelect:
		id := s.hitTest(page, p)
		if id == 0 {
			s.selectAnnotation(0)
			return
		}
		a, _ := s.doc.Annotation(id)
		s.selectAnnotation(id)
		s.doc.BeginChange()
		s.moving = &moveGesture{page: page, start: p, origin: a.Clone()}
	case s.tool.Kind() == annot.Text:
		s.placeText(page, p)
	case s.tool.Kind() == annot.Signature:
		s.placeSignature(page, p)
	default:
		a := annot.New(editor.NewID(), page, s.tool.Kind(), p, s.memory.style(s.tool.Kind()))
		s.drawing = &a
	}
}

// moveGesture extends the shape being drawn or moves the markup held.
func (s *Studio) moveGesture(page int, p geom.Point) {
	p = s.clampToPage(page, p)
	switch {
	case s.drawing != nil && s.drawing.Page == page:
		s.drawing.Extend(p)
	case s.moving != nil && s.moving.page == page:
		a := s.moving.origin.Clone()
		a.Translate(p.Sub(s.moving.start))
		s.doc.ReplaceAnnotation(a)
	}
}

// endGesture finishes the gesture; a cancelled drag leaves no trace.
func (s *Studio) endGesture(cancelled bool) {
	if d := s.drawing; d != nil {
		s.drawing = nil
		if !cancelled && !d.Degenerate() {
			s.doc.AddAnnotation(*d)
		}
	}
	if m := s.moving; m != nil {
		s.moving = nil
		if cancelled {
			s.doc.ReplaceAnnotation(m.origin)
		}
	}
}

func (s *Studio) clampToPage(page int, p geom.Point) geom.Point {
	w, h := s.doc.Doc().Pages[page].Size()
	return geom.Point{X: min(max(p.X, 0), w), Y: min(max(p.Y, 0), h)}
}
