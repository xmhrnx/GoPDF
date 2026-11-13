package studio

import (
	"context"
	"fmt"
	"io"

	"github.com/xmhrnx/giode"

	"gopdf/internal/editor"
	"gopdf/internal/pdf"
)

func (s *Studio) rotatePage(page int, cw bool) {
	s.runEdit("Rotating", "", func(snap editor.Snapshot) (editor.Snapshot, error) {
		return editor.Rotate(snap, []int{page}, cw)
	})
}

func (s *Studio) rotateAll(cw bool) {
	n := s.pageCount()
	s.runEdit("Rotating", "Rotated all pages", func(snap editor.Snapshot) (editor.Snapshot, error) {
		pages := make([]int, n)
		for i := range pages {
			pages[i] = i
		}
		return editor.Rotate(snap, pages, cw)
	})
}

func (s *Studio) deletePage(page int) {
	if s.pageCount() <= 1 {
		s.toasts.Show(giode.ToastWarning, "A document needs at least one page")
		return
	}
	s.runEdit("Deleting page", fmt.Sprintf("Deleted page %d", page+1), func(snap editor.Snapshot) (editor.Snapshot, error) {
		return editor.Delete(snap, []int{page})
	})
}

func (s *Studio) duplicatePage(page int) {
	s.runEdit("Duplicating page", "", func(snap editor.Snapshot) (editor.Snapshot, error) {
		return editor.Duplicate(snap, page)
	})
}

func (s *Studio) movePage(page, to int) {
	if to < 0 || to >= s.pageCount() || to == page {
		return
	}
	s.runEdit("Moving page", "", func(snap editor.Snapshot) (editor.Snapshot, error) {
		return editor.Move(snap, page, to)
	})
	s.view.current = to
}

func (s *Studio) insertBlankPage(after int) {
	s.runEdit("Inserting page", "Inserted a blank page", func(snap editor.Snapshot) (editor.Snapshot, error) {
		return editor.InsertBlank(snap, after)
	})
}

// mergeDialog asks for PDFs and appends their pages.
func (s *Studio) mergeDialog() {
	if !s.canEdit() {
		return
	}
	go func() {
		files, err := s.win.PickFiles(context.Background(), giode.PickOptions{Types: []string{"pdf"}, Multiple: true})
		if err != nil {
			s.win.Post(func() { s.showError(err) })
			return
		}
		var others [][]byte
		for _, f := range files {
			data, err := io.ReadAll(f)
			f.Close()
			if err == nil {
				var doc *pdf.Document
				if doc, err = pdf.Open(data, ""); err == nil {
					data = doc.Data
				}
			}
			if err != nil {
				s.win.Post(func() { s.showError(fmt.Errorf("%s: %w", f.Name, err)) })
				return
			}
			others = append(others, data)
		}
		s.win.Post(func() {
			s.runEdit("Merging", fmt.Sprintf("Added %d file(s)", len(others)), func(snap editor.Snapshot) (editor.Snapshot, error) {
				return editor.Merge(snap, others...)
			})
		})
	}()
}

func (s *Studio) pageCount() int {
	if s.doc == nil {
		return 0
	}
	return s.doc.Doc().PageCount()
}
