package studio

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/xmhrnx/giode"

	"gopdf/internal/editor"
	"gopdf/internal/pdf"
)

// pendingOpen is an encrypted file waiting for its password.
type pendingOpen struct {
	data []byte
	path string
	name string
}

// openDialog asks for a PDF and opens it.
func (s *Studio) openDialog() {
	s.guardUnsaved(func() {
		go func() {
			files, err := s.win.PickFiles(context.Background(), giode.PickOptions{Types: []string{"pdf"}})
			if err != nil || len(files) == 0 {
				if err != nil {
					s.win.Post(func() { s.showError(err) })
				}
				return
			}
			f := files[0]
			data, err := io.ReadAll(f)
			f.Close()
			s.win.Post(func() {
				if err != nil {
					s.showError(err)
					return
				}
				s.openBytes(data, f.Path, f.Name, "")
			})
		}()
	})
}

// openPath opens the file at path.
func (s *Studio) openPath(path string) {
	go func() {
		data, err := os.ReadFile(path)
		s.win.Post(func() {
			if err != nil {
				s.showError(err)
				return
			}
			s.openBytes(data, path, filepath.Base(path), "")
		})
	}()
}

// openBytes parses data in the background and shows it, asking for a
// password first when the file is protected.
func (s *Studio) openBytes(data []byte, path, name, password string) {
	var doc *pdf.Document
	s.busy = ""
	s.runBackground("Opening "+name, func() error {
		var err error
		doc, err = pdf.Open(data, password)
		if errors.Is(err, pdf.ErrPassword) {
			s.win.Post(func() { s.ui.password.ask(pendingOpen{data: data, path: path, name: name}, password != "") })
			return nil
		}
		return err
	}, func() {
		if doc != nil {
			session := editor.NewSession(doc, path, name)
			if doc.Encrypted {
				// Saving keeps the file protected with the same password.
				session.Security = editor.Security{UserPassword: password}
			}
			s.show(session)
		}
	})
}

// show makes session the open document.
func (s *Studio) show(session *editor.Session) {
	s.doc = session
	s.selected, s.drawing = 0, nil
	s.view.current = 0
	s.ui.viewer.scrollTo(0)
	s.recents.add(s, session.Path, session.Doc().PageCount())
	s.win.SetTitle(session.Name + " — GoPDF")
}

// closeDocument returns to the welcome screen, asking about unsaved work.
func (s *Studio) closeDocument() {
	s.guardUnsaved(func() {
		s.doc = nil
		s.selected, s.drawing = 0, nil
		s.win.SetTitle("GoPDF")
	})
}
