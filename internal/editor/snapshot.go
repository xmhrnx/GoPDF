// Package editor is the editing session behind the UI: the document, the
// markup drawn over it and the undo history. Page operations are pure
// functions from one Snapshot to the next, so the UI can run them on a
// background goroutine and commit the result on the UI goroutine.
package editor

import (
	"gopdf/internal/annot"
	"gopdf/internal/pdf"
)

// Snapshot is one immutable state of the document and its markup.
type Snapshot struct {
	Doc    *pdf.Document
	Annots []annot.Annotation
}

// withAnnots returns s with a copy of its annotation list, safe to modify.
func (s Snapshot) withAnnots() Snapshot {
	list := make([]annot.Annotation, len(s.Annots))
	for i, a := range s.Annots {
		list[i] = a.Clone()
	}
	s.Annots = list
	return s
}
