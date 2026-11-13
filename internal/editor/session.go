package editor

import (
	"path/filepath"

	"gopdf/internal/annot"
	"gopdf/internal/history"
	"gopdf/internal/pdf"
)

// historyLimit is how many changes can be undone.
const historyLimit = 100

// Session is one open document. It is not safe for concurrent use: the
// UI goroutine owns it, and background work goes through Snapshot.
type Session struct {
	Path     string // where the file was opened from or last saved to; may be empty
	Name     string // display name
	Security Security

	cur     Snapshot
	hist    *history.Stack[Snapshot]
	dirty   bool
	version int // renewed on every change of the document bytes
}

// NewSession wraps an opened document.
func NewSession(doc *pdf.Document, path, name string) *Session {
	if name == "" {
		name = filepath.Base(path)
	}
	return &Session{Path: path, Name: name, cur: Snapshot{Doc: doc}, hist: history.New[Snapshot](historyLimit), version: nextVersion()}
}

// Snapshot returns the current state.
func (s *Session) Snapshot() Snapshot { return s.cur }

// Doc returns the current document.
func (s *Session) Doc() *pdf.Document { return s.cur.Doc }

// Annots returns the current markup. Callers must not modify it.
func (s *Session) Annots() []annot.Annotation { return s.cur.Annots }

// Version identifies the document bytes, unique across sessions; it
// changes whenever pages do, so rendered images keyed by it go stale
// exactly when they should.
func (s *Session) Version() int { return s.version }

// Dirty reports unsaved changes.
func (s *Session) Dirty() bool { return s.dirty }

// Commit makes next the current state, recording the old one for undo.
func (s *Session) Commit(next Snapshot) {
	s.hist.Push(s.cur)
	s.set(next)
}

func (s *Session) set(next Snapshot) {
	if next.Doc != s.cur.Doc {
		s.version = nextVersion()
	}
	s.cur = next
	s.dirty = true
}

// Undo restores the state before the last change.
func (s *Session) Undo() bool {
	prev, ok := s.hist.Undo(s.cur)
	if ok {
		s.set(prev)
	}
	return ok
}

// Redo reapplies the last undone change.
func (s *Session) Redo() bool {
	next, ok := s.hist.Redo(s.cur)
	if ok {
		s.set(next)
	}
	return ok
}

// CanUndo reports whether Undo would do something.
func (s *Session) CanUndo() bool { return s.hist.CanUndo() }

// CanRedo reports whether Redo would do something.
func (s *Session) CanRedo() bool { return s.hist.CanRedo() }

// MarkSaved records that the current state is on disk at path.
func (s *Session) MarkSaved(path string) {
	if path != "" {
		s.Path, s.Name = path, filepath.Base(path)
	}
	s.dirty = false
}

// Touch marks the document changed without an undo step, for settings
// that only affect how it is written (its password).
func (s *Session) Touch() { s.dirty = true }
