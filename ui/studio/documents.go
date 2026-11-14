package studio

import (
	"gopdf/internal/editor"
	"gopdf/internal/pdf"
)

// openDoc is one open document and how it was being viewed, kept while
// another document is active.
type openDoc struct {
	session  *editor.Session
	renderer *pdf.Renderer
	rendered int
	outline  []pdf.OutlineItem
	selected int64
	zoom     float64
	fitWidth bool
	current  int
}

// addDocument opens session beside the others and shows it; a file
// already open is shown instead of opened twice.
func (s *Studio) addDocument(session *editor.Session) {
	for i, d := range s.docs {
		if session.Path != "" && d.session.Path == session.Path {
			s.switchTo(i)
			return
		}
	}
	s.stash()
	s.docs = append(s.docs, &openDoc{session: session, zoom: 1, fitWidth: true})
	s.active = -1
	s.switchTo(len(s.docs) - 1)
	s.recents.add(s, session.Path, session.Doc().PageCount())
}

// replaceDocument puts session in place of the active document (a file
// reopened after signing).
func (s *Studio) replaceDocument(session *editor.Session) {
	if s.active < 0 {
		s.addDocument(session)
		return
	}
	s.dropRenderer()
	s.docs[s.active] = &openDoc{session: session, zoom: s.view.zoom, fitWidth: s.view.fitWidth, current: s.view.current}
	s.active = -1
	s.switchTo(len(s.docs) - 1)
}

// switchTo makes document i the active one.
func (s *Studio) switchTo(i int) {
	if i < 0 || i >= len(s.docs) || i == s.active {
		return
	}
	s.stash()
	s.active = i
	d := s.docs[i]
	s.doc, s.renderer, s.rendered, s.outline = d.session, d.renderer, d.rendered, d.outline
	s.selected, s.drawing, s.moving = d.selected, nil, nil
	s.view.zoom, s.view.fitWidth, s.view.current = d.zoom, d.fitWidth, d.current
	if s.renderer != nil {
		s.pages.SetRenderer(s.renderer, s.rendered)
		s.thumbs.SetRenderer(s.renderer, s.rendered)
	}
	s.ui.search.reset()
	s.ui.inspector.syncText()
	s.ui.viewer.scrollTo(d.current)
	s.win.SetTitle(d.session.Name + " — GoPDF")
}

// stash keeps the active document's view state in its openDoc.
func (s *Studio) stash() {
	if s.active < 0 || s.active >= len(s.docs) {
		return
	}
	d := s.docs[s.active]
	d.session, d.renderer, d.rendered, d.outline = s.doc, s.renderer, s.rendered, s.outline
	d.selected = s.selected
	d.zoom, d.fitWidth, d.current = s.view.zoom, s.view.fitWidth, s.view.current
}

// closeActive closes the active document, asking about unsaved changes,
// and shows its neighbor, or the welcome screen after the last one.
func (s *Studio) closeActive(then func()) {
	if s.doc == nil {
		if then != nil {
			then()
		}
		return
	}
	s.guardUnsaved(func() {
		s.dropRenderer()
		i := s.active
		s.docs = append(s.docs[:i], s.docs[i+1:]...)
		s.active, s.doc, s.selected, s.drawing, s.moving = -1, nil, 0, nil, nil
		if len(s.docs) == 0 {
			s.win.SetTitle("GoPDF")
		} else {
			s.switchTo(min(i, len(s.docs)-1))
		}
		if then != nil {
			then()
		}
	})
}

// closeAll closes every document, asking about each one's unsaved
// changes, then runs then.
func (s *Studio) closeAll(then func()) {
	if len(s.docs) == 0 {
		then()
		return
	}
	s.closeActive(func() { s.closeAll(then) })
}

// cycleDocuments shows the next (dir 1) or previous (dir -1) document.
func (s *Studio) cycleDocuments(dir int) {
	if n := len(s.docs); n > 1 {
		s.switchTo(((s.active+dir)%n + n) % n)
	}
}

// anyDirty reports whether an open document has unsaved changes.
func (s *Studio) anyDirty() bool {
	s.stash()
	for _, d := range s.docs {
		if d.session.Dirty() {
			return true
		}
	}
	return false
}
