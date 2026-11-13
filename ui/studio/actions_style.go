package studio

import "gopdf/internal/annot"

// The style setters change the style remembered for the kind being
// styled (the selection's, or the active tool's) and, when markup is
// selected, the selection too, as one undoable change.

func (s *Studio) setColor(c annot.Color) {
	s.editStyle(func(st *annot.Style) { st.Color = c })
	s.updateSelected(func(a *annot.Annotation) bool {
		if a.Kind == annot.Whiteout {
			return false
		}
		a.Color = c
		return true
	})
}

func (s *Studio) setStrokeWidth(w float64) {
	s.editStyle(func(st *annot.Style) { st.Width = w })
	s.updateSelected(func(a *annot.Annotation) bool {
		switch {
		case a.Kind == annot.Marker:
			a.Width = max(w*5, 10)
		case a.Kind == annot.Text, a.Kind == annot.Highlight, a.Kind == annot.Whiteout:
			return false
		default:
			a.Width = w
		}
		return true
	})
}

func (s *Studio) setFontSize(size float64) {
	size = min(max(size, annot.MinFontSize), annot.MaxFontSize)
	s.editStyle(func(st *annot.Style) { st.FontSize = size })
	s.updateSelected(func(a *annot.Annotation) bool {
		a.FontSize = size
		return a.Kind == annot.Text
	})
}

func (s *Studio) setSignatureSize(w float64) {
	w = min(max(w, annot.MinSignatureSize), annot.MaxSignatureSize)
	s.editStyle(func(st *annot.Style) { st.Size = w })
	s.updateSelected(func(a *annot.Annotation) bool {
		a.Resize(w)
		return a.Kind == annot.Signature
	})
}

func (s *Studio) setFilled(on bool) {
	s.editStyle(func(st *annot.Style) { st.Filled = on })
	s.updateSelected(func(a *annot.Annotation) bool {
		a.Filled = on
		return a.Kind == annot.Rect || a.Kind == annot.Ellipse
	})
}

// updateSelected applies change to a copy of the selection and commits
// it when change reports that it applies.
func (s *Studio) updateSelected(change func(*annot.Annotation) bool) {
	a, ok := s.selectedAnnot()
	if !ok || !s.idle() {
		return
	}
	a = a.Clone()
	if change(&a) {
		s.doc.UpdateAnnotation(a)
	}
}

// styledKind is the kind the style controls act on: the selection's,
// or the drawing tool's.
func (s *Studio) styledKind() (annot.Kind, bool) {
	if a, ok := s.selectedAnnot(); ok {
		return a.Kind, true
	}
	if s.tool != ToolSelect {
		return s.tool.Kind(), true
	}
	return 0, false
}

// editStyle changes and remembers the style of the kind being styled.
func (s *Studio) editStyle(change func(*annot.Style)) {
	k, ok := s.styledKind()
	if !ok {
		return
	}
	st := s.memory.style(k)
	change(&st)
	s.memory.setStyle(s, k, st)
}
