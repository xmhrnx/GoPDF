package studio

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/styles/properties"

	"gopdf/internal/annot"
	"gopdf/ui/theme"
)

// sealBadges labels the signatures on page that become digital
// signatures when the document is saved.
func sealBadges(s *Studio, page int, z float64) []giode.Element {
	var out []giode.Element
	for _, a := range s.doc.Annots() {
		if a.Page != page || a.Kind != annot.Signature || a.Seal == nil {
			continue
		}
		b := a.Box()
		out = append(out, giode.Box(giode.Styles{
			Position: giode.PositionAbsolute,
			Left:     properties.Left(edge(b.Min.X * z)), Top: properties.Top(edge(b.Max.Y*z + 4)),
		}, giode.Row(giode.Styles{
			Gap: 4, Align: giode.AlignCenter, Padding: giode.SymmetricInset(7, 3), BorderRadius: 99,
			Background: theme.Primary, Wrap: giode.WrapNoWrap,
		},
			theme.Icon("badge-check", 12, theme.OnPrimary),
			giode.Text("Digitally signed on save", giode.Styles{FontSize: 10.5, Color: theme.OnPrimary, FontWeight: giode.FontWeightMedium}),
		)))
	}
	return out
}
