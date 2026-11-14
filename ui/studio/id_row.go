package studio

import (
	"github.com/xmhrnx/giode"

	"gopdf/internal/digitalid"
	"gopdf/ui/components"
	"gopdf/ui/theme"
)

// idRow is one stored digital ID.
type idRow struct {
	info   digitalid.Info
	remove *components.ToolButton
}

func newIDRow(s *Studio, info digitalid.Info) *idRow {
	r := &idRow{info: info}
	r.remove = components.NewToolButton("trash-2", "Delete this digital ID", func() {
		s.ui.confirm.ask("trash-2", "Delete digital ID?",
			"Signatures already made with it stay valid, but you won't be able to sign as "+info.Identity.Name+" with it again.",
			"Delete", func() { s.ids.remove(s, info.Fingerprint) })
	}).Compact()
	return r
}

// View returns the row for this frame.
func (r *idRow) View() giode.Element {
	who := r.info.Identity.Name
	if r.info.Identity.Email != "" {
		who += " · " + r.info.Identity.Email
	}
	issuer := "Self-signed"
	if !r.info.SelfSigned {
		issuer = "Issued by " + r.info.Issuer
	}
	return giode.Row(giode.Styles{
		Width: -1, Gap: 10, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap, Padding: giode.SymmetricInset(12, 8),
		BorderRadius: 10, BorderWidth: 1, BorderColor: theme.Border,
	},
		theme.Icon("badge-check", 18, theme.Primary),
		giode.Stack(giode.Styles{Gap: 2, FlexGrow: 1},
			giode.Text(who, giode.Styles{FontSize: 13, Color: theme.TextColor, FontWeight: giode.FontWeightMedium, TextOverflow: giode.TextOverflowEllipsis}),
			giode.Text(issuer+" · valid until "+r.info.NotAfter.Format("Jan 2006"), theme.Caption),
		),
		r.remove.View(),
	)
}
