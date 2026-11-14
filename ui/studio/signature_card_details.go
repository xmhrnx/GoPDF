package studio

import (
	"fmt"

	"github.com/xmhrnx/giode"

	"gopdf/internal/pdf"
	"gopdf/ui/theme"
)

// signatureCardDetails lays out one signature: its status, then who,
// when, why and where.
func signatureCardDetails(f pdf.SignatureField, st sealStatus) giode.Element {
	color := sealColor(st.level)
	rows := []giode.Element{
		giode.Row(giode.Styles{Gap: 10, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
			theme.Icon(sealIcon(st.level), 22, color),
			giode.Text(st.headline, giode.Styles{FontSize: 14, Color: theme.TextColor, FontWeight: giode.FontWeightSemiBold})),
		giode.Text(st.detail, theme.BodyMute),
	}
	if !f.Signed {
		return detailsCard(color, rows)
	}
	add := func(label, value string) {
		if value != "" {
			rows = append(rows, infoRow(label, value))
		}
	}
	add("Signed by", f.Name)
	add("Email", f.Email)
	add("Organization", f.Org)
	if f.SelfSig {
		add("Certificate", "Self-signed")
	} else {
		add("Certificate", "Issued by "+f.Issuer)
	}
	if !f.When.IsZero() {
		add("Signed on", f.When.Local().Format("Jan 2, 2006 at 15:04"))
	}
	add("Reason", f.Reason)
	add("Location", f.Place)
	if f.Contact != f.Email {
		add("Contact", f.Contact)
	}
	if f.Page >= 0 && f.Rect.Dx() > 0 {
		add("Page", fmt.Sprint(f.Page+1))
	}
	return detailsCard(color, rows)
}

func detailsCard(color string, rows []giode.Element) giode.Element {
	return giode.Stack(giode.Styles{
		Width: -1, Gap: 8, Padding: giode.UniformInset(14), BorderRadius: 12,
		BorderWidth: 1, BorderColor: giode.Mix(color, "theme.border", 0.35),
	}, rows...)
}
