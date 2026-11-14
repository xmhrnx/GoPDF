package studio

import (
	"fmt"

	"github.com/xmhrnx/giode"

	"gopdf/internal/annot"
	"gopdf/ui/theme"
)

// pendingSealCard describes a signature that saving will seal: who will
// sign, and why and where.
func pendingSealCard(s *Studio, a annot.Annotation) giode.Element {
	rows := []giode.Element{
		giode.Row(giode.Styles{Gap: 10, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
			theme.Icon("badge-check", 22, theme.Primary),
			giode.Text("Will be signed when you save", giode.Styles{FontSize: 14, Color: theme.TextColor, FontWeight: giode.FontWeightSemiBold})),
	}
	if info, ok := s.ids.info(a.Seal.ID); ok {
		rows = append(rows, infoRow("Signed by", info.Identity.Name))
		if info.Identity.Email != "" {
			rows = append(rows, infoRow("Email", info.Identity.Email))
		}
		if info.Identity.Organization != "" {
			rows = append(rows, infoRow("Organization", info.Identity.Organization))
		}
	} else {
		rows = append(rows, giode.Text("Its digital ID was deleted: choose another before saving.", giode.Styles{FontSize: 12, Color: theme.Danger}))
	}
	if a.Seal.Reason != "" {
		rows = append(rows, infoRow("Reason", a.Seal.Reason))
	}
	if a.Seal.Location != "" {
		rows = append(rows, infoRow("Location", a.Seal.Location))
	}
	rows = append(rows, infoRow("Page", fmt.Sprint(a.Page+1)))
	return detailsCard("theme.primary", rows)
}
