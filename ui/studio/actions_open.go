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
	data    []byte
	path    string
	name    string
	replace bool // in place of the active document
}

// openDialog asks for PDFs and opens each beside the open documents.
func (s *Studio) openDialog() {
	go func() {
		files, err := s.win.PickFiles(context.Background(), giode.PickOptions{Types: []string{"pdf"}, Multiple: true})
		if err != nil {
			s.win.Post(func() { s.showError(err) })
			return
		}
		for _, f := range files {
			data, err := io.ReadAll(f)
			f.Close()
			s.win.Post(func() {
				if err != nil {
					s.showError(err)
					return
				}
				s.openBytes(data, f.Path, f.Name, "", false)
			})
		}
	}()
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
			s.openBytes(data, path, filepath.Base(path), "", false)
		})
	}()
}

// openBytes parses data in the background and shows it in a new tab, or
// in place of the active document with replace, asking for a password
// first when the file is protected.
func (s *Studio) openBytes(data []byte, path, name, password string, replace bool) {
	var doc *pdf.Document
	s.busy = ""
	s.runBackground("Opening "+name, func() error {
		var err error
		doc, err = pdf.Open(data, password)
		if errors.Is(err, pdf.ErrPassword) {
			s.win.Post(func() {
				s.ui.password.ask(pendingOpen{data: data, path: path, name: name, replace: replace}, password != "")
			})
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
			if replace {
				s.replaceDocument(session)
			} else {
				s.addDocument(session)
			}
		}
	})
}

// closeDocument closes the active document, asking about unsaved work.
func (s *Studio) closeDocument() { s.closeActive(nil) }
