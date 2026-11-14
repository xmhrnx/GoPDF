package studio

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dialog"
	"github.com/xmhrnx/giode/components/scroll"

	"gopdf/internal/annot"
	"gopdf/internal/pdf"
	"gopdf/ui/components"
	"gopdf/ui/theme"
)

// detailsMaxHeight caps the list of signatures before it scrolls, in dp.
const detailsMaxHeight = 440

// signatureDetailsDialog shows who signed and whether each signature
// holds.
type signatureDetailsDialog struct {
	s       *Studio
	open    bool
	fields  []pdf.SignatureField
	pending []annot.Annotation // signatures sealed on the next save
	list    *scroll.Scroll
	close   *giode.ButtonElement
	dialog  *dialog.Dialog
}

func newSignatureDetailsDialog(s *Studio) *signatureDetailsDialog {
	d := &signatureDetailsDialog{s: s}
	d.list = giode.Scroll(d.cards).FitContent(true)
	d.close = components.NewPrimaryButton("Close", func() { d.open = false })
	d.dialog = components.NewDialog(&d.open, d.content)
	return d
}

// show opens the dialog on fields.
func (d *signatureDetailsDialog) show(fields []pdf.SignatureField) {
	d.fields, d.pending, d.open = fields, nil, true
	d.list.ScrollTo(0)
}

// showAll opens the dialog on every signature of the document: those in
// the file and those that saving will add.
func (d *signatureDetailsDialog) showAll() {
	d.show(d.s.doc.Doc().Signatures)
	for _, a := range d.s.doc.Annots() {
		if a.Kind == annot.Signature && a.Seal != nil {
			d.pending = append(d.pending, a)
		}
	}
}

func (d *signatureDetailsDialog) content() giode.Element {
	return components.DialogFrame("badge-check", "Signatures",
		"Checked on this computer, without contacting any server.",
		giode.Box(giode.Styles{Width: -1, MaxHeight: detailsMaxHeight}, d.list), d.close)
}

func (d *signatureDetailsDialog) cards() giode.Element {
	if len(d.fields) == 0 && len(d.pending) == 0 {
		return giode.Text("This document has no signatures yet. Pick the Signature tool (S) and click where you want to sign.", theme.BodyMute)
	}
	all := d.s.doc.Doc().Signatures
	cards := make([]giode.Element, 0, len(d.fields)+len(d.pending))
	for _, f := range d.fields {
		cards = append(cards, signatureCardDetails(f, statusOf(f, all)))
	}
	for _, a := range d.pending {
		cards = append(cards, pendingSealCard(d.s, a))
	}
	return giode.Stack(giode.Styles{Width: -1, Gap: 12, Padding: giode.Inset{Right: 8}}, cards...)
}
