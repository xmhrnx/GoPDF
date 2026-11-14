package studio

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dialog"
	"github.com/xmhrnx/giode/components/scroll"

	"gopdf/internal/pdf"
	"gopdf/ui/components"
)

// detailsMaxHeight caps the list of signatures before it scrolls, in dp.
const detailsMaxHeight = 440

// signatureDetailsDialog shows who signed and whether each signature
// holds.
type signatureDetailsDialog struct {
	s      *Studio
	open   bool
	fields []pdf.SignatureField
	list   *scroll.Scroll
	close  *giode.ButtonElement
	dialog *dialog.Dialog
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
	d.fields, d.open = fields, true
	d.list.ScrollTo(0)
}

func (d *signatureDetailsDialog) content() giode.Element {
	return components.DialogFrame("badge-check", "Signatures",
		"Checked on this computer, without contacting any server.",
		giode.Box(giode.Styles{Width: -1, MaxHeight: detailsMaxHeight}, d.list), d.close)
}

func (d *signatureDetailsDialog) cards() giode.Element {
	all := d.s.doc.Doc().Signatures
	cards := make([]giode.Element, len(d.fields))
	for i, f := range d.fields {
		cards[i] = signatureCardDetails(f, statusOf(f, all))
	}
	return giode.Stack(giode.Styles{Width: -1, Gap: 12, Padding: giode.Inset{Right: 8}}, cards...)
}
