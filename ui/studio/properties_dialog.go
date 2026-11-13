package studio

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dialog"
	"github.com/xmhrnx/giode/components/input"

	"gopdf/internal/editor"
	"gopdf/internal/pdf"
	"gopdf/ui/components"
)

// propertiesDialog edits the document information (title, author, ...).
type propertiesDialog struct {
	s                                *Studio
	open                             bool
	title, author, subject, keywords *input.Input
	apply, cancel                    *giode.ButtonElement
	dialog                           *dialog.Dialog
}

func newPropertiesDialog(s *Studio) *propertiesDialog {
	d := &propertiesDialog{s: s}
	d.title = components.NewTextInput("Untitled")
	d.author = components.NewTextInput("Author")
	d.subject = components.NewTextInput("What the document is about")
	d.keywords = components.NewTextInput("Comma, separated, words")
	d.apply = components.NewPrimaryButton("Apply", d.submit)
	d.cancel = components.NewSecondaryButton("Cancel", func() { d.open = false })
	d.dialog = components.NewDialog(&d.open, d.content)
	return d
}

func (d *propertiesDialog) show() {
	m := d.s.doc.Doc().Meta
	d.title.SetText(m.Title)
	d.author.SetText(m.Author)
	d.subject.SetText(m.Subject)
	d.keywords.SetText(m.Keywords)
	d.open = true
}

func (d *propertiesDialog) submit() {
	d.open = false
	m := d.s.doc.Doc().Meta
	m.Title, m.Author, m.Subject, m.Keywords = d.title.Text(), d.author.Text(), d.subject.Text(), d.keywords.Text()
	d.s.runEdit("Updating properties", "Properties updated", func(snap editor.Snapshot) (editor.Snapshot, error) {
		return editor.SetMetadata(snap, m)
	})
}

func (d *propertiesDialog) content() giode.Element {
	m := d.s.doc.Doc().Meta
	body := giode.Stack(giode.Styles{Width: -1, Gap: 12},
		components.Field("Title", d.title),
		components.Field("Author", d.author),
		components.Field("Subject", d.subject),
		components.Field("Keywords", d.keywords),
		readOnlyMeta(m),
	)
	return components.DialogFrame("info", "Document properties", "Information stored inside the PDF file.", body, d.cancel, d.apply)
}

func readOnlyMeta(m pdf.Metadata) giode.Element {
	var rows []giode.Element
	if m.Creator != "" {
		rows = append(rows, infoRow("Created with", m.Creator))
	}
	if m.Producer != "" {
		rows = append(rows, infoRow("Producer", m.Producer))
	}
	return giode.Stack(giode.Styles{Width: -1, Gap: 6}, rows...)
}
