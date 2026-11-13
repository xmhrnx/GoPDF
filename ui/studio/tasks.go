package studio

import (
	"github.com/xmhrnx/giode"

	"gopdf/internal/editor"
)

// runEdit runs a page operation on a background goroutine and commits
// its result as one undoable change. label shows in the status bar.
func (s *Studio) runEdit(label, done string, op func(editor.Snapshot) (editor.Snapshot, error)) {
	if !s.canEdit() {
		return
	}
	snap := s.doc.Snapshot()
	session := s.doc
	s.busy = label
	go func() {
		next, err := op(snap)
		s.win.Post(func() {
			s.busy = ""
			if err != nil {
				s.toasts.Show(giode.ToastError, label+" failed: "+err.Error())
				return
			}
			if s.doc != session {
				return // the document was closed meanwhile
			}
			s.doc.Commit(next)
			s.selected = 0
			if done != "" {
				s.toasts.Show(giode.ToastSuccess, done)
			}
		})
	}()
}

// runBackground runs work off the UI goroutine with label in the status
// bar, then reports its error, or calls then with no error, on the UI
// goroutine.
func (s *Studio) runBackground(label string, work func() error, then func()) {
	if !s.idle() {
		return
	}
	s.busy = label
	go func() {
		err := work()
		s.win.Post(func() {
			s.busy = ""
			if err != nil {
				s.showError(err)
				return
			}
			if then != nil {
				then()
			}
		})
	}()
}

func (s *Studio) showError(err error) {
	if isCanceled(err) {
		return
	}
	s.toasts.Show(giode.ToastError, err.Error())
}
