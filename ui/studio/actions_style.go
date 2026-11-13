package studio

import "gopdf/internal/annot"

// The style setters change the active tool's style and, when markup is
// selected, the selection too (as one undoable change).

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
	s.editStyle(func(st *annot.Style) { st.FontSize = size })
	s.updateSelected(func(a *annot.Annotation) bool {
		a.FontSize = size
		return a.Kind == annot.Text
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

// editStyle changes the style of the active tool.
func (s *Studio) editStyle(change func(*annot.Style)) {
	st := s.toolStyle(s.tool)
	change(&st)
	s.styles[s.tool] = st
}

// toolStyle returns the style tool draws with.
func (s *Studio) toolStyle(t Tool) annot.Style {
	if st, ok := s.styles[t]; ok {
		return st
	}
	return annot.DefaultStyle
}

// defaultToolStyles gives the highlighters a highlighter color.
func defaultToolStyles() map[Tool]annot.Style {
	yellow := annot.DefaultStyle
	yellow.Color = annot.Swatches[3]
	return map[Tool]annot.Style{
		Tool(annot.Marker):    yellow,
		Tool(annot.Highlight): yellow,
	}
}
