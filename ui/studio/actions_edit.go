package studio

import "gopdf/internal/annot"

func (s *Studio) undo() {
	if s.canEdit() && s.doc.Undo() {
		s.afterHistoryJump()
	}
}

func (s *Studio) redo() {
	if s.canEdit() && s.doc.Redo() {
		s.afterHistoryJump()
	}
}

func (s *Studio) afterHistoryJump() {
	if _, ok := s.doc.Annotation(s.selected); !ok {
		s.selected = 0
	}
	s.ui.inspector.syncText()
}

// deleteSelected removes the selected annotation.
func (s *Studio) deleteSelected() {
	if s.canEdit() && s.selected != 0 {
		s.doc.DeleteAnnotation(s.selected)
		s.selected = 0
	}
}

// setTool switches tools; Escape-like behavior when choosing Select twice.
func (s *Studio) setTool(t Tool) {
	s.tool = t
	if t != ToolSelect {
		s.selected = 0
	}
}

// selectedAnnot returns the selected annotation.
func (s *Studio) selectedAnnot() (annot.Annotation, bool) {
	if s.doc == nil || s.selected == 0 {
		return annot.Annotation{}, false
	}
	return s.doc.Annotation(s.selected)
}

// selectAnnotation selects id and keeps the inspector's text in step.
func (s *Studio) selectAnnotation(id int64) {
	if id != s.selected {
		s.ui.viewer.blurFields()
	}
	s.selected = id
	s.ui.inspector.syncText()
}

// escape cancels the current gesture, then the selection, then the tool.
func (s *Studio) escape() {
	switch {
	case s.drawing != nil:
		s.drawing = nil
	case s.selected != 0:
		s.selected = 0
	default:
		s.tool = ToolSelect
	}
}

// escapeKey leaves a text field first, then works as escape.
func (s *Studio) escapeKey() {
	if s.ui.typing() && !s.ui.dialogOpen() {
		s.ui.viewer.blurFields()
		return
	}
	s.escape()
}
