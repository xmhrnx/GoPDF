package studio

import (
	"context"
	"fmt"

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
	snap, sec, path := s.doc.Snapshot(), s.doc.Security, s.doc.Path
	s.runBackground("Saving", func() error {
		data, err := editor.Export(snap, sec)
		if err != nil {
			return err
		}
		return writeFileAtomic(path, data)
	}, func() {
		s.doc.MarkSaved(path)
		s.toasts.Show(giode.ToastSuccess, "Saved "+s.doc.Name)
		if then != nil {
			then()
		}
	})
}

// saveAs asks for a location and writes the document there.
func (s *Studio) saveAs() { s.saveAsThen(nil) }

func (s *Studio) saveAsThen(then func()) {
	if !s.canEdit() {
		return
	}
	snap, sec := s.doc.Snapshot(), s.doc.Security
	var path string
	s.runBackground("Saving", func() error {
		data, err := editor.Export(snap, sec)
		if err != nil {
			return err
		}
		w, err := s.win.SaveFile(context.Background(), s.doc.Name)
		if err != nil {
			return err
		}
		path, err = writeAndClose(w, data)
		return err
	}, func() {
		s.doc.MarkSaved(path)
		s.recents.add(s, path, snap.Doc.PageCount())
		s.win.SetTitle(s.doc.Name + " — GoPDF")
		s.toasts.Show(giode.ToastSuccess, "Saved "+s.doc.Name)
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
