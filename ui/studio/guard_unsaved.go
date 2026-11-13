package studio

// guardUnsaved runs then right away, or after the user has saved or
// discarded the open document's unsaved changes.
func (s *Studio) guardUnsaved(then func()) {
	if s.doc == nil || !s.doc.Dirty() {
		then()
		return
	}
	s.ui.unsaved.ask(s.doc.Name, then)
}
