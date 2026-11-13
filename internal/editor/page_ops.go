package editor

import (
	"slices"

	"gopdf/internal/annot"
	"gopdf/internal/pdf"
)

// US Letter, in points: the size of a blank page in an empty document.
const letterW, letterH = 612.0, 792.0

func reload(s Snapshot, data []byte, err error) (Snapshot, error) {
	if err != nil {
		return s, err
	}
	doc, err := pdf.Load(data)
	if err != nil {
		return s, err
	}
	s.Doc = doc
	return s, nil
}

// Rotate turns pages a quarter turn clockwise (cw) or counter-clockwise,
// turning their markup with them.
func Rotate(s Snapshot, pages []int, cw bool) (Snapshot, error) {
	deg := 90
	if !cw {
		deg = -90
	}
	data, err := pdf.Rotate(s.Doc.Data, pages, deg)
	next, err := reload(s, data, err)
	if err != nil {
		return s, err
	}
	next = next.withAnnots()
	for _, p := range pages {
		w, h := s.Doc.Pages[p].Size()
		annot.RotatePage(next.Annots, p, cw, w, h)
	}
	return next, nil
}

// Arrange rebuilds the document from the pages in order (old indexes):
// leaving pages out deletes them, repeating one duplicates it. Markup
// follows its page.
func Arrange(s Snapshot, order []int) (Snapshot, error) {
	data, err := pdf.Arrange(s.Doc.Data, order)
	next, err := reload(s, data, err)
	if err != nil {
		return s, err
	}
	next.Annots = remap(s.Annots, order)
	return next, nil
}

func remap(list []annot.Annotation, order []int) []annot.Annotation {
	var out []annot.Annotation
	used := map[int]bool{}
	for newPage, oldPage := range order {
		for _, a := range list {
			if a.Page != oldPage {
				continue
			}
			c := a.Clone()
			c.Page = newPage
			if used[oldPage] {
				c.ID = NewID()
			}
			out = append(out, c)
		}
		used[oldPage] = true
	}
	return out
}

// Delete removes pages; at least one page must remain.
func Delete(s Snapshot, pages []int) (Snapshot, error) {
	var order []int
	for i := range s.Doc.PageCount() {
		if !slices.Contains(pages, i) {
			order = append(order, i)
		}
	}
	return Arrange(s, order)
}

// Move puts page at index to.
func Move(s Snapshot, page, to int) (Snapshot, error) {
	order := make([]int, 0, s.Doc.PageCount())
	for i := range s.Doc.PageCount() {
		if i != page {
			order = append(order, i)
		}
	}
	order = slices.Insert(order, to, page)
	return Arrange(s, order)
}

// Duplicate inserts a copy of page right after it.
func Duplicate(s Snapshot, page int) (Snapshot, error) {
	order := make([]int, 0, s.Doc.PageCount()+1)
	for i := range s.Doc.PageCount() {
		order = append(order, i)
		if i == page {
			order = append(order, i)
		}
	}
	return Arrange(s, order)
}

// InsertBlank adds an empty page, the size of page after, right after it.
func InsertBlank(s Snapshot, after int) (Snapshot, error) {
	w, h := letterW, letterH
	if after >= 0 && after < s.Doc.PageCount() {
		w, h = s.Doc.Pages[after].Size()
	}
	data, err := pdf.InsertBlank(s.Doc.Data, after, w, h)
	next, err := reload(s, data, err)
	if err != nil {
		return s, err
	}
	next = next.withAnnots()
	for i := range next.Annots {
		if next.Annots[i].Page > after {
			next.Annots[i].Page++
		}
	}
	return next, nil
}

// Merge appends the pages of other documents.
func Merge(s Snapshot, others ...[]byte) (Snapshot, error) {
	data, err := pdf.Merge(s.Doc.Data, others...)
	return reload(s, data, err)
}

// Watermark stamps text on every page.
func Watermark(s Snapshot, wm pdf.Watermark) (Snapshot, error) {
	data, err := pdf.AddWatermark(s.Doc.Data, wm)
	return reload(s, data, err)
}

// PageNumbers stamps page numbers on every page.
func PageNumbers(s Snapshot, format string, pos pdf.PageNumberPosition, size int) (Snapshot, error) {
	data, err := pdf.AddPageNumbers(s.Doc.Data, format, pos, size)
	return reload(s, data, err)
}

// SetMetadata replaces the document information.
func SetMetadata(s Snapshot, m pdf.Metadata) (Snapshot, error) {
	data, err := pdf.SetMetadata(s.Doc.Data, m)
	return reload(s, data, err)
}
