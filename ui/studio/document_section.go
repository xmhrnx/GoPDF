package studio

import (
	"fmt"

	"github.com/xmhrnx/giode"

	"gopdf/ui/components"
	"gopdf/ui/theme"
)

// documentSection summarizes the document and links to its settings.
type documentSection struct {
	s          *Studio
	properties *giode.ButtonElement
	security   *giode.ButtonElement
}

func newDocumentSection(s *Studio) *documentSection {
	return &documentSection{
		s:          s,
		properties: components.NewSecondaryButton("Properties", s.ui.openProperties).Icon(theme.Icon("info", 15, "")),
		security:   components.NewSecondaryButton("Password", s.ui.openSecurity).Icon(theme.Icon("lock", 15, "")),
	}
}

// View returns the section for this frame.
func (d *documentSection) View() giode.Element {
	doc := d.s.doc.Doc()
	w, h := doc.Pages[min(d.s.view.current, doc.PageCount()-1)].Size()
	protection := "None"
	if d.s.doc.Security.Enabled() {
		protection = "Password on save"
	}
	rows := []giode.Element{
		infoRow("File", d.s.doc.Name),
		infoRow("Pages", fmt.Sprint(doc.PageCount())),
		infoRow("Page size", formatPageSize(w, h)),
		infoRow("Size", formatBytes(len(doc.Data))),
		infoRow("Markup", fmt.Sprintf("%d item(s)", len(d.s.doc.Annots()))),
		infoRow("Protection", protection),
	}
	if t := doc.Meta.Title; t != "" {
		rows = append(rows[:1], append([]giode.Element{infoRow("Title", t)}, rows[1:]...)...)
	}
	rows = append(rows, giode.Row(giode.Styles{Gap: 8, Padding: giode.Inset{Top: 6}}, d.properties, d.security))
	return components.Section("Document", rows...)
}
