package studio

import (
	"fmt"

	"github.com/xmhrnx/giode"

	"gopdf/ui/components"
	"gopdf/ui/theme"
)

// signaturesBanner tells, above the pages, whether the document is signed
// and whether its signatures hold, as PDF readers do.
type signaturesBanner struct {
	s       *Studio
	details *giode.ButtonElement
}

func newSignaturesBanner(s *Studio) *signaturesBanner {
	b := &signaturesBanner{s: s}
	b.details = components.NewSecondaryButton("Signature Details", func() {
		s.ui.sigDetails.showAll()
	})
	return b
}

// View returns the banner, or nothing when the document isn't signed.
func (b *signaturesBanner) View() giode.Element {
	all := b.s.doc.Doc().Signatures
	level, signed := worstStatus(all)
	if signed == 0 {
		return giode.Box(giode.Styles{Display: giode.DisplayNone})
	}
	text := fmt.Sprintf("Signed, %d signature(s): all valid.", signed)
	switch level {
	case sealUnverified:
		text = fmt.Sprintf("Signed, %d signature(s): the document hasn't changed, but a signer's identity isn't verified.", signed)
	case sealBroken:
		text = "At least one signature is invalid: the document changed after it was signed."
	}
	if b.s.doc.SignaturesAtRisk() > 0 {
		text += " Saving your changes will invalidate the signatures."
	}
	color := sealColor(level)
	return giode.Row(giode.Styles{
		Width: -1, Gap: 10, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap, Padding: giode.SymmetricInset(14, 8),
		Background: giode.Mix(color, "theme.surface", 0.12), BorderWidth: 0,
	},
		theme.Icon(sealIcon(level), 18, color),
		giode.Box(giode.Styles{FlexGrow: 1}, giode.Text(text, giode.Styles{FontSize: 12.5, Color: theme.TextColor})),
		b.details,
	)
}
