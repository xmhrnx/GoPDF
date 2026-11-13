package components

import (
	"github.com/xmhrnx/giode"

	"gopdf/ui/theme"
)

// EmptyHint fills an empty panel with an icon, a title and a hint.
func EmptyHint(icon, title, hint string) giode.Element {
	return giode.Stack(giode.Styles{
		Width: -1, Align: giode.AlignCenter, Gap: 8, Padding: giode.SymmetricInset(20, 36),
	},
		theme.Icon(icon, 28, theme.Faint),
		giode.Text(title, giode.Styles{FontSize: 13, Color: theme.TextColor, FontWeight: giode.FontWeightMedium, TextAlign: giode.TextAlignCenter}),
		giode.Text(hint, giode.Styles{FontSize: 12, Color: theme.Muted, TextAlign: giode.TextAlignCenter}),
	)
}
