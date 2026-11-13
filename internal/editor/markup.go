package editor

import "gopdf/internal/annot"

// AddAnnotation appends a to the markup as one undoable change.
func (s *Session) AddAnnotation(a annot.Annotation) {
	next := s.cur.withAnnots()
	next.Annots = append(next.Annots, a.Clone())
	s.Commit(next)
}

// UpdateAnnotation replaces the annotation with a's ID as one undoable
// change. It does nothing when no annotation has that ID.
func (s *Session) UpdateAnnotation(a annot.Annotation) {
	i := s.indexOf(a.ID)
	if i < 0 {
		return
	}
	next := s.cur.withAnnots()
	next.Annots[i] = a.Clone()
	s.Commit(next)
}

// ReplaceAnnotation swaps in a without recording history: for live
// previews during a drag that began with BeginChange.
func (s *Session) ReplaceAnnotation(a annot.Annotation) {
	if i := s.indexOf(a.ID); i >= 0 {
		s.cur.Annots[i] = a.Clone()
		s.dirty = true
	}
}

// BeginChange records the current state for undo before a series of
// ReplaceAnnotation calls (a drag), and detaches the list so the recorded
// state isn't modified by them.
func (s *Session) BeginChange() {
	s.Commit(s.cur.withAnnots())
}

// DeleteAnnotation removes the annotation with id as one undoable change.
func (s *Session) DeleteAnnotation(id int64) {
	i := s.indexOf(id)
	if i < 0 {
		return
	}
	next := s.cur.withAnnots()
	next.Annots = append(next.Annots[:i], next.Annots[i+1:]...)
	s.Commit(next)
}

// Annotation returns the annotation with id.
func (s *Session) Annotation(id int64) (annot.Annotation, bool) {
	if i := s.indexOf(id); i >= 0 {
		return s.cur.Annots[i], true
	}
	return annot.Annotation{}, false
}

func (s *Session) indexOf(id int64) int {
	for i, a := range s.cur.Annots {
		if a.ID == id {
			return i
		}
	}
	return -1
}
