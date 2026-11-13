package studio

import (
	"github.com/xmhrnx/giode"

	"gopdf/ui/components"
	"gopdf/ui/theme"
)

// signatureSection lists the saved signatures to place, and starts new
// ones.
type signatureSection struct {
	s     *Studio
	added *giode.ButtonElement
}

func newSignatureSection(s *Studio) *signatureSection {
	return &signatureSection{
		s:     s,
		added: components.NewSecondaryButton("New Signature…", func() { s.ui.signaturePad.show() }).Icon(theme.Icon("plus", 15, "")),
	}
}

// View returns the section for this frame.
func (sec *signatureSection) View() giode.Element {
	var rows []giode.Element
	if len(sec.s.signatures.list) == 0 {
		rows = append(rows, giode.Text("No saved signatures yet.", theme.Caption))
	}
	active, _ := sec.s.signatures.active(sec.s)
	for _, sig := range sec.s.signatures.list {
		card := giode.Keyed("signature", sig.ID, func() *signatureCard { return newSignatureCard(sec.s, sig) })
		rows = append(rows, card.View(sig.ID == active.ID))
	}
	rows = append(rows, sec.added)
	return components.Section("Signatures", rows...)
}
