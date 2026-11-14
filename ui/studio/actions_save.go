package studio

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/xmhrnx/giode"

	"gopdf/internal/editor"
)

// save writes the document back to its file, or asks where to save it.
func (s *Studio) save() {
	s.saveThen(nil)
}

// saveThen saves and runs then once the file is written.
func (s *Studio) saveThen(then func()) {
	if !s.canEdit() {
		return
	}
	if s.doc.Path == "" {
		s.saveAsThen(then)
		return
	}
	s.prepareSave(func() {
		path := s.doc.Path
		s.writeDocument(func(data []byte) (string, error) { return path, writeFileAtomic(path, data) }, then)
	})
}

// saveAs asks for a location and writes the document there.
func (s *Studio) saveAs() { s.saveAsThen(nil) }

func (s *Studio) saveAsThen(then func()) {
	if !s.canEdit() {
		return
	}
	s.prepareSave(func() {
		name := s.doc.Name
		s.writeDocument(func(data []byte) (string, error) {
			w, err := s.win.SaveFile(context.Background(), name)
			if err != nil {
				return "", err
			}
			return writeAndClose(w, data)
		}, then)
	})
}

// prepareSave unlocks the digital IDs the document's signatures need and
// confirms breaking existing signatures, then runs write.
func (s *Studio) prepareSave(write func()) {
	for _, id := range editor.SealIDs(s.doc.Snapshot()) {
		if s.ids.unlocked[id] != nil {
			continue
		}
		info, ok := s.ids.info(id)
		if !ok {
			s.toasts.Show(giode.ToastError, "A signature uses a digital ID that was deleted. Choose another ID for it.")
			return
		}
		s.ui.unlockID.ask(info, func() { s.prepareSave(write) })
		return
	}
	if n := s.doc.SignaturesAtRisk(); n > 0 {
		s.ui.confirm.ask("shield-alert", "Invalidate signatures?",
			fmt.Sprintf("This document has %d digital signature(s). Saving your changes rewrites the file, and the signatures will no longer be valid.", n),
			"Save Anyway", write)
		return
	}
	write()
}

// writeDocument exports the document in the background and hands the
// bytes to put, which writes them and returns the path, if known.
func (s *Studio) writeDocument(put func([]byte) (string, error), then func()) {
	snap, sec := s.doc.Snapshot(), s.doc.Security
	signers := s.ids.unlocked
	sealed := len(editor.SealIDs(snap)) > 0
	var (
		path string
		data []byte
	)
	s.runBackground("Saving", func() error {
		var err error
		if data, err = editor.Export(snap, sec, signers); err != nil {
			return err
		}
		path, err = put(data)
		return err
	}, func() {
		name := s.doc.Name
		s.doc.MarkSaved(path)
		if path != "" {
			s.recents.add(s, path, snap.Doc.PageCount())
		}
		s.win.SetTitle(s.doc.Name + " — GoPDF")
		s.toasts.Show(giode.ToastSuccess, "Saved "+s.doc.Name)
		if sealed {
			// Show the signed file as it now is: sealed, with its
			// signatures checked.
			if path != "" {
				name = filepath.Base(path)
			}
			s.openBytes(data, path, name, sec.UserPassword)
		}
		if then != nil {
			then()
		}
	})
}

// exportFile writes data produced in the background to a file the user
// picks, suggesting name.
func (s *Studio) exportFile(label, name string, produce func() ([]byte, error)) {
	if !s.canEdit() {
		return
	}
	s.runBackground(label, func() error {
		data, err := produce()
		if err != nil {
			return err
		}
		w, err := s.win.SaveFile(context.Background(), name)
		if err != nil {
			return err
		}
		_, err = writeAndClose(w, data)
		return err
	}, func() { s.toasts.Show(giode.ToastSuccess, fmt.Sprintf("Exported %s", name)) })
}
